package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

const careAboutColumnsOrder = true // TODO: make this optional

func main() {
	fmt.Print("\n")
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalf("No program arguments, please use: %s -t_1 <db 1 type> -db_1 <db 1 name> -u_1 <user 1> [--password_1 <password 1>] -h_1 <host 1> [-p_1 <port 1>] -t_2 <db 2 type> -db_2 <db 2 name> -u_2 <user 2> [--password_2 <password 2>] -h_2 <host 2> [-p_2 <port 2>]", os.Args[0])
	}

	connType1 := parseConnType(&args, "1")
	db1, err := newDbFrom(connType1, "1")
	if err != nil {
		log.Fatalln(err)
	}
	defer db1.close()

	connType2 := parseConnType(&args, "2")
	db2, err := newDbFrom(connType2, "2")
	if err != nil {
		log.Fatalln(err)
	}
	defer db2.close()

	if connType1.getConnString() == connType2.getConnString() {
		log.Fatalln("The databases are the same")
	}

	fmt.Print("\n")
	tablesInCommon, err := compareTables(db1, db2)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print("\n")
	err = compareTablesContent(db1, db2, tablesInCommon)
	if err != nil {
		log.Fatalln(err)
	}
}

// compareTables prints the differences between the tables of the given databases,
// It returns the tables the databases have in common or any error
func compareTables(first, second *anyDb) ([]string, error) {
	tables1, err := first.listTables()
	if err != nil {
		return nil, fmt.Errorf("failed to list tables of database n°1: %v", err)
	}
	tables2, err := second.listTables()
	if err != nil {
		return nil, fmt.Errorf("failed to list tables of database n°2: %v", err)
	}

	var both []string
	var notIn1 []string
	var notIn2 []string

firstLoop:
	for _, table1 := range tables1 {
		for _, table2 := range tables2 {
			if table1 == table2 {
				both = append(both, table1)
				continue firstLoop
			}
		}
		notIn2 = append(notIn2, table1)
	}

secondLoop:
	for _, table2 := range tables2 {
		for _, table1 := range tables1 {
			if table2 == table1 {
				continue secondLoop
			}
		}
		notIn1 = append(notIn1, table2)
	}

	fmt.Println("Comparison of the tables of the two databases:")
	spacer()
	if len(notIn1) == 0 && len(notIn2) == 0 {
		fmt.Println(colorForSame("Both databases have the same tables."))
	} else {
		if len(both) == 0 {
			fmt.Println(colorForMissing("They are no tables in common between the two databases."))
		} else {
			fmt.Print("Tables in common between the two databases:\n\t", colorForCommon(strings.Join(both, ", ")), "\n")
		}
		if len(notIn1) == 0 {
			fmt.Println(colorForPartiallySame("The first database contains all the tables of the second database."))
		} else {
			fmt.Print("Tables missing from the first database:\n\t", colorForMissing(strings.Join(notIn1, ", ")), "\n")
		}
		if len(notIn2) == 0 {
			fmt.Println(colorForPartiallySame("The second database contains all the tables of the first database."))
		} else {
			fmt.Print("Tables missing from the second database:\n\t", colorForMissing(strings.Join(notIn2, ", ")), "\n")
		}
	}
	spacer()

	return both, nil
}

// compareTablesContent compares the columns of the given tablesInCommon from db1 and db2
func compareTablesContent(db1, db2 *anyDb, tablesInCommon []string) error {
	for _, table := range tablesInCommon {
		columns1, err := db1.listColumns(table)
		if err != nil {
			return fmt.Errorf("failed to list columns of the table %q in the first database: %v", table, err)
		}
		columns2, err := db2.listColumns(table)
		if err != nil {
			return fmt.Errorf("failed to list columns of the table %q in the second database: %v", table, err)
		}

		if careAboutColumnsOrder {
			sort.Slice(columns1, lessColFunc(columns1))
			sort.Slice(columns2, lessColFunc(columns2))
		}

		fmt.Printf("\nComparison of the columns of the table %q:\n", table)
		spacer()
		compareColumns(columns1, columns2)
		spacer()
	}

	return nil
}

// compareColumns compares the two given set of columns and prints their differences
func compareColumns(columns1, columns2 []*sql.ColumnType) {
	var both []string
	var notSameType []string
	var notSameLength []string
	var notSameNullability []string
	var notIn1 []string
	var notIn2 []string

	for i, col1 := range columns1 {
		if i < len(columns2) {
			col2 := columns2[i]
			if col1.Name() == col2.Name() {
				// Same type check
				if col1.DatabaseTypeName() != col2.DatabaseTypeName() {
					notSameType = append(notSameType, col1.Name())
					continue
				}
				// Same length check
				length1, ok1 := col1.Length()
				length2, ok2 := col2.Length()
				if ok1 != ok2 || length1 != length2 {
					notSameLength = append(notSameLength, col1.Name())
					continue
				}
				// Same nullability check
				null1, ok1 := col1.Nullable()
				null2, ok2 := col2.Nullable()
				if ok1 != ok2 || null1 != null2 {
					notSameNullability = append(notSameNullability, col1.Name())
					continue
				}
				both = append(both, col1.Name())
				continue
			}
		}
		notIn2 = append(notIn2, col1.Name())
	}

	for i, col2 := range columns2 {
		if i < len(columns1) {
			if col2.Name() == columns1[i].Name() {
				continue
			}
		}
		notIn1 = append(notIn1, col2.Name())
	}

	if len(notIn1) == 0 && len(notIn2) == 0 {
		fmt.Println(colorForSame("Both tables have the same columns."))
	} else {
		if len(both) == 0 {
			fmt.Println(colorForMissing("They are no columns in common between the two tables."))
		} else {
			fmt.Print("Columns in common between the two tables:\n\t", colorForCommon(strings.Join(both, ", ")), "\n")
		}
		if len(notSameType) != 0 {
			fmt.Print("Columns having the same name but not the same type:\n\t", colorForDifferent(strings.Join(notSameType, ", ")), "\n")
		}
		if len(notSameLength) != 0 {
			fmt.Print("Columns having the same name but not the same length:\n\t", colorForDifferent(strings.Join(notSameLength, ", ")), "\n")
		}
		if len(notSameNullability) != 0 {
			fmt.Print("Columns having the same name but not the same nullability:\n\t", colorForDifferent(strings.Join(notSameNullability, ", ")), "\n")
		}
		if len(notIn1) == 0 {
			fmt.Println(colorForPartiallySame("The first table contains all the columns of the second table."))
		} else {
			fmt.Print("Columns missing from the first table:\n\t", colorForMissing(strings.Join(notIn1, ", ")), "\n")
		}
		if len(notIn2) == 0 {
			fmt.Println(colorForPartiallySame("The second table contains all the columns of the first table."))
		} else {
			fmt.Print("Columns missing from the second table:\n\t", colorForMissing(strings.Join(notIn2, ", ")), "\n")
		}
	}
}

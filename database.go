package main

import (
	"database/sql"
	"fmt"
	"log"
)

func newDbFrom(connType connectionType, id string) (*anyDb, error) {
	fmt.Printf("Database %s: connecting with %s\n", id, connType.getConnString())
	conn, err := sql.Open(connType.getDriver(), connType.getConnString())
	if err != nil {
		return nil, fmt.Errorf("failed to open the connection: %v", err)
	}
	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}
	return &anyDb{
		connType: connType,
		db:       conn,
	}, nil
}

type anyDb struct {
	connType connectionType
	db       *sql.DB
}

// connect initializes the connection to the database and returns any error
func (anyDb *anyDb) connect() error {
	fmt.Println("Connecting with", anyDb.connType.getConnString())
	db, err := sql.Open(anyDb.connType.getDriver(), anyDb.connType.getConnString())
	if err != nil {
		return fmt.Errorf("failed to open the connection: %v", err)
	}
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}
	anyDb.db = db
	return nil
}

// close closes the connection to the database and prints any error
func (anyDb *anyDb) close() {
	err := anyDb.db.Close()
	if err != nil {
		log.Printf("[%s] Failed to close connection: %v", anyDb.connType.getDriver(), err)
	}
}

// listTables returns all the tables contained in the database
func (anyDb *anyDb) listTables() ([]string, error) {
	rows, err := anyDb.db.Query(anyDb.connType.getListTablesScript())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

// listColumns returns all the columns contained in the given table
func (anyDb *anyDb) listColumns(table string) ([]*sql.ColumnType, error) {
	rows, err := anyDb.db.Query(anyDb.connType.getListColumnsScript(table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	return cols, nil
}

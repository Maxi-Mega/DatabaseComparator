package main

import (
	"database/sql"
	"fmt"
)

// lessColFunc returns a function that compare alphabetically the columns at the indexes i and j from the slice cols
func lessColFunc(cols []*sql.ColumnType) func(i, j int) bool {
	return func(i, j int) bool {
		col1Name := cols[i].Name()
		col2Name := cols[j].Name()
		for c := 0; c < len(col1Name) && c < len(col2Name); c++ {
			if col1Name[c] == col2Name[c] {
				continue
			}
			return col1Name[c] < col2Name[c]
		}
		return false
	}
}

// spacer prints 20 dashes each separated by a space
func spacer() {
	fmt.Println("- - - - - - - - - - - - - - - - - - - -")
}

/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"
)

var testlistCmd = &cobra.Command{
	Use:   "list",
	Short: "get all test list",
	Long:  `get all test`,
	Run: func(cmd *cobra.Command, args []string) {
		testdb, err := sql.Open("sqlite3", "./tests.db")
		if err != nil {
			fmt.Println("Unable to open Sqlite connection:", err)
		}
		defer testdb.Close()
		rows, err := testdb.Query("SELECT * FROM list")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var payload string
			if err := rows.Scan(&id, &payload); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(id, payload)
		}
		if err := rows.Err(); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	testCmd.AddCommand(testlistCmd)
}

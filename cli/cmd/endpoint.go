/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config",
	Long:  `config lets you set the endpoint`,
	Run: func(cmd *cobra.Command, args []string) {
		// in future we can implement this like
		// we can get the endpoint value and store in db for persistance
		fmt.Println("The endpoint is set as: ", endPoint)
		db, err := sql.Open("sqlite3", "./value.db")
		if err != nil {
			fmt.Println("Unable to open Sqlite connection:", err)
		}
		defer db.Close()

		sqlStmt := `
	CREATE TABLE IF NOT EXISTS endpoint (name TEXT);
	DELETE FROM endpoint;
	`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			fmt.Println(err)
			return
		}
		tx, err := db.Begin()
		if err != nil {
			fmt.Println(err)
			return
		}
		stmt, err := tx.Prepare("INSERT INTO endpoint(name) VALUES(?)")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(endPoint)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = tx.Commit()
		if err != nil {
			fmt.Println(err)
			return
		}
		rows, err := db.Query("SELECT name FROM endpoint")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var endpoint string
			err = rows.Scan(&endpoint)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("scanned endpoint value: ", endpoint)
		}
		err = rows.Err()
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {

	configCmd.PersistentFlags().StringVarP(&endPoint, "endpoint", "e", "", "--endpoint <endpoint of qt server>")
	rootCmd.AddCommand(configCmd)
}

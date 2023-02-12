/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var testRunCmd = &cobra.Command{
	Use:   "run",
	Short: "run a test",
	Long:  "run a test using a definition or an id",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", "./value.db")
		if err != nil {
			fmt.Println("Unable to open Sqlite connection:", err)
		}
		defer db.Close()
		rows, err := db.Query("SELECT name FROM endpoint")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()
		var endpoint string
		for rows.Next() {
			err = rows.Scan(&endpoint)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		err = rows.Err()
		if err != nil {
			fmt.Println(err)
			return
		}

		// endpoint := os.Getenv("ENDPOINT")
		fmt.Println("run called")
		data, err := ParseYaml(runTestFileDefinition)
		if err != nil {
			log.Println("Unable to parse the yaml file")
		}
		json_data, err := json.Marshal(data)

		if err != nil {
			log.Fatal(err)
		}
		path := fmt.Sprintf("%s/test/%v/run", endpoint, testID)
		fmt.Println(path)
		resp, err := http.Post(path, "application/json",
			bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(body))
		testdb, err := sql.Open("sqlite3", "./tests.db")
		if err != nil {
			fmt.Println("Unable to open Sqlite connection:", err)
		}
		defer testdb.Close()
		testquery := `
CREATE TABLE IF NOT EXISTS list (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	payload TEXT
);
`
		_, err = testdb.Exec(testquery)
		if err != nil {
			fmt.Println("error while create a list",err)
			return
		}

		txnquery := `
		INSERT INTO list (payload)
		VALUES (?);
		`
		_, err = testdb.Exec(txnquery, string(body))
		if err != nil {
			fmt.Println("error while inserting the test in db",err)
			return
		}
	},
}

func init() {
	testRunCmd.PersistentFlags().StringVarP(&runTestFileDefinition, "definition", "d", "", "--definition <definition-file.yml>")
	testRunCmd.PersistentFlags().StringVarP(&testID, "test-id", "i", "", "--testid <test-id>")
	testCmd.AddCommand(testRunCmd)
}

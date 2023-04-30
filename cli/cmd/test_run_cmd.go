/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var testRunCmd = &cobra.Command{
	Use:   "run",
	Short: "run a test",
	Long:  "run a test using a definition or an id",
	Run: func(cmd *cobra.Command, args []string) {
		if Config.QualityTraceUrl == "" {
			log.Println(`Please set the Quality Trace url using this command:
			qt config --set-server <Quality-Trace-url>
			OR
			create $HOME/config/config.yaml and provide the details 
			for example: 
			CH_CONN: http://localhost:9000?username=admin&password=admin
			QT_CONN: http://localhost:8080 `)
			return
		}

		fmt.Println("run called")
		data, err := ParseYaml(runTestFileDefinition)
		if err != nil {
			log.Println("Unable to parse the yaml file")
		}
		json_data, err := json.Marshal(data)

		if err != nil {
			log.Fatal(err)
		}
		path := fmt.Sprintf("%s/test/%v/run", Config.QualityTraceUrl, testID)
		fmt.Println(path)
		resp, err := http.Post(path, "application/json",
			bytes.NewBuffer(json_data))
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(body))

		// this connection is used to store the test repponses
		// this responses will be used with the qt test list
		// it will fetch the test results and show it in cli
		// 		testdb, err := sql.Open("sqlite3", "./tests.db")
		// 		if err != nil {
		// 			fmt.Println("Unable to open Sqlite connection:", err)
		// 		}
		// 		defer testdb.Close()
		// 		testquery := `
		// CREATE TABLE IF NOT EXISTS list (
		// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
		// 	payload TEXT
		// );
		// `
		// 		_, err = testdb.Exec(testquery)
		// 		if err != nil {
		// 			fmt.Println("error while create a list", err)
		// 			return
		// 		}

		// 		txnquery := `
		// 		INSERT INTO list (payload)
		// 		VALUES (?);
		// 		`
		// 		_, err = testdb.Exec(txnquery, string(body))
		// 		if err != nil {
		// 			fmt.Println("error while inserting the test in db", err)
		// 			return
		// 		}
	},
}

func init() {
	testRunCmd.PersistentFlags().StringVarP(&runTestFileDefinition, "definition", "d", "", "--definition <definition-file.yml>")
	testRunCmd.PersistentFlags().StringVarP(&testID, "test-id", "i", "", "--testid <test-id>")
	testCmd.AddCommand(testRunCmd)
}

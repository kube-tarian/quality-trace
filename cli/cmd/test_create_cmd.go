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

func logFatal(args ...interface{}) {
	log.Fatal(args)
}

var testCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a test",
	Long:  "create a test using a definition",
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
		// if endpoint == "" {
		// 	fmt.Println(`please set the endpoint
		// 	qt config --endpoint <your end point>`)
		// 	return
		// }
		fmt.Println("create called")
		data, err := ParseYaml(runTestFileDefinition)
		if err != nil {
			log.Println("Unable to parse the yaml file")
		}
		json_data, err := json.Marshal(data)

		if err != nil {
			log.Fatal(err)
		}
		path := fmt.Sprintf("%s/test/", endpoint)
		resp, err := http.Post(path, "application/json",
			bytes.NewBuffer(json_data))

		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("unable to read the response body: ", err)
		}
		fmt.Println(string(body))
	},
}

func init() {
	testCreateCmd.PersistentFlags().StringVarP(&runTestFileDefinition, "definition", "d", "", "--definition <definition-file.yml>")
	testCmd.AddCommand(testCreateCmd)
}

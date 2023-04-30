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

var testCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a test",
	Long:  "create a test using a definition",
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
		fmt.Println("create called")
		// runTestFileDefinition is populated from the persistent flags
		data, err := ParseYaml(runTestFileDefinition)
		if err != nil {
			log.Println("Unable to parse the yaml file")
		}
		json_data, err := json.Marshal(data)

		if err != nil {
			log.Fatal(err)
		}
		path := fmt.Sprintf("%s/test/", Config.QualityTraceUrl)
		resp, err := http.Post(path, "application/json",
			bytes.NewBuffer(json_data))

		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
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

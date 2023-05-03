package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/spf13/cobra"
)

var generateAssertionCmd = &cobra.Command{
	Use:   "assertion",
	Short: "generates assertion",
	Long:  "generates assertions based on traces generated on previous run",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("ASSERTION GENERATION STARTED...")
		var result []Repo
		reponame := args[0]
		if Config.ClickHouseUrl == "" {
			log.Println(`Please set the Clickhouse url using this command:
			qt config --set-clickhouse <Clickhouse-url>
			OR
			create $HOME/config/config.yaml and provide the details 
			for example: 
			ch_conn: http://localhost:9000?username=admin&password=admin
			qt_conn: http://localhost:8080 `)
			return
		}
		// parsing the clickhouse url
		dsnURL, err := url.Parse(Config.ClickHouseUrl)
		if err != nil {
			logger.Println(err)
		}

		options := &clickhouse.Options{
			Addr: []string{dsnURL.Host},
		}
		if dsnURL.Query().Get("username") == "" || dsnURL.Query().Get("password") == "" {
			logger.Println("url query has credentials missing")
		}
		if dsnURL.Query().Get("username") != "" {
			auth := clickhouse.Auth{
				Database: "signoz_traces",
				Username: dsnURL.Query().Get("username"),
				Password: dsnURL.Query().Get("password"),
			}
			options.Auth = auth
		}

		// creating a clickhouse connection
		db, err := clickhouse.Open(options)
		if err != nil {
			logger.Println(err)
			return
		}
		// closing the clickhouse connection at the end
		defer db.Close()
		// checking the clickhouse connection with a ping
		err = db.Ping(context.Background())
		if err != nil {
			logger.Println(err)
		}
		// we query the database and print the details
		query := "SELECT * FROM signoz_traces.repo WHERE name = ?"
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = db.Select(ctx, &result, query, reponame)
		if err != nil {
			logger.Println("Unable to query the database", err)
		}

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

		jsonData, err := json.Marshal(result[0])
		if err != nil {
			logger.Println("unable to marshal the data:", err)
		}

		path := fmt.Sprintf("%s/dryRun/", Config.QualityTraceUrl)
		resp, err := http.Post(path, "application/json",
			bytes.NewBuffer(jsonData))

		if err != nil {
			logger.Println(err)
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Println("unable to read the response body: ", err)
			return
		}
		strBody := string(body)
		statusStr := "Fetching Data..."
		lastIdx := strings.LastIndex(strBody, statusStr)
		jsonBody := strBody[lastIdx+len(statusStr)+2 : len(strBody)-2]
		data := AssertionData{}
		json.Unmarshal([]byte(jsonBody), &data)
		for _, event := range data.Events {
			assertionFields := map[int]string{}
			assertionFieldValueMap := map[string]string{}
			if event[4] == "/books" {
				// assertion fields
				for idx, v := range event[7].([]interface{}) {
					assertionFields[idx] = v.(string)
				}
				// assertion values
				for idx, v := range event[8].([]interface{}) {
					if field, ok := assertionFields[idx]; ok {
						assertionFieldValueMap[field] = v.(string)
					}
				}
				fmt.Printf("\n assertionFieldValueMap %v", assertionFieldValueMap)
				GenerateYaml(assertionFieldValueMap)
			}
		}
	},
}

type AssertionData struct {
	Columns []string        `json:"columns"`
	Events  [][]interface{} `json:"events"`
}

func init() {
	repoCmd.AddCommand(generateAssertionCmd)
}

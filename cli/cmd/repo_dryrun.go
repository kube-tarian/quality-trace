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
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/spf13/cobra"
)

var repoDryRunCmd = &cobra.Command{
	Use:   "dryrun",
	Short: "runs a dryrun",
	Long:  "runs a dryrun using the check.yaml",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("DRYRUN STARTED...")
		reponame := args[0]
		var result []Repo
		if Config.ClickHouseUrl == "" {
			log.Println(`Please set the Clickhouse url using this command:
			qt config --set-clickhouse <Clickhouse-url>
			OR
			create $HOME/config/config.yaml and provide the details 
			for example: 
			CH_CONN: http://localhost:9000?username=admin&password=admin
			QT_CONN: http://localhost:8080 `)
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
			logger.Println("unable to query the database", err)
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
		json_data, err := json.Marshal(result[0])
		if err != nil {
			logger.Println("unable to marshal the data:", err)
		}
		path := fmt.Sprintf("%s/clonecheck/", Config.QualityTraceUrl)
		resp, err := http.Post(path, "application/json",
			bytes.NewBuffer(json_data))

		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Println("unable to read the response body: ", err)
		}
		fmt.Println(string(body))
	},
}

func init() {
	repoCmd.AddCommand(repoDryRunCmd)
}

package cmd

import (
	"context"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"log"
	"net/url"
	"os"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/spf13/cobra"
)

var repolistCmd = &cobra.Command{
	Use:   "list",
	Short: "list a repo",
	Long:  "gives the list of repo details",
	Run: func(cmd *cobra.Command, args []string) {
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
		dsnURL, err := url.Parse(Config.ClickHouseUrl)
		if err != nil {
			log.Println(err)
		}
		options := &clickhouse.Options{
			Addr: []string{dsnURL.Host},
		}
		if dsnURL.Query().Get("username") == "" || dsnURL.Query().Get("password") == "" {
			log.Println("url query has credentials missing")
		}
		if dsnURL.Query().Get("username") != "" {
			auth := clickhouse.Auth{
				Database: "signoz_traces",
				Username: dsnURL.Query().Get("username"),
				Password: dsnURL.Query().Get("password"),
			}
			options.Auth = auth
		}
		db, err := clickhouse.Open(options)
		if err != nil {
			log.Println(err)
			return
		}
		defer db.Close()
		err = db.Ping(context.Background())
		if err != nil {
			log.Println(err)
		}
		query := "SELECT * FROM signoz_traces.repo"
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = db.Select(ctx, &result, query)
		if err != nil {
			fmt.Errorf("error in processing sql query: %v", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"RepoName", "RepoUrl", "RepoAuth", "GitToken", "CreatedAt", "UpdatedAt"})

		for _, val := range result {
			t.AppendRow(table.Row{val.Name, val.URL, val.AuthType, val.Token, val.CreatedAt, val.UpdatedAt})
		}
		t.SetTitle("REPOLIST")
		t.Render()

	},
}

func init() {
	repoCmd.AddCommand(repolistCmd)
}

package cmd

import (
	"context"
	"fmt"
	"log"
	"net/url"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/spf13/cobra"
)

var dropRepoCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drops the repository table",
	Long:  "Drops the repository table from signoz clickhouse",
	Run: func(cmd *cobra.Command, args []string) {
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
		//dsn := "http://localhost:9000?username=admin&password=admin"
		ctx := context.Background()
		dsnURL, err := url.Parse(Config.ClickHouseUrl)
		if err != nil {
			logger.Println(err)
		}
		fmt.Println("URL:", dsnURL)
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
		db, err := clickhouse.Open(options)
		if err != nil {
			logger.Println(err)
		}
		err = db.Ping(context.Background())
		if err != nil {
			logger.Println(err)
		}
		err = db.Exec(ctx, `DROP TABLE IF EXISTS repo`)

		if err != nil {
			logger.Println(err)
		}
	},
}

func init() {
	repoCmd.AddCommand(dropRepoCmd)
}

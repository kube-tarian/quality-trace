package cmd

import (
	"context"
	"log"
	"net/url"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/spf13/cobra"
)

var repoUpdateCmd = &cobra.Command{
	Use:   "update [reponame]",
	Short: "Update a repository",
	Long:  "Update a repository.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var repoUrl string
		var repoAuthType string
		var gitToken string
		name := args[0]
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
		err = db.Ping(context.Background())
		if err != nil {
			log.Println(err)
		}
		if name == "" {
			log.Println("Please provide the repo name : qt update <reponame>")
			return
		}
		if repoUrl != "" || repoAuthType != "" || repoToken != "" {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			updated_at := time.Now()
			var err error
			var params []interface{}
			query := "ALTER TABLE signoz_traces.repo UPDATE updated_at = ?"
			params = append(params, updated_at)
			if repoUrl != "" {
				query += ", url = ?"
				params = append(params, repoUrl)
			}
			if repoAuthType != "" {
				query += ", auth_type = ?"
				params = append(params, repoAuthType)
			}

			if gitToken != "" {
				query += ", token = ?"
				params = append(params, gitToken)
			}
			query += " WHERE name = ?"
			params = append(params, name)
			err = db.Exec(ctx, query, params...)
			if err != nil {
				logger.Printf("error updating table: %v", err)
				return
			}
		}

		// err = db.Exec(ctx, `
		// ALTER TABLE signoz_traces.repo UPDATE url = ?, auth_type = ?, token = ?, updated_at =  ? WHERE name = ?`, repoUrl, repoAuthType, gitToken, updated_at, name)
		// if err != nil {
		// 	log.Printf("error updating table: %v", err)
		// 	return
		// }
		log.Println("successfully updated the details of repo in signoz_traces table")
	},
}

func init() {
	repoUpdateCmd.PersistentFlags().StringVarP(&repoUrl, "url", "u", "", "url of repo")
	repoUpdateCmd.PersistentFlags().StringVarP(&repoAuthType, "auth-type", "a", "", "authentication type of repo")
	repoUpdateCmd.PersistentFlags().StringVarP(&repoToken, "token", "t", "", "token of repo")

	repoCmd.AddCommand(repoUpdateCmd)
}

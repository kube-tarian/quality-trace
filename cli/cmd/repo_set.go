package cmd

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/spf13/cobra"
)

var (
	repoName, repoUrl, repoAuthType, repoToken string
)

type Repo struct {
	Name      string    `ch:"name"`
	URL       string    `ch:"url"`
	AuthType  string    `ch:"auth_type"`
	Token     string    `ch:"token"`
	CreatedAt time.Time `ch:"created_at"`
	UpdatedAt time.Time `ch:"updated_at"`
}

var repoSetCmd = &cobra.Command{
	Use:   "set",
	Short: "sets the repo details",
	Long:  `set lets you add the reponame, repo url, repo authentication type and repo token`,
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
		// dsn := "http://localhost:9000?username=admin&password=admin"
		ctx := context.Background()
		dsnURL, err := url.Parse(Config.ClickHouseUrl)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("URL:", dsnURL)
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
		}
		defer db.Close()
		err = db.Ping(ctx)
		if err != nil {
			log.Println(err)
		}

		// Create table if not exists
		err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS signoz_traces.repo (
			name String,
			url String,
			auth_type String,
			token String,
			created_at DateTime DEFAULT now(),
			updated_at DateTime DEFAULT now()
		) ENGINE = MergeTree()
		ORDER BY created_at
		`)
		if err != nil {
			logger.Printf("error creating table: %v", err)
		}

		// Insert data into table
		err = db.Exec(ctx, `
			INSERT INTO signoz_traces.repo (name, url, auth_type, token) VALUES (?, ?, ?, ?)
		`, repoName, repoUrl, repoAuthType, repoToken)
		if err != nil {
			logger.Printf("error inserting data: %v", err)
		}
		log.Println("repo details added to signoz clickhouse successfully")

	},
}

func init() {
	repoSetCmd.PersistentFlags().StringVarP(&repoName, "name", "n", "", "name of repo")
	repoSetCmd.PersistentFlags().StringVarP(&repoUrl, "url", "u", "", "url of repo")
	repoSetCmd.PersistentFlags().StringVarP(&repoAuthType, "auth-type", "a", "", "authentication type of repo")
	repoSetCmd.PersistentFlags().StringVarP(&repoToken, "token", "t", "", "token of repo")

	repoCmd.AddCommand(repoSetCmd)
}

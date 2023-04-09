package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	CH_CONN, QT_CONN string
)

var cfgCmd = &cobra.Command{
	Use:   "config",
	Short: "config sets the configuration",
	Long:  "config sets the Clickhouse connection Url and Quality Trace connection Url",
	Run: func(cmd *cobra.Command, args []string) {
		config := viper.New()
		if CH_CONN != "" || QT_CONN != "" {
			if CH_CONN != "" {
				config.Set("CH_CONN", CH_CONN)
			}
			if QT_CONN != "" {
				config.Set("QT_CONN", QT_CONN)
			}
		}
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.Println(err)
			return
		}
		// Create the directory if it does not exist
		configDir := fmt.Sprintf("%s/config", homeDir)
		if err := os.MkdirAll(configDir, 0700); err != nil {
			logger.Println(err)
			return
		}
		configFilePath := fmt.Sprintf("%s/config.yaml", configDir)
		config.SetConfigFile(configFilePath)
		f, err := os.Create(configFilePath)
		if err != nil {
			logger.Println(err)
			return
		}
		defer f.Close()
		if err := config.WriteConfig(); err != nil {
			logger.Println(err)
			return
		}

	},
}

func init() {
	cfgCmd.PersistentFlags().StringVarP(&CH_CONN, "set-clickhouse", "c", "", "sets the clickhouse connection url")
	cfgCmd.PersistentFlags().StringVarP(&QT_CONN, "set-server", "s", "", "sets the quality trace connection url")
	rootCmd.AddCommand(cfgCmd)
}

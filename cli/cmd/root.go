/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/kube-tarian/quality-trace/qt/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var Config model.Env
var logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qt",
	Short: "Quality trace cli",
	Long: `Quality trace cli helps you communicate with QT server:
For example:
	qt --path : 
	lets you set the config path, 
	the default config path is $HOME/config/config.yaml

	qt config --set-clickhouse <clickhouse-url> --set-server <qt-server-url> : 
	lets you set the clickhouse and server url

	qt repo set --name <reponame> --url <repourl> --auth-type <authType> --token <gittoken> : 
	lets you set the git repo details

	qt repo list : lets you list the git repo details 

	qt repo dryrun <reponame> : 
	Lets you run a dry test to get the spans details using the check.yaml from the provided reponame.

	qt repo run <reponame> : 
	Lets you run a Test using the test.yaml from the provided reponame.

	qt repo drop: Lets you drop the repository table from the signoz table.
	
	qt repo update <reponame> --url <repourl> --auth-type <authType> --token <gittoken>: 
	Lets you update the repository details.

	qt test create --definition <definition-file.yml> : lets you create a test

	qt test run --definition <definition-file.yml> --testid <test-id> : lets you run a the test

	qt test create --definition <definition-file.yml> : lets you create a test

	qt test run --definition <definition-file.yml> --testid <test-id> : lets you run a the test
	`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "path", "", "config file (default is $HOME/config/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		homeDir := fmt.Sprintf("%s/config/", home)
		viper.AddConfigPath(homeDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match
	Config = model.Env{}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		logger.Println("Unable to Read The config:", err)
	}
	err := viper.Unmarshal(&Config)
	if err != nil {
		logger.Println("Unable to Read The config:", err)
	}
}

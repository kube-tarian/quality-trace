/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "repo",
	Long: `Repo command manages to Set the repo Details and 
	Tests`,
}

func init() {
	rootCmd.AddCommand(repoCmd)
}

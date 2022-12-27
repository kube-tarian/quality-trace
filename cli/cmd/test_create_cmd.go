/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a test",
	Long:  "create a test using a definition",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		fmt.Println(runTestFileDefinition)
	},
}

func init() {
	testCreateCmd.PersistentFlags().StringVarP(&runTestFileDefinition, "definition", "d", "", "--definition <definition-file.yml>")
	testCmd.AddCommand(testCreateCmd)
}

/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testRunCmd = &cobra.Command{
	Use:   "run",
	Short: "run a test",
	Long:  "run a test using a definition or an id",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
		fmt.Println(runTestFileDefinition)
		fmt.Println(testID)
	},
}

func init() {
	testRunCmd.PersistentFlags().StringVarP(&runTestFileDefinition, "definition", "d", "", "--definition <definition-file.yml>")
	testRunCmd.PersistentFlags().StringVarP(&testID, "test-id", "i", "", "--testid <test-id>")
	testCmd.AddCommand(testRunCmd)
}

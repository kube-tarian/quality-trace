/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var runTestFileDefinition string
var testID string

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "tests",
	Long:  `Test command lets you manage tests.`,
}

func init() {
	rootCmd.AddCommand(testCmd)
}

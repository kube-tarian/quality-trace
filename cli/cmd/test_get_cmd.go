/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get a test",
	Long:  `get a test using an id`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
		fmt.Println(testID)
	},
}

func init() {
	testGetCmd.PersistentFlags().StringVarP(&testID, "test-id", "i", "", "--testid <test-id>")
	testCmd.AddCommand(testGetCmd)
}

/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var testGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get a test",
	Long:  `get a test using an id`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
		fmt.Println(testID)
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/test/%v",testID))

		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(body))
	},
}

func init() {
	testGetCmd.PersistentFlags().StringVarP(&testID, "test-id", "i", "", "--testid <test-id>")
	testCmd.AddCommand(testGetCmd)
}

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	//  "/home/shifnazarnaz/clickhouse-integration/src/rest"
	"github.com/ShashankGirish/clickhouse-integration/src/adapters/clickhouseReader"
	"github.com/ShashankGirish/clickhouse-integration/src/assertions"
	"github.com/ShashankGirish/clickhouse-integration/src/parser"
	"github.com/ShashankGirish/clickhouse-integration/src/rest"
)

func main() {
	ctx := context.Background()
	var err error

	// Step 1: Establish connection with clickhouse.
	reader := clickhouseReader.NewReader("http://localhost:9000?user=default&password=")

	// Step 2: Clear the DB before running our test.
	err = reader.ClearTracesTable(ctx)
	if err != nil {
		fmt.Println("Err: ", err)
		os.Exit(1)
	}

	descriptorPath := os.Args[1]

	fmt.Println("Reading test descriptor:", descriptorPath)

	// Step 3: Read testDescriptor descriptor file and parse it
	testDescriptor, err := parser.ParseYaml(descriptorPath)
	if err != nil {
		fmt.Printf("Error while reading yaml file: %v", err)
		os.Exit(1)
	}
	fmt.Println(testDescriptor)

	// Step 4: Query the sample application
	resp, err := rest.SendQuery(testDescriptor.Trigger.HttpRequest)
	if err != nil {
		fmt.Println("Error while querying endpoint:", err)
		os.Exit(1)
	}

	fmt.Println(string(resp))
	time.Sleep(2 * time.Second) // It takes some time for the traces to reach signoz and then be populated in clickhouse db.

	testStatus := map[bool]string{true: "Success!", false: "Failure!"}

	for _, spec := range testDescriptor.Specs {
		// For each test specification, we need to retry the assertions in case the traces are not yet populated in clickhouse.
		// How many times to retry is configurable in the test descriptor.
		for i := -1; i < spec.MaxRetries; i++ {
			fmt.Println("Running assertion: ", spec.Name)

			// Step 5: Get the traces
			res, err := reader.GetTraces(ctx, spec.Selectors)
			if err != nil {
				fmt.Println("Error while retrieving traces:", err)
				os.Exit(1)
			}
			// fmt.Printf("Result: %v", res)

			// Step 6: Perform assertions
			isSuccess, match, err := assertions.RunAssertions(res, spec.Assertions)
			if err != nil {
				fmt.Println("Error during assertions:", err)
				os.Exit(1)
			}

			fmt.Println("\n\nTest Status:", testStatus[isSuccess])
			if isSuccess {
				fmt.Printf("Found trace %v with passing assertions.\n", match.TraceID)
				break
			} else if i+1 < spec.MaxRetries {
				fmt.Printf("\n\nRetrying in %v seconds...\n", spec.RetryInterval)
				time.Sleep(time.Duration(spec.RetryInterval) * time.Second)
			}
		}
	}

}

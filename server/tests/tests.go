package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kube-tarian/quality-trace/server/adapters/clickhousereader"
	"github.com/kube-tarian/quality-trace/server/assertions"
	"github.com/kube-tarian/quality-trace/server/model"
	"github.com/kube-tarian/quality-trace/server/rest"
)

type TestHandler struct {
	Ctx    context.Context
	Reader clickhousereader.ClickHouseReader
	Tests  map[int]*model.Test
}

func (h TestHandler) GetTest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	fmt.Println(id)
	index, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Err: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	testDescriptor := h.Tests[index]
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&testDescriptor)
	if err != nil {
		fmt.Println("There was an error encoding the initialized struct")
	}
	w.WriteHeader(http.StatusOK)
}

func (h TestHandler) CreateTest(w http.ResponseWriter, r *http.Request) {
	var testBody *model.Test
	json.NewDecoder(r.Body).Decode(&testBody)

	index := len(h.Tests) + 1
	testBody.Id = index
	testBody.Status = "Not Run"

	h.Tests[index] = testBody

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&testBody)
	if err != nil {
		fmt.Println("There was an error encoding the initialized struct")
	}
	w.WriteHeader(http.StatusOK)
}

func (h TestHandler) DeleteTest(w http.ResponseWriter, r *http.Request) {
}

func (h TestHandler) RunTest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	fmt.Println(id)
	index, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Err: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	testDescriptor := h.Tests[index]
	testDescriptor.Status = "Running"

	go func() {

		testStatus := map[bool]string{true: "Success!", false: "Failure!"}

		ctx := context.Background()
		err = h.Reader.ClearTracesTable(ctx)
		if err != nil {
			fmt.Println("Err: ", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Trigger API
		resp, err := rest.SendQuery(testDescriptor.Trigger.HttpRequest)
		if err != nil {
			fmt.Println("Error while querying endpoint:", err)
			os.Exit(1)
		}
		fmt.Println("Resp:", string(resp))

		// Wait for traces to be generated
		time.Sleep(2 * time.Second)

		runSuccess := false

		//  Run Assertions
		for _, spec := range testDescriptor.Specs {
			// For each test specification, we need to retry the assertions in case the traces are not yet populated in clickhouse.
			// How many times to retry is configurable in the test descriptor.
			for i := -1; i < spec.MaxRetries; i++ {
				fmt.Fprintf(w, "Running assertion: %v\n", spec.Name)

				// Step 5: Get the traces
				res, err := h.Reader.GetTraces(h.Ctx, spec.Selectors)
				fmt.Fprintf(w, "Trace: %v\n", res)
				//	fmt.Println("Trace Status is :",trace.ResponseStatusCode)
				if err != nil {
					fmt.Fprintf(w, "Error while retrieving traces: %v\n", err)
					os.Exit(1)
				}
				// fmt.Printf("Result: %v", res)

				// Perform assertions
				isSuccess, match, err := assertions.RunAssertions(res, spec.Assertions)
				if err != nil {
					fmt.Fprintf(w, "Error during assertions: %v\n", err)
					os.Exit(1)
				}

				fmt.Println("\n\nTest Status:", testStatus[isSuccess])
				if isSuccess {
					runSuccess = true
					fmt.Printf("Found trace %v with passing assertions.\n", match.TraceID)
					break
				} else if i+1 < spec.MaxRetries {
					fmt.Printf("\n\nRetrying in %v seconds...\n", spec.RetryInterval)
					time.Sleep(time.Duration(spec.RetryInterval) * time.Second)
				}
			}
		}

		testDescriptor.Status = testStatus[runSuccess]
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

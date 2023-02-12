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
	testBody.Status = "Dry Run Running"
	h.Tests[index] = testBody

	data := make(chan string)
	go func() {
		ctx := context.Background()
		err := h.Reader.ClearTracesTable(ctx)
		if err != nil {
			fmt.Println("Err: ", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Trigger API
		resp, err := rest.SendQuery(testBody.Trigger.HttpRequest)
		if err != nil {
			fmt.Println("Error while querying endpoint:", err)
			os.Exit(1)
		}
		fmt.Println("Resp:", string(resp))
		for {
			var traceId string
			traces, err := h.Reader.GetTrace(ctx)
			if err != nil {
				fmt.Println("Error while querying for traces:", err)
				os.Exit(1)
			}

			for _, trace := range *traces {
				if trace.TraceID != "" {
					traceId = trace.TraceID
				}
			}
			if traceId != "" {
				data <- traceId
				break
			}
			time.Sleep((1 * time.Second))
		}
	}()
	msg1 := <-data
	fmt.Println(msg1)
	resp, err := h.Reader.Searchtrace(context.Background(), msg1)
	if err != nil {
		fmt.Println("unable to search trace", err)
	}
	_ = json.NewEncoder(w).Encode(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h TestHandler) DeleteTest(w http.ResponseWriter, r *http.Request) {
}

func (h TestHandler) RunTest(w http.ResponseWriter, r *http.Request) {
	var testBody *model.Test
	json.NewDecoder(r.Body).Decode(&testBody)
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
	testBody.Id = index
	data := make(chan string)
	go func() {
		ctx := context.Background()
		err = h.Reader.ClearTracesTable(ctx)
		if err != nil {
			fmt.Println("Err: ", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Trigger API
		resp, err := rest.SendQuery(testBody.Trigger.HttpRequest)
		if err != nil {
			fmt.Println("Error while querying endpoint:", err)
			os.Exit(1)
		}
		fmt.Println("Resp:", string(resp))
		for {
			var traceId string
			traces, err := h.Reader.GetTrace(ctx)
			if err != nil {
				fmt.Println("Error while querying for traces:", err)
				os.Exit(1)
			}

			for _, trace := range *traces {
				if trace.TraceID != "" {
					traceId = trace.TraceID
				}
			}
			if traceId != "" {
				data <- traceId
				break
			}
			time.Sleep((1 * time.Second))
		}
	}()

	<-data
	//  Run Assertions
	for _, spec := range testBody.Specs {
		// Step 5: Get the traces
		res, err := h.Reader.GetTraces(h.Ctx, spec.Selectors)
		if err != nil {
			fmt.Fprintf(w, "Error while retrieving traces: %v\n", err)
			os.Exit(1)
		}
		// Perform assertions
		assertRess := assertions.RunAssertions(res, spec.Assertions)
		err = json.NewEncoder(w).Encode(assertRess)
		if err != nil {
			fmt.Println("error while encoding assertres", err)
		}
	}

	testBody.Status = "Success"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

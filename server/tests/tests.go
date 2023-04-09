package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"gopkg.in/yaml.v2"

	"github.com/gorilla/mux"
	"github.com/kube-tarian/quality-trace/server/adapters/clickhousereader"
	"github.com/kube-tarian/quality-trace/server/assertions"
	"github.com/kube-tarian/quality-trace/server/model"
	"github.com/kube-tarian/quality-trace/server/rest"
)

var logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

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

func (h TestHandler) RunAssertion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("run clone test")
	// we can get repo details from the post request
	var repoDetails model.Repo
	err := json.NewDecoder(r.Body).Decode(&repoDetails)
	if err != nil {
		log.Println("unable to fetch the data:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	file := "test.yaml"
	testBody, err := CloneAndParse(w, repoDetails.URL, file, repoDetails.Token)
	if err != nil {
		log.Println("Unable to clone the repo:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// data channel is used to send the testresults
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
		testBody.Status = "Success"
	}
}

func (h TestHandler) DryRun(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DryRun...")
	// we can get repo details from the post request
	var repoDetails model.Repo
	err := json.NewDecoder(r.Body).Decode(&repoDetails)
	if err != nil {
		logger.Println("unable to fetch the data:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	file := "check.yaml"
	fmt.Fprintln(w, "Cloning the repo...")
	testBody, err := CloneAndParse(w, repoDetails.URL, file, repoDetails.Token)
	if err != nil {
		logger.Println("Unable to clone the repo:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	progress := make(chan string)
	response := make(chan string)
	go func() {
		ctx := context.Background()
		err := h.Reader.ClearTracesTable(ctx)
		if err != nil {
			logger.Println("Err: ", err)
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
				logger.Println("Error while querying for traces:", err)
				os.Exit(1)
			}

			for _, trace := range *traces {
				if trace.TraceID != "" {
					traceId = trace.TraceID
				}
			}
			if traceId != "" {
				response <- traceId
				break
			}
			time.Sleep((1 * time.Second))
			progress <- "Fetching Data..."
		}
	}()

	for {
		select {
		case msg := <-progress:
			// send progress updates to the client
			fmt.Fprintln(w, msg)
		case traceID := <-response:
			// send the final response to the client
			fmt.Println(traceID)
			resp, err := h.Reader.Searchtrace(context.Background(), traceID)
			if err != nil {
				logger.Println("unable to search trace", err)
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
	}
}

// var testBody *model.Test
// repoURL should be a https , dont use ssh repourl
// use a classic token with permissions to clone and read
func CloneAndParse(w http.ResponseWriter, repoURL, file, token string) (*model.Test, error) {
	var auth *githttp.BasicAuth

	if token != "" {
		auth = &githttp.BasicAuth{
			Username: "user-name", // yes, this can be anything except an empty string
			Password: token,
		}
	} else {
		auth = &githttp.BasicAuth{}
	}

	// Clone the repository into an in-memory storage
	memStorage := memory.NewStorage()
	_, err := git.Clone(memStorage, nil, &git.CloneOptions{
		URL:      repoURL,
		Auth:     auth,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}

	// Get the HEAD commit's tree
	repo, err := git.Open(memStorage, nil)
	if err != nil {
		return nil, err
	}
	headRef, err := repo.Head()
	if err != nil {
		return nil, err
	}
	commit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, err
	}
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	// Find the YAML file in the repository's tree
	var configData []byte
	for _, entry := range tree.Entries {
		fmt.Printf("\n entry: %v", entry.Name)
		if entry.Name == file {
			blob, err := repo.BlobObject(entry.Hash)
			if err != nil {
				return nil, err
			}
			Data, err := blob.Reader()
			if err != nil {
				return nil, err
			}
			configData, err = io.ReadAll(Data)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	// Parse the YAML data into a Config struct
	var testBody model.Test
	err = yaml.Unmarshal(configData, &testBody)
	if err != nil {
		return nil, err
	}

	return &testBody, nil
}

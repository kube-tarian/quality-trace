package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kube-tarian/quality-trace/server/adapters/clickhousereader"
	"github.com/kube-tarian/quality-trace/server/config"
	"github.com/kube-tarian/quality-trace/server/model"
	"github.com/kube-tarian/quality-trace/server/tests"
)

var cfg = flag.String("config", "config.yaml", "path to the config file")

func main() {
	ctx := context.Background()
	fmt.Println("Context in main func", ctx)

	flag.Parse()
	cfg, err := config.FromFile(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize reader with connection to clickhouse
	connUrl := cfg.ClickhouseConnUrl
	reader := clickhousereader.NewReader(connUrl)

	// Initializing mux router and performing handler function
	s := mux.NewRouter()

	r := &tests.TestHandler{Ctx: ctx, Reader: *reader, Tests: map[int]*model.Test{}}

	s.HandleFunc("/test/", r.CreateTest).Methods("POST")
	// s.HandleFunc("/test/delete", r.DeleteTest)
	s.HandleFunc("/test/{id:[0-9]+}", r.GetTest).Methods("GET")
	s.HandleFunc("/test/{id:[0-9]+}/run", r.RunTest).Methods("POST")

	fmt.Printf("Server started at :8080...")
	http.ListenAndServe(":8080", s)
}

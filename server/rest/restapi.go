package rest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/kube-tarian/quality-trace/server/model"
)

func SendQuery(request model.HttpRequest) ([]byte, error) {
	var response *http.Response
	var err error

	fullUrl := fmt.Sprintf("http://%s%s", request.Url, request.Route)

	switch request.Method {
	case "GET":
		fmt.Println("Getting content..")
		response, err = http.Get(fullUrl)
	case "POST":
		fmt.Println("Posting content..")
		bodyBuf := bytes.NewBuffer([]byte(request.Body))
		response, err = http.Post(fullUrl, "application/json", bodyBuf)
	}

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData, err
}

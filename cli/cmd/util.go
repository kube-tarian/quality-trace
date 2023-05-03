package cmd

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func GenerateYaml(values map[string]string) {
	fmt.Println("\n Generating assertions yaml file")
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error occured while getting working dir %v", err)
	}
	parentDir := filepath.Dir(wd)
	assertionsDir := path.Join(parentDir, "generatedAssertionFiles")
	os.Mkdir(assertionsDir, 0777)
	fileName := fmt.Sprintf("assertions-%v.yaml", time.Now().Format(time.RFC822))
	filePath := path.Join(assertionsDir, fileName)

	assertVals := map[string]map[string]string{"spec": values}

	data, err := yaml.Marshal(&assertVals)

	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filePath, data, 0666)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n assertions yaml file created succesfully at path: %v", filePath)
}

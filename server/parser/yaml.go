package parser

import (
	"io/ioutil"
	"path/filepath"

	"github.com/kube-tarian/quality-trace/server/model"
	"gopkg.in/yaml.v2"
)

func ParseYaml(path string) (*model.Test, error) {
	var testModel model.Test

	filename, _ := filepath.Abs(path)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &testModel)
	if err != nil {
		return nil, err
	}

	return &testModel, nil
}

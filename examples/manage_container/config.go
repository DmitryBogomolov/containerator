package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	ImageName     string            `yaml:"image"`
	ImageRepo     string            `yaml:"image-repo"`
	ContainerName string            `yaml:"container"`
	Volumes       map[string]string `yaml:"volumes"`
	Ports         map[string]string `yaml:"ports"`
	Env           map[string]string `yaml:"environment"`
	Modes         map[string]config `yaml:"$modes"`
}

func readConfig(pathToFile string) (*config, error) {
	bytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}
	data := &config{}
	err = yaml.Unmarshal(bytes, data)
	return data, err
}

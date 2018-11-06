package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	ImageRepo     string `yaml:"image_repo"`
	ContainerName string `yaml:"container_name"`
	Network       string `yaml:"network"`

	BasePort   float64 `yaml:"base_port"`
	PortOffset float64 `yaml:"port_offset'`

	Ports   []float64         `yaml:"ports"`
	Volumes map[string]string `yaml:"volumes"`
	Env     map[string]string `yaml:"env"`

	Modes []string `yaml:"modes"`
}

func readConfig(pathToFile string) (*config, error) {
	bytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}
	var data config
	err = yaml.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	if len(data.Modes) == 0 {
		data.Modes = []string{defaultMode}
	}

	return &data, nil
}

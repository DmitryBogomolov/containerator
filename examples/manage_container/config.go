package main

import (
	"io/ioutil"

	"github.com/DmitryBogomolov/containerator"
	"gopkg.in/yaml.v2"
)

type config struct {
	ImageRepo  string                  `yaml:"image_repo"`
	Network    string                  `yaml:"network"`
	BasePort   float64                 `yaml:"base_port"`
	PortOffset float64                 `yaml:"port_offset"`
	Ports      []float64               `yaml:"ports"`
	Volumes    []containerator.Mapping `yaml:"volumes"`
	Env        []containerator.Mapping `yaml:"env"`
	Modes      []string                `yaml:"modes"`
}

func readConfig(pathToFile string) (*config, error) {
	bytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}
	var data config
	err = yaml.Unmarshal(bytes, &data)
	return &data, err
}

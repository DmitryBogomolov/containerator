package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// TODO: Embedded fields are not unmarshalled.

// type basicConfig struct {
// 	ContainerName string            `yaml:"container"`
// 	Volumes       map[string]string `yaml:"volumes"`
// 	Ports         map[string]string `yaml:"ports"`
// 	Env           map[string]string `yaml:"env"`
// 	EnvFile       string            `yaml:"env_file"`
// }

type modeConfig struct {
	// basicConfig

	Name       string `yaml:"name"`
	PortOffset string `yaml:"port_offset"`

	ContainerName string            `yaml:"container"`
	Volumes       map[string]string `yaml:"volumes"`
	Ports         map[string]string `yaml:"ports"`
	Env           map[string]string `yaml:"env"`
	EnvFile       string            `yaml:"env_file"`
}

type config struct {
	//basicConfig

	ImageName string       `yaml:"image"`
	ImageRepo string       `yaml:"image-repo"`
	Network   string       `yaml:"network"`
	Modes     []modeConfig `yaml:"modes"`

	ContainerName string            `yaml:"container"`
	Volumes       map[string]string `yaml:"volumes"`
	Ports         map[string]string `yaml:"ports"`
	Env           map[string]string `yaml:"env"`
	EnvFile       string            `yaml:"env_file"`
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

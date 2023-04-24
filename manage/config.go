package manage

import (
	"io/ioutil"

	"github.com/DmitryBogomolov/containerator/core"
	"gopkg.in/yaml.v2"
)

// Config contains options for container management.
type Config struct {
	ImageName     string         `yaml:"image_name"`               // Image name; required
	ContainerName string         `yaml:"container_name,omitempty"` // Container name
	Network       string         `yaml:",omitempty"`               // Container network
	BasePort      float64        `yaml:"base_port,omitempty"`
	PortOffset    float64        `yaml:"port_offset,omitempty"`
	Ports         []float64      `yaml:",omitempty"`
	Volumes       []core.Mapping `yaml:",omitempty"` // Volumes mapping
	Env           []core.Mapping `yaml:",omitempty"` // Environment variables
	Modes         []string       `yaml:",omitempty"`
}

// ReadConfig reads config from yaml file.
//
//	ReadConfig("/path/to/config,yaml") -> &config, err
func ReadConfig(pathToFile string) (*Config, error) {
	bytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}
	var conf Config
	err = yaml.Unmarshal(bytes, &conf)
	return &conf, err
}

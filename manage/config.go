package manage

import (
	"io/ioutil"

	"github.com/DmitryBogomolov/containerator/core"
	"gopkg.in/yaml.v2"
)

// Config contains options for container management.
type Config struct {
	ImageName     string         `yaml:"image_repo"`     // Image name; required
	ContainerName string         `yaml:"container_name"` // Container name
	Network       string         `yaml:"network"`        // Container network
	BasePort      float64        `yaml:"base_port"`
	PortOffset    float64        `yaml:"port_offset"`
	Ports         []float64      `yaml:"ports"`
	Volumes       []core.Mapping `yaml:"volumes"` // Volumes mapping
	Env           []core.Mapping `yaml:"env"`     // Environment variables
	Modes         []string       `yaml:"modes"`
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

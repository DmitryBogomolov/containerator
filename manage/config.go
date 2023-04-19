package manage

import (
	"io/ioutil"

	"github.com/DmitryBogomolov/containerator"
	"gopkg.in/yaml.v2"
)

// Config contains options for container management.
//
// `ImageRepo` is required, others are optional.
type Config struct {
	ImageRepo     string                  `yaml:"image_repo"`
	ContainerName string                  `yaml:"container_name"`
	Network       string                  `yaml:"network"`
	BasePort      float64                 `yaml:"base_port"`
	PortOffset    float64                 `yaml:"port_offset"`
	Ports         []float64               `yaml:"ports"`
	Volumes       []containerator.Mapping `yaml:"volumes"`
	Env           []containerator.Mapping `yaml:"env"`
	Modes         []string                `yaml:"modes"`
}

// ReadConfig read config from yaml file.
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

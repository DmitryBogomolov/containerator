package manage

import (
	"io/ioutil"
	"path/filepath"

	"github.com/DmitryBogomolov/containerator/core"
	"gopkg.in/yaml.v2"
)

// Config contains options for container management.
type Config struct {
	ImageName     string         `yaml:"image_name"`               // Image name; required
	ContainerName string         `yaml:"container_name,omitempty"` // Container name
	Network       string         `yaml:",omitempty"`               // Container network
	Ports         []core.Mapping `yaml:",omitempty"`               // Ports mapping
	Volumes       []core.Mapping `yaml:",omitempty"`               // Volumes mapping
	Env           []core.Mapping `yaml:",omitempty"`               // Environment variables
}

// ReadConfig reads config from yaml file.
//
//	ReadConfig("/path/to/config,yaml") -> &config, err
func ReadConfig(pathToFile string) (*Config, error) {
	bytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(bytes, &cfg)
	if err == nil {
		curDir, _ := filepath.Abs(filepath.Dir(pathToFile))
		processConfig(&cfg, curDir)
	}
	return &cfg, err
}

func processConfig(cfg *Config, dir string) {
	for i, mapping := range cfg.Volumes {
		hostPath := filepath.Clean(mapping.Source)
		if !filepath.IsAbs(hostPath) {
			hostPath = filepath.Join(dir, hostPath)
		}
		cfg.Volumes[i].Source = hostPath
	}
}

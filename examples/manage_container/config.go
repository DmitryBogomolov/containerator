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

	Name       string  `yaml:"name"`
	PortOffset float64 `yaml:"port_offset"`

	ContainerName string              `yaml:"container"`
	Volumes       map[string]string   `yaml:"volumes"`
	Ports         map[float64]float64 `yaml:"ports"`
	Env           map[string]string   `yaml:"env"`
	EnvFile       string              `yaml:"env_file"`
}

type config struct {
	//basicConfig

	ImageName string       `yaml:"image"`
	ImageRepo string       `yaml:"image-repo"`
	Network   string       `yaml:"network"`
	Modes     []modeConfig `yaml:"modes"`

	ContainerName string              `yaml:"container"`
	Volumes       map[string]string   `yaml:"volumes"`
	Ports         map[float64]float64 `yaml:"ports"`
	Env           map[string]string   `yaml:"env"`
	EnvFile       string              `yaml:"env_file"`
}

func readConfig(pathToFile string) (*config, error) {
	bytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}
	var data config
	err = yaml.Unmarshal(bytes, &data)

	if len(data.Modes) == 0 {
		data.Modes = append(data.Modes, modeConfig{
			Name: defaultMode,
		})
	}

	return &data, err
}

func selectModeConfig(conf *config, mode string) *modeConfig {
	for i, item := range conf.Modes {
		if item.Name == mode {
			return &conf.Modes[i]
		}
	}
	return nil
}

func buildModeConfig(conf *config, mode string) *config {
	modeConf := selectModeConfig(conf, mode)
	if modeConf == nil {
		return nil
	}

	ret := *conf
	ret.Ports = make(map[float64]float64)
	for key, val := range conf.Ports {
		ret.Ports[key] = val
	}
	ret.Volumes = make(map[string]string)
	for key, val := range conf.Volumes {
		ret.Volumes[key] = val
	}
	ret.Env = make(map[string]string)
	for key, val := range conf.Env {
		ret.Env[key] = val
	}

	if modeConf.ContainerName != "" {
		ret.ContainerName = modeConf.ContainerName
	}
	if modeConf.PortOffset > 0 {
		for key, val := range ret.Ports {
			delete(ret.Ports, key)
			ret.Ports[key+modeConf.PortOffset] = val
		}
	}
	if len(modeConf.Ports) > 0 {
		for key, val := range modeConf.Ports {
			ret.Ports[key] = val
		}
	}
	if len(modeConf.Volumes) > 0 {
		for key, val := range modeConf.Volumes {
			ret.Volumes[key] = val
		}
	}
	if len(modeConf.Env) > 0 {
		for key, val := range modeConf.Env {
			ret.Env[key] = val
		}
	}
	if modeConf.EnvFile != "" {
		ret.EnvFile = modeConf.EnvFile
	}
	return &ret
}

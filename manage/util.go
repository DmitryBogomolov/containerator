package manage

import (
	"fmt"
	"strconv"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
)

func findIndex[T comparable](item T, list []T) int {
	for i, obj := range list {
		if obj == item {
			return i
		}
	}
	return -1
}

// NotValidModeError indicates that specified mode is not found amoung config modes.
type NotValidModeError struct {
	mode  string
	modes []string
}

func (err *NotValidModeError) Error() string {
	return fmt.Sprintf("mode '%s' is not valid", err.mode)
}

// Mode returns mode.
func (err *NotValidModeError) Mode() string {
	return err.mode
}

// Modes returns config modes.
func (err *NotValidModeError) Modes() []string {
	return append([]string(nil), err.modes...)
}

func getContainerName(conf *Config, postfix string) string {
	name := conf.ContainerName
	if name == "" {
		name = conf.ImageName
	}
	if postfix != "" {
		name += "-" + postfix
	}
	return name
}

func removeContainer(cli client.ContainerAPIClient, container core.Container, name string) (core.Container, error) {
	if container == nil {
		return nil, &NoContainerError{name}
	}
	if err := core.RemoveContainer(cli, container); err != nil {
		return nil, err
	}
	return container, nil
}

func findImage(cli client.ImageAPIClient, name string, tag string) (core.Image, error) {
	imageName := name
	if tag != "" {
		imageName += ":" + tag
	}
	return core.FindImageByName(cli.(client.ImageAPIClient), imageName)
}

func buildContainerOptions(
	cfg *Config, imageName string, containerName string, options *Options,
) (*core.RunContainerOptions, error) {
	result := core.RunContainerOptions{
		Image:         imageName,
		Name:          containerName,
		RestartPolicy: core.RestartAlways,
		Network:       cfg.Network,
		Volumes:       cfg.Volumes,
		Env:           cfg.Env,
	}
	if options.PortOffset != 0 && len(cfg.Ports) > 0 {
		ports := make([]core.Mapping, len(cfg.Ports))
		for i, mapping := range cfg.Ports {
			val, _ := strconv.Atoi(mapping.Source)
			ports[i] = core.Mapping{
				Source: strconv.Itoa(val + options.PortOffset),
				Target: mapping.Target,
			}
		}
		result.Ports = ports
	}
	if options.EnvFilePath != "" {
		data, err := godotenv.Read(options.EnvFilePath)
		if err != nil {
			return nil, err
		}
		mappings := make([]core.Mapping, 0, len(data))
		for key, val := range data {
			mappings = append(mappings, core.Mapping{Source: key, Target: val})
		}
		result.Env = append(result.Env, mappings...)
	}
	return &result, nil
}

package manage

import (
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
	image, err := core.FindImageByName(cli.(client.ImageAPIClient), imageName)
	if err != nil {
		return nil, err
	}
	if image == nil {
		return nil, &NoImageError{imageName}
	}
	return image, nil
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
		Ports:         cfg.Ports,
		Env:           cfg.Env,
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

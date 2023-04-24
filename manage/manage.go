// Package manage contains function to run, suspend, resume, remove containers.
package manage

import (
	"fmt"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
)

func updateContainer(
	options *core.RunContainerOptions, currentContainer core.Container, cli client.ContainerAPIClient,
) (container core.Container, err error) {
	if currentContainer != nil {
		if err = core.SuspendContainer(cli, currentContainer); err != nil {
			return
		}
		defer func() {
			if err != nil {
				if otherErr := core.ResumeContainer(cli, currentContainer, options.Name); otherErr != nil {
					err = fmt.Errorf("%v (%v)", err, otherErr)
				}
			} else {
				err = core.RemoveContainer(cli, currentContainer)
			}
		}()
	}
	container, err = core.RunContainer(cli, options)
	return
}

// Options contains additional arguments for Manage function.
type Options struct {
	Mode           string // If set should match one of modes in config
	Tag            string // Image tag; if not set newest image is selected
	Force          bool   // If set running container is replaced
	Remove         bool   // If set running container is removed
	GetEnvFilePath func(string) string
}

// DefaultConfigName defines default name of config file.
const DefaultConfigName = "config.yaml"

// Manage runs containers with the last tag for the specified image repo.
//
//	Manage(cli, "/path/to/config.yaml", &Options{Mode:"dev"}) -> &container, err
func Manage(cli interface{}, cfg *Config, options *Options) (core.Container, error) {
	mode, modeIndex, err := selectMode(options.Mode, cfg)
	if err != nil {
		return nil, err
	}

	containerName := getContainerName(cfg, mode)

	containerCli := cli.(client.ContainerAPIClient)
	currentContainer, err := core.FindContainerByName(containerCli, containerName)
	if err != nil {
		if _, ok := err.(*core.ContainerNotFoundError); !ok {
			return nil, err
		}
		currentContainer = nil
	}

	if options.Remove {
		if currentContainer == nil {
			return nil, &NoContainerError{containerName}
		}
		if err = core.RemoveContainer(containerCli, currentContainer); err != nil {
			return nil, err
		}
		return currentContainer, nil
	}

	imageName := cfg.ImageName
	if options.Tag != "" {
		imageName += ":" + options.Tag
	}
	image, err := core.FindImageByName(cli.(client.ImageAPIClient), cfg.ImageName)
	if err != nil {
		return nil, err
	}

	if currentContainer != nil && currentContainer.ImageID() == image.ID() && !options.Force {
		return nil, &ContainerAlreadyRunningError{currentContainer.Name()}
	}

	runOptions := buildContainerOptions(cfg, image.FullName(), containerName, modeIndex)
	if options.GetEnvFilePath != nil {
		filePath := options.GetEnvFilePath(mode)
		data, err := godotenv.Read(filePath)
		if err != nil {
			return nil, err
		}
		mappings := make([]core.Mapping, len(data))
		for key, val := range data {
			mappings = append(mappings, core.Mapping{Source: key, Target: val})
		}
		runOptions.Env = append(runOptions.Env, mappings...)
	}
	return updateContainer(runOptions, currentContainer, containerCli)
}

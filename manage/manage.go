// Package manage contains function to run, suspend, resume, remove containers.
package manage

import (
	"fmt"
	"io"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func updateContainer(
	options *core.RunContainerOptions, currentContainer *types.Container, cli client.ContainerAPIClient,
) (container *types.Container, err error) {
	if currentContainer != nil {
		currentContainerID := currentContainer.ID
		if err = core.SuspendContainer(cli, currentContainerID); err != nil {
			return
		}
		defer func() {
			if err != nil {
				if otherErr := core.ResumeContainer(cli, currentContainerID, options.Name); otherErr != nil {
					err = fmt.Errorf("%v (%v)", err, otherErr)
				}
			} else {
				err = core.RemoveContainer(cli, currentContainerID)
			}
		}()
	}
	container, err = core.RunContainer(cli, options)
	return
}

// Options contains additional arguments for Manage function.
type Options struct {
	Mode         string // If set should match one of modes in config
	Tag          string // Image tag; if not set newest image is selected
	Force        bool   // If set running container is replaced
	Remove       bool   // If set running container is removed
	GetEnvReader func(string) (io.Reader, error)
}

// DefaultConfigName defines default name of config file.
const DefaultConfigName = "config.yaml"

// Manage runs containers with the last tag for the specified image repo.
//
//	Manage(cli, "/path/to/config.yaml", &Options{Mode:"dev"}) -> &container, err
func Manage(cli interface{}, cfg *Config, options *Options) (*types.Container, error) {
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
		if err = core.RemoveContainer(containerCli, currentContainer.ID); err != nil {
			return nil, err
		}
		return currentContainer, nil
	}

	image, err := findImage(cli.(client.ImageAPIClient), cfg.ImageName, options.Tag)
	if err != nil {
		return nil, err
	}
	imageName := core.GetImageFullName(image)

	if currentContainer != nil && currentContainer.ImageID == image.ID && !options.Force {
		return nil, &ContainerAlreadyRunningError{currentContainer}
	}

	runOptions := buildContainerOptions(cfg, imageName, containerName, modeIndex)
	if options.GetEnvReader != nil {
		reader, err := options.GetEnvReader(mode)
		if err != nil {
			return nil, err
		}
		runOptions.EnvReader = reader
	}
	return updateContainer(runOptions, currentContainer, containerCli)
}

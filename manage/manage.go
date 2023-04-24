// Package manage contains function to run, suspend, resume, remove containers.
package manage

import (
	"fmt"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/docker/docker/client"
)

func updateContainer(
	cli client.ContainerAPIClient, options *core.RunContainerOptions, currentContainer core.Container,
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
	Postfix     string // Container name postfix
	Tag         string // Image tag; if not set newest image is selected
	Force       bool   // If set running container is replaced
	Remove      bool   // If set running container is removed
	PortOffset  int    // Host machine port offset
	EnvFilePath string // Path to file with additional environment variables
}

// DefaultConfigName defines default name of config file.
const DefaultConfigName = "config.yaml"

// RunContainer runs container with the last tag for the specified image.
//
//	RunContainer(cli, "/path/to/config.yaml", &Options{Mode:"dev"}) -> &container, err
func RunContainer(cli interface{}, cfg *Config, options *Options) (core.Container, error) {
	containerName := getContainerName(cfg, options.Postfix)

	containerCli := cli.(client.ContainerAPIClient)
	currentContainer, err := core.FindContainerByName(containerCli, containerName)
	if err != nil {
		return nil, err
	}

	if options.Remove {
		return removeContainer(containerCli, currentContainer, containerName)
	}

	image, err := findImage(cli.(client.ImageAPIClient), cfg.ImageName, options.Tag)
	if err != nil {
		return nil, err
	}

	if currentContainer != nil && currentContainer.ImageID() == image.ID() && !options.Force {
		return nil, &ContainerAlreadyRunningError{currentContainer.Name()}
	}

	runOptions, err := buildContainerOptions(cfg, image.FullName(), containerName, options)
	if err != nil {
		return nil, err
	}
	return updateContainer(containerCli, runOptions, currentContainer)
}

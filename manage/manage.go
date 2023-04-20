/*
Package manage contains function to run, suspend, resume, remove containers.
*/
package manage

import (
	"fmt"
	"io"

	containerator "github.com/DmitryBogomolov/containerator/core"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func updateContainer(options *containerator.RunContainerOptions, currentContainer *types.Container,
	cli client.ContainerAPIClient) (container *types.Container, err error) {
	if currentContainer != nil {
		currentContainerID := currentContainer.ID
		if err = containerator.SuspendContainer(cli, currentContainerID); err != nil {
			return
		}
		defer func() {
			if err != nil {
				otherErr := containerator.ResumeContainer(cli, currentContainerID, options.Name)
				if otherErr != nil {
					err = fmt.Errorf("%v (%v)", err, otherErr)
				}
			} else {
				err = containerator.RemoveContainer(cli, currentContainerID)
			}
		}()
	}
	container, err = containerator.RunContainer(cli, options)
	return
}

/*
Options type contains additional arguments for Manage function.

*Mode* might be required (depends on *modes* in config), all others are optional.
*/
type Options struct {
	Mode         string
	Tag          string
	Force        bool
	Remove       bool
	GetEnvReader func(string) (io.Reader, error)
}

// DefaultConfigName defines default name of config file.
const DefaultConfigName = "config.yaml"

// NoContainerError is returned on attempt to remove container when it is not found.
type NoContainerError struct {
	container string
}

func (err *NoContainerError) Error() string {
	return fmt.Sprintf("container '%s' is not found", err.container)
}

// Container returns container name.
func (err *NoContainerError) Container() string {
	return err.container
}

// ContainerAlreadyRunningError is returned on attempt to run container when similar container is already running.
type ContainerAlreadyRunningError struct {
	container *types.Container
}

func (err *ContainerAlreadyRunningError) Error() string {
	return fmt.Sprintf("container '%s' is already running", containerator.GetContainerName(err.container))
}

// Container returns running container.
func (err *ContainerAlreadyRunningError) Container() *types.Container {
	return err.container
}

/*
Manage runs containers with the last tag for the specified image repo.

*Mode* should match those in config, otherwise shouldn't be defined.
The newest tag is selected (if *Tag* is not defined).
Use *Force* to override currently running container.
Use *Remove* to remove currently running container.

	Manage("/path/to/config.yaml", &Options{Mode:"dev"}) -> &container, err
*/
func Manage(cli interface{}, conf *Config, options *Options) (*types.Container, error) {
	mode, modeIndex, err := selectMode(options.Mode, conf)
	if err != nil {
		return nil, err
	}

	containerName := getContainerName(conf, mode)

	containerCli := cli.(client.ContainerAPIClient)
	currentContainer, err := containerator.FindContainerByName(containerCli, containerName)
	if err != nil {
		if _, ok := err.(*containerator.ContainerNotFoundError); !ok {
			return nil, err
		}
		currentContainer = nil
	}

	if options.Remove {
		if currentContainer == nil {
			return nil, &NoContainerError{containerName}
		}
		err = containerator.RemoveContainer(containerCli, currentContainer.ID)
		if err != nil {
			return nil, err
		}
		return currentContainer, nil
	}

	image, err := findImage(cli.(client.ImageAPIClient), conf.ImageRepo, options.Tag)
	if err != nil {
		return nil, err
	}
	imageName := containerator.GetImageFullName(image)

	if currentContainer != nil && currentContainer.ImageID == image.ID && !options.Force {
		return nil, &ContainerAlreadyRunningError{currentContainer}
	}

	runOptions := buildContainerOptions(conf, imageName, containerName, modeIndex)
	if options.GetEnvReader != nil {
		reader, err := options.GetEnvReader(mode)
		if err != nil {
			return nil, err
		}
		runOptions.EnvReader = reader
	}
	return updateContainer(runOptions, currentContainer, containerCli)
}

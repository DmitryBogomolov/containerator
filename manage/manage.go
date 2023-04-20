/*
Package manage contains function to run, suspend, resume, remove containers.
*/
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
//
// `Mode` might be required (depends on `Modes` in config), all others are optional.
type Options struct {
	Mode         string
	Tag          string
	Force        bool
	Remove       bool
	GetEnvReader func(string) (io.Reader, error)
}

// DefaultConfigName defines default name of config file.
const DefaultConfigName = "config.yaml"

/*
Manage runs containers with the last tag for the specified image repo.

`Mode` should match those in config, otherwise shouldn't be defined.
The newest tag is selected (if `Tag` is not defined).
Use `Force` to override currently running container.
Use `Remove` to remove currently running container.

	Manage("/path/to/config.yaml", &Options{Mode:"dev"}) -> &container, err
*/
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
		err = core.RemoveContainer(containerCli, currentContainer.ID)
		if err != nil {
			return nil, err
		}
		return currentContainer, nil
	}

	image, err := findImage(cli.(client.ImageAPIClient), cfg.ImageRepo, options.Tag)
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

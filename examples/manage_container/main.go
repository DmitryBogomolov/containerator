// Program manage_container shows usage of *containerator* functions that run and remove containers.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func buildContainerOptions(config *config, imageName string, containerName string, modeIndex int) *containerator.RunContainerOptions {
	ret := containerator.RunContainerOptions{
		Image:         imageName,
		Name:          containerName,
		RestartPolicy: containerator.RestartAlways,
		Network:       config.Network,
		Volumes:       config.Volumes,
		Env:           config.Env,
	}
	basePort := int(config.BasePort) + int(config.PortOffset)*modeIndex
	if len(config.Ports) > 0 {
		ports := make([]containerator.Mapping, len(config.Ports))
		for i, port := range config.Ports {
			ports[i] = containerator.Mapping{
				Source: strconv.Itoa(basePort + i),
				Target: strconv.Itoa(int(port)),
			}
		}
		ret.Ports = ports
	}
	return &ret
}

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

const defaultConfigName = "config.yaml"

func run() error {
	var configPathOption string
	flag.StringVar(&configPathOption, "config", defaultConfigName, "configuration file")
	var modeOption string
	flag.StringVar(&modeOption, "mode", "", "mode")
	var tagOption string
	flag.StringVar(&tagOption, "tag", "", "image tag")
	var removeOption bool
	flag.BoolVar(&removeOption, "remove", false, "remove container")
	var forceOption bool
	flag.BoolVar(&forceOption, "force", false, "force container creation")

	flag.Parse()

	config, err := readConfig(configPathOption)
	if err != nil {
		return err
	}

	mode, modeIndex, err := selectMode(modeOption, config)
	if err != nil {
		return err
	}

	containerName := getContainerName(config.ImageRepo, mode)
	log.Printf("Container: %s\n", containerName)

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	currentContainer, err := containerator.FindContainerByName(cli, containerName)
	if err != nil {
		return err
	}

	if removeOption {
		if currentContainer != nil {
			err = containerator.RemoveContainer(cli, currentContainer.ID)
			if err != nil {
				return err
			}
			log.Println("Container is removed")

		} else {
			log.Println("There is no container")
		}
		return nil
	}

	image, err := findImage(cli, config.ImageRepo, tagOption)
	if err != nil {
		return err
	}
	imageName := containerator.GetImageFullName(image)
	log.Printf("Image: %s\n", imageName)

	if currentContainer != nil && currentContainer.ImageID == image.ID && !forceOption {
		log.Println("Container is already running")
		return nil
	}

	options := buildContainerOptions(config, imageName, containerName, modeIndex)
	if reader := getEnvFileReader(configPathOption, mode); reader != nil {
		options.EnvReader = reader
	}

	container, err := updateContainer(options, currentContainer, cli)
	if err != nil {
		return err
	}

	log.Printf("Container: %s %s\n",
		containerator.GetContainerName(container), containerator.GetContainerShortID(container))

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
}

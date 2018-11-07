package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func selectMode(mode string, config *config) string {
	for _, item := range config.Modes {
		if item == mode {
			return mode
		}
	}
	return ""
}

func findImage(cli client.ImageAPIClient, config *config) (*types.ImageSummary, error) {
	list, err := containerator.FindImagesByRepo(cli, config.ImageRepo)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("no %s images", config.ImageRepo)
	}
	return list[0], nil
}

func buildContainerOptions(config *config, imageName string, mode string) *containerator.RunContainerOptions {
	containerName := config.ImageRepo
	if mode != "" {
		containerName += "-" + mode
	}
	ret := containerator.RunContainerOptions{
		Image:         imageName,
		Name:          containerName,
		RestartPolicy: containerator.RestartAlways,
		Network:       config.Network,
		Volumes:       containerator.NewMappingListFromMap(config.Volumes),
		Env:           containerator.NewMappingListFromMap(config.Env),
	}
	modeOffset := 0
	for i, val := range config.Modes {
		if val == mode {
			modeOffset = i
			break
		}
	}
	basePort := int(config.BasePort) + int(config.PortOffset)*modeOffset
	if len(config.Ports) > 0 {
		ports := make([]containerator.Mapping, 0, len(config.Ports))
		for i, port := range config.Ports {
			ports = append(ports, containerator.Mapping{
				Source: strconv.Itoa(basePort + i),
				Target: strconv.Itoa(int(port)),
			})
		}
		ret.Ports = ports
	}
	return &ret
}

func isFileExist(file string) bool {
	if _, err := os.Stat(file); os.IsExist(err) {
		return true
	}
	return false
}

func getEnvFileName(dir string, mode string) string {
	return path.Join(dir, fmt.Sprintf("%s.list", mode))
}

func selectEnvFile(dir string, mode string) string {
	name := getEnvFileName(dir, mode)
	if isFileExist(name) {
		return name
	}
	name = getEnvFileName(dir, "env")
	if isFileExist(name) {
		return name
	}
	return ""
}

func getEnvFileReader(configPath string, mode string) io.Reader {
	dir, _ := filepath.Abs(filepath.Dir(configPath))
	envFileName := selectEnvFile(dir, mode)
	if envFileName == "" {
		log.Println("env file is not found")
		return nil
	}
	data, err := ioutil.ReadFile(envFileName)
	if err != nil {
		log.Printf("failed to read env file: %+v", err)
		return nil
	}
	return strings.NewReader(string(data))
}

func suspendCurrentContainer(container *types.Container, cli client.ContainerAPIClient) error {
	tmpName := containerator.GetContainerName(container) + ".current"
	if err := cli.ContainerRename(context.Background(), container.ID, tmpName); err != nil {
		return err
	}
	if err := cli.ContainerStop(context.Background(), container.ID, nil); err != nil {
		return err
	}
	return nil
}

func resumeCurrentContainer(container *types.Container, name string, cli client.ContainerAPIClient) error {
	if err := cli.ContainerRename(context.Background(), container.ID, name); err != nil {
		return err
	}
	if err := cli.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

func removeCurrentContainer(container *types.Container, cli client.ContainerAPIClient) error {
	return cli.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{})
}

func updateContainer(options *containerator.RunContainerOptions,
	current *types.Container, cli client.ContainerAPIClient) (*types.Container, error) {
	if current != nil {
		if err := suspendCurrentContainer(current, cli); err != nil {
			return nil, err
		}
	}

	nextContainer, err := containerator.RunContainer(cli, options)
	if err != nil {
		if current != nil {
			if err := resumeCurrentContainer(current, containerator.GetContainerName(current), cli); err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	return nextContainer, nil
}

const (
	defaultConfigName = "config.yaml"
)

func run() error {
	var configPathOption string
	flag.StringVar(&configPathOption, "config", defaultConfigName, "configuration file")
	var modeOption string
	flag.StringVar(&modeOption, "mode", "", "mode")
	var forceOption bool
	flag.BoolVar(&forceOption, "force", false, "force container creation")

	flag.Parse()

	workDir, _ := filepath.Abs(filepath.Dir(configPathOption))
	log.Printf("Directory: %s\n", workDir)

	config, err := readConfig(configPathOption)
	if err != nil {
		return err
	}

	mode := selectMode(modeOption, config)
	if mode == "" {
		return fmt.Errorf("'%s' mode is not valid", modeOption)
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	image, err := findImage(cli, config)
	if err != nil {
		return err
	}
	imageName := containerator.GetImageFullName(image)
	log.Printf("Image: %s\n", imageName)

	options := buildContainerOptions(config, imageName, mode)
	if reader := getEnvFileReader(configPathOption, mode); reader != nil {
		options.EnvReader = reader
	}

	log.Printf("Container: %s\n", options.Name)

	currentContainer, err := containerator.FindContainerByName(cli, options.Name)
	if err != nil {
		return err
	}
	if currentContainer != nil && currentContainer.ImageID == image.ID && !forceOption {
		log.Println("Container is already running")
		return nil
	}

	if currentContainer != nil {
		err = suspendCurrentContainer(currentContainer, cli)
		if err != nil {
			return err
		}
	}

	nextContainer, err := containerator.RunContainer(cli, options)
	if err != nil {
		if currentContainer != nil {
			resumeCurrentContainer(currentContainer, options.Name, cli)
		}
		return err
	}

	if currentContainer != nil {
		removeCurrentContainer(currentContainer, cli)
	}

	log.Printf("Container: %s %s\n", containerator.GetContainerName(nextContainer),
		containerator.GetContainerShortID(nextContainer))

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
}

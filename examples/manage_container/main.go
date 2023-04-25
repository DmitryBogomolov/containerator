// Program manage_container shows usage of *containerator* functions that run and remove containers.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/DmitryBogomolov/containerator/manage"
	"github.com/docker/docker/client"
)

func makeConfig(configPath string, containerName string, imageName string) (*manage.Config, error) {
	config, err := manage.ReadConfig(configPath)
	if err != nil && os.IsNotExist(err) {
		config = &manage.Config{}
		err = nil
	}
	if err != nil {
		return nil, err
	}
	if imageName != "" {
		config.ImageName = imageName
	}
	if containerName != "" {
		config.ContainerName = containerName
	}
	if config.ImageName == "" {
		return nil, errors.New("image is not defined")
	}
	return config, err
}

func makeOptions(postfix string, tag string, force bool, remove bool, configPath string) *manage.Options {
	options := manage.Options{
		Postfix: postfix,
		Tag:     tag,
		Force:   force,
		Remove:  remove,
	}
	envFilePath := filepath.Join(filepath.Dir(configPath), fmt.Sprintf("%s.list", postfix))
	if _, err := os.Stat(envFilePath); err == nil {
		options.EnvFilePath = envFilePath
	}
	return &options
}

func displayContainer(container core.Container) string {
	return fmt.Sprintf("%s(%s)", container.Name(), container.ShortID())
}

func run() error {
	var configPath string
	flag.StringVar(&configPath, "config", manage.DefaultConfigName, "configuration file")
	var imageName string
	flag.StringVar(&imageName, "image", "", "image name")
	var containerName string
	flag.StringVar(&containerName, "container", "", "container name")
	var postfix string
	flag.StringVar(&postfix, "postfix", "", "postfix")
	var tag string
	flag.StringVar(&tag, "tag", "", "image tag")
	var remove bool
	flag.BoolVar(&remove, "remove", false, "remove container")
	var force bool
	flag.BoolVar(&force, "force", false, "force container creation")

	flag.Parse()

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	config, err := makeConfig(configPath, containerName, imageName)
	if err != nil {
		return err
	}
	options := makeOptions(postfix, tag, force, remove, configPath)
	container, err := manage.RunContainer(cli, config, options)

	if options.Remove {
		if err != nil {
			return err
		}
		log.Printf("%s: removed\n", displayContainer(container))
		return nil
	}

	if _, ok := err.(*manage.ContainerAlreadyRunningError); ok {
		log.Printf("%s: already running\n", displayContainer(container))
		return nil
	}
	if err != nil {
		return err
	}
	log.Printf("%s: %s\n", displayContainer(container), container.State())

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

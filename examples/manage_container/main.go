// Program manage_container shows usage of *containerator* functions that run and remove containers.
package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/DmitryBogomolov/containerator/manage"
	"github.com/docker/docker/client"
)

func run() error {
	var configPathOption string
	flag.StringVar(&configPathOption, "config", manage.DefaultConfigName, "configuration file")
	var imageRepo string
	flag.StringVar(&imageRepo, "image", "", "image repo")
	var containerName string
	flag.StringVar(&containerName, "container", "", "container name")
	var modeOption string
	flag.StringVar(&modeOption, "mode", "", "mode")
	var tagOption string
	flag.StringVar(&tagOption, "tag", "", "image tag")
	var removeOption bool
	flag.BoolVar(&removeOption, "remove", false, "remove container")
	var forceOption bool
	flag.BoolVar(&forceOption, "force", false, "force container creation")

	flag.Parse()

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	config, err := manage.ReadConfig(configPathOption)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		config = &manage.Config{}
	}
	if imageRepo != "" {
		config.ImageRepo = imageRepo
	}
	if containerName != "" {
		config.ContainerName = containerName
	}

	if config.ImageRepo == "" {
		return errors.New("image repo is not defined")
	}

	options := &manage.Options{
		Mode:   modeOption,
		Tag:    tagOption,
		Force:  forceOption,
		Remove: removeOption,
		GetEnvReader: func(mode string) (io.Reader, error) {
			reader, err := manage.GetEnvFileReader(filepath.Dir(configPathOption), mode)
			if err != nil {
				log.Printf("Failed to load env file: %v\n", err)
			}
			return reader, nil
		},
	}
	container, err := manage.Manage(cli, config, options)

	if options.Remove {
		if _, ok := err.(*manage.NoContainerError); ok {
			log.Println("There is no container")
			return nil
		}
		if err != nil {
			return err
		}
		log.Println("Container is removed")
		return nil
	}

	if _, ok := err.(*manage.ContainerAlreadyRunningError); ok {
		log.Println("Container is already running")
		return nil
	}

	if err != nil {
		return err
	}
	log.Printf("Container: %s %s\n",
		core.GetContainerName(container),
		core.GetContainerShortID(container))

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
}

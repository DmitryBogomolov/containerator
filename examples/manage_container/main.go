package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func getWorkDir(option string) string {
	if _, err := os.Stat(option); os.IsNotExist(err) {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return cwd
	}
	return option
}

func getMode(mode string) string {
	if mode == "" {
		return "dev"
	}
	return strings.ToLower(mode)
}

func selectImage(imageName string, imageRepo string, workDir string, cli client.ImageAPIClient) (*types.ImageSummary, error) {
	if imageName != "" {
		image, err := containerator.FindImageByRepoTag(cli, imageName)
		if image == nil && err == nil {
			err = fmt.Errorf("image %s does not exist", imageName)
		}
		return image, err

	}
	if imageRepo == "" {
		imageRepo = filepath.Base(workDir)
	}
	images, err := containerator.FindImagesByRepo(cli, imageRepo)
	if len(images) == 0 && err == nil {
		err = fmt.Errorf("images %s do not exist", imageRepo)
	}
	if err != nil {
		return nil, err
	}
	return images[0], nil
}

func isFileExist(file string) bool {
	if _, err := os.Stat(file); os.IsExist(err) {
		return true
	}
	return false
}

func getEnvFileName(workDir string, mode string) string {
	return path.Join(workDir, fmt.Sprintf("%s.list", mode))
}

func selectEnvFile(workDir string, mode string) string {
	name := getEnvFileName(workDir, mode)
	if isFileExist(name) {
		return name
	}
	name = getEnvFileName(workDir, "env")
	if isFileExist(name) {
		return name
	}
	return ""
}

func getEnvFileReader(workDir string, mode string) *os.File {
	envFileName := getEnvFileName(workDir, mode)
	if envFileName != "" {
		file, err := os.Open(envFileName)
		if err != nil {
			log.Printf("failed to read env file: %+v", err)
		} else {
			return file
		}
	} else {
		log.Println("env file is not found")
	}
	return nil
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
	if err := cli.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	if err := cli.ContainerRename(context.Background(), container.ID, name); err != nil {
		return err
	}
	return nil
}

func run() error {
	var workDir string
	flag.StringVar(&workDir, "dir", "", "project directory")
	var mode string
	flag.StringVar(&mode, "mode", "", "mode")
	var imageName string
	flag.StringVar(&imageName, "image", "", "full image name")
	var imageRepo string
	flag.StringVar(&imageRepo, "image-repo", "", "image repo")
	var force bool
	flag.BoolVar(&force, "force", false, "force container creation")
	flag.Parse()

	workDir = getWorkDir(workDir)
	log.Printf("Directory: %s\n", workDir)

	mode = getMode(mode)
	log.Printf("Mode: %s\n", mode)

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	image, err := selectImage(imageName, imageRepo, workDir, cli)
	if err != nil {
		return err
	}

	containerName := fmt.Sprintf("%s-%s", containerator.GetImageName(image), mode)

	log.Printf("Image: %s\n", containerator.GetImageFullName(image))
	log.Printf("Container: %s\n", containerName)

	currentContainer, err := containerator.FindContainerByName(cli, containerName)
	if err != nil {
		return err
	}

	if currentContainer != nil && currentContainer.ImageID == image.ID && !force {
		log.Println("Container is already running")
		return nil
	}

	options := containerator.RunContainerOptions{}

	if envReader := getEnvFileReader(workDir, mode); envReader != nil {
		defer envReader.Close()
		options.EnvReader = envReader
	}

	if currentContainer != nil {
		if err := suspendCurrentContainer(currentContainer, cli); err != nil {
			return err
		}
	}

	nextContainer, err := containerator.RunContainer(cli, &options)
	if err != nil {
		if currentContainer != nil {
			if err := resumeCurrentContainer(currentContainer, containerName, cli); err != nil {
				return err
			}
		}
		return err
	}

	log.Printf("Container: %s %s\n", containerName, nextContainer.ID[:8])
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
}

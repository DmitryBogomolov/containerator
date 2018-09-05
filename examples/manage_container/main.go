package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const defaultMode = "dev"

func getMode(mode string) string {
	if mode == "" {
		return defaultMode
	}
	return strings.ToLower(mode)
}

func main() {
	modePtr := flag.String("mode", "", "mode")
	imageNamePtr := flag.String("image", "", "full image name")
	imageRepoPtr := flag.String("image-repo", "", "image repo")
	forcePtr := flag.Bool("force", false, "force container creation")
	flag.Parse()

	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Printf("NewEnvClient: %s", err)
		os.Exit(1)
	}

	mode := getMode(*modePtr)
	fmt.Printf("mode: %s\n", mode)

	var image *types.ImageSummary
	if *imageNamePtr != "" {
		img, err := containerator.FindImageByRepoTag(cli, *imageNamePtr)
		if err != nil {
			fmt.Printf("FindImageByRepoTag: %s\n", err)
		} else {
			image = img
		}
	} else {
		repo := *imageRepoPtr
		if repo == "" {
			cwd, _ := os.Getwd()
			repo = filepath.Base(cwd)
		}
		images, err := containerator.FindImagesByRepo(cli, repo)
		if err != nil {
			fmt.Printf("FindImagesByRepo: %s\n", err)
		} else {
			if len(images) > 0 {
				image = images[0]
			}
		}
	}

	if image == nil {
		fmt.Println("Image is not selected")
		os.Exit(1)
	}

	containerName := fmt.Sprintf("%s-%s", containerator.GetImageName(image), mode)

	fmt.Printf("Image: %s\n", containerator.GetImageFullName(image))
	fmt.Printf("Container: %s\n", containerName)

	currentContainer, err := containerator.FindContainerByName(cli, containerName)
	if err != nil {
		fmt.Printf("FindContainerByName: %s\n", err)
		os.Exit(1)
	}

	if currentContainer != nil && currentContainer.ImageID == image.ID {
		if !*forcePtr {
			fmt.Println("Container is already running")
			os.Exit(0)
		}
	}

	options := containerator.RunContainerOptions{}

	if currentContainer != nil {
		tmpName := containerator.GetContainerName(currentContainer) + ".current"
		err = cli.ContainerRename(context.Background(), currentContainer.ID, tmpName)
		if err != nil {
			fmt.Printf("ContainerRename: %s\n", err)
			os.Exit(1)
		}
		err = cli.ContainerStop(context.Background(), currentContainer.ID, nil)
		if err != nil {
			fmt.Printf("ContainerStop: %s\n", err)
		}
	}

	nextContainer, err := containerator.RunContainer(cli, &options)
	if err != nil {
		fmt.Printf("RunContainer: %s\n", err)
		if currentContainer != nil {
			err = cli.ContainerStart(context.Background(), currentContainer.ID, types.ContainerStartOptions{})
			if err != nil {
				fmt.Printf("ContainerStart: %s\n", err)
				os.Exit(1)
			}
			err = cli.ContainerRename(context.Background(), currentContainer.ID, containerName)
			if err != nil {
				fmt.Printf("ContainerRename: %s\n", err)
				os.Exit(1)
			}
		}
		os.Exit(1)
	}

	fmt.Printf("Container: %s %s\n", containerName, nextContainer.ID[:8])
}

package main

import (
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
	// forcePtr := flag.Bool("force", false, "force container creation")
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
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Getwd: %s\n", err)
			} else {
				repo = filepath.Base(cwd)
			}
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

	container, err := containerator.FindContainerByName(cli, containerName)
	if err != nil {
		fmt.Printf("FindContainerByName: %s\n", err)
		os.Exit(1)
	}

	fmt.Print(container)
}

// Program find_image shows usage of *containerator* functions that find docker images.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/docker/docker/client"
)

func run() error {
	var id string
	flag.StringVar(&id, "id", "", "id")
	var repoTag string
	flag.StringVar(&repoTag, "repo-tag", "", "repo tag")
	var repo string
	flag.StringVar(&repo, "repo", "", "repo")

	flag.Parse()

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	if id != "" {
		image, err := core.FindImageByID(cli, id)
		if err != nil {
			return err
		}
		if image == nil {
			fmt.Println("Image is not found.")
		} else {
			fmt.Printf("Image: %s\n", image.FullName())
		}
	} else if repoTag != "" {
		image, err := core.FindImageByName(cli, repoTag)
		if err != nil {
			return err
		}
		if image == nil {
			fmt.Println("Image is not found.")
		} else {
			fmt.Printf("Image: %s\n", image.FullName())
		}
	} else if repo != "" {
		images, err := core.FindAllImagesByName(cli, repo)
		if err != nil {
			return err
		}
		fmt.Println("Images:")
		tags := core.TransformSlice(images, func(image core.Image) string {
			return image.Tag()
		})
		for i, image := range images {
			fmt.Printf("  %s %s\n", image.FullName(), tags[i])
		}
	} else {
		flag.Usage()
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

// Program find_image shows usage of *containerator* functions that find docker images.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DmitryBogomolov/containerator"
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
		image, err := containerator.FindImageByID(cli, id)
		if err != nil {
			return err
		}
		if image == nil {
			fmt.Println("Image is not found.")
		} else {
			fmt.Printf("Image: %s\n", containerator.GetImageFullName(image))
		}
	} else if repoTag != "" {
		image, err := containerator.FindImageByRepoTag(cli, repoTag)
		if err != nil {
			return err
		}
		if image == nil {
			fmt.Println("Image is not found.")
		} else {
			fmt.Printf("Image: %s\n", containerator.GetImageFullName(image))
		}
	} else if repo != "" {
		images, err := containerator.FindImagesByRepo(cli, repo)
		if err != nil {
			return err
		}
		fmt.Println("Images:")
		tags := containerator.GetImagesTags(images)
		for i, image := range images {
			fmt.Printf("  %s %s\n", containerator.GetImageFullName(image), tags[i])
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

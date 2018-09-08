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
			panic(err)
		}
		if image == nil {
			fmt.Println("Image is not found.")
		} else {
			fmt.Printf("Image: %s\n", containerator.GetImageFullName(image))
		}
	}
	if repoTag != "" {
		image, err := containerator.FindImageByRepoTag(cli, repoTag)
		if err != nil {
			panic(err)
		}
		if image == nil {
			fmt.Println("Image is not found.")
		} else {
			fmt.Printf("Image: %s\n", containerator.GetImageFullName(image))
		}
	}
	if repo != "" {
		images, err := containerator.FindImagesByRepo(cli, repo)
		if err != nil {
			panic(err)
		}
		fmt.Println("Images:")
		for _, image := range images {
			fmt.Printf("  %s\n", containerator.GetImageFullName(image))
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

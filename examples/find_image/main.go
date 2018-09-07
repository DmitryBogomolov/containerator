package main

import (
	"flag"
	"fmt"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/client"
)

// TODO: Use flag.StringVar
func main() {
	id := flag.String("id", "", "id")
	repoTag := flag.String("repo-tag", "", "repo tag")
	repo := flag.String("repo", "", "repo")
	flag.Parse()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	if *id != "" {
		image, err := containerator.FindImageByID(cli, *id)
		if err != nil {
			panic(err)
		}
		if image == nil {
			fmt.Println("Not found")
		} else {
			fmt.Printf("Image: %s\n", containerator.GetImageName(image))
		}
	}
	if *repoTag != "" {
		image, err := containerator.FindImageByRepoTag(cli, *repoTag)
		if err != nil {
			panic(err)
		}
		if image == nil {
			fmt.Println("Not found")
		} else {
			fmt.Printf("Image: %s\n", containerator.GetImageName(image))
		}
	}
	if *repo != "" {
		images, err := containerator.FindImagesByRepo(cli, *repo)
		if err != nil {
			panic(err)
		}
		fmt.Println("Images:")
		for _, image := range images {
			fmt.Printf("  %s\n", containerator.GetImageName(image))
		}
	}
}

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
	name := flag.String("name", "", "name")
	imageID := flag.String("image-id", "", "image id")
	flag.Parse()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	if *id != "" {
		container, err := containerator.FindContainerByID(cli, *id)
		if err != nil {
			panic(err)
		}
		if container == nil {
			fmt.Println("Not found")
		} else {
			fmt.Printf("Container: %s\n", containerator.GetContainerName(container))
		}
	}
	if *name != "" {
		container, err := containerator.FindContainerByName(cli, *name)
		if err != nil {
			panic(err)
		}
		if container == nil {
			fmt.Println("Not found")
		} else {
			fmt.Printf("Container: %s\n", containerator.GetContainerName(container))
		}
	}
	if *imageID != "" {
		containers, err := containerator.FindContainersByImageID(cli, *imageID)
		if err != nil {
			panic(err)
		}
		fmt.Println("Containers:")
		for _, container := range containers {
			fmt.Printf("  %s\n", containerator.GetContainerName(container))
		}
	}
}

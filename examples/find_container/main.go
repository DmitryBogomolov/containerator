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
	var name string
	flag.StringVar(&name, "name", "", "name")
	var imageID string
	flag.StringVar(&imageID, "image-id", "", "image id")

	flag.Parse()

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	if id != "" {
		container, err := containerator.FindContainerByID(cli, id)
		if err != nil {
			panic(err)
		}
		if container == nil {
			fmt.Println("Container is not found.")
		} else {
			fmt.Printf("Container: %s\n", containerator.GetContainerName(container))
		}
	}
	if name != "" {
		container, err := containerator.FindContainerByName(cli, name)
		if err != nil {
			panic(err)
		}
		if container == nil {
			fmt.Println("Container is not found.")
		} else {
			fmt.Printf("Container: %s\n", containerator.GetContainerName(container))
		}
	}
	if imageID != "" {
		containers, err := containerator.FindContainersByImageID(cli, imageID)
		if err != nil {
			panic(err)
		}
		fmt.Println("Containers:")
		for _, container := range containers {
			fmt.Printf("  %s\n", containerator.GetContainerName(container))
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

// Program find_container shows usage of *containerator* functions that find docker containers.
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
			return err
		}
		if container == nil {
			fmt.Println("Container is not found.")
		} else {
			fmt.Printf("Container: %s\n", containerator.GetContainerName(container))
		}
	} else if name != "" {
		container, err := containerator.FindContainerByName(cli, name)
		if err != nil {
			return err
		}
		if container == nil {
			fmt.Println("Container is not found.")
		} else {
			fmt.Printf("Container: %s\n", containerator.GetContainerName(container))
		}
	} else if imageID != "" {
		containers, err := containerator.FindContainersByImageID(cli, imageID)
		if err != nil {
			return err
		}
		fmt.Println("Containers:")
		for _, container := range containers {
			fmt.Printf("  %s\n", containerator.GetContainerName(container))
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

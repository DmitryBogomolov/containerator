// Program find_container shows usage of *containerator* functions that find docker containers.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/docker/docker/client"
)

func displayContainer(container core.Container) string {
	return fmt.Sprintf("%s (%s)", container.Name(), container.ShortID())
}

func findContainerByID(cli *client.Client, id string) error {
	container, err := core.FindContainerByShortID(cli, id)
	if _, ok := err.(*core.ContainerNotFoundError); ok {
		fmt.Println("Container not found")
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Println(displayContainer(container))
	return nil
}

func findContainerByName(cli *client.Client, name string) error {
	container, err := core.FindContainerByName(cli, name)
	if _, ok := err.(*core.ContainerNotFoundError); ok {
		fmt.Println("Container not found")
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Println(displayContainer(container))
	return nil
}

func findContainersByImageID(cli *client.Client, imageID string) error {
	containers, err := core.FindContainersByImageID(cli, imageID)
	if err != nil {
		return err
	}
	for _, container := range containers {
		fmt.Println(displayContainer(container))
	}
	return nil
}

func listAllContainers(cli *client.Client) error {
	containerIDs, err := core.ListAllContainerIDs(cli)
	if err != nil {
		return err
	}
	for _, id := range containerIDs {
		fmt.Println(id)
	}
	return nil
}

func run() error {
	var id string
	flag.StringVar(&id, "id", "", "container by d")
	var name string
	flag.StringVar(&name, "name", "", "container by name")
	var imageID string
	flag.StringVar(&imageID, "image-id", "", "containers by image id")
	var listAll bool
	flag.BoolVar(&listAll, "list-all", false, "list all containers")

	flag.Parse()

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	if id != "" {
		return findContainerByID(cli, id)
	} else if name != "" {
		return findContainerByName(cli, name)
	} else if imageID != "" {
		return findContainersByImageID(cli, imageID)
	} else if listAll {
		return listAllContainers(cli)
	} else {
		flag.Usage()
		return nil
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

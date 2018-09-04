package main

import (
	"fmt"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	container, err := containerator.FindContainer(cli, containerator.FindContainerOptions{
		ID: "test",
	})
	if err != nil {
		fmt.Println(err)
	} else if container == nil {
		fmt.Println("container is not found")
	} else {
		fmt.Printf("%s %s %s %s\n", container.ID, container.Name, container.Image, container.State)
	}
}

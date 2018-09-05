package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/client"
)

func main() {
	imagePtr := flag.String("image", "", "image name")
	namePtr := flag.String("name", "", "container name")
	volumePtr := flag.String("volume", "", "volume")
	portPtr := flag.String("port", "", "port")
	envPtr := flag.String("env", "", "environment")
	flag.Parse()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	options := &containerator.RunContainerOptions{
		Image: *imagePtr,
		Name:  *namePtr,
	}
	if *volumePtr != "" {
		list := strings.Split(*volumePtr, ":")
		options.Volumes = []containerator.Mapping{
			containerator.Mapping{Source: list[0], Target: list[1]},
		}
	}
	if *portPtr != "" {
		list := strings.Split(*portPtr, ":")
		options.Ports = []containerator.Mapping{
			containerator.Mapping{Source: list[0], Target: list[1]},
		}
	}
	if *envPtr != "" {
		list := strings.Split(*envPtr, "=")
		mapping := containerator.Mapping{Source: list[0]}
		if len(list) > 1 {
			mapping.Target = list[1]
		}
		options.Env = []containerator.Mapping{mapping}
	}

	container, err := containerator.RunContainer(cli, options)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s %s %s\n", container.ImageID[7:15], container.ID[7:15],
		containerator.GetContainerName(container), container.State)
}

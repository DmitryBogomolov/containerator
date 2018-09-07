package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/client"
)

// TODO: Use flag.StringVar
// TODO: Use custom flag for `port`, `volume`, `env`
// TODO: Add `env-file`
func main() {
	imagePtr := flag.String("image", "", "image name")
	namePtr := flag.String("name", "", "container name")
	volumePtr := flag.String("volume", "", "volume")
	portPtr := flag.String("port", "", "port")
	envPtr := flag.String("env", "", "environment")
	var restart string
	flag.StringVar(&restart, "restart", "", "restart policy")
	var network string
	flag.StringVar(&network, "network", "", "network")

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
	options.RestartPolicy = containerator.RestartPolicy(restart)
	options.Network = network

	container, err := containerator.RunContainer(cli, options)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s %s %s\n", container.ImageID[7:15], container.ID[7:15],
		containerator.GetContainerName(container), container.State)
}

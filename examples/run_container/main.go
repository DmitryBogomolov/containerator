// Program run_container shows usage of *containerator.RunContainer* function.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/docker/docker/client"
)

func run() error {
	var imageName string
	flag.StringVar(&imageName, "image", "", "image name")
	var containerName string
	flag.StringVar(&containerName, "name", "", "container name")
	volumes := core.NewMappingListFlag(":", false)
	flag.Var(volumes, "volume", "volume")
	ports := core.NewMappingListFlag(":", false)
	flag.Var(ports, "port", "port")
	env := core.NewMappingListFlag("=", true)
	flag.Var(env, "env", "environment")
	var restart string
	flag.StringVar(&restart, "restart", "", "restart policy")
	var network string
	flag.StringVar(&network, "network", "", "network")

	flag.Parse()

	if imageName == "" || containerName == "" {
		flag.Usage()
		return nil
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	options := core.RunContainerOptions{
		Image:         imageName,
		Name:          containerName,
		Network:       network,
		RestartPolicy: core.RestartPolicy(restart),
		Ports:         ports.Get(),
		Volumes:       volumes.Get(),
		Env:           env.Get(),
	}
	container, err := core.RunContainer(cli, &options)
	if err != nil {
		return err
	}
	fmt.Printf("%s/%s (%s): %s\n", container.Name(), container.ShortID(), imageName, container.State())

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

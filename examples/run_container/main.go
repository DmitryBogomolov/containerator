// Program run_container shows usage of *containerator.RunContainer* function.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/client"
)

func run() error {
	var imageName string
	flag.StringVar(&imageName, "image", "", "image name")
	var containerName string
	flag.StringVar(&containerName, "name", "", "container name")
	volumes := containerator.NewMappingListVar(":", false)
	flag.Var(volumes, "volume", "volume")
	ports := containerator.NewMappingListVar(":", false)
	flag.Var(ports, "port", "port")
	env := containerator.NewMappingListVar("=", true)
	flag.Var(env, "env", "environment")
	var envFile string
	flag.StringVar(&envFile, "env-file", "", "env file")
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

	options := &containerator.RunContainerOptions{
		Image: imageName,
		Name:  containerName,
	}
	options.Volumes = volumes.Get()
	options.Ports = ports.Get()
	options.Env = env.Get()
	options.RestartPolicy = containerator.RestartPolicy(restart)
	options.Network = network

	container, err := containerator.RunContainer(cli, options)
	if err != nil {
		return err
	}
	fmt.Printf("%s %s %s %s\n",
		imageName,
		containerator.GetContainerShortID(container),
		containerator.GetContainerName(container),
		container.State)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

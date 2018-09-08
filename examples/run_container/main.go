package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/client"
)

type mapping struct {
	list []containerator.Mapping
}

func (m *mapping) String() string {
	return fmt.Sprintf("%v", m.list)
}

func (m *mapping) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) < 2 {
		return errors.New("not a pair")
	}
	m.list = append(m.list, containerator.Mapping{Source: parts[0], Target: parts[1]})
	return nil
}

type envMapping struct {
	mapping
}

func (m *envMapping) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	obj := containerator.Mapping{Source: parts[0]}
	if len(parts) > 1 {
		obj.Target = parts[1]
	}
	m.list = append(m.list, obj)
	return nil
}

func run() error {
	var imageName string
	flag.StringVar(&imageName, "image", "", "image name")
	var containerName string
	flag.StringVar(&containerName, "name", "", "container name")
	var volumes mapping
	flag.Var(&volumes, "volume", "volume")
	var ports mapping
	flag.Var(&ports, "port", "port")
	var env envMapping
	flag.Var(&env, "env", "environment")
	var envFile string
	flag.StringVar(&envFile, "env-file", "", "env file")
	var restart string
	flag.StringVar(&restart, "restart", "", "restart policy")
	var network string
	flag.StringVar(&network, "network", "", "network")

	flag.Parse()

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	options := &containerator.RunContainerOptions{
		Image: imageName,
		Name:  containerName,
	}
	if len(volumes.list) > 0 {
		options.Volumes = volumes.list
	}
	if len(ports.list) > 0 {
		options.Ports = ports.list
	}
	if len(env.mapping.list) > 0 {
		options.Env = env.list
	}
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

package containerator

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func cliContainerCreate(cli *client.Client, containerName string, config *container.Config, hostConfig *container.HostConfig) (container.ContainerCreateCreatedBody, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return cli.ContainerCreate(ctx, config, hostConfig, nil, containerName)
}

func cliContainerStart(cli *client.Client, containerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

type containerOptions struct {
	Image   string
	Name    string
	Volumes map[string]string
	Ports   map[int]int
}

func buildPortBindings(options map[int]int) (nat.PortSet, nat.PortMap) {
	var ports nat.PortSet
	var bindings nat.PortMap
	dummy := struct{}{}
	for from, to := range options {
		key := nat.Port(fmt.Sprintf("%d/tcp", from))
		ports[key] = dummy
		val := nat.PortBinding{
			HostPort: fmt.Sprintf("%d", to),
			HostIP:   "0.0.0.0",
		}
		bindings[key] = []nat.PortBinding{val}
	}
	return ports, bindings
}

func buildVolumes(options map[string]string) []string {
	var volumes []string
	for from, to := range options {
		volumes = append(volumes, fmt.Sprintf("%s:%s", from, to))
	}
	return volumes
}

func runContainer(cli *client.Client, options *containerOptions) (string, error) {
	config := container.Config{}
	hostConfig := container.HostConfig{}
	config.ExposedPorts, hostConfig.PortBindings = buildPortBindings(options.Ports)
	hostConfig.Binds = buildVolumes(options.Volumes)
	body, err := cliContainerCreate(cli, options.Name, &config, &hostConfig)
	if err != nil {
		return "", err
	}
	err = cliContainerStart(cli, body.ID)
	if err != nil {
		removeContainer(cli, body.ID)
		return "", err
	}
	return body.ID, nil
}

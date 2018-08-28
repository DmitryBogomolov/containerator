package containerator

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

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

func runContainer(cli client.ContainerAPIClient, options *containerOptions) (containerInfo, error) {
	config := container.Config{}
	hostConfig := container.HostConfig{}
	config.ExposedPorts, hostConfig.PortBindings = buildPortBindings(options.Ports)
	hostConfig.Binds = buildVolumes(options.Volumes)
	body, err := cliContainerCreate(cli, &config, &hostConfig, options.Name)
	emptyInfo := containerInfo{}
	if err != nil {
		return emptyInfo, err
	}
	err = cliContainerStart(cli, body.ID)
	if err != nil {
		removeContainer(cli, body.ID)
		return emptyInfo, err
	}
	return inspectContainer(cli, body.ID)
}

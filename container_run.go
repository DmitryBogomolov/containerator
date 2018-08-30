package containerator

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// ContainerOptions contains options for container.
type ContainerOptions struct {
	Image   string
	Name    string
	Volumes map[string]string
	Ports   map[int]int
}

func buildPortBindings(options map[int]int) (nat.PortSet, nat.PortMap) {
	if len(options) == 0 {
		return nil, nil
	}
	ports := make(nat.PortSet)
	bindings := make(nat.PortMap)
	var dummy struct{}
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

// RunContainer creates and starts container.
func RunContainer(cli client.ContainerAPIClient, options *ContainerOptions) (*ContainerInfo, error) {
	config := container.Config{}
	hostConfig := container.HostConfig{}
	config.Image = options.Image
	config.ExposedPorts, hostConfig.PortBindings = buildPortBindings(options.Ports)
	hostConfig.Binds = buildVolumes(options.Volumes)
	body, err := cliContainerCreate(cli, &config, &hostConfig, options.Name)
	if err != nil {
		return nil, err
	}
	err = cliContainerStart(cli, body.ID)
	if err != nil {
		removeContainer(cli, body.ID)
		return nil, err
	}
	return inspectContainer(cli, body.ID)
}

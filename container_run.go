package containerator

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// RunContainerOptions contains options for container.
type RunContainerOptions struct {
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

func buildMounts(options map[string]string) []mount.Mount {
	var mounts []mount.Mount
	for from, to := range options {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: from,
			Target: to,
		})
	}
	return mounts
}

// RunContainer creates and starts container.
func RunContainer(cli client.ContainerAPIClient, options *RunContainerOptions) (*types.Container, error) {
	config := container.Config{}
	hostConfig := container.HostConfig{}
	config.Image = options.Image
	config.ExposedPorts, hostConfig.PortBindings = buildPortBindings(options.Ports)
	hostConfig.Mounts = buildMounts(options.Volumes)
	body, err := cliContainerCreate(cli, &config, &hostConfig, options.Name)
	if err != nil {
		return nil, err
	}
	err = cliContainerStart(cli, body.ID)
	if err != nil {
		cliContainerRemove(cli, body.ID)
		return nil, err
	}
	return FindContainerByID(cli, body.ID)
}

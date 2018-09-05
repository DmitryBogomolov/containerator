package containerator

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Mapping defines source-target pair.
type Mapping struct {
	Source string
	Target string
}

// RunContainerOptions contains options for container.
type RunContainerOptions struct {
	Image   string
	Name    string
	Volumes []Mapping
	Ports   []Mapping
	Env     []Mapping
}

func buildPortBindings(options []Mapping) (nat.PortSet, nat.PortMap) {
	if len(options) == 0 {
		return nil, nil
	}
	ports := make(nat.PortSet)
	bindings := make(nat.PortMap)
	var dummy struct{}
	for _, mapping := range options {
		key := nat.Port(fmt.Sprintf("%s/tcp", mapping.Target))
		ports[key] = dummy
		val := nat.PortBinding{
			HostPort: fmt.Sprintf("%s", mapping.Source),
			HostIP:   "0.0.0.0",
		}
		bindings[key] = []nat.PortBinding{val}
	}
	return ports, bindings
}

func buildMounts(options []Mapping) []mount.Mount {
	var mounts []mount.Mount
	for _, mapping := range options {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: mapping.Source,
			Target: mapping.Target,
		})
	}
	return mounts
}

func buildEnvironment(options []Mapping) []string {
	var env []string
	for _, mapping := range options {
		name := mapping.Source
		value := mapping.Target
		if value == "" {
			value = os.Getenv(name)
		}
		env = append(env, fmt.Sprintf("%s=%s", name, value))
	}
	return env
}

// RunContainer creates and starts container.
func RunContainer(cli client.ContainerAPIClient, options *RunContainerOptions) (*types.Container, error) {
	config := container.Config{}
	hostConfig := container.HostConfig{}
	config.Image = options.Image
	config.ExposedPorts, hostConfig.PortBindings = buildPortBindings(options.Ports)
	config.Env = buildEnvironment(options.Env)
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

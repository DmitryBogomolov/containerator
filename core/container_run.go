package core

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// RunContainerOptions contains options used to create and start container.
type RunContainerOptions struct {
	Image         string                      `json:"image,omitempty" yaml:",omitempty"`          // Image name; required
	Name          string                      `json:"name,omitempty" yaml:",omitempty"`           // Container name
	Volumes       []Mapping                   `json:"volumes,omitempty" yaml:",omitempty"`        // List of volume mappings
	Ports         []Mapping                   `json:"ports,omitempty" yaml:",omitempty"`          // List of port mappings
	Env           []Mapping                   `json:"env,omitempty" yaml:",omitempty"`            // List of environment variables; has priority over `EnvReader`
	RestartPolicy container.RestartPolicyMode `json:"restart,omitempty" yaml:"restart,omitempty"` // Container restart policy
	Network       string                      `json:"network,omitempty" yaml:",omitempty"`        // Container network
}

func buildPortBindings(mappings []Mapping) (nat.PortSet, nat.PortMap) {
	if len(mappings) == 0 {
		return nil, nil
	}
	ports := make([]string, len(mappings))
	for i, mapping := range mappings {
		ports[i] = fmt.Sprintf("0.0.0.0:%s:%s/tcp", mapping.Source, mapping.Target)
	}
	exposedPorts, portBindings, _ := nat.ParsePortSpecs(ports)
	return exposedPorts, portBindings
}

func buildMounts(mappings []Mapping) []mount.Mount {
	if len(mappings) == 0 {
		return nil
	}
	result := make([]mount.Mount, 0, len(mappings))
	for _, mapping := range mappings {
		result = append(result, mount.Mount{
			Type:   mount.TypeBind,
			Source: mapping.Source,
			Target: mapping.Target,
		})
	}
	return result
}

func buildEnvironment(mappings []Mapping) []string {
	if len(mappings) == 0 {
		return nil
	}
	result := make([]string, 0, len(mappings))
	for _, mapping := range mappings {
		name := mapping.Source
		value := mapping.Target
		if value == "" {
			value = os.Getenv(name)
		}
		result = append(result, fmt.Sprintf("%s=%s", name, value))
	}
	return result
}

/*
RunContainer creates and starts container.

Roughly duplicates `docker run` command.
If created container fails at start it is removed.

	RunContainer(cli, &RunContainerOptions{
		Image: "my-image:1",
		Name: "my-container-1",
		RestartPolicy: RestartAlways,
		Network: "my-network-1",
		Volumes: []Mapping{
			{"/tmp", "/usr/app"},
		},
		Ports: []Mapping{
			{"50001", "3000"},
		},
		Env: []Mapping{
			{"A", "1"},
		},
	}) -> &container
*/
func RunContainer(cli client.ContainerAPIClient, options *RunContainerOptions) (Container, error) {
	config := container.Config{}
	hostConfig := container.HostConfig{}

	config.Image = options.Image
	config.ExposedPorts, hostConfig.PortBindings = buildPortBindings(options.Ports)
	config.Env = buildEnvironment(options.Env)
	hostConfig.Mounts = buildMounts(options.Volumes)
	if options.RestartPolicy != "" {
		hostConfig.RestartPolicy.Name = options.RestartPolicy
	}
	if options.Network != "" {
		hostConfig.NetworkMode = container.NetworkMode(options.Network)
	}

	body, err := cliContainerCreate(cli, &config, &hostConfig, options.Name)
	if err != nil {
		return nil, err
	}
	containerID := body.ID
	err = cliContainerStart(cli, containerID)
	if err != nil {
		cliContainerRemove(cli, containerID)
		return nil, err
	}
	return FindContainerByID(cli, containerID)
}

package core

import (
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// RunContainerOptions contains options used to create and start container.
type RunContainerOptions struct {
	Image         string        `json:"image,omitempty" yaml:",omitempty"`          // Image name; required
	Name          string        `json:"name,omitempty" yaml:",omitempty"`           // Container name
	Volumes       []Mapping     `json:"volumes,omitempty" yaml:",omitempty"`        // List of volume mappings
	Ports         []Mapping     `json:"ports,omitempty" yaml:",omitempty"`          // List of port mappings
	Env           []Mapping     `json:"env,omitempty" yaml:",omitempty"`            // List of environment variables; has priority over `EnvReader`
	EnvReader     io.Reader     `json:"-" yaml:"-"`                                 // Environment variables in yaml format
	RestartPolicy RestartPolicy `json:"restart,omitempty" yaml:"restart,omitempty"` // Container restart policy
	Network       string        `json:"network,omitempty" yaml:",omitempty"`        // Container network
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

func buildEnvironment(env []Mapping, envReader io.Reader) ([]string, error) {
	var ret []string
	if envReader != nil {
		obj, err := godotenv.Parse(envReader)
		if err != nil {
			return nil, err
		}
		for name, value := range obj {
			ret = append(ret, fmt.Sprintf("%s=%s", name, value))
		}
	}
	for _, mapping := range env {
		name := mapping.Source
		value := mapping.Target
		if value == "" {
			value = os.Getenv(name)
		}
		ret = append(ret, fmt.Sprintf("%s=%s", name, value))
	}
	return ret, nil
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
func RunContainer(cli client.ContainerAPIClient, options *RunContainerOptions) (*types.Container, error) {
	config := container.Config{}
	hostConfig := container.HostConfig{}

	config.Image = options.Image
	config.ExposedPorts, hostConfig.PortBindings = buildPortBindings(options.Ports)
	env, err := buildEnvironment(options.Env, options.EnvReader)
	if err != nil {
		return nil, err
	}
	config.Env = env
	hostConfig.Mounts = buildMounts(options.Volumes)
	if options.RestartPolicy != "" {
		hostConfig.RestartPolicy.Name = string(options.RestartPolicy)
	}
	if options.Network != "" {
		hostConfig.NetworkMode = container.NetworkMode(options.Network)
	}

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

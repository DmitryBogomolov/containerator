package containerator

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

// Mapping defines source-target pair.
type Mapping struct {
	Source string
	Target string
}

// RestartPolicy defines container restart policy.
type RestartPolicy string

// RestartPolicy values.
const (
	RestartOnFailure     RestartPolicy = "on-failure"
	RestartUnlessStopped RestartPolicy = "unless-stopped"
	RestartAlways        RestartPolicy = "always"
)

// RunContainerOptions contains options for container.
type RunContainerOptions struct {
	Image         string
	Name          string
	Volumes       []Mapping
	Ports         []Mapping
	Env           []Mapping
	EnvReader     io.Reader
	RestartPolicy RestartPolicy
	Network       string
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

// RunContainer creates and starts container.
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

package core

import (
	"encoding/json"
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

// Mapping stores key-value pair. Used for volumes, ports, environment variables.
type Mapping struct {
	Source string
	Target string
}

func (mapping Mapping) toMap() map[string]string {
	ret := map[string]string{}
	ret[mapping.Source] = mapping.Target
	return ret
}

func (mapping *Mapping) fromMap(data map[string]string) {
	for key, val := range data {
		mapping.Source = key
		mapping.Target = val
	}
}

// MarshalJSON implements `json.Marshaler` interface.
func (mapping Mapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(mapping.toMap())
}

// UnmarshalJSON implements `json.Unmarshaler` interface.
func (mapping *Mapping) UnmarshalJSON(data []byte) error {
	var tmp map[string]string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	mapping.fromMap(tmp)
	return nil
}

// MarshalYAML implements `yaml.Marshaler` interface.
func (mapping Mapping) MarshalYAML() (interface{}, error) {
	return mapping.toMap(), nil
}

// UnmarshalYAML implements `yaml.Unmarshaler` interface.
func (mapping *Mapping) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp map[string]string
	if err := unmarshal(&tmp); err != nil {
		return err
	}
	mapping.fromMap(tmp)
	return nil
}

// RestartPolicy defines container restart policy.
type RestartPolicy string

// RestartPolicy values.
const (
	RestartOnFailure     RestartPolicy = "on-failure"
	RestartUnlessStopped RestartPolicy = "unless-stopped"
	RestartAlways        RestartPolicy = "always"
)

// RunContainerOptions contains options used to create and start container.
//
//`Image` is required, rest are optional.
//`Env` has priority over `EnvReader`.
//`EnvReader` expects *yaml*.
type RunContainerOptions struct {
	Image         string        `json:"image,omitempty" yaml:",omitempty"`
	Name          string        `json:"name,omitempty" yaml:",omitempty"`
	Volumes       []Mapping     `json:"volumes,omitempty" yaml:",omitempty"`
	Ports         []Mapping     `json:"ports,omitempty" yaml:",omitempty"`
	Env           []Mapping     `json:"env,omitempty" yaml:",omitempty"`
	EnvReader     io.Reader     `json:"-" yaml:"-"`
	RestartPolicy RestartPolicy `json:"restart,omitempty" yaml:"restart,omitempty"`
	Network       string        `json:"network,omitempty" yaml:",omitempty"`
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
		Image: "my-image:1", // or "sha256:<guid>"
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

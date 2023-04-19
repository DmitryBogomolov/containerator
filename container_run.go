package containerator

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Mapping stores key-value pair. Used for volumes, ports, envrionment variables.
type Mapping struct {
	Source string
	Target string
}

// MarshalJSON implements `json.Marshaler` interface.
func (m Mapping) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf(`{"%s":"%s"}`, m.Source, m.Target)
	return []byte(str), nil
}

// UnmarshalJSON implements `json.Unmarshaler` interface.
func (m *Mapping) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(str), "{"), "}")
	parts := strings.SplitN(str, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("not valid JSON Mapping: %s", data)
	}
	m.Source = strings.Trim(strings.TrimSpace(parts[0]), `"`)
	m.Target = strings.Trim(strings.TrimSpace(parts[1]), `"`)
	return nil
}

// MarshalYAML implements `yaml.Marshaler` interface.
func (m Mapping) MarshalYAML() (interface{}, error) {
	ret := make(map[string]string)
	ret[m.Source] = m.Target
	return ret, nil
}

// UnmarshalYAML implements `yaml.Unmarshaler` interface.
func (m *Mapping) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tmp := make(map[string]string)
	err := unmarshal(tmp)
	if err != nil {
		return err
	}
	for key, val := range tmp {
		m.Source = key
		m.Target = val
		break
	}
	return nil
}

// NewMappingListFromMap creates list of Mapping instances from dictionary.
func NewMappingListFromMap(data map[string]string) []Mapping {
	ret := make([]Mapping, 0, len(data))
	for source, target := range data {
		ret = append(ret, Mapping{Source: source, Target: target})
	}
	return ret
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

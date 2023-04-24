package manage

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestConfig(t *testing.T) {
	t.Run("YAMLMarshal", func(t *testing.T) {
		config := Config{
			ImageName:     "test-image",
			ContainerName: "test-container",
			Network:       "test-network",
			Modes:         []string{"dev", "test", "prod"},
			BasePort:      1001,
			PortOffset:    42,
			Ports:         []int{11, 12, 13},
			Volumes: []core.Mapping{
				{Source: "/a", Target: "/b"},
			},
			Env: []core.Mapping{
				{Source: "A", Target: "1"},
				{Source: "B", Target: "2"},
			},
		}

		bytes, err := yaml.Marshal(config)
		data := string(bytes)

		assert.NoError(t, err)
		expected := strings.Join([]string{
			"image_name: test-image",
			"container_name: test-container",
			"network: test-network",
			"base_port: 1001",
			"port_offset: 42",
			"ports:",
			"- 11",
			"- 12",
			"- 13",
			"volumes:",
			"- /a: /b",
			"env:",
			`- A: "1"`,
			`- B: "2"`,
			"modes:",
			"- dev",
			"- test",
			"- prod",
			"",
		}, "\n")
		assert.Equal(t, expected, data)
	})

	t.Run("YAMLUnmarshal", func(t *testing.T) {
		data := strings.Join([]string{
			"image_name: test-image",
			"container_name: test-container",
			"env:",
			`- A: "1"`,
			`- B: "2"`,
			"modes:",
			"- dev",
			"- prod",
		}, "\n")
		var config Config

		err := yaml.Unmarshal([]byte(data), &config)

		assert.NoError(t, err)
		assert.Equal(t, Config{
			ImageName:     "test-image",
			ContainerName: "test-container",
			Env: []core.Mapping{
				{Source: "A", Target: "1"},
				{Source: "B", Target: "2"},
			},
			Modes: []string{
				"dev",
				"prod",
			},
		}, config)
	})
}

func TestReadConfig(t *testing.T) {
	t.Run("Read file", func(t *testing.T) {
		ioutil.WriteFile("test.yaml", []byte("image_name: my-image\nmodes: ['a', 'b']\n"), os.ModePerm)
		defer os.Remove("test.yaml")

		config, err := ReadConfig("test.yaml")

		assert.NoError(t, err, "error")
		assert.Equal(t, &Config{
			ImageName: "my-image",
			Modes:     []string{"a", "b"},
		}, config, "config")
	})

	t.Run("No file", func(t *testing.T) {
		config, err := ReadConfig("test.yaml")

		assert.Error(t, err, "error")
		assert.Equal(t, (err.(*os.PathError)).Path, "test.yaml", "error data")
		assert.Equal(t, (*Config)(nil), config, "config")
	})
}

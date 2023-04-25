package manage

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
			Ports: []core.Mapping{
				{Source: "5001", Target: "11"},
				{Source: "5002", Target: "12"},
			},
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
			"ports:",
			`- "5001": "11"`,
			`- "5002": "12"`,
			"volumes:",
			"- /a: /b",
			"env:",
			`- A: "1"`,
			`- B: "2"`,
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
		}, config)
	})
}

func TestReadConfig(t *testing.T) {
	t.Run("Read file", func(t *testing.T) {
		testContent := strings.Join([]string{
			"image_name: test-image",
			"container_name: test-container",
		}, "\n")
		testFile := "test.yaml"
		ioutil.WriteFile(testFile, []byte(testContent), os.ModePerm)
		defer os.Remove(testFile)

		config, err := ReadConfig(testFile)

		assert.NoError(t, err, "error")
		assert.Equal(t, &Config{
			ImageName:     "test-image",
			ContainerName: "test-container",
		}, config, "config")
	})

	t.Run("No file", func(t *testing.T) {
		config, err := ReadConfig("test.yaml")

		assert.Error(t, err, "error")
		assert.Equal(t, (err.(*os.PathError)).Path, "test.yaml", "error data")
		assert.Equal(t, (*Config)(nil), config, "config")
	})

	t.Run("Process mounts", func(t *testing.T) {
		testContent := strings.Join([]string{
			"volumes:",
			"- /a/b: /dir1",
			"- ./a/b: /dir2",
		}, "\n")
		testFile := "test.yaml"
		ioutil.WriteFile(testFile, []byte(testContent), os.ModePerm)
		defer os.Remove(testFile)

		config, err := ReadConfig(testFile)

		assert.NoError(t, err, "error")
		curDir, _ := filepath.Abs(".")
		assert.Equal(t, &Config{
			Volumes: []core.Mapping{
				{Source: "/a/b", Target: "/dir1"},
				{Source: filepath.Join(curDir, "./a/b"), Target: "/dir2"},
			},
		}, config, "config")
	})
}

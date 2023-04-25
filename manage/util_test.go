package manage

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/DmitryBogomolov/containerator/core"

	"github.com/stretchr/testify/assert"
)

func TestGetContainerName(t *testing.T) {
	t.Run("Without postfix", func(t *testing.T) {
		cfg := Config{ImageName: "test-name"}
		assert.Equal(t, "test-name", getContainerName(&cfg, ""))
	})

	t.Run("With postfix", func(t *testing.T) {
		cfg := Config{ImageName: "test-name"}
		assert.Equal(t, "test-name-dev", getContainerName(&cfg, "dev"))
	})

	t.Run("Prefer container name", func(t *testing.T) {
		cfg := Config{ImageName: "test-image", ContainerName: "test-container"}
		assert.Equal(t, "test-container-m1", getContainerName(&cfg, "m1"))
	})
}

func TestBuildContainerOptions(t *testing.T) {
	t.Run("With volumes", func(t *testing.T) {
		actual, err := buildContainerOptions(
			&Config{
				Network: "test-net",
				Volumes: []core.Mapping{
					{Source: "/a", Target: "/usr/a"},
					{Source: "/b", Target: "/usr/b"},
				},
			},
			"test-image",
			"test-container",
			&Options{},
		)

		assert.NoError(t, err)
		assert.Equal(
			t,
			&core.RunContainerOptions{
				Image:         "test-image",
				Name:          "test-container",
				RestartPolicy: core.RestartAlways,
				Network:       "test-net",
				Volumes: []core.Mapping{
					{Source: "/a", Target: "/usr/a"},
					{Source: "/b", Target: "/usr/b"},
				},
			},
			actual,
		)
	})

	t.Run("With ports", func(t *testing.T) {
		actual, err := buildContainerOptions(
			&Config{
				Ports: []core.Mapping{
					{Source: "5001", Target: "11"},
					{Source: "5002", Target: "12"},
				},
			},
			"test-image",
			"test-container",
			&Options{},
		)

		assert.NoError(t, err)
		assert.Equal(
			t,
			&core.RunContainerOptions{
				Image:         "test-image",
				Name:          "test-container",
				RestartPolicy: core.RestartAlways,
				Ports: []core.Mapping{
					{Source: "5001", Target: "11"},
					{Source: "5002", Target: "12"},
				},
			},
			actual,
		)
	})

	t.Run("With env file", func(t *testing.T) {
		testContent := strings.Join([]string{
			"B=2x",
		}, "\n")
		testFile := "test.yaml"
		ioutil.WriteFile(testFile, []byte(testContent), os.ModePerm)
		defer os.Remove(testFile)

		actual, err := buildContainerOptions(
			&Config{
				Env: []core.Mapping{
					{Source: "A", Target: "1"},
					{Source: "B", Target: "2"},
				},
			},
			"test-image",
			"test-container",
			&Options{
				EnvFilePath: testFile,
			},
		)

		assert.NoError(t, err)
		assert.Equal(
			t,
			&core.RunContainerOptions{
				Image:         "test-image",
				Name:          "test-container",
				RestartPolicy: core.RestartAlways,
				Env: []core.Mapping{
					{Source: "A", Target: "1"},
					{Source: "B", Target: "2"},
					{Source: "B", Target: "2x"},
				},
			},
			actual,
		)
	})
}

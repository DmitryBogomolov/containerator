package manage

import (
	"testing"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/docker/docker/api/types/container"

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
	actual, err := buildContainerOptions(
		&Config{
			ImageName:     "base-image",
			ContainerName: "base-container",
			Network:       "test-net",
			Ports: []core.Mapping{
				{Source: "5001", Target: "11"},
				{Source: "5002", Target: "12"},
			},
			Volumes: []core.Mapping{
				{Source: "/a", Target: "/usr/a"},
				{Source: "/b", Target: "/usr/b"},
			},
			Env: []core.Mapping{
				{Source: "A", Target: "1"},
				{Source: "B", Target: "2"},
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
			RestartPolicy: container.RestartPolicyAlways,
			Network:       "test-net",
			Ports: []core.Mapping{
				{Source: "5001", Target: "11"},
				{Source: "5002", Target: "12"},
			},
			Volumes: []core.Mapping{
				{Source: "/a", Target: "/usr/a"},
				{Source: "/b", Target: "/usr/b"},
			},
			Env: []core.Mapping{
				{Source: "A", Target: "1"},
				{Source: "B", Target: "2"},
			},
		},
		actual,
	)
}

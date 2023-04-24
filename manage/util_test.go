package manage

import (
	"testing"

	"github.com/DmitryBogomolov/containerator/core"

	"github.com/stretchr/testify/assert"
)

func TestSelectMode(t *testing.T) {
	conf := &Config{
		Modes: []string{"m1", "m2", "m3"},
	}

	t.Run("Unknown mode", func(t *testing.T) {
		mode, i, err := selectMode("m4", conf)

		assert.Error(t, err, "error")
		assert.Equal(t, (err.(*NotValidModeError)).mode, "m4", "error data")
		assert.Equal(t, "", mode, "mode")
		assert.Equal(t, 0, i, "index")
	})

	t.Run("Known mode", func(t *testing.T) {
		mode, i, err := selectMode("m2", conf)

		assert.NoError(t, nil, err, "error")
		assert.Equal(t, "m2", mode, "mode")
		assert.Equal(t, 1, i, "index")
	})

	t.Run("Empty mode", func(t *testing.T) {
		mode, i, err := selectMode("", &Config{})

		assert.NoError(t, nil, err, "error")
		assert.Equal(t, "", mode, "mode")
		assert.Equal(t, 0, i, "index")
	})
}

func TestGetContainerName(t *testing.T) {
	t.Run("Without mode", func(t *testing.T) {
		name := getContainerName(&Config{ImageName: "test-name"}, "")
		assert.Equal(t, "test-name", name)
	})

	t.Run("With mode", func(t *testing.T) {
		name := getContainerName(&Config{ImageName: "test-name"}, "dev")
		assert.Equal(t, "test-name-dev", name)
	})

	t.Run("Prefer container name", func(t *testing.T) {
		name := getContainerName(&Config{ImageName: "test-image", ContainerName: "test-container"}, "m1")
		assert.Equal(t, "test-container-m1", name)
	})
}

func TestBuildContainerOptions(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		expected := &core.RunContainerOptions{
			Image:         "test-image",
			Name:          "test-container",
			RestartPolicy: core.RestartAlways,
			Network:       "test-net",
		}
		actual := buildContainerOptions(&Config{
			Network: "test-net",
		}, "test-image", "test-container", 1)

		assert.Equal(t, expected, actual)
	})

	t.Run("With ports", func(t *testing.T) {
		expected := &core.RunContainerOptions{
			Image:         "test-image",
			Name:          "test-container",
			RestartPolicy: core.RestartAlways,
			Network:       "test-net",
			Ports: []core.Mapping{
				{Source: "210", Target: "1"},
				{Source: "211", Target: "2"},
			},
		}
		actual := buildContainerOptions(&Config{
			Network:    "test-net",
			BasePort:   200,
			PortOffset: 10,
			Ports:      []int{1, 2},
		}, "test-image", "test-container", 1)

		assert.Equal(t, expected, actual)
	})
}

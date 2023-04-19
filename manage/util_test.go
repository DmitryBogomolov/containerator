package manage

import (
	"testing"

	"github.com/DmitryBogomolov/containerator"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSelectMode(t *testing.T) {
	conf := &Config{
		Modes: []string{"m1", "m2", "m3"},
	}

	t.Run("Unknonwn mode", func(t *testing.T) {
		mode, i, err := selectMode("m4", conf)

		assert.Error(t, err, "error")
		assert.Equal(t, (err.(*NotValidModeError)).Mode, "m4", "error data")
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
		name := getContainerName(&Config{ImageRepo: "test-name"}, "")
		assert.Equal(t, "test-name", name)
	})

	t.Run("With mode", func(t *testing.T) {
		name := getContainerName(&Config{ImageRepo: "test-name"}, "dev")
		assert.Equal(t, "test-name-dev", name)
	})

	t.Run("Prefer container name", func(t *testing.T) {
		name := getContainerName(&Config{ImageRepo: "test-image", ContainerName: "test-container"}, "m1")
		assert.Equal(t, "test-container-m1", name)
	})
}

func TestFindImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testImages := []types.ImageSummary{
		{
			RepoTags: []string{"test-image:1"},
		},
		{
			RepoTags: []string{"test-image:2"},
		},
		{
			RepoTags: []string{"test-image:3"},
		},
	}
	cli := test_mocks.NewMockImageAPIClient(ctrl)
	cli.EXPECT().ImageList(gomock.Any(), gomock.Any()).Return(testImages, nil).AnyTimes()

	t.Run("With tag", func(t *testing.T) {
		image, err := findImage(cli, "test-image", "2")

		assert.NoError(t, err, "error")
		assert.Equal(t, &testImages[1], image, "image")
	})

	t.Run("With tag - not found", func(t *testing.T) {
		image, err := findImage(cli, "test-image", "4")

		assert.Error(t, err, "error")
		assert.Equal(t, (err.(*containerator.ImageNotFoundError)).Image, "test-image:4", "error data")
		assert.Equal(t, (*types.ImageSummary)(nil), image, "image")
	})

	t.Run("Without tag", func(t *testing.T) {
		image, err := findImage(cli, "test-image", "")

		assert.NoError(t, err, "error")
		assert.Equal(t, &testImages[0], image, "image")
	})

	t.Run("Without tag - not found", func(t *testing.T) {
		image, err := findImage(cli, "test-image-other", "")

		assert.Error(t, err, "error")
		assert.Equal(t, (err.(*containerator.ImageNotFoundError)).Image, "test-image-other", "error data")
		assert.Equal(t, (*types.ImageSummary)(nil), image, "image")
	})
}

func TestBuildContainerOptions(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		expected := &containerator.RunContainerOptions{
			Image:         "test-image",
			Name:          "test-container",
			RestartPolicy: containerator.RestartAlways,
			Network:       "test-net",
		}
		actual := buildContainerOptions(&Config{
			Network: "test-net",
		}, "test-image", "test-container", 1)

		assert.Equal(t, expected, actual)
	})

	t.Run("With ports", func(t *testing.T) {
		expected := &containerator.RunContainerOptions{
			Image:         "test-image",
			Name:          "test-container",
			RestartPolicy: containerator.RestartAlways,
			Network:       "test-net",
			Ports: []containerator.Mapping{
				{Source: "210", Target: "1"},
				{Source: "211", Target: "2"},
			},
		}
		actual := buildContainerOptions(&Config{
			Network:    "test-net",
			BasePort:   200,
			PortOffset: 10,
			Ports:      []float64{1, 2},
		}, "test-image", "test-container", 1)

		assert.Equal(t, expected, actual)
	})
}

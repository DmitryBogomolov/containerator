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

	mode, i, err := selectMode("", conf)
	assert.Equal(t, "", mode, "mode")
	assert.Equal(t, 0, i, "index")
	assert.EqualErrorf(t, err, "mode '' is not valid", "error")

	mode, i, err = selectMode("m2", conf)
	assert.Equal(t, "m2", mode, "mode")
	assert.Equal(t, 1, i, "index")
	assert.Equal(t, nil, err, "error")
}

func TestGetContainerName(t *testing.T) {
	var name string

	name = getContainerName(&Config{ImageRepo: "cont"}, "")
	assert.Equal(t, "cont", name, "without mode")

	name = getContainerName(&Config{ImageRepo: "cont"}, "m1")
	assert.Equal(t, "cont-m1", name, "with mode")
}

func TestFindImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testImages := []types.ImageSummary{
		types.ImageSummary{
			RepoTags: []string{"test-image:1"},
		},
		types.ImageSummary{
			RepoTags: []string{"test-image:2"},
		},
		types.ImageSummary{
			RepoTags: []string{"test-image:3"},
		},
	}
	cli := test_mocks.NewMockImageAPIClient(ctrl)
	cli.EXPECT().ImageList(gomock.Any(), gomock.Any()).Return(testImages, nil).AnyTimes()

	t.Run("With tag", func(t *testing.T) {
		image, err := findImage(cli, "test-image", "2")

		assert.Equal(t, &testImages[1], image, "image")
		assert.Equal(t, nil, err, "error")

		image, err = findImage(cli, "test-image", "4")

		assert.Equal(t, (*types.ImageSummary)(nil), image, "image")
		assert.EqualError(t, err, "no 'test-image:4' image (image 'test-image:4' is not found)", "error")
	})

	t.Run("Without tag", func(t *testing.T) {
		image, err := findImage(cli, "test-image", "")

		assert.Equal(t, &testImages[0], image, "image")
		assert.Equal(t, nil, err, "error")

		image, err = findImage(cli, "test-image-other", "")

		assert.Equal(t, (*types.ImageSummary)(nil), image, "image")
		assert.EqualError(t, err, "no 'test-image-other' images", "error")
	})
}

func TestBuildContainerOptions(t *testing.T) {
	assert.Equal(t, &containerator.RunContainerOptions{
		Image:         "test-image",
		Name:          "test-container",
		RestartPolicy: containerator.RestartAlways,
		Network:       "test-net",
	}, buildContainerOptions(&Config{
		Network: "test-net",
	}, "test-image", "test-container", 1), "without ports")

	assert.Equal(t, &containerator.RunContainerOptions{
		Image:         "test-image",
		Name:          "test-container",
		RestartPolicy: containerator.RestartAlways,
		Network:       "test-net",
		Ports: []containerator.Mapping{
			containerator.Mapping{Source: "210", Target: "1"},
			containerator.Mapping{Source: "211", Target: "2"},
		},
	}, buildContainerOptions(&Config{
		Network:    "test-net",
		BasePort:   200,
		PortOffset: 10,
		Ports:      []float64{1, 2},
	}, "test-image", "test-container", 1), "with ports")
}

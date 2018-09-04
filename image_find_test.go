package containerator

import (
	"testing"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestFindImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testImages := []types.ImageSummary{
		types.ImageSummary{
			ID:       "i1",
			RepoTags: []string{"test", "test:1"},
			Created:  10,
		},
		types.ImageSummary{
			ID:       "i2",
			RepoTags: []string{},
			Created:  12,
		},
		types.ImageSummary{
			ID:       "i3",
			RepoTags: []string{"test:2"},
			Created:  14,
		},
		types.ImageSummary{
			ID:       "i4",
			RepoTags: []string{"test:3", "test:4"},
			Created:  12,
		},
	}

	cli := test_mocks.NewMockImageAPIClient(ctrl)
	cli.EXPECT().ImageList(gomock.Any(), gomock.Any()).Return(testImages, nil).AnyTimes()

	t.Run("ByID", func(t *testing.T) {
		var image *ImageInfo
		var err error

		image, err = FindImage(cli, FindImageOptions{ID: "i1"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, *image, ImageInfo{ID: "i1", Name: "test", Created: 10}, "image")

		image, err = FindImage(cli, FindImageOptions{ID: "i2"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, *image, ImageInfo{ID: "i2", Created: 12}, "image")

		image, err = FindImage(cli, FindImageOptions{ID: "unknown"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, image, nil, "image")
	})

	t.Run("ByRepo", func(t *testing.T) {
		var image *ImageInfo
		var err error

		image, err = FindImage(cli, FindImageOptions{Repo: "test"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, *image, ImageInfo{ID: "i3", Name: "test:2", Created: 14}, "image")

		image, err = FindImage(cli, FindImageOptions{Repo: "unknown"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, image, nil, "image")
	})

	t.Run("ByRepoAndTag", func(t *testing.T) {
		var image *ImageInfo
		var err error

		image, err = FindImage(cli, FindImageOptions{Repo: "test", Tag: "3"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, *image, ImageInfo{ID: "i4", Name: "test:3", Created: 12}, "image")

		image, err = FindImage(cli, FindImageOptions{Repo: "test", Tag: "4"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, *image, ImageInfo{ID: "i4", Name: "test:4", Created: 12}, "image")

		image, err = FindImage(cli, FindImageOptions{Repo: "test", Tag: "5"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, image, nil, "image")
	})

	t.Run("RaiseError", func(t *testing.T) {
		_, err := FindImage(cli, FindImageOptions{})
		assertEqual(t, err, ErrFindImage, "error")
	})
}

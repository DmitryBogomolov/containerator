package containerator

import (
	"testing"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestFindImageRepoTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	getCli := func() *test_mocks.MockImageAPIClient {
		cli := test_mocks.NewMockImageAPIClient(ctrl)
		list := []types.ImageSummary{
			types.ImageSummary{
				ID:       "i1",
				RepoTags: []string{"test:1", "test:2"},
				Created:  2,
			},
			types.ImageSummary{},
			types.ImageSummary{
				ID:       "i2",
				RepoTags: []string{"test:3", "test:4"},
				Created:  4,
			},
		}
		cli.EXPECT().ImageList(gomock.Any(), gomock.Any()).Return(list, nil)
		return cli
	}

	t.Run("FindByTag", func(t *testing.T) {
		cli := getCli()
		image, err := FindImageByTag(cli, "test:2")

		assertEqual(t, err, nil, "error")
		assertEqual(t, *image, ImageInfo{ID: "i1", Tag: "test:2", Created: 2}, "image")
	})

	t.Run("SortByCreationTime", func(t *testing.T) {
		cli := getCli()
		image, err := FindImageByTag(cli, "test")

		assertEqual(t, err, nil, "error")
		assertEqual(t, *image, ImageInfo{ID: "i2", Tag: "test:3", Created: 4}, "image")
	})

	t.Run("ErrorIfNotFound", func(t *testing.T) {
		cli := getCli()
		_, err := FindImageByTag(cli, "test:5")

		assertEqual(t, err.Error(), "image test:5 is not found", "error")
	})
}

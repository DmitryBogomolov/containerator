package core

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestFindImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testImages := []types.ImageSummary{
		{
			ID:       "sha256:00112233445566778899",
			RepoTags: []string{"test:latest", "test:1"},
		},
		{
			ID:       "sha256:11223344556677889900",
			RepoTags: []string{},
		},
		{
			ID:       "sha256:22334455667788990011",
			RepoTags: []string{"test:2"},
		},
		{
			ID:       "sha256:33445566778899001122",
			RepoTags: []string{"test:3", "test:4"},
		},
	}

	cli := test_mocks.NewMockImageAPIClient(ctrl)
	cli.EXPECT().ImageList(gomock.Any(), gomock.Any()).Return(testImages, nil).AnyTimes()

	t.Run("ByID", func(t *testing.T) {
		image, err := FindImageByID(cli, "sha256:00112233445566778899")
		assert.NoError(t, err)
		assert.Equal(t, makeImage(&testImages[0]), image)
	})

	t.Run("ByID / not found", func(t *testing.T) {
		image, err := FindImageByID(cli, "unknown")
		assert.NoError(t, err)
		assert.Nil(t, image)
	})

	t.Run("ByShortID", func(t *testing.T) {
		image, err := FindImageByShortID(cli, "0011")
		assert.NoError(t, err)
		assert.Equal(t, makeImage(&testImages[0]), image)
	})

	t.Run("ByShortID / not found", func(t *testing.T) {
		image, err := FindImageByShortID(cli, "unknown")
		assert.NoError(t, err)
		assert.Nil(t, image)
	})

	t.Run("ByRepoTag / no tag", func(t *testing.T) {
		image, err := FindImageByName(cli, "test")
		assert.NoError(t, err)
		assert.Equal(t, makeImage(&testImages[0]), image)
	})

	t.Run("ByRepoTag / tag", func(t *testing.T) {
		image, err := FindImageByName(cli, "test:4")
		assert.NoError(t, err)
		assert.Equal(t, makeImage(&testImages[3]), image)
	})

	t.Run("ByRepoTag / not found", func(t *testing.T) {
		image, err := FindImageByName(cli, "unknown")
		assert.NoError(t, err)
		assert.Nil(t, image)
	})

	t.Run("ByRepo", func(t *testing.T) {
		images, err := FindAllImagesByName(cli, "test")
		assert.NoError(t, err)
		expected := []Image{
			makeImage(&testImages[0]),
			makeImage(&testImages[2]),
			makeImage(&testImages[3]),
		}
		assert.Equal(t, expected, images)
	})

	t.Run("ByRepo / not found", func(t *testing.T) {
		images, err := FindAllImagesByName(cli, "unknown")
		assert.NoError(t, err)
		assert.Equal(t, []Image(nil), images)
	})
}

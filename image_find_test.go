package containerator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestGetImageFullName(t *testing.T) {
	t.Run("No RepoTags", func(t *testing.T) {
		name := GetImageFullName(&types.ImageSummary{RepoTags: []string{}})
		assert.Equal(t, "", name)
	})

	t.Run("First RepoTag", func(t *testing.T) {
		name := GetImageFullName(&types.ImageSummary{RepoTags: []string{"a:1", "b:2"}})
		assert.Equal(t, "a:1", name)
	})
}

func TestSplitImageNameTag(t *testing.T) {
	t.Run("With tag", func(t *testing.T) {
		name, tag := SplitImageNameTag("a:1")
		assert.Equal(t, "a", name)
		assert.Equal(t, "1", tag)
	})

	t.Run("Without tag", func(t *testing.T) {
		name, tag := SplitImageNameTag("b")
		assert.Equal(t, "b", name)
		assert.Equal(t, "latest", tag)
	})

	t.Run("Panic on empty string", func(t *testing.T) {
		assert.Panics(t, func() {
			SplitImageNameTag("")
		})
	})
}

func TestJoinImageNameTag(t *testing.T) {
	t.Run("With tag", func(t *testing.T) {
		name := JoinImageNameTag("a", "1")
		assert.Equal(t, "a:1", name)
	})

	t.Run("Without tag", func(t *testing.T) {
		name := JoinImageNameTag("b", "")
		assert.Equal(t, "b:latest", name)
	})

	t.Run("Panic on empty string", func(t *testing.T) {
		assert.Panics(t, func() {
			JoinImageNameTag("", "")
		})
	})
}

func TestGetImageShortID(t *testing.T) {
	id := GetImageShortID(&types.ImageSummary{ID: "sha256:01234567890123456789"})
	assert.Equal(t, "012345678901", id)
}

func TestFindImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testImages := []types.ImageSummary{
		types.ImageSummary{
			ID:       "sha256:00112233445566778899",
			RepoTags: []string{"test:latest", "test:1"},
		},
		types.ImageSummary{
			ID:       "sha256:11223344556677889900",
			RepoTags: []string{},
		},
		types.ImageSummary{
			ID:       "sha256:22334455667788990011",
			RepoTags: []string{"test:2"},
		},
		types.ImageSummary{
			ID:       "sha256:33445566778899001122",
			RepoTags: []string{"test:3", "test:4"},
		},
	}

	cli := test_mocks.NewMockImageAPIClient(ctrl)
	cli.EXPECT().ImageList(gomock.Any(), gomock.Any()).Return(testImages, nil).AnyTimes()

	t.Run("ByID", func(t *testing.T) {
		image, err := FindImageByID(cli, "sha256:00112233445566778899")
		assert.NoError(t, err)
		assert.Equal(t, &testImages[0], image)
	})

	t.Run("ByID / not found", func(t *testing.T) {
		image, err := FindImageByID(cli, "unknown")
		assert.Error(t, err)
		imageErr, ok := err.(*ImageNotFoundError)
		assert.Equal(t, true, ok && imageErr.Image == "unknown")
		assert.Nil(t, image)
	})

	t.Run("ByShortID", func(t *testing.T) {
		image, err := FindImageByShortID(cli, "0011")
		assert.NoError(t, err)
		assert.Equal(t, &testImages[0], image)
	})

	t.Run("ByShortID / not found", func(t *testing.T) {
		image, err := FindImageByShortID(cli, "unknown")
		assert.Error(t, err)
		imageErr, ok := err.(*ImageNotFoundError)
		assert.Equal(t, true, ok && imageErr.Image == "unknown")
		assert.Nil(t, image)
	})

	t.Run("ByRepoTag / no tag", func(t *testing.T) {
		image, err := FindImageByRepoTag(cli, "test")
		assert.NoError(t, err)
		assert.Equal(t, &testImages[0], image)
	})

	t.Run("ByRepoTag / tag", func(t *testing.T) {
		image, err := FindImageByRepoTag(cli, "test:4")
		assert.NoError(t, err)
		assert.Equal(t, &testImages[3], image)
	})

	t.Run("ByRepoTag / not found", func(t *testing.T) {
		image, err := FindImageByRepoTag(cli, "unknown")
		assert.Error(t, err)
		imageErr, ok := err.(*ImageNotFoundError)
		assert.Equal(t, true, ok && imageErr.Image == "unknown")
		assert.Nil(t, image)
	})

	t.Run("ByRepo", func(t *testing.T) {
		images, err := FindImagesByRepo(cli, "test")
		assert.NoError(t, err)
		expected := []*types.ImageSummary{&testImages[0], &testImages[2], &testImages[3]}
		assert.Equal(t, expected, images)
	})

	t.Run("ByRepo / not found", func(t *testing.T) {
		images, err := FindImagesByRepo(cli, "unknown")
		assert.NoError(t, err)
		var expected []*types.ImageSummary
		assert.Equal(t, expected, images)
	})
}

func TestGetImagesTags(t *testing.T) {
	tags := GetImagesTags([]*types.ImageSummary{
		&types.ImageSummary{RepoTags: []string{"a:1"}},
		&types.ImageSummary{RepoTags: []string{"a:2"}},
		&types.ImageSummary{RepoTags: []string{"a"}},
	})

	assert.Equal(t, []string{"1", "2", "latest"}, tags)
}

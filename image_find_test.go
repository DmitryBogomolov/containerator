package containerator

import (
	"testing"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestGetImageFullName(t *testing.T) {
	var name string

	name = GetImageFullName(&types.ImageSummary{RepoTags: []string{}})
	assertEqual(t, name, "", "name")

	name = GetImageFullName(&types.ImageSummary{RepoTags: []string{"a:1", "b:2"}})
	assertEqual(t, name, "a:1", "name")
}

func TestSplitImageNameTag(t *testing.T) {
	var name, tag string

	name, tag = SplitImageNameTag("a:1")
	assertEqual(t, name, "a", "name")
	assertEqual(t, tag, "1", "tag")

	name, tag = SplitImageNameTag("b")
	assertEqual(t, name, "b", "name")
	assertEqual(t, tag, "latest", "tag")
}

func TestJoinImageNameTag(t *testing.T) {
	var name string

	name = JoinImageNameTag("a", "1")
	assertEqual(t, name, "a:1", "name")

	name = JoinImageNameTag("b", "")
	assertEqual(t, name, "b:latest", "name")
}

func TestGetImageShortID(t *testing.T) {
	var id string

	id = GetImageShortID(&types.ImageSummary{ID: "sha256:01234567890123456789"})
	assertEqual(t, id, "012345678901", "id")
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
		var image *types.ImageSummary
		var err error

		image, err = FindImageByID(cli, "sha256:00112233445566778899")
		assertEqual(t, err, nil, "error")
		assertEqual(t, image, &testImages[0], "image")

		image, err = FindImageByID(cli, "sha256:11223344556677889900")
		assertEqual(t, err, nil, "error")
		assertEqual(t, image, &testImages[1], "image")

		image, err = FindImageByID(cli, "unknown")
		assertEqual(t, err, ErrImageNotFound, "error")
		assertEqual(t, image, nil, "image")
	})

	t.Run("ByShortID", func(t *testing.T) {
		var image *types.ImageSummary
		var err error

		image, err = FindImageByShortID(cli, "0011")
		assertEqual(t, err, nil, "error")
		assertEqual(t, image, &testImages[0], "image")

		image, err = FindImageByShortID(cli, "unknown")
		assertEqual(t, err, ErrImageNotFound, "error")
		assertEqual(t, image, nil, "image")
	})

	t.Run("ByRepoTag", func(t *testing.T) {
		var image *types.ImageSummary
		var err error

		image, err = FindImageByRepoTag(cli, "test")
		assertEqual(t, err, nil, "error")
		assertEqual(t, image, &testImages[0], "image")

		image, err = FindImageByRepoTag(cli, "test:4")
		assertEqual(t, err, nil, "error")
		assertEqual(t, image, &testImages[3], "image")

		image, err = FindImageByRepoTag(cli, "unknown")
		assertEqual(t, err, ErrImageNotFound, "error")
		assertEqual(t, image, nil, "image")
	})

	t.Run("ByRepo", func(t *testing.T) {
		var images []*types.ImageSummary
		var err error

		images, err = FindImagesByRepo(cli, "test")
		assertEqual(t, err, nil, "error")
		assertEqual(t, len(images), 3, "images count")
		assertEqual(t, images[0], &testImages[0], "image 1")
		assertEqual(t, images[1], &testImages[2], "image 2")
		assertEqual(t, images[2], &testImages[3], "image 3")

		images, err = FindImagesByRepo(cli, "unknown")
		assertEqual(t, err, nil, "error")
		assertEqual(t, len(images), 0, "images count")
	})
}

func TestGetImagesTags(t *testing.T) {
	tags := GetImagesTags([]*types.ImageSummary{
		&types.ImageSummary{RepoTags: []string{"a:1"}},
		&types.ImageSummary{RepoTags: []string{"a:2"}},
		&types.ImageSummary{RepoTags: []string{"a"}},
	})

	assertEqual(t, len(tags), 3, "count")
	assertEqual(t, tags[0], "1", "item 1")
	assertEqual(t, tags[1], "2", "item 1")
	assertEqual(t, tags[2], "latest", "item 1")
}

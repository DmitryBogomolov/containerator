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

	t.Run("Searches for tag", func(t *testing.T) {
		cli := getCli()
		image, err := FindImageByTag(cli, "test:2")
		if err != nil {
			t.Fatal(err)
		}
		expected := ImageInfo{ID: "i1", Tag: "test:2", Created: 2}
		if *image != expected {
			t.Fatalf("tag: %v / expected: %v", image, expected)
		}
	})

	t.Run("Sorts by creation time", func(t *testing.T) {
		cli := getCli()
		image, err := FindImageByTag(cli, "test")
		if err != nil {
			t.Fatal(err)
		}
		expected := ImageInfo{ID: "i2", Tag: "test:3", Created: 4}
		if *image != expected {
			t.Fatalf("tag: %v / expected: %v", image, expected)
		}
	})

	t.Run("Returns error if nothing is found", func(t *testing.T) {
		cli := getCli()
		_, err := FindImageByTag(cli, "test:5")
		if err == nil {
			t.Fatal("Error is expected")
		}
		const expectedErr = "image test:5 is not found"
		if err.Error() != expectedErr {
			t.Fatalf("error: %s / expected: %s", err, expectedErr)
		}
	})
}

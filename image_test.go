package containerator

import (
	"context"
	"io"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
)

type stubFunc func(args ...interface{}) (interface{}, error)

type testImageAPIClient struct {
	stub stubFunc
}

func (cli *testImageAPIClient) ImageBuild(ctx context.Context, context io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return types.ImageBuildResponse{}, nil
}

func (cli *testImageAPIClient) BuildCachePrune(ctx context.Context) (*types.BuildCachePruneReport, error) {
	return nil, nil
}

func (cli *testImageAPIClient) BuildCancel(ctx context.Context, id string) error {
	return nil
}

func (cli *testImageAPIClient) ImageCreate(ctx context.Context, parentReference string, options types.ImageCreateOptions) (io.ReadCloser, error) {
	return nil, nil
}

func (cli *testImageAPIClient) ImageHistory(ctx context.Context, image string) ([]image.HistoryResponseItem, error) {
	return nil, nil
}

func (cli *testImageAPIClient) ImageImport(ctx context.Context, source types.ImageImportSource, ref string, options types.ImageImportOptions) (io.ReadCloser, error) {
	return nil, nil
}
func (cli *testImageAPIClient) ImageInspectWithRaw(ctx context.Context, image string) (types.ImageInspect, []byte, error) {
	return types.ImageInspect{}, nil, nil
}

func (cli *testImageAPIClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	ret, err := cli.stub(ctx, options)
	return ret.([]types.ImageSummary), err
}

func (cli *testImageAPIClient) ImageLoad(ctx context.Context, input io.Reader, quiet bool) (types.ImageLoadResponse, error) {
	return types.ImageLoadResponse{}, nil
}

func (cli *testImageAPIClient) ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return nil, nil
}

func (cli *testImageAPIClient) ImagePush(ctx context.Context, ref string, options types.ImagePushOptions) (io.ReadCloser, error) {
	return nil, nil
}

func (cli *testImageAPIClient) ImageRemove(ctx context.Context, image string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	return nil, nil
}

func (cli *testImageAPIClient) ImageSearch(ctx context.Context, term string, options types.ImageSearchOptions) ([]registry.SearchResult, error) {
	return nil, nil
}

func (cli *testImageAPIClient) ImageSave(ctx context.Context, images []string) (io.ReadCloser, error) {
	return nil, nil
}

func (cli *testImageAPIClient) ImageTag(ctx context.Context, image, ref string) error {
	return nil
}

func (cli *testImageAPIClient) ImagesPrune(ctx context.Context, pruneFilter filters.Args) (types.ImagesPruneReport, error) {
	return types.ImagesPruneReport{}, nil
}

func TestFindImageRepoTag(t *testing.T) {

	getImages := func(args ...interface{}) (interface{}, error) {
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
		return list, nil
	}

	cli := &testImageAPIClient{stub: getImages}

	t.Run("Searches for tag", func(t *testing.T) {
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

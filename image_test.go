package containerator

import (
	"context"
	"io"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

type testImageAPIClient struct{}

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
	return nil, nil
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
	var cli client.ImageAPIClient = &testImageAPIClient{}

	tag, err := findImageRepoTag(cli, "test")
	if err != nil {
		t.Fatal(err)
	}
	if tag != "test" {
		t.Fatalf("Tag: %s", tag)
	}
}

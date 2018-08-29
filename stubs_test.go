package containerator

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
)

type stubFunc func(args ...interface{}) (interface{}, error)

type testImageAPIClient struct {
	stub stubFunc
}

type testContainerAPIClient struct {
	stub stubFunc
}

func (cli *testImageAPIClient) ImageBuild(ctx context.Context, context io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	ret, err := cli.stub(ctx, context, options)
	return ret.(types.ImageBuildResponse), err
}

func (cli *testImageAPIClient) BuildCachePrune(ctx context.Context) (*types.BuildCachePruneReport, error) {
	ret, err := cli.stub(ctx)
	return ret.(*types.BuildCachePruneReport), err
}

func (cli *testImageAPIClient) BuildCancel(ctx context.Context, id string) error {
	_, err := cli.stub(ctx, id)
	return err
}

func (cli *testImageAPIClient) ImageCreate(ctx context.Context, parentReference string, options types.ImageCreateOptions) (io.ReadCloser, error) {
	ret, err := cli.stub(ctx, parentReference, options)
	return ret.(io.ReadCloser), err
}

func (cli *testImageAPIClient) ImageHistory(ctx context.Context, img string) ([]image.HistoryResponseItem, error) {
	ret, err := cli.stub(ctx, img)
	return ret.([]image.HistoryResponseItem), err
}

func (cli *testImageAPIClient) ImageImport(ctx context.Context, source types.ImageImportSource, ref string, options types.ImageImportOptions) (io.ReadCloser, error) {
	ret, err := cli.stub(ctx, source, ref, options)
	return ret.(io.ReadCloser), err
}
func (cli *testImageAPIClient) ImageInspectWithRaw(ctx context.Context, image string) (types.ImageInspect, []byte, error) {
	ret, err := cli.stub(ctx, image)
	return ret.(types.ImageInspect), nil, err
}

func (cli *testImageAPIClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	ret, err := cli.stub(ctx, options)
	return ret.([]types.ImageSummary), err
}

func (cli *testImageAPIClient) ImageLoad(ctx context.Context, input io.Reader, quiet bool) (types.ImageLoadResponse, error) {
	ret, err := cli.stub(ctx, input, quiet)
	return ret.(types.ImageLoadResponse), err
}

func (cli *testImageAPIClient) ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error) {
	ret, err := cli.stub(ctx, ref, options)
	return ret.(io.ReadCloser), err
}

func (cli *testImageAPIClient) ImagePush(ctx context.Context, ref string, options types.ImagePushOptions) (io.ReadCloser, error) {
	ret, err := cli.stub(ctx, ref, options)
	return ret.(io.ReadCloser), err
}

func (cli *testImageAPIClient) ImageRemove(ctx context.Context, image string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	ret, err := cli.stub(ctx, image, options)
	return ret.([]types.ImageDeleteResponseItem), err
}

func (cli *testImageAPIClient) ImageSearch(ctx context.Context, term string, options types.ImageSearchOptions) ([]registry.SearchResult, error) {
	ret, err := cli.stub(ctx, term, options)
	return ret.([]registry.SearchResult), err
}

func (cli *testImageAPIClient) ImageSave(ctx context.Context, images []string) (io.ReadCloser, error) {
	ret, err := cli.stub(ctx, images)
	return ret.(io.ReadCloser), err
}

func (cli *testImageAPIClient) ImageTag(ctx context.Context, image, ref string) error {
	_, err := cli.stub(ctx, image, ref)
	return err
}

func (cli *testImageAPIClient) ImagesPrune(ctx context.Context, pruneFilter filters.Args) (types.ImagesPruneReport, error) {
	ret, err := cli.stub(ctx, pruneFilter)
	return ret.(types.ImagesPruneReport), err
}

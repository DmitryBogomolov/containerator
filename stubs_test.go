package containerator

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
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

func (cli *testContainerAPIClient) ContainerAttach(ctx context.Context, container string, options types.ContainerAttachOptions) (types.HijackedResponse, error) {
	ret, err := cli.stub(ctx, container, options)
	return ret.(types.HijackedResponse), err
}

func (cli *testContainerAPIClient) ContainerCommit(ctx context.Context, container string, options types.ContainerCommitOptions) (types.IDResponse, error) {
	ret, err := cli.stub(ctx, container, options)
	return ret.(types.IDResponse), err
}

func (cli *testContainerAPIClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error) {
	ret, err := cli.stub(ctx, config, hostConfig, networkingConfig, containerName)
	return ret.(container.ContainerCreateCreatedBody), err
}

func (cli *testContainerAPIClient) ContainerDiff(ctx context.Context, cont string) ([]container.ContainerChangeResponseItem, error) {
	ret, err := cli.stub(ctx, cont)
	return ret.([]container.ContainerChangeResponseItem), err
}

func (cli *testContainerAPIClient) ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error) {
	ret, err := cli.stub(ctx, execID, config)
	return ret.(types.HijackedResponse), err
}

func (cli *testContainerAPIClient) ContainerExecCreate(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error) {
	ret, err := cli.stub(ctx, container, config)
	return ret.(types.IDResponse), err
}

func (cli *testContainerAPIClient) ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error) {
	ret, err := cli.stub(ctx, execID)
	return ret.(types.ContainerExecInspect), err
}

func (cli *testContainerAPIClient) ContainerExecResize(ctx context.Context, execID string, options types.ResizeOptions) error {
	_, err := cli.stub(ctx, execID, options)
	return err
}

func (cli *testContainerAPIClient) ContainerExecStart(ctx context.Context, execID string, config types.ExecStartCheck) error {
	_, err := cli.stub(ctx, execID, config)
	return err
}

func (cli *testContainerAPIClient) ContainerExport(ctx context.Context, container string) (io.ReadCloser, error) {
	ret, err := cli.stub(ctx, container)
	return ret.(io.ReadCloser), err
}

func (cli *testContainerAPIClient) ContainerInspect(ctx context.Context, container string) (types.ContainerJSON, error) {
	ret, err := cli.stub(ctx, container)
	return ret.(types.ContainerJSON), err
}

func (cli *testContainerAPIClient) ContainerInspectWithRaw(ctx context.Context, container string, getSize bool) (types.ContainerJSON, []byte, error) {
	ret, err := cli.stub(ctx, container, getSize)
	return ret.(types.ContainerJSON), nil, err
}

func (cli *testContainerAPIClient) ContainerKill(ctx context.Context, container, signal string) error {
	_, err := cli.stub(ctx, container, signal)
	return err
}

func (cli *testContainerAPIClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	ret, err := cli.stub(ctx, options)
	return ret.([]types.Container), err
}

func (cli *testContainerAPIClient) ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	ret, err := cli.stub(ctx, container, options)
	return ret.(io.ReadCloser), err
}

func (cli *testContainerAPIClient) ContainerPause(ctx context.Context, container string) error {
	_, err := cli.stub(ctx, container)
	return err
}

func (cli *testContainerAPIClient) ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error {
	_, err := cli.stub(ctx, container, options)
	return err
}

func (cli *testContainerAPIClient) ContainerRename(ctx context.Context, container, newContainerName string) error {
	_, err := cli.stub(ctx, container, newContainerName)
	return err
}

func (cli *testContainerAPIClient) ContainerResize(ctx context.Context, container string, options types.ResizeOptions) error {
	_, err := cli.stub(ctx, container, options)
	return err
}

func (cli *testContainerAPIClient) ContainerRestart(ctx context.Context, container string, timeout *time.Duration) error {
	_, err := cli.stub(ctx, container, timeout)
	return err
}

func (cli *testContainerAPIClient) ContainerStatPath(ctx context.Context, container, path string) (types.ContainerPathStat, error) {
	ret, err := cli.stub(ctx, container, path)
	return ret.(types.ContainerPathStat), err
}

func (cli *testContainerAPIClient) ContainerStats(ctx context.Context, container string, stream bool) (types.ContainerStats, error) {
	ret, err := cli.stub(ctx, container, stream)
	return ret.(types.ContainerStats), err
}

func (cli *testContainerAPIClient) ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error {
	_, err := cli.stub(ctx, container, options)
	return err
}

func (cli *testContainerAPIClient) ContainerStop(ctx context.Context, container string, timeout *time.Duration) error {
	_, err := cli.stub(ctx, container, timeout)
	return err
}

func (cli *testContainerAPIClient) ContainerTop(ctx context.Context, cont string, arguments []string) (container.ContainerTopOKBody, error) {
	ret, err := cli.stub(ctx, cont, arguments)
	return ret.(container.ContainerTopOKBody), err
}

func (cli *testContainerAPIClient) ContainerUnpause(ctx context.Context, container string) error {
	_, err := cli.stub(ctx, container)
	return err
}

func (cli *testContainerAPIClient) ContainerUpdate(ctx context.Context, cont string, updateConfig container.UpdateConfig) (container.ContainerUpdateOKBody, error) {
	ret, err := cli.stub(ctx, cont, updateConfig)
	return ret.(container.ContainerUpdateOKBody), err
}

func (cli *testContainerAPIClient) ContainerWait(ctx context.Context, cont string, condition container.WaitCondition) (<-chan container.ContainerWaitOKBody, <-chan error) {
	ret, err := cli.stub(ctx, cont, condition)
	return ret.(<-chan container.ContainerWaitOKBody), err.(interface{}).(<-chan error)
}

func (cli *testContainerAPIClient) CopyFromContainer(ctx context.Context, container, srcPath string) (io.ReadCloser, types.ContainerPathStat, error) {
	ret, err := cli.stub(ctx, container, srcPath)
	return ret.(io.ReadCloser), types.ContainerPathStat{}, err
}

func (cli *testContainerAPIClient) CopyToContainer(ctx context.Context, container, path string, content io.Reader, options types.CopyToContainerOptions) error {
	_, err := cli.stub(ctx, container, path, content, options)
	return err
}

func (cli *testContainerAPIClient) ContainersPrune(ctx context.Context, pruneFilters filters.Args) (types.ContainersPruneReport, error) {
	ret, err := cli.stub(ctx, pruneFilters)
	return ret.(types.ContainersPruneReport), err
}

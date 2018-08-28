package containerator

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	contextTimeout = 3 * time.Second
)

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), contextTimeout)
}

func cliImageList(cli *client.Client) ([]types.ImageSummary, error) {
	ctx, cancel := getContext()
	defer cancel()
	return cli.ImageList(ctx, types.ImageListOptions{})
}

func cliContainerCreate(cli *client.Client, config *container.Config, hostConfig *container.HostConfig, containerName string) (container.ContainerCreateCreatedBody, error) {
	ctx, cancel := getContext()
	defer cancel()
	return cli.ContainerCreate(ctx, config, hostConfig, nil, containerName)
}

func cliContainerStart(cli *client.Client, containerID string) error {
	ctx, cancel := getContext()
	defer cancel()
	return cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

func cliContainerRemove(cli *client.Client, containerID string) error {
	ctx, cancel := getContext()
	defer cancel()
	return cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
}

func cliContainerInspect(cli *client.Client, containerID string) (types.ContainerJSON, error) {
	ctx, cancel := getContext()
	defer cancel()
	return cli.ContainerInspect(ctx, containerID)
}

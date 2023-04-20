// Package core contains functions to work with docker containers.
package core

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

//go:generate mockgen -destination ../test_mocks/mock_imageapiclient.go -package test_mocks github.com/docker/docker/client ImageAPIClient
//go:generate mockgen -destination ../test_mocks/mock_containerapiclient.go -package test_mocks github.com/docker/docker/client ContainerAPIClient

const (
	contextTimeout = 10 * time.Second
)

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), contextTimeout)
}

func cliImageList(cli client.ImageAPIClient) ([]types.ImageSummary, error) {
	ctx, cancel := getContext()
	defer cancel()
	return cli.ImageList(ctx, types.ImageListOptions{})
}

func cliContainerList(cli client.ContainerAPIClient) ([]types.Container, error) {
	ctx, cancel := getContext()
	defer cancel()
	return cli.ContainerList(ctx, types.ContainerListOptions{All: true})
}

func cliContainerCreate(
	cli client.ContainerAPIClient,
	config *container.Config, hostConfig *container.HostConfig, name string,
) (container.ContainerCreateCreatedBody, error) {
	ctx, cancel := getContext()
	defer cancel()
	return cli.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
}

func cliContainerStart(cli client.ContainerAPIClient, name string) error {
	ctx, cancel := getContext()
	defer cancel()
	return cli.ContainerStart(ctx, name, types.ContainerStartOptions{})
}

func cliContainerStop(cli client.ContainerAPIClient, name string) error {
	ctx, cancel := getContext()
	defer cancel()
	return cli.ContainerStop(ctx, name, container.StopOptions{})
}

func cliContainerRename(cli client.ContainerAPIClient, name string, newName string) error {
	ctx, cancel := getContext()
	defer cancel()
	return cli.ContainerRename(ctx, name, newName)
}

func cliContainerRemove(cli client.ContainerAPIClient, name string) error {
	ctx, cancel := getContext()
	defer cancel()
	return cli.ContainerRemove(ctx, name, types.ContainerRemoveOptions{Force: true})
}

func cliContainerInspect(cli client.ContainerAPIClient, name string) (types.ContainerJSON, error) {
	ctx, cancel := getContext()
	defer cancel()
	return cli.ContainerInspect(ctx, name)
}

package containerator

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func cliContainerRemove(cli *client.Client, containerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
}

func removeContainer(cli *client.Client, containerID string) error {
	return cliContainerRemove(cli, containerID)
}

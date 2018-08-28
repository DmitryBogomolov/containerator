package containerator

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func cliContainerInspect(cli *client.Client, containerID string) (types.ContainerJSON, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return cli.ContainerInspect(ctx, containerID)
}

func inspectContainer(cli *client.Client, containerID string) (string, error) {
	data, err := cliContainerInspect(cli, containerID)
	if err != nil {
		return "", err
	}
	return data.State.Status, nil
}

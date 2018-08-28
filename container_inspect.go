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

func inspectContainer(cli *client.Client, containerID string) (containerInfo, error) {
	data, err := cliContainerInspect(cli, containerID)
	if err != nil {
		return containerInfo{}, err
	}
	return containerInfo{
		ID:    data.ID,
		Name:  data.Name,
		Image: data.Image,
		State: data.State.Status,
	}, nil
}

package containerator

import (
	"github.com/docker/docker/client"
)

func inspectContainer(cli client.ContainerAPIClient, containerID string) (*ContainerInfo, error) {
	data, err := cliContainerInspect(cli, containerID)
	if err != nil {
		return nil, err
	}
	return &ContainerInfo{
		ID:    data.ID,
		Name:  data.Name,
		Image: data.Image,
		State: data.State.Status,
	}, nil
}

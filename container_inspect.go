package containerator

import (
	"github.com/docker/docker/client"
)

func inspectContainer(cli client.ContainerAPIClient, containerID string) (containerInfo, error) {
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

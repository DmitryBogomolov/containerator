package containerator

import (
	"github.com/docker/docker/client"
)

func removeContainer(cli client.ContainerAPIClient, containerID string) error {
	return cliContainerRemove(cli, containerID)
}

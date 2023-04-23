package core

import (
	"github.com/docker/docker/client"
)

// RemoveContainer removes container.
func RemoveContainer(cli client.ContainerAPIClient, container Container) error {
	return cliContainerRemove(cli, container.ID())
}

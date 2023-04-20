package core

import (
	"github.com/docker/docker/client"
)

// RemoveContainer removes container.
//
//	RemoveContainer(cli, "my-container")
func RemoveContainer(cli client.ContainerAPIClient, name string) error {
	return cliContainerRemove(cli, name)
}

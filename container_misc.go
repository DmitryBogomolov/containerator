package containerator

import (
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/namesgenerator"
)

/*
RemoveContainer removes container.

	RemoveContainer(cli, "my-container")
*/
func RemoveContainer(cli client.ContainerAPIClient, container string) error {
	return cliContainerRemove(cli, container)
}

/*
SuspendContainer renames and stops container.

Uses docker names generator to acquire temporary container name.
Accepts container ID (not a name) because container ID persists.

	SuspendContainer(cli, "1234")
*/
func SuspendContainer(cli client.ContainerAPIClient, containerID string) error {
	tmpName := namesgenerator.GetRandomName(2)
	if err := cliContainerRename(cli, containerID, tmpName); err != nil {
		return err
	}
	if err := cliContainerStop(cli, containerID); err != nil {
		return err
	}
	return nil
}

/*
ResumeContainer renames and starts container.

	ResumeContainer(cli, "1234", "my-container")
*/
func ResumeContainer(cli client.ContainerAPIClient, containerID string, name string) error {
	if err := cliContainerRename(cli, containerID, name); err != nil {
		return err
	}
	if err := cliContainerStart(cli, containerID); err != nil {
		return err
	}
	return nil
}

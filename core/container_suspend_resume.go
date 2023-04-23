package core

import (
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/namesgenerator"
)

// SuspendContainer renames and stops container.
//
// Uses docker names generator to acquire temporary container name.
func SuspendContainer(cli client.ContainerAPIClient, container Container) error {
	tmpName := namesgenerator.GetRandomName(2)
	if err := cliContainerRename(cli, container.ID(), tmpName); err != nil {
		return err
	}
	if err := cliContainerStop(cli, container.ID()); err != nil {
		return err
	}
	return nil
}

// ResumeContainer renames and starts container.
func ResumeContainer(cli client.ContainerAPIClient, container Container, name string) error {
	if err := cliContainerRename(cli, container.ID(), name); err != nil {
		return err
	}
	if err := cliContainerStart(cli, container.ID()); err != nil {
		return err
	}
	return nil
}

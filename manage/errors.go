package manage

import (
	"fmt"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/docker/docker/api/types"
)

// NoContainerError is returned on attempt to remove container when it is not found.
type NoContainerError struct {
	container string
}

func (err NoContainerError) Error() string {
	return fmt.Sprintf("container '%s' is not found", err.container)
}

// Container returns container name.
func (err NoContainerError) Container() string {
	return err.container
}

// ContainerAlreadyRunningError is returned on attempt to run container when similar container is already running.
type ContainerAlreadyRunningError struct {
	container *types.Container
}

func (err ContainerAlreadyRunningError) Error() string {
	return fmt.Sprintf("container '%s' is already running", core.GetContainerName(err.container))
}

// Container returns running container.
func (err ContainerAlreadyRunningError) Container() *types.Container {
	return err.container
}

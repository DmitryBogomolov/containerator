package manage

import (
	"fmt"
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

// NoImageError is returned when image is not found.
type NoImageError struct {
	image string
}

func (err NoImageError) Error() string {
	return fmt.Sprintf("image '%s' is not found", err.image)
}

// Image returns image name.
func (err NoImageError) Image() string {
	return err.image
}

// ContainerAlreadyRunningError is returned on attempt to run container when similar container is already running.
type ContainerAlreadyRunningError struct {
	container string
}

func (err ContainerAlreadyRunningError) Error() string {
	return fmt.Sprintf("container '%s' is already running", err.container)
}

// Container returns running container.
func (err ContainerAlreadyRunningError) Container() string {
	return err.container
}

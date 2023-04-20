package errors

import "fmt"

// ContainerNotFoundError indicates that container with specified ID or name is not found.
type ContainerNotFoundError struct {
	container string
}

// ContainerNotFound creates ContainerNotFoundError.
func ContainerNotFound(container string) ContainerNotFoundError {
	return ContainerNotFoundError{container}
}

func (err ContainerNotFoundError) Error() string {
	return fmt.Sprintf("container '%s' is not found", err.container)
}

// Container returns container ID or name.
func (err ContainerNotFoundError) Container() string {
	return err.container
}

// ImageNotFoundError indicates that image with specified ID or full name is not found.
type ImageNotFoundError struct {
	image string
}

// ImageNotFound creates ImageNotFoundError.
func ImageNotFound(image string) ImageNotFoundError {
	return ImageNotFoundError{image}
}

func (err ImageNotFoundError) Error() string {
	return fmt.Sprintf("image '%s' is not found", err.image)
}

// Image returns image ID or name.
func (err ImageNotFoundError) Image() string {
	return err.image
}

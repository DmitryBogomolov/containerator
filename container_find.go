package containerator

import (
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func getContainerName(container *types.Container) string {
	if len(container.Names) > 0 {
		return container.Names[0][1:]
	}
	return ""
}

func findContainerByID(id string, containers []types.Container) *ContainerInfo {
	for _, container := range containers {
		if container.ID == id {
			return &ContainerInfo{
				ID:    id,
				Image: container.Image,
				Name:  getContainerName(&container),
				State: container.State,
			}
		}
	}
	return nil
}

func findContainerByName(name string, containers []types.Container) *ContainerInfo {
	for _, container := range containers {
		for _, val := range container.Names {
			if val[1:] == name {
				return &ContainerInfo{
					ID:    container.ID,
					Image: container.Image,
					Name:  name,
					State: container.State,
				}
			}
		}
	}
	return nil
}

func findContainerByImageID(imageID string, containers []types.Container) *ContainerInfo {
	for _, container := range containers {
		if container.ImageID == imageID {
			return &ContainerInfo{
				ID:    container.ID,
				Image: container.Image,
				Name:  getContainerName(&container),
				State: container.State,
			}
		}
	}
	return nil
}

func findContainerByImage(image string, containers []types.Container) *ContainerInfo {
	for _, container := range containers {
		if container.Image == image {
			return &ContainerInfo{
				ID:    container.ID,
				Image: image,
				Name:  getContainerName(&container),
				State: container.State,
			}
		}
	}
	return nil
}

// ErrFindContainer shows that search options are not valid.
var ErrFindContainer = errors.New("search options are not valid")

// FindContainerOptions defines search options.
type FindContainerOptions struct {
	ID      string
	Name    string
	ImageID string
	Image   string
}

// FindContainer searches container.
func FindContainer(cli client.ContainerAPIClient, options FindContainerOptions) (*ContainerInfo, error) {
	containers, err := cliContainerList(cli)
	if err != nil {
		return nil, err
	}
	if options.ID != "" {
		return findContainerByID(options.ID, containers), nil
	}
	if options.Name != "" {
		return findContainerByName(options.Name, containers), nil
	}
	if options.ImageID != "" {
		return findContainerByImageID(options.ImageID, containers), nil
	}
	if options.Image != "" {
		return findContainerByImage(options.Image, containers), nil
	}
	return nil, ErrFindContainer
}

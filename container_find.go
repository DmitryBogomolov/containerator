package containerator

import (
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// GetContainerName returns friendly image name.
func GetContainerName(container *types.Container) string {
	if len(container.Names) > 0 {
		return container.Names[0][1:]
	}
	return ""
}

// GetContainerShortID return short container id.
func GetContainerShortID(container *types.Container) string {
	return container.ID[:shortIDLength]
}

// FindContainerByID searches container by id.
func FindContainerByID(cli client.ContainerAPIClient, id string) (*types.Container, error) {
	containers, err := cliContainerList(cli)
	if err != nil {
		return nil, err
	}
	for i, container := range containers {
		if container.ID == id {
			return &containers[i], nil
		}
	}
	return nil, nil
}

// FindContainerByShortID searches container by short id.
func FindContainerByShortID(cli client.ContainerAPIClient, id string) (*types.Container, error) {
	containers, err := cliContainerList(cli)
	if err != nil {
		return nil, err
	}
	for i, container := range containers {
		if strings.HasPrefix(container.ID, id) {
			return &containers[i], nil
		}
	}
	return nil, nil
}

// FindContainerByName searches container by name.
func FindContainerByName(cli client.ContainerAPIClient, name string) (*types.Container, error) {
	containers, err := cliContainerList(cli)
	if err != nil {
		return nil, err
	}
	val := "/" + name
	for i, container := range containers {
		for _, item := range container.Names {
			if item == val {
				return &containers[i], nil
			}
		}
	}
	return nil, nil
}

// FindContainersByImageID searches containers by image id.
func FindContainersByImageID(cli client.ContainerAPIClient, imageID string) ([]*types.Container, error) {
	containers, err := cliContainerList(cli)
	if err != nil {
		return nil, err
	}
	var ret []*types.Container
	for i, container := range containers {
		if container.ImageID == imageID {
			ret = append(ret, &containers[i])
		}
	}
	return ret, nil
}

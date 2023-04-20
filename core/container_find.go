package core

import (
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// GetContainerName returns container name.
//
// Takes first of container names and trims leading  "/" character.
//
//	GetContainerName(&container) -> "my-container"
func GetContainerName(container *types.Container) string {
	if len(container.Names) > 0 {
		return container.Names[0][1:]
	}
	return ""
}

// GetContainerShortID return short variant of container id.
//
// Takes first 12 characters of identifier.
// Produces same value as the first columns of `docker ps` output.
//
//	GetContainerShortID(&container) -> "12345678abcd"
func GetContainerShortID(container *types.Container) string {
	return container.ID[:shortIDLength]
}

// FindContainerByID searches container by id.
//
// `id` is a full (64 characters) identifier.
//
//	FindContainerByID(cli, "<guid>") -> &container
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
	return nil, &ContainerNotFoundError{id}
}

// FindContainerByShortID searches container by short id.
//
// Uses `strings.HasPrefix` to compare container identifiers.
// Any substring of actual identifier can be passed.
//
//	FindContainerByShortID(cli, "1234") -> &container
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
	return nil, &ContainerNotFoundError{id}
}

// FindContainerByName searches container by name.
//
// Adds leading "/" character to passed value.
//
//	FindContainerByName(cli, "my-container") -> &container
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
	return nil, &ContainerNotFoundError{name}
}

// FindContainersByImageID searches containers by image id.
//
// `imageID` is a full image identifier - 64 characters with leading "sha256:".
//
//	FindContainersByImageID(cli, "sha256:<guid>") -> &container
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

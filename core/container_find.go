package core

import (
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// FindContainerByID searches container by id.
//
// `id` is a full (64 characters) identifier.
//
//	FindContainerByID(cli, "<guid>") -> container
func FindContainerByID(cli client.ContainerAPIClient, id string) (Container, error) {
	containers, err := cliContainerList(cli)
	if err != nil {
		return nil, err
	}
	for i, container := range containers {
		if container.ID == id {
			return makeContainer(&containers[i]), nil
		}
	}
	return nil, nil
}

// FindContainerByShortID searches container by short id.
//
// Uses `strings.HasPrefix` to compare container identifiers.
// Any substring of actual identifier can be passed.
//
//	FindContainerByShortID(cli, "1234") -> container
func FindContainerByShortID(cli client.ContainerAPIClient, id string) (Container, error) {
	containers, err := cliContainerList(cli)
	if err != nil {
		return nil, err
	}
	for i, container := range containers {
		if strings.HasPrefix(container.ID, id) {
			return makeContainer(&containers[i]), nil
		}
	}
	return nil, nil
}

// FindContainerByName searches container by name.
//
// Adds leading "/" character to passed value.
//
//	FindContainerByName(cli, "my-container") -> container
func FindContainerByName(cli client.ContainerAPIClient, name string) (Container, error) {
	containers, err := cliContainerList(cli)
	if err != nil {
		return nil, err
	}
	targetName := "/" + name
	for i, container := range containers {
		for _, name := range container.Names {
			if name == targetName {
				return makeContainer(&containers[i]), nil
			}
		}
	}
	return nil, nil
}

// FindContainersByImageID searches containers by image id.
//
// `imageID` is a full image identifier - 64 characters with leading "sha256:".
//
//	FindContainersByImageID(cli, "sha256:<guid>") -> []container
func FindContainersByImageID(cli client.ContainerAPIClient, imageID string) ([]Container, error) {
	containers, err := cliContainerList(cli)
	if err != nil {
		return nil, err
	}
	var objects []*types.Container
	for i, container := range containers {
		if container.ImageID == imageID {
			objects = append(objects, &containers[i])
		}
	}
	return TransformSlice(objects, makeContainer), nil
}

package core

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ListAllContainerIDs returns all container ids.
//
//	ListAllContainerIDs(cli) -> []string
func ListAllContainerIDs(cli client.ContainerAPIClient) ([]string, error) {
	containers, err := cliContainerList(cli)
	if err != nil {
		return nil, err
	}
	return TransformSlice(containers, func(container types.Container) string {
		return container.ID
	}), nil
}

package core

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ListAllImageIDs returns all images ids.
//
//	ListAllImageIDs(cli) -> []string
func ListAllImageIDs(cli client.ImageAPIClient) ([]string, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return nil, err
	}
	return TransformSlice(images, func(image types.ImageSummary) string {
		return image.ID
	}), nil
}

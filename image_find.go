package containerator

import (
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// GetImageName returns friendly image name.
func GetImageName(image *types.ImageSummary) string {
	if len(image.RepoTags) > 0 {
		return image.RepoTags[0]
	}
	return ""
}

// FindImageByID searches image by id.
func FindImageByID(cli client.ImageAPIClient, id string) (*types.ImageSummary, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return nil, err
	}
	for i, image := range images {
		if image.ID == id {
			return &images[i], nil
		}
	}
	return nil, nil
}

// FindImageByRepoTag searches image by repo tag.
func FindImageByRepoTag(cli client.ImageAPIClient, repoTag string) (*types.ImageSummary, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return nil, err
	}
	for i, image := range images {
		for _, item := range image.RepoTags {
			if item == repoTag {
				return &images[i], nil
			}
		}
	}
	return nil, nil
}

// FindImagesByRepo searches images by repo.
func FindImagesByRepo(cli client.ImageAPIClient, repo string) ([]*types.ImageSummary, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return nil, err
	}
	var ret []*types.ImageSummary
	for i, image := range images {
		for _, repoTag := range image.RepoTags {
			if strings.Split(repoTag, ":")[0] == repo {
				ret = append(ret, &images[i])
				break
			}
		}
	}
	return ret, nil
}

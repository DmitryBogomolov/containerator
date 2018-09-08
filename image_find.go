package containerator

import (
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// GetImageFullName returns full image name.
func GetImageFullName(image *types.ImageSummary) string {
	if len(image.RepoTags) > 0 {
		return image.RepoTags[0]
	}
	return ""
}

func extractRepo(repoTag string) string {
	return strings.SplitN(repoTag, ":", 2)[0]
}

// GetImageName returns friendly image name.
func GetImageName(image *types.ImageSummary) string {
	return extractRepo(GetImageFullName(image))
}

const (
	imageIDPrefix = "sha256:"
	shortIDLength = 12
)

// GetImageShortID returns short image id.
func GetImageShortID(image *types.ImageSummary) string {
	return image.ID[len(imageIDPrefix) : len(imageIDPrefix)+shortIDLength]
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

// FindImageByShortID searches image by short id.
func FindImageByShortID(cli client.ImageAPIClient, id string) (*types.ImageSummary, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return nil, err
	}
	val := imageIDPrefix + id
	for i, image := range images {
		if strings.HasPrefix(image.ID, val) {
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
	val := repoTag
	if extractRepo(val) == val {
		val = val + ":latest"
	}
	for i, image := range images {
		for _, item := range image.RepoTags {
			if item == val {
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

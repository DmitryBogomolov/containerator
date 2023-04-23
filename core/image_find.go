package core

import (
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// FindImageByID searches image by full id.
//
// `id` is a full (64 characters) identifier with "sha256:" prefix.
//
//	FindImageByID(cli, "sha256:<guid>") -> image
func FindImageByID(cli client.ImageAPIClient, id string) (Image, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return nil, err
	}
	for i, image := range images {
		if image.ID == id {
			return makeImage(&images[i]), nil
		}
	}
	return nil, &ImageNotFoundError{id}
}

// FindImageByShortID searches image by short id.
//
// Adds "sha256:" prefix and uses `string.HasPrefix` to compare identifiers.
// Any substring of actual identifier can be passed.
//
//	FindImageByShortID(cli, "1234") -> &image
func FindImageByShortID(cli client.ImageAPIClient, id string) (Image, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return nil, err
	}
	val := imageIDPrefix + id
	for i, image := range images {
		if strings.HasPrefix(image.ID, val) {
			return makeImage(&images[i]), nil
		}
	}
	return nil, &ImageNotFoundError{id}
}

func normalizeImageName(name string) string {
	if strings.Contains(name, ":") {
		return name
	}
	return name + ":" + imageTagLatest
}

// FindImageByName searches image by repo:tag.
//
// If tag is not provided then ":latest" is assumed.
//
//	FindImageByName(cli, "my-image:1") -> image
func FindImageByName(cli client.ImageAPIClient, name string) (Image, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return nil, err
	}
	targetName := normalizeImageName(name)
	for i, image := range images {
		for _, repoTag := range image.RepoTags {
			if repoTag == targetName {
				return makeImage(&images[i]), nil
			}
		}
	}
	return nil, &ImageNotFoundError{name}
}

// FindAllImagesByName searches images by repo.
//
// Finds all images with matching repository name.
//
//	FindAllImagesByName(cli, "my-image") -> []image
func FindAllImagesByName(cli client.ImageAPIClient, repo string) ([]Image, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return nil, err
	}
	var objects []*types.ImageSummary
	for i, image := range images {
		for _, repoTag := range image.RepoTags {
			if takeImageName(repoTag) == repo {
				objects = append(objects, &images[i])
				break
			}
		}
	}
	return TransformSlice(objects, makeImage), nil
}

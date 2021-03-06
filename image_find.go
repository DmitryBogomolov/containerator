package containerator

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	tagLatest     = "latest"
	imageIDPrefix = "sha256:"
	shortIDLength = 12
)

/*
GetImageFullName returns full image name.

Takes first of repository-tag pairs.

	GetImageFullName(&image) -> "my-image:1"
*/
func GetImageFullName(image *types.ImageSummary) string {
	if len(image.RepoTags) > 0 {
		return image.RepoTags[0]
	}
	return ""
}

/*
SplitImageNameTag splits full image name into repository and tag parts.

	SplitImageNameTag("my-image:1") -> "my-image", "1"
*/
func SplitImageNameTag(fullName string) (name string, tag string) {
	if fullName == "" {
		panic("empty string")
	}
	items := strings.SplitN(fullName, ":", 2)
	name = items[0]
	tag = tagLatest
	if len(items) > 1 {
		tag = items[1]
	}
	return
}

/*
JoinImageNameTag joins image name and tag into full image name.

	JoinImageNameTag("my-image", "1") -> "my-image:1"
*/
func JoinImageNameTag(name string, tag string) string {
	if name == "" {
		panic("empty string")
	}
	if tag == "" {
		tag = tagLatest
	}
	return name + ":" + tag
}

/*
GetImageShortID returns short image id.

Removes "sha256:" prefix and takes first 12 characters of identifier,

	GetImageShortID(&image) -> "12345678abcd"
*/
func GetImageShortID(image *types.ImageSummary) string {
	return image.ID[len(imageIDPrefix) : len(imageIDPrefix)+shortIDLength]
}

// ImageNotFoundError indicates that image with specified ID or full name is not found.
type ImageNotFoundError struct {
	Image string
}

func (e *ImageNotFoundError) Error() string {
	return fmt.Sprintf("image '%s' is not found", e.Image)
}

/*
FindImageByID searches image by id.

`id` is a full (64 characters) identifier with "sha256:" prefix.

	FindImageByID(cli, "sha256:<guid>") -> &image
*/
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
	return nil, &ImageNotFoundError{id}
}

/*
FindImageByShortID searches image by short id.

Adds "sha256:" prefix and uses `string.HasPrefix` to compare identifiers.
Any substring of actual identifier can be passed.

	FindImageByShortID(cli, "1234") -> &image
*/
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
	return nil, &ImageNotFoundError{id}
}

/*
FindImageByRepoTag searches image by repo tag.

`repoTag` contains repository and tag separated by ":".
If `repoTag` does not contain ":" then ":latest" postfix is added.

	FindImageByRepoTag(cli, "my-image:1") -> &image
*/
func FindImageByRepoTag(cli client.ImageAPIClient, repoTag string) (*types.ImageSummary, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return nil, err
	}
	val := JoinImageNameTag(SplitImageNameTag(repoTag))
	for i, image := range images {
		for _, item := range image.RepoTags {
			if item == val {
				return &images[i], nil
			}
		}
	}
	return nil, &ImageNotFoundError{repoTag}
}

/*
FindImagesByRepo searches images by repo.

Finds all images with matching repository.

	FindImagesByRepo(cli, "my-image") -> []&image
*/
func FindImagesByRepo(cli client.ImageAPIClient, repo string) ([]*types.ImageSummary, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return nil, err
	}
	var ret []*types.ImageSummary
	for i, image := range images {
		for _, repoTag := range image.RepoTags {
			name, _ := SplitImageNameTag(repoTag)
			if name == repo {
				ret = append(ret, &images[i])
				break
			}
		}
	}
	return ret, nil
}

/*
GetImagesTags extracts tags from list of images.

	GetImagesTags(images) -> []string{"1", "2", "3"}
*/
func GetImagesTags(images []*types.ImageSummary) []string {
	items := make([]string, len(images))
	for i, image := range images {
		_, tag := SplitImageNameTag(GetImageFullName(image))
		items[i] = tag
	}
	return items
}

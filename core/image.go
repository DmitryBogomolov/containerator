package core

import (
	"strings"

	"github.com/docker/docker/api/types"
)

const (
	imageTagLatest     = "latest"
	imageIDPrefix      = "sha256:"
	imageShortIDLength = 12
)

// Image provides information about image.
type Image interface {
	ID() string
	ShortID() string
	FullName() string
	Name() string
	Tag() string
}

type _Image struct {
	object *types.ImageSummary
}

// ID returns image id.
//
//	image.ID() -> "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
func (image *_Image) ID() string {
	return image.object.ID
}

// ShortID returns short image id.
//
// Removes "sha256:" prefix and takes first 12 characters of identifier,
//
//	image.ShortID() -> "0123456789ab"
func (image *_Image) ShortID() string {
	id := image.ID()
	return id[len(imageIDPrefix) : len(imageIDPrefix)+imageShortIDLength]
}

func takeFirst[T any](list []T) T {
	if len(list) > 0 {
		return list[0]
	}
	var zero T
	return zero
}

// FullName returns full image name.
//
// Takes first of repo:tag pairs.
//
//	image.FullName() -> "my-image:1"
func (image *_Image) FullName() string {
	return takeFirst(image.object.RepoTags)
}

// Name returns image name.
//
// Takes first part of full name.
//
//	image.Name() -> "my-image"
func (image *_Image) Name() string {
	return takeImageName(image.FullName())
}

// Tag returns image tag.
//
// Takes second part of full name.
//
//	image.Tag() -> "1"
func (image *_Image) Tag() string {
	return takeImageTag(image.FullName())
}

func makeImage(object *types.ImageSummary) Image {
	return &_Image{object}
}

func takeImageName(repoTag string) string {
	parts := strings.Split(repoTag, ":")
	return takeFirst(parts)
}

func takeImageTag(repoTag string) string {
	parts := strings.Split(repoTag, ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

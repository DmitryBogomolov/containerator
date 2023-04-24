package core

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func testImage(id string, repoTags ...string) Image {
	object := &types.ImageSummary{
		ID:       id,
		RepoTags: repoTags,
	}
	return &_Image{object}
}

func TestImage_ID(t *testing.T) {
	image := testImage("sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	assert.Equal(t, "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", image.ID())
}

func TestImage_ShortID(t *testing.T) {
	image := testImage("sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	assert.Equal(t, "0123456789ab", image.ShortID())
}

func TestImage_FullName(t *testing.T) {
	image1 := testImage("")
	assert.Equal(t, "", image1.FullName(), "no repo tags")

	image2 := testImage("", "a:1", "b:2")
	assert.Equal(t, "a:1", image2.FullName(), "take first repo tag")
}

func TestImage_Name(t *testing.T) {
	image1 := testImage("")
	assert.Equal(t, "", image1.Name(), "no repo tags")

	image2 := testImage("", "a:1", "b:2")
	assert.Equal(t, "a", image2.Name(), "take first repo tag")
}

func TestImage_Tag(t *testing.T) {
	image1 := testImage("")
	assert.Equal(t, "", image1.Tag(), "no repo tags")

	image2 := testImage("", "a:1", "b:2")
	assert.Equal(t, "1", image2.Tag(), "take first repo tag")
}

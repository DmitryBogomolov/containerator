package core

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func testContainer(id string, imageID string, names ...string) Container {
	object := &types.Container{
		ID:      id,
		ImageID: imageID,
		Names:   names,
	}
	return &_Container{object}
}

func TestContainer_ID(t *testing.T) {
	container := testContainer("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", "")
	assert.Equal(t, "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", container.ID())
}

func TestContainer_ShortID(t *testing.T) {
	container := testContainer("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", "")
	assert.Equal(t, "0123456789ab", container.ShortID())
}

func TestContainer_Name(t *testing.T) {
	container1 := testContainer("", "")
	assert.Equal(t, "", container1.Name(), "empty name")

	container2 := testContainer("", "", "/c1", "/c2")
	assert.Equal(t, "c1", container2.Name(), "take first name")
}

func TestContainer_ImageID(t *testing.T) {
	container := testContainer("", "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	assert.Equal(t, "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", container.ImageID())
}

package core

import (
	"strings"

	"github.com/docker/docker/api/types"
)

const (
	containerShortIDLength = 12
)

type Container interface {
	ID() string
	ShortID() string
	Name() string
	ImageID() string
	State() string
}

type _Container struct {
	object *types.Container
}

// ID returns container id.
//
//	container.ID() -> "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
func (container *_Container) ID() string {
	return container.object.ID
}

// ShortID returns short variant of container id.
//
// Takes first 12 characters of identifier.
// Produces same value as the first columns of `docker ps` output.
//
//	container.ShortID() -> "0123456789ab"
func (container *_Container) ShortID() string {
	id := container.ID()
	return id[:containerShortIDLength]
}

// Name returns container name.
//
// Takes first of container names and trims leading  "/" character.
//
//	container.Name() -> "my-container"
func (container *_Container) Name() string {
	return strings.TrimLeft(takeFirst(container.object.Names), "/")
}

// ImageID returns container imageid.
//
//	container.ImageID() -> "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
func (container *_Container) ImageID() string {
	return container.object.ImageID
}

// State returns container state.
//
//	container.State() -> "running"
func (container *_Container) State() string {
	return container.object.State
}

func makeContainer(object *types.Container) Container {
	return &_Container{object}
}

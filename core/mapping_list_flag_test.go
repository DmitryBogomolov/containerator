package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMappingListFlagValue(t *testing.T) {
	obj := NewMappingListFlag("-", true)
	actual, ok := obj.(*_MappingListFlag)
	assert.True(t, ok)
	assert.Equal(t, "-", actual.separator)
	assert.Equal(t, true, actual.allowOne)
}

func TestMappingListFlagValue(t *testing.T) {
	t.Run("String / empty", func(t *testing.T) {
		obj := _MappingListFlag{}
		assert.Equal(t, "[]", obj.String())
	})

	t.Run("String / list", func(t *testing.T) {
		obj := _MappingListFlag{mappings: []Mapping{
			{Source: "a"},
			{Source: "m", Target: "n"},
		}}
		assert.Equal(t, "[{a } {m n}]", obj.String())
	})

	t.Run("Set", func(t *testing.T) {
		obj := _MappingListFlag{separator: "#", allowOne: true}

		obj.Set("a")
		assert.Equal(t, []Mapping{
			{Source: "a"},
		}, obj.mappings)

		obj.Set("x#y")
		assert.Equal(t, []Mapping{
			{Source: "a"},
			{Source: "x", Target: "y"},
		}, obj.mappings)
	})

	t.Run("Get", func(t *testing.T) {
		list := []Mapping{{}, {}}
		obj := _MappingListFlag{mappings: list}

		assert.Equal(t, list, obj.Get())
	})
}

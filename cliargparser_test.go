package containerator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMappingListVar(t *testing.T) {
	obj := NewMappingListVar("-", true)
	actual, ok := obj.(*_MappingListVar)
	assert.True(t, ok)
	assert.Equal(t, "-", actual.sep)
	assert.Equal(t, true, actual.allowOne)
}

func TestMappingListVar(t *testing.T) {
	t.Run("String / empty", func(t *testing.T) {
		obj := _MappingListVar{}
		assert.Equal(t, "[]", obj.String())
	})

	t.Run("String / list", func(t *testing.T) {
		obj := _MappingListVar{list: []Mapping{
			{Source: "a"},
			{Source: "m", Target: "n"},
		}}
		assert.Equal(t, "[{a } {m n}]", obj.String())
	})

	t.Run("Set", func(t *testing.T) {
		obj := _MappingListVar{sep: "#", allowOne: true}

		obj.Set("a")
		assert.Equal(t, []Mapping{
			{Source: "a"},
		}, obj.list)

		obj.Set("x#y")
		assert.Equal(t, []Mapping{
			{Source: "a"},
			{Source: "x", Target: "y"},
		}, obj.list)
	})

	t.Run("Get", func(t *testing.T) {
		list := []Mapping{{}, {}}
		obj := _MappingListVar{list: list}

		assert.Equal(t, list, obj.Get())
	})
}

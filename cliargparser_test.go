package containerator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMappingListVar(t *testing.T) {
	obj := NewMappingListVar("-", true)
	actual, ok := obj.(*mappingListVar)
	assert.True(t, ok)
	assert.Equal(t, "-", actual.sep)
	assert.Equal(t, true, actual.allowOne)
}

func TestMappingListVar(t *testing.T) {
	t.Run("String / empty", func(t *testing.T) {
		obj := mappingListVar{}
		assert.Equal(t, "[]", obj.String())
	})

	t.Run("String / list", func(t *testing.T) {
		obj := mappingListVar{list: []Mapping{
			Mapping{Source: "a"},
			Mapping{Source: "m", Target: "n"},
		}}
		assert.Equal(t, "[{a } {m n}]", obj.String())
	})

	t.Run("Set", func(t *testing.T) {
		obj := mappingListVar{sep: "#", allowOne: true}

		obj.Set("a")
		assert.Equal(t, []Mapping{
			Mapping{Source: "a"},
		}, obj.list)

		obj.Set("x#y")
		assert.Equal(t, []Mapping{
			Mapping{Source: "a"},
			Mapping{Source: "x", Target: "y"},
		}, obj.list)
	})

	t.Run("Get", func(t *testing.T) {
		list := []Mapping{Mapping{}, Mapping{}}
		obj := mappingListVar{list: list}

		assert.Equal(t, list, obj.Get())
	})
}

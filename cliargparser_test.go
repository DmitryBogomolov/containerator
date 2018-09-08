package containerator

import (
	"testing"
)

func TestNewMappingListVar(t *testing.T) {
	obj := NewMappingListVar("-", true)
	actual, ok := obj.(*mappingListVar)
	assertEqual(t, ok, true, "type asserted")
	assertEqual(t, actual.sep, "-", "sep")
	assertEqual(t, actual.allowOne, true, "allowOne")
}

func TestMappingListVar(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		obj := mappingListVar{}
		assertEqual(t, obj.String(), "[]", "String")

		list := []Mapping{
			Mapping{Source: "a"},
			Mapping{Source: "m", Target: "n"},
		}
		obj.list = list
		assertEqual(t, obj.String(), "[{a } {m n}]", "String")
	})

	t.Run("Set", func(t *testing.T) {
		obj := mappingListVar{sep: "#", allowOne: true}

		obj.Set("a")
		assertEqual(t, obj.list[0], Mapping{Source: "a"}, "Set")

		obj.Set("x#y")
		assertEqual(t, obj.list[1], Mapping{Source: "x", Target: "y"}, "Set")
	})

	t.Run("Get", func(t *testing.T) {
		list := []Mapping{Mapping{}, Mapping{}}
		obj := mappingListVar{list: list}

		assertEqual(t, len(obj.Get()), 2, "list length")
		assertEqual(t, obj.Get()[0], list[0], "list item 1")
		assertEqual(t, obj.Get()[1], list[1], "list item 2")
	})
}

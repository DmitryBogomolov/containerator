package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestMapping(t *testing.T) {
	t.Run("JSONMarshal", func(t *testing.T) {
		mapping := Mapping{Source: "src", Target: "dst"}

		bytes, err := json.Marshal(mapping)
		data := string(bytes)

		assert.NoError(t, err)
		assert.Equal(t, "{\"src\":\"dst\"}", data)
	})

	t.Run("JSONUnmarshal", func(t *testing.T) {
		data := " { \"src\": \"dst\" } "
		var mapping Mapping

		err := json.Unmarshal([]byte(data), &mapping)

		assert.NoError(t, err)
		assert.Equal(t, Mapping{Source: "src", Target: "dst"}, mapping)
	})

	t.Run("YAMLMarshal", func(t *testing.T) {
		mapping := Mapping{Source: "src", Target: "dst"}

		bytes, err := yaml.Marshal(mapping)
		data := string(bytes)

		assert.NoError(t, err)
		assert.Equal(t, "src: dst\n", data)
	})

	t.Run("YAMLUnmarshal", func(t *testing.T) {
		data := "src: dst"
		var mapping Mapping

		err := yaml.Unmarshal([]byte(data), &mapping)

		assert.NoError(t, err)
		assert.Equal(t, Mapping{Source: "src", Target: "dst"}, mapping)
	})
}

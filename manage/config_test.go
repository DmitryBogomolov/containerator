package manage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	t.Run("Read file", func(t *testing.T) {
		ioutil.WriteFile("test.yaml", []byte("image_repo: my-image\nmodes: ['a', 'b']\n"), os.ModePerm)
		defer os.Remove("test.yaml")

		config, err := ReadConfig("test.yaml")

		assert.NoError(t, err, "error")
		assert.Equal(t, &Config{
			ImageName: "my-image",
			Modes:     []string{"a", "b"},
		}, config, "config")
	})

	t.Run("No file", func(t *testing.T) {
		config, err := ReadConfig("test.yaml")

		assert.Error(t, err, "error")
		assert.Equal(t, (err.(*os.PathError)).Path, "test.yaml", "error data")
		assert.Equal(t, (*Config)(nil), config, "config")
	})
}

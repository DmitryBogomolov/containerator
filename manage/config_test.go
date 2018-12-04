package manage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	ioutil.WriteFile("test.yaml", []byte("image_repo: my-image\nmodes: ['a', 'b']\n"), os.ModePerm)
	defer os.Remove("test.yaml")

	config, err := ReadConfig("test.yaml")

	assert.Equal(t, nil, err, "error")
	assert.Equal(t, Config{
		ImageRepo: "my-image",
		Modes:     []string{"a", "b"},
	}, *config, "config")
}

func TestReadConfigNoFile(t *testing.T) {
	config, err := ReadConfig("test.yaml")

	assert.EqualError(t, err, "open test.yaml: no such file or directory", "error")
	assert.Equal(t, (*Config)(nil), config, "config")
}

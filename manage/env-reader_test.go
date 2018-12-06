package manage

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvFileReader(t *testing.T) {
	ioutil.WriteFile("env.list", []byte("common list"), os.ModePerm)
	defer os.Remove("env.list")
	ioutil.WriteFile("m1.list", []byte("m1-mode list"), os.ModePerm)
	defer os.Remove("m1.list")

	read := func(reader io.Reader) string {
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}

	t.Run("Empty mode", func(t *testing.T) {
		reader, err := GetEnvFileReader(".", "")

		assert.NoError(t, err, "error")
		assert.Equal(t, "common list", read(reader), "data")
	})

	t.Run("Valid mode", func(t *testing.T) {
		reader, err := GetEnvFileReader(".", "m1")

		assert.NoError(t, err, "error")
		assert.Equal(t, "m1-mode list", read(reader), "data")
	})

	t.Run("Not valid mode", func(t *testing.T) {
		reader, err := GetEnvFileReader(".", "m2")

		assert.NoError(t, err, "error")
		assert.Equal(t, "common list", read(reader), "data")
	})

	t.Run("No files", func(t *testing.T) {
		reader, err := GetEnvFileReader(".test", "m1")

		assert.Error(t, err, "error")
		assert.Equal(t, err, ErrNoEnvFile, "error data")
		assert.Equal(t, nil, reader, "data")
	})
}

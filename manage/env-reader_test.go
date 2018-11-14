package manage

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvFileReader(t *testing.T) {
	ioutil.WriteFile("env.list", []byte("c1"), os.ModePerm)
	ioutil.WriteFile("m1.list", []byte("c2"), os.ModePerm)
	defer func() {
		os.Remove("env.list")
		os.Remove("m1.list")
	}()

	read := func(reader io.Reader) string {
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}

	reader, err := GetEnvFileReader(".", "")
	assert.Equal(t, nil, err, "error")
	assert.Equal(t, "c1", read(reader), "data")

	reader, err = GetEnvFileReader(".", "m1")
	assert.Equal(t, nil, err, "error")
	assert.Equal(t, "c2", read(reader), "data")

	reader, err = GetEnvFileReader(".", "m2")
	assert.Equal(t, nil, err, "error")
	assert.Equal(t, "c1", read(reader), "data")
}

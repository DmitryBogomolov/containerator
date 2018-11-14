package manage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func isFileExist(file string) bool {
	if _, err := os.Stat(file); os.IsExist(err) {
		return true
	}
	return false
}

func getEnvFileName(dir string, mode string) string {
	return filepath.Join(dir, fmt.Sprintf("%s.list", mode))
}

func selectEnvFile(dir string, mode string) string {
	name := getEnvFileName(dir, mode)
	if isFileExist(name) {
		return name
	}
	name = getEnvFileName(dir, "env")
	if isFileExist(name) {
		return name
	}
	return ""
}

// ErrNoEnvFile is returned when neither <mode>.list no env.list is found in a directory.
var ErrNoEnvFile = errors.New("env file is not found")

/*
GetEnvFileReader creates *EnvReader* searching directory with specified *mode*.

	GetEnvFileReader("/path/to/dir", mode) -> reader, err
*/
func GetEnvFileReader(configPath string, mode string) (io.Reader, error) {
	dir, _ := filepath.Abs(filepath.Dir(configPath))
	envFileName := selectEnvFile(dir, mode)
	if envFileName == "" {
		return nil, ErrNoEnvFile
	}
	data, err := ioutil.ReadFile(envFileName)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

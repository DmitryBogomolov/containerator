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

func doesFileExist(file string) bool {
	stat, err := os.Stat(file)
	return err == nil && !stat.IsDir()
}

func getEnvFileName(dir string, mode string) string {
	return filepath.Join(dir, fmt.Sprintf("%s.list", mode))
}

func selectEnvFile(dir string, mode string) string {
	name := getEnvFileName(dir, mode)
	if doesFileExist(name) {
		return name
	}
	name = getEnvFileName(dir, "env")
	if doesFileExist(name) {
		return name
	}
	return ""
}

// ErrNoEnvFile is returned when neither <mode>.list no env.list is found in a directory.
var ErrNoEnvFile = errors.New("env file is not found")

// GetEnvFileReader creates *EnvReader* searching directory with specified *mode*.
//
//	GetEnvFileReader("/path/to/dir", mode) -> reader, err
func GetEnvFileReader(dirPath string, mode string) (io.Reader, error) {
	envFileName := selectEnvFile(dirPath, mode)
	if envFileName == "" {
		return nil, ErrNoEnvFile
	}
	data, err := ioutil.ReadFile(envFileName)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

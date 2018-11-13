package manage

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func findIndex(item string, list []string) int {
	for i, obj := range list {
		if obj == item {
			return i
		}
	}
	return -1
}

func selectMode(modeOption string, conf *Config) (string, int, error) {
	index := findIndex(modeOption, conf.Modes)
	if index >= 0 {
		return modeOption, index, nil
	}
	if modeOption == "" && len(conf.Modes) == 0 {
		return modeOption, 0, nil
	}
	return "", 0, fmt.Errorf("mode '%s' is not valid", modeOption)
}

func getContainerName(imageRepo string, mode string) string {
	name := imageRepo
	if mode != "" {
		name += "-" + mode
	}
	return name
}

func findImage(cli client.ImageAPIClient, imageRepo string, imageTag string) (*types.ImageSummary, error) {
	if imageTag != "" {
		repoTag := imageRepo + ":" + imageTag
		item, err := containerator.FindImageByRepoTag(cli, repoTag)
		if err != nil {
			return nil, err
		}
		if item == nil {
			return nil, fmt.Errorf("no '%s' image", repoTag)
		}
		return item, nil
	}
	list, err := containerator.FindImagesByRepo(cli, imageRepo)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("no '%s' images", imageRepo)
	}
	return list[0], nil
}

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
	return strings.NewReader(string(data)), nil
}

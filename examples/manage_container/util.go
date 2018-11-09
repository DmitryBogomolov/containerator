package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
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

func selectMode(modeOption string, config *config) (string, int, error) {
	index := findIndex(modeOption, config.Modes)
	if index >= 0 {
		return modeOption, index, nil
	}
	if modeOption == "" && len(config.Modes) == 0 {
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
	return path.Join(dir, fmt.Sprintf("%s.list", mode))
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

func getEnvFileReader(configPath string, mode string) io.Reader {
	dir, _ := filepath.Abs(filepath.Dir(configPath))
	envFileName := selectEnvFile(dir, mode)
	if envFileName == "" {
		log.Println("env file is not found")
		return nil
	}
	data, err := ioutil.ReadFile(envFileName)
	if err != nil {
		log.Printf("failed to read env file: %+v", err)
		return nil
	}
	return strings.NewReader(string(data))
}

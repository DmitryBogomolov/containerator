package manage

import (
	"fmt"
	"strconv"

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

func getContainerName(conf *Config, mode string) string {
	name := conf.ContainerName
	if name == "" {
		name = conf.ImageRepo
	}
	if mode != "" {
		name += "-" + mode
	}
	return name
}

func findImage(cli client.ImageAPIClient, imageRepo string, imageTag string) (*types.ImageSummary, error) {
	if imageTag != "" {
		repoTag := imageRepo + ":" + imageTag
		item, err := containerator.FindImageByRepoTag(cli, repoTag)
		if _, ok := err.(*containerator.ImageNotFoundError); ok {
			err = fmt.Errorf("no '%s' image (%v)", repoTag, err)
		}
		return item, err
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

func buildContainerOptions(conf *Config, imageName string, containerName string,
	modeIndex int) *containerator.RunContainerOptions {
	ret := containerator.RunContainerOptions{
		Image:         imageName,
		Name:          containerName,
		RestartPolicy: containerator.RestartAlways,
		Network:       conf.Network,
		Volumes:       conf.Volumes,
		Env:           conf.Env,
	}
	if len(conf.Ports) > 0 {
		basePort := int(conf.BasePort) + int(conf.PortOffset)*modeIndex
		ports := make([]containerator.Mapping, len(conf.Ports))
		for i, port := range conf.Ports {
			ports[i] = containerator.Mapping{
				Source: strconv.Itoa(basePort + i),
				Target: strconv.Itoa(int(port)),
			}
		}
		ret.Ports = ports
	}
	return &ret
}

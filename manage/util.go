package manage

import (
	"fmt"
	"strconv"

	"github.com/DmitryBogomolov/containerator/core"
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

// NotValidModeError indicates that specified mode is not found amoung config modes.
type NotValidModeError struct {
	mode  string
	modes []string
}

func (err *NotValidModeError) Error() string {
	return fmt.Sprintf("mode '%s' is not valid", err.mode)
}

// Mode returns mode.
func (err *NotValidModeError) Mode() string {
	return err.mode
}

// Modes returns config modes.
func (err *NotValidModeError) Modes() []string {
	return append([]string(nil), err.modes...)
}

func selectMode(modeOption string, conf *Config) (string, int, error) {
	index := findIndex(modeOption, conf.Modes)
	if index >= 0 {
		return modeOption, index, nil
	}
	if modeOption == "" && len(conf.Modes) == 0 {
		return modeOption, 0, nil
	}
	return "", 0, &NotValidModeError{modeOption, conf.Modes}
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
		item, err := core.FindImageByRepoTag(cli, repoTag)
		return item, err
	}
	list, err := core.FindImagesByRepo(cli, imageRepo)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil // TODO: core.ImageNotFound(imageRepo)
	}
	return list[0], nil
}

func buildContainerOptions(conf *Config, imageName string, containerName string,
	modeIndex int) *core.RunContainerOptions {
	ret := core.RunContainerOptions{
		Image:         imageName,
		Name:          containerName,
		RestartPolicy: core.RestartAlways,
		Network:       conf.Network,
		Volumes:       conf.Volumes,
		Env:           conf.Env,
	}
	if len(conf.Ports) > 0 {
		basePort := int(conf.BasePort) + int(conf.PortOffset)*modeIndex
		ports := make([]core.Mapping, len(conf.Ports))
		for i, port := range conf.Ports {
			ports[i] = core.Mapping{
				Source: strconv.Itoa(basePort + i),
				Target: strconv.Itoa(int(port)),
			}
		}
		ret.Ports = ports
	}
	return &ret
}

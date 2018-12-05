package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/DmitryBogomolov/containerator"

	"github.com/DmitryBogomolov/containerator/manage"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func parseBool(value string) bool {
	ret, _ := strconv.ParseBool(value)
	return ret
}

func getTag(cli client.ImageAPIClient, cont *types.Container) string {
	image, err := containerator.FindImageByID(cli, cont.ImageID)
	if err != nil {
		return fmt.Sprintf("Error(%+v)", err)
	}
	_, tag := containerator.SplitImageNameTag(containerator.GetImageFullName(image))
	return tag
}

func invokeManage(cli interface{}, configPath string, r *http.Request) (map[string]string, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	options := &manage.Options{
		Mode:   r.PostFormValue("mode"),
		Tag:    r.PostFormValue("tag"),
		Remove: parseBool(r.PostFormValue("remove")),
		Force:  parseBool(r.PostFormValue("force")),
	}
	config, err := manage.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}
	cont, err := manage.Manage(cli, config, options)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"name":  containerator.GetContainerName(cont),
		"image": config.ImageRepo,
		"tag":   getTag(cli.(client.ImageAPIClient), cont),
	}, nil
}

func getImageTags(cli interface{}, configPath string) ([]string, error) {
	config, err := manage.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}
	images, err := containerator.FindImagesByRepo(cli.(client.ImageAPIClient), config.ImageRepo)
	if err != nil {
		return nil, err
	}
	return containerator.GetImagesTags(images), nil
}

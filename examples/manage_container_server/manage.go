package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/DmitryBogomolov/containerator/manage"
	"github.com/docker/docker/client"
)

func parseBool(value string) bool {
	ret, _ := strconv.ParseBool(value)
	return ret
}

func getTag(cli client.ImageAPIClient, container core.Container) string {
	image, err := core.FindImageByID(cli, container.ImageID())
	if err != nil {
		return fmt.Sprintf("error(%+v)", err)
	}
	return image.Tag()
}

func parseRequestBody(body io.ReadCloser) *manage.Options {
	options := manage.Options{}
	var data map[string]any
	defer body.Close()
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return &options
	}
	if val, ok := data["postfix"]; ok {
		if postfix, ok := val.(string); ok {
			options.Postfix = postfix
		}
	}
	if val, ok := data["tag"]; ok {
		if tag, ok := val.(string); ok {
			options.Tag = tag
		}
	}
	if val, ok := data["force"]; ok {
		if force, ok := val.(bool); ok {
			options.Force = force
		}
	}
	if val, ok := data["remove"]; ok {
		if remove, ok := val.(bool); ok {
			options.Remove = remove
		}
	}
	return &options
}

func invokeManage(cli any, configPath string, r *http.Request) (map[string]any, error) {
	options := parseRequestBody(r.Body)
	config, err := manage.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}
	container, err := manage.RunContainer(cli, config, options)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"name":  container.Name(),
		"image": config.ImageName,
		"tag":   getTag(cli.(client.ImageAPIClient), container),
	}, nil
}

func getImageInfo(cli any, configPath string) (map[string]any, error) {
	config, err := manage.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}
	images, err := core.FindAllImagesByName(cli.(client.ImageAPIClient), config.ImageName)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"tags": core.TransformSlice(images, func(image core.Image) string {
			return image.Tag()
		}),
	}, nil
}

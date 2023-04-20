package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/DmitryBogomolov/containerator/core"

	"github.com/DmitryBogomolov/containerator/manage"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func parseBool(value string) bool {
	ret, _ := strconv.ParseBool(value)
	return ret
}

func getTag(cli client.ImageAPIClient, cont *types.Container) string {
	image, err := core.FindImageByID(cli, cont.ImageID)
	if err != nil {
		return fmt.Sprintf("Error(%+v)", err)
	}
	_, tag := core.SplitImageNameTag(core.GetImageFullName(image))
	return tag
}

func parseRequestBody(body io.ReadCloser) *manage.Options {
	ret := &manage.Options{}
	var data map[string]interface{}
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return ret
	}
	defer body.Close()
	if val, ok := data["mode"]; ok {
		if mode, ok := val.(string); ok {
			ret.Mode = mode
		}
	}
	if val, ok := data["tag"]; ok {
		if tag, ok := val.(string); ok {
			ret.Tag = tag
		}
	}
	if val, ok := data["force"]; ok {
		if force, ok := val.(bool); ok {
			ret.Force = force
		}
	}
	if val, ok := data["remove"]; ok {
		if remove, ok := val.(bool); ok {
			ret.Remove = remove
		}
	}
	return ret
}

func invokeManage(cli interface{}, configPath string, r *http.Request) (map[string]string, error) {
	options := parseRequestBody(r.Body)
	config, err := manage.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}
	options.GetEnvReader = func(mode string) (io.Reader, error) {
		reader, err := manage.GetEnvFileReader(filepath.Dir(configPath), mode)
		if err != nil {
			logger.Printf("failed to read env for '%s' (%+v)", configPath, err)
		}
		return reader, nil
	}
	cont, err := manage.Manage(cli, config, options)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"name":  core.GetContainerName(cont),
		"image": config.ImageName,
		"tag":   getTag(cli.(client.ImageAPIClient), cont),
	}, nil
}

func getImageInfo(cli interface{}, configPath string) (map[string]interface{}, error) {
	config, err := manage.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}
	images, err := core.FindImagesByRepo(cli.(client.ImageAPIClient), config.ImageName)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"modes": config.Modes,
		"tags":  core.GetImagesTags(images),
	}, nil
}

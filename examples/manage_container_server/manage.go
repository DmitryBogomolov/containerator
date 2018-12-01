package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/DmitryBogomolov/containerator"
	"github.com/DmitryBogomolov/containerator/manage"
	"github.com/docker/docker/client"
)

func parseBool(value string) bool {
	ret, _ := strconv.ParseBool(value)
	return ret
}

func invokeManage(cli client.CommonAPIClient, configPath string, r *http.Request) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", err
	}
	options := &manage.Options{
		Mode:   r.PostFormValue("mode"),
		Tag:    r.PostFormValue("tag"),
		Remove: parseBool(r.PostFormValue("remove")),
		Force:  parseBool(r.PostFormValue("force")),
	}
	config, err := manage.ReadConfig(configPath)
	if err != nil {
		return "", err
	}
	cont, err := manage.Manage(cli, config, options)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Container: %s", containerator.GetContainerName(cont)), nil
}

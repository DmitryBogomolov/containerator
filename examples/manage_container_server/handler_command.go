package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/DmitryBogomolov/containerator"
	"github.com/DmitryBogomolov/containerator/manage"
	"github.com/docker/docker/client"
)

type commandHandler struct {
	cache *projectsCache
	cli   client.CommonAPIClient
}

func newCommandHandler(cache *projectsCache) (*commandHandler, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &commandHandler{
		cache: cache,
		cli:   cli,
	}, nil
}

func (h *commandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !checkHTTPMethod(http.MethodPost, w, r) {
		return
	}
	name := strings.Replace(r.URL.Path, "/manage/", "", 1)
	item := h.cache.get(name)
	if item == nil {
		http.Error(w, fmt.Sprintf("'%s' is not found\n", name), http.StatusNotFound)
		return
	}
	body, err := invokeManage(h.cli, item.ConfigPath, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v\n", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, body)
}

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

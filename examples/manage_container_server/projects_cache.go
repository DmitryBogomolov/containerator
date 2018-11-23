package main

import (
	"crypto/sha256"
	"os"
)

type projectItem struct {
	ID         string
	Name       string
	ConfigPath string
}

type projectsCache struct {
	Dir   string
	Items []projectItem
}

func (obj *projectsCache) refresh() {
	obj.Items = []projectItem{
		newProjectItem("Project 1", "/at"),
		newProjectItem("Project 2", "/gv"),
	}
}

func (obj *projectsCache) get(name string) *projectItem {
	for i, item := range obj.Items {
		if item.Name == name {
			return &obj.Items[i]
		}
	}
	return nil
}

func newProjectsCache() *projectsCache {
	workDir, _ := os.Getwd()
	cache := &projectsCache{Dir: workDir}
	cache.refresh()
	return cache
}

func newProjectItem(name string, configPath string) projectItem {
	h := sha256.New()
	h.Write([]byte(configPath))
	id := string(h.Sum(nil))
	return projectItem{
		ID:         id,
		Name:       name,
		ConfigPath: configPath,
	}
}

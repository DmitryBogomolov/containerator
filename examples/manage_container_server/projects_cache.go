package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

type projectItem struct {
	Name       string
	ConfigPath string
}

type projectsCache struct {
	Dir      string
	Projects []projectItem
}

func (obj *projectsCache) refresh() {
	obj.Projects = []projectItem{
		newProjectItem("project-1", "/at"),
		newProjectItem("project-2", "/gv"),
	}
}

func (obj *projectsCache) get(name string) *projectItem {
	for i, item := range obj.Projects {
		if item.Name == name {
			return &obj.Projects[i]
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

func getProjectID(configPath string) string {
	h := sha256.New()
	h.Write([]byte(configPath))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash[:8]
}

func newProjectItem(name string, configPath string) projectItem {
	return projectItem{
		Name:       name,
		ConfigPath: configPath,
	}
}

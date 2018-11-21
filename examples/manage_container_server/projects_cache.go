package main

import "os"

type projectItem struct {
	Name       string
	ConfigPath string
}

type projectsCache struct {
	Dir   string
	Items []projectItem
}

func (obj *projectsCache) refresh() {
	obj.Items = []projectItem{
		projectItem{
			Name:       "Project 1",
			ConfigPath: "/at",
		},
		projectItem{
			Name:       "Project 2",
			ConfigPath: "/gv",
		},
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

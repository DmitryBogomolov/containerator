package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"time"

	"github.com/DmitryBogomolov/containerator/batcher"
)

type projectItem struct {
	Name       string
	configPath string
}

type projectsCache struct {
	workspace string
	Projects  []projectItem
	batcher   *batcher.Batcher
}

func collectProjects(pattern string) []string {
	for i := 0; i < 3; i++ {
		matches, err := filepath.Glob(pattern)
		if err == nil {
			return matches
		}
		logger.Printf("%+v", err)
		time.Sleep(3 * time.Second)
	}
	logger.Panicln("failed to refresh projects")
	return nil
}

func (obj *projectsCache) refreshCore() {
	matches := collectProjects(filepath.Join(obj.workspace, "*", "*.yaml"))
	items := make([]projectItem, len(matches))
	for i, m := range matches {
		items[i] = projectItem{
			Name:       filepath.Base(filepath.Dir(m)),
			configPath: m,
		}
	}
	obj.Projects = items
}

func (obj *projectsCache) refresh() {
	obj.batcher.Invoke()
}

func (obj *projectsCache) get(name string) (projectItem, error) {
	for i, item := range obj.Projects {
		if item.Name == name {
			return obj.Projects[i], nil
		}
	}
	var notFound projectItem
	return notFound, fmt.Errorf("project '%s' is not found", name)
}

func newProjectsCache(workspace string) *projectsCache {
	cache := &projectsCache{}
	cache.workspace = workspace
	cache.batcher = batcher.NewBatcher(cache.refreshCore)
	cache.refresh()
	return cache
}

func getProjectID(configPath string) string {
	h := sha256.New()
	h.Write([]byte(configPath))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash[:8]
}

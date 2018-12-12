package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/DmitryBogomolov/containerator/batcher"
)

const (
	refreshInterval = 10 * time.Second
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
	logger.Printf("Refresh\n  %s\n", strings.Join(obj.getProjectNames(), ", "))
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

func (obj *projectsCache) getProjectNames() []string {
	names := make([]string, len(obj.Projects))
	for i, p := range obj.Projects {
		names[i] = p.Name
	}
	return names
}

func (obj *projectsCache) refreshByInterval(ch <-chan time.Time) {
	for range ch {
		obj.refresh()
	}
}

func (obj *projectsCache) beginIntervalRefresh(d time.Duration) {
	ticker := time.NewTicker(d)
	go obj.refreshByInterval(ticker.C)
}

func newProjectsCache(workspace string) *projectsCache {
	cache := &projectsCache{}
	cache.workspace = workspace
	cache.batcher = batcher.NewBatcher(cache.refreshCore)
	cache.refresh()
	cache.beginIntervalRefresh(refreshInterval)
	return cache
}

func getProjectID(configPath string) string {
	h := sha256.New()
	h.Write([]byte(configPath))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash[:8]
}

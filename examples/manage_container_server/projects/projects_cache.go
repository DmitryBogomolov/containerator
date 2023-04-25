package projects

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/DmitryBogomolov/containerator/batcher"
	"github.com/DmitryBogomolov/containerator/examples/manage_container_server/logger"
)

const (
	refreshInterval = 10 * time.Second
)

type ProjectItem struct {
	Name       string
	ConfigPath string
}

type ProjectsCache struct {
	workspace string
	Projects  []ProjectItem
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
	logger.Panicf("failed to refresh projects\n")
	return nil
}

func (obj *ProjectsCache) refreshCore() {
	matches := collectProjects(filepath.Join(obj.workspace, "*", "*.yaml"))
	items := make([]ProjectItem, len(matches))
	for i, m := range matches {
		items[i] = ProjectItem{
			Name:       filepath.Base(filepath.Dir(m)),
			ConfigPath: m,
		}
	}
	obj.Projects = items
}

func (obj *ProjectsCache) Refresh() {
	obj.batcher.Invoke()
	logger.Printf("Refresh\n  %s\n", strings.Join(obj.getProjectNames(), ", "))
}

func (obj *ProjectsCache) Get(name string) (ProjectItem, error) {
	for i, item := range obj.Projects {
		if item.Name == name {
			return obj.Projects[i], nil
		}
	}
	var notFound ProjectItem
	return notFound, fmt.Errorf("project '%s' is not found", name)
}

func (obj *ProjectsCache) getProjectNames() []string {
	names := make([]string, len(obj.Projects))
	for i, p := range obj.Projects {
		names[i] = p.Name
	}
	return names
}

func (obj *ProjectsCache) refreshByInterval(ch <-chan time.Time) {
	for range ch {
		obj.Refresh()
	}
}

func (obj *ProjectsCache) beginIntervalRefresh(d time.Duration) {
	ticker := time.NewTicker(d)
	go obj.refreshByInterval(ticker.C)
}

func NewProjectsCache(workspace string) *ProjectsCache {
	cache := &ProjectsCache{}
	cache.workspace = workspace
	cache.batcher = batcher.NewBatcher(cache.refreshCore)
	cache.Refresh()
	cache.beginIntervalRefresh(refreshInterval)
	return cache
}

func getProjectID(configPath string) string {
	h := sha256.New()
	h.Write([]byte(configPath))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash[:8]
}

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

type projectItem struct {
	Name       string
	ConfigPath string
}

type projectsCache struct {
	locker    int32
	Workspace string
	Projects  []projectItem
}

func collectProjects(pattern string) []string {
	for i := 0; i < 3; i++ {
		matches, err := filepath.Glob(pattern)
		if err == nil {
			return matches
		}
		log.Printf("%+v", err)
		time.Sleep(3 * time.Second)
	}
	log.Fatalf("failed to refresh projects")
	return nil
}

func (obj *projectsCache) refresh() {
	if !atomic.CompareAndSwapInt32(&obj.locker, 0, 1) {
		return
	}
	defer atomic.StoreInt32(&obj.locker, 0)
	matches := collectProjects(filepath.Join(obj.Workspace, "*", "*.yaml"))
	items := make([]projectItem, len(matches))
	for i, m := range matches {
		items[i] = projectItem{
			Name:       filepath.Base(filepath.Dir(m)),
			ConfigPath: m,
		}
	}
	obj.Projects = items
}

func (obj *projectsCache) get(name string) *projectItem {
	for i, item := range obj.Projects {
		if item.Name == name {
			return &obj.Projects[i]
		}
	}
	return nil
}

func newProjectsCache(workspace string) *projectsCache {
	workDir, _ := os.Getwd()
	cache := &projectsCache{Workspace: workDir}
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

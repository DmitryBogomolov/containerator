package registry

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

type Item struct {
	Name       string
	ConfigPath string
}

type Registry struct {
	workspace string
	Projects  []Item
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

func (registry *Registry) refreshCore() {
	matches := collectProjects(filepath.Join(registry.workspace, "*", "*.yaml"))
	items := make([]Item, len(matches))
	for i, match := range matches {
		items[i] = Item{
			Name:       filepath.Base(filepath.Dir(match)),
			ConfigPath: match,
		}
	}
	registry.Projects = items
}

func (registry *Registry) Refresh() {
	registry.batcher.Invoke()
	logger.Printf("refresh\n  %s\n", strings.Join(registry.getProjectNames(), ", "))
}

func (registry *Registry) Get(name string) (*Item, error) {
	for i, item := range registry.Projects {
		if item.Name == name {
			return &registry.Projects[i], nil
		}
	}
	return nil, fmt.Errorf("project '%s' is not found", name)
}

func (registry *Registry) getProjectNames() []string {
	names := make([]string, len(registry.Projects))
	for i, p := range registry.Projects {
		names[i] = p.Name
	}
	return names
}

func (registry *Registry) refreshByInterval(ch <-chan time.Time) {
	for range ch {
		registry.Refresh()
	}
}

func (registry *Registry) beginIntervalRefresh(d time.Duration) {
	ticker := time.NewTicker(d)
	go registry.refreshByInterval(ticker.C)
}

func New(workspace string) *Registry {
	cache := &Registry{}
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

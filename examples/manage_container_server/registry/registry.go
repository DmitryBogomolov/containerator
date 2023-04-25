package registry

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/DmitryBogomolov/containerator/examples/manage_container_server/logger"
)

const (
	registryRefreshInterval   = 10 * time.Second
	collectProjectsRetryCount = 3
)

type Item struct {
	Name       string
	ConfigPath string
}

type Registry struct {
	workspace     string
	items         []Item
	refreshHandle sync.WaitGroup
	refreshLock   sync.Mutex
	refreshState  bool
}

func collectConfigFiles(pattern string) []string {
	for i := 0; i < collectProjectsRetryCount; i++ {
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

func collectItems(workspace string) []Item {
	configFiles := collectConfigFiles(filepath.Join(workspace, "*", "*.yaml"))
	items := make([]Item, len(configFiles))
	names := make([]string, len(configFiles))
	for i, configPath := range configFiles {
		name := filepath.Base(filepath.Dir(configPath))
		items[i] = Item{
			Name:       name,
			ConfigPath: configPath,
		}
		names[i] = name
	}
	logger.Printf("refresh\n  %s\n", strings.Join(names, ", "))
	return items
}

func invokeRegistryRefresh(registry *Registry) {
	registry.refreshLock.Lock()
	defer registry.refreshLock.Unlock()
	if !registry.refreshState {
		registry.refreshState = true
		registry.refreshHandle.Add(1)
		go func() {
			registry.items = collectItems(registry.workspace)
			registry.refreshHandle.Done()
			registry.refreshState = false
		}()
	}
}

func runIntervalRefresh(registry *Registry, duration time.Duration) {
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			registry.Refresh()
		}
	}()
}

func (registry *Registry) Refresh() {
	invokeRegistryRefresh(registry)
	registry.refreshHandle.Wait()
}

func (registry *Registry) Items() []Item {
	return registry.items
}

func (registry *Registry) GetItem(name string) (Item, error) {
	for _, item := range registry.items {
		if item.Name == name {
			return item, nil
		}
	}
	return Item{}, fmt.Errorf("'%s' not found", name)
}

func New(workspace string) *Registry {
	registry := Registry{
		workspace: workspace,
	}
	registry.Refresh()
	runIntervalRefresh(&registry, registryRefreshInterval)
	return &registry
}

func getProjectID(configPath string) string {
	h := sha256.New()
	h.Write([]byte(configPath))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash[:8]
}

package watcher

import (
	"slices"
	"sync"
)

type ActiveGameIds struct {
	ActiveIds []string
}

var (
	mu              sync.RWMutex
	ConnectionStore map[string]ActiveGameIds
	BackendServices []string
)

func UpdateBackendServices(backendServices []string) {
	mu.Lock()
	defer mu.Lock()
	BackendServices = backendServices
}

func GetBackEndServices() []string {
	mu.RLock()
	defer mu.RUnlock()
	return slices.Clone(BackendServices)
}

package datasource

import (
	"context"
	"sync"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

// InstanceManager handles the lifecycle of datasource instances
type InstanceManager struct {
	mu        sync.RWMutex
	instances map[string]instancemgmt.Instance
}

// NewInstanceManager creates a new instance manager
func NewInstanceManager() *InstanceManager {
	return &InstanceManager{
		instances: make(map[string]instancemgmt.Instance),
	}
}

// Get returns an existing instance or creates a new one
func (im *InstanceManager) Get(ctx context.Context, pluginContext backend.PluginContext) (instancemgmt.Instance, error) {
	im.mu.RLock()
	instance, exists := im.instances[pluginContext.DataSourceInstanceSettings.UID]
	im.mu.RUnlock()

	if exists {
		return instance, nil
	}

	im.mu.Lock()
	defer im.mu.Unlock()

	// Double check after acquiring write lock
	if instance, exists = im.instances[pluginContext.DataSourceInstanceSettings.UID]; exists {
		return instance, nil
	}

	instance, err := NewAutotaskDataSource(*pluginContext.DataSourceInstanceSettings)
	if err != nil {
		return nil, err
	}

	im.instances[pluginContext.DataSourceInstanceSettings.UID] = instance
	return instance, nil
}

// Dispose disposes of an instance
func (im *InstanceManager) Dispose(ctx context.Context, pluginContext backend.PluginContext) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if instance, exists := im.instances[pluginContext.DataSourceInstanceSettings.UID]; exists {
		if ds, ok := instance.(*AutotaskDataSource); ok {
			ds.Dispose()
		}
		delete(im.instances, pluginContext.DataSourceInstanceSettings.UID)
	}

	return nil
}

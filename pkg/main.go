package main

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/wyretech/autotask-datasource/pkg/datasource"
)

type handler struct {
	im *datasource.InstanceManager
}

func (h *handler) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	instance, err := h.im.Get(ctx, req.PluginContext)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("Failed to get instance: %v", err),
		}, nil
	}
	return instance.(*datasource.AutotaskDataSource).CheckHealth(ctx, req)
}

func (h *handler) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	instance, err := h.im.Get(ctx, req.PluginContext)
	if err != nil {
		return nil, err
	}
	return instance.(*datasource.AutotaskDataSource).QueryData(ctx, req)
}

func main() {
	// Create a new instance manager
	im := datasource.NewInstanceManager()

	// Create a new plugin
	plugin := backend.ServeOpts{
		CheckHealthHandler: &handler{im: im},
		QueryDataHandler:   &handler{im: im},
	}

	// Start the plugin
	if err := backend.Serve(plugin); err != nil {
		log.DefaultLogger.Error("Failed to start plugin", "error", err)
	}
}

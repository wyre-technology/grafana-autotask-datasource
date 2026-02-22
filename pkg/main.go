package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	ds "github.com/wyre-technology/grafana-autotask-datasource/pkg/datasource"
)

type handler struct {
	im instancemgmt.InstanceManager
}

func (h *handler) getInstance(pluginContext backend.PluginContext) (*ds.AutotaskDatasource, error) {
	instance, err := h.im.Get(context.Background(), pluginContext)
	if err != nil {
		return nil, err
	}
	return instance.(*ds.AutotaskDatasource), nil
}

func (h *handler) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	instance, err := h.getInstance(req.PluginContext)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("Failed to get instance: %v", err),
		}, nil
	}
	return instance.CheckHealth(ctx, req)
}

func (h *handler) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	instance, err := h.getInstance(req.PluginContext)
	if err != nil {
		return nil, err
	}
	return instance.QueryData(ctx, req)
}

func (h *handler) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	instance, err := h.getInstance(req.PluginContext)
	if err != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status: 500,
			Body:   []byte(fmt.Sprintf("Failed to get instance: %v", err)),
		})
	}

	switch req.Path {
	case "zoneinfo":
		return h.handleZoneInfo(ctx, req, sender, instance)
	case "query":
		return h.handleQuery(ctx, req, sender, instance)
	case "test":
		return h.handleTest(ctx, sender, instance)
	default:
		return sender.Send(&backend.CallResourceResponse{
			Status: 404,
			Body:   []byte(fmt.Sprintf("Unknown resource path: %s", req.Path)),
		})
	}
}

func (h *handler) handleZoneInfo(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender, instance *ds.AutotaskDatasource) error {
	if req.Method != "POST" {
		return sender.Send(&backend.CallResourceResponse{Status: 405, Body: []byte("Method not allowed")})
	}

	zoneInfo, err := instance.GetZoneInfo(ctx)
	if err != nil {
		log.DefaultLogger.Error("Failed to get zone info", "error", err)
		return sender.Send(&backend.CallResourceResponse{
			Status: 500,
			Body:   []byte(fmt.Sprintf("Failed to get zone info: %v", err)),
		})
	}

	body, err := json.Marshal(zoneInfo)
	if err != nil {
		return sender.Send(&backend.CallResourceResponse{Status: 500, Body: []byte("Failed to marshal response")})
	}

	return sender.Send(&backend.CallResourceResponse{
		Status:  200,
		Headers: map[string][]string{"Content-Type": {"application/json"}},
		Body:    body,
	})
}

func (h *handler) handleQuery(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender, instance *ds.AutotaskDatasource) error {
	if req.Method != "POST" {
		return sender.Send(&backend.CallResourceResponse{Status: 405, Body: []byte("Method not allowed")})
	}

	query := backend.DataQuery{
		RefID: "A",
		JSON:  req.Body,
	}
	queryReq := &backend.QueryDataRequest{
		PluginContext: req.PluginContext,
		Queries:      []backend.DataQuery{query},
	}

	response, err := instance.QueryData(ctx, queryReq)
	if err != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status: 500,
			Body:   []byte(fmt.Sprintf("Failed to execute query: %v", err)),
		})
	}

	body, err := json.Marshal(response.Responses["A"].Frames)
	if err != nil {
		return sender.Send(&backend.CallResourceResponse{Status: 500, Body: []byte("Failed to marshal response")})
	}

	return sender.Send(&backend.CallResourceResponse{
		Status:  200,
		Headers: map[string][]string{"Content-Type": {"application/json"}},
		Body:    body,
	})
}

func (h *handler) handleTest(ctx context.Context, sender backend.CallResourceResponseSender, instance *ds.AutotaskDatasource) error {
	_, err := instance.GetZoneInfo(ctx)
	if err != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status: 500,
			Body:   []byte(fmt.Sprintf("Failed to test connection: %v", err)),
		})
	}

	return sender.Send(&backend.CallResourceResponse{
		Status:  200,
		Headers: map[string][]string{"Content-Type": {"application/json"}},
		Body:    []byte(`{"status":"success","message":"Successfully connected to Autotask API"}`),
	})
}

func main() {
	im := datasource.NewInstanceManager(ds.NewAutotaskDataSource)

	h := &handler{im: im}

	opts := backend.ServeOpts{
		CheckHealthHandler:  h,
		QueryDataHandler:    h,
		CallResourceHandler: h,
	}

	if err := backend.Serve(opts); err != nil {
		log.DefaultLogger.Error("Failed to start plugin", "error", err)
	}
}

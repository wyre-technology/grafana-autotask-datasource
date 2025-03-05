package datasource

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/asachs01/autotask-go/pkg/autotask"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/wyretech/autotask-datasource/pkg/config"
)

// AutotaskDataSource represents the Autotask datasource
type AutotaskDataSource struct {
	client autotask.Client
	config *config.AutotaskConfig
}

// NewAutotaskDataSource creates a new instance of the Autotask datasource
func NewAutotaskDataSource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	cfg, err := config.LoadSettings(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to load settings: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	client := autotask.NewClient(cfg.Username, cfg.Secret, cfg.IntegrationCode)

	return &AutotaskDataSource{
		client: client,
		config: cfg,
	}, nil
}

// QueryData handles multiple queries and returns multiple responses
func (ds *AutotaskDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		res := ds.queryData(ctx, q)
		response.Responses[q.RefID] = res
	}

	return response, nil
}

// queryData handles a single query and returns a single response
func (ds *AutotaskDataSource) queryData(ctx context.Context, query backend.DataQuery) backend.DataResponse {
	var queryModel struct {
		QueryType string `json:"queryType"`
		Filter    string `json:"filter"`
	}

	if err := json.Unmarshal(query.JSON, &queryModel); err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to unmarshal query: %v", err))
	}

	switch queryModel.QueryType {
	case "tickets":
		return ds.queryTickets(ctx, query, queryModel.Filter)
	case "resources":
		return ds.queryResources(ctx, query, queryModel.Filter)
	case "companies":
		return ds.queryCompanies(ctx, query, queryModel.Filter)
	default:
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("unknown query type: %s", queryModel.QueryType))
	}
}

// CheckHealth handles health checks sent from Grafana to the plugin
func (ds *AutotaskDataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	// Test the connection by making a simple API call
	var ticketResponse struct {
		Items       []autotask.Ticket    `json:"items"`
		PageDetails autotask.PageDetails `json:"pageDetails"`
	}
	err := ds.client.Tickets().Query(ctx, "", &ticketResponse)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("Failed to connect to Autotask: %v", err),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Successfully connected to Autotask",
	}, nil
}

// Dispose cleans up datasource instance resources
func (ds *AutotaskDataSource) Dispose() {
	// Clean up any resources if needed
}

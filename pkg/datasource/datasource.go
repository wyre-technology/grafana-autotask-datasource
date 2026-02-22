package datasource

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/asachs01/autotask-go/pkg/autotask"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/wyre-technology/grafana-autotask-datasource/pkg/config"
)

// AutotaskDatasource handles communication with the Autotask API
type AutotaskDatasource struct {
	client autotask.Client
	cfg    *config.AutotaskConfig
}

// NewAutotaskDataSource creates a new datasource instance.
func NewAutotaskDataSource(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	cfg, err := config.LoadSettings(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to load settings: %w", err)
	}

	client := autotask.NewClient(cfg.Username, cfg.Secret, cfg.IntegrationCode)

	log.DefaultLogger.Debug("Created Autotask datasource", "username", cfg.Username, "url", cfg.URL)

	return &AutotaskDatasource{
		client: client,
		cfg:    cfg,
	}, nil
}

// QueryData handles multiple queries and returns multiple responses.
func (d *AutotaskDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		res := d.query(ctx, q)
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (d *AutotaskDatasource) query(ctx context.Context, query backend.DataQuery) backend.DataResponse {
	var qm QueryModel
	if err := json.Unmarshal(query.JSON, &qm); err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to unmarshal query: %v", err))
	}

	log.DefaultLogger.Debug("Query", "type", qm.QueryType, "filter", qm.Filter)

	switch qm.QueryType {
	case "tickets":
		return d.queryTickets(ctx, query, qm)
	case "resources":
		return d.queryResources(ctx, query, qm)
	case "companies":
		return d.queryCompanies(ctx, query, qm)
	case "contacts":
		return d.queryContacts(ctx, query, qm)
	default:
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("unknown query type: %s", qm.QueryType))
	}
}

// CheckHealth tests the connection to the Autotask API
func (d *AutotaskDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	zoneInfo, err := d.GetZoneInfo(ctx)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("Failed to connect to Autotask: %v", err),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: fmt.Sprintf("Connected to Autotask (Zone: %s)", zoneInfo.ZoneName),
	}, nil
}

// Dispose cleans up datasource instance resources.
func (d *AutotaskDatasource) Dispose() {}

// GetZoneInfo returns the zone information for the configured Autotask account.
// Uses direct HTTP rather than the client library to avoid Grafana proxy issues.
func (d *AutotaskDatasource) GetZoneInfo(ctx context.Context) (*autotask.ZoneInfo, error) {
	zoneURL := fmt.Sprintf("%s/atservicesrest/v1.0/ZoneInformation?user=%s", d.cfg.URL, url.QueryEscape(d.cfg.Username))

	req, err := http.NewRequestWithContext(ctx, "GET", zoneURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", d.cfg.Username, d.cfg.Secret)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	req.Header.Set("UserName", d.cfg.Username)
	req.Header.Set("Secret", d.cfg.Secret)
	req.Header.Set("ApiIntegrationCode", d.cfg.IntegrationCode)

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-OK response: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var zoneInfo autotask.ZoneInfo
	if err := json.NewDecoder(resp.Body).Decode(&zoneInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &zoneInfo, nil
}

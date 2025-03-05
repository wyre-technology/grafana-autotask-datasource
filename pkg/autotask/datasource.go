package autotask

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/asachs01/autotask-go/pkg/autotask"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/wyretech/autotask-datasource/pkg/config"
)

// AutotaskDatasource handles the communication with the Autotask API
type AutotaskDatasource struct {
	client autotask.Client
	config *config.AutotaskConfig
}

// NewAutotaskDatasource creates a new instance of the Autotask datasource
func NewAutotaskDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var cfg config.AutotaskConfig
	if err := json.Unmarshal(settings.JSONData, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	client := autotask.NewClient(
		cfg.Username,
		cfg.Secret,
		cfg.IntegrationCode,
	)

	return &AutotaskDatasource{
		client: client,
		config: &cfg,
	}, nil
}

// QueryData handles queries from Grafana
func (d *AutotaskDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		res := d.handleQuery(ctx, q)
		response.Responses[q.RefID] = res
	}

	return response, nil
}

// handleQuery processes a single query
func (d *AutotaskDatasource) handleQuery(ctx context.Context, query backend.DataQuery) backend.DataResponse {
	var queryModel struct {
		QueryType string `json:"queryType"`
		Filter    string `json:"filter"`
	}

	if err := json.Unmarshal(query.JSON, &queryModel); err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to unmarshal query: %v", err))
	}

	switch queryModel.QueryType {
	case "tickets":
		return d.queryTickets(ctx, query, queryModel.Filter)
	case "resources":
		return d.queryResources(ctx, query, queryModel.Filter)
	case "companies":
		return d.queryCompanies(ctx, query, queryModel.Filter)
	default:
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("unknown query type: %s", queryModel.QueryType))
	}
}

// queryTickets handles ticket queries
func (d *AutotaskDatasource) queryTickets(ctx context.Context, query backend.DataQuery, filter string) backend.DataResponse {
	var ticketResponse struct {
		Items []struct {
			ID           int64     `json:"id"`
			TicketNumber string    `json:"ticketNumber"`
			Title        string    `json:"title"`
			Status       string    `json:"status"`
			Priority     string    `json:"priority"`
			CreateDate   time.Time `json:"createDate"`
		} `json:"items"`
		PageDetails struct {
			Count int `json:"count"`
		} `json:"pageDetails"`
	}

	err := d.client.Tickets().Query(ctx, filter, &ticketResponse)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to query tickets: %v", err))
	}

	frame := data.NewFrame("tickets",
		data.NewField("id", nil, make([]int64, len(ticketResponse.Items))),
		data.NewField("ticketNumber", nil, make([]string, len(ticketResponse.Items))),
		data.NewField("title", nil, make([]string, len(ticketResponse.Items))),
		data.NewField("status", nil, make([]string, len(ticketResponse.Items))),
		data.NewField("priority", nil, make([]string, len(ticketResponse.Items))),
		data.NewField("createDate", nil, make([]time.Time, len(ticketResponse.Items))),
	)

	for i, ticket := range ticketResponse.Items {
		frame.SetRow(i, ticket.ID, ticket.TicketNumber, ticket.Title, ticket.Status, ticket.Priority, ticket.CreateDate)
	}

	return backend.DataResponse{
		Frames: data.Frames{frame},
	}
}

// queryResources handles resource queries
func (d *AutotaskDatasource) queryResources(ctx context.Context, query backend.DataQuery, filter string) backend.DataResponse {
	var resourceResponse struct {
		Items []struct {
			ID        int64  `json:"id"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Email     string `json:"email"`
			Active    bool   `json:"active"`
		} `json:"items"`
		PageDetails struct {
			Count int `json:"count"`
		} `json:"pageDetails"`
	}

	err := d.client.Resources().Query(ctx, filter, &resourceResponse)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to query resources: %v", err))
	}

	frame := data.NewFrame("resources",
		data.NewField("id", nil, make([]int64, len(resourceResponse.Items))),
		data.NewField("firstName", nil, make([]string, len(resourceResponse.Items))),
		data.NewField("lastName", nil, make([]string, len(resourceResponse.Items))),
		data.NewField("email", nil, make([]string, len(resourceResponse.Items))),
		data.NewField("active", nil, make([]bool, len(resourceResponse.Items))),
	)

	for i, resource := range resourceResponse.Items {
		frame.SetRow(i, resource.ID, resource.FirstName, resource.LastName, resource.Email, resource.Active)
	}

	return backend.DataResponse{
		Frames: data.Frames{frame},
	}
}

// queryCompanies handles company queries
func (d *AutotaskDatasource) queryCompanies(ctx context.Context, query backend.DataQuery, filter string) backend.DataResponse {
	var companyResponse struct {
		Items []struct {
			ID          int64  `json:"id"`
			CompanyName string `json:"companyName"`
			Name        string `json:"name"`
			Active      bool   `json:"active"`
		} `json:"items"`
		PageDetails struct {
			Count int `json:"count"`
		} `json:"pageDetails"`
	}

	err := d.client.Companies().Query(ctx, filter, &companyResponse)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to query companies: %v", err))
	}

	frame := data.NewFrame("companies",
		data.NewField("id", nil, make([]int64, len(companyResponse.Items))),
		data.NewField("companyName", nil, make([]string, len(companyResponse.Items))),
		data.NewField("name", nil, make([]string, len(companyResponse.Items))),
		data.NewField("active", nil, make([]bool, len(companyResponse.Items))),
	)

	for i, company := range companyResponse.Items {
		frame.SetRow(i, company.ID, company.CompanyName, company.Name, company.Active)
	}

	return backend.DataResponse{
		Frames: data.Frames{frame},
	}
}

// CheckHealth handles health checks sent from Grafana to the plugin
func (d *AutotaskDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	// Test the connection by making a simple API call
	var ticketResponse struct {
		Items []struct {
			ID           int64     `json:"id"`
			TicketNumber string    `json:"ticketNumber"`
			Title        string    `json:"title"`
			Status       string    `json:"status"`
			Priority     string    `json:"priority"`
			CreateDate   time.Time `json:"createDate"`
		} `json:"items"`
		PageDetails struct {
			Count int `json:"count"`
		} `json:"pageDetails"`
	}

	err := d.client.Tickets().Query(ctx, "", &ticketResponse)
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
func (d *AutotaskDatasource) Dispose() {
	// Clean up any resources if needed
}

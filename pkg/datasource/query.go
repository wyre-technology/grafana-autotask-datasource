package datasource

import (
	"context"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// QueryModel represents the query parameters
type QueryModel struct {
	QueryType string `json:"queryType"`
	Filter    string `json:"filter"`
}

// queryTickets handles ticket queries
func (ds *AutotaskDataSource) queryTickets(ctx context.Context, query backend.DataQuery, filter string) backend.DataResponse {
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

	err := ds.client.Tickets().Query(ctx, filter, &ticketResponse)
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
func (ds *AutotaskDataSource) queryResources(ctx context.Context, query backend.DataQuery, filter string) backend.DataResponse {
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

	err := ds.client.Resources().Query(ctx, filter, &resourceResponse)
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
func (ds *AutotaskDataSource) queryCompanies(ctx context.Context, query backend.DataQuery, filter string) backend.DataResponse {
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

	err := ds.client.Companies().Query(ctx, filter, &companyResponse)
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

package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// QueryModel represents the query parameters from the frontend
type QueryModel struct {
	QueryType   string `json:"queryType"`
	Filter      string `json:"filter"`
	TimeField   string `json:"timeField"`
	MaxRecords  int    `json:"maxRecords"`
}

// buildFilter combines a user filter with Grafana's time range if a timeField is set
func buildFilter(qm QueryModel, timeRange backend.TimeRange) string {
	filter := qm.Filter

	if qm.TimeField != "" {
		timeFilter := fmt.Sprintf(`{"op":"and","items":[{"op":"gte","field":"%s","value":"%s"},{"op":"lte","field":"%s","value":"%s"}]}`,
			qm.TimeField, timeRange.From.UTC().Format(time.RFC3339),
			qm.TimeField, timeRange.To.UTC().Format(time.RFC3339),
		)
		if filter != "" {
			// Wrap both in an AND
			filter = fmt.Sprintf(`{"op":"and","items":[%s,%s]}`, filter, timeFilter)
		} else {
			filter = timeFilter
		}
	}

	return filter
}

func (ds *AutotaskDatasource) queryTickets(ctx context.Context, query backend.DataQuery, qm QueryModel) backend.DataResponse {
	filter := buildFilter(qm, query.TimeRange)

	var resp struct {
		Items []struct {
			ID           int64  `json:"id"`
			TicketNumber string `json:"ticketNumber"`
			Title        string `json:"title"`
			Status       int    `json:"status"`
			Priority     int    `json:"priority"`
			CreateDate   string `json:"createDate"`
			DueDateTime  string `json:"dueDateTime"`
			CompanyID    int64  `json:"companyID"`
			QueueID      int64  `json:"queueID"`
		} `json:"items"`
		PageDetails struct {
			Count int `json:"count"`
		} `json:"pageDetails"`
	}

	if err := ds.client.Tickets().Query(ctx, filter, &resp); err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to query tickets: %v", err))
	}

	n := len(resp.Items)
	ids := make([]int64, n)
	ticketNumbers := make([]string, n)
	titles := make([]string, n)
	statuses := make([]int64, n)
	priorities := make([]int64, n)
	createDates := make([]*time.Time, n)
	dueDates := make([]*time.Time, n)
	companyIDs := make([]int64, n)
	queueIDs := make([]int64, n)

	for i, t := range resp.Items {
		ids[i] = t.ID
		ticketNumbers[i] = t.TicketNumber
		titles[i] = t.Title
		statuses[i] = int64(t.Status)
		priorities[i] = int64(t.Priority)
		createDates[i] = parseTime(t.CreateDate)
		dueDates[i] = parseTime(t.DueDateTime)
		companyIDs[i] = t.CompanyID
		queueIDs[i] = t.QueueID
	}

	frame := data.NewFrame("tickets",
		data.NewField("id", nil, ids),
		data.NewField("ticketNumber", nil, ticketNumbers),
		data.NewField("title", nil, titles),
		data.NewField("status", nil, statuses),
		data.NewField("priority", nil, priorities),
		data.NewField("createDate", nil, createDates),
		data.NewField("dueDateTime", nil, dueDates),
		data.NewField("companyID", nil, companyIDs),
		data.NewField("queueID", nil, queueIDs),
	)

	return backend.DataResponse{Frames: data.Frames{frame}}
}

func (ds *AutotaskDatasource) queryResources(ctx context.Context, query backend.DataQuery, qm QueryModel) backend.DataResponse {
	filter := buildFilter(qm, query.TimeRange)

	var resp struct {
		Items []struct {
			ID        int64  `json:"id"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Email     string `json:"email"`
			Active    bool   `json:"active"`
		} `json:"items"`
	}

	if err := ds.client.Resources().Query(ctx, filter, &resp); err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to query resources: %v", err))
	}

	n := len(resp.Items)
	ids := make([]int64, n)
	firstNames := make([]string, n)
	lastNames := make([]string, n)
	emails := make([]string, n)
	actives := make([]bool, n)

	for i, r := range resp.Items {
		ids[i] = r.ID
		firstNames[i] = r.FirstName
		lastNames[i] = r.LastName
		emails[i] = r.Email
		actives[i] = r.Active
	}

	frame := data.NewFrame("resources",
		data.NewField("id", nil, ids),
		data.NewField("firstName", nil, firstNames),
		data.NewField("lastName", nil, lastNames),
		data.NewField("email", nil, emails),
		data.NewField("active", nil, actives),
	)

	return backend.DataResponse{Frames: data.Frames{frame}}
}

func (ds *AutotaskDatasource) queryCompanies(ctx context.Context, query backend.DataQuery, qm QueryModel) backend.DataResponse {
	filter := buildFilter(qm, query.TimeRange)

	var resp struct {
		Items []struct {
			ID          int64  `json:"id"`
			CompanyName string `json:"companyName"`
			Phone       string `json:"phone"`
			Active      bool   `json:"active"`
			City        string `json:"city"`
			State       string `json:"state"`
		} `json:"items"`
	}

	if err := ds.client.Companies().Query(ctx, filter, &resp); err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to query companies: %v", err))
	}

	n := len(resp.Items)
	ids := make([]int64, n)
	names := make([]string, n)
	phones := make([]string, n)
	actives := make([]bool, n)
	cities := make([]string, n)
	states := make([]string, n)

	for i, c := range resp.Items {
		ids[i] = c.ID
		names[i] = c.CompanyName
		phones[i] = c.Phone
		actives[i] = c.Active
		cities[i] = c.City
		states[i] = c.State
	}

	frame := data.NewFrame("companies",
		data.NewField("id", nil, ids),
		data.NewField("companyName", nil, names),
		data.NewField("phone", nil, phones),
		data.NewField("active", nil, actives),
		data.NewField("city", nil, cities),
		data.NewField("state", nil, states),
	)

	return backend.DataResponse{Frames: data.Frames{frame}}
}

func (ds *AutotaskDatasource) queryContacts(ctx context.Context, query backend.DataQuery, qm QueryModel) backend.DataResponse {
	filter := buildFilter(qm, query.TimeRange)

	var resp struct {
		Items []struct {
			ID        int64  `json:"id"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Email     string `json:"emailAddress"`
			Phone     string `json:"phone"`
			CompanyID int64  `json:"companyID"`
			Active    bool   `json:"isActive"`
		} `json:"items"`
	}

	if err := ds.client.Contacts().Query(ctx, filter, &resp); err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to query contacts: %v", err))
	}

	n := len(resp.Items)
	ids := make([]int64, n)
	firstNames := make([]string, n)
	lastNames := make([]string, n)
	emails := make([]string, n)
	phones := make([]string, n)
	companyIDs := make([]int64, n)
	actives := make([]bool, n)

	for i, c := range resp.Items {
		ids[i] = c.ID
		firstNames[i] = c.FirstName
		lastNames[i] = c.LastName
		emails[i] = c.Email
		phones[i] = c.Phone
		companyIDs[i] = c.CompanyID
		actives[i] = c.Active
	}

	frame := data.NewFrame("contacts",
		data.NewField("id", nil, ids),
		data.NewField("firstName", nil, firstNames),
		data.NewField("lastName", nil, lastNames),
		data.NewField("email", nil, emails),
		data.NewField("phone", nil, phones),
		data.NewField("companyID", nil, companyIDs),
		data.NewField("active", nil, actives),
	)

	return backend.DataResponse{Frames: data.Frames{frame}}
}

// parseTime attempts to parse common Autotask date formats, returning nil on failure
func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}

	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return &t
		}
	}

	log.DefaultLogger.Warn("Failed to parse time", "value", s)
	return nil
}

// MarshalQueryModel serializes a QueryModel to JSON (used by tests)
func MarshalQueryModel(qm QueryModel) (json.RawMessage, error) {
	return json.Marshal(qm)
}

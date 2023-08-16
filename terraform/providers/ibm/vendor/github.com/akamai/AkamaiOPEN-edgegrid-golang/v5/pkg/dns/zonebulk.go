package dns

import (
	"context"
	"fmt"
	"net/http"
)

// BulkZonesCreate contains a list of one or more new Zones to create
type BulkZonesCreate struct {
	Zones []*ZoneCreate `json:"zones"`
}

// BulkZonesResponse contains response from bulk-create request
type BulkZonesResponse struct {
	RequestId      string `json:"requestId"`
	ExpirationDate string `json:"expirationDate"`
}

// BulkStatusResponse contains current status of a running or completed bulk-create request
type BulkStatusResponse struct {
	RequestId      string `json:"requestId"`
	ZonesSubmitted int    `json:"zonesSubmitted"`
	SuccessCount   int    `json:"successCount"`
	FailureCount   int    `json:"failureCount"`
	IsComplete     bool   `json:"isComplete"`
	ExpirationDate string `json:"expirationDate"`
}

// BulkFailedZone contains information about failed zone
type BulkFailedZone struct {
	Zone          string `json:"zone"`
	FailureReason string `json:"failureReason"`
}

// BulkCreateResultResponse contains the response from a completed bulk-create request
type BulkCreateResultResponse struct {
	RequestId                string            `json:"requestId"`
	SuccessfullyCreatedZones []string          `json:"successfullyCreatedZones"`
	FailedZones              []*BulkFailedZone `json:"failedZones"`
}

// BulkDeleteResultResponse contains the response from a completed bulk-delete request
type BulkDeleteResultResponse struct {
	RequestId                string            `json:"requestId"`
	SuccessfullyDeletedZones []string          `json:"successfullyDeletedZones"`
	FailedZones              []*BulkFailedZone `json:"failedZones"`
}

func (p *dns) GetBulkZoneCreateStatus(ctx context.Context, requestid string) (*BulkStatusResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetBulkZoneCreateStatus")

	bulkzonesURL := fmt.Sprintf("/config-dns/v2/zones/create-requests/%s", requestid)
	var status BulkStatusResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bulkzonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBulkZoneCreateStatus request: %w", err)
	}

	resp, err := p.Exec(req, &status)
	if err != nil {
		return nil, fmt.Errorf("GetBulkZoneCreateStatus request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &status, nil
}

func (p *dns) GetBulkZoneDeleteStatus(ctx context.Context, requestid string) (*BulkStatusResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetBulkZoneDeleteStatus")

	bulkzonesURL := fmt.Sprintf("/config-dns/v2/zones/delete-requests/%s", requestid)
	var status BulkStatusResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bulkzonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBulkZoneDeleteStatus request: %w", err)
	}

	resp, err := p.Exec(req, &status)
	if err != nil {
		return nil, fmt.Errorf("GetBulkZoneDeleteStatus request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &status, nil
}

func (p *dns) GetBulkZoneCreateResult(ctx context.Context, requestid string) (*BulkCreateResultResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetBulkZoneCreateResult")

	bulkzonesURL := fmt.Sprintf("/config-dns/v2/zones/create-requests/%s/result", requestid)
	var status BulkCreateResultResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bulkzonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBulkZoneCreateResult request: %w", err)
	}

	resp, err := p.Exec(req, &status)
	if err != nil {
		return nil, fmt.Errorf("GetBulkZoneCreateResult request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &status, nil
}

func (p *dns) GetBulkZoneDeleteResult(ctx context.Context, requestid string) (*BulkDeleteResultResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetBulkZoneDeleteResult")

	bulkzonesURL := fmt.Sprintf("/config-dns/v2/zones/delete-requests/%s/result", requestid)
	var status BulkDeleteResultResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bulkzonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetBulkZoneDeleteResult request: %w", err)
	}

	resp, err := p.Exec(req, &status)
	if err != nil {
		return nil, fmt.Errorf("GetBulkZoneDeleteResult request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &status, nil
}

func (p *dns) CreateBulkZones(ctx context.Context, bulkzones *BulkZonesCreate, zonequerystring ZoneQueryString) (*BulkZonesResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateBulkZones")

	bulkzonesURL := "/config-dns/v2/zones/create-requests?contractId=" + zonequerystring.Contract
	if len(zonequerystring.Group) > 0 {
		bulkzonesURL += "&gid=" + zonequerystring.Group
	}

	var status BulkZonesResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, bulkzonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateBulkZones request: %w", err)
	}

	resp, err := p.Exec(req, &status, bulkzones)
	if err != nil {
		return nil, fmt.Errorf("CreateBulkZones request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &status, nil
}

func (p *dns) DeleteBulkZones(ctx context.Context, zoneslist *ZoneNameListResponse, bypassSafetyChecks ...bool) (*BulkZonesResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("DeleteBulkZones")

	bulkzonesURL := "/config-dns/v2/zones/delete-requests"
	if len(bypassSafetyChecks) > 0 {
		bulkzonesURL += fmt.Sprintf("?bypassSafetyChecks=%t", bypassSafetyChecks[0])
	}

	var status BulkZonesResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, bulkzonesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DeleteBulkZones request: %w", err)
	}

	resp, err := p.Exec(req, &status, zoneslist)
	if err != nil {
		return nil, fmt.Errorf("DeleteBulkZones request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &status, nil
}

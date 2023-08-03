package dns

import (
	"context"
	"fmt"
	"net/http"
)

type (
	// Authorities contains operations available on Authorities data sources.
	Authorities interface {
		// GetAuthorities provides a list of structured read-only list of name servers.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-data-authorities
		GetAuthorities(context.Context, string) (*AuthorityResponse, error)
		// GetNameServerRecordList provides a list of name server records.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-data-authorities
		GetNameServerRecordList(context.Context, string) ([]string, error)
		// NewAuthorityResponse instantiate authority response.
		NewAuthorityResponse(context.Context, string) *AuthorityResponse
	}

	// Contract contains contractID and a list of currently assigned Akamai authoritative nameservers
	Contract struct {
		ContractID  string   `json:"contractId"`
		Authorities []string `json:"authorities"`
	}

	// AuthorityResponse contains response with a list of one or more Contracts
	AuthorityResponse struct {
		Contracts []Contract `json:"contracts"`
	}
)

func (p *dns) NewAuthorityResponse(ctx context.Context, _ string) *AuthorityResponse {

	logger := p.Log(ctx)
	logger.Debug("NewAuthorityResponse")

	authorities := &AuthorityResponse{}
	return authorities
}

func (p *dns) GetAuthorities(ctx context.Context, contractID string) (*AuthorityResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetAuthorities")

	if contractID == "" {
		return nil, fmt.Errorf("%w: GetAuthorities reqs valid contractId", ErrBadRequest)
	}

	getURL := fmt.Sprintf("/config-dns/v2/data/authorities?contractIds=%s", contractID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getauthorities request: %w", err)
	}

	var authNames AuthorityResponse
	resp, err := p.Exec(req, &authNames)
	if err != nil {
		return nil, fmt.Errorf("getauthorities request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &authNames, nil
}

func (p *dns) GetNameServerRecordList(ctx context.Context, contractID string) ([]string, error) {

	logger := p.Log(ctx)
	logger.Debug("GetNameServerRecordList")

	if contractID == "" {
		return nil, fmt.Errorf("%w: GetAuthorities reqs valid contractId", ErrBadRequest)
	}

	NSrecords, err := p.GetAuthorities(ctx, contractID)

	if err != nil {
		return nil, err
	}

	var arrLength int
	for _, c := range NSrecords.Contracts {
		arrLength = len(c.Authorities)
	}

	ns := make([]string, 0, arrLength)

	for _, r := range NSrecords.Contracts {
		for _, n := range r.Authorities {
			ns = append(ns, n)
		}
	}
	return ns, nil
}

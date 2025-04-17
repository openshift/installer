package dns

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"reflect"
	"strings"
	"sync"
)

var (
	tsigWriteLock sync.Mutex
)

type (
	// TSIGKeys contains operations available on TSIKeyG resource.
	TSIGKeys interface {
		// NewTsigKey returns bare bones tsig key struct.
		NewTsigKey(context.Context, string) *TSIGKey
		// NewTsigQueryString returns empty query string struct. No elements required.
		NewTsigQueryString(context.Context) *TSIGQueryString
		// ListTsigKeys lists the TSIG keys used by zones that you are allowed to manage.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-keys
		ListTsigKeys(context.Context, *TSIGQueryString) (*TSIGReportResponse, error)
		// GetTsigKeyZones retrieves DNS Zones using tsig key.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-keys-used-by
		GetTsigKeyZones(context.Context, *TSIGKey) (*ZoneNameListResponse, error)
		// GetTsigKeyAliases retrieves a DNS Zone's aliases.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-key-used-by
		GetTsigKeyAliases(context.Context, string) (*ZoneNameListResponse, error)
		// TsigKeyBulkUpdate updates Bulk Zones tsig key.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-keys-bulk-update
		TsigKeyBulkUpdate(context.Context, *TSIGKeyBulkPost) error
		// GetTsigKey retrieves a Tsig key for zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-key
		GetTsigKey(context.Context, string) (*TSIGKeyResponse, error)
		// DeleteTsigKey deletes tsig key for zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/delete-zones-zone-key
		DeleteTsigKey(context.Context, string) error
		// UpdateTsigKey updates tsig key for zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/put-zones-zone-key
		UpdateTsigKey(context.Context, *TSIGKey, string) error
	}

	// TSIGQueryString contains TSIG query parameters
	TSIGQueryString struct {
		ContractIds []string `json:"contractIds,omitempty"`
		Search      string   `json:"search,omitempty"`
		SortBy      []string `json:"sortBy,omitempty"`
		Gid         int64    `json:"gid,omitempty"`
	}

	// TSIGKey contains TSIG key POST response
	TSIGKey struct {
		Name      string `json:"name"`
		Algorithm string `json:"algorithm,omitempty"`
		Secret    string `json:"secret,omitempty"`
	}
	// TSIGKeyResponse contains TSIG key GET response
	TSIGKeyResponse struct {
		TSIGKey
		ZoneCount int64 `json:"zonesCount,omitempty"`
	}

	// TSIGKeyBulkPost contains TSIG key and a list of names of zones that should use the key. Used with update function.
	TSIGKeyBulkPost struct {
		Key   *TSIGKey `json:"key"`
		Zones []string `json:"zones"`
	}

	// TSIGZoneAliases contains list of zone aliases
	TSIGZoneAliases struct {
		Aliases []string `json:"aliases"`
	}

	// TSIGReportMeta contains metadata for TSIGReport response
	TSIGReportMeta struct {
		TotalElements int64    `json:"totalElements"`
		Search        string   `json:"search,omitempty"`
		Contracts     []string `json:"contracts,omitempty"`
		Gid           int64    `json:"gid,omitempty"`
		SortBy        []string `json:"sortBy,omitempty"`
	}

	// TSIGReportResponse contains response with a list of the TSIG keys used by zones.
	TSIGReportResponse struct {
		Metadata *TSIGReportMeta    `json:"metadata"`
		Keys     []*TSIGKeyResponse `json:"keys,omitempty"`
	}
)

// Validate validates RecordBody
func (key *TSIGKey) Validate() error {

	return validation.Errors{
		"Name":      validation.Validate(key.Name, validation.Required),
		"Algorithm": validation.Validate(key.Algorithm, validation.Required),
		"Secret":    validation.Validate(key.Secret, validation.Required),
	}.Filter()
}

// Validate validates TSIGKeyBulkPost
func (bulk *TSIGKeyBulkPost) Validate() error {
	return validation.Errors{
		"Key":   validation.Validate(bulk.Key, validation.Required),
		"Zones": validation.Validate(bulk.Zones, validation.Required),
	}.Filter()
}

func (p *dns) NewTsigKey(ctx context.Context, name string) *TSIGKey {

	logger := p.Log(ctx)
	logger.Debug("NewTsigKey")

	key := &TSIGKey{Name: name}
	return key
}

func (p *dns) NewTsigQueryString(ctx context.Context) *TSIGQueryString {

	logger := p.Log(ctx)
	logger.Debug("NewTsigQueryString")

	tsigquerystring := &TSIGQueryString{}
	return tsigquerystring
}

func constructTsigQueryString(tsigquerystring *TSIGQueryString) string {

	queryString := ""
	qsElems := reflect.ValueOf(tsigquerystring).Elem()
	for i := 0; i < qsElems.NumField(); i++ {
		varName := qsElems.Type().Field(i).Name
		varValue := qsElems.Field(i).Interface()
		keyVal := fmt.Sprint(varValue)
		switch varName {
		case "ContractIds":
			contractList := ""
			for j, id := range varValue.([]string) {
				contractList += id
				if j < len(varValue.([]string))-1 {
					contractList += "%2C"
				}
			}
			if len(varValue.([]string)) > 0 {
				queryString += "contractIds=" + contractList
			}
		case "SortBy":
			sortByList := ""
			for j, sb := range varValue.([]string) {
				sortByList += sb
				if j < len(varValue.([]string))-1 {
					sortByList += "%2C"
				}
			}
			if len(varValue.([]string)) > 0 {
				queryString += "sortBy=" + sortByList
			}
		case "Search":
			if keyVal != "" {
				queryString += "search=" + keyVal
			}
		case "Gid":
			if varValue.(int64) != 0 {
				queryString += "gid=" + keyVal
			}
		}
		if i < qsElems.NumField()-1 {
			queryString += "&"
		}
	}
	queryString = strings.TrimRight(queryString, "&")
	if len(queryString) > 0 {
		return "?" + queryString
	}
	return ""
}

func (p *dns) ListTsigKeys(ctx context.Context, tsigquerystring *TSIGQueryString) (*TSIGReportResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("ListTsigKeys")

	var tsigList TSIGReportResponse
	getURL := fmt.Sprintf("/config-dns/v2/keys%s", constructTsigQueryString(tsigquerystring))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListTsigKeyss request: %w", err)
	}

	resp, err := p.Exec(req, &tsigList)
	if err != nil {
		return nil, fmt.Errorf(" ListTsigKeys request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &tsigList, nil

}

func (p *dns) GetTsigKeyZones(ctx context.Context, tsigKey *TSIGKey) (*ZoneNameListResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetTsigKeyZones")

	if err := tsigKey.Validate(); err != nil {
		return nil, err
	}

	reqbody, err := convertStructToReqBody(tsigKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate request body: %w", err)
	}

	var zonesList ZoneNameListResponse
	postURL := "/config-dns/v2/keys/used-by"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqbody)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetTsigKeyZones request: %w", err)
	}

	resp, err := p.Exec(req, &zonesList)
	if err != nil {
		return nil, fmt.Errorf("GetTsigKeyZones request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &zonesList, nil
}

// TODO: Reconcile
func (p *dns) GetTsigKeyAliases(ctx context.Context, zone string) (*ZoneNameListResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetTsigKeyAliases")

	var zonesList ZoneNameListResponse
	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/key/used-by", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetTsigKeyAliases request: %w", err)
	}

	resp, err := p.Exec(req, &zonesList)
	if err != nil {
		return nil, fmt.Errorf("GetTsigKeyAliases request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &zonesList, nil
}

func (p *dns) TsigKeyBulkUpdate(ctx context.Context, tsigBulk *TSIGKeyBulkPost) error {

	logger := p.Log(ctx)
	logger.Debug("TsigKeyBulkUpdate")

	if err := tsigBulk.Validate(); err != nil {
		return err
	}

	reqbody, err := convertStructToReqBody(tsigBulk)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := "/config-dns/v2/keys/bulk-update"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqbody)
	if err != nil {
		return fmt.Errorf("failed to create TsigKeyBulkUpdate request: %w", err)
	}

	resp, err := p.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("TsigKeyBulkUpdate request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return p.Error(resp)
	}

	return nil
}

func (p *dns) GetTsigKey(ctx context.Context, zone string) (*TSIGKeyResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetTsigKey")

	var zonekey TSIGKeyResponse
	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/key", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetTsigKey request: %w", err)
	}

	resp, err := p.Exec(req, &zonekey)
	if err != nil {
		return nil, fmt.Errorf("GetTsigKey request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &zonekey, nil
}

func (p *dns) DeleteTsigKey(ctx context.Context, zone string) error {

	logger := p.Log(ctx)
	logger.Debug("DeleteTsigKey")

	delURL := fmt.Sprintf("/config-dns/v2/zones/%s/key", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create DeleteTsigKey request: %w", err)
	}

	resp, err := p.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("DeleteTsigKey request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return p.Error(resp)
	}

	return nil
}

func (p *dns) UpdateTsigKey(ctx context.Context, tsigKey *TSIGKey, zone string) error {

	logger := p.Log(ctx)
	logger.Debug("UpdateTsigKey")

	if err := tsigKey.Validate(); err != nil {
		return err
	}

	reqbody, err := convertStructToReqBody(tsigKey)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	putURL := fmt.Sprintf("/config-dns/v2/zones/%s/key", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, reqbody)
	if err != nil {
		return fmt.Errorf("failed to create UpdateTsigKey request: %w", err)
	}

	resp, err := p.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("UpdateTsigKey request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return p.Error(resp)
	}

	return nil
}

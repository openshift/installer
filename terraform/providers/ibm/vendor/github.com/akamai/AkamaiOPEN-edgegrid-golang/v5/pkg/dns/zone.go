package dns

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var (
	zoneWriteLock sync.Mutex
)

type (
	// Zones contains operations available on Zone resources.
	Zones interface {
		// ListZones retrieves a list of all zones user can access.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones
		ListZones(context.Context, ...ZoneListQueryArgs) (*ZoneListResponse, error)
		// NewZone returns a new ZoneCreate object.
		NewZone(context.Context, ZoneCreate) *ZoneCreate
		// NewZoneResponse returns a new ZoneResponse object.
		NewZoneResponse(context.Context, string) *ZoneResponse
		// NewChangeListResponse returns a new ChangeListResponse object.
		NewChangeListResponse(context.Context, string) *ChangeListResponse
		// NewZoneQueryString returns a new ZoneQueryString object.
		NewZoneQueryString(context.Context, string, string) *ZoneQueryString
		// GetZone retrieves Zone metadata.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zone
		GetZone(context.Context, string) (*ZoneResponse, error)
		//GetChangeList retrieves Zone changelist.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-changelists-zone
		GetChangeList(context.Context, string) (*ChangeListResponse, error)
		// GetMasterZoneFile retrieves master zone file.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-zone-file
		GetMasterZoneFile(context.Context, string) (string, error)
		// PostMasterZoneFile updates master zone file.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-zone-zone-file
		PostMasterZoneFile(context.Context, string, string) error
		// CreateZone creates new zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zone
		CreateZone(context.Context, *ZoneCreate, ZoneQueryString, ...bool) error
		// SaveChangelist creates a new Change List based on the most recent version of a zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-changelists
		SaveChangelist(context.Context, *ZoneCreate) error
		// SubmitChangelist submits changelist for the Zone to create default NS SOA records.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-changelists-zone-submit
		SubmitChangelist(context.Context, *ZoneCreate) error
		// UpdateZone updates zone.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/put-zone
		UpdateZone(context.Context, *ZoneCreate, ZoneQueryString) error
		// DeleteZone deletes zone.
		//
		// See: N/A
		DeleteZone(context.Context, *ZoneCreate, ZoneQueryString) error
		// ValidateZone validates zone metadata based on type.
		ValidateZone(context.Context, *ZoneCreate) error
		// GetZoneNames retrieves a list of a zone's record names.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zone-names
		GetZoneNames(context.Context, string) (*ZoneNamesResponse, error)
		// GetZoneNameTypes retrieves a zone name's record types.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zone-name-types
		GetZoneNameTypes(context.Context, string, string) (*ZoneNameTypesResponse, error)
		// CreateBulkZones submits create bulk zone request.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-create-requests
		CreateBulkZones(context.Context, *BulkZonesCreate, ZoneQueryString) (*BulkZonesResponse, error)
		// DeleteBulkZones submits delete bulk zone request.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-delete-requests
		DeleteBulkZones(context.Context, *ZoneNameListResponse, ...bool) (*BulkZonesResponse, error)
		// GetBulkZoneCreateStatus retrieves submit request status.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-create-requests-requestid
		GetBulkZoneCreateStatus(context.Context, string) (*BulkStatusResponse, error)
		//GetBulkZoneDeleteStatus retrieves submit request status.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-delete-requests-requestid
		GetBulkZoneDeleteStatus(context.Context, string) (*BulkStatusResponse, error)
		// GetBulkZoneCreateResult retrieves create request result.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-create-requests-requestid-result
		GetBulkZoneCreateResult(ctx context.Context, requestid string) (*BulkCreateResultResponse, error)
		// GetBulkZoneDeleteResult retrieves delete request result.
		//
		// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-delete-requests-requestid-result
		GetBulkZoneDeleteResult(context.Context, string) (*BulkDeleteResultResponse, error)
	}

	// ZoneQueryString contains zone query parameters
	ZoneQueryString struct {
		Contract string
		Group    string
	}

	// ZoneCreate contains zone create request
	ZoneCreate struct {
		Zone                  string   `json:"zone"`
		Type                  string   `json:"type"`
		Masters               []string `json:"masters,omitempty"`
		Comment               string   `json:"comment,omitempty"`
		SignAndServe          bool     `json:"signAndServe"`
		SignAndServeAlgorithm string   `json:"signAndServeAlgorithm,omitempty"`
		TsigKey               *TSIGKey `json:"tsigKey,omitempty"`
		Target                string   `json:"target,omitempty"`
		EndCustomerID         string   `json:"endCustomerId,omitempty"`
		ContractID            string   `json:"contractId,omitempty"`
	}

	// ZoneResponse contains zone create response
	ZoneResponse struct {
		Zone                  string   `json:"zone,omitempty"`
		Type                  string   `json:"type,omitempty"`
		Masters               []string `json:"masters,omitempty"`
		Comment               string   `json:"comment,omitempty"`
		SignAndServe          bool     `json:"signAndServe"`
		SignAndServeAlgorithm string   `json:"signAndServeAlgorithm,omitempty"`
		TsigKey               *TSIGKey `json:"tsigKey,omitempty"`
		Target                string   `json:"target,omitempty"`
		EndCustomerID         string   `json:"endCustomerId,omitempty"`
		ContractID            string   `json:"contractId,omitempty"`
		AliasCount            int64    `json:"aliasCount,omitempty"`
		ActivationState       string   `json:"activationState,omitempty"`
		LastActivationDate    string   `json:"lastActivationDate,omitempty"`
		LastModifiedBy        string   `json:"lastModifiedBy,omitempty"`
		LastModifiedDate      string   `json:"lastModifiedDate,omitempty"`
		VersionId             string   `json:"versionId,omitempty"`
	}

	// ZoneListQueryArgs contains parameters for List Zones query
	ZoneListQueryArgs struct {
		ContractIDs string
		Page        int
		PageSize    int
		Search      string
		ShowAll     bool
		SortBy      string
		Types       string
	}

	// ListMetadata contains metadata for List Zones request
	ListMetadata struct {
		ContractIDs   []string `json:"contractIds"`
		Page          int      `json:"page"`
		PageSize      int      `json:"pageSize"`
		ShowAll       bool     `json:"showAll"`
		TotalElements int      `json:"totalElements"`
	} //`json:"metadata"`

	// ZoneListResponse contains response for List Zones request
	ZoneListResponse struct {
		Metadata *ListMetadata   `json:"metadata,omitempty"`
		Zones    []*ZoneResponse `json:"zones,omitempty"`
	}

	// ChangeListResponse contains metadata about a change list
	ChangeListResponse struct {
		Zone             string `json:"zone,omitempty"`
		ChangeTag        string `json:"changeTag,omitempty"`
		ZoneVersionID    string `json:"zoneVersionId,omitempty"`
		LastModifiedDate string `json:"lastModifiedDate,omitempty"`
		Stale            bool   `json:"stale,omitempty"`
	}

	// ZoneNameListResponse contains response with a list of zone's names and aliases
	ZoneNameListResponse struct {
		Zones   []string `json:"zones"`
		Aliases []string `json:"aliases,omitempty"`
	}

	// ZoneNamesResponse contains record set names for zone
	ZoneNamesResponse struct {
		Names []string `json:"names"`
	}

	// ZoneNameTypesResponse contains record set types for zone
	ZoneNameTypesResponse struct {
		Types []string `json:"types"`
	}
)

var zoneStructMap = map[string]string{
	"Zone":                  "zone",
	"Type":                  "type",
	"Masters":               "masters",
	"Comment":               "comment",
	"SignAndServe":          "signAndServe",
	"SignAndServeAlgorithm": "signAndServeAlgorithm",
	"TsigKey":               "tsigKey",
	"Target":                "target",
	"EndCustomerID":         "endCustomerId",
	"ContractId":            "contractId"}

// Util to convert struct to http request body, eg. io.reader
func convertStructToReqBody(srcstruct interface{}) (io.Reader, error) {

	reqbody, err := json.Marshal(srcstruct)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(reqbody), nil
}

func (p *dns) ListZones(ctx context.Context, queryArgs ...ZoneListQueryArgs) (*ZoneListResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("ListZones")

	// construct GET url
	getURL := fmt.Sprintf("/config-dns/v2/zones")
	if len(queryArgs) > 1 {
		return nil, fmt.Errorf("ListZones QueryArgs invalid")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create listzones request: %w", err)
	}

	q := req.URL.Query()
	if len(queryArgs) > 0 {
		if queryArgs[0].Page > 0 {
			q.Add("page", strconv.Itoa(queryArgs[0].Page))
		}
		if queryArgs[0].PageSize > 0 {
			q.Add("pageSize", strconv.Itoa(queryArgs[0].PageSize))
		}
		if queryArgs[0].Search != "" {
			q.Add("search", queryArgs[0].Search)
		}
		q.Add("showAll", strconv.FormatBool(queryArgs[0].ShowAll))
		if queryArgs[0].SortBy != "" {
			q.Add("sortBy", queryArgs[0].SortBy)
		}
		if queryArgs[0].Types != "" {
			q.Add("types", queryArgs[0].Types)
		}
		if queryArgs[0].ContractIDs != "" {
			q.Add("contractIds", queryArgs[0].ContractIDs)
		}
		req.URL.RawQuery = q.Encode()
	}

	var zonelist ZoneListResponse
	resp, err := p.Exec(req, &zonelist)
	if err != nil {
		return nil, fmt.Errorf("listzones request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &zonelist, nil
}

func (p *dns) NewZone(ctx context.Context, params ZoneCreate) *ZoneCreate {

	logger := p.Log(ctx)
	logger.Debug("NewZone")

	zone := &ZoneCreate{Zone: params.Zone,
		Type:                  params.Type,
		Masters:               params.Masters,
		TsigKey:               params.TsigKey,
		Target:                params.Target,
		EndCustomerID:         params.EndCustomerID,
		ContractID:            params.ContractID,
		Comment:               params.Comment,
		SignAndServe:          params.SignAndServe,
		SignAndServeAlgorithm: params.SignAndServeAlgorithm}

	logger.Debugf("Created zone: %v", zone)
	return zone
}

func (p *dns) NewZoneResponse(ctx context.Context, zonename string) *ZoneResponse {

	logger := p.Log(ctx)
	logger.Debug("NewZoneResponse")

	zone := &ZoneResponse{Zone: zonename}
	return zone
}

func (p *dns) NewChangeListResponse(ctx context.Context, zone string) *ChangeListResponse {

	logger := p.Log(ctx)
	logger.Debug("NewChangeListResponse")

	changelist := &ChangeListResponse{Zone: zone}
	return changelist
}

func (p *dns) NewZoneQueryString(ctx context.Context, contract string, group string) *ZoneQueryString {

	logger := p.Log(ctx)
	logger.Debug("NewZoneQueryString")

	zonequerystring := &ZoneQueryString{Contract: contract, Group: group}
	return zonequerystring
}

func (p *dns) GetZone(ctx context.Context, zonename string) (*ZoneResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetZone")

	var zone ZoneResponse

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s", zonename)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetZone request: %w", err)
	}

	resp, err := p.Exec(req, &zone)
	if err != nil {
		return nil, fmt.Errorf("GetZone request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &zone, nil
}

func (p *dns) GetChangeList(ctx context.Context, zone string) (*ChangeListResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetChangeList")

	var changelist ChangeListResponse
	getURL := fmt.Sprintf("/config-dns/v2/changelists/%s", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetChangeList request: %w", err)
	}

	resp, err := p.Exec(req, &changelist)
	if err != nil {
		return nil, fmt.Errorf("GetChangeList request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &changelist, nil
}

func (p *dns) GetMasterZoneFile(ctx context.Context, zone string) (string, error) {

	logger := p.Log(ctx)
	logger.Debug("GetMasterZoneFile")

	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/zone-file", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create GetMasterZoneFile request: %w", err)
	}
	req.Header.Add("Accept", "text/dns")

	resp, err := p.Exec(req, nil)
	if err != nil {
		return "", fmt.Errorf("GetMasterZoneFile request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", p.Error(resp)
	}

	masterfile, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("GetMasterZoneFile request failed: %w", err)
	}

	return string(masterfile), nil
}

func (p *dns) PostMasterZoneFile(ctx context.Context, zone string, filedata string) error {

	logger := p.Log(ctx)
	logger.Debug("PostMasterZoneFile")

	mtresp := ""
	pmzfURL := fmt.Sprintf("/config-dns/v2/zones/%s/zone-file", zone)
	buf := bytes.NewReader([]byte(filedata))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, pmzfURL, buf)
	if err != nil {
		return fmt.Errorf("failed to create PostMasterZoneFile request: %w", err)
	}

	req.Header.Set("Content-Type", "text/dns")

	resp, err := p.Exec(req, &mtresp)
	if err != nil {
		return fmt.Errorf("Create PostMasterZoneFile failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return p.Error(resp)
	}

	return nil
}

func (p *dns) CreateZone(ctx context.Context, zone *ZoneCreate, zonequerystring ZoneQueryString, clearConn ...bool) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	zoneWriteLock.Lock()
	defer zoneWriteLock.Unlock()

	logger := p.Log(ctx)
	logger.Debug("Zone Create")

	if err := p.ValidateZone(ctx, zone); err != nil {
		return err
	}

	zoneMap := filterZoneCreate(zone)

	var zoneresponse ZoneResponse
	zoneURL := "/config-dns/v2/zones/?contractId=" + zonequerystring.Contract
	if len(zonequerystring.Group) > 0 {
		zoneURL += "&gid=" + zonequerystring.Group
	}

	reqbody, err := convertStructToReqBody(zoneMap)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, zoneURL, reqbody)
	if err != nil {
		return fmt.Errorf("failed to create Zone Create request: %w", err)
	}

	resp, err := p.Exec(req, &zoneresponse)
	if err != nil {
		return fmt.Errorf("Create Zone request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return p.Error(resp)
	}

	if strings.ToUpper(zone.Type) == "PRIMARY" {
		// Timing issue with Create immediately followed by SaveChangelist
		for _, clear := range clearConn {
			// should only be one entry
			if clear {
				logger.Info("Clearing Idle Connections")
				p.Client().CloseIdleConnections()
			}
		}
	}

	return nil
}

func (p *dns) SaveChangelist(ctx context.Context, zone *ZoneCreate) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	zoneWriteLock.Lock()
	defer zoneWriteLock.Unlock()

	logger := p.Log(ctx)
	logger.Debug("SaveChangeList")

	reqbody, err := convertStructToReqBody("")
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := fmt.Sprintf("/config-dns/v2/changelists/?zone=%s", zone.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqbody)
	if err != nil {
		return fmt.Errorf("failed to create SaveChangeList request: %w", err)
	}

	resp, err := p.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("SaveChangeList request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return p.Error(resp)
	}

	return nil
}

func (p *dns) SubmitChangelist(ctx context.Context, zone *ZoneCreate) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	zoneWriteLock.Lock()
	defer zoneWriteLock.Unlock()

	logger := p.Log(ctx)
	logger.Debug("SubmitChangeList")

	reqbody, err := convertStructToReqBody("")
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := fmt.Sprintf("/config-dns/v2/changelists/%s/submit", zone.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqbody)
	if err != nil {
		return fmt.Errorf("failed to create SubmitChangeList request: %w", err)
	}

	resp, err := p.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("SubmitChangeList request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return p.Error(resp)
	}

	return nil
}

func (p *dns) UpdateZone(ctx context.Context, zone *ZoneCreate, _ ZoneQueryString) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	zoneWriteLock.Lock()
	defer zoneWriteLock.Unlock()

	logger := p.Log(ctx)
	logger.Debug("Zone Update")

	if err := p.ValidateZone(ctx, zone); err != nil {
		return err
	}

	zoneMap := filterZoneCreate(zone)
	reqbody, err := convertStructToReqBody(zoneMap)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	putURL := fmt.Sprintf("/config-dns/v2/zones/%s", zone.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, reqbody)
	if err != nil {
		return fmt.Errorf("failed to create Get Update request: %w", err)
	}

	var zoneresp ZoneResponse
	resp, err := p.Exec(req, &zoneresp)
	if err != nil {
		return fmt.Errorf("Zone Update request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return p.Error(resp)
	}

	return nil

}

func (p *dns) DeleteZone(ctx context.Context, zone *ZoneCreate, _ ZoneQueryString) error {
	// remove all the records except for SOA
	// which is required and save the zone

	zoneWriteLock.Lock()
	defer zoneWriteLock.Unlock()

	logger := p.Log(ctx)
	logger.Debug("Zone Delete")

	if zone.Zone == "" {
		return fmt.Errorf("Zone name missing")
	}

	deleteURL := fmt.Sprintf("/config-dns/v2/zones/%s", zone.Zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create Zone Delete request: %w", err)
	}

	resp, err := p.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("Zone Delete request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	if resp.StatusCode != http.StatusNoContent {
		return p.Error(resp)
	}

	return nil

}

func filterZoneCreate(zone *ZoneCreate) map[string]interface{} {

	zoneType := strings.ToUpper(zone.Type)
	filteredZone := make(map[string]interface{})
	zoneElems := reflect.ValueOf(zone).Elem()
	for i := 0; i < zoneElems.NumField(); i++ {
		varName := zoneElems.Type().Field(i).Name
		varLower := zoneStructMap[varName]
		varValue := zoneElems.Field(i).Interface()
		switch varName {
		case "Target":
			if zoneType == "ALIAS" {
				filteredZone[varLower] = varValue
			}
		case "TsigKey":
			if zoneType == "SECONDARY" {
				filteredZone[varLower] = varValue
			}
		case "Masters":
			if zoneType == "SECONDARY" {
				filteredZone[varLower] = varValue
			}
		case "SignAndServe":
			if zoneType != "ALIAS" {
				filteredZone[varLower] = varValue
			}
		case "SignAndServeAlgorithm":
			if zoneType != "ALIAS" {
				filteredZone[varLower] = varValue
			}
		default:
			filteredZone[varLower] = varValue
		}
	}

	return filteredZone
}

// ValidateZone validates ZoneCreate Object
func (p *dns) ValidateZone(ctx context.Context, zone *ZoneCreate) error {

	logger := p.Log(ctx)
	logger.Debug("ValidateZone")

	if len(zone.Zone) == 0 {
		return fmt.Errorf("Zone name is required")
	}
	ztype := strings.ToUpper(zone.Type)
	if ztype != "PRIMARY" && ztype != "SECONDARY" && ztype != "ALIAS" {
		return fmt.Errorf("Invalid zone type")
	}
	if ztype != "SECONDARY" && zone.TsigKey != nil {
		return fmt.Errorf("TsigKey is invalid for %s zone type", ztype)
	}
	if ztype == "ALIAS" {
		if len(zone.Target) == 0 {
			return fmt.Errorf("Target is required for Alias zone type")
		}
		if zone.Masters != nil && len(zone.Masters) > 0 {
			return fmt.Errorf("Masters is invalid for Alias zone type")
		}
		if zone.SignAndServe {
			return fmt.Errorf("SignAndServe is invalid for Alias zone type")
		}
		if len(zone.SignAndServeAlgorithm) > 0 {
			return fmt.Errorf("SignAndServeAlgorithm is invalid for Alias zone type")
		}
		return nil
	}
	// Primary or Secondary
	if len(zone.Target) > 0 {
		return fmt.Errorf("Target is invalid for %s zone type", ztype)
	}
	if zone.Masters != nil && len(zone.Masters) > 0 && ztype == "PRIMARY" {
		return fmt.Errorf("Masters is invalid for Primary zone type")
	}

	return nil
}

func (p *dns) GetZoneNames(ctx context.Context, zone string) (*ZoneNamesResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetZoneNames")

	var znresponse ZoneNamesResponse
	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/names", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetZoneNames request: %w", err)
	}

	resp, err := p.Exec(req, &znresponse)
	if err != nil {
		return nil, fmt.Errorf("GetZoneNames request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &znresponse, nil
}

func (p *dns) GetZoneNameTypes(ctx context.Context, zname string, zone string) (*ZoneNameTypesResponse, error) {

	logger := p.Log(ctx)
	logger.Debug(" GetZoneNameTypes")

	var zntypes ZoneNameTypesResponse
	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types", zone, zname)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetZoneNameTypes request: %w", err)
	}

	resp, err := p.Exec(req, &zntypes)
	if err != nil {
		return nil, fmt.Errorf("GetZoneNameTypes request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &zntypes, nil
}

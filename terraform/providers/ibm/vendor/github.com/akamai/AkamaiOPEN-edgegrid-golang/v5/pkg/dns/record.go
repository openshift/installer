package dns

import (
	"context"
	"fmt"
	"net/http"

	"net"
	"sync"
)

// Records contains operations available on a Record resource.
type Records interface {
	// RecordToMap returns a map containing record content.
	RecordToMap(context.Context, *RecordBody) map[string]interface{}
	// NewRecordBody returns bare bones tsig key struct.
	NewRecordBody(context.Context, RecordBody) *RecordBody
	// GetRecordList retrieves recordset list based on type.
	//
	// See: https://techdocs.akamai.com/edge-dns/reference/get-zones-zone-recordsets
	GetRecordList(context.Context, string, string, string) (*RecordSetResponse, error)
	// GetRdata retrieves record rdata, e.g. target.
	GetRdata(context.Context, string, string, string) ([]string, error)
	// ProcessRdata process rdata.
	ProcessRdata(context.Context, []string, string) []string
	// ParseRData parses rdata. returning map.
	ParseRData(context.Context, string, []string) map[string]interface{}
	// GetRecord retrieves a recordset and returns as RecordBody.
	//
	// See:  https://techdocs.akamai.com/edge-dns/reference/get-zone-name-type
	GetRecord(context.Context, string, string, string) (*RecordBody, error)
	// CreateRecord creates recordset.
	//
	// See: https://techdocs.akamai.com/edge-dns/reference/post-zones-zone-names-name-types-type
	CreateRecord(context.Context, *RecordBody, string, ...bool) error
	// DeleteRecord removes recordset.
	//
	// See: https://techdocs.akamai.com/edge-dns/reference/delete-zone-name-type
	DeleteRecord(context.Context, *RecordBody, string, ...bool) error
	// UpdateRecord replaces the recordset.
	//
	// See: https://techdocs.akamai.com/edge-dns/reference/put-zones-zone-names-name-types-type
	UpdateRecord(context.Context, *RecordBody, string, ...bool) error
	// FullIPv6 is utility method to convert IP to string.
	FullIPv6(context.Context, net.IP) string
	// PadCoordinates is utility method to convert IP to normalize coordinates.
	PadCoordinates(context.Context, string) string
}

// RecordBody contains request body for dns record
type RecordBody struct {
	Name       string `json:"name,omitempty"`
	RecordType string `json:"type,omitempty"`
	TTL        int    `json:"ttl,omitempty"`
	// Active field no longer used in v2
	Active bool     `json:"active,omitempty"`
	Target []string `json:"rdata,omitempty"`
}

var (
	zoneRecordWriteLock sync.Mutex
)

// Validate validates RecordBody
func (rec *RecordBody) Validate() error {

	if len(rec.Name) < 1 {
		return fmt.Errorf("Record body is missing Name")
	}
	if len(rec.RecordType) < 1 {
		return fmt.Errorf("Record body is missing RecordType")
	}
	if rec.TTL == 0 {
		return fmt.Errorf("Record body is missing TTL")
	}
	if rec.Target == nil || len(rec.Target) < 1 {
		return fmt.Errorf("Record body is missing Target")
	}

	return nil
}

func (p *dns) RecordToMap(ctx context.Context, record *RecordBody) map[string]interface{} {

	logger := p.Log(ctx)
	logger.Debug("RecordToMap")

	if err := record.Validate(); err != nil {
		logger.Errorf("Record to map failed. %w", err)
		return nil
	}

	return map[string]interface{}{
		"name":       record.Name,
		"ttl":        record.TTL,
		"recordtype": record.RecordType,
		// active no longer used
		"active": record.Active,
		"target": record.Target,
	}
}

func (p *dns) NewRecordBody(ctx context.Context, params RecordBody) *RecordBody {

	logger := p.Log(ctx)
	logger.Debug("NewRecordBody")

	recordbody := &RecordBody{Name: params.Name}
	return recordbody
}

// Eval option lock arg passed into writable endpoints. Default is true, e.g. lock
func localLock(lockArg []bool) bool {

	for _, lock := range lockArg {
		// should only be one entry
		return lock
	}

	return true

}

func (p *dns) CreateRecord(ctx context.Context, record *RecordBody, zone string, recLock ...bool) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(recLock) {
		zoneRecordWriteLock.Lock()
		defer zoneRecordWriteLock.Unlock()
	}

	logger := p.Log(ctx)
	logger.Debug("CreateRecord")
	logger.Debugf("DNS Lib Create Record: [%v]", record)
	if err := record.Validate(); err != nil {
		logger.Errorf("Record content not valid: %w", err)
		return fmt.Errorf("Record content not valid. [%w]", err)
	}

	reqbody, err := convertStructToReqBody(record)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	var rec RecordBody
	postURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/%s", zone, record.Name, record.RecordType)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqbody)
	if err != nil {
		return fmt.Errorf("failed to create CreateRecord request: %w", err)
	}

	resp, err := p.Exec(req, &rec)
	if err != nil {
		return fmt.Errorf("CreateRecord request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return p.Error(resp)
	}

	return nil
}

func (p *dns) UpdateRecord(ctx context.Context, record *RecordBody, zone string, recLock ...bool) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(recLock) {
		zoneRecordWriteLock.Lock()
		defer zoneRecordWriteLock.Unlock()
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateRecord")
	logger.Debugf("DNS Lib Update Record: [%v]", record)
	if err := record.Validate(); err != nil {
		logger.Errorf("Record content not valid: %s", err.Error())
		return fmt.Errorf("Record content not valid. [%w]", err)
	}

	reqbody, err := convertStructToReqBody(record)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	var rec RecordBody
	putURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/%s", zone, record.Name, record.RecordType)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, reqbody)
	if err != nil {
		return fmt.Errorf("failed to create UpdateRecord request: %w", err)
	}

	resp, err := p.Exec(req, &rec)
	if err != nil {
		return fmt.Errorf("UpdateRecord request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return p.Error(resp)
	}

	return nil
}

func (p *dns) DeleteRecord(ctx context.Context, record *RecordBody, zone string, recLock ...bool) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(recLock) {
		zoneRecordWriteLock.Lock()
		defer zoneRecordWriteLock.Unlock()
	}

	logger := p.Log(ctx)
	logger.Debug("DeleteRecord")

	if err := record.Validate(); err != nil {
		logger.Errorf("Record content not valid: %w", err)
		return fmt.Errorf("Record content not valid. [%w]", err)
	}

	//var mtbody string
	deleteURL := fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/%s", zone, record.Name, record.RecordType)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create DeleteRecord request: %w", err)
	}

	resp, err := p.Exec(req, nil) //, &mtbody)
	if err != nil {
		return fmt.Errorf("DeleteRecord request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return p.Error(resp)
	}

	return nil
}

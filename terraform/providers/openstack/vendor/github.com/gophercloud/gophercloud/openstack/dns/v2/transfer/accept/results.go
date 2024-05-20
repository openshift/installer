package accept

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract interprets a GetResult, CreateResult as a TransferAccept.
// An error is returned if the original call or the extraction failed.
func (r commonResult) Extract() (*TransferAccept, error) {
	var s *TransferAccept
	err := r.ExtractInto(&s)
	return s, err
}

// CreateResult is the result of a Create request. Call its Extract method
// to interpret the result as a TransferAccept.
type CreateResult struct {
	commonResult
}

// GetResult is the result of a Get request. Call its Extract method
// to interpret the result as a TransferAccept.
type GetResult struct {
	commonResult
}

// TransferAcceptPage is a single page of TransferAccept results.
type TransferAcceptPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if the page contains no results.
func (r TransferAcceptPage) IsEmpty() (bool, error) {
	s, err := ExtractTransferAccepts(r)
	return len(s) == 0, err
}

// ExtractTransferAccepts extracts a slice of TransferAccept from a List result.
func ExtractTransferAccepts(r pagination.Page) ([]TransferAccept, error) {
	var s struct {
		TransferAccepts []TransferAccept `json:"transfer_accepts"`
	}
	err := (r.(TransferAcceptPage)).ExtractInto(&s)
	return s.TransferAccepts, err
}

// TransferAccept represents a Zone transfer accept task.
type TransferAccept struct {
	// ID for this zone transfer accept.
	ID string `json:"id"`

	// Status is current status of the zone transfer request.
	Status string `json:"status"`

	// ProjectID identifies the project/tenant owning this resource.
	ProjectID string `json:"project_id"`

	// ZoneID is the ID for the zone that was being exported.
	ZoneID string `json:"zone_id"`

	// Key is used as part of the zone transfer accept process.
	// This is only shown to the creator, and must be communicated out of band.
	Key string `json:"key"`

	// ZoneTransferRequestID is ID for this zone transfer request
	ZoneTransferRequestID string `json:"zone_transfer_request_id"`

	// CreatedAt is the date when the zone transfer accept was created.
	CreatedAt time.Time `json:"-"`

	// UpdatedAt is the date when the last change was made to the zone transfer accept.
	UpdatedAt time.Time `json:"-"`

	// Links includes HTTP references to the itself, useful for passing along
	// to other APIs that might want a server reference.
	Links map[string]interface{} `json:"links"`
}

func (r *TransferAccept) UnmarshalJSON(b []byte) error {
	type tmp TransferAccept
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = TransferAccept(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return err
}

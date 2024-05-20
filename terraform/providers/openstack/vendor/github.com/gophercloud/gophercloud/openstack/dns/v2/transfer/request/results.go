package request

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract interprets a GetResult, CreateResult or UpdateResult as a TransferRequest.
// An error is returned if the original call or the extraction failed.
func (r commonResult) Extract() (*TransferRequest, error) {
	var s *TransferRequest
	err := r.ExtractInto(&s)
	return s, err
}

// CreateResult is the result of a Create request. Call its Extract method
// to interpret the result as a TransferRequest.
type CreateResult struct {
	commonResult
}

// GetResult is the result of a Get request. Call its Extract method
// to interpret the result as a TransferRequest.
type GetResult struct {
	commonResult
}

// UpdateResult is the result of an Update request. Call its Extract method
// to interpret the result as a TransferRequest.
type UpdateResult struct {
	commonResult
}

// DeleteResult is the result of a Delete request. Call its ExtractErr method
// to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// TransferRequestPage is a single page of TransferRequest results.
type TransferRequestPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if the page contains no results.
func (r TransferRequestPage) IsEmpty() (bool, error) {
	s, err := ExtractTransferRequests(r)
	return len(s) == 0, err
}

// ExtractTransferRequests extracts a slice of TransferRequest from a List result.
func ExtractTransferRequests(r pagination.Page) ([]TransferRequest, error) {
	var s struct {
		TransferRequests []TransferRequest `json:"transfer_requests"`
	}
	err := (r.(TransferRequestPage)).ExtractInto(&s)
	return s.TransferRequests, err
}

// TransferRequest represents a Zone transfer request task.
type TransferRequest struct {
	// ID uniquely identifies this transfer request zone amongst all other transfer requests,
	// including those not accessible to the current tenant.
	ID string `json:"id"`

	// ZoneID is the ID for the zone that is being exported.
	ZoneID string `json:"zone_id"`

	// Name is the name of the zone that is being exported.
	ZoneName string `json:"zone_name"`

	// ProjectID identifies the project/tenant owning this resource.
	ProjectID string `json:"project_id"`

	// TargetProjectID identifies the project/tenant to transfer this resource.
	TargetProjectID string `json:"target_project_id"`

	// Key is used as part of the zone transfer accept process.
	// This is only shown to the creator, and must be communicated out of band.
	Key string `json:"key"`

	// Description for the resource.
	Description string `json:"description"`

	// Status is the status of the resource.
	Status string `json:"status"`

	// CreatedAt is the date when the zone transfer request was created.
	CreatedAt time.Time `json:"-"`

	// UpdatedAt is the date when the last change was made to the zone transfer request.
	UpdatedAt time.Time `json:"-"`

	// Links includes HTTP references to the itself, useful for passing along
	// to other APIs that might want a server reference.
	Links map[string]interface{} `json:"links"`
}

func (r *TransferRequest) UnmarshalJSON(b []byte) error {
	type tmp TransferRequest
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = TransferRequest(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return err
}

package fc

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

const (
	reservedCapacitiesPath = "/reservedCapacities"
)

type reservedCapacityMetadata struct {
	InstanceID       *string `json:"instanceId"`
	CU               *int64  `json:"cu"`
	Deadline         *string `json:"deadline"`
	CreatedTime      *string `json:"createdTime"`
	LastModifiedTime *string `json:"lastModifiedTime"`
	IsRefunded       *string `json:"isRefunded"`
}

// ListReservedCapacitiesOutput : ...
type ListReservedCapacitiesOutput struct {
	Header            http.Header
	ReservedCapacities []*reservedCapacityMetadata `json:"reservedCapacities"`
	NextToken         *string                     `json:"nextToken,omitempty"`
}

// ListReservedCapacitiesInput : ...
type ListReservedCapacitiesInput struct {
	Query
}

// NewListReservedCapacitiesInput : ...
func NewListReservedCapacitiesInput() *ListReservedCapacitiesInput {
	return &ListReservedCapacitiesInput{}
}

func (i *ListReservedCapacitiesInput) WithNextToken(nextToken string) *ListReservedCapacitiesInput {
	i.NextToken = &nextToken
	return i
}

func (i *ListReservedCapacitiesInput) WithLimit(limit int32) *ListReservedCapacitiesInput {
	i.Limit = &limit
	return i
}

func (i *ListReservedCapacitiesInput) GetQueryParams() url.Values {
	out := url.Values{}
	if i.NextToken != nil {
		out.Set("nextToken", *i.NextToken)
	}

	if i.Limit != nil {
		out.Set("limit", strconv.FormatInt(int64(*i.Limit), 10))
	}

	return out
}

func (i *ListReservedCapacitiesInput) GetPath() string {
	return reservedCapacitiesPath
}

func (i *ListReservedCapacitiesInput) GetHeaders() Header {
	return make(Header, 0)
}

func (i *ListReservedCapacitiesInput) GetPayload() interface{} {
	return nil
}

func (i *ListReservedCapacitiesInput) Validate() error {
	return nil
}

func (o ListReservedCapacitiesOutput) String() string {
	b, err := json.MarshalIndent(o, "", printIndent)
	if err != nil {
		return ""
	}
	return string(b)
}

func (o ListReservedCapacitiesOutput) GetRequestID() string {
	return GetRequestID(o.Header)
}

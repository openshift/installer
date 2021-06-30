package alibabacloud

import (
	"sync"
)

// Metadata holds additional metadata for InstallConfig resources that
// does not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	client     *Client
	BaseDomain string
	Region     string `json:"region,omitempty"`

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(region string, baseDomain string) *Metadata {
	return &Metadata{Region: region, BaseDomain: baseDomain}
}

// Client returns a client used for making API calls to Alibaba Cloud services.
func (m *Metadata) Client(regionId string) (*Client, error) {
	if m.client == nil {
		client, err := NewClient(regionId)
		if err != nil {
			return nil, err
		}
		m.client = client
	}
	return m.client, nil
}

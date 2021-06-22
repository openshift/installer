package ibmcloud

import (
	"context"
	"fmt"
	"sync"
)

// Metadata holds additional metadata for InstallConfig resources that
// does not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	BaseDomain string

	accountID      string
	cisInstanceCRN string
	client         *Client

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(baseDomain string) *Metadata {
	return &Metadata{BaseDomain: baseDomain}
}

// AccountID returns the IBM Cloud account ID associated with the authentication
// credentials.
func (m *Metadata) AccountID(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.accountID == "" {
		client, err := m.Client()
		if err != nil {
			return "", err
		}

		apiKeyDetails, err := client.GetAuthenticatorAPIKeyDetails(ctx)
		if err != nil {
			return "", err
		}

		m.accountID = *apiKeyDetails.AccountID
	}
	return m.accountID, nil
}

// CISInstanceCRN returns the Cloud Internet Services instance CRN that is
// managing the DNS zone for the base domain.
func (m *Metadata) CISInstanceCRN(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.cisInstanceCRN == "" {
		client, err := m.Client()
		if err != nil {
			return "", err
		}

		zones, err := client.GetDNSZones(ctx)
		if err != nil {
			return "", err
		}

		for _, z := range zones {
			if z.Name == m.BaseDomain {
				m.cisInstanceCRN = z.CISInstanceCRN
				return m.cisInstanceCRN, nil
			}
		}
		return "", fmt.Errorf("cisInstanceCRN unknown due to DNS zone %q not found", m.BaseDomain)
	}
	return m.cisInstanceCRN, nil
}

// Client returns a client used for making API calls to IBM Cloud services.
func (m *Metadata) Client() (*Client, error) {
	if m.client == nil {
		client, err := NewClient()
		if err != nil {
			return nil, err
		}
		m.client = client
	}
	return m.client, nil
}

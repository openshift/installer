package ibmcloud

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/go-sdk-core/v5/core"
)

// Metadata holds additional metadata for InstallConfig resources that
// does not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	BaseDomain              string
	ComputeSubnetNames      []string
	ControlPlaneSubnetNames []string
	Region                  string

	accountID           string
	cisInstanceCRN      string
	client              *Client
	computeSubnets      map[string]Subnet
	controlPlaneSubnets map[string]Subnet

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(baseDomain string, region string, controlPlaneSubnets []string, computeSubnets []string) *Metadata {
	return &Metadata{
		BaseDomain:              baseDomain,
		ComputeSubnetNames:      computeSubnets,
		ControlPlaneSubnetNames: controlPlaneSubnets,
		Region:                  region,
	}
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
				m.SetCISInstanceCRN(z.CISInstanceCRN)
				return m.cisInstanceCRN, nil
			}
		}
		return "", fmt.Errorf("cisInstanceCRN unknown due to DNS zone %q not found", m.BaseDomain)
	}
	return m.cisInstanceCRN, nil
}

// SetCISInstanceCRN sets Cloud Internet Services instance CRN to a string value.
func (m *Metadata) SetCISInstanceCRN(crn string) {
	m.cisInstanceCRN = crn
}

// ComputeSubnets gets the Subnet details for compute subnets
func (m *Metadata) ComputeSubnets(ctx context.Context) (map[string]Subnet, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.ComputeSubnetNames) > 0 && len(m.computeSubnets) == 0 {
		client, err := m.Client()
		if err != nil {
			return nil, err
		}
		m.computeSubnets, err = getSubnets(ctx, client, m.Region, m.ComputeSubnetNames)
		if err != nil {
			return nil, err
		}
	}

	return m.computeSubnets, nil
}

// ControlPlaneSubnets gets the Subnet details for control plane subnets
func (m *Metadata) ControlPlaneSubnets(ctx context.Context) (map[string]Subnet, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.ControlPlaneSubnetNames) > 0 && len(m.controlPlaneSubnets) == 0 {
		client, err := m.Client()
		if err != nil {
			return nil, err
		}
		m.controlPlaneSubnets, err = getSubnets(ctx, client, m.Region, m.ControlPlaneSubnetNames)
		if err != nil {
			return nil, err
		}
	}

	return m.controlPlaneSubnets, nil
}

// Client returns a client used for making API calls to IBM Cloud services.
func (m *Metadata) Client() (*Client, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.client == nil {
		client, err := NewClient()
		if err != nil {
			return nil, err
		}
		err = client.SetVPCServiceURLForRegion(context.TODO(), m.Region)
		if err != nil {
			return nil, err
		}
		m.client = client
	}
	return m.client, nil
}

// NewIamAuthenticator returns a new IamAuthenticator for using IBM Cloud services.
func NewIamAuthenticator(apiKey string) (*core.IamAuthenticator, error) {
	return core.NewIamAuthenticatorBuilder().SetApiKey(apiKey).Build()
}

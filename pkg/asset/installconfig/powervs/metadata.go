package powervs

import (
	"context"
	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
	"sync"
)

//go:generate mockgen -source=./metadata.go -destination=./mock/powervsmetadata_generated.go -package=mock

// MetadataAPI represents functions that eventually call out to the API
type MetadataAPI interface {
	AccountID(ctx context.Context) (string, error)
	APIKey(ctx context.Context) (string, error)
	CISInstanceCRN(ctx context.Context) (string, error)
	DNSInstanceCRN(ctx context.Context) (string, error)
}

// Metadata holds additional metadata for InstallConfig resources that
// do not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	BaseDomain string

	accountID      string
	apiKey         string
	cisInstanceCRN string
	dnsInstanceCRN string
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

	if m.client == nil {
		client, err := NewClient()
		if err != nil {
			return "", err
		}

		m.client = client
	}

	if m.accountID == "" {
		apiKeyDetails, err := m.client.GetAuthenticatorAPIKeyDetails(ctx)
		if err != nil {
			return "", err
		}

		m.accountID = *apiKeyDetails.AccountID
	}

	return m.accountID, nil
}

// APIKey returns the IBM Cloud account API Key associated with the authentication
// credentials.
func (m *Metadata) APIKey(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.client == nil {
		client, err := NewClient()
		if err != nil {
			return "", err
		}

		m.client = client
	}

	if m.apiKey == "" {
		m.apiKey = m.client.GetAPIKey()
	}

	return m.apiKey, nil
}

// CISInstanceCRN returns the Cloud Internet Services instance CRN that is
// managing the DNS zone for the base domain.
func (m *Metadata) CISInstanceCRN(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var err error
	if m.client == nil {
		client, err := NewClient()
		if err != nil {
			return "", err
		}

		m.client = client
	}

	if m.cisInstanceCRN == "" {
		m.cisInstanceCRN, err = m.client.GetInstanceCRNByName(ctx, m.BaseDomain, types.ExternalPublishingStrategy)
		if err != nil {
			return "", err
		}
	}
	return m.cisInstanceCRN, nil
}

// SetCISInstanceCRN sets Cloud Internet Services instance CRN to a string value.
func (m *Metadata) SetCISInstanceCRN(crn string) {
	m.cisInstanceCRN = crn
}

// DNSInstanceCRN returns the IBM DNS Service instance CRN that is
// managing the DNS zone for the base domain.
func (m *Metadata) DNSInstanceCRN(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var err error
	if m.client == nil {
		client, err := NewClient()
		if err != nil {
			return "", err
		}

		m.client = client
	}

	if m.dnsInstanceCRN == "" {
		m.dnsInstanceCRN, err = m.client.GetInstanceCRNByName(ctx, m.BaseDomain, types.InternalPublishingStrategy)
		if err != nil {
			return "", err
		}
	}

	return m.dnsInstanceCRN, nil
}

// SetDNSInstanceCRN sets IBM DNS Service instance CRN to a string value.
func (m *Metadata) SetDNSInstanceCRN(crn string) {
	m.dnsInstanceCRN = crn
}

// GetExistingVPCGateway checks if the VPC is a Permitted Network for the DNS Zone
func (m *Metadata) GetExistingVPCGateway(ctx context.Context, vpcName string, vpcSubnet string) (string, bool, error) {
	if vpcName == "" {
		return "", false, nil
	}

	vpc, err := m.client.GetVPCByName(ctx, vpcName)
	if err != nil {
		return "", false, errors.Wrap(err, "failed to get VPC")
	}

	vpcCRN, err := crn.Parse(*vpc.CRN)
	if err != nil {
		return "", false, errors.Wrap(err, "failed to parse VPC CRN")
	}

	subnet, err := m.client.GetSubnetByName(ctx, vpcSubnet, vpcCRN.Region)
	if err != nil {
		return "", false, errors.Wrap(err, "failed to get subnet")
	}
	// Check if subnet has an attached public gateway. If it does, we're done.
	if subnet.PublicGateway != nil {
		return *subnet.PublicGateway.Name, true, nil
	}

	// Check if a gateway exists in the VPN that isn't attached
	gw, err := m.client.GetPublicGatewayByVPC(ctx, vpcName)
	if err != nil {
		return "", false, errors.Wrap(err, "failed to get find gw")
	}
	// Found an unattached gateway
	if gw != nil {
		return *gw.Name, false, nil
	}
	return "", false, nil
}

// IsVPCPermittedNetwork checks if the VPC is a Permitted Network for the DNS Zone
func (m *Metadata) IsVPCPermittedNetwork(ctx context.Context, vpcName string, baseDomain string) (bool, error) {
	// An empty pre-existing VPC Name signifies a new VPC will be created (not pre-existing), so it won't be permitted
	if vpcName == "" {
		return false, nil
	}

	// Collect DNSInstance details if not already collected
	if m.dnsInstanceCRN == "" {
		_, err := m.DNSInstanceCRN(ctx)
		if err != nil {
			return false, errors.Wrap(err, "cannot collect DNS permitted networks without DNS Instance")
		}
	}

	if m.client == nil {
		client, err := NewClient()
		if err != nil {
			return false, err
		}

		m.client = client
	}

	// Get CIS zone ID by name
	zoneID, err := m.client.GetDNSZoneIDByName(context.TODO(), baseDomain, types.InternalPublishingStrategy)
	if err != nil {
		return false, errors.Wrap(err, "failed to get DNS zone ID")
	}

	dnsCRN, err := crn.Parse(m.dnsInstanceCRN)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse DNSInstanceCRN")
	}

	networks, err := m.client.GetDNSInstancePermittedNetworks(ctx, dnsCRN.ServiceInstance, zoneID)
	if err != nil {
		return false, err
	}
	if len(networks) < 1 {
		return false, nil
	}

	vpc, err := m.client.GetVPCByName(ctx, vpcName)
	for _, network := range networks {
		if network == *vpc.CRN {
			return true, nil
		}
	}

	return false, nil
}

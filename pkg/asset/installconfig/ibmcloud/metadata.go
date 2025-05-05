package ibmcloud

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
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
	client              API
	computeSubnets      map[string]Subnet
	controlPlaneSubnets map[string]Subnet
	dnsInstance         *DNSInstance

	mutex       sync.Mutex
	clientMutex sync.Mutex
}

// DNSInstance holds information for a DNS Services instance
type DNSInstance struct {
	ID   string
	CRN  string
	Zone string
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

		zones, err := client.GetDNSZones(ctx, types.ExternalPublishingStrategy)
		if err != nil {
			return "", err
		}

		for _, z := range zones {
			if z.Name == m.BaseDomain {
				m.cisInstanceCRN = z.InstanceCRN
				return m.cisInstanceCRN, nil
			}
		}
		return "", fmt.Errorf("cisInstanceCRN unknown due to DNS zone %q not found", m.BaseDomain)
	}
	return m.cisInstanceCRN, nil
}

// DNSInstance returns a DNSInstance holding information about the DNS Services instance
// managing the DNS zone for the base domain.
func (m *Metadata) DNSInstance(ctx context.Context) (*DNSInstance, error) {
	if m.dnsInstance != nil {
		return m.dnsInstance, nil
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Prevent multiple attempts to retrieve (set) the dnsInstance if it hasn't been set (multiple threads reach mutex concurrently)
	if m.dnsInstance == nil {
		client, err := m.Client()
		if err != nil {
			return nil, err
		}

		zones, err := client.GetDNSZones(ctx, types.InternalPublishingStrategy)
		if err != nil {
			return nil, err
		}

		for _, z := range zones {
			if z.Name == m.BaseDomain {
				if z.InstanceID == "" || z.InstanceCRN == "" {
					return nil, fmt.Errorf("dnsInstance has unknown ID/CRN: %q - %q", z.InstanceID, z.InstanceCRN)
				}
				m.dnsInstance = &DNSInstance{
					ID:   z.InstanceID,
					CRN:  z.InstanceCRN,
					Zone: z.ID,
				}
				return m.dnsInstance, nil
			}
		}
		return nil, fmt.Errorf("dnsInstance unknown due to DNS zone %q not found", m.BaseDomain)
	}
	return m.dnsInstance, nil
}

// IsVPCPermittedNetwork checks if the VPC is a Permitted Network for the DNS Zone
func (m *Metadata) IsVPCPermittedNetwork(ctx context.Context, vpcName string) (bool, error) {
	// An empty pre-existing VPC Name signifies a new VPC will be created (not pre-existing), so it won't be permitted
	if vpcName == "" {
		return false, nil
	}
	// Collect DNSInstance details if not already collected
	if m.dnsInstance == nil {
		_, err := m.DNSInstance(ctx)
		if err != nil {
			return false, errors.Wrap(err, "cannot collect DNS permitted networks without DNS Instance")
		}
	}

	client, err := m.Client()
	if err != nil {
		return false, err
	}

	networks, err := client.GetDNSInstancePermittedNetworks(ctx, m.dnsInstance.ID, m.dnsInstance.Zone)
	if err != nil {
		return false, err
	}
	if len(networks) < 1 {
		return false, nil
	}

	vpc, err := client.GetVPCByName(ctx, vpcName)
	if err != nil {
		return false, err
	}
	for _, network := range networks {
		if network == *vpc.CRN {
			return true, nil
		}
	}

	return false, nil
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
func (m *Metadata) Client() (API, error) {
	if m.client != nil {
		return m.client, nil
	}

	m.clientMutex.Lock()
	defer m.clientMutex.Unlock()

	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	err = client.SetVPCServiceURLForRegion(context.TODO(), m.Region)
	if err != nil {
		return nil, err
	}
	m.client = client
	return m.client, nil
}

// NewIamAuthenticator returns a new IamAuthenticator for using IBM Cloud services.
func NewIamAuthenticator(apiKey string) (*core.IamAuthenticator, error) {
	return core.NewIamAuthenticatorBuilder().SetApiKey(apiKey).Build()
}

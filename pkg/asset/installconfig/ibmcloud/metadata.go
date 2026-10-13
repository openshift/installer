package ibmcloud

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
)

const (
	// PrivateHostPrefix is the prefix for private API traffic used for DNS Records.
	PrivateHostPrefix = "api-int."
	// PublicHostPrefix is the prefix for public API traffic used for DNS Records.
	PublicHostPrefix = "api."
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
	publishStrategy     types.PublishingStrategy
	serviceEndpoints    []configv1.IBMCloudServiceEndpoint

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
func NewMetadata(config *types.InstallConfig) *Metadata {
	return &Metadata{
		BaseDomain:              config.BaseDomain,
		ComputeSubnetNames:      config.Platform.IBMCloud.ComputeSubnets,
		ControlPlaneSubnetNames: config.Platform.IBMCloud.ControlPlaneSubnets,
		publishStrategy:         config.Publish,
		Region:                  config.Platform.IBMCloud.Region,
		serviceEndpoints:        config.Platform.IBMCloud.ServiceEndpoints,
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

// AddVPCToPermittedNetworks adds a VPC to the DNS Services Zone of Permitted Networks.
func (m *Metadata) AddVPCToPermittedNetworks(ctx context.Context, vpcID string) error {
	// The following values are required to add the VPC to Permitted Networks:
	// - DNS Services Instance ID
	// - DNS Services Zone ID
	// - VPC CRN
	// We should already have the Instance ID and Zone ID, or they will be fetched.
	dnsInstance, err := m.DNSInstance(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve dns services instance: %w", err)
	}

	client, err := m.Client()
	if err != nil {
		return fmt.Errorf("failed to create ibmcloud client: %w", err)
	}
	// Lookup the VPC CRN.
	vpcDetails, err := client.GetVPC(ctx, vpcID)
	if err != nil {
		return fmt.Errorf("failed to retrieve vpc %s: %w", vpcID, err)
	} else if vpcDetails.CRN == nil {
		return fmt.Errorf("error vpc crn not set for vpc %s", vpcID)
	}

	if err = client.CreateDNSServicesPermittedNetwork(ctx, dnsInstance.ID, dnsInstance.Zone, *vpcDetails.CRN); err != nil {
		return fmt.Errorf("failed to add vpc %s to permitted networks of dns services zone %s and instance %s: %w", vpcID, dnsInstance.Zone, dnsInstance.ID, err)
	}

	return nil
}

// CreateDNSRecord creates a CNAME DNS Record in the IBM Cloud Internet Services zone or DNS Services zone for a Load Balancer hostname, based on the PublishStrategy.
func (m *Metadata) CreateDNSRecord(ctx context.Context, recordName string, loadBalancer *vpcv1.LoadBalancer) error {
	client, err := m.Client()
	if err != nil {
		return fmt.Errorf("failed to create ibmcloud client: %w", err)
	}

	zoneID, err := client.GetDNSZoneIDByName(ctx, m.BaseDomain, m.publishStrategy)
	if err != nil {
		return fmt.Errorf("failed to retrieve dns zone by base domain %s for %s cluster: %w", m.BaseDomain, m.publishStrategy, err)
	}

	switch m.publishStrategy {
	case types.ExternalPublishingStrategy:
		cisInstanceCRN, err := m.CISInstanceCRN(ctx)
		if err != nil {
			return fmt.Errorf("failed to retrieve cis instance crn for dns record: %w", err)
		}
		return client.CreateCISDNSRecord(ctx, cisInstanceCRN, zoneID, recordName, *loadBalancer.Hostname)
	case types.InternalPublishingStrategy:
		dnsInstance, err := m.DNSInstance(ctx)
		if err != nil {
			return fmt.Errorf("failed to retrieve dns instance for dns record: %w", err)
		}
		return client.CreateDNSServicesDNSRecord(ctx, dnsInstance.ID, zoneID, recordName, *loadBalancer.Hostname)
	default:
		return fmt.Errorf("failed to create dns record, invalid publish strategy: %s", m.publishStrategy)
	}
}

// CISInstanceCRN returns the Cloud Internet Services instance CRN that is
// managing the DNS zone for the base domain.
func (m *Metadata) CISInstanceCRN(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Only attempt to find the CIS instance if using ExternalPublishingStrategy and we have not collected it already
	if m.publishStrategy == types.ExternalPublishingStrategy && m.cisInstanceCRN == "" {
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

	// Only attempt to find the DNS Services instance if using InternalPublishingStrategy and also
	// prevent multiple attempts to retrieve (set) the dnsInstance if it hasn't been set (multiple threads reach mutex concurrently)
	if m.publishStrategy == types.InternalPublishingStrategy && m.dnsInstance == nil {
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

// GetIAMToken will retrieve an IAM access token using an IAM Authenticator and API Key.
func (m *Metadata) GetIAMToken(apiKey string) (*string, error) {
	// Get the IAM Service endpoint override, if one was supplied for the authenticator.
	authenticator, err := NewIamAuthenticator(apiKey, ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceIAM, m.serviceEndpoints))
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticator to get iam token: %w", err)
	}

	token, err := authenticator.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get iam token: %w", err)
	}
	return ptr.To(token), nil
}

// Client returns a client used for making API calls to IBM Cloud services.
func (m *Metadata) Client() (API, error) {
	if m.client != nil {
		return m.client, nil
	}

	m.clientMutex.Lock()
	defer m.clientMutex.Unlock()

	client, err := NewClient(m.serviceEndpoints)
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
func NewIamAuthenticator(apiKey string, iamServiceEndpointOverride string) (*core.IamAuthenticator, error) {
	if iamServiceEndpointOverride != "" {
		return core.NewIamAuthenticatorBuilder().SetApiKey(apiKey).SetURL(iamServiceEndpointOverride).Build()
	}
	return core.NewIamAuthenticatorBuilder().SetApiKey(apiKey).Build()
}

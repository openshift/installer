package powervs

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
)

const (
	cosCrnTempl = "crn:v1:bluemix:public:cloud-object-storage:global:::endpoint:s3.direct.%s.cloud-object-storage.appdomain.cloud"
	dnsCrn      = "crn:v1:bluemix:public:dns-svcs:global::::"
	iamCrn      = "crn:v1:bluemix:public:iam-svcs:global:::endpoint:private.iam.cloud.ibm.com"
	rcCrn       = "crn:v1:bluemix:public:resource-controller:global:::endpoint:private.resource-controller.cloud.ibm.com"
	vpcCrnTempl = "crn:v1:bluemix:public:is:%s:::endpoint:%s.private.iaas.cloud.ibm.com"
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
	BaseDomain      string
	PublishStrategy types.PublishingStrategy

	accountID      string
	apiKey         string
	cisInstanceCRN string
	dnsInstanceCRN string
	sessionClient  *Client

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(config *types.InstallConfig) *Metadata {
	return &Metadata{BaseDomain: config.BaseDomain, PublishStrategy: config.Publish}
}

func (m *Metadata) client() (*Client, error) {
	if m.sessionClient != nil {
		return m.sessionClient, nil
	}

	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	m.sessionClient = client

	return m.sessionClient, nil
}

// AccountID returns the IBM Cloud account ID associated with the authentication
// credentials.
func (m *Metadata) AccountID(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.accountID == "" {
		client, err := m.client()
		if err != nil {
			return "", err
		}

		if client.BXCli.User == nil || client.BXCli.User.Account == "" {
			return "", fmt.Errorf("failed to get find account ID: %+v", client.BXCli.User)
		}
		m.accountID = client.BXCli.User.Account
	}

	return m.accountID, nil
}

// APIKey returns the IBM Cloud account API Key associated with the authentication
// credentials.
func (m *Metadata) APIKey(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.apiKey == "" {
		client, err := m.client()
		if err != nil {
			return "", err
		}

		m.apiKey = client.GetAPIKey()
	}

	return m.apiKey, nil
}

// CISInstanceCRN returns the Cloud Internet Services instance CRN that is
// managing the DNS zone for the base domain.
func (m *Metadata) CISInstanceCRN(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.PublishStrategy == types.ExternalPublishingStrategy && m.cisInstanceCRN == "" {
		client, err := m.client()
		if err != nil {
			return "", err
		}

		m.cisInstanceCRN, err = client.GetInstanceCRNByName(ctx, m.BaseDomain, types.ExternalPublishingStrategy)
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

	if m.PublishStrategy == types.InternalPublishingStrategy && m.dnsInstanceCRN == "" {
		client, err := m.client()
		if err != nil {
			return "", err
		}

		m.dnsInstanceCRN, err = client.GetInstanceCRNByName(ctx, m.BaseDomain, types.InternalPublishingStrategy)
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
	if vpcName == "" || vpcSubnet == "" {
		return "", false, nil
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, err := m.client()
	if err != nil {
		return "", false, err
	}

	vpc, err := client.GetVPCByName(ctx, vpcName)
	if err != nil {
		return "", false, fmt.Errorf("failed to get VPC: %w", err)
	}

	vpcCRN, err := crn.Parse(*vpc.CRN)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse VPC CRN: %w", err)
	}

	subnet, err := client.GetSubnetByName(ctx, vpcSubnet, vpcCRN.Region)
	if err != nil {
		return "", false, fmt.Errorf("failed to get subnet: %w", err)
	}
	// Check if subnet has an attached public gateway. If it does, we're done.
	if subnet.PublicGateway != nil {
		return *subnet.PublicGateway.Name, true, nil
	}

	// Check if a gateway exists in the VPN that isn't attached
	gw, err := client.GetPublicGatewayByVPC(ctx, vpcName)
	if err != nil {
		return "", false, fmt.Errorf("failed to get find gw: %w", err)
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
			return false, fmt.Errorf("cannot collect DNS permitted networks without DNS Instance: %w", err)
		}
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, err := m.client()
	if err != nil {
		return false, err
	}

	// Get CIS zone ID by name
	zoneID, err := client.GetDNSZoneIDByName(context.TODO(), baseDomain, types.InternalPublishingStrategy)
	if err != nil {
		return false, fmt.Errorf("failed to get DNS zone ID: %w", err)
	}

	dnsCRN, err := crn.Parse(m.dnsInstanceCRN)
	if err != nil {
		return false, fmt.Errorf("failed to parse DNSInstanceCRN: %w", err)
	}

	networks, err := client.GetDNSInstancePermittedNetworks(ctx, dnsCRN.ServiceInstance, zoneID)
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

// EnsureVPCIsPermittedNetwork checks if a VPC is permitted to the DNS zone and adds it if it is not.
func (m *Metadata) EnsureVPCIsPermittedNetwork(ctx context.Context, vpcName string) error {
	dnsCRN, err := crn.Parse(m.dnsInstanceCRN)
	if err != nil {
		return fmt.Errorf("failed to parse DNSInstanceCRN: %w", err)
	}

	isVPCPermittedNetwork, err := m.IsVPCPermittedNetwork(ctx, vpcName, m.BaseDomain)
	if err != nil {
		return fmt.Errorf("failed to determine if VPC is permitted network: %w", err)
	}

	if !isVPCPermittedNetwork {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		client, err := m.client()
		if err != nil {
			return err
		}

		vpc, err := client.GetVPCByName(ctx, vpcName)
		if err != nil {
			return fmt.Errorf("failed to find VPC by name: %w", err)
		}

		zoneID, err := client.GetDNSZoneIDByName(ctx, m.BaseDomain, types.InternalPublishingStrategy)
		if err != nil {
			return fmt.Errorf("failed to get DNS zone ID: %w", err)
		}
		err = client.AddVPCToPermittedNetworks(ctx, *vpc.CRN, dnsCRN.ServiceInstance, zoneID)
		if err != nil {
			return fmt.Errorf("failed to add permitted network: %w", err)
		}
	}
	return nil
}

// GetSubnetID gets the ID of a VPC subnet by name and region.
func (m *Metadata) GetSubnetID(ctx context.Context, subnetName string, vpcRegion string) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, err := m.client()
	if err != nil {
		return "", err
	}
	subnet, err := client.GetSubnetByName(ctx, subnetName, vpcRegion)
	if err != nil {
		return "", err
	}
	return *subnet.ID, err
}

// GetVPCSubnets gets a list of subnets in a VPC.
func (m *Metadata) GetVPCSubnets(ctx context.Context, vpcName string) ([]vpcv1.Subnet, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, err := m.client()
	if err != nil {
		return nil, err
	}

	vpc, err := client.GetVPCByName(ctx, vpcName)
	if err != nil {
		return nil, err
	}
	subnets, err := client.GetVPCSubnets(ctx, *vpc.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get VPC subnets: %w", err)
	}
	return subnets, err
}

// GetDNSServerIP gets the IP of a custom resolver for DNS use.
func (m *Metadata) GetDNSServerIP(ctx context.Context, vpcName string) (string, error) {
	if m.dnsInstanceCRN == "" {
		_, err := m.DNSInstanceCRN(ctx)
		if err != nil {
			return "", fmt.Errorf("unable to locate DNS instance: %w", err)
		}
	}
	dnsCRN, err := crn.Parse(m.dnsInstanceCRN)
	if err != nil {
		return "", fmt.Errorf("failed to parse DNSInstanceCRN: %w", err)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, err := m.client()
	if err != nil {
		return "", err
	}
	vpc, err := client.GetVPCByName(ctx, vpcName)
	if err != nil {
		return "", err
	}

	dnsServerIP, err := client.GetDNSCustomResolverIP(ctx, dnsCRN.ServiceInstance, *vpc.ID)
	if err != nil {
		// There is no custom resolver, try to create one.
		customResolverName := fmt.Sprintf("%s-custom-resolver", vpcName)
		customResolver, err := client.CreateDNSCustomResolver(ctx, customResolverName, dnsCRN.ServiceInstance, *vpc.ID)
		if err != nil {
			return "", err
		}
		// Wait for the custom resolver to be enabled.
		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}

		customResolverID := *customResolver.ID
		var lastErr error
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			customResolver, lastErr = client.EnableDNSCustomResolver(ctx, dnsCRN.ServiceInstance, customResolverID)
			if lastErr == nil {
				return true, nil
			}
			return false, nil
		})
		if err != nil {
			if lastErr != nil {
				err = lastErr
			}
			return "", fmt.Errorf("failed to enable custom resolver %s: %w", *customResolver.ID, err)
		}
		dnsServerIP = *customResolver.Locations[0].DnsServerIp
	}
	return dnsServerIP, nil
}

// CreateDNSRecord creates a CNAME record for the specified hostname and destination hostname.
func (m *Metadata) CreateDNSRecord(ctx context.Context, hostname string, destHostname string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, err := m.client()
	if err != nil {
		return err
	}

	instanceCRN, err := client.GetInstanceCRNByName(ctx, m.BaseDomain, m.PublishStrategy)
	if err != nil {
		return fmt.Errorf("failed to get InstanceCRN (%s) by name: %w", m.PublishStrategy, err)
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}

	var lastErr error
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		lastErr = client.CreateDNSRecord(ctx, m.PublishStrategy, instanceCRN, m.BaseDomain, hostname, destHostname)
		if lastErr == nil {
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		if lastErr != nil {
			err = lastErr
		}
		return fmt.Errorf("failed to create a DNS CNAME record (%s, %s): %w",
			hostname,
			destHostname,
			err)
	}
	return err
}

// ListSecurityGroupRules lists the rules created in the specified VPC.
func (m *Metadata) ListSecurityGroupRules(ctx context.Context, vpcID string) (*vpcv1.SecurityGroupRuleCollection, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, err := m.client()
	if err != nil {
		return nil, err
	}

	return client.ListSecurityGroupRules(ctx, vpcID)
}

// SetVPCServiceURLForRegion sets the URL for the VPC based on the specified region.
func (m *Metadata) SetVPCServiceURLForRegion(ctx context.Context, vpcRegion string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, err := m.client()
	if err != nil {
		return err
	}

	return client.SetVPCServiceURLForRegion(ctx, vpcRegion)
}

// AddSecurityGroupRule adds a security group rule to the specified VPC.
func (m *Metadata) AddSecurityGroupRule(ctx context.Context, rule *vpcv1.SecurityGroupRulePrototype, vpcID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}

	client, err := m.client()
	if err != nil {
		return err
	}

	var lastErr error
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		lastErr = client.AddSecurityGroupRule(ctx, vpcID, rule)
		return lastErr == nil, nil
	})

	if err != nil {
		if lastErr != nil {
			err = lastErr
		}
		return fmt.Errorf("failed to add security group rule: %w", err)
	}
	return err
}

func leftInContext(ctx context.Context) time.Duration {
	deadline, ok := ctx.Deadline()
	if !ok {
		return math.MaxInt64
	}

	return time.Until(deadline)
}

// SetDefaultPrivateServiceEndpoints sets service endpoint overrides as needed for Disconnected install.
func (m *Metadata) SetDefaultPrivateServiceEndpoints(ctx context.Context, overrides []configv1.PowerVSServiceEndpoint, cosRegion string, vpcRegion string) []configv1.PowerVSServiceEndpoint {
	overrides = addOverride(overrides, string(configv1.IBMCloudServiceCOS), fmt.Sprintf("https://s3.direct.%s.cloud-object-storage.appdomain.cloud", cosRegion))
	overrides = addOverride(overrides, string(configv1.IBMCloudServiceDNSServices), "https://api.private.dns-svcs.cloud.ibm.com")
	overrides = addOverride(overrides, string(configv1.IBMCloudServiceIAM), "https://private.iam.cloud.ibm.com")
	overrides = addOverride(overrides, "Power", fmt.Sprintf("https://private.%s.power-iaas.cloud.ibm.com", vpcRegion)) // FIXME confiv1.IBMCloudServicePower?
	overrides = addOverride(overrides, string(configv1.IBMCloudServiceResourceController), "https://private.resource-controller.cloud.ibm.com")
	overrides = addOverride(overrides, string(configv1.IBMCloudServiceResourceManager), "https://private.resource-controller.cloud.ibm.com")
	overrides = addOverride(overrides, string(configv1.IBMCloudServiceVPC), fmt.Sprintf("https://%s.private.iaas.cloud.ibm.com", vpcRegion))
	return overrides
}

func addOverride(overrides []configv1.PowerVSServiceEndpoint, name string, url string) []configv1.PowerVSServiceEndpoint {
	found := false
	for _, service := range overrides {
		if service.Name == name {
			found = true
			break
		}
	}
	if !found {
		return append(overrides, configv1.PowerVSServiceEndpoint{Name: name, URL: url})
	}
	return overrides
}

// CreateVirtualPrivateEndpointGateways checks and creates necessary VPEs in given target VPC for Disconnected install.
func (m *Metadata) CreateVirtualPrivateEndpointGateways(ctx context.Context, infraID string, region string, vpcID string, subnetID string, groupID string, endpointOverrides []configv1.PowerVSServiceEndpoint) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, err := m.client()
	if err != nil {
		return err
	}

	endpoint, err := m.fetchEndpointOverride(endpointOverrides, "IAM")
	if err != nil {
		return fmt.Errorf("failed to fetch endpoint override found for IAM: %w", err)
	}
	if endpoint == "" {
		name := fmt.Sprintf("%s-vpe-iam", infraID)
		logrus.Debugf("InfraReady: Ensuring VPE gateway for IAM %v", iamCrn)
		_, err = client.CreateVirtualPrivateEndpointGateway(ctx, name, vpcID, subnetID, groupID, iamCrn)
		if err != nil {
			return fmt.Errorf("failed to create VPE: %w", err)
		}
	}

	endpoint, err = m.fetchEndpointOverride(endpointOverrides, "COS")
	if err != nil {
		return fmt.Errorf("failed to fetch endpoint override found for COS: %w", err)
	}
	if endpoint == "" {
		name := fmt.Sprintf("%s-vpe-cos", infraID)
		cosCrn := fmt.Sprintf(cosCrnTempl, region)
		logrus.Debugf("InfraReady: Ensuring VPE gateway for COS %v", cosCrn)
		_, err = client.CreateVirtualPrivateEndpointGateway(ctx, name, vpcID, subnetID, groupID, cosCrn)
		if err != nil {
			return fmt.Errorf("failed to create VPE: %w", err)
		}
	}

	endpoint, err = m.fetchEndpointOverride(endpointOverrides, "DNSServices")
	if err != nil {
		return fmt.Errorf("failed to fetch endpoint override found for DNS: %w", err)
	}
	if endpoint == "" {
		name := fmt.Sprintf("%s-vpe-dns", infraID)
		logrus.Debugf("InfraReady: Ensuring VPE gateway for DNS services %v", dnsCrn)
		_, err = client.CreateVirtualPrivateEndpointGateway(ctx, name, vpcID, subnetID, groupID, dnsCrn)
		if err != nil {
			return fmt.Errorf("failed to create VPE: %w", err)
		}
	}

	endpoint, err = m.fetchEndpointOverride(endpointOverrides, string(configv1.IBMCloudServiceResourceController))
	if err != nil {
		return fmt.Errorf("failed to fetch endpoint override found for RC: %w", err)
	}
	if endpoint == "" {
		name := fmt.Sprintf("%s-vpe-rc", infraID)
		logrus.Debugf("InfraReady: Ensuring VPE gateway for RC %v", rcCrn)
		_, err = client.CreateVirtualPrivateEndpointGateway(ctx, name, vpcID, subnetID, groupID, rcCrn)
		if err != nil {
			return fmt.Errorf("failed to create VPE: %w", err)
		}
	}

	endpoint, err = m.fetchEndpointOverride(endpointOverrides, string(configv1.IBMCloudServiceVPC))
	if err != nil {
		return fmt.Errorf("failed to fetch endpoint override found for VPC: %w", err)
	}
	if endpoint == "" {
		name := fmt.Sprintf("%s-vpe-vpc", infraID)
		vpcCrn := fmt.Sprintf(vpcCrnTempl, region, region)
		logrus.Debugf("InfraReady: Ensuring VPE gateway for VPC %v", vpcCrn)
		_, err = client.CreateVirtualPrivateEndpointGateway(ctx, name, vpcID, subnetID, groupID, vpcCrn)
		if err != nil {
			return fmt.Errorf("failed to create VPE: %w", err)
		}
	}

	return nil
}

func (m *Metadata) fetchEndpointOverride(overrides []configv1.PowerVSServiceEndpoint, svcID string) (string, error) {
	for _, endpoint := range overrides {
		if endpoint.Name == svcID {
			return endpoint.URL, nil
		}
	}
	return "", nil
}

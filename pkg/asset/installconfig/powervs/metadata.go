package powervs

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/installer/pkg/types"
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
	client         *Client

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(config *types.InstallConfig) *Metadata {
	return &Metadata{BaseDomain: config.BaseDomain, PublishStrategy: config.Publish}
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
		if m.client.BXCli.User == nil || m.client.BXCli.User.Account == "" {
			return "", fmt.Errorf("failed to get find account ID: %+v", m.client.BXCli.User)
		}
		m.accountID = m.client.BXCli.User.Account
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

	if m.PublishStrategy == types.ExternalPublishingStrategy && m.cisInstanceCRN == "" {
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

	if m.PublishStrategy == types.InternalPublishingStrategy && m.dnsInstanceCRN == "" {
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
	if vpcName == "" || vpcSubnet == "" {
		return "", false, nil
	}

	vpc, err := m.client.GetVPCByName(ctx, vpcName)
	if err != nil {
		return "", false, fmt.Errorf("failed to get VPC: %w", err)
	}

	vpcCRN, err := crn.Parse(*vpc.CRN)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse VPC CRN: %w", err)
	}

	subnet, err := m.client.GetSubnetByName(ctx, vpcSubnet, vpcCRN.Region)
	if err != nil {
		return "", false, fmt.Errorf("failed to get subnet: %w", err)
	}
	// Check if subnet has an attached public gateway. If it does, we're done.
	if subnet.PublicGateway != nil {
		return *subnet.PublicGateway.Name, true, nil
	}

	// Check if a gateway exists in the VPN that isn't attached
	gw, err := m.client.GetPublicGatewayByVPC(ctx, vpcName)
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
		return false, fmt.Errorf("failed to get DNS zone ID: %w", err)
	}

	dnsCRN, err := crn.Parse(m.dnsInstanceCRN)
	if err != nil {
		return false, fmt.Errorf("failed to parse DNSInstanceCRN: %w", err)
	}

	networks, err := m.client.GetDNSInstancePermittedNetworks(ctx, dnsCRN.ServiceInstance, zoneID)
	if err != nil {
		return false, err
	}
	if len(networks) < 1 {
		return false, nil
	}

	vpc, err := m.client.GetVPCByName(ctx, vpcName)
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
		vpc, err := m.client.GetVPCByName(ctx, vpcName)
		if err != nil {
			return fmt.Errorf("failed to find VPC by name: %w", err)
		}

		zoneID, err := m.client.GetDNSZoneIDByName(ctx, m.BaseDomain, types.InternalPublishingStrategy)
		if err != nil {
			return fmt.Errorf("failed to get DNS zone ID: %w", err)
		}
		err = m.client.AddVPCToPermittedNetworks(ctx, *vpc.CRN, dnsCRN.ServiceInstance, zoneID)
		if err != nil {
			return fmt.Errorf("failed to add permitted network: %w", err)
		}
	}
	return nil
}

// GetSubnetID gets the ID of a VPC subnet by name and region.
func (m *Metadata) GetSubnetID(ctx context.Context, subnetName string, vpcRegion string) (string, error) {
	subnet, err := m.client.GetSubnetByName(ctx, subnetName, vpcRegion)
	if err != nil {
		return "", err
	}
	return *subnet.ID, err
}

// GetVPCSubnets gets a list of subnets in a VPC.
func (m *Metadata) GetVPCSubnets(ctx context.Context, vpcName string) ([]vpcv1.Subnet, error) {
	vpc, err := m.client.GetVPCByName(ctx, vpcName)
	if err != nil {
		return nil, err
	}
	subnets, err := m.client.GetVPCSubnets(ctx, *vpc.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get VPC subnets: %w", err)
	}
	return subnets, err
}

// GetDNSServerIP gets the IP of a custom resolver for DNS use.
func (m *Metadata) GetDNSServerIP(ctx context.Context, vpcName string) (string, error) {
	vpc, err := m.client.GetVPCByName(ctx, vpcName)
	if err != nil {
		return "", err
	}

	dnsCRN, err := crn.Parse(m.dnsInstanceCRN)
	if err != nil {
		return "", fmt.Errorf("failed to parse DNSInstanceCRN: %w", err)
	}
	dnsServerIP, err := m.client.GetDNSCustomResolverIP(ctx, dnsCRN.ServiceInstance, *vpc.ID)
	if err != nil {
		// There is no custom resolver, try to create one.
		customResolverName := fmt.Sprintf("%s-custom-resolver", vpcName)
		customResolver, err := m.client.CreateDNSCustomResolver(ctx, customResolverName, dnsCRN.ServiceInstance, *vpc.ID)
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
			customResolver, lastErr = m.client.EnableDNSCustomResolver(ctx, dnsCRN.ServiceInstance, customResolverID)
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
	instanceCRN, err := m.client.GetInstanceCRNByName(ctx, m.BaseDomain, m.PublishStrategy)
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
		lastErr = m.client.CreateDNSRecord(ctx, m.PublishStrategy, instanceCRN, m.BaseDomain, hostname, destHostname)
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
	return m.client.ListSecurityGroupRules(ctx, vpcID)
}

// SetVPCServiceURLForRegion sets the URL for the VPC based on the specified region.
func (m *Metadata) SetVPCServiceURLForRegion(ctx context.Context, vpcRegion string) error {
	return m.client.SetVPCServiceURLForRegion(ctx, vpcRegion)
}

// AddSecurityGroupRule adds a security group rule to the specified VPC.
func (m *Metadata) AddSecurityGroupRule(ctx context.Context, rule *vpcv1.SecurityGroupRulePrototype, vpcID string) error {
	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}

	var lastErr error
	err := wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		lastErr = m.client.AddSecurityGroupRule(ctx, vpcID, rule)
		if lastErr == nil {
			return true, nil
		}
		return false, nil
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

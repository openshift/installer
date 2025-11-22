package powervs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/datacenters"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/dnsrecordsv1"
	"github.com/IBM/networking-go-sdk/dnssvcsv1"
	"github.com/IBM/networking-go-sdk/dnszonesv1"
	"github.com/IBM/networking-go-sdk/resourcerecordsv1"
	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/types"
)

//go:generate mockgen -source=./client.go -destination=./mock/powervsclient_generated.go -package=mock

// API represents the calls made to the API.
type API interface {
	// DNS
	GetDNSRecordsByName(ctx context.Context, crnstr string, zoneID string, recordName string, publish types.PublishingStrategy) ([]DNSRecordResponse, error)
	GetDNSZoneIDByName(ctx context.Context, name string, publish types.PublishingStrategy) (string, error)
	GetDNSZones(ctx context.Context, publish types.PublishingStrategy) ([]DNSZoneResponse, error)
	GetDNSInstancePermittedNetworks(ctx context.Context, dnsID string, dnsZone string) ([]string, error)
	GetDNSCustomResolverIP(ctx context.Context, dnsID string, vpcID string) (string, error)
	CreateDNSCustomResolver(ctx context.Context, name string, dnsID string, vpcID string) (*dnssvcsv1.CustomResolver, error)
	EnableDNSCustomResolver(ctx context.Context, dnsID string, resolverID string) (*dnssvcsv1.CustomResolver, error)
	CreateDNSRecord(ctx context.Context, publish types.PublishingStrategy, crnstr string, baseDomain string, hostname string, cname string) error
	AddVPCToPermittedNetworks(ctx context.Context, vpcCRN string, dnsID string, dnsZone string) error

	// VPC
	GetVPCByName(ctx context.Context, vpcName string) (*vpcv1.VPC, error)
	GetVPCByID(ctx context.Context, vpcID string, region string) (*vpcv1.VPC, error)
	GetPublicGatewayByVPC(ctx context.Context, vpcName string) (*vpcv1.PublicGateway, error)
	SetVPCServiceURLForRegion(ctx context.Context, region string) error
	GetVPCs(ctx context.Context, region string) ([]vpcv1.VPC, error)
	GetVPCSubnets(ctx context.Context, vpcID string) ([]vpcv1.Subnet, error)

	// TG
	TransitGatewayNameToID(ctx context.Context, name string) (string, error)
	TransitGatewayIDValid(ctx context.Context, id string) error
	GetTGConnectionVPC(ctx context.Context, gatewayID string, vpcSubnetID string) (string, error)
	GetAttachedTransitGateway(ctx context.Context, svcInsID string) (string, error)

	// Data Center
	GetDatacenterCapabilities(ctx context.Context, region string) (map[string]bool, error)
	GetDatacenterSupportedSystems(ctx context.Context, region string) ([]string, error)

	// API
	GetAuthenticatorAPIKeyDetails(ctx context.Context) (*iamidentityv1.APIKey, error)
	GetAPIKey() string

	// Subnet
	GetSubnetByName(ctx context.Context, subnetName string, region string) (*vpcv1.Subnet, error)

	// Resource Groups
	ListResourceGroups(ctx context.Context) (*resourcemanagerv2.ResourceGroupList, error)

	// Service Instance
	ListServiceInstances(ctx context.Context) ([]string, error)
	ServiceInstanceGUIDToName(ctx context.Context, id string) (string, error)
	ServiceInstanceNameToGUID(ctx context.Context, name string) (string, error)

	// Security Group
	ListSecurityGroups(ctx context.Context, vpcID string, regions string) ([]vpcv1.SecurityGroup, error)
	ListSecurityGroupRules(ctx context.Context, securityGroupID string) (*vpcv1.SecurityGroupRuleCollection, error)
	AddSecurityGroupRule(ctx context.Context, vpcID string, sgID string, rule *vpcv1.SecurityGroupRulePrototype) error

	// SSH
	CreateSSHKey(ctx context.Context, serviceInstance string, zone string, sshKeyName string, sshKey string) error

	// Load Balancer
	AddIPToLoadBalancerPool(ctx context.Context, lbID string, poolName string, port int64, ip string) error

	// Virtual Private Endpoint Gateway
	CreateVirtualPrivateEndpointGateway(ctx context.Context, name string, vpcID string, subnetID string, rgID string, targetCRN string) (*vpcv1.EndpointGateway, error)
}

// Client makes calls to the PowerVS API.
type Client struct {
	APIKey            string
	BXCli             *BxClient
	managementAPI     *resourcemanagerv2.ResourceManagerV2
	controllerAPI     *resourcecontrollerv2.ResourceControllerV2
	vpcAPI            *vpcv1.VpcV1
	dnsServicesAPI    *dnssvcsv1.DnsSvcsV1
	transitGatewayAPI *transitgatewayapisv1.TransitGatewayApisV1
}

// cisServiceID is the Cloud Internet Services' catalog service ID.
const (
	cisServiceID          = "75874a60-cb12-11e7-948e-37ac098eb1b9"
	dnsServiceID          = "b4ed8a30-936f-11e9-b289-1d079699cbe5"
	serviceInstanceType   = "service_instance"
	compositeInstanceType = "composite_instance"
)

// DNSZoneResponse represents a DNS zone response.
type DNSZoneResponse struct {
	// Name is the domain name of the zone.
	Name string

	// ID is the zone's ID.
	ID string

	// CISInstanceCRN is the IBM Cloud Resource Name for the CIS instance where
	// the DNS zone is managed.
	InstanceCRN string

	// CISInstanceName is the display name of the CIS instance where the DNS zone
	// is managed.
	InstanceName string

	// ResourceGroupID is the resource group ID of the CIS instance.
	ResourceGroupID string
}

// DNSRecordResponse represents a DNS record response.
type DNSRecordResponse struct {
	Name string
	Type string
}

// NewClient initializes a client with a session.
func NewClient() (*Client, error) {
	bxCli, err := NewBxClient(false)
	if err != nil {
		return nil, err
	}

	client := &Client{
		APIKey: bxCli.APIKey,
		BXCli:  bxCli,
	}

	if err := client.loadSDKServices(); err != nil {
		return nil, fmt.Errorf("failed to load IBM SDK services: %w", err)
	}

	if bxCli.PowerVSResourceGroup == "Default" {
		// Here we are initialized enough to handle a default resource group
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
		defer cancel()

		resourceGroups, err := client.ListResourceGroups(ctx)
		if err != nil {
			return nil, fmt.Errorf("client.ListResourceGroups failed: %w", err)
		}
		if resourceGroups == nil {
			return nil, errors.New("client.ListResourceGroups returns nil")
		}

		found := false
		for _, resourceGroup := range resourceGroups.Resources {
			if resourceGroup.Default != nil && *resourceGroup.Default {
				bxCli.PowerVSResourceGroup = *resourceGroup.Name
				found = true
				break
			}
		}

		if !found {
			return nil, errors.New("no default resource group found")
		}
	}

	return client, nil
}

func (c *Client) loadSDKServices() error {
	servicesToLoad := []func() error{
		c.loadResourceManagementAPI,
		c.loadResourceControllerAPI,
		c.loadVPCV1API,
		c.loadDNSServicesAPI,
		c.loadTransitGatewayAPI,
	}

	// Call all the load functions.
	for _, fn := range servicesToLoad {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

// GetDNSRecordsByName gets DNS records in specific Cloud Internet Services instance
// by its CRN, zone ID, and DNS record name.
func (c *Client) GetDNSRecordsByName(ctx context.Context, crnstr string, zoneID string, recordName string, publish types.PublishingStrategy) ([]DNSRecordResponse, error) {
	authenticator := &core.IamAuthenticator{
		ApiKey: c.APIKey,
	}
	dnsRecords := []DNSRecordResponse{}
	switch publish {
	case types.ExternalPublishingStrategy:
		// Set CIS DNS record service
		dnsService, err := dnsrecordsv1.NewDnsRecordsV1(&dnsrecordsv1.DnsRecordsV1Options{
			Authenticator:  authenticator,
			Crn:            core.StringPtr(crnstr),
			ZoneIdentifier: core.StringPtr(zoneID),
		})
		if err != nil {
			return nil, err
		}

		// Get CIS DNS records by name
		records, _, err := dnsService.ListAllDnsRecordsWithContext(ctx, &dnsrecordsv1.ListAllDnsRecordsOptions{
			Name: core.StringPtr(recordName),
		})
		if err != nil {
			return nil, fmt.Errorf("could not retrieve DNS records: %w", err)
		}
		for _, record := range records.Result {
			dnsRecords = append(dnsRecords, DNSRecordResponse{Name: *record.Name, Type: *record.Type})
		}
	case types.InternalPublishingStrategy:
		// Set DNS record service
		dnsService, err := resourcerecordsv1.NewResourceRecordsV1(&resourcerecordsv1.ResourceRecordsV1Options{
			Authenticator: authenticator,
		})
		if err != nil {
			return nil, err
		}

		dnsCRN, err := crn.Parse(crnstr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DNSInstanceCRN: %w", err)
		}

		// Get DNS records by name
		records, _, err := dnsService.ListResourceRecords(&resourcerecordsv1.ListResourceRecordsOptions{
			InstanceID: &dnsCRN.ServiceInstance,
			DnszoneID:  &zoneID,
		})
		for _, record := range records.ResourceRecords {
			if *record.Name == recordName {
				dnsRecords = append(dnsRecords, DNSRecordResponse{Name: *record.Name, Type: *record.Type})
			}
		}
		if err != nil {
			return nil, fmt.Errorf("could not retrieve DNS records: %w", err)
		}
	}

	return dnsRecords, nil
}

// GetInstanceCRNByName finds the CRN of the instance with the specified name.
func (c *Client) GetInstanceCRNByName(ctx context.Context, name string, publish types.PublishingStrategy) (string, error) {
	zones, err := c.GetDNSZones(ctx, publish)
	if err != nil {
		return "", err
	}

	for _, z := range zones {
		if z.Name == name {
			return z.InstanceCRN, nil
		}
	}

	return "", fmt.Errorf("DNS zone %q not found", name)
}

// GetDNSCustomResolverIP gets the DNS Server IP of a custom resolver associated with the specified VPC subnet in the specified DNS zone.
func (c *Client) GetDNSCustomResolverIP(ctx context.Context, dnsID string, vpcID string) (string, error) {
	listCustomResolversOptions := c.dnsServicesAPI.NewListCustomResolversOptions(dnsID)
	customResolvers, _, err := c.dnsServicesAPI.ListCustomResolversWithContext(ctx, listCustomResolversOptions)
	if err != nil {
		return "", err
	}

	subnets, err := c.GetVPCSubnets(ctx, vpcID)
	if err != nil {
		return "", err
	}

	for _, customResolver := range customResolvers.CustomResolvers {
		for _, location := range customResolver.Locations {
			for _, subnet := range subnets {
				if *subnet.CRN == *location.SubnetCrn {
					return *location.DnsServerIp, nil
				}
			}
		}
	}
	return "", fmt.Errorf("DNS server IP of custom resolver for %q not found", dnsID)
}

// CreateDNSCustomResolver creates a custom resolver associated with the specified VPC in the specified DNS zone.
func (c *Client) CreateDNSCustomResolver(ctx context.Context, name string, dnsID string, vpcID string) (*dnssvcsv1.CustomResolver, error) {
	createCustomResolverOptions := c.dnsServicesAPI.NewCreateCustomResolverOptions(dnsID, name)

	subnets, err := c.GetVPCSubnets(ctx, vpcID)
	if err != nil {
		return nil, err
	}

	locations := []dnssvcsv1.LocationInput{}
	for _, subnet := range subnets {
		location, err := c.dnsServicesAPI.NewLocationInput(*subnet.CRN)
		if err != nil {
			return nil, err
		}
		location.Enabled = core.BoolPtr(true)
		locations = append(locations, *location)
	}
	createCustomResolverOptions.SetLocations(locations)

	customResolver, _, err := c.dnsServicesAPI.CreateCustomResolverWithContext(ctx, createCustomResolverOptions)
	if err != nil {
		return nil, err
	}
	return customResolver, nil
}

// EnableDNSCustomResolver enables a specified custom resolver.
func (c *Client) EnableDNSCustomResolver(ctx context.Context, dnsID string, resolverID string) (*dnssvcsv1.CustomResolver, error) {
	updateCustomResolverOptions := c.dnsServicesAPI.NewUpdateCustomResolverOptions(dnsID, resolverID)
	updateCustomResolverOptions.SetEnabled(true)

	customResolver, _, err := c.dnsServicesAPI.UpdateCustomResolverWithContext(ctx, updateCustomResolverOptions)
	if err != nil {
		return nil, err
	}
	return customResolver, nil
}

// GetDNSZoneIDByName gets the CIS zone ID from its domain name.
func (c *Client) GetDNSZoneIDByName(ctx context.Context, name string, publish types.PublishingStrategy) (string, error) {
	zones, err := c.GetDNSZones(ctx, publish)
	if err != nil {
		return "", err
	}

	for _, z := range zones {
		if z.Name == name {
			return z.ID, nil
		}
	}

	return "", fmt.Errorf("DNS zone %q not found", name)
}

// GetDNSZones returns all of the active DNS zones managed by CIS.
func (c *Client) GetDNSZones(ctx context.Context, publish types.PublishingStrategy) ([]DNSZoneResponse, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	options := c.controllerAPI.NewListResourceInstancesOptions()
	switch publish {
	case types.ExternalPublishingStrategy:
		options.SetResourceID(cisServiceID)
	case types.InternalPublishingStrategy:
		options.SetResourceID(dnsServiceID)
	default:
		return nil, errors.New("unknown publishing strategy")
	}

	listResourceInstancesResponse, _, err := c.controllerAPI.ListResourceInstances(options)
	if err != nil {
		return nil, fmt.Errorf("failed to get cis instance: %w", err)
	}

	var allZones []DNSZoneResponse
	for _, instance := range listResourceInstancesResponse.Resources {
		authenticator := &core.IamAuthenticator{
			ApiKey: c.APIKey,
		}

		switch publish {
		case types.ExternalPublishingStrategy:
			zonesService, err := zonesv1.NewZonesV1(&zonesv1.ZonesV1Options{
				Authenticator: authenticator,
				Crn:           instance.CRN,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to list DNS zones: %w", err)
			}

			options := zonesService.NewListZonesOptions()
			listZonesResponse, _, err := zonesService.ListZones(options)

			if listZonesResponse == nil {
				return nil, err
			}

			for _, zone := range listZonesResponse.Result {
				if *zone.Status == "active" {
					zoneStruct := DNSZoneResponse{
						Name:            *zone.Name,
						ID:              *zone.ID,
						InstanceCRN:     *instance.CRN,
						InstanceName:    *instance.Name,
						ResourceGroupID: *instance.ResourceGroupID,
					}
					allZones = append(allZones, zoneStruct)
				}
			}
		case types.InternalPublishingStrategy:
			dnsZonesService, err := dnszonesv1.NewDnsZonesV1(&dnszonesv1.DnsZonesV1Options{
				Authenticator: authenticator,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to list DNS zones: %w", err)
			}

			options := dnsZonesService.NewListDnszonesOptions(*instance.GUID)
			listZonesResponse, _, err := dnsZonesService.ListDnszones(options)

			if listZonesResponse == nil {
				return nil, err
			}

			for _, zone := range listZonesResponse.Dnszones {
				if *zone.State == "ACTIVE" || *zone.State == "PENDING_NETWORK_ADD" {
					zoneStruct := DNSZoneResponse{
						Name:            *zone.Name,
						ID:              *zone.ID,
						InstanceCRN:     *instance.CRN,
						InstanceName:    *instance.Name,
						ResourceGroupID: *instance.ResourceGroupID,
					}
					allZones = append(allZones, zoneStruct)
				}
			}
		}
	}
	return allZones, nil
}

// GetDNSInstancePermittedNetworks gets the permitted VPC networks for a DNS Services instance.
func (c *Client) GetDNSInstancePermittedNetworks(ctx context.Context, dnsID string, dnsZone string) ([]string, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	listPermittedNetworksOptions := c.dnsServicesAPI.NewListPermittedNetworksOptions(dnsID, dnsZone)
	permittedNetworks, _, err := c.dnsServicesAPI.ListPermittedNetworksWithContext(ctx, listPermittedNetworksOptions)
	if err != nil {
		return nil, err
	}

	networks := []string{}
	for _, network := range permittedNetworks.PermittedNetworks {
		networks = append(networks, *network.PermittedNetwork.VpcCrn)
	}
	return networks, nil
}

// AddVPCToPermittedNetworks adds the specified VPC to the specified DNS zone.
func (c *Client) AddVPCToPermittedNetworks(ctx context.Context, vpcCRN string, dnsID string, dnsZone string) error {
	permittedNetwork, err := c.dnsServicesAPI.NewPermittedNetworkVpc(vpcCRN)
	if err != nil {
		return err
	}

	createPermittedNetworkOptions := c.dnsServicesAPI.NewCreatePermittedNetworkOptions(dnsID, dnsZone, dnssvcsv1.CreatePermittedNetworkOptions_Type_Vpc, permittedNetwork)

	_, _, err = c.dnsServicesAPI.CreatePermittedNetworkWithContext(ctx, createPermittedNetworkOptions)
	if err != nil {
		return err
	}
	return nil
}

// CreateDNSRecord Creates a DNS CNAME record in the given base domain and CRN.
func (c *Client) CreateDNSRecord(ctx context.Context, publish types.PublishingStrategy, crnstr string, baseDomain string, hostname string, cname string) error {
	switch publish {
	case types.InternalPublishingStrategy:
		return c.createPrivateDNSRecord(ctx, crnstr, baseDomain, hostname, cname)
	case types.ExternalPublishingStrategy:
		return c.createPublicDNSRecord(ctx, crnstr, baseDomain, hostname, cname)
	default:
		return fmt.Errorf("publish strategy %q not supported", publish)
	}
}

func (c *Client) createPublicDNSRecord(ctx context.Context, crnstr string, baseDomain string, hostname string, cname string) error {
	logrus.Debugf("createDNSRecord: crnstr = %s, hostname = %s, cname = %s", crnstr, hostname, cname)

	var (
		zoneID           string
		err              error
		authenticator    *core.IamAuthenticator
		globalOptions    *dnsrecordsv1.DnsRecordsV1Options
		dnsRecordService *dnsrecordsv1.DnsRecordsV1
	)

	// Get CIS zone ID by name
	zoneID, err = c.GetDNSZoneIDByName(ctx, baseDomain, types.ExternalPublishingStrategy)
	if err != nil {
		logrus.Errorf("c.GetDNSZoneIDByName returns %v", err)
		return err
	}
	logrus.Debugf("CreatePublicDNSRecord: zoneID = %s", zoneID)

	authenticator = &core.IamAuthenticator{
		ApiKey: c.APIKey,
	}
	globalOptions = &dnsrecordsv1.DnsRecordsV1Options{
		Authenticator:  authenticator,
		Crn:            ptr.To(crnstr),
		ZoneIdentifier: ptr.To(zoneID),
	}
	dnsRecordService, err = dnsrecordsv1.NewDnsRecordsV1(globalOptions)
	if err != nil {
		logrus.Errorf("dnsrecordsv1.NewDnsRecordsV1 returns %v", err)
		return err
	}
	logrus.Debugf("CreatePublicDNSRecord: dnsRecordService = %+v", dnsRecordService)

	createOptions := dnsRecordService.NewCreateDnsRecordOptions()
	createOptions.SetName(hostname)
	createOptions.SetType(dnsrecordsv1.CreateDnsRecordOptions_Type_Cname)
	createOptions.SetContent(cname)

	result, response, err := dnsRecordService.CreateDnsRecord(createOptions)
	if err != nil {
		logrus.Errorf("dnsRecordService.CreateDnsRecord returns %v", err)
		return err
	}
	logrus.Debugf("createPublicDNSRecord: Result.ID = %v, RawResult = %v", *result.Result.ID, response.RawResult)

	return nil
}

func (c *Client) createPrivateDNSRecord(ctx context.Context, crnstr string, baseDomain string, hostname string, cname string) error {
	logrus.Debugf("createPrivateDNSRecord: crnstr = %s, hostname = %s, cname = %s", crnstr, hostname, cname)

	zoneID, err := c.GetDNSZoneIDByName(ctx, baseDomain, types.InternalPublishingStrategy)
	if err != nil {
		logrus.Errorf("c.GetDNSZoneIDByName returns %v", err)
		return err
	}
	logrus.Debugf("createPrivateDNSRecord: zoneID = %s", zoneID)

	dnsCRN, err := crn.Parse(crnstr)
	if err != nil {
		return fmt.Errorf("failed to parse DNSInstanceCRN: %w", err)
	}

	rdataCnameRecord, err := c.dnsServicesAPI.NewResourceRecordInputRdataRdataCnameRecord(cname)
	if err != nil {
		return fmt.Errorf("NewResourceRecordInputRdataRdataCnameRecord failed: %w", err)
	}
	createOptions := c.dnsServicesAPI.NewCreateResourceRecordOptions(dnsCRN.ServiceInstance, zoneID, dnssvcsv1.CreateResourceRecordOptions_Type_Cname)
	createOptions.SetRdata(rdataCnameRecord)
	createOptions.SetTTL(120)
	createOptions.SetName(hostname)
	result, resp, err := c.dnsServicesAPI.CreateResourceRecord(createOptions)
	if err != nil {
		logrus.Errorf("dnsRecordService.CreateResourceRecord returns %v", err)
		return err
	}
	logrus.Debugf("createPrivateDNSRecord: result.ID = %v, resp.RawResult = %v", *result.ID, resp.RawResult)

	return nil
}

// GetVPCByName gets a VPC by its name.
func (c *Client) GetVPCByName(ctx context.Context, vpcName string) (*vpcv1.VPC, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	listRegionsOptions := c.vpcAPI.NewListRegionsOptions()
	listRegionsResponse, _, err := c.vpcAPI.ListRegionsWithContext(ctx, listRegionsOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list vpc regions: %w", err)
	}
	var vpcNamesList []string
	for _, region := range listRegionsResponse.Regions {
		err := c.vpcAPI.SetServiceURL(fmt.Sprintf("%s/v1", *region.Endpoint))
		if err != nil {
			return nil, fmt.Errorf("failed to set vpc api service url: %w", err)
		}
		vpcs, detailedResponse, err := c.vpcAPI.ListVpcsWithContext(ctx, c.vpcAPI.NewListVpcsOptions())
		if err != nil {
			if detailedResponse.GetStatusCode() != http.StatusNotFound {
				return nil, err
			}
		} else {
			for _, vpc := range vpcs.Vpcs {
				vpcNamesList = append(vpcNamesList, *vpc.Name)
				if *vpc.Name == vpcName {
					return &vpc, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("failed to find VPC %q. Available VPCs: %v", vpcName, vpcNamesList)
}

// GetVPCByID checks if an id is a valid VPC id and, if so, returns the VPC.
func (c *Client) GetVPCByID(ctx context.Context, vpcID string, region string) (*vpcv1.VPC, error) {
	vpcs, err := c.GetVPCs(ctx, region)
	if err != nil {
		return nil, err
	}

	for _, vpc := range vpcs {
		if *vpc.ID == vpcID {
			return &vpc, nil
		}
	}

	return nil, fmt.Errorf("VPC with id (%s) does not exist in region (%s)", vpcID, region)
}

// GetPublicGatewayByVPC gets all PublicGateways in a region
func (c *Client) GetPublicGatewayByVPC(ctx context.Context, vpcName string) (*vpcv1.PublicGateway, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	vpc, err := c.GetVPCByName(ctx, vpcName)
	if err != nil {
		return nil, fmt.Errorf("failed to get VPC: %w", err)
	}

	vpcCRN, err := crn.Parse(*vpc.CRN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse VPC CRN: %w", err)
	}

	err = c.SetVPCServiceURLForRegion(ctx, vpcCRN.Region)
	if err != nil {
		return nil, err
	}

	listPublicGatewaysOptions := c.vpcAPI.NewListPublicGatewaysOptions()
	publicGatewayCollection, detailedResponse, err := c.vpcAPI.ListPublicGatewaysWithContext(ctx, listPublicGatewaysOptions)
	if err != nil {
		return nil, err
	} else if detailedResponse.GetStatusCode() == http.StatusNotFound {
		return nil, errors.New("failed to find publicGateways")
	}
	for _, gw := range publicGatewayCollection.PublicGateways {
		if *vpc.ID == *gw.VPC.ID {
			return &gw, nil
		}
	}

	return nil, nil
}

// GetVPCSubnets retrieves all subnets in the given VPC.
func (c *Client) GetVPCSubnets(ctx context.Context, vpcID string) ([]vpcv1.Subnet, error) {
	listSubnetsOptions := c.vpcAPI.NewListSubnetsOptions()
	listSubnetsOptions.VPCID = &vpcID
	subnets, _, err := c.vpcAPI.ListSubnetsWithContext(ctx, listSubnetsOptions)
	if err != nil {
		return nil, err
	}

	return subnets.Subnets, nil
}

// GetSubnetByName gets a VPC Subnet by its name and region.
func (c *Client) GetSubnetByName(ctx context.Context, subnetName string, region string) (*vpcv1.Subnet, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	err := c.SetVPCServiceURLForRegion(ctx, region)
	if err != nil {
		return nil, err
	}

	listSubnetsOptions := c.vpcAPI.NewListSubnetsOptions()
	subnetCollection, detailedResponse, err := c.vpcAPI.ListSubnetsWithContext(ctx, listSubnetsOptions)
	if err != nil {
		return nil, err
	} else if detailedResponse.GetStatusCode() == http.StatusNotFound {
		return nil, errors.New("failed to find VPC Subnet")
	}
	for _, subnet := range subnetCollection.Subnets {
		if subnetName == *subnet.Name {
			return &subnet, nil
		}
	}

	return nil, errors.New("failed to find VPC Subnet")
}

func (c *Client) loadResourceManagementAPI() error {
	authenticator := &core.IamAuthenticator{
		ApiKey: c.APIKey,
	}
	options := &resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: authenticator,
	}
	resourceManagerV2Service, err := resourcemanagerv2.NewResourceManagerV2(options)
	if err != nil {
		return err
	}
	c.managementAPI = resourceManagerV2Service
	return nil
}

func (c *Client) loadResourceControllerAPI() error {
	authenticator := &core.IamAuthenticator{
		ApiKey: c.APIKey,
	}
	options := &resourcecontrollerv2.ResourceControllerV2Options{
		Authenticator: authenticator,
	}
	resourceControllerV2Service, err := resourcecontrollerv2.NewResourceControllerV2(options)
	if err != nil {
		return err
	}
	c.controllerAPI = resourceControllerV2Service
	return nil
}

func (c *Client) loadVPCV1API() error {
	authenticator := &core.IamAuthenticator{
		ApiKey: c.APIKey,
	}
	vpcService, err := vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		Authenticator: authenticator,
	})
	if err != nil {
		return err
	}
	c.vpcAPI = vpcService
	return nil
}

func (c *Client) loadDNSServicesAPI() error {
	authenticator := &core.IamAuthenticator{
		ApiKey: c.APIKey,
	}
	dnsService, err := dnssvcsv1.NewDnsSvcsV1(&dnssvcsv1.DnsSvcsV1Options{
		Authenticator: authenticator,
	})
	if err != nil {
		return err
	}
	c.dnsServicesAPI = dnsService
	return nil
}

func (c *Client) loadTransitGatewayAPI() error {
	authenticator := &core.IamAuthenticator{
		ApiKey: c.APIKey,
	}
	versionDate := "2023-07-04"
	tgSvc, err := transitgatewayapisv1.NewTransitGatewayApisV1(&transitgatewayapisv1.TransitGatewayApisV1Options{
		Authenticator: authenticator,
		Version:       &versionDate,
	})
	if err != nil {
		return err
	}
	c.transitGatewayAPI = tgSvc
	return nil
}

// SetVPCServiceURLForRegion will set the VPC Service URL to a specific IBM Cloud Region, in order to access Region scoped resources
func (c *Client) SetVPCServiceURLForRegion(ctx context.Context, region string) error {
	regionOptions := c.vpcAPI.NewGetRegionOptions(region)
	vpcRegion, _, err := c.vpcAPI.GetRegionWithContext(ctx, regionOptions)
	if err != nil {
		return err
	}
	err = c.vpcAPI.SetServiceURL(fmt.Sprintf("%s/v1", *vpcRegion.Endpoint))
	if err != nil {
		return err
	}
	return nil
}

// GetAuthenticatorAPIKeyDetails gets detailed information on the API key used
// for authentication to the IBM Cloud APIs.
func (c *Client) GetAuthenticatorAPIKeyDetails(ctx context.Context) (*iamidentityv1.APIKey, error) {
	authenticator := &core.IamAuthenticator{
		ApiKey: c.APIKey,
	}
	iamIdentityService, err := iamidentityv1.NewIamIdentityV1(&iamidentityv1.IamIdentityV1Options{
		Authenticator: authenticator,
	})
	if err != nil {
		return nil, err
	}

	options := iamIdentityService.NewGetAPIKeysDetailsOptions()
	options.SetIamAPIKey(c.APIKey)
	details, _, err := iamIdentityService.GetAPIKeysDetailsWithContext(ctx, options)
	if err != nil {
		return nil, err
	}
	// NOTE: details.Apikey
	// https://cloud.ibm.com/apidocs/iam-identity-token-api?code=go#get-api-keys-details
	// This property only contains the API key value for the following cases: create an API key,
	// update a service ID API key that stores the API key value as retrievable, or get a service
	// ID API key that stores the API key value as retrievable. All other operations don't return
	// the API key value, for example all user API key related operations, except for create,
	// don't contain the API key value.
	return details, nil
}

// GetAPIKey returns the PowerVS API key
func (c *Client) GetAPIKey() string {
	return c.APIKey
}

// GetVPCs gets all VPCs in a region.
func (c *Client) GetVPCs(ctx context.Context, region string) ([]vpcv1.VPC, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	err := c.SetVPCServiceURLForRegion(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("failed to set vpc api service url: %w", err)
	}

	vpcs, _, err := c.vpcAPI.ListVpcs(c.vpcAPI.NewListVpcsOptions())
	if err != nil {
		return nil, err
	}

	return vpcs.Vpcs, nil
}

// ListResourceGroups returns a list of resource groups.
func (c *Client) ListResourceGroups(ctx context.Context) (*resourcemanagerv2.ResourceGroupList, error) {
	listResourceGroupsOptions := c.managementAPI.NewListResourceGroupsOptions()
	listResourceGroupsOptions.AccountID = &c.BXCli.User.Account

	resourceGroups, _, err := c.managementAPI.ListResourceGroups(listResourceGroupsOptions)
	if err != nil {
		return nil, err
	}

	return resourceGroups, err
}

const (
	// resource Id for Power Systems Virtual Server in the Global catalog.
	powerIAASResourceID = "abd259f0-9990-11e8-acc8-b9f54a8f1661"
)

// ListServiceInstances lists all service instances in the cloud.
func (c *Client) ListServiceInstances(ctx context.Context) ([]string, error) {
	var (
		serviceInstances []string
		options          *resourcecontrollerv2.ListResourceInstancesOptions
		resources        *resourcecontrollerv2.ResourceInstancesList
		err              error
		perPage          int64 = 10
		moreData               = true
		nextURL          *string
		groupID          = c.BXCli.PowerVSResourceGroup
	)

	// If the user passes in a human readable group id, then we need to convert it to a UUID
	listGroupOptions := c.managementAPI.NewListResourceGroupsOptions()
	listGroupOptions.AccountID = &c.BXCli.User.Account
	groups, _, err := c.managementAPI.ListResourceGroupsWithContext(ctx, listGroupOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list resource groups: %w", err)
	}
	for _, group := range groups.Resources {
		if *group.Name == groupID {
			groupID = *group.ID
		}
	}

	options = c.controllerAPI.NewListResourceInstancesOptions()
	options.SetResourceGroupID(groupID)
	// resource ID for Power Systems Virtual Server in the Global catalog
	options.SetResourceID(powerIAASResourceID)
	options.SetLimit(perPage)

	for moreData {
		resources, _, err = c.controllerAPI.ListResourceInstancesWithContext(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("failed to list resource instances: %w", err)
		}

		for _, resource := range resources.Resources {
			var (
				getResourceOptions *resourcecontrollerv2.GetResourceInstanceOptions
				resourceInstance   *resourcecontrollerv2.ResourceInstance
				response           *core.DetailedResponse
			)

			getResourceOptions = c.controllerAPI.NewGetResourceInstanceOptions(*resource.ID)

			resourceInstance, response, err = c.controllerAPI.GetResourceInstance(getResourceOptions)
			if err != nil {
				return nil, fmt.Errorf("failed to get instance: %w", err)
			}
			if response != nil && response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusInternalServerError {
				continue
			}

			if resourceInstance.Type != nil && (*resourceInstance.Type == serviceInstanceType || *resourceInstance.Type == compositeInstanceType) {
				serviceInstances = append(serviceInstances, fmt.Sprintf("%s %s", *resource.Name, *resource.GUID))
			}
		}

		// Based on: https://cloud.ibm.com/apidocs/resource-controller/resource-controller?code=go#list-resource-instances
		nextURL, err = core.GetQueryParam(resources.NextURL, "start")
		if err != nil {
			return nil, fmt.Errorf("failed to GetQueryParam on start: %w", err)
		}
		if nextURL == nil {
			options.SetStart("")
		} else {
			options.SetStart(*nextURL)
		}

		moreData = *resources.RowsCount == perPage
	}

	return serviceInstances, nil
}

// ServiceInstanceGUIDToName returns the name of the matching service instance GUID which was passed in.
func (c *Client) ServiceInstanceGUIDToName(ctx context.Context, id string) (string, error) {
	var (
		options   *resourcecontrollerv2.ListResourceInstancesOptions
		resources *resourcecontrollerv2.ResourceInstancesList
		err       error
		perPage   int64 = 10
		moreData        = true
		nextURL   *string
		groupID   = c.BXCli.PowerVSResourceGroup
	)

	// If the user passes in a human readable group id, then we need to convert it to a UUID
	listGroupOptions := c.managementAPI.NewListResourceGroupsOptions()
	listGroupOptions.AccountID = &c.BXCli.User.Account
	groups, _, err := c.managementAPI.ListResourceGroupsWithContext(ctx, listGroupOptions)
	if err != nil {
		return "", fmt.Errorf("failed to list resource groups: %w", err)
	}
	for _, group := range groups.Resources {
		if *group.Name == groupID {
			groupID = *group.ID
		}
	}

	options = c.controllerAPI.NewListResourceInstancesOptions()
	options.SetResourceGroupID(groupID)
	// resource ID for Power Systems Virtual Server in the Global catalog
	options.SetResourceID(powerIAASResourceID)
	options.SetLimit(perPage)

	for moreData {
		resources, _, err = c.controllerAPI.ListResourceInstancesWithContext(ctx, options)
		if err != nil {
			return "", fmt.Errorf("failed to list resource instances: %w", err)
		}

		for _, resource := range resources.Resources {
			var (
				getResourceOptions *resourcecontrollerv2.GetResourceInstanceOptions
				resourceInstance   *resourcecontrollerv2.ResourceInstance
				response           *core.DetailedResponse
			)

			getResourceOptions = c.controllerAPI.NewGetResourceInstanceOptions(*resource.ID)

			resourceInstance, response, err = c.controllerAPI.GetResourceInstance(getResourceOptions)
			if err != nil {
				return "", fmt.Errorf("failed to get instance: %w", err)
			}
			if response != nil && response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusInternalServerError {
				continue
			}

			if resourceInstance.Type != nil && (*resourceInstance.Type == serviceInstanceType || *resourceInstance.Type == compositeInstanceType) {
				if resourceInstance.GUID != nil && *resourceInstance.GUID == id {
					if resourceInstance.Name == nil {
						return "", nil
					}
					return *resourceInstance.Name, nil
				}
			}
		}

		// Based on: https://cloud.ibm.com/apidocs/resource-controller/resource-controller?code=go#list-resource-instances
		nextURL, err = core.GetQueryParam(resources.NextURL, "start")
		if err != nil {
			return "", fmt.Errorf("failed to GetQueryParam on start: %w", err)
		}
		if nextURL == nil {
			options.SetStart("")
		} else {
			options.SetStart(*nextURL)
		}

		moreData = *resources.RowsCount == perPage
	}

	return "", nil
}

// ServiceInstanceNameToGUID returns the name of the matching service instance GUID which was passed in.
func (c *Client) ServiceInstanceNameToGUID(ctx context.Context, name string) (string, error) {
	var (
		options   *resourcecontrollerv2.ListResourceInstancesOptions
		resources *resourcecontrollerv2.ResourceInstancesList
		err       error
		perPage   int64 = 10
		moreData        = true
		nextURL   *string
		groupID   = c.BXCli.PowerVSResourceGroup
	)

	// If the user passes in a human readable group id, then we need to convert it to a UUID
	listGroupOptions := c.managementAPI.NewListResourceGroupsOptions()
	listGroupOptions.AccountID = &c.BXCli.User.Account
	groups, _, err := c.managementAPI.ListResourceGroupsWithContext(ctx, listGroupOptions)
	if err != nil {
		return "", fmt.Errorf("failed to list resource groups: %w", err)
	}
	for _, group := range groups.Resources {
		if *group.Name == groupID {
			groupID = *group.ID
		}
	}

	options = c.controllerAPI.NewListResourceInstancesOptions()
	options.SetResourceGroupID(groupID)
	// resource ID for Power Systems Virtual Server in the Global catalog
	options.SetResourceID(powerIAASResourceID)
	options.SetLimit(perPage)

	for moreData {
		resources, _, err = c.controllerAPI.ListResourceInstancesWithContext(ctx, options)
		if err != nil {
			return "", fmt.Errorf("failed to list resource instances: %w", err)
		}

		for _, resource := range resources.Resources {
			var (
				getResourceOptions *resourcecontrollerv2.GetResourceInstanceOptions
				resourceInstance   *resourcecontrollerv2.ResourceInstance
				response           *core.DetailedResponse
			)

			getResourceOptions = c.controllerAPI.NewGetResourceInstanceOptions(*resource.ID)

			resourceInstance, response, err = c.controllerAPI.GetResourceInstance(getResourceOptions)
			if err != nil {
				return "", fmt.Errorf("failed to get instance: %w", err)
			}
			if response != nil && response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusInternalServerError {
				continue
			}

			if resourceInstance.Type != nil && (*resourceInstance.Type == serviceInstanceType || *resourceInstance.Type == compositeInstanceType) {
				if resourceInstance.Name != nil && *resourceInstance.Name == name {
					if resourceInstance.GUID == nil {
						return "", nil
					}
					return *resourceInstance.GUID, nil
				}
			}
		}

		// Based on: https://cloud.ibm.com/apidocs/resource-controller/resource-controller?code=go#list-resource-instances
		nextURL, err = core.GetQueryParam(resources.NextURL, "start")
		if err != nil {
			return "", fmt.Errorf("failed to GetQueryParam on start: %w", err)
		}
		if nextURL == nil {
			options.SetStart("")
		} else {
			options.SetStart(*nextURL)
		}

		moreData = *resources.RowsCount == perPage
	}

	return "", nil
}

// GetDatacenterCapabilities retrieves the capabilities of the specified datacenter.
func (c *Client) GetDatacenterCapabilities(ctx context.Context, region string) (map[string]bool, error) {
	var err error
	if c.BXCli.PISession == nil {
		err = c.BXCli.NewPISession()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize PISession in GetDatacenterCapabilities: %w", err)
		}
	}
	params := datacenters.NewV1DatacentersGetParamsWithContext(ctx).WithDatacenterRegion(region)
	getOk, err := c.BXCli.PISession.Power.Datacenters.V1DatacentersGet(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get datacenter capabilities: %w", err)
	}
	return getOk.Payload.Capabilities, nil
}

// GetDatacenterSupportedSystems retrieves the capabilities of the specified datacenter.
func (c *Client) GetDatacenterSupportedSystems(ctx context.Context, region string) ([]string, error) {
	var err error
	if c.BXCli.PISession == nil {
		err = c.BXCli.NewPISession()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize PISession in GetDatacenterSupportedSystems: %w", err)
		}
	}
	params := datacenters.NewV1DatacentersGetParamsWithContext(ctx).WithDatacenterRegion(region)
	getOk, err := c.BXCli.PISession.Power.Datacenters.V1DatacentersGet(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get datacenter supported systems: %w", err)
	}
	return getOk.Payload.CapabilitiesDetails.SupportedSystems.General, nil
}

// TransitGatewayNameToID checks to see if the name is an existing transit gateway name.
func (c *Client) TransitGatewayNameToID(ctx context.Context, name string) (string, error) {
	var (
		gateways []transitgatewayapisv1.TransitGateway
		gateway  transitgatewayapisv1.TransitGateway
		err      error
	)

	gateways, err = c.getTransitGateways(ctx)
	if err != nil {
		return "", err
	}
	for _, gateway = range gateways {
		if *gateway.Name == name {
			return *gateway.ID, nil
		}
	}

	return "", nil
}

// TransitGatewayIDValid checks to see if the id is an existing transit gateway id.
func (c *Client) TransitGatewayIDValid(ctx context.Context, id string) error {
	var (
		gateways []transitgatewayapisv1.TransitGateway
		gateway  transitgatewayapisv1.TransitGateway
		err      error
	)

	gateways, err = c.getTransitGateways(ctx)
	if err != nil {
		return err
	}
	for _, gateway = range gateways {
		if *gateway.ID == id {
			return nil
		}
	}

	return fmt.Errorf("transit gateway id (%s) not found", id)
}

// GetAttachedTransitGateway finds an existing Transit Gateway attached to the provided PowerVS cloud instance.
func (c *Client) GetAttachedTransitGateway(ctx context.Context, svcInsID string) (string, error) {
	var (
		gateways []transitgatewayapisv1.TransitGateway
		gateway  transitgatewayapisv1.TransitGateway
		err      error
		conns    []transitgatewayapisv1.TransitConnection
		conn     transitgatewayapisv1.TransitConnection
	)
	gateways, err = c.getTransitGateways(ctx)
	if err != nil {
		return "", err
	}
	for _, gateway = range gateways {
		conns, err = c.getTransitConnections(ctx, *gateway.ID)
		if err != nil {
			return "", err
		}
		for _, conn = range conns {
			if *conn.NetworkType == "power_virtual_server" && strings.Contains(*conn.NetworkID, svcInsID) {
				return *conn.TransitGateway.ID, nil
			}
		}
	}
	return "", nil
}

// GetTGConnectionVPC checks if the VPC subnet is attached to the provided Transit Gateway.
func (c *Client) GetTGConnectionVPC(ctx context.Context, gatewayID string, vpcSubnetID string) (string, error) {
	conns, err := c.getTransitConnections(ctx, gatewayID)
	if err != nil {
		return "", err
	}
	for _, conn := range conns {
		if *conn.NetworkType == "vpc" && strings.Contains(*conn.NetworkID, vpcSubnetID) {
			return *conn.ID, nil
		}
	}
	return "", nil
}

func (c *Client) getTransitGateways(ctx context.Context) ([]transitgatewayapisv1.TransitGateway, error) {
	var (
		listTransitGatewaysOptions *transitgatewayapisv1.ListTransitGatewaysOptions
		gatewayCollection          *transitgatewayapisv1.TransitGatewayCollection
		response                   *core.DetailedResponse
		err                        error
		perPage                    int64 = 32
		moreData                         = true
	)

	listTransitGatewaysOptions = c.transitGatewayAPI.NewListTransitGatewaysOptions()
	listTransitGatewaysOptions.Limit = &perPage

	result := []transitgatewayapisv1.TransitGateway{}

	for moreData {
		// https://github.com/IBM/networking-go-sdk/blob/master/transitgatewayapisv1/transit_gateway_apis_v1.go#L184
		gatewayCollection, response, err = c.transitGatewayAPI.ListTransitGatewaysWithContext(ctx, listTransitGatewaysOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list transit gateways: %w and the respose is: %s", err, response)
		}

		result = append(result, gatewayCollection.TransitGateways...)

		if gatewayCollection.Next != nil {
			listTransitGatewaysOptions.SetStart(*gatewayCollection.Next.Start)
		}

		moreData = gatewayCollection.Next != nil
	}

	return result, nil
}

func (c *Client) getTransitConnections(ctx context.Context, tgID string) ([]transitgatewayapisv1.TransitConnection, error) {
	var (
		listConnectionsOptions *transitgatewayapisv1.ListConnectionsOptions
		connectionCollection   *transitgatewayapisv1.TransitConnectionCollection
		response               *core.DetailedResponse
		err                    error
		perPage                int64 = 32
		moreData                     = true
	)

	listConnectionsOptions = c.transitGatewayAPI.NewListConnectionsOptions()
	listConnectionsOptions.Limit = &perPage

	result := []transitgatewayapisv1.TransitConnection{}

	for moreData {
		connectionCollection, response, err = c.transitGatewayAPI.ListConnectionsWithContext(ctx, listConnectionsOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list transit gateways: %w and the respose is: %s", err, response)
		}

		result = append(result, connectionCollection.Connections...)

		if connectionCollection.Next != nil {
			listConnectionsOptions.SetStart(*connectionCollection.Next.Start)
		}

		moreData = connectionCollection.Next != nil
	}

	return result, nil
}

func (c *Client) ListSecurityGroups(ctx context.Context, vpcID string, region string) ([]vpcv1.SecurityGroup, error) {
	var groupID string
	localContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	groups, err := c.ListResourceGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list resource groups: %w", err)
	}

	err = c.SetVPCServiceURLForRegion(ctx, region)
	if err != nil {
		return nil, err
	}

	for _, group := range groups.Resources {
		if *group.Name == c.BXCli.PowerVSResourceGroup {
			groupID = *group.ID
		}
	}
	listSecurityGroupOptions := c.vpcAPI.NewListSecurityGroupsOptions()
	listSecurityGroupOptions.SetVPCID(vpcID)
	listSecurityGroupOptions.SetResourceGroupID(groupID)
	securityGroupsPager, err := c.vpcAPI.NewSecurityGroupsPager(listSecurityGroupOptions)
	if err != nil {
		return nil, fmt.Errorf("failed creating pager for security group lookup: %w", err)
	}

	securityGroups, err := securityGroupsPager.GetAllWithContext(localContext)
	logrus.Debugf("%v security groups found", len(securityGroups))
	if err != nil {
		return nil, fmt.Errorf("failed collecting all security groups with pager: %w", err)
	}
	for _, sg := range securityGroups {
		logrus.Debugf("SG Name %v found", *sg.Name)
	}
	return securityGroups, nil
}

// ListSecurityGroupRules returns a list of the security group rules.
func (c *Client) ListSecurityGroupRules(ctx context.Context, securityGroupID string) (*vpcv1.SecurityGroupRuleCollection, error) {
	logrus.Debugf("ListSecurityGroupRules: securityGroupID = %s", securityGroupID)

	var (
		vpcOptions  *vpcv1.GetVPCOptions
		vpc         *vpcv1.VPC
		optionsLSGR *vpcv1.ListSecurityGroupRulesOptions
		result      *vpcv1.SecurityGroupRuleCollection
		response    *core.DetailedResponse
		err         error
	)

	vpcOptions = c.vpcAPI.NewGetVPCOptions(securityGroupID)

	vpc, response, err = c.vpcAPI.GetVPC(vpcOptions)
	if err != nil {
		return nil, fmt.Errorf("failure ListSecurityGroupRules GetVPC returns %w, response is %+v", err, response)
	}
	logrus.Debugf("ListSecurityGroupRules: vpc = %+v", vpc)

	optionsLSGR = c.vpcAPI.NewListSecurityGroupRulesOptions(*vpc.DefaultSecurityGroup.ID)

	result, response, err = c.vpcAPI.ListSecurityGroupRulesWithContext(ctx, optionsLSGR)
	if err != nil {
		logrus.Debugf("ListSecurityGroupRules: result = %+v, response = %+v, err = %v", result, response, err)
	}

	return result, err
}

// AddSecurityGroupRule adds a security group rule to an existing security group.
func (c *Client) AddSecurityGroupRule(ctx context.Context, vpcID string, sgID string, rule *vpcv1.SecurityGroupRulePrototype) error {
	logrus.Debugf("AddSecurityGroupRule: vpcID = %s, rule = %+v", vpcID, *rule)

	var (
		vpcOptions  *vpcv1.GetVPCOptions
		vpc         *vpcv1.VPC
		optionsCSGR *vpcv1.CreateSecurityGroupRuleOptions
		result      vpcv1.SecurityGroupRuleIntf
		response    *core.DetailedResponse
		err         error
	)

	vpcOptions = c.vpcAPI.NewGetVPCOptions(vpcID)

	vpc, response, err = c.vpcAPI.GetVPC(vpcOptions)
	if err != nil {
		return fmt.Errorf("failure AddSecurityGroupRule GetVPC returns %w, response is %+v", err, response)
	}
	logrus.Debugf("AddSecurityGroupRule: vpc = %+v", vpc)

	optionsCSGR = &vpcv1.CreateSecurityGroupRuleOptions{}
	if sgID == "" {
		optionsCSGR.SetSecurityGroupID(*vpc.DefaultSecurityGroup.ID)
	} else {
		optionsCSGR.SetSecurityGroupID(sgID)
	}
	optionsCSGR.SetSecurityGroupRulePrototype(rule)

	result, response, err = c.vpcAPI.CreateSecurityGroupRuleWithContext(ctx, optionsCSGR)
	if err != nil {
		logrus.Debugf("AddSecurityGroupRule: result = %+v, response = %+v, err = %v", result, response, err)
	}

	return err
}

// CreateSSHKey creates a SSH key in the PowerVS Workspace for the workers to use.
func (c *Client) CreateSSHKey(ctx context.Context, serviceInstance string, zone string, sshKeyName string, sshKey string) error {
	logrus.Debugf("CreateSSHKey: serviceInstance = %s, sshKeyName = %s sshKey = %s", serviceInstance, sshKeyName, sshKey)

	var (
		user          *User
		authenticator core.Authenticator
		options       *ibmpisession.IBMPIOptions
		piSession     *ibmpisession.IBMPISession
		keyClient     *instance.IBMPIKeyClient
		sshKeyInput   *models.SSHKey
		err           error
	)

	user, err = FetchUserDetails(c.APIKey)
	if err != nil {
		return fmt.Errorf("createSSHKey: failed to fetch user details %w", err)
	}

	authenticator = &core.IamAuthenticator{
		ApiKey: c.APIKey,
	}
	err = authenticator.Validate()
	if err != nil {
		return fmt.Errorf("createSSHKey: authenticator failed validate %w", err)
	}

	options = &ibmpisession.IBMPIOptions{
		Authenticator: authenticator,
		Debug:         false,
		UserAccount:   user.Account,
		Zone:          zone,
	}

	piSession, err = ibmpisession.NewIBMPISession(options)
	if err != nil {
		return fmt.Errorf("createSSHKey: ibmpisession.New: %w", err)
	}
	if piSession == nil {
		return fmt.Errorf("createSSHKey: piSession is nil")
	}
	logrus.Debugf("CreateSSHKey: piSession = %+v", piSession)

	keyClient = instance.NewIBMPIKeyClient(ctx, piSession, serviceInstance)
	logrus.Debugf("CreateSSHKey: keyClient = %+v", keyClient)

	sshKeyInput = &models.SSHKey{
		Name:   ptr.To(sshKeyName),
		SSHKey: ptr.To(sshKey),
	}
	logrus.Debugf("CreateSSHKey: sshKeyInput = %+v", sshKeyInput)

	_, err = keyClient.Create(sshKeyInput)
	if err != nil {
		return fmt.Errorf("createSSHKey: failed to create the ssh key %w", err)
	}

	return nil
}

// AddIPToLoadBalancerPool adds a server to a load balancer pool for the specified port.
// @TODO Remove once https://github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/issues/1679 is fixed.
func (c *Client) AddIPToLoadBalancerPool(ctx context.Context, lbID string, poolName string, port int64, ip string) error {
	var (
		glbOptions    *vpcv1.GetLoadBalancerOptions
		llbpOptions   *vpcv1.ListLoadBalancerPoolsOptions
		llbpmOptions  *vpcv1.ListLoadBalancerPoolMembersOptions
		clbpmOptions  *vpcv1.CreateLoadBalancerPoolMemberOptions
		lb            *vpcv1.LoadBalancer
		lbPools       *vpcv1.LoadBalancerPoolCollection
		lbPool        vpcv1.LoadBalancerPool
		lbPoolMembers *vpcv1.LoadBalancerPoolMemberCollection
		lbpmtp        *vpcv1.LoadBalancerPoolMemberTargetPrototypeIP
		lbpm          *vpcv1.LoadBalancerPoolMember
		response      *core.DetailedResponse
		err           error
	)

	// Make sure the load balancer exists
	glbOptions = c.vpcAPI.NewGetLoadBalancerOptions(lbID)

	lb, response, err = c.vpcAPI.GetLoadBalancerWithContext(ctx, glbOptions)
	if err != nil {
		return fmt.Errorf("failed to get load balancer and the response = %+v, err = %w", response, err)
	}
	logrus.Debugf("AddIPToLoadBalancerPool: GLBWC lb = %+v", lb)

	// Query the existing load balancer pools
	llbpOptions = c.vpcAPI.NewListLoadBalancerPoolsOptions(lbID)

	lbPools, response, err = c.vpcAPI.ListLoadBalancerPoolsWithContext(ctx, llbpOptions)
	if err != nil {
		return fmt.Errorf("failed to list the load balancer pools and the response = %+v, err = %w",
			response,
			err)
	}

	// Find the pool with the specified name
	for _, pool := range lbPools.Pools {
		logrus.Debugf("AddIPToLoadBalancerPool: pool.ID = %v", *pool.ID)
		logrus.Debugf("AddIPToLoadBalancerPool: pool.Name = %v", *pool.Name)

		if *pool.Name == poolName {
			lbPool = pool
			break
		}
	}
	if lbPool.ID == nil {
		return fmt.Errorf("could not find loadbalancer pool with name %s", poolName)
	}

	// Query the load balancer pool members
	llbpmOptions = c.vpcAPI.NewListLoadBalancerPoolMembersOptions(lbID, *lbPool.ID)

	lbPoolMembers, response, err = c.vpcAPI.ListLoadBalancerPoolMembersWithContext(ctx, llbpmOptions)
	if err != nil {
		return fmt.Errorf("failed to list load balancer pool members and the response = %+v, err = %w",
			response,
			err)
	}

	// See if a member already exists with that IP
	for _, poolMember := range lbPoolMembers.Members {
		logrus.Debugf("AddIPToLoadBalancerPool: poolMember.ID = %s", *poolMember.ID)

		switch pmt := poolMember.Target.(type) {
		case *vpcv1.LoadBalancerPoolMemberTarget:
			logrus.Debugf("AddIPToLoadBalancerPool: pmt.Address = %+v", *pmt.Address)
			if ip == *pmt.Address {
				logrus.Debugf("AddIPToLoadBalancerPool: found %s", ip)
				return nil
			}
		case *vpcv1.LoadBalancerPoolMemberTargetIP:
			logrus.Debugf("AddIPToLoadBalancerPool: pmt.Address = %+v", *pmt.Address)
			if ip == *pmt.Address {
				logrus.Debugf("AddIPToLoadBalancerPool: found %s", ip)
				return nil
			}
		case *vpcv1.LoadBalancerPoolMemberTargetInstanceReference:
			// No IP address, ignore
		default:
			logrus.Debugf("AddIPToLoadBalancerPool: unhandled type %T", poolMember.Target)
		}
	}

	// Create a new member
	lbpmtp, err = c.vpcAPI.NewLoadBalancerPoolMemberTargetPrototypeIP(ip)
	if err != nil {
		return fmt.Errorf("could not create a new load balancer pool member, err = %w", err)
	}
	logrus.Debugf("AddIPToLoadBalancerPool: lbpmtp = %+v", *lbpmtp)

	// Add that member to the pool
	clbpmOptions = c.vpcAPI.NewCreateLoadBalancerPoolMemberOptions(lbID, *lbPool.ID, port, lbpmtp)
	logrus.Debugf("AddIPToLoadBalancerPool: clbpmOptions = %+v", clbpmOptions)

	return wait.PollUntilContextCancel(ctx,
		time.Second*10,
		false,
		func(ctx context.Context) (bool, error) {
			lbpm, response, err = c.vpcAPI.CreateLoadBalancerPoolMemberWithContext(ctx, clbpmOptions)
			if err != nil {
				logrus.Debugf("AddIPToLoadBalancerPool: could not add the load balancer pool member yet, err = %v", err)
				return false, nil
			}

			logrus.Debugf("AddIPToLoadBalancerPool: CLBPMWC lbpm = %+v", lbpm)

			return true, nil
		})
}

// CreateVirtualPrivateEndpointGateway creates a VPE gateway with given target resource type and CRN.
func (c *Client) CreateVirtualPrivateEndpointGateway(ctx context.Context, name string, vpcID string, subnetID string, rgID string, targetCRN string) (*vpcv1.EndpointGateway, error) {
	var (
		resp   *core.DetailedResponse
		err    error
		ok     bool
		egs    *vpcv1.EndpointGatewayCollection
		egRef  *vpcv1.EndpointGatewayTarget
		idIntf *vpcv1.VPCIdentityByID
		target *vpcv1.EndpointGatewayTargetPrototypeEndpointGatewayTargetResourceTypeProviderCloudServicePrototype
		rgIntf *vpcv1.ResourceGroupIdentityByID
		ipIntf *vpcv1.EndpointGatewayReservedIPReservedIPIdentityByID
	)

	listOpts := c.vpcAPI.NewListEndpointGatewaysOptions()
	listOpts.SetVPCID(vpcID)
	egs, _, err = c.vpcAPI.ListEndpointGateways(listOpts)
	if err != nil {
		return nil, err
	}

	for _, eg := range egs.EndpointGateways {
		egRef, ok = eg.Target.(*vpcv1.EndpointGatewayTarget)
		if !ok {
			return nil, fmt.Errorf("invalid target inside returned EndpointGateway: %v", eg.Target)
		}
		if *egRef.CRN == targetCRN {
			return &eg, nil
		}
	}

	target, err = c.vpcAPI.NewEndpointGatewayTargetPrototypeEndpointGatewayTargetResourceTypeProviderCloudServicePrototype(targetCRN, vpcv1.EndpointGatewayTargetPrototypeResourceTypeProviderCloudServiceConst)
	if err != nil {
		return nil, err
	}
	idIntf, err = c.vpcAPI.NewVPCIdentityByID(vpcID)
	if err != nil {
		return nil, err
	}
	createOpts := c.vpcAPI.NewCreateEndpointGatewayOptions(target, idIntf)
	createOpts.SetName(name)
	createOpts.SetAllowDnsResolutionBinding(true)
	rgIntf, err = c.vpcAPI.NewResourceGroupIdentityByID(rgID)
	if err != nil {
		return nil, err
	}
	createOpts.SetResourceGroup(rgIntf)
	ipName := fmt.Sprintf("%s-ip", name)
	createIPOpts := c.vpcAPI.NewCreateSubnetReservedIPOptions(subnetID)
	createIPOpts.SetName(ipName)
	createIPOpts.SetSubnetID(subnetID)
	reservedIP, _, err := c.vpcAPI.CreateSubnetReservedIPWithContext(ctx, createIPOpts)
	if err != nil {
		return nil, err
	}
	ipIntf, err = c.vpcAPI.NewEndpointGatewayReservedIPReservedIPIdentityByID(*reservedIP.ID)
	if err != nil {
		return nil, err
	}
	ips := []vpcv1.EndpointGatewayReservedIPIntf{ipIntf}
	createOpts.SetIps(ips)

	eg, resp, err := c.vpcAPI.CreateEndpointGatewayWithContext(ctx, createOpts)
	if err != nil {
		logrus.Debugf("CreateEndpointGatewayWithContext returned %v", resp)
		return nil, err
	}
	return eg, nil
}

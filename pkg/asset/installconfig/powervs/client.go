package powervs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM-Cloud/power-go-client/power/client/datacenters"
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

	"github.com/openshift/installer/pkg/types"
)

//go:generate mockgen -source=./client.go -destination=./mock/powervsclient_generated.go -package=mock

// API represents the calls made to the API.
type API interface {
	GetDNSRecordsByName(ctx context.Context, crnstr string, zoneID string, recordName string, publish types.PublishingStrategy) ([]DNSRecordResponse, error)
	GetDNSZoneIDByName(ctx context.Context, name string, publish types.PublishingStrategy) (string, error)
	GetDNSZones(ctx context.Context, publish types.PublishingStrategy) ([]DNSZoneResponse, error)
	GetDNSInstancePermittedNetworks(ctx context.Context, dnsID string, dnsZone string) ([]string, error)
	GetVPCByName(ctx context.Context, vpcName string) (*vpcv1.VPC, error)
	GetPublicGatewayByVPC(ctx context.Context, vpcName string) (*vpcv1.PublicGateway, error)
	GetSubnetByName(ctx context.Context, subnetName string, region string) (*vpcv1.Subnet, error)
	GetAuthenticatorAPIKeyDetails(ctx context.Context) (*iamidentityv1.APIKey, error)
	GetAPIKey() string
	SetVPCServiceURLForRegion(ctx context.Context, region string) error
	GetVPCs(ctx context.Context, region string) ([]vpcv1.VPC, error)
	ListResourceGroups(ctx context.Context) (*resourcemanagerv2.ResourceGroupList, error)
	ListServiceInstances(ctx context.Context) ([]string, error)
	ServiceInstanceGUIDToName(ctx context.Context, id string) (string, error)
	GetDatacenterCapabilities(ctx context.Context, region string) (map[string]bool, error)
	GetAttachedTransitGateway(ctx context.Context, svcInsID string) (string, error)
	GetTGConnectionVPC(ctx context.Context, gatewayID string, vpcSubnetID string) (string, error)
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
	cisServiceID = "75874a60-cb12-11e7-948e-37ac098eb1b9"
	dnsServiceID = "b4ed8a30-936f-11e9-b289-1d079699cbe5"
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
				if *zone.State == "ACTIVE" {
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

// GetDNSInstancePermittedNetworks gets the permitted VPC networks for a DNS Services instance
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

// GetVPCByName gets a VPC by its name.
func (c *Client) GetVPCByName(ctx context.Context, vpcName string) (*vpcv1.VPC, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	listRegionsOptions := c.vpcAPI.NewListRegionsOptions()
	listRegionsResponse, _, err := c.vpcAPI.ListRegionsWithContext(ctx, listRegionsOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list vpc regions: %w", err)
	}

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
				if *vpc.Name == vpcName {
					return &vpc, nil
				}
			}
		}
	}

	return nil, errors.New("failed to find VPC")
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

			if resourceInstance.Type != nil && *resourceInstance.Type == "service_instance" {
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

			if resourceInstance.Type != nil && *resourceInstance.Type == "service_instance" {
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

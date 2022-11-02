package powervs

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/dnsrecordsv1"
	"github.com/IBM/networking-go-sdk/dnssvcsv1"
	"github.com/IBM/networking-go-sdk/dnszonesv1"
	"github.com/IBM/networking-go-sdk/resourcerecordsv1"
	"github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
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
}

// Client makes calls to the PowerVS API.
type Client struct {
	APIKey         string
	managementAPI  *resourcemanagerv2.ResourceManagerV2
	controllerAPI  *resourcecontrollerv2.ResourceControllerV2
	vpcAPI         *vpcv1.VpcV1
	dnsServicesAPI *dnssvcsv1.DnsSvcsV1
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
	bxCli, err := NewBxClient()
	if err != nil {
		return nil, err
	}

	client := &Client{
		APIKey: bxCli.APIKey,
	}

	if err := client.loadSDKServices(); err != nil {
		return nil, errors.Wrap(err, "failed to load IBM SDK services")
	}

	return client, nil
}

func (c *Client) loadSDKServices() error {
	servicesToLoad := []func() error{
		c.loadResourceManagementAPI,
		c.loadResourceControllerAPI,
		c.loadVPCV1API,
		c.loadDNSServicesAPI,
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
			return nil, errors.Wrap(err, "could not retrieve DNS records")
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
			return nil, errors.Wrap(err, "Failed to parse DNSInstanceCRN")
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
			return nil, errors.Wrap(err, "could not retrieve DNS records")
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
		return nil, errors.Wrap(err, "failed to get cis instance")
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
				return nil, errors.Wrap(err, "failed to list DNS zones")
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
				return nil, errors.Wrap(err, "failed to list DNS zones")
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
		return nil, errors.Wrap(err, "failed to list vpc regions")
	}

	for _, region := range listRegionsResponse.Regions {
		err := c.vpcAPI.SetServiceURL(fmt.Sprintf("%s/v1", *region.Endpoint))
		if err != nil {
			return nil, errors.Wrap(err, "failed to set vpc api service url")
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
		return nil, errors.Wrap(err, "failed to get VPC")
	}

	vpcCRN, err := crn.Parse(*vpc.CRN)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse VPC CRN")
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

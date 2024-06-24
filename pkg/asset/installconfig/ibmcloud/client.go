package ibmcloud

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/dnsrecordsv1"
	"github.com/IBM/networking-go-sdk/dnssvcsv1"
	"github.com/IBM/networking-go-sdk/dnszonesv1"
	"github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud/responses"
	"github.com/openshift/installer/pkg/types"
)

//go:generate mockgen -source=./client.go -destination=./mock/ibmcloudclient_generated.go -package=mock

// API represents the calls made to the API.
type API interface {
	GetAPIKey() string
	GetAuthenticatorAPIKeyDetails(ctx context.Context) (*iamidentityv1.APIKey, error)
	GetCISInstance(ctx context.Context, crnstr string) (*resourcecontrollerv2.ResourceInstance, error)
	GetDNSInstance(ctx context.Context, crnstr string) (*resourcecontrollerv2.ResourceInstance, error)
	GetDNSInstancePermittedNetworks(ctx context.Context, dnsID string, dnsZone string) ([]string, error)
	GetDedicatedHostByName(ctx context.Context, name string, region string) (*vpcv1.DedicatedHost, error)
	GetDedicatedHostProfiles(ctx context.Context, region string) ([]vpcv1.DedicatedHostProfile, error)
	GetDNSRecordsByName(ctx context.Context, crnstr string, zoneID string, recordName string) ([]dnsrecordsv1.DnsrecordDetails, error)
	GetDNSZoneIDByName(ctx context.Context, name string, publish types.PublishingStrategy) (string, error)
	GetDNSZones(ctx context.Context, publish types.PublishingStrategy) ([]responses.DNSZoneResponse, error)
	GetEncryptionKey(ctx context.Context, keyCRN string) (*responses.EncryptionKeyResponse, error)
	GetResourceGroups(ctx context.Context) ([]resourcemanagerv2.ResourceGroup, error)
	GetResourceGroup(ctx context.Context, nameOrID string) (*resourcemanagerv2.ResourceGroup, error)
	GetSubnet(ctx context.Context, subnetID string) (*vpcv1.Subnet, error)
	GetSubnetByName(ctx context.Context, subnetName string, region string) (*vpcv1.Subnet, error)
	GetVSIProfiles(ctx context.Context) ([]vpcv1.InstanceProfile, error)
	GetVPC(ctx context.Context, vpcID string) (*vpcv1.VPC, error)
	GetVPCs(ctx context.Context, region string) ([]vpcv1.VPC, error)
	GetVPCByName(ctx context.Context, vpcName string) (*vpcv1.VPC, error)
	GetVPCZonesForRegion(ctx context.Context, region string) ([]string, error)
	SetVPCServiceURLForRegion(ctx context.Context, region string) error
}

// Client makes calls to the IBM Cloud API.
type Client struct {
	apiKey         string
	managementAPI  *resourcemanagerv2.ResourceManagerV2
	controllerAPI  *resourcecontrollerv2.ResourceControllerV2
	vpcAPI         *vpcv1.VpcV1
	dnsServicesAPI *dnssvcsv1.DnsSvcsV1
}

// InstanceType is the IBM Cloud network services type being used
type InstanceType string

const (
	// CISInstanceType is a Cloud Internet Services InstanceType
	CISInstanceType InstanceType = "CIS"
	// DNSInstanceType is a DNS Services InstanceType
	DNSInstanceType InstanceType = "DNS"

	// cisServiceID is the Cloud Internet Services' catalog service ID.
	cisServiceID = "75874a60-cb12-11e7-948e-37ac098eb1b9"
	// dnsServiceID is the DNS Services' catalog service ID.
	dnsServiceID = "b4ed8a30-936f-11e9-b289-1d079699cbe5"
)

// VPCResourceNotFoundError represents an error for a VPC resoruce that is not found.
type VPCResourceNotFoundError struct{}

// Error returns the error message for the VPCResourceNotFoundError error type.
func (e *VPCResourceNotFoundError) Error() string {
	return "Not Found"
}

// NewClient initializes a client with a session.
func NewClient() (*Client, error) {
	apiKey := os.Getenv("IC_API_KEY")

	client := &Client{
		apiKey: apiKey,
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

// GetAPIKey gets the API Key.
func (c *Client) GetAPIKey() string {
	return c.apiKey
}

// GetAuthenticatorAPIKeyDetails gets detailed information on the API key used
// for authentication to the IBM Cloud APIs
func (c *Client) GetAuthenticatorAPIKeyDetails(ctx context.Context) (*iamidentityv1.APIKey, error) {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey())
	if err != nil {
		return nil, err
	}
	iamIdentityService, err := iamidentityv1.NewIamIdentityV1(&iamidentityv1.IamIdentityV1Options{
		Authenticator: authenticator,
	})
	if err != nil {
		return nil, err
	}

	options := iamIdentityService.NewGetAPIKeysDetailsOptions()
	options.SetIamAPIKey(c.GetAPIKey())
	details, _, err := iamIdentityService.GetAPIKeysDetailsWithContext(ctx, options)
	if err != nil {
		return nil, err
	}
	return details, nil
}

// getInstance gets a specific DNS or CIS instance by its CRN.
func (c *Client) getInstance(ctx context.Context, crnstr string, iType InstanceType) (*resourcecontrollerv2.ResourceInstance, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	options := c.controllerAPI.NewGetResourceInstanceOptions(crnstr)
	resourceInstance, _, err := c.controllerAPI.GetResourceInstance(options)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s instances", iType)
	}

	return resourceInstance, nil
}

// GetCISInstance gets a specific Cloud Internet Services by its CRN.
func (c *Client) GetCISInstance(ctx context.Context, crnstr string) (*resourcecontrollerv2.ResourceInstance, error) {
	return c.getInstance(ctx, crnstr, CISInstanceType)
}

// GetDNSInstance gets a specific DNS Services instance by its CRN.
func (c *Client) GetDNSInstance(ctx context.Context, crnstr string) (*resourcecontrollerv2.ResourceInstance, error) {
	return c.getInstance(ctx, crnstr, DNSInstanceType)
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

// GetDedicatedHostByName gets dedicated host by name.
func (c *Client) GetDedicatedHostByName(ctx context.Context, name string, region string) (*vpcv1.DedicatedHost, error) {
	err := c.SetVPCServiceURLForRegion(ctx, region)
	if err != nil {
		return nil, err
	}

	options := c.vpcAPI.NewListDedicatedHostsOptions()
	dhosts, _, err := c.vpcAPI.ListDedicatedHostsWithContext(ctx, options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list dedicated hosts")
	}

	for _, dhost := range dhosts.DedicatedHosts {
		if *dhost.Name == name {
			return &dhost, nil
		}
	}

	return nil, fmt.Errorf("dedicated host %q not found", name)
}

// GetDedicatedHostProfiles gets a list of profiles supported in a region.
func (c *Client) GetDedicatedHostProfiles(ctx context.Context, region string) ([]vpcv1.DedicatedHostProfile, error) {
	err := c.SetVPCServiceURLForRegion(ctx, region)
	if err != nil {
		return nil, err
	}

	profilesOptions := c.vpcAPI.NewListDedicatedHostProfilesOptions()
	profiles, _, err := c.vpcAPI.ListDedicatedHostProfilesWithContext(ctx, profilesOptions)
	if err != nil {
		return nil, err
	}

	return profiles.Profiles, nil
}

// GetDNSRecordsByName gets DNS records in specific Cloud Internet Services instance
// by its CRN, zone ID, and DNS record name.
func (c *Client) GetDNSRecordsByName(ctx context.Context, crnstr string, zoneID string, recordName string) ([]dnsrecordsv1.DnsrecordDetails, error) {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey())
	if err != nil {
		return nil, err
	}
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

	return records.Result, nil
}

// GetDNSZoneIDByName gets the DNS (Internal) or CIS zone ID from its domain name.
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

// GetDNSZones returns all of the active DNS zones managed by DNS or CIS.
func (c *Client) GetDNSZones(ctx context.Context, publish types.PublishingStrategy) ([]responses.DNSZoneResponse, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	if publish == types.InternalPublishingStrategy {
		return c.getDNSDNSZones(ctx)
	}
	return c.getCISDNSZones(ctx)
}

func (c *Client) getDNSDNSZones(ctx context.Context) ([]responses.DNSZoneResponse, error) {
	options := c.controllerAPI.NewListResourceInstancesOptions()
	options.SetResourceID(dnsServiceID)

	listResourceInstancesResponse, _, err := c.controllerAPI.ListResourceInstances(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get dns instance")
	}

	var allZones []responses.DNSZoneResponse
	for _, instance := range listResourceInstancesResponse.Resources {
		authenticator, err := NewIamAuthenticator(c.GetAPIKey())
		if err != nil {
			return nil, err
		}
		dnsZoneService, err := dnszonesv1.NewDnsZonesV1(&dnszonesv1.DnsZonesV1Options{
			Authenticator: authenticator,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to list DNS zones")
		}

		options := dnsZoneService.NewListDnszonesOptions(*instance.GUID)
		result, _, err := dnsZoneService.ListDnszones(options)
		if result == nil {
			return nil, err
		}

		for _, zone := range result.Dnszones {
			stateLower := strings.ToLower(*zone.State)
			// DNS Zones can be 'pending_network_add' (without a permitted network, added during TF)
			if stateLower == dnszonesv1.Dnszone_State_Active || stateLower == dnszonesv1.Dnszone_State_PendingNetworkAdd {
				zoneStruct := responses.DNSZoneResponse{
					Name:            *zone.Name,
					ID:              *zone.ID,
					InstanceID:      *instance.GUID,
					InstanceCRN:     *instance.CRN,
					InstanceName:    *instance.Name,
					ResourceGroupID: *instance.ResourceGroupID,
				}
				allZones = append(allZones, zoneStruct)
			}
		}
	}

	return allZones, nil
}

func (c *Client) getCISDNSZones(ctx context.Context) ([]responses.DNSZoneResponse, error) {
	options := c.controllerAPI.NewListResourceInstancesOptions()
	options.SetResourceID(cisServiceID)

	listResourceInstancesResponse, _, err := c.controllerAPI.ListResourceInstances(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cis instance")
	}

	var allZones []responses.DNSZoneResponse
	for _, instance := range listResourceInstancesResponse.Resources {
		authenticator, err := NewIamAuthenticator(c.GetAPIKey())
		if err != nil {
			return nil, err
		}
		crnstr := instance.CRN
		zonesService, err := zonesv1.NewZonesV1(&zonesv1.ZonesV1Options{
			Authenticator: authenticator,
			Crn:           crnstr,
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
				zoneStruct := responses.DNSZoneResponse{
					Name:            *zone.Name,
					ID:              *zone.ID,
					InstanceID:      *instance.GUID,
					InstanceCRN:     *instance.CRN,
					InstanceName:    *instance.Name,
					ResourceGroupID: *instance.ResourceGroupID,
				}
				allZones = append(allZones, zoneStruct)
			}
		}
	}

	return allZones, nil
}

// GetEncryptionKey gets data for an encryption key
func (c *Client) GetEncryptionKey(ctx context.Context, keyCRN string) (*responses.EncryptionKeyResponse, error) {
	// TODO: IBM: Call KMS / Hyperprotect Crpyto APIs.
	return &responses.EncryptionKeyResponse{}, nil
}

// GetResourceGroup gets a resource group by its name or ID.
func (c *Client) GetResourceGroup(ctx context.Context, nameOrID string) (*resourcemanagerv2.ResourceGroup, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	groups, err := c.GetResourceGroups(ctx)
	if err != nil {
		return nil, err
	}

	for idx, rg := range groups {
		if *rg.ID == nameOrID || *rg.Name == nameOrID {
			return &groups[idx], nil
		}
	}
	return nil, fmt.Errorf("resource group %q not found", nameOrID)
}

// GetResourceGroups gets the list of resource groups.
func (c *Client) GetResourceGroups(ctx context.Context) ([]resourcemanagerv2.ResourceGroup, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	apikey, err := c.GetAuthenticatorAPIKeyDetails(ctx)
	if err != nil {
		return nil, err
	}

	options := c.managementAPI.NewListResourceGroupsOptions()
	options.SetAccountID(*apikey.AccountID)
	listResourceGroupsResponse, _, err := c.managementAPI.ListResourceGroupsWithContext(ctx, options)
	if err != nil {
		return nil, err
	}
	return listResourceGroupsResponse.Resources, nil
}

// GetSubnet gets a subnet by its ID.
func (c *Client) GetSubnet(ctx context.Context, subnetID string) (*vpcv1.Subnet, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	subnet, detailedResponse, err := c.vpcAPI.GetSubnet(&vpcv1.GetSubnetOptions{ID: &subnetID})
	if detailedResponse.GetStatusCode() == http.StatusNotFound {
		return nil, &VPCResourceNotFoundError{}
	}
	return subnet, err
}

// GetSubnetByName gets a subnet by its Name.
func (c *Client) GetSubnetByName(ctx context.Context, subnetName string, region string) (*vpcv1.Subnet, error) {
	localContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	err := c.SetVPCServiceURLForRegion(ctx, region)
	if err != nil {
		return nil, err
	}

	listSubnetsOptions := c.vpcAPI.NewListSubnetsOptions()
	subnetsPager, err := c.vpcAPI.NewSubnetsPager(listSubnetsOptions)
	if err != nil {
		return nil, err
	}

	subnets, err := subnetsPager.GetAllWithContext(localContext)
	if err != nil {
		return nil, err
	}
	for _, subnet := range subnets {
		if subnetName == *subnet.Name {
			return &subnet, nil
		}
	}
	return nil, &VPCResourceNotFoundError{}
}

// GetVSIProfiles gets a list of all VSI profiles.
func (c *Client) GetVSIProfiles(ctx context.Context) ([]vpcv1.InstanceProfile, error) {
	listInstanceProfilesOptions := c.vpcAPI.NewListInstanceProfilesOptions()
	profiles, _, err := c.vpcAPI.ListInstanceProfilesWithContext(ctx, listInstanceProfilesOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list vpc vsi profiles")
	}
	return profiles.Profiles, nil
}

// GetVPC gets a VPC by its ID.
func (c *Client) GetVPC(ctx context.Context, vpcID string) (*vpcv1.VPC, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	regions, err := c.getVPCRegions(ctx)
	if err != nil {
		return nil, err
	}

	for _, region := range regions {
		err := c.vpcAPI.SetServiceURL(fmt.Sprintf("%s/v1", *region.Endpoint))
		if err != nil {
			return nil, errors.Wrap(err, "failed to set vpc api service url")
		}

		if vpc, detailedResponse, err := c.vpcAPI.GetVPC(c.vpcAPI.NewGetVPCOptions(vpcID)); err != nil {
			if detailedResponse.GetStatusCode() != http.StatusNotFound {
				return nil, err
			}
		} else if vpc != nil {
			return vpc, nil
		}
	}

	return nil, &VPCResourceNotFoundError{}
}

// GetVPCs gets all VPCs in a region
func (c *Client) GetVPCs(ctx context.Context, region string) ([]vpcv1.VPC, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	err := c.SetVPCServiceURLForRegion(ctx, region)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set vpc api service url")
	}

	allVPCs := []vpcv1.VPC{}
	if vpcs, detailedResponse, err := c.vpcAPI.ListVpcs(c.vpcAPI.NewListVpcsOptions()); err != nil {
		if detailedResponse.GetStatusCode() != http.StatusNotFound {
			return nil, err
		}
	} else if vpcs != nil {
		allVPCs = append(allVPCs, vpcs.Vpcs...)
	}
	return allVPCs, nil
}

// GetVPCByName gets a VPC by its name.
func (c *Client) GetVPCByName(ctx context.Context, vpcName string) (*vpcv1.VPC, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	regions, err := c.getVPCRegions(ctx)
	if err != nil {
		return nil, err
	}

	for _, region := range regions {
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

	return nil, &VPCResourceNotFoundError{}
}

// GetVPCZonesForRegion gets the supported zones for a VPC region.
func (c *Client) GetVPCZonesForRegion(ctx context.Context, region string) ([]string, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	regionZonesOptions := c.vpcAPI.NewListRegionZonesOptions(region)
	zones, _, err := c.vpcAPI.ListRegionZonesWithContext(ctx, regionZonesOptions)
	if err != nil {
		return nil, err
	}

	response := make([]string, len(zones.Zones))
	for idx, zone := range zones.Zones {
		response[idx] = *zone.Name
	}
	return response, err
}

func (c *Client) getVPCRegions(ctx context.Context) ([]vpcv1.Region, error) {
	listRegionsOptions := c.vpcAPI.NewListRegionsOptions()
	listRegionsResponse, _, err := c.vpcAPI.ListRegionsWithContext(ctx, listRegionsOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list vpc regions")
	}

	return listRegionsResponse.Regions, nil
}

func (c *Client) loadResourceManagementAPI() error {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey())
	if err != nil {
		return err
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
	authenticator, err := NewIamAuthenticator(c.GetAPIKey())
	if err != nil {
		return err
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
	authenticator, err := NewIamAuthenticator(c.GetAPIKey())
	if err != nil {
		return err
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
	authenticator, err := NewIamAuthenticator(c.GetAPIKey())
	if err != nil {
		return err
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

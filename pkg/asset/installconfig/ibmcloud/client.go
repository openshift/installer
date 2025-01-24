package ibmcloud

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	ibmcrn "github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM/go-sdk-core/v5/core"
	kpclient "github.com/IBM/keyprotect-go-client"
	"github.com/IBM/networking-go-sdk/dnsrecordsv1"
	"github.com/IBM/networking-go-sdk/dnssvcsv1"
	"github.com/IBM/networking-go-sdk/dnszonesv1"
	"github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/globalcatalogv1"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud/responses"
	"github.com/openshift/installer/pkg/types"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
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
	GetIBMCloudRegions(ctx context.Context) (map[string]string, error)
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
	apiKey        string
	catalogAPI    *globalcatalogv1.GlobalCatalogV1
	managementAPI *resourcemanagerv2.ResourceManagerV2
	controllerAPI *resourcecontrollerv2.ResourceControllerV2
	vpcAPI        *vpcv1.VpcV1

	// A set of overriding endpoints for IBM Cloud Services
	serviceEndpoints []configv1.IBMCloudServiceEndpoint
	// Cache endpoint override for IBM Cloud IAM, if one was provided in serviceEndpoints
	iamServiceEndpointOverride string
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

	// hyperProtectCRNServiceName is the service name within the IBM Cloud CRN for the Hyper Protect service.
	hyperProtectCRNServiceName = "hs-crypto"
	// keyProtectCRNServiceName is the service name within the IBM Cloud CRN for the Key Protect service.
	keyProtectCRNServiceName = "kms"

	// hyperProtectDefaultURLTemplate is the default URL endpoint template, with region substitution, for IBM Cloud Hyper Protect service.
	hyperProtectDefaultURLTemplate = "https://api.%s.hs-crypto.cloud.ibm.com"
	// iamTokenDefaultURL is the default URL endpoint for IBM Cloud IAM token service.
	iamTokenDefaultURL = "https://iam.cloud.ibm.com/identity/token" // #nosec G101 // this is the URL for IBM Cloud IAM tokens, not a credential
	// iamTokenPath is the URL path, to add to override IAM endpoints, for the IBM Cloud IAM token service.
	iamTokenPath = "identity/token" // #nosec G101 // this is the URI path for IBM Cloud IAM tokens, not a credential
	// keyProtectDefaultURLTemplate is the default URL endpoint template, with region substitution, for IBM Cloud Key Protect service.
	keyProtectDefaultURLTemplate = "https://%s.kms.cloud.ibm.com"
)

// VPCResourceNotFoundError represents an error for a VPC resoruce that is not found.
type VPCResourceNotFoundError struct{}

// Error returns the error message for the VPCResourceNotFoundError error type.
func (e *VPCResourceNotFoundError) Error() string {
	return "Not Found"
}

// NewClient initializes a client with any provided endpoint overrides.
func NewClient(endpoints []configv1.IBMCloudServiceEndpoint) (*Client, error) {
	apiKey := os.Getenv("IC_API_KEY")

	client := &Client{
		apiKey:           apiKey,
		serviceEndpoints: endpoints,
	}

	// Look for an override to IBM Cloud IAM service, preventing searching each time its necessary
	for _, endpoint := range endpoints {
		if endpoint.Name == configv1.IBMCloudServiceIAM {
			client.iamServiceEndpointOverride = endpoint.URL
			break
		}
	}

	if err := client.loadSDKServices(); err != nil {
		return nil, errors.Wrap(err, "failed to load IBM SDK services")
	}

	return client, nil
}

func (c *Client) loadSDKServices() error {
	servicesToLoad := []func() error{
		c.loadGlobalCatalogAPI,
		c.loadResourceManagementAPI,
		c.loadResourceControllerAPI,
		c.loadVPCV1API,
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
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return nil, err
	}
	newIamIdentityV1Options := iamidentityv1.IamIdentityV1Options{
		Authenticator: authenticator,
	}
	// If an IAM service endpoint override was provided, pass it along to override the default IAM service endpoint
	if c.iamServiceEndpointOverride != "" {
		newIamIdentityV1Options.URL = c.iamServiceEndpointOverride
	}
	iamIdentityService, err := iamidentityv1.NewIamIdentityV1(&newIamIdentityV1Options)
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

	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return nil, err
	}
	options := &dnssvcsv1.DnsSvcsV1Options{
		Authenticator: authenticator,
	}
	// If a DNS Services service endpoint override was provided, pass it along to override the default DNS Services service endpoint
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceDNSServices, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}
	// Isolate DNS Services service usage for Internal (Private) clusters; within this function
	dnsService, err := dnssvcsv1.NewDnsSvcsV1(options)
	if err != nil {
		return nil, err
	}

	listPermittedNetworksOptions := dnsService.NewListPermittedNetworksOptions(dnsID, dnsZone)
	permittedNetworks, _, err := dnsService.ListPermittedNetworksWithContext(ctx, listPermittedNetworksOptions)
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
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return nil, err
	}
	// Set CIS DNS record service options
	options := &dnsrecordsv1.DnsRecordsV1Options{
		Authenticator:  authenticator,
		Crn:            core.StringPtr(crnstr),
		ZoneIdentifier: core.StringPtr(zoneID),
	}
	// If a CIS service endpoint override was provided, pass it along to override the default DNS Records service
	// dnsrecordsv1 is provided via IBM CIS endpoint
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceCIS, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}
	// Isolate DNS Records service (for IBM Cloud CIS) usage for External (Public) clusters; within this function
	dnsRecordsService, err := dnsrecordsv1.NewDnsRecordsV1(options)
	if err != nil {
		return nil, err
	}

	// Get CIS DNS records by name
	records, _, err := dnsRecordsService.ListAllDnsRecordsWithContext(ctx, &dnsrecordsv1.ListAllDnsRecordsOptions{
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
		authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
		if err != nil {
			return nil, err
		}
		options := &dnszonesv1.DnsZonesV1Options{
			Authenticator: authenticator,
		}
		// If a DNS Services service endpoint override was provided, pass it along to override the default DNS Zones service
		// dnszonesv1 is provided via IBM Cloud DNS Services endpoint
		if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceDNSServices, c.serviceEndpoints); overrideURL != "" {
			options.URL = overrideURL
		}
		// Isolate DNS Zones service (for IBM Cloud DNS Services) usage for Internal (Private) clusters; within this function
		dnsZoneService, err := dnszonesv1.NewDnsZonesV1(options)
		if err != nil {
			return nil, fmt.Errorf("failed to list DNS zones: %w", err)
		}

		listZonesOptions := dnsZoneService.NewListDnszonesOptions(*instance.GUID)
		result, _, err := dnsZoneService.ListDnszones(listZonesOptions)
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
		authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
		if err != nil {
			return nil, err
		}
		options := &zonesv1.ZonesV1Options{
			Authenticator: authenticator,
			Crn:           instance.CRN,
		}
		// If a CIS service endpoint override was provided, pass it along to override the default Zones service
		// zonesv1 is provided via IBM Cloud CIS endpoint
		if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceCIS, c.serviceEndpoints); overrideURL != "" {
			options.URL = overrideURL
		}
		// Isolate Zones service (for IBM Cloud CIS) usage for External (Public) clusters; within this function
		zonesService, err := zonesv1.NewZonesV1(options)
		if err != nil {
			return nil, fmt.Errorf("failed to list DNS zones: %w", err)
		}

		listZonesOptions := zonesService.NewListZonesOptions()
		listZonesResponse, _, err := zonesService.ListZones(listZonesOptions)

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
	// Parse out the IBM Cloud details from CRN
	crn, err := ibmcrn.Parse(keyCRN)
	if err != nil {
		return nil, err
	}

	var keyResponse *responses.EncryptionKeyResponse

	keyAPI, err := c.getKeyServiceAPI(crn)
	if err != nil {
		return nil, err
	}

	key, err := keyAPI.GetKey(ctx, crn.Resource)
	if err != nil {
		return nil, err
	}

	if key != nil {
		keyResponse = &responses.EncryptionKeyResponse{
			ID:      key.ID,
			Type:    key.Type,
			CRN:     key.CRN,
			State:   key.State,
			Deleted: key.Deleted,
		}
	}

	return keyResponse, nil
}

// GetIBMCloudRegions gets the Regions for IBM Cloud, mapped by shortname to descriptive name.
func (c *Client) GetIBMCloudRegions(ctx context.Context) (map[string]string, error) {
	regions := make(map[string]string, 0)
	// Global Catalog is used to collect the IBM Cloud Regions, via multiple commands.
	geographyOptions := c.catalogAPI.NewListCatalogEntriesOptions()
	geographyOptions.SetQ("kind:geography")

	geographyResult, _, err := c.catalogAPI.ListCatalogEntriesWithContext(ctx, geographyOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list ibmcloud georaphic regions: %w", err)
	}

	countryList, err := c.GetChildrenFromParents(ctx, geographyResult.Resources, "country")
	if err != nil {
		return nil, fmt.Errorf("failed to list ibmcloud countrys: %w", err)
	}
	metroList, err := c.GetChildrenFromParents(ctx, countryList, "metro")
	if err != nil {
		return nil, fmt.Errorf("failed to list ibmcloud metros: %w", err)
	}

	regionList, err := c.GetChildrenFromParents(ctx, metroList, "region")
	if err != nil {
		return nil, fmt.Errorf("failed to list ibmcloud regions: %w", err)
	}

	for _, child := range regionList {
		if child.ID == nil {
			continue
		}
		var description string
		// Attempt to collect the descriptive Region name, although this is based on language, so leaves room for improvement.
		if _, ok := child.OverviewUI["en"]; ok {
			description = *child.OverviewUI["en"].Description
		}
		regions[*child.ID] = description
	}

	return regions, nil
}

// GetChildrenFromParents fetches the children from the IBM Catalog using the given list of parents and the specified kind.
func (c *Client) GetChildrenFromParents(ctx context.Context, parentList []globalcatalogv1.CatalogEntry, kind string) ([]globalcatalogv1.CatalogEntry, error) {
	var childrenList []globalcatalogv1.CatalogEntry

	for num, parent := range parentList {
		if parent.ID == nil {
			fmt.Printf("skipping parent num %d of type %s due to unexpected nil ID\n", num, kind)
			continue
		}

		childrenOptions := c.catalogAPI.NewGetChildObjectsOptions(*parent.ID, kind)
		childrenOptions.SetInclude("*")
		childrenResult, _, err := c.catalogAPI.GetChildObjectsWithContext(ctx, childrenOptions)
		if err != nil {
			return childrenList, fmt.Errorf("failed to retrieve child entries of type %s for %s: %w", kind, *parent.ID, err)
		}
		childrenList = append(childrenList, childrenResult.Resources...)
	}
	return childrenList, nil
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
		return nil, errors.Wrapf(err, "failed to list vpc vsi profiles using: %s", c.vpcAPI.Service.Options.URL)
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
			if detailedResponse != nil {
				// If the response code signifies the VPC was not found, simply move on to the next region; otherwise we log the response
				if detailedResponse.GetStatusCode() != http.StatusNotFound {
					logrus.Warnf("Unexpected response while checking VPC %s in %s region: %s", vpcID, *region.Name, detailedResponse)
				}
			} else {
				logrus.Warnf("Failure collecting VPC %s in %s: %q", vpcID, *region.Name, err)
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
			return nil, fmt.Errorf("failed to set vpc api service url: %w", err)
		}

		if vpcs, detailedResponse, err := c.vpcAPI.ListVpcsWithContext(ctx, c.vpcAPI.NewListVpcsOptions()); err != nil {
			if detailedResponse != nil {
				// If the response code signifies no VPCs were not found, we simply move on to the next region; otherwise log the response
				if detailedResponse.GetStatusCode() != http.StatusNotFound {
					logrus.Warnf("Unexpected response while checking %s region: %s", *region.Name, detailedResponse)
				}
			} else {
				logrus.Warnf("Failure collecting VPCs in %s: %q", *region.Name, err)
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

func (c *Client) loadGlobalCatalogAPI() error {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return fmt.Errorf("failed generating authenticator for global catalog api: %w", err)
	}
	options := &globalcatalogv1.GlobalCatalogV1Options{
		Authenticator: authenticator,
	}
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceGlobalCatalog, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}
	globalCatalogV1Service, err := globalcatalogv1.NewGlobalCatalogV1(options)
	if err != nil {
		return fmt.Errorf("failed creating global catalog api service: %w", err)
	}
	c.catalogAPI = globalCatalogV1Service
	return nil
}

func (c *Client) loadResourceManagementAPI() error {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return err
	}
	options := &resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: authenticator,
	}
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceResourceManager, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}
	resourceManagerV2Service, err := resourcemanagerv2.NewResourceManagerV2(options)
	if err != nil {
		return err
	}
	c.managementAPI = resourceManagerV2Service
	return nil
}

func (c *Client) loadResourceControllerAPI() error {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return err
	}
	options := &resourcecontrollerv2.ResourceControllerV2Options{
		Authenticator: authenticator,
	}
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceResourceController, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}
	resourceControllerV2Service, err := resourcecontrollerv2.NewResourceControllerV2(options)
	if err != nil {
		return err
	}
	c.controllerAPI = resourceControllerV2Service
	return nil
}

func (c *Client) loadVPCV1API() error {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return err
	}
	options := &vpcv1.VpcV1Options{
		Authenticator: authenticator,
	}
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceVPC, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}
	vpcService, err := vpcv1.NewVpcV1(options)
	if err != nil {
		return err
	}
	c.vpcAPI = vpcService
	return nil
}

func (c *Client) getKeyServiceAPI(crn ibmcrn.CRN) (*kpclient.Client, error) {
	var clientConfig kpclient.ClientConfig

	switch crn.ServiceName {
	case hyperProtectCRNServiceName:
		clientConfig = kpclient.ClientConfig{
			BaseURL:    fmt.Sprintf(hyperProtectDefaultURLTemplate, crn.Region),
			APIKey:     c.apiKey,
			TokenURL:   iamTokenDefaultURL,
			InstanceID: crn.ServiceInstance,
		}

		// Override HyperProtect service URL, if one was provided
		if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceHyperProtect, c.serviceEndpoints); overrideURL != "" {
			clientConfig.BaseURL = overrideURL
		}
	case keyProtectCRNServiceName:
		clientConfig = kpclient.ClientConfig{
			BaseURL:    fmt.Sprintf(keyProtectDefaultURLTemplate, crn.Region),
			APIKey:     c.apiKey,
			TokenURL:   iamTokenDefaultURL,
			InstanceID: crn.ServiceInstance,
		}

		// Override KeyProtect service URL, if one was provided
		if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceKeyProtect, c.serviceEndpoints); overrideURL != "" {
			clientConfig.BaseURL = overrideURL
		}
	default:
		return nil, fmt.Errorf("unknown key service for provided encryption key: %s", crn)
	}

	// Override IAM token URL, if an IAM service override URL was provided
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceIAM, c.serviceEndpoints); overrideURL != "" {
		// Construct the token URL using the overridden IAM URL and the token path
		clientConfig.TokenURL = fmt.Sprintf("%s/%s", overrideURL, iamTokenPath)
	}

	return kpclient.New(clientConfig, kpclient.DefaultTransport())
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

package ibmcloud

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
)

//go:generate mockgen -source=./client.go -destination=./mock/ibmcloudclient_generated.go -package=mock

// API represents the calls made to the API.
type API interface {
	GetCISInstance(ctx context.Context, crnstr string) (*resourcecontrollerv2.ResourceInstance, error)
	GetCustomImageByName(ctx context.Context, imageName string, region string) (*vpcv1.Image, error)
	GetCustomImages(ctx context.Context, region string) ([]vpcv1.Image, error)
	GetDNSZones(ctx context.Context) ([]DNSZoneResponse, error)
	GetEncryptionKey(ctx context.Context, keyCRN string) (*EncryptionKeyResponse, error)
	GetResourceGroups(ctx context.Context) ([]resourcemanagerv2.ResourceGroup, error)
	GetResourceGroup(ctx context.Context, nameOrID string) (*resourcemanagerv2.ResourceGroup, error)
	GetSubnet(ctx context.Context, subnetID string) (*vpcv1.Subnet, error)
	GetVSIProfiles(ctx context.Context) ([]vpcv1.InstanceProfile, error)
	GetVPC(ctx context.Context, vpcID string) (*vpcv1.VPC, error)
	GetVPCZonesForRegion(ctx context.Context, region string) ([]string, error)
	GetZoneIDByName(ctx context.Context, name string) (string, error)
}

// Client makes calls to the IBM Cloud API.
type Client struct {
	managementAPI *resourcemanagerv2.ResourceManagerV2
	controllerAPI *resourcecontrollerv2.ResourceControllerV2
	vpcAPI        *vpcv1.VpcV1
	Authenticator *core.IamAuthenticator
}

// cisServiceID is the Cloud Internet Services' catalog service ID.
const cisServiceID = "75874a60-cb12-11e7-948e-37ac098eb1b9"

// VPCResourceNotFoundError represents an error for a VPC resoruce that is not found.
type VPCResourceNotFoundError struct{}

// Error returns the error message for the VPCResourceNotFoundError error type.
func (e *VPCResourceNotFoundError) Error() string {
	return "Not Found"
}

// DNSZoneResponse represents a DNS zone response.
type DNSZoneResponse struct {
	// Name is the domain name of the zone.
	Name string

	// ID is the zone's ID.
	ID string

	// CISInstanceCRN is the IBM Cloud Resource Name for the CIS instance where
	// the DNS zone is managed.
	CISInstanceCRN string

	// CISInstanceName is the display name of the CIS instance where the DNS zone
	// is managed.
	CISInstanceName string

	// ResourceGroupID is the resource group ID of the CIS instance.
	ResourceGroupID string
}

// EncryptionKeyResponse represents an encryption key response.
type EncryptionKeyResponse struct{}

// NewClient initializes a client with a session.
func NewClient() (*Client, error) {
	apiKey := os.Getenv("IC_API_KEY")
	authenticator := &core.IamAuthenticator{
		ApiKey: apiKey,
	}

	client := &Client{
		Authenticator: authenticator,
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
	}

	// Call all the load functions.
	for _, fn := range servicesToLoad {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

// GetCISInstance gets a specific Cloud Internet Services instance by its CRN.
func (c *Client) GetCISInstance(ctx context.Context, crnstr string) (*resourcecontrollerv2.ResourceInstance, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	options := c.controllerAPI.NewGetResourceInstanceOptions(crnstr)
	resourceInstance, _, err := c.controllerAPI.GetResourceInstance(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cis instances")
	}

	return resourceInstance, nil
}

// GetDNSZones returns all of the active DNS zones managed by CIS.
func (c *Client) GetDNSZones(ctx context.Context) ([]DNSZoneResponse, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	options := c.controllerAPI.NewListResourceInstancesOptions()
	options.SetResourceID(cisServiceID)

	listResourceInstancesResponse, _, err := c.controllerAPI.ListResourceInstances(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cis instance")
	}

	var allZones []DNSZoneResponse
	for _, instance := range listResourceInstancesResponse.Resources {
		crnstr := instance.CRN
		zonesService, err := zonesv1.NewZonesV1(&zonesv1.ZonesV1Options{
			Authenticator: c.Authenticator,
			Crn:           crnstr,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to list DNS zones")
		}

		options := zonesService.NewListZonesOptions()
		listZonesResponse, _, _ := zonesService.ListZones(options)

		for _, zone := range listZonesResponse.Result {
			if *zone.Status == "active" {
				zoneStruct := DNSZoneResponse{
					Name:            *zone.Name,
					ID:              *zone.ID,
					CISInstanceCRN:  *instance.CRN,
					CISInstanceName: *instance.Name,
					ResourceGroupID: *instance.ResourceGroupID,
				}
				allZones = append(allZones, zoneStruct)
			}
		}
	}

	return allZones, nil
}

// GetEncryptionKey gets data for an encryption key
func (c *Client) GetEncryptionKey(ctx context.Context, keyCRN string) (*EncryptionKeyResponse, error) {
	// TODO: IBM: Call KMS / Hyperprotect Crpyto APIs.
	return &EncryptionKeyResponse{}, nil
}

// GetZoneIDByName gets the CIS zone ID from its domain name.
func (c *Client) GetZoneIDByName(ctx context.Context, name string) (string, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	zones, err := c.GetDNSZones(ctx)
	if err != nil {
		return "", err
	}

	for _, z := range zones {
		if z.Name == name {
			return z.ID, nil
		}
	}

	return "", fmt.Errorf("zone %q not found", name)
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

	options := c.managementAPI.NewListResourceGroupsOptions()
	listResourceGroupsResponse, _, err := c.managementAPI.ListResourceGroups(options)
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

// GetCustomImages gets a list of custom images within a region. If the image
// status is not "available" it is omitted.
func (c *Client) GetCustomImages(ctx context.Context, region string) ([]vpcv1.Image, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	vpcRegion, err := c.getVPCRegionByName(ctx, region)
	if err != nil {
		return nil, err
	}

	var images []vpcv1.Image
	privateImages, err := c.listPrivateImagesForRegion(ctx, *vpcRegion)
	if err != nil {
		return nil, err
	}
	for _, image := range privateImages {
		if *image.Status == vpcv1.ImageStatusAvailableConst {
			images = append(images, image)
		}
	}
	return images, nil
}

// GetCustomImageByName gets a custom image using its name and region.
func (c *Client) GetCustomImageByName(ctx context.Context, imageName string, region string) (*vpcv1.Image, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	customImages, err := c.GetCustomImages(ctx, region)
	if err != nil {
		return nil, err
	}

	for _, image := range customImages {
		if *image.Name == imageName && *image.Status == vpcv1.ImageStatusAvailableConst {
			return &image, nil
		}
	}
	return nil, fmt.Errorf("image %q not found", imageName)
}

// GetVSIProfiles gets a list of all VSI profiles.
func (c *Client) GetVSIProfiles(ctx context.Context) ([]vpcv1.InstanceProfile, error) {
	listInstanceProfilesOptions := c.vpcAPI.NewListInstanceProfilesOptions()
	profiles, _, err := c.vpcAPI.ListInstanceProfilesWithContext(ctx, listInstanceProfilesOptions)
	return profiles.Profiles, err
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

func (c *Client) getVPCRegionByName(ctx context.Context, regionName string) (*vpcv1.Region, error) {
	region, _, err := c.vpcAPI.GetRegionWithContext(ctx, c.vpcAPI.NewGetRegionOptions(regionName))
	return region, err
}

func (c *Client) listPrivateImagesForRegion(ctx context.Context, region vpcv1.Region) ([]vpcv1.Image, error) {
	listImageOptions := c.vpcAPI.NewListImagesOptions()
	listImageOptions.SetVisibility(vpcv1.ImageVisibilityPrivateConst)

	err := c.vpcAPI.SetServiceURL(fmt.Sprintf("%s/v1", *region.Endpoint))
	if err != nil {
		return nil, errors.Wrap(err, "failed to set vpc api service url")
	}

	listImagesResponse, _, err := c.vpcAPI.ListImagesWithContext(ctx, listImageOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list vpc images")
	}

	return listImagesResponse.Images, nil
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
	options := &resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: c.Authenticator,
	}
	resourceManagerV2Service, err := resourcemanagerv2.NewResourceManagerV2(options)
	if err != nil {
		return err
	}
	c.managementAPI = resourceManagerV2Service
	return nil
}

func (c *Client) loadResourceControllerAPI() error {
	options := &resourcecontrollerv2.ResourceControllerV2Options{
		Authenticator: c.Authenticator,
	}
	resourceControllerV2Service, err := resourcecontrollerv2.NewResourceControllerV2(options)
	if err != nil {
		return err
	}
	c.controllerAPI = resourceControllerV2Service
	return nil
}

func (c *Client) loadVPCV1API() error {
	vpcService, err := vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		Authenticator: c.Authenticator,
	})
	if err != nil {
		return err
	}
	c.vpcAPI = vpcService
	return nil
}

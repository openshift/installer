package validation

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/availabilityzones"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
	"github.com/gophercloud/gophercloud/v2/openstack/common/extensions"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	computequotasets "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/quotasets"
	tokensv2 "github.com/gophercloud/gophercloud/v2/openstack/identity/v2/tokens"
	tokensv3 "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/mtu"
	networkquotasets "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/quotas"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	azutils "github.com/gophercloud/utils/v2/openstack/compute/v2/availabilityzones"
	flavorutils "github.com/gophercloud/utils/v2/openstack/compute/v2/flavors"
	imageutils "github.com/gophercloud/utils/v2/openstack/image/v2/images"
	networkutils "github.com/gophercloud/utils/v2/openstack/networking/v2/networks"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/quota"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
	"github.com/openshift/installer/pkg/types/openstack/validation/networkextensions"
)

// CloudInfo caches data fetched from the user's openstack cloud
type CloudInfo struct {
	APIFIP                  *floatingips.FloatingIP
	ExternalNetwork         *Network
	Flavors                 map[string]Flavor
	IngressFIP              *floatingips.FloatingIP
	ControlPlanePortSubnets []*subnets.Subnet
	ControlPlanePortNetwork *Network
	OSImage                 *images.Image
	ComputeZones            []string
	VolumeZones             []string
	VolumeTypes             []string
	NetworkExtensions       []extensions.Extension
	Quotas                  []quota.Quota
	Networks                []string
	SecurityGroups          []string

	clients *clients
}

// Network holds a gophercloud network with additional info such as MTU.
type Network struct {
	networks.Network
	mtu.NetworkMTUExt
}

type clients struct {
	networkClient  *gophercloud.ServiceClient
	computeClient  *gophercloud.ServiceClient
	imageClient    *gophercloud.ServiceClient
	identityClient *gophercloud.ServiceClient
	volumeClient   *gophercloud.ServiceClient
}

// Flavor embeds information from the Gophercloud Flavor struct and adds
// information on whether a flavor is of baremetal type.
type Flavor struct {
	flavors.Flavor
	Baremetal bool
}

var ci *CloudInfo

// GetCloudInfo fetches and caches metadata from openstack
func GetCloudInfo(ctx context.Context, ic *types.InstallConfig) (*CloudInfo, error) {
	var err error

	if ci != nil {
		return ci, nil
	}

	ci = &CloudInfo{
		clients: &clients{},
		Flavors: map[string]Flavor{},
	}

	opts := openstackdefaults.DefaultClientOpts(ic.OpenStack.Cloud)

	ci.clients.networkClient, err = openstackdefaults.NewServiceClient(ctx, "network", opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create a network client: %w", err)
	}

	ci.clients.computeClient, err = openstackdefaults.NewServiceClient(ctx, "compute", opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create a compute client: %w", err)
	}

	ci.clients.imageClient, err = openstackdefaults.NewServiceClient(ctx, "image", opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create an image client: %w", err)
	}

	ci.clients.identityClient, err = openstackdefaults.NewServiceClient(ctx, "identity", opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create an identity client: %w", err)
	}

	ci.clients.volumeClient, err = openstackdefaults.NewServiceClient(ctx, "volume", opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create a volume client: %w", err)
	}

	err = ci.collectInfo(ctx, ic)
	if err != nil {
		logrus.Warnf("Failed to generate OpenStack cloud info: %v", err)
		return nil, nil
	}

	return ci, nil
}

// I see no reason to artificially split this function into chunks just to make
// the linter happy
//
//nolint:gocyclo
func (ci *CloudInfo) collectInfo(ctx context.Context, ic *types.InstallConfig) error {
	var err error

	ci.ExternalNetwork, err = ci.getNetworkByName(ctx, ic.OpenStack.ExternalNetwork)
	if err != nil {
		return fmt.Errorf("failed to fetch external network info: %w", err)
	}

	// Fetch the image info if the user provided a Glance image name
	imagePtr := ic.OpenStack.ClusterOSImage
	if imagePtr != "" {
		if _, err := url.ParseRequestURI(imagePtr); err != nil {
			ci.OSImage, err = ci.getImage(ctx, imagePtr)
			if err != nil {
				return err
			}
		}
	}

	// Get flavor info
	if ic.Platform.OpenStack.DefaultMachinePlatform != nil {
		if flavorName := ic.Platform.OpenStack.DefaultMachinePlatform.FlavorName; flavorName != "" {
			if _, seen := ci.Flavors[flavorName]; !seen {
				flavor, err := ci.getFlavor(ctx, flavorName)
				if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
					if err != nil {
						return err
					}
					ci.Flavors[flavorName] = flavor
				}
			}
		}
	}

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.OpenStack != nil {
		if flavorName := ic.ControlPlane.Platform.OpenStack.FlavorName; flavorName != "" {
			if _, seen := ci.Flavors[flavorName]; !seen {
				flavor, err := ci.getFlavor(ctx, flavorName)
				if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
					if err != nil {
						return err
					}
					ci.Flavors[flavorName] = flavor
				}
			}
		}
	}

	for _, machine := range ic.Compute {
		if machine.Platform.OpenStack != nil {
			if flavorName := machine.Platform.OpenStack.FlavorName; flavorName != "" {
				if _, seen := ci.Flavors[flavorName]; !seen {
					flavor, err := ci.getFlavor(ctx, flavorName)
					if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
						if err != nil {
							return err
						}
						ci.Flavors[flavorName] = flavor
					}
				}
			}
		}
	}
	if ic.OpenStack.ControlPlanePort != nil {
		controlPlanePort := ic.OpenStack.ControlPlanePort
		ci.ControlPlanePortSubnets, err = ci.getSubnets(ctx, controlPlanePort)
		if err != nil {
			return err
		}

		ci.ControlPlanePortNetwork, err = ci.getNetwork(ctx, controlPlanePort)
		if err != nil {
			return err
		}
	}

	ci.APIFIP, err = ci.getFloatingIP(ctx, ic.OpenStack.APIFloatingIP)
	if err != nil {
		return err
	}

	ci.IngressFIP, err = ci.getFloatingIP(ctx, ic.OpenStack.IngressFloatingIP)
	if err != nil {
		return err
	}

	ci.ComputeZones, err = ci.getComputeZones(ctx)
	if err != nil {
		return err
	}

	ci.VolumeZones, err = ci.getVolumeZones(ctx)
	if err != nil {
		return err
	}

	ci.VolumeTypes, err = ci.getVolumeTypes(ctx)
	if err != nil {
		return err
	}

	ci.Quotas, err = loadQuotas(ctx, ci)
	if err != nil {
		switch {
		case gophercloud.ResponseCodeIs(err, http.StatusForbidden):
			logrus.Warnf("Missing permissions to fetch Quotas and therefore will skip checking them: %v", err)
		case gophercloud.ResponseCodeIs(err, http.StatusNotFound):
			logrus.Warnf("Quota API is not available and therefore will skip checking them: %v", err)
		default:
			return fmt.Errorf("failed to load Quota: %w", err)
		}
	}

	ci.NetworkExtensions, err = networkextensions.Get(ctx, ci.clients.networkClient)
	if err != nil {
		return fmt.Errorf("failed to fetch network extensions: %w", err)
	}

	ci.Networks, err = ci.getNetworks(ctx)
	if err != nil {
		return err
	}

	ci.SecurityGroups, err = ci.getSecurityGroups(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (ci *CloudInfo) getSubnets(ctx context.Context, controlPlanePort *openstack.PortTarget) ([]*subnets.Subnet, error) {
	controlPlaneSubnets := make([]*subnets.Subnet, 0, len(controlPlanePort.FixedIPs))
	for _, fixedIP := range controlPlanePort.FixedIPs {
		page, err := subnets.List(ci.clients.networkClient, subnets.ListOpts{ID: fixedIP.Subnet.ID, Name: fixedIP.Subnet.Name}).AllPages(ctx)
		if err != nil {
			return controlPlaneSubnets, err
		}
		subnetList, err := subnets.ExtractSubnets(page)
		if err != nil {
			return controlPlaneSubnets, err
		}
		if len(subnetList) == 1 {
			controlPlaneSubnets = append(controlPlaneSubnets, &subnetList[0])
		} else if len(subnetList) > 1 {
			return controlPlaneSubnets, fmt.Errorf("found multiple subnets")
		}
	}
	return controlPlaneSubnets, nil
}

func (ci *CloudInfo) getFlavor(ctx context.Context, flavorName string) (Flavor, error) {
	flavorID, err := flavorutils.IDFromName(ctx, ci.clients.computeClient, flavorName)
	if err != nil {
		return Flavor{}, err
	}

	flavor, err := flavors.Get(ctx, ci.clients.computeClient, flavorID).Extract()
	if err != nil {
		return Flavor{}, err
	}

	var baremetal bool
	{
		const baremetalProperty = "baremetal"

		m, err := flavors.GetExtraSpec(ctx, ci.clients.computeClient, flavorID, baremetalProperty).Extract()
		if err != nil && !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			return Flavor{}, err
		}

		if m != nil && m[baremetalProperty] == "true" {
			baremetal = true
		}
	}

	// NOTE(mdbooth): The dereference of flavor is safe here because
	// flavors.Get().Extract() should have raised an error above if the flavor
	// was not found.
	return Flavor{
		Flavor:    *flavor,
		Baremetal: baremetal,
	}, nil
}

// getNetworks returns all the network IDs available on the cloud.
func (ci *CloudInfo) getNetworks(ctx context.Context) ([]string, error) {
	pages, err := networks.List(ci.clients.networkClient, nil).AllPages(ctx)
	if err != nil {
		return nil, err
	}

	networks, err := networks.ExtractNetworks(pages)
	if err != nil {
		return nil, err
	}

	networkIDs := make([]string, len(networks))
	for i := range networks {
		networkIDs[i] = networks[i].ID
	}

	return networkIDs, nil
}

// getSecurityGroups returns all the security group IDs available on the cloud.
func (ci *CloudInfo) getSecurityGroups(ctx context.Context) ([]string, error) {
	pages, err := groups.List(ci.clients.networkClient, groups.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := groups.ExtractGroups(pages)
	if err != nil {
		return nil, err
	}

	sgIDs := make([]string, len(groups))
	for i := range groups {
		sgIDs[i] = groups[i].ID
	}

	return sgIDs, nil
}

func (ci *CloudInfo) getNetworkByName(ctx context.Context, networkName string) (*Network, error) {
	if networkName == "" {
		return nil, nil
	}
	networkID, err := networkutils.IDFromName(ctx, ci.clients.networkClient, networkName)
	if err != nil {
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var network Network
	err = networks.Get(ctx, ci.clients.networkClient, networkID).ExtractInto(&network)
	if err != nil {
		return nil, err
	}

	return &network, nil
}

func (ci *CloudInfo) getNetwork(ctx context.Context, controlPlanePort *openstack.PortTarget) (*Network, error) {
	networkName := controlPlanePort.Network.Name
	networkID := controlPlanePort.Network.ID
	if networkName == "" && networkID == "" {
		if len(ci.ControlPlanePortSubnets) > 0 && ci.ControlPlanePortSubnets[0].NetworkID != "" {
			networkID = ci.ControlPlanePortSubnets[0].NetworkID
		} else {
			return nil, nil
		}
	}
	opts := networks.ListOpts{}
	if networkID != "" {
		opts.ID = networkID
	}
	if networkName != "" {
		opts.Name = controlPlanePort.Network.Name
	}
	allPages, err := networks.List(ci.clients.networkClient, opts).AllPages(ctx)
	if err != nil {
		return nil, err
	}

	var allNetworks []Network
	err = networks.ExtractNetworksInto(allPages, &allNetworks)
	if err != nil {
		return nil, err
	}

	if len(allNetworks) == 0 {
		return nil, nil
	} else if len(allNetworks) > 1 {
		return nil, fmt.Errorf("found multiple networks")
	}

	return &allNetworks[0], nil
}

func (ci *CloudInfo) getFloatingIP(ctx context.Context, fip string) (*floatingips.FloatingIP, error) {
	if fip != "" {
		opts := floatingips.ListOpts{
			FloatingIP: fip,
		}
		allPages, err := floatingips.List(ci.clients.networkClient, opts).AllPages(ctx)
		if err != nil {
			return nil, err
		}

		allFIPs, err := floatingips.ExtractFloatingIPs(allPages)
		if err != nil {
			return nil, err
		}

		if len(allFIPs) == 0 {
			return nil, nil
		}
		return &allFIPs[0], nil
	}
	return nil, nil
}

func (ci *CloudInfo) getImage(ctx context.Context, imageName string) (*images.Image, error) {
	imageID, err := imageutils.IDFromName(ctx, ci.clients.imageClient, imageName)
	if err != nil {
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			return nil, nil
		}
		return nil, err
	}

	image, err := images.Get(ctx, ci.clients.imageClient, imageID).Extract()
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (ci *CloudInfo) getComputeZones(ctx context.Context) ([]string, error) {
	zones, err := azutils.ListAvailableAvailabilityZones(ctx, ci.clients.computeClient)
	if err != nil {
		return nil, fmt.Errorf("failed to list compute availability zones: %w", err)
	}

	if len(zones) == 0 {
		return nil, fmt.Errorf("could not find an available compute availability zone")
	}

	return zones, nil
}

func (ci *CloudInfo) getVolumeZones(ctx context.Context) ([]string, error) {
	allPages, err := availabilityzones.List(ci.clients.volumeClient).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list volume availability zones: %w", err)
	}

	availabilityZoneInfo, err := availabilityzones.ExtractAvailabilityZones(allPages)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response with volume availability zone list: %w", err)
	}

	if len(availabilityZoneInfo) == 0 {
		return nil, fmt.Errorf("could not find an available volume availability zone")
	}

	var zones []string
	for _, zone := range availabilityZoneInfo {
		if zone.ZoneState.Available {
			zones = append(zones, zone.ZoneName)
		}
	}

	return zones, nil
}

func (ci *CloudInfo) getVolumeTypes(ctx context.Context) ([]string, error) {
	allPages, err := volumetypes.List(ci.clients.volumeClient, volumetypes.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list volume types: %w", err)
	}

	volumeTypeInfo, err := volumetypes.ExtractVolumeTypes(allPages)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response with volume types list: %w", err)
	}

	if len(volumeTypeInfo) == 0 {
		return nil, fmt.Errorf("could not find an available block storage volume type")
	}

	var types []string
	for _, volumeType := range volumeTypeInfo {
		types = append(types, volumeType.Name)
	}

	return types, nil
}

// loadQuotas loads the quota information for a project and provided services. It provides information
// about the usage and limit for each resource quota.
func loadQuotas(ctx context.Context, ci *CloudInfo) ([]quota.Quota, error) {
	var quotas []quota.Quota

	projectID, err := getProjectID(ci)
	if err != nil {
		return nil, fmt.Errorf("failed to get keystone project ID: %w", err)
	}

	computeRecords, err := getComputeLimits(ctx, ci, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get compute quota records: %w", err)
	}
	quotas = append(quotas, computeRecords...)

	networkRecords, err := getNetworkLimits(ctx, ci, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network quota records: %w", err)
	}
	quotas = append(quotas, networkRecords...)

	return quotas, nil
}

func getComputeLimits(ctx context.Context, ci *CloudInfo, projectID string) ([]quota.Quota, error) {
	qs, err := computequotasets.GetDetail(ctx, ci.clients.computeClient, projectID).Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to get QuotaSets from OpenStack Compute API: %w", err)
	}

	var quotas []quota.Quota
	addQuota := func(name string, quotaDetail computequotasets.QuotaDetail) {
		quotas = append(quotas, quota.Quota{
			Service:   "compute",
			Name:      name,
			InUse:     int64(quotaDetail.InUse),
			Limit:     int64(quotaDetail.Limit - quotaDetail.Reserved),
			Unlimited: quotaDetail.Limit < 0,
		})
	}
	addQuota("Cores", qs.Cores)
	addQuota("Instances", qs.Instances)
	addQuota("RAM", qs.RAM)
	addQuota("ServerGroups", qs.ServerGroups)
	addQuota("ServerGroupMembers", qs.ServerGroupMembers)

	return quotas, nil
}

func getNetworkLimits(ctx context.Context, ci *CloudInfo, projectID string) ([]quota.Quota, error) {
	qs, err := networkquotasets.GetDetail(ctx, ci.clients.networkClient, projectID).Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to get QuotaSets from OpenStack Network API: %w", err)
	}

	var quotas []quota.Quota
	addQuota := func(name string, quotaDetail networkquotasets.QuotaDetail) {
		quotas = append(quotas, quota.Quota{
			Service:   "network",
			Name:      name,
			InUse:     int64(quotaDetail.Used),
			Limit:     int64(quotaDetail.Limit - quotaDetail.Reserved),
			Unlimited: quotaDetail.Limit < 0,
		})
	}
	addQuota("Port", qs.Port)
	addQuota("Router", qs.Router)
	addQuota("Subnet", qs.Subnet)
	addQuota("Network", qs.Network)
	addQuota("SecurityGroup", qs.SecurityGroup)
	addQuota("SecurityGroupRule", qs.SecurityGroupRule)

	return quotas, nil
}

func getProjectID(ci *CloudInfo) (string, error) {
	authResult := ci.clients.identityClient.GetAuthResult()
	if authResult == nil {
		return "", fmt.Errorf("client did not use openstack.Authenticate()")
	}

	switch authResult.(type) {
	case tokensv2.CreateResult:
		// Gophercloud has support for v2, but keystone has deprecated
		// and it's not even documented.
		return "", fmt.Errorf("extracting project ID using the keystone v2 API is not supported")

	case tokensv3.CreateResult:
		v3Result := authResult.(tokensv3.CreateResult)
		project, err := v3Result.ExtractProject()
		if err != nil {
			return "", fmt.Errorf("extracting project from v3 authResult: %w", err)
		} else if project == nil {
			return "", fmt.Errorf("token is not scoped to a project")
		}
		return project.ID, nil

	default:
		return "", fmt.Errorf("unsupported AuthResult type: %T", authResult)
	}
}

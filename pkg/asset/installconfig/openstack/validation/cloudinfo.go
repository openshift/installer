package validation

import (
	"net/url"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/availabilityzones"
	computequotasets "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	tokensv2 "github.com/gophercloud/gophercloud/openstack/identity/v2/tokens"
	tokensv3 "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	networkquotasets "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/quotas"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/utils/openstack/clientconfig"
	azutils "github.com/gophercloud/utils/openstack/compute/v2/availabilityzones"
	flavorutils "github.com/gophercloud/utils/openstack/compute/v2/flavors"
	imageutils "github.com/gophercloud/utils/openstack/imageservice/v2/images"
	networkutils "github.com/gophercloud/utils/openstack/networking/v2/networks"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/quota"
	"github.com/openshift/installer/pkg/types"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// CloudInfo caches data fetched from the user's openstack cloud
type CloudInfo struct {
	APIFIP          *floatingips.FloatingIP
	ExternalNetwork *networks.Network
	Flavors         map[string]Flavor
	IngressFIP      *floatingips.FloatingIP
	MachinesSubnet  *subnets.Subnet
	OSImage         *images.Image
	ComputeZones    []string
	VolumeZones     []string
	Quotas          []quota.Quota

	clients *clients
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

// record stores the data from quota limits and usages.
type record struct {
	Service string
	Name    string
	Value   int64
}

var ci *CloudInfo

// GetCloudInfo fetches and caches metadata from openstack
func GetCloudInfo(ic *types.InstallConfig) (*CloudInfo, error) {
	var err error

	if ci != nil {
		return ci, nil
	}

	ci = &CloudInfo{
		clients: &clients{},
		Flavors: map[string]Flavor{},
	}

	opts := openstackdefaults.DefaultClientOpts(ic.OpenStack.Cloud)

	ci.clients.networkClient, err = clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a network client")
	}

	ci.clients.computeClient, err = clientconfig.NewServiceClient("compute", opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a compute client")
	}

	ci.clients.imageClient, err = clientconfig.NewServiceClient("image", opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create an image client")
	}

	ci.clients.identityClient, err = clientconfig.NewServiceClient("identity", opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create an identity client")
	}

	ci.clients.volumeClient, err = clientconfig.NewServiceClient("volume", opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a volume client")
	}

	err = ci.collectInfo(ic, opts)
	if err != nil {
		logrus.Warnf("Failed to generate OpenStack cloud info: %v", err)
		return nil, nil
	}

	return ci, nil
}

func (ci *CloudInfo) collectInfo(ic *types.InstallConfig, opts *clientconfig.ClientOpts) error {
	var err error

	ci.ExternalNetwork, err = ci.getNetwork(ic.OpenStack.ExternalNetwork)
	if err != nil {
		return errors.Wrap(err, "failed to fetch external network info")
	}

	// Fetch the image info if the user provided a Glance image name
	imagePtr := ic.OpenStack.ClusterOSImage
	if imagePtr != "" {
		if _, err := url.ParseRequestURI(imagePtr); err != nil {
			ci.OSImage, err = ci.getImage(imagePtr)
			if err != nil {
				return err
			}
		}
	}

	// Get flavor info
	if ic.Platform.OpenStack.DefaultMachinePlatform != nil {
		if flavorName := ic.Platform.OpenStack.DefaultMachinePlatform.FlavorName; flavorName != "" {
			if _, seen := ci.Flavors[flavorName]; !seen {
				ci.Flavors[flavorName], err = ci.getFlavor(flavorName)
				if err != nil {
					return err
				}
			}
		}
	}

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.OpenStack != nil {
		if flavorName := ic.ControlPlane.Platform.OpenStack.FlavorName; flavorName != "" {
			if _, seen := ci.Flavors[flavorName]; !seen {
				ci.Flavors[flavorName], err = ci.getFlavor(flavorName)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, machine := range ic.Compute {
		if machine.Platform.OpenStack != nil {
			if flavorName := machine.Platform.OpenStack.FlavorName; flavorName != "" {
				if _, seen := ci.Flavors[flavorName]; !seen {
					ci.Flavors[flavorName], err = ci.getFlavor(flavorName)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	ci.MachinesSubnet, err = ci.getSubnet(ic.OpenStack.MachinesSubnet)
	if err != nil {
		return errors.Wrap(err, "failed to fetch machine subnet info")
	}

	ci.APIFIP, err = ci.getFloatingIP(ic.OpenStack.APIFloatingIP)
	if err != nil {
		return err
	}

	ci.IngressFIP, err = ci.getFloatingIP(ic.OpenStack.IngressFloatingIP)
	if err != nil {
		return err
	}

	ci.ComputeZones, err = ci.getComputeZones()
	if err != nil {
		return err
	}

	ci.VolumeZones, err = ci.getVolumeZones()
	if err != nil {
		return err
	}

	ci.Quotas, err = loadQuotas(ci)
	if err != nil {
		if isUnauthorized(err) {
			logrus.Warnf("Missing permissions to fetch Quotas and therefore will skip checking them: %v", err)
			return nil
		}
		if isNotFoundError(err) {
			logrus.Warnf("Quota API is not available and therefore will skip checking them: %v", err)
			return nil
		}
		return errors.Wrap(err, "failed to load Quota")
	}

	return nil
}

func (ci *CloudInfo) getSubnet(subnetID string) (*subnets.Subnet, error) {
	if subnetID == "" {
		return nil, nil
	}
	subnet, err := subnets.Get(ci.clients.networkClient, subnetID).Extract()
	if err != nil {
		if isNotFoundError(err) {
			return nil, err
		}
		return nil, err
	}

	return subnet, nil
}

func isNotFoundError(err error) bool {
	var errNotFound gophercloud.ErrResourceNotFound
	var pErrNotFound *gophercloud.ErrResourceNotFound
	var errDefault404 gophercloud.ErrDefault404
	var pErrDefault404 *gophercloud.ErrDefault404

	return errors.As(err, &errNotFound) || errors.As(err, &pErrNotFound) || errors.As(err, &errDefault404) || errors.As(err, &pErrDefault404)
}

func (ci *CloudInfo) getFlavor(flavorName string) (Flavor, error) {
	flavorID, err := flavorutils.IDFromName(ci.clients.computeClient, flavorName)
	if err != nil {
		if isNotFoundError(err) {
			return Flavor{}, nil
		}
		return Flavor{}, err
	}

	flavor, err := flavors.Get(ci.clients.computeClient, flavorID).Extract()
	if err != nil {
		return Flavor{}, err
	}

	var baremetal bool
	{
		const baremetalProperty = "baremetal"

		m, err := flavors.GetExtraSpec(ci.clients.computeClient, flavorID, baremetalProperty).Extract()
		if err != nil && !isNotFoundError(err) {
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

func (ci *CloudInfo) getNetwork(networkName string) (*networks.Network, error) {
	if networkName == "" {
		return nil, nil
	}
	networkID, err := networkutils.IDFromName(ci.clients.networkClient, networkName)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	network, err := networks.Get(ci.clients.networkClient, networkID).Extract()
	if err != nil {
		return nil, err
	}

	return network, nil
}

func (ci *CloudInfo) getFloatingIP(fip string) (*floatingips.FloatingIP, error) {
	if fip != "" {
		opts := floatingips.ListOpts{
			FloatingIP: fip,
		}
		allPages, err := floatingips.List(ci.clients.networkClient, opts).AllPages()
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

func (ci *CloudInfo) getImage(imageName string) (*images.Image, error) {
	imageID, err := imageutils.IDFromName(ci.clients.imageClient, imageName)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	image, err := images.Get(ci.clients.imageClient, imageID).Extract()
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (ci *CloudInfo) getComputeZones() ([]string, error) {
	zones, err := azutils.ListAvailableAvailabilityZones(ci.clients.computeClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list compute availability zones")
	}

	if len(zones) == 0 {
		return nil, errors.New("could not find an available compute availability zone")
	}

	return zones, nil
}

func (ci *CloudInfo) getVolumeZones() ([]string, error) {
	allPages, err := availabilityzones.List(ci.clients.volumeClient).AllPages()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list volume availability zones")
	}

	availabilityZoneInfo, err := availabilityzones.ExtractAvailabilityZones(allPages)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse response with volume availability zone list")
	}

	if len(availabilityZoneInfo) == 0 {
		return nil, errors.New("could not find an available compute availability zone")
	}

	var zones []string
	for _, zone := range availabilityZoneInfo {
		if zone.ZoneState.Available {
			zones = append(zones, zone.ZoneName)
		}
	}

	return zones, nil
}

// loadLimits loads the consumer quota metric.
func loadLimits(ci *CloudInfo) ([]record, error) {
	var limits []record

	projectID, err := getProjectID(ci)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get keystone project ID")
	}

	computeRecords, err := getComputeLimits(ci, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get compute quota records")
	}
	for _, r := range computeRecords {
		limits = append(limits, r)
	}

	networkRecords, err := getNetworkLimits(ci, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get network quota records")
	}
	for _, r := range networkRecords {
		limits = append(limits, r)
	}

	return limits, nil
}

func getComputeLimits(ci *CloudInfo, projectID string) ([]record, error) {
	qs, err := computequotasets.GetDetail(ci.clients.computeClient, projectID).Extract()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get QuotaSets from OpenStack Compute API")
	}

	var records []record
	addRecord := func(name string, quota computequotasets.QuotaDetail) {
		qval := int64(quota.Limit - quota.InUse - quota.Reserved)
		// -1 means unlimited in OpenStack so we will ignore that record.
		if quota.Limit == -1 {
			qval = -1
		}
		records = append(records, record{
			Service: "compute",
			Name:    name,
			Value:   qval,
		})
	}
	addRecord("Cores", qs.Cores)
	addRecord("Instances", qs.Instances)
	addRecord("RAM", qs.RAM)

	return records, nil
}

func getNetworkLimits(ci *CloudInfo, projectID string) ([]record, error) {
	qs, err := networkquotasets.GetDetail(ci.clients.networkClient, projectID).Extract()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get QuotaSets from OpenStack Network API")
	}

	var records []record
	addRecord := func(name string, quota networkquotasets.QuotaDetail) {
		qval := int64(quota.Limit - quota.Used - quota.Reserved)
		// -1 means unlimited in OpenStack so we will ignore that record.
		if quota.Limit == -1 {
			qval = -1
		}
		records = append(records, record{
			Service: "network",
			Name:    name,
			Value:   qval,
		})
	}
	addRecord("Port", qs.Port)
	addRecord("Router", qs.Router)
	addRecord("Subnet", qs.Subnet)
	addRecord("Network", qs.Network)
	addRecord("SecurityGroup", qs.SecurityGroup)
	addRecord("SecurityGroupRule", qs.SecurityGroupRule)

	return records, nil
}

func getProjectID(ci *CloudInfo) (string, error) {
	authResult := ci.clients.identityClient.GetAuthResult()
	if authResult == nil {
		return "", errors.Errorf("Client did not use openstack.Authenticate()")
	}

	switch authResult.(type) {
	case tokensv2.CreateResult:
		// Gophercloud has support for v2, but keystone has deprecated
		// and it's not even documented.
		return "", errors.Errorf("Extracting project ID using the keystone v2 API is not supported")

	case tokensv3.CreateResult:
		v3Result := authResult.(tokensv3.CreateResult)
		project, err := v3Result.ExtractProject()
		if err != nil {
			return "", errors.Wrap(err, "Extracting project from v3 authResult")
		} else if project == nil {
			return "", errors.Errorf("Token is not scoped to a project")
		}
		return project.ID, nil

	default:
		return "", errors.Errorf("Unsupported AuthResult type: %T", authResult)
	}
}

// loadQuotas loads the quota information for a project and provided services. It provides information
// about the usage and limit for each resource quota.
func loadQuotas(ci *CloudInfo) ([]quota.Quota, error) {
	records, err := loadLimits(ci)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load quota limits")
	}
	return newQuota(records), nil
}

func newQuota(limits []record) []quota.Quota {
	var ret []quota.Quota
	for _, limit := range limits {
		isUnlimited := limit.Value == -1
		q := quota.Quota{
			Service: limit.Service,
			Name:    limit.Name,
			// Since limit.Value contains the actual
			// available resources, we can set InUse to 0.
			InUse:     0,
			Limit:     limit.Value,
			Unlimited: isUnlimited,
		}
		ret = append(ret, q)
	}
	return ret
}

// isUnauthorized checks if the error is unauthorized (http code 403)
func isUnauthorized(err error) bool {
	if err == nil {
		return false
	}
	var gErr gophercloud.ErrDefault403
	return errors.As(err, &gErr)
}

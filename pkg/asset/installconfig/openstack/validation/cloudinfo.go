package validation

import (
	"net/url"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/availabilityzones"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/utils/openstack/clientconfig"
	flavorutils "github.com/gophercloud/utils/openstack/compute/v2/flavors"
	imageutils "github.com/gophercloud/utils/openstack/imageservice/v2/images"
	networkutils "github.com/gophercloud/utils/openstack/networking/v2/networks"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

// CloudInfo caches data fetched from the user's openstack cloud
type CloudInfo struct {
	APIFIP          *floatingips.FloatingIP
	ExternalNetwork *networks.Network
	Flavors         map[string]Flavor
	IngressFIP      *floatingips.FloatingIP
	MachinesSubnet  *subnets.Subnet
	OSImage         *images.Image
	Zones           []string

	clients *clients
}

type clients struct {
	networkClient *gophercloud.ServiceClient
	computeClient *gophercloud.ServiceClient
	imageClient   *gophercloud.ServiceClient
}

// Flavor embeds information from the Gophercloud Flavor struct and adds
// information on whether a flavor is of baremetal type.
type Flavor struct {
	*flavors.Flavor
	Baremetal bool
}

// GetCloudInfo fetches and caches metadata from openstack
func GetCloudInfo(ic *types.InstallConfig) (*CloudInfo, error) {
	var err error
	ci := CloudInfo{
		clients: &clients{},
		Flavors: map[string]Flavor{},
	}

	opts := &clientconfig.ClientOpts{Cloud: ic.OpenStack.Cloud}

	// We should unset OS_CLOUD env variable here, because the real cloud name was
	// defined on the previous step. OS_CLOUD has more priority, so the value from
	// "opts" variable will be ignored if OS_CLOUD contains something.
	os.Unsetenv("OS_CLOUD")

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

	err = ci.collectInfo(ic)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate OpenStack cloud info")
	}

	return &ci, nil
}

func (ci *CloudInfo) collectInfo(ic *types.InstallConfig) error {
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

	ci.Zones, err = ci.getZones()
	if err != nil {
		return err
	}

	return nil
}

func (ci *CloudInfo) getSubnet(subnetID string) (*subnets.Subnet, error) {
	subnet, err := subnets.Get(ci.clients.networkClient, subnetID).Extract()
	if err != nil {
		var gerr *gophercloud.ErrResourceNotFound
		if errors.As(err, &gerr) {
			return nil, nil
		}
		return nil, err
	}

	return subnet, nil
}

func isNotFoundError(err error) bool {
	var errNotFound *gophercloud.ErrResourceNotFound
	return errors.As(err, &errNotFound) || strings.Contains(err.Error(), "Resource not found")
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

	return Flavor{
		Flavor:    flavor,
		Baremetal: baremetal,
	}, nil
}

func (ci *CloudInfo) getNetwork(networkName string) (*networks.Network, error) {
	if networkName == "" {
		return nil, nil
	}
	networkID, err := networkutils.IDFromName(ci.clients.networkClient, networkName)
	if err != nil {
		var gerr *gophercloud.ErrResourceNotFound
		if errors.As(err, &gerr) {
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
		var gerr *gophercloud.ErrResourceNotFound
		if errors.As(err, &gerr) {
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

func (ci *CloudInfo) getZones() ([]string, error) {
	zones := []string{}
	allPages, err := availabilityzones.List(ci.clients.computeClient).AllPages()
	if err != nil {
		return nil, err
	}

	availabilityZoneInfo, err := availabilityzones.ExtractAvailabilityZones(allPages)
	if err != nil {
		return nil, err
	}

	for _, zoneInfo := range availabilityZoneInfo {
		if zoneInfo.ZoneState.Available {
			zones = append(zones, zoneInfo.ZoneName)
		}
	}

	if len(zones) == 0 {
		return nil, errors.New("could not find an available compute availability zone")
	}

	return zones, nil
}

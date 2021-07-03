// Package openstack contains OpenStack-specific Terraform-variable logic.
package openstack

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	"github.com/openshift/installer/pkg/types"
	types_openstack "github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"
)

type config struct {
	BaseImageName                    string                            `json:"openstack_base_image_name,omitempty"`
	ExternalNetwork                  string                            `json:"openstack_external_network,omitempty"`
	Cloud                            string                            `json:"openstack_credentials_cloud,omitempty"`
	FlavorName                       string                            `json:"openstack_master_flavor_name,omitempty"`
	APIFloatingIP                    string                            `json:"openstack_api_floating_ip,omitempty"`
	IngressFloatingIP                string                            `json:"openstack_ingress_floating_ip,omitempty"`
	APIVIP                           string                            `json:"openstack_api_int_ip,omitempty"`
	IngressVIP                       string                            `json:"openstack_ingress_ip,omitempty"`
	TrunkSupport                     bool                              `json:"openstack_trunk_support,omitempty"`
	OctaviaSupport                   bool                              `json:"openstack_octavia_support,omitempty"`
	RootVolumeSize                   int                               `json:"openstack_master_root_volume_size,omitempty"`
	RootVolumeType                   string                            `json:"openstack_master_root_volume_type,omitempty"`
	BootstrapShim                    string                            `json:"openstack_bootstrap_shim_ignition,omitempty"`
	ExternalDNS                      []string                          `json:"openstack_external_dns,omitempty"`
	MasterServerGroupName            string                            `json:"openstack_master_server_group_name,omitempty"`
	MasterServerGroupPolicy          types_openstack.ServerGroupPolicy `json:"openstack_master_server_group_policy"`
	AdditionalNetworkIDs             []string                          `json:"openstack_additional_network_ids,omitempty"`
	AdditionalSecurityGroupIDs       []string                          `json:"openstack_master_extra_sg_ids,omitempty"`
	MachinesSubnet                   string                            `json:"openstack_machines_subnet_id,omitempty"`
	MachinesNetwork                  string                            `json:"openstack_machines_network_id,omitempty"`
	MasterAvailabilityZones          []string                          `json:"openstack_master_availability_zones,omitempty"`
	MasterRootVolumeAvalabilityZones []string                          `json:"openstack_master_root_volume_availability_zones,omitempty"`
}

// TFVars generates OpenStack-specific Terraform variables.
func TFVars(masterConfigs []*v1alpha1.OpenstackProviderSpec, cloud string, externalNetwork string, externalDNS []string, apiFloatingIP string, ingressFloatingIP string, apiVIP string, ingressVIP string, baseImage string, baseImageProperties map[string]string, infraID string, userCA string, bootstrapIgn string, mpool, defaultmpool *types_openstack.MachinePool, machinesSubnet string, proxy *types.Proxy) ([]byte, error) {
	zones := []string{}
	seen := map[string]bool{}
	for _, config := range masterConfigs {
		if !seen[config.AvailabilityZone] {
			zones = append(zones, config.AvailabilityZone)
			seen[config.AvailabilityZone] = true
		}
	}

	cfg := &config{
		ExternalNetwork:         externalNetwork,
		Cloud:                   cloud,
		FlavorName:              masterConfigs[0].Flavor,
		APIFloatingIP:           apiFloatingIP,
		IngressFloatingIP:       ingressFloatingIP,
		APIVIP:                  apiVIP,
		IngressVIP:              ingressVIP,
		ExternalDNS:             externalDNS,
		MachinesSubnet:          machinesSubnet,
		MasterAvailabilityZones: zones,
	}

	if defaultmpool != nil && defaultmpool.RootVolume != nil {
		cfg.MasterRootVolumeAvalabilityZones = defaultmpool.RootVolume.Zones
	}
	if mpool != nil && mpool.RootVolume != nil && mpool.RootVolume.Zones != nil {
		cfg.MasterRootVolumeAvalabilityZones = mpool.RootVolume.Zones
	}

	serviceCatalog, err := getServiceCatalog(cloud)
	if err != nil {
		return nil, errors.Errorf("Could not retrieve service catalog: %v", err)
	}

	// Normally baseImage contains a URL that we will use to create a new Glance image, but for testing
	// purposes we also allow to set a custom Glance image name to skip the uploading. Here we check
	// whether baseImage is a URL or not. If this is the first case, it means that the image should be
	// created by the installer from the URL. Otherwise, it means that we are given the name of the pre-created
	// Glance image, which we should use for instances.
	imageName, isURL := rhcos.GenerateOpenStackImageName(baseImage, infraID)
	cfg.BaseImageName = imageName
	if isURL {
		// Valid URL -> use baseImage as a URL that will be used to create new Glance image with name "<infraID>-rhcos".
		var localFilePath string

		url, err := url.Parse(baseImage)
		if err != nil {
			return nil, err
		}

		// We support 'http(s)' and 'file' schemes. If the scheme is http(s), then we will upload a file from that
		// location. Otherwise will take local file path from the URL.
		switch url.Scheme {
		case "http", "https":
			localFilePath, err = cache.DownloadImageFile(baseImage)
			if err != nil {
				return nil, err
			}
		case "file":
			localFilePath = filepath.FromSlash(url.Path)
		default:
			return nil, errors.Errorf("Unsupported URL scheme: '%v'", url.Scheme)
		}

		err = uploadBaseImage(cloud, localFilePath, imageName, infraID, baseImageProperties)
		if err != nil {
			return nil, err
		}
	}

	clientConfigCloud, err := clientconfig.GetCloudFromYAML(openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return nil, err
	}

	glancePublicURL, err := getGlancePublicURL(serviceCatalog, clientConfigCloud.RegionName)
	if err != nil {
		return nil, err
	}

	configLocation, err := uploadBootstrapConfig(cloud, bootstrapIgn, infraID)
	if err != nil {
		return nil, err
	}

	tokenID, err := getAuthToken(cloud)
	if err != nil {
		return nil, err
	}

	bootstrapConfigURL := fmt.Sprintf("%s%s", glancePublicURL, configLocation)
	userCAIgnition, err := generateIgnitionShim(userCA, infraID, bootstrapConfigURL, tokenID, proxy)
	if err != nil {
		return nil, err
	}

	cfg.BootstrapShim = userCAIgnition
	masterConfig := masterConfigs[0]

	cfg.TrunkSupport = masterConfig.Trunk

	cfg.OctaviaSupport, err = isOctaviaSupported(serviceCatalog)
	if err != nil {
		return nil, err
	}

	if masterConfig.RootVolume != nil {
		cfg.RootVolumeSize = masterConfig.RootVolume.Size
		cfg.RootVolumeType = masterConfig.RootVolume.VolumeType
	}

	cfg.MasterServerGroupName = masterConfig.ServerGroupName

	if mpool != nil && mpool.ServerGroupPolicy != types_openstack.SGPolicyUnset {
		cfg.MasterServerGroupPolicy = mpool.ServerGroupPolicy
	} else if defaultmpool != nil && defaultmpool.ServerGroupPolicy != types_openstack.SGPolicyUnset {
		cfg.MasterServerGroupPolicy = defaultmpool.ServerGroupPolicy
	} else {
		cfg.MasterServerGroupPolicy = types_openstack.SGPolicySoftAntiAffinity
	}

	if masterConfig.ServerGroupID != "" {
		return nil, errors.Errorf("ServerGroupID is not implemented in the Installer. Please use ServerGroupName for automatic creation of the Control Plane server group.")
	}

	cfg.AdditionalNetworkIDs = []string{}
	if mpool != nil {
		cfg.AdditionalNetworkIDs = append(cfg.AdditionalNetworkIDs, mpool.AdditionalNetworkIDs...)
	}

	cfg.AdditionalSecurityGroupIDs = []string{}
	if mpool != nil {
		cfg.AdditionalSecurityGroupIDs = append(cfg.AdditionalSecurityGroupIDs, mpool.AdditionalSecurityGroupIDs...)
	}

	if machinesSubnet != "" {
		cfg.MachinesNetwork, err = getNetworkFromSubnet(cloud, machinesSubnet)
		if err != nil {
			return nil, err
		}
	}

	return json.MarshalIndent(cfg, "", "  ")
}

// We need to obtain Glance public endpoint that will be used by Ignition to download bootstrap ignition files.
// By design this should be done by using https://www.terraform.io/docs/providers/openstack/d/identity_endpoint_v3.html
// but OpenStack default policies forbid to use this API for regular users.
// On the other hand when a user authenticates in OpenStack (i.e. gets a token), it includes the whole service
// catalog in the output json. So we are able to parse the data and get the endpoint from there
// https://docs.openstack.org/api-ref/identity/v3/?expanded=token-authentication-with-scoped-authorization-detail#token-authentication-with-scoped-authorization
// Unfortunately this feature is not currently supported by Terraform, so we had to implement it here.
// We do next:
// 1. In "getServiceCatalog" we authenticate in OpenStack (tokens.Create(..)),
//    parse the token and extract the service catalog: (ExtractServiceCatalog())
// 2. In getGlancePublicURL we iterate through the catalog and find "public" endpoint for "image".

// getGlancePublicURL obtains Glance public endpoint URL
func getGlancePublicURL(serviceCatalog *tokens.ServiceCatalog, region string) (string, error) {
	glancePublicURL, err := openstack.V3EndpointURL(serviceCatalog, gophercloud.EndpointOpts{
		Type:         "image",
		Availability: gophercloud.AvailabilityPublic,
		Region:       region,
	})
	if err != nil {
		return "", errors.Errorf("cannot retrieve Glance URL from the service catalog: %v", err)
	}

	return glancePublicURL, nil
}

// getServiceCatalog fetches OpenStack service catalog with service endpoints
func getServiceCatalog(cloud string) (*tokens.ServiceCatalog, error) {
	conn, err := clientconfig.NewServiceClient("identity", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return nil, err
	}

	authResult := conn.GetAuthResult()
	auth, ok := authResult.(tokens.CreateResult)
	if !ok {
		return nil, errors.New("unable to extract service catalog")
	}

	serviceCatalog, err := auth.ExtractServiceCatalog()
	if err != nil {
		return nil, err
	}

	return serviceCatalog, nil
}

// getNetworkFromSubnet looks up a subnet in openstack and returns the ID of the network it's a part of
func getNetworkFromSubnet(cloud string, subnetID string) (string, error) {
	networkClient, err := clientconfig.NewServiceClient("network", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return "", err
	}

	subnet, err := subnets.Get(networkClient, subnetID).Extract()
	if err != nil {
		return "", err
	}

	return subnet.NetworkID, nil
}

func isOctaviaSupported(serviceCatalog *tokens.ServiceCatalog) (bool, error) {
	_, err := openstack.V3EndpointURL(serviceCatalog, gophercloud.EndpointOpts{
		Type:         "load-balancer",
		Name:         "octavia",
		Availability: gophercloud.AvailabilityPublic,
	})
	if err != nil {
		if _, ok := err.(*gophercloud.ErrEndpointNotFound); ok {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

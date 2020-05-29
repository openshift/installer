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
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	types_openstack "github.com/openshift/installer/pkg/types/openstack"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"
)

type config struct {
	BaseImageName              string   `json:"openstack_base_image_name,omitempty"`
	ExternalNetwork            string   `json:"openstack_external_network,omitempty"`
	Cloud                      string   `json:"openstack_credentials_cloud,omitempty"`
	FlavorName                 string   `json:"openstack_master_flavor_name,omitempty"`
	LbFloatingIP               string   `json:"openstack_lb_floating_ip,omitempty"`
	APIVIP                     string   `json:"openstack_api_int_ip,omitempty"`
	DNSVIP                     string   `json:"openstack_node_dns_ip,omitempty"`
	IngressVIP                 string   `json:"openstack_ingress_ip,omitempty"`
	TrunkSupport               string   `json:"openstack_trunk_support,omitempty"`
	OctaviaSupport             string   `json:"openstack_octavia_support,omitempty"`
	RootVolumeSize             int      `json:"openstack_master_root_volume_size,omitempty"`
	RootVolumeType             string   `json:"openstack_master_root_volume_type,omitempty"`
	BootstrapShim              string   `json:"openstack_bootstrap_shim_ignition,omitempty"`
	ExternalDNS                []string `json:"openstack_external_dns,omitempty"`
	MasterServerGroupName      string   `json:"openstack_master_server_group_name,omitempty"`
	AdditionalNetworkIDs       []string `json:"openstack_additional_network_ids,omitempty"`
	AdditionalSecurityGroupIDs []string `json:"openstack_master_extra_sg_ids,omitempty"`
	MachinesSubnet             string   `json:"openstack_machines_subnet_id,omitempty"`
	MachinesNetwork            string   `json:"openstack_machines_network_id,omitempty"`
}

// TFVars generates OpenStack-specific Terraform variables.
func TFVars(masterConfig *v1alpha1.OpenstackProviderSpec, cloud string, externalNetwork string, externalDNS []string, lbFloatingIP string, apiVIP string, dnsVIP string, ingressVIP string, trunkSupport string, octaviaSupport string, baseImage string, infraID string, userCA string, bootstrapIgn string, mpool *types_openstack.MachinePool, machinesSubnet string) ([]byte, error) {

	cfg := &config{
		ExternalNetwork: externalNetwork,
		Cloud:           cloud,
		FlavorName:      masterConfig.Flavor,
		LbFloatingIP:    lbFloatingIP,
		APIVIP:          apiVIP,
		DNSVIP:          dnsVIP,
		IngressVIP:      ingressVIP,
		ExternalDNS:     externalDNS,
		TrunkSupport:    trunkSupport,
		OctaviaSupport:  octaviaSupport,
		MachinesSubnet:  machinesSubnet,
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

		err = uploadBaseImage(cloud, localFilePath, imageName, infraID)
		if err != nil {
			return nil, err
		}
	} else {
		// Not a URL -> use baseImage value as an overridden Glance image name.
		// Need to check if this image exists and there are no other images with this name.
		err := validateOverriddenImageName(imageName, cloud)
		if err != nil {
			return nil, err
		}
	}

	glancePublicURL, err := getGlancePublicURL(cloud)
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
	userCAIgnition, err := generateIgnitionShim(userCA, infraID, bootstrapConfigURL, tokenID)
	if err != nil {
		return nil, err
	}

	cfg.BootstrapShim = userCAIgnition

	if masterConfig.RootVolume != nil {
		cfg.RootVolumeSize = masterConfig.RootVolume.Size
		cfg.RootVolumeType = masterConfig.RootVolume.VolumeType
	}

	cfg.MasterServerGroupName = masterConfig.ServerGroupName

	if masterConfig.ServerGroupID != "" {
		return nil, errors.Errorf("ServerGroupID is not implemented in the Installer. Please use ServerGroupName for automatic creation of the Control Plane server group.")
	}

	cfg.AdditionalNetworkIDs = []string{}
	if mpool.AdditionalNetworkIDs != nil {
		for _, networkID := range mpool.AdditionalNetworkIDs {
			cfg.AdditionalNetworkIDs = append(cfg.AdditionalNetworkIDs, networkID)
		}
	}

	cfg.AdditionalSecurityGroupIDs = []string{}
	if mpool.AdditionalSecurityGroupIDs != nil {
		for _, sgID := range mpool.AdditionalSecurityGroupIDs {
			cfg.AdditionalSecurityGroupIDs = append(cfg.AdditionalSecurityGroupIDs, sgID)
		}
	}

	if machinesSubnet != "" {
		cfg.MachinesNetwork, err = getNetworkFromSubnet(cloud, machinesSubnet)
		if err != nil {
			return nil, err
		}

		// Make sure that the network has the primary cluster network tag.
		// In the case of multiple networks this tag is required for
		// cluster-api-provider-openstack to define which ip address to set as
		// the primary one.
		err = setNetworkTag(cloud, cfg.MachinesNetwork, infraID+"-primaryClusterNetwork")
		if err != nil {
			return nil, err
		}
	}

	return json.MarshalIndent(cfg, "", "  ")
}

func validateOverriddenImageName(imageName, cloud string) error {
	opts := &clientconfig.ClientOpts{
		Cloud: cloud,
	}

	client, err := clientconfig.NewServiceClient("image", opts)
	if err != nil {
		return err
	}

	listOpts := images.ListOpts{
		Name: imageName,
	}

	allPages, err := images.List(client, listOpts).AllPages()
	if err != nil {
		return err
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return err
	}

	if len(allImages) == 0 {
		return errors.Errorf("image '%v' doesn't exist", imageName)
	}

	if len(allImages) > 1 {
		return errors.Errorf("there's more than one image with the name '%v'", imageName)
	}

	return nil
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
func getGlancePublicURL(cloud string) (string, error) {
	serviceCatalog, err := getServiceCatalog(cloud)
	if err != nil {
		return "", err
	}

	glancePublicURL, err := openstack.V3EndpointURL(serviceCatalog, gophercloud.EndpointOpts{
		Type:         "image",
		Availability: gophercloud.AvailabilityPublic,
	})
	if err != nil {
		return "", errors.Errorf("cannot retrieve Glance URL from the service catalog: %v", err)
	}

	return glancePublicURL, nil
}

// getServiceCatalog fetches OpenStack service catalog with service endpoints
func getServiceCatalog(cloud string) (*tokens.ServiceCatalog, error) {
	opts := &clientconfig.ClientOpts{
		Cloud: cloud,
	}

	conn, err := clientconfig.NewServiceClient("identity", opts)
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
	opts := &clientconfig.ClientOpts{
		Cloud: cloud,
	}

	networkClient, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return "", err
	}

	subnet, err := subnets.Get(networkClient, subnetID).Extract()
	if err != nil {
		return "", err
	}

	return subnet.NetworkID, nil
}

// setNetworkTag sets a tag for the network
func setNetworkTag(cloud string, networkID string, networkTag string) error {
	opts := &clientconfig.ClientOpts{
		Cloud: cloud,
	}

	networkClient, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return err
	}

	err = attributestags.Add(networkClient, "networks", networkID, networkTag).ExtractErr()
	if err != nil {
		return err
	}

	return nil
}

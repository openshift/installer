// Package openstack contains OpenStack-specific Terraform-variable logic.
package openstack

import (
	"encoding/json"

	"github.com/openshift/installer/pkg/rhcos"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"
)

type config struct {
	BaseImageName   string `json:"openstack_base_image_name,omitempty"`
	BaseImageURL    string `json:"openstack_base_image_url,omitempty"`
	ExternalNetwork string `json:"openstack_external_network,omitempty"`
	Cloud           string `json:"openstack_credentials_cloud,omitempty"`
	FlavorName      string `json:"openstack_master_flavor_name,omitempty"`
	LbFloatingIP    string `json:"openstack_lb_floating_ip,omitempty"`
	APIVIP          string `json:"openstack_api_int_ip,omitempty"`
	DNSVIP          string `json:"openstack_node_dns_ip,omitempty"`
	IngressVIP      string `json:"openstack_ingress_ip,omitempty"`
	TrunkSupport    string `json:"openstack_trunk_support,omitempty"`
	OctaviaSupport  string `json:"openstack_octavia_support,omitempty"`
	RootVolumeSize  int    `json:"openstack_master_root_volume_size,omitempty"`
	RootVolumeType  string `json:"openstack_master_root_volume_type,omitempty"`
}

// TFVars generates OpenStack-specific Terraform variables.
func TFVars(masterConfig *v1alpha1.OpenstackProviderSpec, cloud string, externalNetwork string, lbFloatingIP string, apiVIP string, dnsVIP string, ingressVIP string, trunkSupport string, octaviaSupport string, baseImage string, infraID string) ([]byte, error) {
	cfg := &config{
		ExternalNetwork: externalNetwork,
		Cloud:           cloud,
		FlavorName:      masterConfig.Flavor,
		LbFloatingIP:    lbFloatingIP,
		APIVIP:          apiVIP,
		DNSVIP:          dnsVIP,
		IngressVIP:      ingressVIP,
		TrunkSupport:    trunkSupport,
		OctaviaSupport:  octaviaSupport,
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
		cfg.BaseImageURL = baseImage
	} else {
		// Not a URL -> use baseImage value as a Glance image name.

		// TODO(mfedosin): add validations that this image exists and there are no other images with this name.
	}

	if masterConfig.RootVolume != nil {
		cfg.RootVolumeSize = masterConfig.RootVolume.Size
		cfg.RootVolumeType = masterConfig.RootVolume.VolumeType
	}

	return json.MarshalIndent(cfg, "", "  ")
}

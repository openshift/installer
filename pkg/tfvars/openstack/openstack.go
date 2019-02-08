// Package openstack contains OpenStack-specific Terraform-variable logic.
package openstack

import (
	"encoding/json"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"
)

type config struct {
	Region          string `json:"openstack_region,omitempty"`
	BaseImage       string `json:"openstack_base_image,omitempty"`
	ExternalNetwork string `json:"openstack_external_network,omitempty"`
	Cloud           string `json:"openstack_credentials_cloud,omitempty"`
	FlavorName      string `json:"openstack_control_plane_flavor_name,omitempty"`
	LbFloatingIP    string `json:"openstack_lb_floating_ip,omitempty"`
	TrunkSupport    string `json:"openstack_trunk_support,omitempty"`
}

// TFVars generates OpenStack-specific Terraform variables.
func TFVars(controlPlaneConfig *v1alpha1.OpenstackProviderSpec, region string, externalNetwork string, lbFloatingIP string, trunkSupport string) ([]byte, error) {
	cfg := &config{
		Region:          region,
		BaseImage:       controlPlaneConfig.Image,
		ExternalNetwork: externalNetwork,
		Cloud:           controlPlaneConfig.CloudName,
		FlavorName:      controlPlaneConfig.Flavor,
		LbFloatingIP:    lbFloatingIP,
		TrunkSupport:    trunkSupport,
	}

	return json.MarshalIndent(cfg, "", "  ")
}

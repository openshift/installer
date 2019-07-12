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
	FlavorName      string `json:"openstack_master_flavor_name,omitempty"`
	LbFloatingIP    string `json:"openstack_lb_floating_ip,omitempty"`
	APIVIP          string `json:"openstack_api_int_ip,omitempty"`
	DNSVIP          string `json:"openstack_node_dns_ip,omitempty"`
	TrunkSupport    string `json:"openstack_trunk_support,omitempty"`
	OctaviaSupport  string `json:"openstack_octavia_support,omitempty"`
}

// TFVars generates OpenStack-specific Terraform variables.
func TFVars(masterConfig *v1alpha1.OpenstackProviderSpec, region string, externalNetwork string, lbFloatingIP string, apiVIP string, dnsVIP string, trunkSupport string, octaviaSupport string) ([]byte, error) {
	cfg := &config{
		Region:          region,
		BaseImage:       masterConfig.Image,
		ExternalNetwork: externalNetwork,
		Cloud:           masterConfig.CloudName,
		FlavorName:      masterConfig.Flavor,
		LbFloatingIP:    lbFloatingIP,
		APIVIP:          apiVIP,
		DNSVIP:          dnsVIP,
		TrunkSupport:    trunkSupport,
		OctaviaSupport:  octaviaSupport,
	}

	return json.MarshalIndent(cfg, "", "  ")
}

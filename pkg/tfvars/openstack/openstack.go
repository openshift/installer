// Package openstack contains OpenStack-specific Terraform-variable logic.
package openstack

import (
	"encoding/json"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"
)

type config struct {
	BaseImage       string `json:"openstack_base_image,omitempty"`
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
func TFVars(masterConfig *v1alpha1.OpenstackProviderSpec, cloud string, externalNetwork string, lbFloatingIP string, apiVIP string, dnsVIP string, ingressVIP string, trunkSupport string, octaviaSupport string) ([]byte, error) {
	cfg := &config{
		BaseImage:       masterConfig.Image,
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
	if masterConfig.RootVolume != nil {
		cfg.RootVolumeSize = masterConfig.RootVolume.Size
		cfg.RootVolumeType = masterConfig.RootVolume.VolumeType
	}

	return json.MarshalIndent(cfg, "", "  ")
}

package vmworkflow

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func VirtualMachineOvfDeploySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"local_ovf_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The absolute path to the ovf/ova file in the local system.",
			ForceNew:    true,
		},
		"remote_ovf_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "URL to the remote ovf/ova file to be deployed.",
			ForceNew:    true,
		},
		"ip_allocation_policy": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The IP allocation policy.",
			ForceNew:    true,
		},
		"ip_protocol": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The IP protocol.",
			ForceNew:    true,
		},
		"disk_provisioning": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An optional disk provisioning. If set, all the disks in the deployed ovf will have the same specified disk type (e.g., thin provisioned).",
			ForceNew:    true,
		},
		"deployment_option": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The Deployment option to be chosen. If empty, the default option is used.",
			ForceNew:    true,
		},
		"ovf_network_map": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The mapping of name of network identifiers from the ovf descriptor to network UUID in the VI infrastructure.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			ForceNew:    true,
		},
		"allow_unverified_ssl_cert": {
			Type:        schema.TypeBool,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("VSPHERE_ALLOW_UNVERIFIED_SSL", false),
			Description: "Allow unverified ssl certificates while deploying ovf/ova from url.",
		},
		"enable_hidden_properties": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Allow properties with ovf:userConfigurable=false to be set.",
			ForceNew:    true,
		},
	}
}

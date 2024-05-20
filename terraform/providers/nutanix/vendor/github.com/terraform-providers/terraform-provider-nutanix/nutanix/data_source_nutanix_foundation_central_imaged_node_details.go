package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fc "github.com/terraform-providers/terraform-provider-nutanix/client/fc"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceFCImagedNodeDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFCImagedNodeDetailsRead,
		Schema: map[string]*schema.Schema{
			"imaged_node_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aos_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_key_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"block_serial": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cvm_gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cvm_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cvm_ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cvm_netmask": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cvm_up": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cvm_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cvm_vlan_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"foundation_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hardware_attributes": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"hypervisor_gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor_netmask": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"imaged_cluster_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipmi_gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipmi_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipmi_netmask": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_interface": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_position": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_serial": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_version": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"latest_hb_ts_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceFCImagedNodeDetailsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).FoundationCentral
	req := fc.ImagedNodeDetailsInput{}

	nodeUUID, ok := d.GetOk("imaged_node_uuid")
	if !ok {
		return diag.Errorf("please provide the imaged_node_uuid")
	}
	req.ImagedNodeUUID = nodeUUID.(string)

	res, err := conn.Service.GetImagedNode(ctx, req.ImagedNodeUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("aos_version", res.AosVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("api_key_uuid", res.APIKeyUUID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("available", res.Available); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("block_serial", res.BlockSerial); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_timestamp", res.CreatedTimestamp); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("current_time", res.CurrentTime); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cvm_gateway", res.CvmGateway); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cvm_ip", res.CvmIP); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cvm_ipv6", res.CvmIpv6); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cvm_netmask", res.CvmNetmask); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cvm_up", res.CvmUp); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cvm_uuid", res.CvmUUID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cvm_vlan_id", res.CvmVlanID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("foundation_version", res.FoundationVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hypervisor_gateway", res.HypervisorGateway); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hypervisor_hostname", res.HypervisorHostname); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hypervisor_ip", res.HypervisorIP); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hypervisor_netmask", res.HypervisorNetmask); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hypervisor_type", res.HypervisorType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hypervisor_version", res.HypervisorVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("imaged_cluster_uuid", res.ImagedClusterUUID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ipmi_gateway", res.IpmiGateway); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ipmi_ip", res.IpmiIP); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ipmi_netmask", res.IpmiNetmask); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ipv6_interface", res.Ipv6Interface); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("model", res.Model); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_position", res.NodePosition); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_serial", res.NodeSerial); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_state", res.NodeState); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_type", res.NodeType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("object_version", res.ObjectVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(("imaged_node_uuid"), res.ImagedNodeUUID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(("latest_hb_ts_list"), utils.StringValueSlice(res.LatestHbTSList)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(("hardware_attributes"), flattenHardwareAttributes(res.HardwareAttributes)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

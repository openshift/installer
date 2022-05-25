package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fc "github.com/terraform-providers/terraform-provider-nutanix/client/fc"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixFCImagedNodesList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixFCImagedNodesListRead,
		Schema: map[string]*schema.Schema{
			"length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_state": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"imaged_nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"imaged_node_uuid": {
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
						"supported_features": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"latest_hb_ts_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_matches": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"length": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"offset": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixFCImagedNodesListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).FoundationCentral
	req := fc.ImagedNodesListInput{}

	if len, lenok := d.GetOk("length"); lenok {
		req.Length = utils.IntPtr(len.(int))
	}
	if offset, offok := d.GetOk("offset"); offok {
		req.Offset = utils.IntPtr(offset.(int))
	}

	if filter, fok := d.GetOk("filters"); fok {
		filt := &fc.ImagedNodeListFilter{}
		filter := filter.([]interface{})[0].(map[string]interface{})
		filt.NodeState = utils.StringPtr(filter["node_state"].(string))
		req.Filters = filt
	}

	resp, err := conn.Service.ListImagedNodes(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	imgNodes := flattenImagedNodes(resp.ImagedNodes)

	d.Set("imaged_nodes", imgNodes)

	if resp.Metadata != nil {
		metalist := make([]map[string]interface{}, 0)
		meta := make(map[string]interface{})
		meta["length"] = (resp.Metadata.Length)
		meta["offset"] = (resp.Metadata.Offset)
		meta["total_matches"] = (resp.Metadata.TotalMatches)

		metalist = append(metalist, meta)
		d.Set("metadata", metalist)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenImagedNodes(imgcls []*fc.ImagedNodeDetails) []map[string]interface{} {
	imgClsList := make([]map[string]interface{}, len(imgcls))
	if len(imgcls) > 0 {
		for i, v := range imgcls {
			imgClsList[i] = map[string]interface{}{
				"aos_version":         v.AosVersion,
				"api_key_uuid":        v.APIKeyUUID,
				"available":           v.Available,
				"block_serial":        v.BlockSerial,
				"created_timestamp":   v.CreatedTimestamp,
				"current_time":        v.CurrentTime,
				"cvm_gateway":         v.CvmGateway,
				"cvm_ip":              v.CvmIP,
				"cvm_ipv6":            v.CvmIpv6,
				"cvm_netmask":         v.CvmNetmask,
				"cvm_up":              v.CvmUp,
				"cvm_uuid":            v.CvmUUID,
				"cvm_vlan_id":         v.CvmVlanID,
				"foundation_version":  v.FoundationVersion,
				"hypervisor_gateway":  v.HypervisorGateway,
				"hypervisor_hostname": v.HypervisorHostname,
				"hypervisor_ip":       v.HypervisorIP,
				"hypervisor_netmask":  v.HypervisorNetmask,
				"hypervisor_type":     v.HypervisorType,
				"hypervisor_version":  v.HypervisorVersion,
				"imaged_cluster_uuid": v.ImagedClusterUUID,
				"imaged_node_uuid":    v.ImagedNodeUUID,
				"ipmi_gateway":        v.IpmiGateway,
				"ipmi_ip":             v.IpmiIP,
				"ipmi_netmask":        v.IpmiNetmask,
				"ipv6_interface":      v.Ipv6Interface,
				"model":               v.Model,
				"node_position":       v.NodePosition,
				"node_serial":         v.NodeSerial,
				"node_state":          v.NodeState,
				"node_type":           v.NodeType,
				"object_version":      v.ObjectVersion,
				"supported_features":  flattenFeature(v.SupportedFeatures),
				"latest_hb_ts_list":   flattenFeature(v.LatestHbTSList),
				"hardware_attributes": flattenHardwareAttributes(v.HardwareAttributes),
			}
		}
	}
	return imgClsList
}

func flattenFeature(pr []*string) []string {
	res := []string{}

	for _, v := range pr {
		res = append(res, *v)
	}
	return res
}

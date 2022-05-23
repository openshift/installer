package nutanix

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/foundation"
)

// dataSourceFoundationDiscoverNodes datasource gets discovered nodes within Ipv6 network of foundation
func dataSourceFoundationDiscoverNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFoundationDiscoverNodesRead,
		Schema: map[string]*schema.Schema{
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"model": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chassis_n": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"block_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"foundation_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipv6_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"node_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"current_network_interface": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"node_position": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hypervisor": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"configured": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"nos_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"current_cvm_vlan_tag": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hypervisor_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"svm_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"model": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"node_serial": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// dataSourceFoundationDiscoverNodesRead performs get operation on /discover_nodes api and sets it to resource data schema appropriately
func dataSourceFoundationDiscoverNodesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// create foundation client connection
	conn := meta.(*Client).FoundationClientAPI

	resp, err := conn.Networking.DiscoverNodes(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	entities := make([]map[string]interface{}, len(*resp))
	for k, v := range *resp {
		entity := make(map[string]interface{})
		entity["model"] = v.Model
		entity["chassis_n"] = v.ChassisN
		entity["block_id"] = v.BlockID

		// construct node details map
		entity["nodes"] = flattenDiscoveredNodes(v.Nodes)

		entities[k] = entity
	}

	if setErr := d.Set("entities", entities); setErr != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

// flattenDiscoveredNodes constructs node details map for each discovered node
func flattenDiscoveredNodes(nodesList []foundation.DiscoveredNode) []map[string]interface{} {
	nodes := make([]map[string]interface{}, len(nodesList))
	for k, v := range nodesList {
		node := make(map[string]interface{})

		node["foundation_version"] = v.FoundationVersion
		node["ipv6_address"] = v.Ipv6Address
		node["node_uuid"] = v.NodeUUID
		node["current_network_interface"] = v.CurrentNetworkInterface
		node["node_position"] = v.NodePosition
		node["hypervisor"] = v.Hypervisor
		node["configured"] = v.Configured
		node["nos_version"] = v.NosVersion
		node["node_serial"] = v.NodeSerial

		//ClusterID is of interface{} type so making sure we only accept integer values
		if val, ok := v.ClusterID.(float64); ok {
			node["cluster_id"] = strconv.FormatInt(int64(val), 10)
		} else {
			node["cluster_id"] = ""
		}

		//CurrentCvmVlanTag is of interface{} type so making sure we only accept integer values
		if val, ok := v.CurrentCvmVlanTag.(float64); ok {
			node["current_cvm_vlan_tag"] = strconv.FormatInt(int64(val), 10)
		} else {
			node["current_cvm_vlan_tag"] = ""
		}

		node["hypervisor_version"] = v.HypervisorVersion
		node["svm_ip"] = v.SvmIP
		node["model"] = v.Model

		nodes[k] = node
	}
	return nodes
}

package nutanix

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/foundation"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// dataSourceNodeNetworkDetails datasource for getting node network details as per ipv6 addresses
func dataSourceNodeNetworkDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodeNetworkDetailsRead,
		Schema: map[string]*schema.Schema{
			"ipv6_addresses": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"timeout": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cvm_gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipmi_netmask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cvm_vlan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hypervisor_hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hypervisor_netmask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cvm_netmask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipmi_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hypervisor_gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cvm_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hypervisor_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipmi_gateway": {
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
	}
}

// dataSourceNodeNetworkDetailsRead will get the node network details and set to schema appropriately
func dataSourceNodeNetworkDetailsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// get foundation client api
	conn := meta.(*Client).FoundationClientAPI

	v, ok := d.GetOk("ipv6_addresses")
	if !ok && len(v.([]interface{})) == 0 {
		return diag.Errorf("please provide the ipv6_addresses")
	}

	// create input struct for api call
	input := new(foundation.NodeNetworkDetailsInput)
	ipv6Addresses := expandStringList(v.([]interface{}))
	for _, val := range ipv6Addresses {
		input.Nodes = append(input.Nodes, foundation.NodeIpv6Input{Ipv6Address: *val})
	}
	tout, ok := d.GetOk("timeout")
	if ok {
		input.Timeout = (tout.(string))
	}

	resp, err := conn.Networking.NodeNetworkDetails(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	// create empty diagnostics struct
	var diags diag.Diagnostics

	// set response data appropriately to resource data
	nodes := make([]map[string]string, len(resp.Nodes))
	for k, v := range resp.Nodes {
		node := make(map[string]string)

		if v.Error != "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Node network details for IPv6 address: %s failed", utils.StringValue(ipv6Addresses[k])),
				Detail:   v.Error,
			})
			continue
		}
		node["cvm_gateway"] = v.CvmGateway
		node["ipmi_netmask"] = v.IpmiNetmask
		node["ipv6_address"] = v.Ipv6Address
		node["cvm_vlan_id"] = v.CvmVlanID
		node["hypervisor_hostname"] = v.HypervisorHostname
		node["hypervisor_netmask"] = v.HypervisorNetmask
		node["cvm_netmask"] = v.CvmNetmask
		node["ipmi_ip"] = v.IpmiIP
		node["hypervisor_gateway"] = v.HypervisorGateway
		node["error"] = v.Error
		node["cvm_ip"] = v.CvmIP
		node["hypervisor_ip"] = v.HypervisorIP
		node["ipmi_gateway"] = v.IpmiGateway
		node["node_serial"] = v.NodeSerial
		nodes[k] = node
	}

	// if any errors found
	if len(diags) > 0 {
		return diags
	}

	if setErr := d.Set("nodes", nodes); setErr != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

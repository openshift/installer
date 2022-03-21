package nutanix

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixSubnetsRead,
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"last_update_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"creation_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"spec_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"spec_hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"categories": categoriesSchema(),
						"owner_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"project_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cluster_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"message_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"reason": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"details": {
										Type:     schema.TypeMap,
										Computed: true,
									},
								},
							},
						},
						"cluster_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix_length": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"subnet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dhcp_server_address": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"fqdn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipv6": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"dhcp_server_address_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ip_config_pool_list_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dhcp_options": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"boot_file_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"domain_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tftp_server_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"dhcp_domain_name_server_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dhcp_domain_search_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vlan_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_function_chain_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sort_order": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"offset": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"length": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"sort_attribute": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API
	req := &v3.DSMetadata{}

	metadata, filtersOk := d.GetOk("metadata")
	if filtersOk {
		req = buildDataSourceListMetadata(metadata.(*schema.Set))
	}

	resp, err := conn.V3.ListAllSubnet(utils.StringValue(req.Filter))
	if err != nil {
		return err
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}

	entities := make([]map[string]interface{}, len(resp.Entities))
	for k, v := range resp.Entities {
		entity := make(map[string]interface{})

		m, c := setRSEntityMetadata(v.Metadata)

		entity["metadata"] = m
		entity["project_reference"] = flattenReferenceValues(v.Metadata.ProjectReference)
		entity["owner_reference"] = flattenReferenceValues(v.Metadata.OwnerReference)
		entity["categories"] = c
		entity["api_version"] = utils.StringValue(v.APIVersion)
		entity["name"] = utils.StringValue(v.Status.Name)
		entity["state"] = utils.StringValue(v.Status.State)
		entity["availability_zone_reference"] = flattenReferenceValues(v.Status.AvailabilityZoneReference)
		entity["cluster_reference"] = flattenReferenceValues(v.Status.ClusterReference)

		dgIP := ""
		sIP := ""
		pl := int64(0)
		port := int64(0)
		dhcpSA := make(map[string]interface{})
		dOptions := make(map[string]interface{})
		ipcpl := make([]string, 0)
		dnsList := make([]string, 0)
		dsList := make([]string, 0)

		if v.Status.Resources.IPConfig != nil {
			dgIP = utils.StringValue(v.Status.Resources.IPConfig.DefaultGatewayIP)
			pl = utils.Int64Value(v.Status.Resources.IPConfig.PrefixLength)
			sIP = utils.StringValue(v.Status.Resources.IPConfig.SubnetIP)

			if v.Status.Resources.IPConfig.DHCPServerAddress != nil {
				dhcpSA["ip"] = utils.StringValue(v.Status.Resources.IPConfig.DHCPServerAddress.IP)
				dhcpSA["fqdn"] = utils.StringValue(v.Status.Resources.IPConfig.DHCPServerAddress.FQDN)
				dhcpSA["ipv6"] = utils.StringValue(v.Status.Resources.IPConfig.DHCPServerAddress.IPV6)
				port = utils.Int64Value(v.Status.Resources.IPConfig.DHCPServerAddress.Port)
			}

			if v.Status.Resources.IPConfig.PoolList != nil {
				pl := v.Status.Resources.IPConfig.PoolList
				poolList := make([]string, len(pl))
				for k, v := range pl {
					poolList[k] = utils.StringValue(v.Range)
				}
				ipcpl = poolList
			}
			if v.Status.Resources.IPConfig.DHCPOptions != nil {
				dOptions["boot_file_name"] = utils.StringValue(v.Status.Resources.IPConfig.DHCPOptions.BootFileName)
				dOptions["domain_name"] = utils.StringValue(v.Status.Resources.IPConfig.DHCPOptions.DomainName)
				dOptions["tftp_server_name"] = utils.StringValue(v.Status.Resources.IPConfig.DHCPOptions.TFTPServerName)

				if v.Status.Resources.IPConfig.DHCPOptions.DomainNameServerList != nil {
					dnsList = utils.StringValueSlice(v.Status.Resources.IPConfig.DHCPOptions.DomainNameServerList)
				}
				if v.Status.Resources.IPConfig.DHCPOptions.DomainSearchList != nil {
					dsList = utils.StringValueSlice(v.Status.Resources.IPConfig.DHCPOptions.DomainSearchList)
				}
			}
		}

		entity["dhcp_server_address"] = dhcpSA
		entity["ip_config_pool_list_ranges"] = ipcpl
		entity["dhcp_options"] = dOptions
		entity["dhcp_domain_name_server_list"] = dnsList
		entity["dhcp_domain_search_list"] = dsList

		entity["api_version"] = utils.StringValue(v.APIVersion)
		entity["name"] = utils.StringValue(v.Status.Name)
		entity["description"] = utils.StringValue(v.Status.Description)
		entity["state"] = utils.StringValue(v.Status.State)
		entity["vswitch_name"] = utils.StringValue(v.Status.Resources.VswitchName)
		entity["subnet_type"] = utils.StringValue(v.Status.Resources.SubnetType)
		entity["default_gateway_ip"] = dgIP
		entity["prefix_length"] = pl
		entity["subnet_ip"] = sIP
		entity["dhcp_server_address_port"] = port
		entity["vlan_id"] = utils.Int64Value(v.Status.Resources.VlanID)
		entity["network_function_chain_reference"] = flattenReferenceValues(v.Status.Resources.NetworkFunctionChainReference)

		entities[k] = entity
	}

	if err := d.Set("entities", entities); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return nil
}

func buildDataSourceListMetadata(set *schema.Set) *v3.DSMetadata {
	filters := v3.DSMetadata{}
	for _, v := range set.List() {
		m := v.(map[string]interface{})

		if m["filter"].(string) != "" {
			filters.Filter = utils.StringPtr(m["filter"].(string))
		}
		if m["kind"].(string) != "" {
			filters.Kind = utils.StringPtr(m["kind"].(string))
		}
		if m["sort_order"].(string) != "" {
			filters.SortOrder = utils.StringPtr(m["sort_order"].(string))
		}
		if m["offset"].(int) != 0 {
			filters.Offset = utils.Int64Ptr(int64(m["offset"].(int)))
		}
		if m["length"].(int) != 0 {
			filters.Length = utils.Int64Ptr(int64(m["length"].(int)))
		}
		if m["sort_attribute"].(string) != "" {
			filters.SortAttribute = utils.StringPtr(m["sort_attribute"].(string))
		}
	}
	return &filters
}

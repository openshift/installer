package nutanix

import (
	"fmt"
	"log"
	"strconv"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNutanixNetworkSecurityRule() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceNutanixNetworkSecurityRuleRead,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceNutanixDatasourceNetworkSecurityRuleResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceDatasourceNetworkSecurityRuleStateUpgradeV0,
				Version: 0,
			},
		},
		Schema: map[string]*schema.Schema{
			"network_security_rule_id": {
				Type:     schema.TypeString,
				Required: true,
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
							Type: schema.TypeString,
						},
						"uuid": {
							Type: schema.TypeString,
						},
						"name": {
							Type: schema.TypeString,
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
							Type: schema.TypeString,
						},
						"uuid": {
							Type: schema.TypeString,
						},
						"name": {
							Type: schema.TypeString,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quarantine_rule_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quarantine_rule_outbound_allow_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tcp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"udp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"filter_kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
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
										Type: schema.TypeString,

										Computed: true,
									},
								},
							},
						},
						"icmp_type_code_list": {
							Type: schema.TypeList,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"quarantine_rule_target_group_default_internal_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quarantine_rule_target_group_peer_specification_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quarantine_rule_target_group_filter_kind_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"quarantine_rule_target_group_filter_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quarantine_rule_target_group_filter_params": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"quarantine_rule_inbound_allow_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tcp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"udp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"filter_kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
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
						"icmp_type_code_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"app_rule_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_rule_outbound_allow_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tcp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"udp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"filter_kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
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
						"icmp_type_code_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"app_rule_target_group_default_internal_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_rule_target_group_peer_specification_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_rule_target_group_filter_kind_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"app_rule_target_group_filter_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_rule_target_group_filter_params": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"app_rule_inbound_allow_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tcp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"udp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"filter_kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
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
						"icmp_type_code_list": {
							Type: schema.TypeList,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"isolation_rule_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isolation_rule_first_entity_filter_kind_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"isolation_rule_first_entity_filter_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isolation_rule_first_entity_filter_params": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"isolation_rule_second_entity_filter_kind_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"isolation_rule_second_entity_filter_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isolation_rule_second_entity_filter_params": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"ad_rule_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ad_rule_outbound_allow_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tcp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"udp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"filter_kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
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
						"icmp_type_code_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"ad_rule_target_group_default_internal_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ad_rule_target_group_peer_specification_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ad_rule_target_group_filter_kind_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ad_rule_target_group_filter_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ad_rule_target_group_filter_params": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"ad_rule_inbound_allow_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tcp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"udp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"filter_kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
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
						"icmp_type_code_list": {
							Type: schema.TypeList,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
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

func dataSourceNutanixNetworkSecurityRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Network Security Rule: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*Client).API

	networkSecurityRuleID, ok := d.GetOk("network_security_rule_id")

	if !ok {
		return fmt.Errorf("please provide the required attribute network_security_rule_id")
	}

	// Make request to the API
	resp, err := conn.V3.GetNetworkSecurityRule(networkSecurityRuleID.(string))

	if err != nil {
		return err
	}

	// set metadata values
	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}

	if err := d.Set("categories", c); err != nil {
		return err
	}

	if err := d.Set("project_reference", flattenReferenceValues(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", flattenReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return err
	}

	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Spec.Name))
	d.Set("description", utils.StringValue(resp.Spec.Description))

	if resp.Status == nil {
		return fmt.Errorf("error reading Status from network security rule %s", networkSecurityRuleID.(string))
	}

	if resp.Status.Resources == nil {
		return fmt.Errorf("error reading Status.Resources from network security rule %s", networkSecurityRuleID.(string))
	}

	rules := resp.Status.Resources

	if rules.AllowIpv6Traffic != nil {
		d.Set("allow_ipv6_traffic", utils.BoolValue(rules.AllowIpv6Traffic))
	}

	if rules.IsPolicyHitlogEnabled != nil {
		d.Set("is_policy_hitlog_enabled", utils.BoolValue(rules.IsPolicyHitlogEnabled))
	}

	if rules.QuarantineRule != nil {
		if err := d.Set("quarantine_rule_action", utils.StringValue(rules.QuarantineRule.Action)); err != nil {
			return err
		}

		if rules.QuarantineRule.OutboundAllowList != nil {
			oal := rules.QuarantineRule.OutboundAllowList
			qroaList := make([]map[string]interface{}, len(oal))
			for k, v := range oal {
				qroaItem := make(map[string]interface{})
				qroaItem["protocol"] = utils.StringValue(v.Protocol)

				if v.IPSubnet != nil {
					qroaItem["ip_subnet"] = utils.StringValue(v.IPSubnet.IP)
					qroaItem["ip_subnet_prefix_length"] = strconv.FormatInt(utils.Int64Value(v.IPSubnet.PrefixLength), 10)
				}

				if v.TCPPortRangeList != nil {
					tcpprl := v.TCPPortRangeList
					tcpprList := make([]map[string]interface{}, len(tcpprl))
					for i, tcp := range tcpprl {
						tcpItem := make(map[string]interface{})
						tcpItem["end_port"] = strconv.FormatInt(utils.Int64Value(tcp.EndPort), 10)
						tcpItem["start_port"] = strconv.FormatInt(utils.Int64Value(tcp.StartPort), 10)
						tcpprList[i] = tcpItem
					}
					qroaItem["tcp_port_range_list"] = tcpprList
				}

				if v.UDPPortRangeList != nil {
					udpprl := v.UDPPortRangeList
					udpprList := make([]map[string]interface{}, len(udpprl))
					for i, udp := range udpprl {
						udpItem := make(map[string]interface{})
						udpItem["end_port"] = strconv.FormatInt(utils.Int64Value(udp.EndPort), 10)
						udpItem["start_port"] = strconv.FormatInt(utils.Int64Value(udp.StartPort), 10)
						udpprList[i] = udpItem
					}
					qroaItem["udp_port_range_list"] = udpprList
				}

				if v.Filter != nil {
					qroaItem["filter_kind_list"] = utils.StringValueSlice(v.Filter.KindList)
					qroaItem["filter_type"] = utils.StringValue(v.Filter.Type)
					qroaItem["filter_params"] = expandFilterParams(v.Filter.Params)
				}

				qroaItem["peer_specification_type"] = utils.StringValue(v.PeerSpecificationType)
				qroaItem["expiration_time"] = utils.StringValue(v.ExpirationTime)

				// set network_function_chain_reference
				if v.NetworkFunctionChainReference != nil {
					nfcr := make(map[string]interface{})
					nfcr["kind"] = utils.StringValue(v.NetworkFunctionChainReference.Kind)
					nfcr["name"] = utils.StringValue(v.NetworkFunctionChainReference.Name)
					nfcr["uuid"] = utils.StringValue(v.NetworkFunctionChainReference.UUID)
					qroaItem["network_function_chain_reference"] = nfcr
				}

				if v.IcmpTypeCodeList != nil {
					icmptcl := v.IcmpTypeCodeList
					icmptcList := make([]map[string]interface{}, len(icmptcl))
					for i, icmp := range icmptcl {
						icmpItem := make(map[string]interface{})
						icmpItem["end_port"] = strconv.FormatInt(utils.Int64Value(icmp.Code), 10)
						icmpItem["start_port"] = strconv.FormatInt(utils.Int64Value(icmp.Type), 10)
						icmptcList[i] = icmpItem
					}
					qroaItem["icmp_type_code_list"] = icmptcList
				}

				qroaList[k] = qroaItem
			}

			// Set quarantine_rule_outbound_allow_list
			if err := d.Set("quarantine_rule_outbound_allow_list", qroaList); err != nil {
				return err
			}
		}

		if rules.QuarantineRule.TargetGroup != nil {
			if err := d.Set("quarantine_rule_target_group_default_internal_policy",
				utils.StringValue(rules.QuarantineRule.TargetGroup.DefaultInternalPolicy)); err != nil {
				return err
			}
			if err := d.Set("quarantine_rule_target_group_peer_specification_type",
				utils.StringValue(rules.QuarantineRule.TargetGroup.PeerSpecificationType)); err != nil {
				return err
			}

			if rules.QuarantineRule.TargetGroup.Filter != nil {
				v := rules.QuarantineRule.TargetGroup
				if v.Filter != nil {
					if err := d.Set("quarantine_rule_target_group_filter_kind_list", utils.StringValueSlice(v.Filter.KindList)); err != nil {
						return err
					}

					if err := d.Set("quarantine_rule_target_group_filter_type", utils.StringValue(v.Filter.Type)); err != nil {
						return err
					}
					if err := d.Set("quarantine_rule_target_group_filter_params", expandFilterParams(v.Filter.Params)); err != nil {
						return err
					}
				}
			}
		}

		if rules.QuarantineRule.InboundAllowList != nil {
			ial := rules.QuarantineRule.InboundAllowList
			qriaList := make([]map[string]interface{}, len(ial))
			for k, v := range ial {
				qriaItem := make(map[string]interface{})
				qriaItem["protocol"] = utils.StringValue(v.Protocol)

				if v.IPSubnet != nil {
					qriaItem["ip_subnet"] = utils.StringValue(v.IPSubnet.IP)
					qriaItem["ip_subnet_prefix_length"] = strconv.FormatInt(utils.Int64Value(v.IPSubnet.PrefixLength), 10)
				}

				if v.TCPPortRangeList != nil {
					tcpprl := v.TCPPortRangeList
					tcpprList := make([]map[string]interface{}, len(tcpprl))
					for i, tcp := range tcpprl {
						tcpItem := make(map[string]interface{})
						tcpItem["end_port"] = strconv.FormatInt(utils.Int64Value(tcp.EndPort), 10)
						tcpItem["start_port"] = strconv.FormatInt(utils.Int64Value(tcp.StartPort), 10)
						tcpprList[i] = tcpItem
					}
					qriaItem["tcp_port_range_list"] = tcpprList
				}

				if v.UDPPortRangeList != nil {
					udpprl := v.UDPPortRangeList
					udpprList := make([]map[string]interface{}, len(udpprl))
					for i, udp := range udpprl {
						udpItem := make(map[string]interface{})
						udpItem["end_port"] = strconv.FormatInt(utils.Int64Value(udp.EndPort), 10)
						udpItem["start_port"] = strconv.FormatInt(utils.Int64Value(udp.StartPort), 10)
						udpprList[i] = udpItem
					}
					qriaItem["udp_port_range_list"] = udpprList
				}

				if v.Filter != nil {
					if v.Filter.KindList != nil {
						fkl := v.Filter.KindList
						fkList := make([]string, len(fkl))
						for i, f := range fkl {
							fkList[i] = utils.StringValue(f)
						}
						qriaItem["filter_kind_list"] = fkList
					}

					qriaItem["filter_type"] = utils.StringValue(v.Filter.Type)
					qriaItem["filter_params"] = expandFilterParams(v.Filter.Params)
				}

				qriaItem["peer_specification_type"] = utils.StringValue(v.PeerSpecificationType)
				qriaItem["expiration_time"] = utils.StringValue(v.ExpirationTime)

				// set network_function_chain_reference
				if v.NetworkFunctionChainReference != nil {
					nfcr := make(map[string]interface{})
					nfcr["kind"] = utils.StringValue(v.NetworkFunctionChainReference.Kind)
					nfcr["name"] = utils.StringValue(v.NetworkFunctionChainReference.Name)
					nfcr["uuid"] = utils.StringValue(v.NetworkFunctionChainReference.UUID)
					qriaItem["network_function_chain_reference"] = nfcr
				}

				if v.IcmpTypeCodeList != nil {
					icmptcl := v.IcmpTypeCodeList
					icmptcList := make([]map[string]interface{}, len(icmptcl))
					for i, icmp := range icmptcl {
						icmpItem := make(map[string]interface{})
						icmpItem["end_port"] = strconv.FormatInt(utils.Int64Value(icmp.Code), 10)
						icmpItem["start_port"] = strconv.FormatInt(utils.Int64Value(icmp.Type), 10)
						icmptcList[i] = icmpItem
					}
					qriaItem["icmp_type_code_list"] = icmptcList
				}

				qriaList[k] = qriaItem
			}

			// Set quarantine_rule_inbound_allow_list
			if err := d.Set("quarantine_rule_inbound_allow_list", qriaList); err != nil {
				return err
			}
		}
	} else {
		if err := d.Set("quarantine_rule_inbound_allow_list", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("quarantine_rule_outbound_allow_list", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("quarantine_rule_target_group_filter_kind_list", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("quarantine_rule_target_group_filter_params", make([]string, 0)); err != nil {
			return err
		}
	}

	if err := flattenNetworkRule("app_rule", rules.AppRule, d); err != nil {
		return err
	}

	if err := flattenNetworkRule("ad_rule", rules.AdRule, d); err != nil {
		return err
	}

	if rules.IsolationRule != nil {
		if err := d.Set("isolation_rule_action", utils.StringValue(rules.IsolationRule.Action)); err != nil {
			return err
		}

		if rules.IsolationRule.FirstEntityFilter != nil {
			firstFilter := rules.IsolationRule.FirstEntityFilter
			if err := d.Set("isolation_rule_first_entity_filter_kind_list", utils.StringValueSlice(firstFilter.KindList)); err != nil {
				return err
			}

			if err := d.Set("isolation_rule_first_entity_filter_type", utils.StringValue(firstFilter.Type)); err != nil {
				return err
			}

			if err := d.Set("isolation_rule_first_entity_filter_params", expandFilterParams(firstFilter.Params)); err != nil {
				return err
			}
		}

		if rules.IsolationRule.SecondEntityFilter != nil {
			secondFilter := rules.IsolationRule.SecondEntityFilter

			if err := d.Set("isolation_rule_second_entity_filter_kind_list", utils.StringValueSlice(secondFilter.KindList)); err != nil {
				return err
			}
			if err := d.Set("isolation_rule_second_entity_filter_type", utils.StringValue(secondFilter.Type)); err != nil {
				return err
			}
			if err := d.Set("isolation_rule_second_entity_filter_params", expandFilterParams(secondFilter.Params)); err != nil {
				return err
			}
		}
	} else {
		if err := d.Set("isolation_rule_first_entity_filter_kind_list", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("isolation_rule_first_entity_filter_params", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("isolation_rule_second_entity_filter_kind_list", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("isolation_rule_second_entity_filter_params", make([]string, 0)); err != nil {
			return err
		}
	}

	d.SetId(utils.StringValue(resp.Metadata.UUID))

	return nil
}

func resourceDatasourceNetworkSecurityRuleStateUpgradeV0(is map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Entering resourceDatasourceNetworkSecurityRuleStateUpgradeV0")
	return resourceNutanixCategoriesMigrateState(is, meta)
}

func resourceNutanixDatasourceNetworkSecurityRuleResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"network_security_rule_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"categories": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"owner_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type: schema.TypeString,
						},
						"uuid": {
							Type: schema.TypeString,
						},
						"name": {
							Type: schema.TypeString,
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
							Type: schema.TypeString,
						},
						"uuid": {
							Type: schema.TypeString,
						},
						"name": {
							Type: schema.TypeString,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_ipv6_traffic": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_policy_hitlog_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"quarantine_rule_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quarantine_rule_outbound_allow_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tcp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"udp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"filter_kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
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
										Type: schema.TypeString,

										Computed: true,
									},
								},
							},
						},
						"icmp_type_code_list": {
							Type: schema.TypeList,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"quarantine_rule_target_group_default_internal_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quarantine_rule_target_group_peer_specification_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quarantine_rule_target_group_filter_kind_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"quarantine_rule_target_group_filter_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quarantine_rule_target_group_filter_params": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"quarantine_rule_inbound_allow_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tcp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"udp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"filter_kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
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
						"icmp_type_code_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"app_rule_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_rule_outbound_allow_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tcp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"udp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"filter_kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
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
						"icmp_type_code_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"app_rule_target_group_default_internal_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_rule_target_group_peer_specification_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_rule_target_group_filter_kind_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"app_rule_target_group_filter_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_rule_target_group_filter_params": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"app_rule_inbound_allow_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tcp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"udp_port_range_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"filter_kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
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
						"icmp_type_code_list": {
							Type: schema.TypeList,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"isolation_rule_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isolation_rule_first_entity_filter_kind_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"isolation_rule_first_entity_filter_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isolation_rule_first_entity_filter_params": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"isolation_rule_second_entity_filter_kind_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"isolation_rule_second_entity_filter_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isolation_rule_second_entity_filter_params": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

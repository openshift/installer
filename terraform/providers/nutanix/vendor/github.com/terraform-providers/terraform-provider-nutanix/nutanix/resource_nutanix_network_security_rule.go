package nutanix

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

var (
	netTimeout    = 10 * time.Minute
	netDelay      = 10 * time.Second
	netMinTimeout = 3 * time.Second
)

func resourceNutanixNetworkSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixNetworkSecurityRuleCreate,
		Read:   resourceNutanixNetworkSecurityRuleRead,
		Update: resourceNutanixNetworkSecurityRuleUpdate,
		Delete: resourceNutanixNetworkSecurityRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceNutanixSecurityRuleInstanceResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceSecurityRuleInstanceStateUpgradeV0,
				Version: 0,
			},
		},
		Schema: map[string]*schema.Schema{
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
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"allow_ipv6_traffic": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_policy_hitlog_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"app_rule_action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"app_rule_outbound_allow_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tcp_port_range_list": portRangeSchema(),
						"udp_port_range_list": portRangeSchema(),
						"filter_kind_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"network_function_chain_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Required: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"icmp_type_code_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
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
				Optional: true,
			},
			"app_rule_target_group_peer_specification_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_rule_target_group_filter_kind_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"app_rule_target_group_filter_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"app_rule_target_group_filter_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"app_rule_inbound_allow_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tcp_port_range_list": portRangeSchema(),
						"udp_port_range_list": portRangeSchema(),
						"filter_kind_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"network_function_chain_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Required: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"icmp_type_code_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
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
				Optional: true,
				Computed: true,
			},
			"isolation_rule_first_entity_filter_kind_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"isolation_rule_first_entity_filter_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"isolation_rule_first_entity_filter_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"isolation_rule_second_entity_filter_kind_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"isolation_rule_second_entity_filter_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"isolation_rule_second_entity_filter_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"ad_rule_action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ad_rule_outbound_allow_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tcp_port_range_list": portRangeSchema(),
						"udp_port_range_list": portRangeSchema(),
						"filter_kind_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"network_function_chain_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Required: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"icmp_type_code_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
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
				Optional: true,
			},
			"ad_rule_target_group_peer_specification_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ad_rule_target_group_filter_kind_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ad_rule_target_group_filter_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ad_rule_target_group_filter_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"ad_rule_inbound_allow_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tcp_port_range_list": portRangeSchema(),
						"udp_port_range_list": portRangeSchema(),
						"filter_kind_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"network_function_chain_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Required: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"icmp_type_code_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
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

func resourceNutanixNetworkSecurityRuleCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	// Prepare request
	request := &v3.NetworkSecurityRuleIntentInput{}
	spec := &v3.NetworkSecurityRule{}
	metadata := &v3.Metadata{}
	networkSecurityRule := &v3.NetworkSecurityRuleResources{}

	// Read arguments and set request values
	name, nok := d.GetOk("name")
	desc, descok := d.GetOk("description")

	if !nok {
		return fmt.Errorf("please provide the required attribute name")
	}

	// Read arguments and set request values

	// only set kind
	if errMetad := getMetadataAttributes(d, metadata, "network_security_rule"); errMetad != nil {
		return errMetad
	}

	if descok {
		spec.Description = utils.StringPtr(desc.(string))
	}

	// get resources
	if err := getNetworkSecurityRuleResources(d, networkSecurityRule); err != nil {
		return err
	}

	if descok {
		spec.Description = utils.StringPtr(desc.(string))
	}

	// set request

	spec.Resources = networkSecurityRule

	spec.Name = utils.StringPtr(name.(string))

	// set request attrs
	request.Metadata = metadata
	request.Spec = spec

	// Make request to API
	resp, err := conn.V3.CreateNetworkSecurityRule(request)

	if err != nil {
		return fmt.Errorf("error creating Nutanix Network Security Rule %s: %+v", utils.StringValue(spec.Name), err)
	}

	d.SetId(*resp.Metadata.UUID)

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    netTimeout,
		Delay:      netDelay,
		MinTimeout: netMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for network_security_rule (%s) to create: %s", d.Id(), err)
	}

	return resourceNutanixNetworkSecurityRuleRead(d, meta)
}

func resourceNutanixNetworkSecurityRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Network Security Rule: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*Client).API

	// Make request to the API
	resp, errNet := conn.V3.GetNetworkSecurityRule(d.Id())
	if errNet != nil {
		if strings.Contains(fmt.Sprint(errNet), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return errNet
	}

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

	if resp.Status == nil {
		return fmt.Errorf("error reading Status from network security rule %s", d.Id())
	}

	if resp.Status.Resources == nil {
		return fmt.Errorf("error reading Status.Resources from network security rule %s", d.Id())
	}

	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Status.Name))
	d.Set("description", utils.StringValue(resp.Status.Description))

	rules := resp.Spec.Resources

	if rules.AllowIpv6Traffic != nil {
		d.Set("allow_ipv6_traffic", utils.BoolValue(rules.AllowIpv6Traffic))
	}

	if rules.IsPolicyHitlogEnabled != nil {
		d.Set("is_policy_hitlog_enabled", utils.BoolValue(rules.IsPolicyHitlogEnabled))
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

	return nil
}

func flattenNetworkRule(prefix string, rule *v3.NetworkSecurityRuleResourcesRule, d *schema.ResourceData) error {
	if rule != nil {
		if err := d.Set(fmt.Sprintf("%s_action", prefix), utils.StringValue(rule.Action)); err != nil {
			return err
		}

		if rule.TargetGroup != nil {
			if err := d.Set(fmt.Sprintf("%s_target_group_default_internal_policy", prefix),
				utils.StringValue(rule.TargetGroup.DefaultInternalPolicy)); err != nil {
				return err
			}
			if err := d.Set(fmt.Sprintf("%s_target_group_peer_specification_type", prefix),
				utils.StringValue(rule.TargetGroup.PeerSpecificationType)); err != nil {
				return err
			}

			if rule.TargetGroup.Filter != nil {
				v := rule.TargetGroup
				if v.Filter != nil {
					if err := d.Set(fmt.Sprintf("%s_target_group_filter_kind_list", prefix), utils.StringValueSlice(v.Filter.KindList)); err != nil {
						return err
					}

					if err := d.Set(fmt.Sprintf("%s_target_group_filter_type", prefix), utils.StringValue(v.Filter.Type)); err != nil {
						return err
					}

					if err := d.Set(fmt.Sprintf("%s_target_group_filter_params", prefix), expandFilterParams(v.Filter.Params)); err != nil {
						return err
					}
				}
			}
		}

		// Set app_rule_outbound_allow_list
		if err := d.Set(fmt.Sprintf("%s_outbound_allow_list", prefix), flattenNetworkRuleList(rule.OutboundAllowList)); err != nil {
			return err
		}

		// Set app_rule_inbound_allow_list
		if err := d.Set(fmt.Sprintf("%s_inbound_allow_list", prefix), flattenNetworkRuleList(rule.InboundAllowList)); err != nil {
			return err
		}
	} else if err := d.Set(fmt.Sprintf("%s_target_group_filter_kind_list", prefix), make([]string, 0)); err != nil {
		return err
	}
	return nil
}

func resourceNutanixNetworkSecurityRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	// Prepare request
	request := &v3.NetworkSecurityRuleIntentInput{}
	spec := &v3.NetworkSecurityRule{}
	metadata := &v3.Metadata{}
	networkSecurityRule := &v3.NetworkSecurityRuleResources{}

	response, err := conn.V3.GetNetworkSecurityRule(d.Id())

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}
		return err
	}

	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if response.Spec != nil {
		spec = response.Spec

		if response.Spec.Resources != nil {
			networkSecurityRule = response.Spec.Resources
		}
	}

	if d.HasChange("categories") {
		metadata.Categories = expandCategories(d.Get("categories"))
	}

	if d.HasChange("owner_reference") {
		or := d.Get("owner_reference").(map[string]interface{})
		metadata.OwnerReference = validateRef(or)
	}

	if d.HasChange("project_reference") {
		pr := d.Get("project_reference").(map[string]interface{})
		metadata.ProjectReference = validateRef(pr)
	}

	if d.HasChange("name") {
		spec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		spec.Description = utils.StringPtr(d.Get("description").(string))
	}

	// TODO: Change
	if d.HasChange("allow_ipv6_traffic") ||
		d.HasChange("is_policy_hitlog_enabled") ||
		d.HasChange("app_rule_action") ||
		d.HasChange("app_rule_outbound_allow_list") ||
		d.HasChange("app_rule_target_group_default_internal_policy") ||
		d.HasChange("app_rule_target_group_peer_specification_type") ||
		d.HasChange("app_rule_target_group_filter_kind_list") ||
		d.HasChange("app_rule_target_group_filter_type") ||
		d.HasChange("app_rule_target_group_filter_params") ||
		d.HasChange("app_rule_inbound_allow_list") ||
		d.HasChange("ad_rule_action") ||
		d.HasChange("ad_rule_outbound_allow_list") ||
		d.HasChange("ad_rule_target_group_default_internal_policy") ||
		d.HasChange("ad_rule_target_group_peer_specification_type") ||
		d.HasChange("ad_rule_target_group_filter_kind_list") ||
		d.HasChange("ad_rule_target_group_filter_type") ||
		d.HasChange("ad_rule_target_group_filter_params") ||
		d.HasChange("ad_rule_inbound_allow_list") ||
		d.HasChange("isolation_rule_action") ||
		d.HasChange("isolation_rule_first_entity_filter_kind_list") ||
		d.HasChange("isolation_rule_first_entity_filter_type") ||
		d.HasChange("isolation_rule_first_entity_filter_params") ||
		d.HasChange("isolation_rule_second_entity_filter_kind_list") ||
		d.HasChange("isolation_rule_second_entity_filter_type") ||
		d.HasChange("isolation_rule_second_entity_filter_params") {
		if err := getNetworkSecurityRuleResources(d, networkSecurityRule); err != nil {
			return err
		}
		spec.Resources = networkSecurityRule
	}

	request.Spec = spec
	request.Metadata = metadata

	resp, errUpdate := conn.V3.UpdateNetworkSecurityRule(d.Id(), request)

	if errUpdate != nil {
		return errUpdate
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    netTimeout,
		Delay:      netDelay,
		MinTimeout: netMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for network_security_rule (%s) to update: %s", d.Id(), err)
	}

	return resourceNutanixNetworkSecurityRuleRead(d, meta)
}

func resourceNutanixNetworkSecurityRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting Network Security Rule: %s", d.Get("name").(string))

	conn := meta.(*Client).API
	UUID := d.Id()

	resp, err := conn.V3.DeleteNetworkSecurityRule(UUID)
	if err != nil {
		return err
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    netTimeout,
		Delay:      netDelay,
		MinTimeout: netMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for network_security_rule (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func getNetworkSecurityRuleResources(d *schema.ResourceData, networkSecurityRule *v3.NetworkSecurityRuleResources) error {
	isolationRule := &v3.NetworkSecurityRuleIsolationRule{}

	iRuleFirstEntityFilter := &v3.CategoryFilter{}
	iRuleSecondEntityFilter := &v3.CategoryFilter{}

	if allowIpv6Traffic, aok := d.GetOk("allow_ipv6_traffic"); aok {
		networkSecurityRule.AllowIpv6Traffic = utils.BoolPtr(allowIpv6Traffic.(bool))
	}

	if isPolicyHitlogEnabled, iok := d.GetOk("is_policy_hitlog_enabled"); iok {
		networkSecurityRule.IsPolicyHitlogEnabled = utils.BoolPtr(isPolicyHitlogEnabled.(bool))
	}

	appRule := expandNetworkRule("app_rule", d)
	adRule := expandNetworkRule("ad_rule", d)

	if ira, ok := d.GetOk("isolation_rule_action"); ok && ira.(string) != "" {
		isolationRule.Action = utils.StringPtr(ira.(string))
	}

	if f, fok := d.GetOk("isolation_rule_first_entity_filter_kind_list"); fok && f != nil {
		iRuleFirstEntityFilter.KindList = expandStringList(f.([]interface{}))
	}

	if ft, ftok := d.GetOk("isolation_rule_first_entity_filter_type"); ftok && ft.(string) != "" {
		iRuleFirstEntityFilter.Type = utils.StringPtr(ft.(string))
	}

	if fp, fpok := d.GetOk("isolation_rule_first_entity_filter_params"); fpok {
		fpl := fp.(*schema.Set).List()

		if len(fpl) > 0 {
			fl := make(map[string][]string)
			for _, v := range fpl {
				item := v.(map[string]interface{})

				if i, ok := item["name"]; ok && i.(string) != "" {
					if k, kok := item["values"]; kok && len(k.([]interface{})) > 0 {
						var values []string
						for _, item := range k.([]interface{}) {
							values = append(values, item.(string))
						}
						fl[i.(string)] = values
					}
				}
			}
			iRuleFirstEntityFilter.Params = fl
		} else {
			iRuleFirstEntityFilter.Params = nil
		}
	}

	if f, fok := d.GetOk("isolation_rule_second_entity_filter_kind_list"); fok && f != nil {
		iRuleSecondEntityFilter.KindList = expandStringList(f.([]interface{}))
	}

	if ft, ftok := d.GetOk("isolation_rule_second_entity_filter_type"); ftok && ft.(string) != "" {
		iRuleSecondEntityFilter.Type = utils.StringPtr(ft.(string))
	}

	if fp, fpok := d.GetOk("isolation_rule_second_entity_filter_params"); fpok {
		fpl := fp.(*schema.Set).List()

		if len(fpl) > 0 {
			fl := make(map[string][]string)
			for _, v := range fpl {
				item := v.(map[string]interface{})

				if i, ok := item["name"]; ok && i.(string) != "" {
					if k, kok := item["values"]; kok && len(k.([]interface{})) > 0 {
						var values []string
						for _, item := range k.([]interface{}) {
							values = append(values, item.(string))
						}
						fl[i.(string)] = values
					}
				}
			}
			iRuleSecondEntityFilter.Params = fl
		} else {
			iRuleSecondEntityFilter.Params = nil
		}
	}

	if !reflect.DeepEqual(*appRule, (v3.NetworkSecurityRuleResourcesRule{})) {
		networkSecurityRule.AppRule = appRule
	}

	if !reflect.DeepEqual(*adRule, (v3.NetworkSecurityRuleResourcesRule{})) {
		networkSecurityRule.AdRule = adRule
	}

	if !reflect.DeepEqual(*iRuleFirstEntityFilter, (v3.CategoryFilter{})) {
		isolationRule.FirstEntityFilter = iRuleFirstEntityFilter
	}

	if !reflect.DeepEqual(*iRuleSecondEntityFilter, (v3.CategoryFilter{})) {
		isolationRule.SecondEntityFilter = iRuleSecondEntityFilter
	}

	if !reflect.DeepEqual(*isolationRule, (v3.NetworkSecurityRuleIsolationRule{})) {
		networkSecurityRule.IsolationRule = isolationRule
	}
	return nil
}

func expandNetworkRule(prefix string, d *schema.ResourceData) *v3.NetworkSecurityRuleResourcesRule {
	appRule := &v3.NetworkSecurityRuleResourcesRule{}
	aRuleTargetGroup := &v3.TargetGroup{}
	aRuleTargetGroupFilter := &v3.CategoryFilter{}

	if ara, ok := d.GetOk(fmt.Sprintf("%s_action", prefix)); ok && ara.(string) != "" {
		appRule.Action = utils.StringPtr(ara.(string))
	}

	if qroal, ok := d.GetOk(fmt.Sprintf("%s_outbound_allow_list", prefix)); ok {
		oal := qroal.([]interface{})
		outbound := make([]*v3.NetworkRule, len(oal))

		for k, v := range oal {
			nr := v.(map[string]interface{})
			nrItem := &v3.NetworkRule{}
			iPSubnet := &v3.IPSubnet{}
			filter := &v3.CategoryFilter{}

			if proto, pok := nr["protocol"]; pok && proto.(string) != "" {
				nrItem.Protocol = utils.StringPtr(proto.(string))
			}

			if ip, ipok := nr["ip_subnet"]; ipok && ip.(string) != "" {
				iPSubnet.IP = utils.StringPtr(ip.(string))
			}

			if ippl, ipok := nr["ip_subnet_prefix_length"]; ipok && ippl.(string) != "" {
				if i, err := strconv.Atoi(ippl.(string)); err == nil {
					iPSubnet.PrefixLength = utils.Int64Ptr(int64(i))
				}
			}

			if t, tcpok := nr["tcp_port_range_list"]; tcpok {
				nrItem.TCPPortRangeList = expandPortRangeList(t)
			}

			if u, udpok := nr["udp_port_range_list"]; udpok {
				nrItem.UDPPortRangeList = expandPortRangeList(u)
			}

			if f, fok := nr["filter_kind_list"]; fok && len(f.([]interface{})) > 0 {
				filter.KindList = expandStringList(f.([]interface{}))
			}

			if ft, ftok := nr["filter_type"]; ftok && ft != "" {
				filter.Type = utils.StringPtr(ft.(string))
			}

			if fp, fpok := nr["filter_params"]; fpok {
				fpl := fp.(*schema.Set).List()

				if len(fpl) > 0 {
					fl := make(map[string][]string)
					for _, v := range fpl {
						item := v.(map[string]interface{})

						if i, ok := item["name"]; ok && i.(string) != "" {
							if k2, kok := item["values"]; kok && len(k2.([]interface{})) > 0 {
								var values []string
								for _, item := range k2.([]interface{}) {
									values = append(values, item.(string))
								}
								fl[i.(string)] = values
							}
						}
					}
					filter.Params = fl
				} else {
					filter.Params = nil
				}
			}

			if pet, petok := nr["peer_specification_type"]; petok && pet.(string) != "" {
				nrItem.PeerSpecificationType = utils.StringPtr(pet.(string))
			}

			if et, etok := nr["expiration_time"]; etok && et.(string) != "" {
				nrItem.ExpirationTime = utils.StringPtr(et.(string))
			}

			if nfcr, nfcrok := nr["network_function_chain_reference"]; nfcrok && len(nfcr.(map[string]interface{})) > 0 {
				a := nfcr.(map[string]interface{})
				nrItem.NetworkFunctionChainReference = validateRef(a)
			}

			if icmp, icmpok := nr["icmp_type_code_list"]; icmpok {
				nrItem.IcmpTypeCodeList = expandIcmpTypeCodeList(icmp)
			}

			nrItem.IPSubnet = iPSubnet
			if !reflect.DeepEqual(*filter, v3.CategoryFilter{}) {
				nrItem.Filter = filter
			}
			outbound[k] = nrItem
		}
		appRule.OutboundAllowList = outbound
	}

	if qroal, ok := d.GetOk(fmt.Sprintf("%s_target_group_default_internal_policy", prefix)); ok && qroal != nil {
		aRuleTargetGroup.DefaultInternalPolicy = utils.StringPtr(qroal.(string))
	}

	if qroal, ok := d.GetOk(fmt.Sprintf("%s_target_group_peer_specification_type", prefix)); ok && qroal != nil {
		aRuleTargetGroup.PeerSpecificationType = utils.StringPtr(qroal.(string))
	}

	if f, fok := d.GetOk(fmt.Sprintf("%s_target_group_filter_kind_list", prefix)); fok && f != nil {
		aRuleTargetGroupFilter.KindList = expandStringList(f.([]interface{}))
	}

	if ft, ftok := d.GetOk(fmt.Sprintf("%s_target_group_filter_type", prefix)); ftok && ft.(string) != "" {
		aRuleTargetGroupFilter.Type = utils.StringPtr(ft.(string))
	}

	if fp, fpok := d.GetOk(fmt.Sprintf("%s_target_group_filter_params", prefix)); fpok {
		fpl := fp.(*schema.Set).List()

		if len(fpl) > 0 {
			fl := make(map[string][]string)
			for _, v := range fpl {
				item := v.(map[string]interface{})

				if i, ok := item["name"]; ok && i.(string) != "" {
					if k, kok := item["values"]; kok && len(k.([]interface{})) > 0 {
						var values []string
						for _, item := range k.([]interface{}) {
							values = append(values, item.(string))
						}
						fl[i.(string)] = values
					}
				}
			}
			aRuleTargetGroupFilter.Params = fl
		} else {
			aRuleTargetGroupFilter.Params = nil
		}
	}

	if qrial, ok := d.GetOk(fmt.Sprintf("%s_inbound_allow_list", prefix)); ok {
		oal := qrial.([]interface{})
		inbound := make([]*v3.NetworkRule, len(oal))

		for k, v := range oal {
			nr := v.(map[string]interface{})
			nrItem := &v3.NetworkRule{}
			iPSubnet := &v3.IPSubnet{}
			filter := &v3.CategoryFilter{}

			if proto, pok := nr["protocol"]; pok && proto.(string) != "" {
				nrItem.Protocol = utils.StringPtr(proto.(string))
			}

			if ip, ipok := nr["ip_subnet"]; ipok && ip.(string) != "" {
				iPSubnet.IP = utils.StringPtr(ip.(string))
			}

			if ippl, ipok := nr["ip_subnet_prefix_length"]; ipok && ippl.(string) != "" {
				if i, err := strconv.Atoi(ippl.(string)); err == nil {
					iPSubnet.PrefixLength = utils.Int64Ptr(int64(i))
				}
			}

			if t, tcpok := nr["tcp_port_range_list"]; tcpok {
				nrItem.TCPPortRangeList = expandPortRangeList(t)
			}

			if u, udpok := nr["udp_port_range_list"]; udpok {
				nrItem.UDPPortRangeList = expandPortRangeList(u)
			}

			if f, fok := nr["filter_kind_list"]; fok && len(f.([]interface{})) > 0 {
				filter.KindList = expandStringList(f.([]interface{}))
			}

			if ft, ftok := nr["filter_type"]; ftok && ft != "" {
				filter.Type = utils.StringPtr(ft.(string))
			}

			if fp, fpok := nr["filter_params"]; fpok {
				fpl := fp.(*schema.Set).List()

				if len(fpl) > 0 {
					fl := make(map[string][]string)
					for _, v := range fpl {
						item := v.(map[string]interface{})

						if i, ok := item["name"]; ok && i.(string) != "" {
							if k2, kok := item["values"]; kok && len(k2.([]interface{})) > 0 {
								var values []string
								for _, item := range k2.([]interface{}) {
									values = append(values, item.(string))
								}
								fl[i.(string)] = values
							}
						}
					}
					filter.Params = fl
				} else {
					filter.Params = nil
				}
			}

			if pet, petok := nr["peer_specification_type"]; petok && pet.(string) != "" {
				nrItem.PeerSpecificationType = utils.StringPtr(pet.(string))
			}

			if et, etok := nr["expiration_time"]; etok && et.(string) != "" {
				nrItem.ExpirationTime = utils.StringPtr(et.(string))
			}

			if nfcr, nfcrok := nr["network_function_chain_reference"]; nfcrok && len(nfcr.(map[string]interface{})) > 0 {
				a := nfcr.(map[string]interface{})
				nrItem.NetworkFunctionChainReference = validateRef(a)
			}

			if icmp, icmpok := nr["icmp_type_code_list"]; icmpok {
				nrItem.IcmpTypeCodeList = expandIcmpTypeCodeList(icmp)
			}

			nrItem.IPSubnet = iPSubnet
			if !reflect.DeepEqual(*filter, v3.CategoryFilter{}) {
				nrItem.Filter = filter
			}
			inbound[k] = nrItem
		}
		appRule.InboundAllowList = inbound
	}

	if !reflect.DeepEqual(*aRuleTargetGroupFilter, (v3.CategoryFilter{})) {
		aRuleTargetGroup.Filter = aRuleTargetGroupFilter
	}

	if !reflect.DeepEqual(*aRuleTargetGroup, (v3.TargetGroup{})) {
		appRule.TargetGroup = aRuleTargetGroup
	}

	return appRule
}

func expandFilterParams(fp map[string][]string) []map[string]interface{} {
	fpList := make([]map[string]interface{}, 0)
	if len(fp) > 0 {
		for name, values := range fp {
			fpItem := make(map[string]interface{})
			fpItem["name"] = name
			fpItem["values"] = values
			fpList = append(fpList, fpItem)
		}
	}
	return fpList
}

func expandPortRangeList(pr interface{}) []*v3.PortRange {
	portRange := pr.([]interface{})
	ports := make([]*v3.PortRange, len(portRange))

	for i, p := range portRange {
		port := p.(map[string]interface{})
		portRange := &v3.PortRange{}

		if endp, epok := port["end_port"]; epok {
			portRange.EndPort = utils.Int64Ptr(int64(endp.(int)))
		}

		if stp, stpok := port["start_port"]; stpok {
			portRange.StartPort = utils.Int64Ptr(int64(stp.(int)))
		}
		ports[i] = portRange
	}
	return ports
}

func expandIcmpTypeCodeList(icmp interface{}) []*v3.NetworkRuleIcmpTypeCodeList {
	ic := icmp.([]interface{})
	icmpList := make([]*v3.NetworkRuleIcmpTypeCodeList, 0)
	for _, v := range ic {
		icmpm := v.(map[string]interface{})
		icmpItem := &v3.NetworkRuleIcmpTypeCodeList{}

		if c, cok := icmpm["code"]; cok && c.(string) != "" {
			if i, err := strconv.ParseInt(c.(string), 10, 64); err == nil {
				icmpItem.Code = utils.Int64Ptr(i)
			}
		}

		if t, tok := icmpm["type"]; tok && t.(string) != "" {
			if i, err := strconv.Atoi(t.(string)); err == nil {
				icmpItem.Type = utils.Int64Ptr(int64(i))
			}
		}
		icmpList = append(icmpList, icmpItem)
	}
	return icmpList
}

func filterParamsHash(v interface{}) int {
	params := v.(map[string]interface{})
	return hashcode.String(params["name"].(string))
}

func flattenNetworkRuleList(networkRules []*v3.NetworkRule) []map[string]interface{} {
	ruleList := make([]map[string]interface{}, 0)
	for _, v := range networkRules {
		ruleItem := make(map[string]interface{})
		ruleItem["protocol"] = utils.StringValue(v.Protocol)

		if v.IPSubnet != nil {
			ruleItem["ip_subnet"] = utils.StringValue(v.IPSubnet.IP)
			ruleItem["ip_subnet_prefix_length"] = strconv.FormatInt(utils.Int64Value(v.IPSubnet.PrefixLength), 10)
		}

		if v.TCPPortRangeList != nil {
			tcpprl := v.TCPPortRangeList
			tcpprList := make([]map[string]interface{}, len(tcpprl))
			for i, tcp := range tcpprl {
				tcpItem := make(map[string]interface{})
				tcpItem["end_port"] = utils.Int64Value(tcp.EndPort)
				tcpItem["start_port"] = utils.Int64Value(tcp.StartPort)
				tcpprList[i] = tcpItem
			}
			ruleItem["tcp_port_range_list"] = tcpprList
		}

		if v.UDPPortRangeList != nil {
			udpprl := v.UDPPortRangeList
			udpprList := make([]map[string]interface{}, len(udpprl))
			for i, udp := range udpprl {
				udpItem := make(map[string]interface{})
				udpItem["end_port"] = utils.Int64Value(udp.EndPort)
				udpItem["start_port"] = utils.Int64Value(udp.StartPort)
				udpprList[i] = udpItem
			}
			ruleItem["udp_port_range_list"] = udpprList
		}

		if v.Filter != nil {
			ruleItem["filter_kind_list"] = utils.StringValueSlice(v.Filter.KindList)
			ruleItem["filter_type"] = utils.StringValue(v.Filter.Type)
			ruleItem["filter_params"] = expandFilterParams(v.Filter.Params)
		}

		ruleItem["peer_specification_type"] = utils.StringValue(v.PeerSpecificationType)
		ruleItem["expiration_time"] = utils.StringValue(v.ExpirationTime)

		// set network_function_chain_reference
		if v.NetworkFunctionChainReference != nil {
			nfcr := make(map[string]interface{})
			nfcr["kind"] = utils.StringValue(v.NetworkFunctionChainReference.Kind)
			nfcr["name"] = utils.StringValue(v.NetworkFunctionChainReference.Name)
			nfcr["uuid"] = utils.StringValue(v.NetworkFunctionChainReference.UUID)
			ruleItem["network_function_chain_reference"] = nfcr
		}

		if v.IcmpTypeCodeList != nil {
			icmptcl := v.IcmpTypeCodeList
			icmptcList := make([]map[string]interface{}, len(icmptcl))
			for i, icmp := range icmptcl {
				icmpItem := make(map[string]interface{})
				icmpItem["code"] = strconv.FormatInt(utils.Int64Value(icmp.Code), 10)
				icmpItem["type"] = strconv.FormatInt(utils.Int64Value(icmp.Type), 10)
				icmptcList[i] = icmpItem
			}
			ruleItem["icmp_type_code_list"] = icmptcList
		}
		ruleList = append(ruleList, ruleItem)
	}
	return ruleList
}

func portRangeSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"end_port": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"start_port": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func resourceSecurityRuleInstanceStateUpgradeV0(is map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Entering resourceSecurityRuleInstanceStateUpgradeV0")
	return resourceNutanixCategoriesMigrateState(is, meta)
}

func resourceNutanixSecurityRuleInstanceResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"app_rule_action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"app_rule_outbound_allow_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tcp_port_range_list": portRangeSchema(),
						"udp_port_range_list": portRangeSchema(),
						"filter_kind_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"network_function_chain_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Required: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"icmp_type_code_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
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
				Optional: true,
			},
			"app_rule_target_group_peer_specification_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_rule_target_group_filter_kind_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"app_rule_target_group_filter_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"app_rule_target_group_filter_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"app_rule_inbound_allow_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_subnet": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_subnet_prefix_length": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tcp_port_range_list": portRangeSchema(),
						"udp_port_range_list": portRangeSchema(),
						"filter_kind_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"filter_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"filter_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"peer_specification_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"expiration_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"network_function_chain_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Required: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"icmp_type_code_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
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
				Optional: true,
				Computed: true,
			},
			"isolation_rule_first_entity_filter_kind_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"isolation_rule_first_entity_filter_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"isolation_rule_first_entity_filter_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"isolation_rule_second_entity_filter_kind_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"isolation_rule_second_entity_filter_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"isolation_rule_second_entity_filter_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      filterParamsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

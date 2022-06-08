package nutanix

import (
	"context"
	"fmt"
	"strings"

	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNutanixProtectionRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixProtectionRuleRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"protection_rule_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"protection_rule_name"},
			},
			"protection_rule_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"protection_rule_id"},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"categories": categoriesSchema(),
			"owner_reference": {
				Type:     schema.TypeList,
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
				Type:     schema.TypeList,
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
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_connectivity_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_availability_zone_index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source_availability_zone_index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"snapshot_schedule_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"recovery_point_objective_secs": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"local_snapshot_retention_policy": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"num_snapshots": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"rollup_retention_policy_multiple": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"rollup_retention_policy_snapshot_interval_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"auto_suspend_timeout_secs": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"snapshot_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"remote_snapshot_retention_policy": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"num_snapshots": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"rollup_retention_policy_multiple": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"rollup_retention_policy_snapshot_interval_type": {
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
				},
			},
			"ordered_availability_zone_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"category_filter": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"params": {
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
				},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixProtectionRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API
	protectionRuleID, iOk := d.GetOk("protection_rule_id")
	protectionRuleName, nOk := d.GetOk("protection_rule_name")

	if !iOk && !nOk {
		return diag.Errorf("please provide `protection_rule_id` or `role_name`")
	}

	var err error
	var resp *v3.ProtectionRuleResponse

	if iOk {
		resp, err = conn.V3.GetProtectionRule(protectionRuleID.(string))
	}
	if nOk {
		resp, err = findProtectionRuleByName(conn, protectionRuleName.(string))
	}

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("categories", c); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_reference", flattenReferenceValuesList(resp.Metadata.ProjectReference)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_reference", flattenReferenceValuesList(resp.Metadata.OwnerReference)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", resp.Spec.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("start_time", resp.Spec.Resources.StartTime); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("category_filter", flattenCategoriesFilter(resp.Spec.Resources.CategoryFilter)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("availability_zone_connectivity_list",
		flattenAvailabilityZoneConnectivityList(resp.Spec.Resources.AvailabilityZoneConnectivityList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ordered_availability_zone_list",
		flattenOrderAvailibilityList(resp.Spec.Resources.OrderedAvailabilityZoneList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", resp.Status.State); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}

func findProtectionRuleByName(conn *v3.Client, name string) (*v3.ProtectionRuleResponse, error) {
	filter := fmt.Sprintf("name==%s", name)
	resp, err := conn.V3.ListAllProtectionRules(filter)
	if err != nil {
		return nil, err
	}

	entities := resp.Entities

	found := make([]*v3.ProtectionRuleResponse, 0)
	for _, v := range entities {
		if v.Spec.Name == name {
			found = append(found, v)
		}
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("your query returned more than one result. Please use role_id argument instead")
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("role with the given name, not found")
	}

	return found[0], nil
}

package nutanix

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixProtectionRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixProtectionRulesRead,
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
							Type:     schema.TypeList,
							MaxItems: 1,
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
							MaxItems: 1,
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
													MaxItems: 1,
													MinItems: 1,
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
													MaxItems: 1,
													MinItems: 1,
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
							MaxItems: 1,
							MinItems: 1,
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
				},
			},
		},
	}
}

func dataSourceNutanixProtectionRulesRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API
	req := &v3.DSMetadata{}

	metadata, filtersOk := d.GetOk("metadata")
	if filtersOk {
		req = buildDataSourceListMetadata(metadata.(*schema.Set))
	}
	resp, err := conn.V3.ListAllProtectionRules(utils.StringValue(req.Filter))
	if err != nil {
		return err
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}
	if err := d.Set("entities", flattenProtectionRuleEntities(resp.Entities)); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenProtectionRuleEntities(protectionRules []*v3.ProtectionRuleResponse) []map[string]interface{} {
	entities := make([]map[string]interface{}, len(protectionRules))

	for i, protectionRule := range protectionRules {
		metadata, categories := setRSEntityMetadata(protectionRule.Metadata)

		entities[i] = map[string]interface{}{
			"name":                                protectionRule.Status.Name,
			"description":                         protectionRule.Spec.Description,
			"metadata":                            metadata,
			"categories":                          categories,
			"project_reference":                   flattenReferenceValuesList(protectionRule.Metadata.ProjectReference),
			"owner_reference":                     flattenReferenceValuesList(protectionRule.Metadata.OwnerReference),
			"start_time":                          protectionRule.Status.Resources.StartTime,
			"availability_zone_connectivity_list": flattenAvailabilityZoneConnectivityList(protectionRule.Spec.Resources.AvailabilityZoneConnectivityList),
			"ordered_availability_zone_list":      flattenOrderAvailibilityList(protectionRule.Spec.Resources.OrderedAvailabilityZoneList),
			"category_filter":                     flattenCategoriesFilter(protectionRule.Spec.Resources.CategoryFilter),
			"state":                               protectionRule.Status.State,
			"api_version":                         protectionRule.APIVersion,
		}
	}
	return entities
}

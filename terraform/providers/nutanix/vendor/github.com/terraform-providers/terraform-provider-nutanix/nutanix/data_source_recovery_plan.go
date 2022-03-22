package nutanix

import (
	"fmt"
	"strings"

	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNutanixRecoveryPlan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixRecoveryPlanRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"recovery_plan_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"recovery_plan_name"},
			},
			"recovery_plan_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"recovery_plan_id"},
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
			"stage_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"stage_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay_time_secs": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"stage_work": {
							Type:     schema.TypeList,
							Computed: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"recover_entities": {
										Type:     schema.TypeList,
										Computed: true,
										MinItems: 1,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"entity_info_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"any_entity_reference_kind": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"any_entity_reference_uuid": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"any_entity_reference_name": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"categories": categoriesSchema(),
															"script_list": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"enable_script_exec": {
																			Type:     schema.TypeBool,
																			Computed: true,
																		},
																		"timeout": {
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
								},
							},
						},
					},
				},
			},
			"parameters": {
				Type:     schema.TypeList,
				Computed: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"floating_ip_assignment_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vm_ip_assignment_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"test_floating_ip_config": {
													Type:     schema.TypeList,
													MinItems: 1,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"should_allocate_dynamically": {
																Type:     schema.TypeBool,
																Computed: true,
															},
														},
													},
												},
												"recovery_floating_ip_config": {
													Type:     schema.TypeList,
													MinItems: 1,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"should_allocate_dynamically": {
																Type:     schema.TypeBool,
																Computed: true,
															},
														},
													},
												},
												"vm_reference": {
													Type:     schema.TypeList,
													MinItems: 1,
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
												"vm_nic_information": {
													Type:     schema.TypeList,
													MinItems: 1,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"uuid": {
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
						"network_mapping_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone_network_mapping_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"availability_zone_url": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"recovery_network": {
													Type:     schema.TypeList,
													Computed: true,
													MinItems: 1,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"virtual_network_reference": {
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
															"vpc_reference": {
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
															"subnet_list": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"gateway_ip": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"external_connectivity_state": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																	},
																},
															},
															"use_vpc_reference": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"name": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"test_network": {
													Type:     schema.TypeList,
													Computed: true,
													MinItems: 1,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"virtual_network_reference": {
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
															"vpc_reference": {
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
															"subnet_list": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"gateway_ip": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"external_connectivity_state": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																	},
																},
															},
															"use_vpc_reference": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"name": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"recovery_ip_assignment_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"vm_reference": {
																Type:     schema.TypeList,
																MaxItems: 1,
																Required: true,
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
															"ip_config_list": {
																Type:     schema.TypeList,
																Required: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip_address": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																	},
																},
															},
														},
													},
												},
												"test_ip_assignment_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"vm_reference": {
																Type:     schema.TypeList,
																MaxItems: 1,
																Required: true,
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
															"ip_config_list": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip_address": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																	},
																},
															},
														},
													},
												},
												"cluster_reference_list": {
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
											},
										},
									},
									"are_networks_stretched": {
										Type:     schema.TypeBool,
										Computed: true,
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

func dataSourceNutanixRecoveryPlanRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	recoveryPlanID, iOk := d.GetOk("recovery_plan_id")
	recoveryPlanName, nOk := d.GetOk("recovery_plan_name")

	if !iOk && !nOk {
		return fmt.Errorf("please provide `recovery_plan_id` or `recovery_plan_name`")
	}

	var err error
	var resp *v3.RecoveryPlanResponse

	if iOk {
		resp, err = conn.V3.GetRecoveryPlan(recoveryPlanID.(string))
	}
	if nOk {
		resp, err = findRecoveryPlanByName(conn, recoveryPlanName.(string))
	}
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}
	if err := d.Set("project_reference", flattenReferenceValuesList(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", flattenReferenceValuesList(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	if err := d.Set("name", resp.Spec.Name); err != nil {
		return err
	}
	if err := d.Set("description", resp.Spec.Description); err != nil {
		return err
	}
	if err := d.Set("stage_list", flattenStageList(resp.Spec.Resources.StageList)); err != nil {
		return err
	}
	if err := d.Set("parameters", flattenParameters(resp.Spec.Resources.Parameters)); err != nil {
		return err
	}
	if err := d.Set("state", resp.Status.State); err != nil {
		return err
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}

func findRecoveryPlanByName(conn *v3.Client, name string) (*v3.RecoveryPlanResponse, error) {
	filter := fmt.Sprintf("name==%s", name)
	resp, err := conn.V3.ListAllRecoveryPlans(filter)
	if err != nil {
		return nil, err
	}

	entities := resp.Entities

	found := make([]*v3.RecoveryPlanResponse, 0)
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

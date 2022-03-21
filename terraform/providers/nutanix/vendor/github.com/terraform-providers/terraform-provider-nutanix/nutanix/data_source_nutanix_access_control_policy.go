package nutanix

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixAccessControlPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixAccessControlPolicyRead,
		Schema: map[string]*schema.Schema{
			"access_control_policy_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"access_control_policy_name"},
			},
			"access_control_policy_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"access_control_policy_id"},
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
			"user_reference_list": {
				Type:     schema.TypeSet,
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
			"user_group_reference_list": {
				Type:     schema.TypeSet,
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
			"role_reference": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
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
			"context_filter_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope_filter_expression_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"left_hand_side": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"right_hand_side": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"collection": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"categories": {
													Type:     schema.TypeList,
													MaxItems: 1,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"value": {
																Type:     schema.TypeList,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
												"uuid_list": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
						"entity_filter_expression_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"left_hand_side_entity_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"right_hand_side": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"collection": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"categories": {
													Type:     schema.TypeList,
													MaxItems: 1,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"value": {
																Type:     schema.TypeList,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
												"uuid_list": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
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
	}
}

func dataSourceNutanixAccessControlPolicyRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	id, iok := d.GetOk("access_control_policy_id")
	name, nOk := d.GetOk("access_control_policy_name")

	if !iok && !nOk {
		return fmt.Errorf("please provide `access_control_policy_id` or `access_control_policy_name`")
	}

	var reqErr error
	var resp *v3.AccessControlPolicy

	if iok {
		resp, reqErr = conn.V3.GetAccessControlPolicy(id.(string))
	}
	if nOk {
		resp, reqErr = findACPByName(conn, name.(string))
	}

	if reqErr != nil {
		return reqErr
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
	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}

	if status := resp.Status; status != nil {
		if err := d.Set("name", utils.StringValue(resp.Status.Name)); err != nil {
			return err
		}
		if err := d.Set("description", utils.StringValue(resp.Status.Description)); err != nil {
			return err
		}
		if err := d.Set("state", utils.StringValue(resp.Status.State)); err != nil {
			return err
		}

		if res := status.Resources; res != nil {
			if err := d.Set("user_reference_list", flattenArrayReferenceValues(status.Resources.UserReferenceList)); err != nil {
				return err
			}
			if err := d.Set("user_group_reference_list", flattenArrayReferenceValues(status.Resources.UserGroupReferenceList)); err != nil {
				return err
			}
			if err := d.Set("role_reference", flattenReferenceValuesList(status.Resources.RoleReference)); err != nil {
				return err
			}
			if status.Resources.FilterList.ContextList != nil {
				if err := d.Set("context_filter_list", flattenContextList(status.Resources.FilterList.ContextList)); err != nil {
					return err
				}
			}
		}
	}
	d.SetId(utils.StringValue(resp.Metadata.UUID))

	return nil
}

func findACPByName(conn *v3.Client, name string) (*v3.AccessControlPolicy, error) {
	filter := fmt.Sprintf("name==%s", name)
	resp, err := conn.V3.ListAllAccessControlPolicy(filter)
	if err != nil {
		return nil, err
	}

	entities := resp.Entities

	found := make([]*v3.AccessControlPolicy, 0)
	for _, v := range entities {
		if *v.Spec.Name == name {
			found = append(found, v)
		}
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("your query returned more than one result. Please use access_control_policy_id argument instead")
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("access control policy with the given name, not found")
	}

	return found[0], nil
}

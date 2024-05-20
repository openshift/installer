package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// dataSourceNutanixServiceGroups returns schema for datasource nutanix service groups
// https://www.nutanix.dev/api_references/prism-central-v3/#/b3A6MjU1ODc2OTg-list-the-service-groups
// v3/service_groups/list
func dataSourceNutanixServiceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixServiceGroupsRead,
		Schema: map[string]*schema.Schema{
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
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_group": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"protocol": {
													Type:     schema.TypeString,
													Computed: true,
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
											},
										},
									},
									"is_system_defined": {
										Type:        schema.TypeBool,
										Description: "specifying whether it is a system defined service group",
										Computed:    true,
									},
								},
							},
						},
						"associated_policies_list": {
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
		},
	}
}

func flattenServiceGroups(entries []*v3.ServiceGroupListEntry) interface{} {
	entities := make([]map[string]interface{}, len(entries))
	for i, entry := range entries {
		entities[i] = map[string]interface{}{
			"uuid": utils.StringValue(entry.UUID),
			"service_group": []map[string]interface{}{
				{
					"name":              utils.StringValue(entry.ServiceGroup.Name),
					"description":       utils.StringValue(entry.ServiceGroup.Description),
					"is_system_defined": utils.BoolValue(entry.ServiceGroup.SystemDefined),
					"service_list":      flattenServiceEntry(entry.ServiceGroup),
				},
			},
			"associated_policies_list": flattenReferenceList(entry.AssociatedPoliciesList),
		}
	}
	return entities
}

func dataSourceNutanixServiceGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

	req := &v3.DSMetadata{}

	metadata, filtersOk := d.GetOk("metadata")
	if filtersOk {
		req = buildDataSourceListMetadata(metadata.(*schema.Set))
	}

	resp, err := conn.V3.ListAllServiceGroups(utils.StringValue(req.Filter))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("entities", flattenServiceGroups(resp.Entities)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

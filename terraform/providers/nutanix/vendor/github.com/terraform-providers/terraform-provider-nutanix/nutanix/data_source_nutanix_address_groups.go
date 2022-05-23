package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixAddressGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixAddressGroupsRead,
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
						"address_group": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ip_address_block_list": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:     schema.TypeString,
													Required: true,
												},
												"prefix_length": {
													Type:     schema.TypeInt,
													Required: true,
												},
											},
										},
									},
									"address_group_string": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceNutanixAddressGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection
	conn := meta.(*Client).API
	req := &v3.DSMetadata{}

	metadata, filtersOk := d.GetOk("metadata")
	if filtersOk {
		req = buildDataSourceListMetadata(metadata.(*schema.Set))
	}

	resp, err := conn.V3.ListAllAddressGroups(utils.StringValue(req.Filter))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("entities", flattenAddressGroup(resp.Entities)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenAddressGroup(entries []*v3.AddressGroupListEntry) interface{} {
	entities := make([]map[string]interface{}, len(entries))

	for i, entry := range entries {
		entities[i] = map[string]interface{}{
			"address_group": []map[string]interface{}{
				{
					"name":                  entry.AddressGroup.Name,
					"description":           entry.AddressGroup.Description,
					"address_group_string":  entry.AddressGroup.AddressGroupString,
					"ip_address_block_list": flattenAddressEntry(entry.AddressGroup.BlockList),
				},
			},
			"associated_policies_list": flattenReferenceList(entry.AssociatedPoliciesList),
		}
	}
	return entities
}

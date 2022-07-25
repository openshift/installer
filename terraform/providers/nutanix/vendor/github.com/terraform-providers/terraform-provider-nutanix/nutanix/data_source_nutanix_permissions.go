package nutanix

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixPermissionsRead,
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
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
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
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
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"project_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
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
						"operation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fields": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"field_name_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
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

func dataSourceNutanixPermissionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading Permissions: %s", d.Id())

	// Get client connection
	conn := meta.(*Client).API

	resp, err := conn.V3.ListAllPermission("")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return diag.FromErr(err)
	}

	entities := make([]map[string]interface{}, len(resp.Entities))
	for k, v := range resp.Entities {
		entity := make(map[string]interface{})

		m, c := setRSEntityMetadata(v.Metadata)

		entity["metadata"] = m
		entity["categories"] = c
		entity["project_reference"] = flattenReferenceValues(v.Metadata.ProjectReference)
		entity["owner_reference"] = flattenReferenceValues(v.Metadata.OwnerReference)
		entity["api_version"] = utils.StringValue(v.APIVersion)
		entity["state"] = utils.StringValue(v.Status.State)
		entity["name"] = utils.StringValue(v.Status.Name)
		entity["description"] = utils.StringValue(v.Status.Description)
		entity["operation"] = utils.StringValue(v.Status.Resources.Operation)
		entity["kind"] = utils.StringValue(v.Status.Resources.Kind)
		entity["fields"] = flattenFieldsPermission(v.Status.Resources.Fields)

		entities[k] = entity
	}

	if err := d.Set("entities", entities); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

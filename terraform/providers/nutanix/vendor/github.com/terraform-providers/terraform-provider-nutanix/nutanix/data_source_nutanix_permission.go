package nutanix

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixPermission() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixPermissionRead,
		Schema: map[string]*schema.Schema{
			"permission_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"permission_name"},
			},
			"permission_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"permission_id"},
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
				MaxItems: 1,
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
	}
}

func dataSourceNutanixPermissionRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	permissionID, iok := d.GetOk("permission_id")
	permissionName, rnOk := d.GetOk("permission_name")

	if !iok && !rnOk {
		return fmt.Errorf("please provide `permission_id` or `permission_name`")
	}

	var err error
	var resp *v3.PermissionIntentResponse

	if iok {
		resp, err = conn.V3.GetPermission(permissionID.(string))
	}
	if rnOk {
		resp, err = findPermissionByName(conn, permissionName.(string))
	}

	if err != nil {
		return err
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
			if err := d.Set("operation", utils.StringValue(res.Operation)); err != nil {
				return err
			}
			if err := d.Set("kind", utils.StringValue(res.Kind)); err != nil {
				return err
			}
			if err := d.Set("fields", flattenFieldsPermission(res.Fields)); err != nil {
				return err
			}
		}
	}
	d.SetId(utils.StringValue(resp.Metadata.UUID))

	return nil
}

func flattenFieldsPermission(fieldPermissions *v3.FieldsPermission) []map[string]interface{} {
	flatFieldsPermissions := make([]map[string]interface{}, 0)
	n := map[string]interface{}{
		"field_mode":      fieldPermissions.FieldMode,
		"field_name_list": fieldPermissions.FieldNameList,
	}
	flatFieldsPermissions = append(flatFieldsPermissions, n)
	return flatFieldsPermissions
}

func findPermissionByName(conn *v3.Client, name string) (*v3.PermissionIntentResponse, error) {
	filter := fmt.Sprintf("name==%s", name)
	resp, err := conn.V3.ListAllPermission(filter)
	if err != nil {
		return nil, err
	}

	entities := resp.Entities

	found := make([]*v3.PermissionIntentResponse, 0)
	for _, v := range entities {
		if *v.Spec.Name == name {
			found = append(found, v)
		}
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("your query returned more than one result. Please use permission_id argument instead")
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("permission with the given name, not found")
	}
	found[0].APIVersion = resp.APIVersion
	return found[0], nil
}

package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/customattribute"
	"github.com/vmware/govmomi/object"
)

func dataSourceVSphereCustomAttribute() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereCustomAttributeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the custom attribute.",
				Required:    true,
			},
			"managed_object_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Object type for which the custom attribute is valid. If not specified, the attribute is valid for all managed object types.",
			},
		},
	}
}

func dataSourceVSphereCustomAttributeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	err := customattribute.VerifySupport(client)
	if err != nil {
		return err
	}

	fm, err := object.GetCustomFieldsManager(client.Client)
	if err != nil {
		return err
	}

	field, err := customattribute.ByName(fm, d.Get("name").(string))
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprint(field.Key))
	d.Set("managed_object_type", field.ManagedObjectType)
	return nil
}

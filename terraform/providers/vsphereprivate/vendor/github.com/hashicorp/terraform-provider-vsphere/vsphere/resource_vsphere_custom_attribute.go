package vsphere

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/customattribute"
	"github.com/vmware/govmomi/object"
)

func resourceVSphereCustomAttribute() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereCustomAttributeCreate,
		Read:   resourceVSphereCustomAttributeRead,
		Update: resourceVSphereCustomAttributeUpdate,
		Delete: resourceVSphereCustomAttributeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereCustomAttributeImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the custom attribute.",
				Required:    true,
			},
			"managed_object_type": {
				Type:        schema.TypeString,
				Description: "Object type for which the custom attribute is valid. If not specified, the attribute is valid for all managed object types.",
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceVSphereCustomAttributeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	if err := customattribute.VerifySupport(client); err != nil {
		return err
	}

	fm, err := object.GetCustomFieldsManager(client.Client)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	field, err := fm.Add(ctx, d.Get("name").(string), d.Get("managed_object_type").(string), nil, nil)
	if err != nil {
		return fmt.Errorf("could not create custom attribute: %s", err)
	}

	d.SetId(fmt.Sprint(field.Key))
	d.Set("name", field.Name)
	d.Set("managed_object_type", field.ManagedObjectType)
	return nil
}

func resourceVSphereCustomAttributeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	if err := customattribute.VerifySupport(client); err != nil {
		return err
	}

	fm, err := object.GetCustomFieldsManager(client.Client)
	if err != nil {
		return err
	}
	key, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	fields, err := fm.Field(ctx)
	if err != nil {
		return err
	}
	field := fields.ByKey(int32(key))
	if field == nil {
		return fmt.Errorf("could not locate category with id '%d'", key)
	}
	d.Set("name", field.Name)
	d.Set("managed_object_type", field.ManagedObjectType)
	return nil
}

func resourceVSphereCustomAttributeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	if err := customattribute.VerifySupport(client); err != nil {
		return err
	}

	fm, err := object.GetCustomFieldsManager(client.Client)
	if err != nil {
		return err
	}
	key, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	return fm.Rename(ctx, int32(key), d.Get("name").(string))
}

func resourceVSphereCustomAttributeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	if err := customattribute.VerifySupport(client); err != nil {
		return err
	}

	fm, err := object.GetCustomFieldsManager(client.Client)
	if err != nil {
		return err
	}
	key, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	return fm.Remove(ctx, int32(key))
}

func resourceVSphereCustomAttributeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*VSphereClient).vimClient
	if err := customattribute.VerifySupport(client); err != nil {
		return nil, err
	}

	fm, err := object.GetCustomFieldsManager(client.Client)
	if err != nil {
		return nil, err
	}

	field, err := customattribute.ByName(fm, d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(fmt.Sprint(field.Key))
	return []*schema.ResourceData{d}, nil
}

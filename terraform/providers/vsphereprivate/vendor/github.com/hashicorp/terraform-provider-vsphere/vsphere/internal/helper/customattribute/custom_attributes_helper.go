package customattribute

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// ConfigKey is the string that should always be used as the
// key for specifying custom attribute values for a resource. Various resource
// custom attribute helpers will depend on this value being consistent across
// resources.
//
// When adding custom attributes to a resource schema, the easiest way to do
// that (for now) will be to use the following line:
//
//   customattribute.ConfigKey: customattribute.ConfigSchema(),
//
// This will ensure that the correct key and schema is used across all
// resources.
const ConfigKey = "custom_attributes"

// ConfigSchema returns the schema for custom attribute configuration
// for each resource that needs it.
//
// The key should be set to the ConfigKey constant and should be a
// map of custom attribute ids to values.
func ConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Description: "A list of custom attributes to set on this resource.",
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	}
}

func VerifySupport(client *govmomi.Client) error {
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return errors.New("Custom attributes are only supported on vCenter")
	}
	return nil
}

func IsSupported(client *govmomi.Client) bool {
	return VerifySupport(client) == nil
}

// ReadFromResource reads the custom attributes from an object and saves the
// data into the supplied ResourceData.
//
// TODO: Add error handling and reporting to this method.
func ReadFromResource(entity *mo.ManagedEntity, d *schema.ResourceData) {
	customAttrs := make(map[string]interface{})
	if len(entity.CustomValue) > 0 {
		for _, fv := range entity.CustomValue {
			value := fv.(*types.CustomFieldStringValue).Value
			if value != "" {
				customAttrs[fmt.Sprint(fv.GetCustomFieldValue().Key)] = value
			}
		}
	}
	_ = d.Set(ConfigKey, customAttrs)
}

type DiffProcessor struct {
	// The field manager
	fm *object.CustomFieldsManager

	// Old map of custom attribute key to values
	oldAttributes map[string]interface{}

	// New map of custom attribute key to values
	newAttributes map[string]interface{}
}

func (p *DiffProcessor) clearRemovedAttributes(subject object.Reference) error {
	for k := range p.oldAttributes {
		_, ok := p.newAttributes[k]
		if !ok {
			key, err := strconv.ParseInt(k, 10, 32)
			if err != nil {
				return err
			}
			err = p.fm.Set(context.TODO(), subject.Reference(), int32(key), "")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *DiffProcessor) setNewAttributes(subject object.Reference) error {
	for k, v := range p.newAttributes {
		key, err := strconv.ParseInt(k, 10, 32)
		if err != nil {
			return err
		}
		err = p.fm.Set(context.TODO(), subject.Reference(), int32(key), v.(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *DiffProcessor) ProcessDiff(subject object.Reference) error {
	if err := p.clearRemovedAttributes(subject); err != nil {
		return fmt.Errorf("error clearing removed attributes for object ID %q: %s", subject.Reference().Value, err)
	}
	if err := p.setNewAttributes(subject); err != nil {
		return fmt.Errorf("error setting attributes for object ID %q: %s", subject.Reference().Value, err)
	}
	return nil
}

func GetDiffProcessorIfAttributesDefined(client *govmomi.Client, d *schema.ResourceData) (*DiffProcessor, error) {
	old, newValue := d.GetChange(ConfigKey)
	if len(old.(map[string]interface{})) > 0 || len(newValue.(map[string]interface{})) > 0 {
		if err := VerifySupport(client); err != nil {
			return nil, err
		}
	} else {
		// No custom attributes defined
		return nil, nil
	}
	fm, err := object.GetCustomFieldsManager(client.Client)
	if err != nil {
		return nil, err
	}
	return &DiffProcessor{
		fm:            fm,
		oldAttributes: old.(map[string]interface{}),
		newAttributes: newValue.(map[string]interface{}),
	}, nil
}

func ByName(fm *object.CustomFieldsManager, name string) (*types.CustomFieldDef, error) {
	fields, err := fm.Field(context.TODO())
	if err != nil {
		return nil, err
	}

	for _, def := range fields {
		if def.Name == name {
			return &def, nil
		}
	}

	return nil, object.ErrKeyNameNotFound
}

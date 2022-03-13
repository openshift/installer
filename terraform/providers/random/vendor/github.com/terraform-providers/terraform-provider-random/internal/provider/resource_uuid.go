package provider

import (
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUuid() *schema.Resource {
	return &schema.Resource{
		Description: "The resource `random_uuid` generates random uuid string that is intended to be " +
			"used as unique identifiers for other resources.\n" +
			"\n" +
			"This resource uses [hashicorp/go-uuid](https://github.com/hashicorp/go-uuid) to generate a " +
			"UUID-formatted string for use with services needed a unique string identifier.",
		Create: CreateUuid,
		Read:   schema.Noop,
		Delete: schema.RemoveFromState,
		Importer: &schema.ResourceImporter{
			State: ImportUuid,
		},

		Schema: map[string]*schema.Schema{
			"keepers": {
				Description: "Arbitrary map of values that, when changed, will trigger recreation of " +
					"resource. See [the main provider documentation](../index.html) for more information.",
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"result": {
				Description: "The generated uuid presented in string format.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"id": {
				Description: "The generated uuid presented in string format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func CreateUuid(d *schema.ResourceData, meta interface{}) error {
	result, err := uuid.GenerateUUID()
	if err != nil {
		return errwrap.Wrapf("error generating uuid: {{err}}", err)
	}
	d.Set("result", result)
	d.SetId(result)
	return nil
}

func ImportUuid(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	bytes, err := uuid.ParseUUID(id)
	if err != nil {
		return nil, errwrap.Wrapf("error parsing uuid bytes: {{err}}", err)
	}
	result, err2 := uuid.FormatUUID(bytes)
	if err2 != nil {
		return nil, errwrap.Wrapf("error formatting uuid bytes: {{err2}}", err2)
	}

	d.Set("result", result)
	d.SetId(result)

	return []*schema.ResourceData{d}, nil
}

package ignition

import (
	"encoding/json"

	"github.com/coreos/ignition/v2/config/v3_1/types"
	"github.com/coreos/vcontext/path"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSystemdUnit() *schema.Resource {
	return &schema.Resource{
		Exists: resourceSystemdUnitExists,
		Read:   resourceSystemdUnitRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"mask": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dropin": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"rendered": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSystemdUnitRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildSystemdUnit(d)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceSystemdUnitExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildSystemdUnit(d)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildSystemdUnit(d *schema.ResourceData) (string, error) {
	enabled := d.Get("enabled").(bool)
	unit := &types.Unit{
		Name:    d.Get("name").(string),
		Enabled: &enabled,
	}

	content, hasContent := d.GetOk("content")
	if hasContent {
		str := content.(string)
		unit.Contents = &str
	}

	mask, hasMask := d.GetOk("mask")
	if hasMask {
		bmask := mask.(bool)
		unit.Mask = &bmask
	}

	for _, raw := range d.Get("dropin").([]interface{}) {
		value := raw.(map[string]interface{})

		d := types.Dropin{
			Name: value["name"].(string),
		}

		contents := value["content"].(string)
		if contents != "" {
			d.Contents = &contents
		}

		if err := handleReport(d.Validate(path.ContextPath{})); err != nil {
			return "", err
		}

		unit.Dropins = append(unit.Dropins, d)
	}

	if err := handleReport(unit.Validate(path.ContextPath{})); err != nil {
		return "", err
	}

	b, err := json.Marshal(unit)
	if err != nil {
		return "", err
	}
	err = d.Set("rendered", string(b))
	if err != nil {
		return "", err
	}

	return hash(string(b)), nil
}

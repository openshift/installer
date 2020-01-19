package ignition

import (
	"github.com/coreos/ignition/v2/config/v3_0/types"
	"github.com/coreos/vcontext/path"
	"github.com/hashicorp/terraform/helper/schema"
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
		},
	}
}

func resourceSystemdUnitRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildSystemdUnit(d, globalCache)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceSystemdUnitExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildSystemdUnit(d, globalCache)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildSystemdUnit(d *schema.ResourceData, c *cache) (string, error) {
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

	return c.addSystemdUnit(unit), handleReport(unit.Validate(path.ContextPath{}))
}

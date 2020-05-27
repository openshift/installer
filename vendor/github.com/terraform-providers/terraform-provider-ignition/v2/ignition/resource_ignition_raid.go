package ignition

import (
	"encoding/json"

	"github.com/coreos/ignition/v2/config/v3_0/types"
	"github.com/coreos/vcontext/path"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRaid() *schema.Resource {
	return &schema.Resource{
		Exists: resourceRaidExists,
		Read:   resourceRaidRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"devices": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"spares": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"rendered": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRaidRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildRaid(d)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceRaidExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildRaid(d)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildRaid(d *schema.ResourceData) (string, error) {
	raid := &types.Raid{
		Name:  d.Get("name").(string),
		Level: d.Get("level").(string),
	}
	spares, hasSpares := d.GetOk("spares")
	if hasSpares {
		ispares := spares.(int)
		raid.Spares = &ispares
	}

	for _, value := range d.Get("devices").([]interface{}) {
		raid.Devices = append(raid.Devices, types.Device(value.(string)))
	}

	if err := handleReport(raid.Validate(path.ContextPath{})); err != nil {
		return "", err
	}

	b, err := json.Marshal(raid)
	if err != nil {
		return "", err
	}
	d.Set("rendered", string(b))

	return hash(string(b)), nil
}

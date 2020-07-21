package ignition

import (
	"encoding/json"

	"github.com/coreos/ignition/v2/config/v3_1/types"
	"github.com/coreos/vcontext/path"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDisk() *schema.Resource {
	return &schema.Resource{
		Exists: resourceDiskExists,
		Read:   resourceDiskRead,
		Schema: map[string]*schema.Schema{
			"device": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"wipe_table": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"partition": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"sizemib": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"startmib": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"type_guid": {
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

func resourceDiskRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildDisk(d)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceDiskExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildDisk(d)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildDisk(d *schema.ResourceData) (string, error) {
	disk := &types.Disk{
		Device: d.Get("device").(string),
	}
	wipe, hasWipeTable := d.GetOk("wipe_table")
	if hasWipeTable {
		bwipe := wipe.(bool)
		disk.WipeTable = &bwipe
	}

	for _, raw := range d.Get("partition").([]interface{}) {
		v := raw.(map[string]interface{})
		p := types.Partition{
			Number: v["number"].(int),
		}
		tlabel := v["label"].(string)
		if tlabel != "" {
			p.Label = &tlabel
		}
		tsize := v["sizemib"].(int)
		if tsize != 0 {
			p.SizeMiB = &tsize
		}
		tstart := v["startmib"].(int)
		if tstart != 0 {
			p.StartMiB = &tstart
		}
		tguid := v["type_guid"].(string)
		p.TypeGUID = &tguid

		disk.Partitions = append(disk.Partitions, p)
	}

	if err := handleReport(disk.Validate(path.ContextPath{})); err != nil {
		return "", err
	}

	b, err := json.Marshal(disk)
	if err != nil {
		return "", err
	}
	err = d.Set("rendered", string(b))
	if err != nil {
		return "", err
	}

	return hash(string(b)), nil
}

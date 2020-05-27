package ignition

import (
	"encoding/json"

	"github.com/coreos/ignition/v2/config/v3_0/types"
	"github.com/coreos/ignition/v2/config/validate"
	"github.com/coreos/vcontext/path"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceLink() *schema.Resource {
	return &schema.Resource{
		Exists: resourceLinkExists,
		Read:   resourceLinkRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hard": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"uid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"gid": {
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

func resourceLinkRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildLink(d)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceLinkExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildLink(d)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildLink(d *schema.ResourceData) (string, error) {
	link := &types.Link{}
	link.Path = d.Get("path").(string)
	if err := handleReport(link.Node.Validate(path.ContextPath{})); err != nil {
		return "", err
	}

	link.Target = d.Get("target").(string)

	hard, hasHard := d.GetOk("hard")
	if hasHard {
		bhard := hard.(bool)
		link.Hard = &bhard
	}

	uid := d.Get("uid").(int)
	if uid != 0 {
		link.User = types.NodeUser{ID: &uid}
	}

	gid := d.Get("gid").(int)
	if gid != 0 {
		link.Group = types.NodeGroup{ID: &gid}
	}

	b, err := json.Marshal(link)
	if err != nil {
		return "", err
	}
	d.Set("rendered", string(b))

	return hash(string(b)), handleReport(validate.ValidateWithContext(new(*types.Link), b))
}

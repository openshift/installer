package ignition

import (
	"github.com/coreos/ignition/v2/config/v3_0/types"
	"github.com/coreos/vcontext/path"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDirectory() *schema.Resource {
	return &schema.Resource{
		Exists: resourceDirectoryExists,
		Read:   resourceDirectoryRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mode": {
				Type:     schema.TypeInt,
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
		},
	}
}

func resourceDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildDirectory(d, globalCache)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceDirectoryExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildDirectory(d, globalCache)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildDirectory(d *schema.ResourceData, c *cache) (string, error) {
	dir := &types.Directory{}
	dir.Path = d.Get("path").(string)
	mode := d.Get("mode").(int)
	if mode != 0 {
		dir.Mode = &mode
	}
	uid := d.Get("uid").(int)
	if uid != 0 {
		dir.User = types.NodeUser{ID: &uid}
	}

	gid := d.Get("gid").(int)
	if gid != 0 {
		dir.Group = types.NodeGroup{ID: &gid}
	}

	if err := handleReport(dir.Validate(path.ContextPath{})); err != nil {
		return "", err
	}

	return c.addDirectory(dir), nil
}

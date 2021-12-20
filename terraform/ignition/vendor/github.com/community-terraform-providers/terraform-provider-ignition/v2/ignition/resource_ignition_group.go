package ignition

import (
	"encoding/json"

	"github.com/coreos/ignition/v2/config/v3_1/types"
	"github.com/coreos/ignition/v2/config/validate"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Exists: resourceGroupExists,
		Read:   resourceGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"password_hash": {
				Type:     schema.TypeString,
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

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildGroup(d)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildGroup(d)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildGroup(d *schema.ResourceData) (string, error) {
	group := &types.PasswdGroup{
		Name: d.Get("name").(string),
		Gid:  getInt(d, "gid"),
	}

	passhash, hasPasshash := d.GetOk("password_hash")
	if hasPasshash {
		str := passhash.(string)
		group.PasswordHash = &str
	}

	b, err := json.Marshal(group)
	if err != nil {
		return "", err
	}
	err = d.Set("rendered", string(b))
	if err != nil {
		return "", err
	}

	return hash(string(b)), handleReport(validate.ValidateWithContext(new(*types.PasswdGroup), b))
}

package ignition

import (
	"encoding/json"

	"github.com/coreos/ignition/v2/config/v3_3/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKargs() *schema.Resource {
	return &schema.Resource{
		Exists: resourceKargsExists,
		Read:   resourceKargsRead,
		Schema: map[string]*schema.Schema{
			"shouldexist": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"shouldnotexist": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rendered": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKargsRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildKargs(d)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceKargsExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildKargs(d)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildKargs(d *schema.ResourceData) (string, error) {
	kargs := &types.KernelArguments{}

	for _, value := range d.Get("shouldexist").([]interface{}) {
		kargs.ShouldExist = append(kargs.ShouldExist, types.KernelArgument(value.(string)))
	}

	for _, value := range d.Get("shouldnotexist").([]interface{}) {
		kargs.ShouldNotExist = append(kargs.ShouldNotExist, types.KernelArgument(value.(string)))
	}

	b, err := json.Marshal(kargs)
	if err != nil {
		return "", err
	}
	err = d.Set("rendered", string(b))
	if err != nil {
		return "", err
	}

	return hash(string(b)), nil
}

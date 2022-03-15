package vsphere

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/spbm"
)

func dataSourceVSphereStoragePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereStoragePolicyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the storage policy.",
				Required:    true,
			},
		},
	}
}

func dataSourceVSphereStoragePolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient

	id, err := spbm.PolicyIDByName(client, d.Get("name").(string))
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

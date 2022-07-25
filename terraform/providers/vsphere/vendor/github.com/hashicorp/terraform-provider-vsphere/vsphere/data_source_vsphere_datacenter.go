package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVSphereDatacenter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereDatacenterRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Description: "The name of the datacenter. This can be a name or path.	Can be omitted if there is only one datacenter in your inventory.",
				Optional: true,
			},
		},
	}
}

func dataSourceVSphereDatacenterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient
	datacenter := d.Get("name").(string)
	dc, err := getDatacenter(client, datacenter)
	if err != nil {
		return fmt.Errorf("error fetching datacenter: %s", err)
	}
	id := dc.Reference().Value
	d.SetId(id)

	return nil
}

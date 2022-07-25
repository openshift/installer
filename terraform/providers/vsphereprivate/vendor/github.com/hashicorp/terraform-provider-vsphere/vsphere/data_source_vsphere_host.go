package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
)

func dataSourceVSphereHost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereHostRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Description: "The name of the host. This can be a name or path.	If not provided, the default host is used.",
				Optional: true,
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The managed object ID of the datacenter to look for the host in.",
				Required:    true,
			},
			"resource_pool_id": {
				Type:        schema.TypeString,
				Description: "The managed object ID of the host's root resource pool.",
				Computed:    true,
			},
		},
	}
}

func dataSourceVSphereHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient
	name := d.Get("name").(string)
	dcID := d.Get("datacenter_id").(string)
	dc, err := datacenterFromID(client, dcID)
	if err != nil {
		return fmt.Errorf("error fetching datacenter: %s", err)
	}
	hs, err := hostsystem.SystemOrDefault(client, name, dc)
	if err != nil {
		return fmt.Errorf("error fetching host: %s", err)
	}
	rp, err := hostsystem.ResourcePool(hs)
	if err != nil {
		return err
	}
	err = d.Set("resource_pool_id", rp.Reference().Value)
	if err != nil {
		return err
	}
	id := hs.Reference().Value
	d.SetId(id)
	return nil
}

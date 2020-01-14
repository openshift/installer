package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/datastore"
	"github.com/vmware/govmomi/object"
)

func dataSourceVSphereDatastore() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereDatastoreRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name or path of the datastore.",
				Required:    true,
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The managed object ID of the datacenter the datastore is in. This is not required when using ESXi directly, or if there is only one datacenter in your infrastructure.",
				Optional:    true,
			},
		},
	}
}

func dataSourceVSphereDatastoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient

	name := d.Get("name").(string)
	var dc *object.Datacenter
	if dcID, ok := d.GetOk("datacenter_id"); ok {
		var err error
		dc, err = datacenterFromID(client, dcID.(string))
		if err != nil {
			return fmt.Errorf("cannot locate datacenter: %s", err)
		}
	}
	ds, err := datastore.FromPath(client, name, dc)
	if err != nil {
		return fmt.Errorf("error fetching datastore: %s", err)
	}

	d.SetId(ds.Reference().Value)
	return nil
}

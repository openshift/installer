package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/vappcontainer"
)

func dataSourceVSphereVAppContainer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereVAppContainerRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the vApp container.",
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Managed Object ID of the datacenter.",
			},
		},
	}
}

func dataSourceVSphereVAppContainerRead(d *schema.ResourceData, meta interface{}) error {
	client, err := resourceVSphereVAppContainerClient(meta)
	if err != nil {
		return err
	}
	dc, err := datacenterFromID(client, d.Get("datacenter_id").(string))
	if err != nil {
		return fmt.Errorf("cannot locate datacenter: %s", err)
	}
	vc, err := vappcontainer.FromPath(client, d.Get("name").(string), dc)
	if err != nil {
		return fmt.Errorf("cannot locate vApp Container: %s", err)
	}
	d.SetId(vc.Reference().Value)
	return nil
}

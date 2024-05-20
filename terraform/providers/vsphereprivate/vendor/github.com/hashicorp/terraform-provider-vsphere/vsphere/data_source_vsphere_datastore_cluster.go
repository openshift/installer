package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVSphereDatastoreCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereDatastoreClusterRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name or absolute path to the datastore cluster.",
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The managed object ID of the datacenter the cluster is located in. Not required if using an absolute path.",
			},
		},
	}
}

func dataSourceVSphereDatastoreClusterRead(d *schema.ResourceData, meta interface{}) error {
	pod, err := resourceVSphereDatastoreClusterGetPodFromPath(meta, d.Get("name").(string), d.Get("datacenter_id").(string))
	if err != nil {
		return fmt.Errorf("error loading datastore cluster: %s", err)
	}
	d.SetId(pod.Reference().Value)
	return nil
}

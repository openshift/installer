package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/clustercomputeresource"
)

func dataSourceVSphereComputeCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereComputeClusterRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name or absolute path to the cluster.",
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The managed object ID of the datacenter the cluster is located in. Not required if using an absolute path.",
			},
			"resource_pool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The managed object ID of the cluster's root resource pool.",
			},
		},
	}
}

func dataSourceVSphereComputeClusterRead(d *schema.ResourceData, meta interface{}) error {
	cluster, err := resourceVSphereComputeClusterGetClusterFromPath(meta, d.Get("name").(string), d.Get("datacenter_id").(string))
	if err != nil {
		return fmt.Errorf("error loading cluster: %s", err)
	}
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return fmt.Errorf("error loading cluster properties: %s", err)
	}

	d.SetId(cluster.Reference().Value)
	_ = d.Set("resource_pool_id", props.ResourcePool.Value)
	return nil
}

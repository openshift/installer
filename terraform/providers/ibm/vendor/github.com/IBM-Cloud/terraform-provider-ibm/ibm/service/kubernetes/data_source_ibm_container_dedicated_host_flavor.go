// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMContainerDedicatedHostFlavor() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMContainerDedicatedHostFlavorRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The zone of the dedicated host flavor",
			},
			"host_flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the dedicated host flavor",
			},
			"flavor_class": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The class of the dedicated host flavor",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region of the dedicated host flavor",
			},
			"deprecated": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Describes if the dedicated host flavor is deprecated",
			},
			"max_vcpus": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum available vcpus in the dedicated host flavor",
			},
			"max_memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum available memory in the dedicated host flavor",
			},
			"instance_storage": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The instance storage of the dedicated host flavor",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceIBMContainerDedicatedHostFlavorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	flavorID := d.Get("host_flavor_id").(string)
	client, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	dedicatedHostFlavorAPI := client.DedicatedHostFlavor()
	targetEnv := v2.ClusterTargetHeader{}

	dedicatedHostFlavors, err := dedicatedHostFlavorAPI.ListDedicatedHostFlavors(d.Get("zone").(string), targetEnv)
	if err != nil {
		return diag.Errorf("[ERROR] Listing dedicated host flavors in zone %s failed: %v", d.Get("zone").(string), err)
	}

	var dedicatedHostFlavor *v2.GetDedicatedHostFlavor

	for _, dh := range dedicatedHostFlavors {
		if dh.ID == flavorID {
			dedicatedHostFlavor = &dh
			break
		}
	}

	if dedicatedHostFlavor == nil {
		return diag.Errorf("[ERROR] Dedicated host flavor is not found, id : %s", d.Get("id").(string))
	}

	d.SetId(dedicatedHostFlavor.ID)
	d.Set("flavor_class", dedicatedHostFlavor.FlavorClass)
	d.Set("region", dedicatedHostFlavor.Region)
	d.Set("deprecated", dedicatedHostFlavor.Deprecated)
	d.Set("max_vcpus", dedicatedHostFlavor.MaxVCPUs)
	d.Set("max_memory", dedicatedHostFlavor.MaxMemory)

	instanceStorage := make([]map[string]interface{}, len(dedicatedHostFlavor.InstanceStorage))
	for i, is := range dedicatedHostFlavor.InstanceStorage {
		instanceStorage[i] = map[string]interface{}{
			"count": is.Count,
			"size":  is.Size,
		}
	}
	d.Set("instance_storage", instanceStorage)

	return nil
}

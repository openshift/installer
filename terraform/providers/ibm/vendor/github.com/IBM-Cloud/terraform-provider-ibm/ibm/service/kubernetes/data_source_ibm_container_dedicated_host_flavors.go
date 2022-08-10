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

func DataSourceIBMContainerDedicatedHostFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMContainerDedicatedHostFlavorsRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The zone of the dedicated host flavors",
			},
			"host_flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_flavor_id": {
							Type:        schema.TypeString,
							Computed:    true,
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
				},
			},
		},
	}
}
func dataSourceIBMContainerDedicatedHostFlavorsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	dedicatedHostFlavorsAPI := client.DedicatedHostFlavor()
	targetEnv := v2.ClusterTargetHeader{}

	dedicatedHostFlavors, err := dedicatedHostFlavorsAPI.ListDedicatedHostFlavors(d.Get("zone").(string), targetEnv)
	if err != nil {
		return diag.Errorf("[ERROR] Listing dedicated host flavors in zone %s failed: %v", d.Get("zone").(string), err)
	}

	flavors := make([]interface{}, len(dedicatedHostFlavors))

	for j, dh := range dedicatedHostFlavors {
		instanceStorage := make([]map[string]interface{}, len(dh.InstanceStorage))
		for i, is := range dh.InstanceStorage {
			instanceStorage[i] = map[string]interface{}{
				"count": is.Count,
				"size":  is.Size,
			}
		}
		flavors[j] = map[string]interface{}{
			"host_flavor_id":   dh.ID,
			"flavor_class":     dh.FlavorClass,
			"region":           dh.Region,
			"deprecated":       dh.Deprecated,
			"max_vcpus":        dh.MaxVCPUs,
			"max_memory":       dh.MaxMemory,
			"instance_storage": instanceStorage,
		}
	}

	d.SetId(d.Get("zone").(string))
	d.Set("host_flavors", flavors)
	return nil
}

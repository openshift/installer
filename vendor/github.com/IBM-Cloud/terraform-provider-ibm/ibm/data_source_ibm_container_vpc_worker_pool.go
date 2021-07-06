// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMContainerVpcClusterWorkerPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerVpcClusterWorkerPoolRead,
		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster name",
			},
			"worker_pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "worker pool name",
			},
			"flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zones": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"worker_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"isolation": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceIBMContainerVpcClusterWorkerPoolRead(d *schema.ResourceData, meta interface{}) error {
	wpClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	clusterName := d.Get("cluster").(string)
	workerPoolName := d.Get("worker_pool_name").(string)
	workerPoolsAPI := wpClient.WorkerPools()
	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	workerPool, err := workerPoolsAPI.GetWorkerPool(clusterName, workerPoolName, targetEnv)
	if err != nil {
		return err
	}

	var zones = make([]map[string]interface{}, 0)
	for _, zone := range workerPool.Zones {
		for _, subnet := range zone.Subnets {
			zoneInfo := map[string]interface{}{
				"name":      zone.ID,
				"subnet_id": subnet.ID,
			}
			zones = append(zones, zoneInfo)
		}
	}
	d.Set("worker_pool_name", workerPool.PoolName)
	d.Set("flavor", workerPool.Flavor)
	d.Set("worker_count", workerPool.WorkerCount)
	d.Set("provider", workerPool.Provider)
	d.Set("labels", workerPool.Labels)
	d.Set("zones", zones)
	d.Set("cluster", clusterName)
	d.Set("vpc_id", workerPool.VpcID)
	d.Set("isolation", workerPool.Isolation)
	d.Set("resource_group_id", targetEnv.ResourceGroup)
	d.SetId(workerPool.ID)
	return nil
}

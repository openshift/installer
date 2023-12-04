// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"fmt"

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMSatelliteClusterWorkerPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMSatelliteClusterWorkerPoolRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "worker pool name",
			},
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster name",
			},
			"flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The flavor of the satellite worker node",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the worker pool",
			},
			"zones": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"worker_pool_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Labels on all the workers in the worker pool",
			},
			"host_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Host labels on the workers",
			},
			"operating_system": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The operating system of the hosts in the worker pool",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the resource group",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the region",
			},
			"worker_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of workers that are attached",
			},
			"isolation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Isolation of the worker node",
			},
			"auto_scale_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable auto scalling for worker pool",
			},
			"openshift_license_source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "License source for Openshift",
			},
		},
	}
}
func dataSourceIBMSatelliteClusterWorkerPoolRead(d *schema.ResourceData, meta interface{}) error {
	var name, cluster, resourceGrp string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	if v, ok := d.GetOk("cluster"); ok {
		cluster = v.(string)
	}

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	getSatWorkerPoolOptions := &kubernetesserviceapiv1.GetWorkerPoolOptions{}
	getSatWorkerPoolOptions.Workerpool = &name
	getSatWorkerPoolOptions.Cluster = &cluster

	if v, ok := d.GetOk("resource_group_id"); ok && v != "" {
		resourceGrp = v.(string)
		getSatWorkerPoolOptions.XAuthResourceGroup = &resourceGrp
	}

	if v, ok := d.GetOk("region"); ok && v != "" {
		wpRegion := v.(string)
		getSatWorkerPoolOptions.XRegion = &wpRegion
	}

	workerPool, response, err := satClient.GetWorkerPool(getSatWorkerPoolOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving worker pool  %s: %s\n%s", name, err, response)
	}

	var zones = make([]map[string]interface{}, 0)
	for _, zone := range workerPool.Zones {
		zoneInfo := map[string]interface{}{
			"zone":         *zone.ID,
			"worker_count": *zone.WorkerCount,
		}
		zones = append(zones, zoneInfo)
	}

	d.Set("name", *workerPool.PoolName)
	d.Set("flavor", *workerPool.Flavor)
	d.Set("worker_count", *workerPool.WorkerCount)
	d.Set("worker_pool_labels", workerPool.Labels)
	d.Set("host_labels", workerPool.HostLabels)
	d.Set("operating_system", *workerPool.OperatingSystem)
	d.Set("zones", zones)
	d.Set("cluster", cluster)
	d.Set("auto_scale_enabled", *workerPool.AutoscaleEnabled)
	d.Set("state", *workerPool.Lifecycle.ActualState)
	d.Set("isolation", *workerPool.Isolation)
	d.Set("openshift_license_source", *workerPool.OpenshiftLicense)
	d.SetId(*workerPool.ID)

	return nil
}

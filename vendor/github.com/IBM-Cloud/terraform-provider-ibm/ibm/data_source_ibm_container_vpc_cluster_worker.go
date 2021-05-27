// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMContainerVPCClusterWorker() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerVPCClusterWorkerRead,

		Schema: map[string]*schema.Schema{
			"worker_id": {
				Description: "ID of the worker",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cluster_name_id": {
				Description: "Name or ID of the cluster",
				Type:        schema.TypeString,
				Required:    true,
			},
			"flavor": {
				Description: "flavor of the worker",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"kube_version": {
				Description: "kube version of the worker",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"state": {
				Description: "State of the worker",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pool_id": {
				Description: "worker pool id",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pool_name": {
				Description: "worker pool name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"network_interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
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
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the resource group.",
			},
			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this cluster",
			},
		},
	}
}

func dataSourceIBMContainerVPCClusterWorkerRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	wrkAPI := csClient.Workers()
	workerID := d.Get("worker_id").(string)
	clusterID := d.Get("cluster_name_id").(string)

	workerFields, err := wrkAPI.Get(clusterID, workerID, targetEnv)
	if err != nil {
		return fmt.Errorf("Error retrieving worker: %s", err)
	}

	d.SetId(workerFields.ID)
	d.Set("flavor", workerFields.Flavor)
	d.Set("kube_version", workerFields.KubeVersion.Actual)
	d.Set("state", workerFields.Health.State)
	d.Set("pool_id", workerFields.PoolID)
	d.Set("pool_name", workerFields.PoolName)
	d.Set("network_interfaces", flattenNetworkInterfaces(workerFields.NetworkInterfaces))
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/kubernetes/clusters")

	return nil
}

// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMContainerClusterWorker() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerClusterWorkerRead,

		Schema: map[string]*schema.Schema{
			"worker_id": {
				Description: "ID of the worker",
				Type:        schema.TypeString,
				Required:    true,
			},
			"state": {
				Description: "State of the worker",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: "Status of the worker",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"private_vlan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_vlan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_guid": {
				Description: "The bluemix organization guid this cluster belongs to",
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated",
			},
			"space_guid": {
				Description: "The bluemix space guid this cluster belongs to",
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated",
			},
			"account_guid": {
				Description: "The bluemix account guid this cluster belongs to",
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cluster region",
				Deprecated:  "This field is deprecated",
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

func dataSourceIBMContainerClusterWorkerRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	wrkAPI := csClient.Workers()
	workerID := d.Get("worker_id").(string)
	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	workerFields, err := wrkAPI.Get(workerID, targetEnv)
	if err != nil {
		return fmt.Errorf("Error retrieving worker: %s", err)
	}

	d.SetId(workerFields.ID)
	d.Set("state", workerFields.State)
	d.Set("status", workerFields.Status)
	d.Set("private_vlan", workerFields.PrivateVlan)
	d.Set("public_vlan", workerFields.PublicVlan)
	d.Set("private_ip", workerFields.PrivateIP)
	d.Set("public_ip", workerFields.PublicIP)
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/kubernetes/clusters")

	return nil
}

// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMContainerWorkerPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerWorkerPoolRead,

		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name or ID of the cluster",
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_container_worker_pool",
					"cluster"),
			},

			"worker_pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "worker pool name",
			},

			"machine_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "worker nodes machine type",
			},

			"size_per_zone": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of nodes per zone",
			},

			"hardware": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hardware type",
			},

			"disk_encryption": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "worker node disk encrypted if set to true",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "worker pool state",
			},

			"zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "worker pool zones",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "worker pool zone name",
						},

						"private_vlan": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "worker pool zone private vlan",
						},

						"public_vlan": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "worker pool zone public vlan",
						},

						"worker_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "worker pool zone worker count",
						},
					},
				},
			},

			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "list of labels to worker pool",
			},

			"operating_system": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The operating system of the workers in the worker pool",
			},

			"resource_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the resource group.",
			},

			"autoscale_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Autoscaling is enabled on the workerpool",
			},
		},
	}
}
func DataSourceIBMContainerWorkerPoolValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cluster",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "cluster",
			CloudDataRange:             []string{"resolved_to:id"}})

	iBMContainerWorkerPoolValidator := validate.ResourceValidator{ResourceName: "ibm_container_worker_pool", Schema: validateSchema}
	return &iBMContainerWorkerPoolValidator
}
func dataSourceIBMContainerWorkerPoolRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	workerPoolName := d.Get("worker_pool_name").(string)
	cluster := d.Get("cluster").(string)

	workerPoolsAPI := csClient.WorkerPools()
	targetEnv, err := getWorkerPoolTargetHeader(d, meta)
	if err != nil {
		return err
	}

	workerPool, err := workerPoolsAPI.GetWorkerPool(cluster, workerPoolName, targetEnv)
	if err != nil {
		return err
	}

	machineType := workerPool.MachineType
	d.SetId(workerPool.ID)
	d.Set("machine_type", strings.Split(machineType, ".encrypted")[0])
	d.Set("size_per_zone", workerPool.Size)
	hardware := workerPool.Isolation
	switch strings.ToLower(hardware) {
	case "":
		hardware = hardwareShared
	case isolationPrivate:
		hardware = hardwareDedicated
	case isolationPublic:
		hardware = hardwareShared
	}
	d.Set("hardware", hardware)
	d.Set("state", workerPool.State)
	if workerPool.Labels != nil {
		d.Set("labels", workerPool.Labels)
	}
	d.Set("operating_system", workerPool.OperatingSystem)
	d.Set("zones", flex.FlattenZones(workerPool.Zones))
	if strings.Contains(machineType, "encrypted") {
		d.Set("disk_encryption", true)
	} else {
		d.Set("disk_encryption", false)
	}
	d.Set("resource_group_id", targetEnv.ResourceGroup)
	d.Set("autoscale_enabled", workerPool.AutoscaleEnabled)
	return nil
}

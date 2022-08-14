// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPIInstances() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPIInstancesAllRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"pvm_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pvm_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"processors": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"health_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"networks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"macaddress": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"external_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"proctype": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"minproc": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"minmem": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"maxproc": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"maxmem": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"pin_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_cores_assigned": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_virtual_cores": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_virtual_cores": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_pool": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_pool_affinity": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"license_repository_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						PIPlacementGroupID: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIInstancesAllRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).IBMPISession()

	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	powerC := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
	powervmdata, err := powerC.GetAll()

	if err != nil {
		return diag.FromErr(err)
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set("pvm_instances", flattenPvmInstances(powervmdata.PvmInstances))

	return nil
}

func flattenPvmInstances(list []*models.PVMInstanceReference) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {

		l := map[string]interface{}{
			"pvm_instance_id":             *i.PvmInstanceID,
			"memory":                      *i.Memory,
			"processors":                  *i.Processors,
			"proctype":                    *i.ProcType,
			"status":                      *i.Status,
			"minproc":                     i.Minproc,
			"minmem":                      i.Minmem,
			"maxproc":                     i.Maxproc,
			"maxmem":                      i.Maxmem,
			"pin_policy":                  i.PinPolicy,
			"virtual_cores_assigned":      i.VirtualCores.Assigned,
			"max_virtual_cores":           i.VirtualCores.Max,
			"min_virtual_cores":           i.VirtualCores.Min,
			"storage_type":                i.StorageType,
			"storage_pool":                i.StoragePool,
			"storage_pool_affinity":       i.StoragePoolAffinity,
			"license_repository_capacity": i.LicenseRepositoryCapacity,
			PIPlacementGroupID:            i.PlacementGroup,
			"networks":                    flattenPvmInstanceNetworks(i.Networks),
		}

		if i.Health != nil {
			l["health_status"] = i.Health.Status
		}

		result = append(result, l)

	}
	return result
}

func flattenPvmInstanceNetworks(list []*models.PVMInstanceNetwork) (networks []map[string]interface{}) {
	if list != nil {
		networks = make([]map[string]interface{}, len(list))
		for i, pvmip := range list {

			p := make(map[string]interface{})
			p["ip"] = pvmip.IPAddress
			p["network_name"] = pvmip.NetworkName
			p["network_id"] = pvmip.NetworkID
			p["macaddress"] = pvmip.MacAddress
			p["type"] = pvmip.Type
			p["external_ip"] = pvmip.ExternalIP
			networks[i] = p
		}
		return networks
	}
	return
}

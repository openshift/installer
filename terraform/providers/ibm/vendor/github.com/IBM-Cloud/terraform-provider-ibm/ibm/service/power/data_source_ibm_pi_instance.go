// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPIInstance() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPIInstancesRead,
		Schema: map[string]*schema.Schema{

			helpers.PIInstanceName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Server Name to be used for pvminstances",
				ValidateFunc: validation.NoZeroValues,
			},

			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"addresses": {
				Type:       schema.TypeList,
				Computed:   true,
				Deprecated: "This field is deprecated, use networks instead",
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
						/*"version": {
							Type:     schema.TypeFloat,
							Computed: true,
						},*/
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
	}
}

func dataSourceIBMPIInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).IBMPISession()

	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	powerC := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
	powervmdata, err := powerC.Get(d.Get(helpers.PIInstanceName).(string))

	if err != nil {
		return diag.FromErr(err)
	}

	pvminstanceid := *powervmdata.PvmInstanceID
	d.SetId(pvminstanceid)
	d.Set("memory", powervmdata.Memory)
	d.Set("processors", powervmdata.Processors)
	d.Set("status", powervmdata.Status)
	d.Set("proctype", powervmdata.ProcType)
	d.Set("volumes", powervmdata.VolumeIDs)
	d.Set("minproc", powervmdata.Minproc)
	d.Set("minmem", powervmdata.Minmem)
	d.Set("maxproc", powervmdata.Maxproc)
	d.Set("maxmem", powervmdata.Maxmem)
	d.Set("pin_policy", powervmdata.PinPolicy)
	d.Set("virtual_cores_assigned", powervmdata.VirtualCores.Assigned)
	d.Set("max_virtual_cores", powervmdata.VirtualCores.Max)
	d.Set("min_virtual_cores", powervmdata.VirtualCores.Min)
	d.Set("storage_type", powervmdata.StorageType)
	d.Set("storage_pool", powervmdata.StoragePool)
	d.Set("storage_pool_affinity", powervmdata.StoragePoolAffinity)
	d.Set("license_repository_capacity", powervmdata.LicenseRepositoryCapacity)
	d.Set("networks", flattenPvmInstanceNetworks(powervmdata.Networks))
	if *powervmdata.PlacementGroup != "none" {
		d.Set(PIPlacementGroupID, powervmdata.PlacementGroup)
	}

	if powervmdata.Addresses != nil {
		pvmaddress := make([]map[string]interface{}, len(powervmdata.Addresses))
		for i, pvmip := range powervmdata.Addresses {

			p := make(map[string]interface{})
			p["ip"] = pvmip.IPAddress
			p["network_name"] = pvmip.NetworkName
			p["network_id"] = pvmip.NetworkID
			p["macaddress"] = pvmip.MacAddress
			p["type"] = pvmip.Type
			p["external_ip"] = pvmip.ExternalIP
			pvmaddress[i] = p
		}
		d.Set("addresses", pvmaddress)

	}

	if powervmdata.Health != nil {

		d.Set("health_status", powervmdata.Health.Status)

	}

	return nil
}

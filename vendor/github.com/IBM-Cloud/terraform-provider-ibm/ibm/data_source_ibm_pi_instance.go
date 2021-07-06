// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
)

func dataSourceIBMPIInstance() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMPIInstancesRead,
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"processors": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"health_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"addresses": {
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
				Type:     schema.TypeInt,
				Computed: true,
			},
			"minmem": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"maxproc": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"maxmem": {
				Type:     schema.TypeInt,
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
		},
	}
}

func dataSourceIBMPIInstancesRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()

	if err != nil {
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)

	powerC := instance.NewIBMPIInstanceClient(sess, powerinstanceid)
	powervmdata, err := powerC.Get(d.Get(helpers.PIInstanceName).(string), powerinstanceid, getTimeOut)

	if err != nil {
		return err
	}

	pvminstanceid := *powervmdata.PvmInstanceID
	d.SetId(pvminstanceid)
	d.Set("memory", powervmdata.Memory)
	d.Set("processors", powervmdata.Processors)
	d.Set("status", powervmdata.Status)
	d.Set("proctype", powervmdata.ProcType)
	d.Set("volumes", powervmdata.VolumeIds)
	d.Set("minproc", powervmdata.Minproc)
	d.Set("minmem", powervmdata.Minmem)
	d.Set("maxproc", powervmdata.Maxproc)
	d.Set("maxmem", powervmdata.Maxmem)
	d.Set("pin_policy", powervmdata.PinPolicy)
	d.Set("virtual_cores_assigned", powervmdata.VirtualCores.Assigned)
	d.Set("max_virtual_cores", powervmdata.VirtualCores.Max)
	d.Set("min_virtual_cores", powervmdata.VirtualCores.Min)

	if powervmdata.Addresses != nil {
		pvmaddress := make([]map[string]interface{}, len(powervmdata.Addresses))
		for i, pvmip := range powervmdata.Addresses {

			p := make(map[string]interface{})
			p["ip"] = pvmip.IP
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

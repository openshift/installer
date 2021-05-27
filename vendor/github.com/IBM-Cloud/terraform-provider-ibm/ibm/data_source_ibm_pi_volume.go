// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	//"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
)

func dataSourceIBMPIVolume() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMPIVolumeRead,
		Schema: map[string]*schema.Schema{

			helpers.PIVolumeName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Volume Name to be used for pvminstances",
				ValidateFunc: validation.NoZeroValues,
			},

			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"shareable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"bootable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"disk_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wwn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMPIVolumeRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	volumeC := instance.NewIBMPIVolumeClient(sess, powerinstanceid)
	volumedata, err := volumeC.Get(d.Get(helpers.PIVolumeName).(string), powerinstanceid, getTimeOut)
	if err != nil {
		return err
	}

	d.SetId(*volumedata.VolumeID)
	if volumedata.Size != nil {
		d.Set("size", volumedata.Size)
	}
	if &volumedata.DiskType != nil {
		d.Set("disk_type", volumedata.DiskType)
	}
	if &volumedata.Bootable != nil {
		d.Set("bootable", volumedata.Bootable)
	}
	if &volumedata.State != nil {
		d.Set("state", volumedata.State)
	}
	if &volumedata.Wwn != nil {
		d.Set("wwn", volumedata.Wwn)
	}
	return nil

}

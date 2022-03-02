// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	//"fmt"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
)

func dataSourceIBMPIVolume() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeRead,
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
			"bootable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_pool": {
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

func dataSourceIBMPIVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	volumeC := instance.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
	volumedata, err := volumeC.Get(d.Get(helpers.PIVolumeName).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*volumedata.VolumeID)
	d.Set("size", volumedata.Size)
	d.Set("state", volumedata.State)
	d.Set("shareable", volumedata.Shareable)
	d.Set("bootable", volumedata.Bootable)
	d.Set("disk_type", volumedata.DiskType)
	d.Set("volume_pool", volumedata.VolumePool)
	d.Set("wwn", volumedata.Wwn)

	return nil
}

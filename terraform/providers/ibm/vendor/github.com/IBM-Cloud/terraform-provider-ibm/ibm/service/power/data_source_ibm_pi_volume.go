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

func DataSourceIBMPIVolume() *schema.Resource {
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
			"replication_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the volume should be replication enabled or not",
			},
			"replication_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Replication type(metro,global)",
			},
			"replication_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Replication status of a volume",
			},
			"auxiliary": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "true if volume is auxiliary otherwise false",
			},
			"consistency_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Consistency Group Name if volume is a part of volume group",
			},
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume Group ID",
			},
			"mirroring_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mirroring state for replication enabled volume",
			},
			"primary_role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates whether master/aux volume is playing the primary role",
			},
			"auxiliary_volume_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates auxiliary volume name",
			},
			"master_volume_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates master volume name",
			},
		},
	}
}

func dataSourceIBMPIVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
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
	d.Set("replication_enabled", volumedata.ReplicationEnabled)
	d.Set("replication_type", volumedata.ReplicationType)
	d.Set("replication_status", volumedata.ReplicationStatus)
	d.Set("auxiliary", volumedata.Auxiliary)
	d.Set("consistency_group_name", volumedata.ConsistencyGroupName)
	d.Set("group_id", volumedata.GroupID)
	d.Set("mirroring_state", volumedata.MirroringState)
	d.Set("primary_role", volumedata.PrimaryRole)
	d.Set("auxiliary_volume_name", volumedata.AuxVolumeName)
	d.Set("master_volume_name", volumedata.MasterVolumeName)

	return nil
}

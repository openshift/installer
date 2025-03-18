// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIVolume() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_VolumeName: {
				Description:  "Volume Name to be used for pvminstances",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_Auxiliary: {
				Computed:    true,
				Description: "Indicates if the volume is auxiliary or not.",
				Type:        schema.TypeBool,
			},
			Attr_AuxiliaryVolumeName: {
				Computed:    true,
				Description: "The auxiliary volume name.",
				Type:        schema.TypeString,
			},
			Attr_Bootable: {
				Computed:    true,
				Description: "Indicates if the volume is boot capable.",
				Type:        schema.TypeBool,
			},
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			Attr_ConsistencyGroupName: {
				Computed:    true,
				Description: "Consistency group name if volume is a part of volume group.",
				Type:        schema.TypeString,
			},
			Attr_CreationDate: {
				Computed:    true,
				Description: "Date volume was created.",
				Type:        schema.TypeString,
			},
			Attr_DiskType: {
				Computed:    true,
				Description: "The disk type that is used for the volume.",
				Type:        schema.TypeString,
			},
			Attr_FreezeTime: {
				Computed:    true,
				Description: "The freeze time of remote copy.",
				Type:        schema.TypeString,
			},
			Attr_GroupID: {
				Computed:    true,
				Description: "The volume group id in which the volume belongs.",
				Type:        schema.TypeString,
			},
			Attr_IOThrottleRate: {
				Computed:    true,
				Description: "Amount of iops assigned to the volume",
				Type:        schema.TypeString,
			},
			Attr_LastUpdateDate: {
				Computed:    true,
				Description: "The last updated date of the volume.",
				Type:        schema.TypeString,
			},
			Attr_MasterVolumeName: {
				Computed:    true,
				Description: "The master volume name.",
				Type:        schema.TypeString,
			},
			Attr_MirroringState: {
				Computed:    true,
				Description: "Mirroring state for replication enabled volume.",
				Type:        schema.TypeString,
			},
			Attr_PrimaryRole: {
				Computed:    true,
				Description: "Indicates whether master/auxiliary volume is playing the primary role.",
				Type:        schema.TypeString,
			},
			Attr_ReplicationEnabled: {
				Computed:    true,
				Description: "Indicates if the volume should be replication enabled or not.",
				Type:        schema.TypeBool,
			},
			Attr_ReplicationSites: {
				Computed:    true,
				Description: "List of replication sites for volume replication.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
			Attr_ReplicationStatus: {
				Computed:    true,
				Description: "The replication status of the volume.",
				Type:        schema.TypeString,
			},
			Attr_ReplicationType: {
				Computed:    true,
				Description: "The replication type of the volume, metro or global.",
				Type:        schema.TypeString,
			},
			Attr_Shareable: {
				Computed:    true,
				Description: "Indicates if the volume is shareable between VMs.",
				Type:        schema.TypeBool,
			},
			Attr_Size: {
				Computed:    true,
				Description: "The size of the volume in GB.",
				Type:        schema.TypeInt,
			},
			Attr_State: {
				Computed:    true,
				Description: "The state of the volume.",
				Type:        schema.TypeString,
			},
			Attr_UserTags: {
				Computed:    true,
				Description: "List of user tags attached to the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
			Attr_VolumePool: {
				Computed:    true,
				Description: "Volume pool, name of storage pool where the volume is located.",
				Type:        schema.TypeString,
			},
			Attr_WWN: {
				Computed:    true,
				Description: "The world wide name of the volume.",
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceIBMPIVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	volumeC := instance.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
	volumedata, err := volumeC.Get(d.Get(Arg_VolumeName).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*volumedata.VolumeID)
	d.Set(Attr_Auxiliary, volumedata.Auxiliary)
	d.Set(Attr_AuxiliaryVolumeName, volumedata.AuxVolumeName)
	d.Set(Attr_Bootable, volumedata.Bootable)
	d.Set(Attr_ConsistencyGroupName, volumedata.ConsistencyGroupName)
	d.Set(Attr_CreationDate, volumedata.CreationDate.String())
	if volumedata.Crn != "" {
		d.Set(Attr_CRN, volumedata.Crn)
		tags, err := flex.GetTagsUsingCRN(meta, string(volumedata.Crn))
		if err != nil {
			log.Printf("Error on get of pi volume (%s) user_tags: %s", *volumedata.VolumeID, err)
		}
		d.Set(Attr_UserTags, tags)
	}
	d.Set(Attr_DiskType, volumedata.DiskType)
	if volumedata.FreezeTime != nil {
		d.Set(Attr_FreezeTime, volumedata.FreezeTime.String())
	}
	d.Set(Attr_GroupID, volumedata.GroupID)
	d.Set(Attr_IOThrottleRate, volumedata.IoThrottleRate)
	d.Set(Attr_LastUpdateDate, volumedata.LastUpdateDate.String())
	d.Set(Attr_MasterVolumeName, volumedata.MasterVolumeName)
	d.Set(Attr_MirroringState, volumedata.MirroringState)
	d.Set(Attr_PrimaryRole, volumedata.PrimaryRole)
	d.Set(Attr_ReplicationEnabled, volumedata.ReplicationEnabled)
	d.Set(Attr_ReplicationType, volumedata.ReplicationType)
	if len(volumedata.ReplicationSites) > 0 {
		d.Set(Attr_ReplicationSites, volumedata.ReplicationSites)
	}
	d.Set(Attr_ReplicationStatus, volumedata.ReplicationStatus)
	d.Set(Attr_State, volumedata.State)
	d.Set(Attr_Shareable, volumedata.Shareable)
	d.Set(Attr_Size, volumedata.Size)
	d.Set(Attr_VolumePool, volumedata.VolumePool)
	d.Set(Attr_WWN, volumedata.Wwn)

	return nil
}

// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIVolumeRemoteCopyRelationship() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeRemoteCopyRelationshipsReads,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_VolumeID: {
				Description:  "The ID of the volume for which you want to retrieve detailed information.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_AuxiliaryChangedVolumeName: {
				Computed:    true,
				Description: "The name of the volume that is acting as the auxiliary change volume for the relationship.",
				Type:        schema.TypeString,
			},
			Attr_AuxiliaryVolumeName: {
				Computed:    true,
				Description: "The auxiliary volume name at storage host level.",
				Type:        schema.TypeString,
			},
			Attr_ConsistencyGroupName: {
				Computed:    true,
				Description: "The consistency group name if volume is a part of volume group.",
				Type:        schema.TypeString,
			},
			Attr_CopyType: {
				Computed:    true,
				Description: "The copy type.",
				Type:        schema.TypeString,
			},
			Attr_CyclingMode: {
				Computed:    true,
				Description: "The type of cycling mode used.",
				Type:        schema.TypeString,
			},
			Attr_CyclePeriodSeconds: {
				Computed:    true,
				Description: "The minimum period in seconds between multiple cycles.",
				Type:        schema.TypeInt,
			},
			Attr_FreezeTime: {
				Computed:    true,
				Description: "The freeze time of remote copy relationship.",
				Type:        schema.TypeString,
			},
			Attr_MasterChangedVolumeName: {
				Computed:    true,
				Description: "The name of the volume that is acting as the master change volume for the relationship.",
				Type:        schema.TypeString,
			},
			Attr_MasterVolumeName: {
				Computed:    true,
				Description: "The master volume name at storage host level.",
				Type:        schema.TypeString,
			},
			Attr_Name: {
				Computed:    true,
				Description: "The remote copy relationship name.",
				Type:        schema.TypeString,
			},
			Attr_PrimaryRole: {
				Computed:    true,
				Description: "Indicates whether master/aux volume is playing the primary role.",
				Type:        schema.TypeString,
			},
			Attr_Progress: {
				Computed:    true,
				Description: "The relationship progress.",
				Type:        schema.TypeInt,
			},
			Attr_RemoteCopyID: {
				Computed:    true,
				Description: "The remote copy relationship ID.",
				Type:        schema.TypeString,
			},
			Attr_State: {
				Computed:    true,
				Description: "The relationship state.",
				Type:        schema.TypeString,
			},
			Attr_Synchronized: {
				Computed:    true,
				Description: "Indicates whether the relationship is synchronized.",
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceIBMPIVolumeRemoteCopyRelationshipsReads(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	volClient := instance.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
	volData, err := volClient.GetVolumeRemoteCopyRelationships(d.Get(Arg_VolumeID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(volData.ID)
	d.Set(Attr_AuxiliaryChangedVolumeName, volData.AuxChangedVolumeName)
	d.Set(Attr_AuxiliaryVolumeName, volData.AuxVolumeName)
	d.Set(Attr_ConsistencyGroupName, volData.ConsistencyGroupName)
	d.Set(Attr_CopyType, volData.CopyType)
	d.Set(Attr_CyclingMode, volData.CyclingMode)
	d.Set(Attr_CyclePeriodSeconds, volData.CyclePeriodSeconds)
	d.Set(Attr_FreezeTime, volData.FreezeTime.String())
	d.Set(Attr_MasterChangedVolumeName, volData.MasterChangedVolumeName)
	d.Set(Attr_MasterVolumeName, volData.MasterVolumeName)
	if volData.Name != nil {
		d.Set(Attr_Name, volData.Name)
	}
	d.Set(Attr_PrimaryRole, volData.PrimaryRole)
	d.Set(Attr_Progress, volData.Progress)
	if volData.RemoteCopyID != nil {
		d.Set(Attr_RemoteCopyID, volData.RemoteCopyID)
	}
	d.Set(Attr_State, volData.State)
	d.Set(Attr_Synchronized, volData.Sync)

	return nil
}

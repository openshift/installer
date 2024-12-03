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

func DataSourceIBMPIVolumeGroupStorageDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeGroupStorageDetailsReads,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_VolumeGroupID: {
				Description:  "The ID of the volume group.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_ConsistencyGroupName: {
				Computed:    true,
				Description: "The name of consistency group at storage controller level.",
				Type:        schema.TypeString,
			},
			Attr_CyclePeriodSeconds: {
				Computed:    true,
				Description: "The minimum period in seconds between multiple cycles.",
				Type:        schema.TypeInt,
			},
			Attr_CyclingMode: {
				Computed:    true,
				Description: "The type of cycling mode used.",
				Type:        schema.TypeString,
			},
			Attr_NumberOfVolumes: {
				Computed:    true,
				Description: "The number of volumes in volume group.",
				Type:        schema.TypeInt,
			},
			Attr_PrimaryRole: {
				Computed:    true,
				Description: "Indicates whether master/aux volume is playing the primary role.",
				Type:        schema.TypeString,
			},
			Attr_RemoteCopyRelationshipNames: {
				Computed:    true,
				Description: "List of remote-copy relationship names in a volume group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
			Attr_ReplicationType: {
				Computed:    true,
				Description: "The type of replication (metro, global).",
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

func dataSourceIBMPIVolumeGroupStorageDetailsReads(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	vgClient := instance.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)
	vgID := d.Get(Arg_VolumeGroupID).(string)
	vgData, err := vgClient.GetVolumeGroupLiveDetails(vgID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vgID)
	d.Set(Attr_ConsistencyGroupName, vgData.ConsistencyGroupName)
	d.Set(Attr_CyclePeriodSeconds, vgData.CyclePeriodSeconds)
	d.Set(Attr_CyclingMode, vgData.CyclingMode)
	d.Set(Attr_NumberOfVolumes, vgData.NumOfvols)
	d.Set(Attr_PrimaryRole, vgData.PrimaryRole)
	d.Set(Attr_RemoteCopyRelationshipNames, vgData.RemoteCopyRelationshipNames)
	d.Set(Attr_ReplicationType, vgData.ReplicationType)
	d.Set(Attr_State, vgData.State)
	d.Set(Attr_Synchronized, vgData.Sync)

	return nil
}

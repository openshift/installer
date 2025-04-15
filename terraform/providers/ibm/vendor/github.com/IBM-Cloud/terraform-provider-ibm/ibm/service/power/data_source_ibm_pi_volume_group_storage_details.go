// Copyright IBM Corp. 2022 All Rights Reserved.
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

func DataSourceIBMPIVolumeGroupStorageDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeGroupStorageDetailsReads,
		Schema: map[string]*schema.Schema{
			PIVolumeGroupID: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Volume group ID",
				ValidateFunc: validation.NoZeroValues,
			},
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"consistency_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of consistency group at storage controller level",
			},
			"cycle_period_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates the minimum period in seconds between multiple cycles",
			},
			"cycling_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the type of cycling mode used",
			},
			"number_of_volumes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of volumes in volume group",
			},
			"primary_role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates whether master/aux volume is playing the primary role",
			},
			"remote_copy_relationship_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of remote-copy relationship names in a volume group",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"replication_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of replication(metro,global)",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the relationship state",
			},
			"synchronized": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates whether the relationship is synchronized",
			},
		},
	}
}

func dataSourceIBMPIVolumeGroupStorageDetailsReads(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	vgClient := instance.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)
	vgID := d.Get(PIVolumeGroupID).(string)
	vgData, err := vgClient.GetVolumeGroupLiveDetails(vgID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vgID)
	d.Set("consistency_group_name", vgData.ConsistencyGroupName)
	d.Set("cycle_period_seconds", vgData.CyclePeriodSeconds)
	d.Set("cycling_mode", vgData.CyclingMode)
	d.Set("number_of_volumes", vgData.NumOfvols)
	d.Set("primary_role", vgData.PrimaryRole)
	d.Set("remote_copy_relationship_names", vgData.RemoteCopyRelationshipNames)
	d.Set("replication_type", vgData.ReplicationType)
	d.Set("state", vgData.State)
	d.Set("synchronized", vgData.Sync)

	return nil
}

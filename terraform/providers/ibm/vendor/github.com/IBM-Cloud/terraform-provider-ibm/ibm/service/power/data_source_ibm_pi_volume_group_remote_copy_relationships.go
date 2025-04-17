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

func DataSourceIBMPIVolumeGroupRemoteCopyRelationships() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeGroupRemoteCopyRelationshipsReads,
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
			"remote_copy_relationships": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of remote copy relationships",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auxiliary_changed_volume_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the volume that is acting as the auxiliary change volume for the relationship",
						},
						"auxiliary_volume_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Auxiliary volume name at storage host level",
						},
						"consistency_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Consistency Group Name if volume is a part of volume group",
						},
						"copy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the copy type.",
						},
						"cycling_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the type of cycling mode used.",
						},
						"freeze_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Freeze time of remote copy relationship",
						},
						"master_changed_volume_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the volume that is acting as the master change volume for the relationship",
						},
						"master_volume_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Master volume name at storage host level",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remote copy relationship name",
						},
						"primary_role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates whether master/aux volume is playing the primary role",
						},
						"progress": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the relationship progress",
						},
						"remote_copy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remote copy relationship ID",
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
				},
			},
		},
	}
}

func dataSourceIBMPIVolumeGroupRemoteCopyRelationshipsReads(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	vgClient := instance.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)
	vgData, err := vgClient.GetVolumeGroupRemoteCopyRelationships(d.Get(PIVolumeGroupID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	results := make([]map[string]interface{}, 0, len(vgData.RemoteCopyRelationships))
	for _, i := range vgData.RemoteCopyRelationships {
		if i != nil {
			l := map[string]interface{}{
				"auxiliary_changed_volume_name": i.AuxChangedVolumeName,
				"auxiliary_volume_name":         i.AuxVolumeName,
				"consistency_group_name":        i.ConsistencyGroupName,
				"copy_type":                     i.CopyType,
				"cycling_mode":                  i.CyclingMode,
				"freeze_time":                   i.FreezeTime.String(),
				"master_changed_volume_name":    i.MasterChangedVolumeName,
				"master_volume_name":            i.MasterVolumeName,
				"primary_role":                  i.PrimaryRole,
				"progress":                      i.Progress,
				"state":                         i.State,
				"synchronized":                  i.Sync,
			}
			if i.Name != nil {
				l["name"] = i.Name
			}
			if i.RemoteCopyID != nil {
				l["remote_copy_id"] = i.RemoteCopyID
			}

			results = append(results, l)
		}
	}

	d.SetId(vgData.ID)
	d.Set("remote_copy_relationships", results)

	return nil
}

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

func DataSourceIBMPIVolumeGroupRemoteCopyRelationships() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeGroupRemoteCopyRelationshipsReads,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_VolumeGroupID: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The ID of the volume group for which you want to retrieve detailed information.",
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_RemoteCopyRelationships: {
				Computed:    true,
				Description: "List of remote copy relationships",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPIVolumeGroupRemoteCopyRelationshipsReads(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	vgClient := instance.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)
	vgData, err := vgClient.GetVolumeGroupRemoteCopyRelationships(d.Get(Arg_VolumeGroupID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	results := make([]map[string]interface{}, 0, len(vgData.RemoteCopyRelationships))
	for _, i := range vgData.RemoteCopyRelationships {
		if i != nil {
			l := map[string]interface{}{
				Attr_AuxiliaryChangedVolumeName: i.AuxChangedVolumeName,
				Attr_AuxiliaryVolumeName:        i.AuxVolumeName,
				Attr_ConsistencyGroupName:       i.ConsistencyGroupName,
				Attr_CopyType:                   i.CopyType,
				Attr_CyclingMode:                i.CyclingMode,
				Attr_FreezeTime:                 i.FreezeTime.String(),
				Attr_MasterChangedVolumeName:    i.MasterChangedVolumeName,
				Attr_MasterVolumeName:           i.MasterVolumeName,
				Attr_PrimaryRole:                i.PrimaryRole,
				Attr_Progress:                   i.Progress,
				Attr_State:                      i.State,
				Attr_Synchronized:               i.Sync,
			}
			if i.Name != nil {
				l[Attr_Name] = i.Name
			}
			if i.RemoteCopyID != nil {
				l[Attr_RemoteCopyID] = i.RemoteCopyID
			}

			results = append(results, l)
		}
	}

	d.SetId(vgData.ID)
	d.Set(Attr_RemoteCopyRelationships, results)

	return nil
}

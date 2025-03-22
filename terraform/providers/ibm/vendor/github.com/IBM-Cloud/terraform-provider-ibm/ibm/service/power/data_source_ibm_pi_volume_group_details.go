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

func DataSourceIBMPIVolumeGroupDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeGroupDetailsRead,
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
				Description:  "Name of the volume group",
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_Auxiliary: {
				Computed:    true,
				Description: "Indicates if the volume is auxiliary or not.",
				Type:        schema.TypeBool,
			},
			Attr_ConsistencyGroupName: {
				Computed:    true,
				Description: "The name of consistency group at storage controller level.",
				Type:        schema.TypeString,
			},
			Attr_ReplicationStatus: {
				Computed:    true,
				Description: "The replication status of volume group.",
				Type:        schema.TypeString,
			},
			Attr_ReplicationSites: {
				Computed:    true,
				Description: "Indicates the replication sites of the volume group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
			Attr_Status: {
				Computed:    true,
				Description: "The status of the volume group.",
				Type:        schema.TypeString,
			},
			Attr_StatusDescriptionErrors: {
				Computed:    true,
				Description: "The status details of the volume group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Key: {
							Computed:    true,
							Description: "The volume group error key.",
							Type:        schema.TypeString,
						},
						Attr_Message: {
							Computed:    true,
							Description: "The failure message providing more details about the error key.",
							Type:        schema.TypeString,
						},
						Attr_VolumeIDs: {
							Computed:    true,
							Description: "List of volume IDs, which failed to be added/removed to/from the volume group, with the given error.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Type:        schema.TypeList,
						},
					},
				},
				Type: schema.TypeSet,
			},
			Attr_StoragePool: {
				Computed:    true,
				Description: "Indicates the storage pool of the volume group",
				Type:        schema.TypeString,
			},
			Attr_VolumeIDs: {
				Computed:    true,
				Description: "List of volume IDs, member of volume group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
			Attr_VolumeGroupName: {
				Computed:    true,
				Description: "The name of the volume group.",
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceIBMPIVolumeGroupDetailsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	vgClient := instance.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)
	vgData, err := vgClient.GetDetails(d.Get(Arg_VolumeGroupID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*vgData.ID)
	d.Set(Attr_Auxiliary, vgData.Auxiliary)
	d.Set(Attr_ConsistencyGroupName, vgData.ConsistencyGroupName)
	d.Set(Attr_ReplicationStatus, vgData.ReplicationStatus)
	if len(vgData.ReplicationSites) > 0 {
		d.Set(Attr_ReplicationSites, vgData.ReplicationSites)
	}
	d.Set(Attr_Status, vgData.Status)
	if vgData.StatusDescription != nil {
		d.Set(Attr_StatusDescriptionErrors, flattenVolumeGroupStatusDescription(vgData.StatusDescription.Errors))
	}
	d.Set(Attr_StoragePool, vgData.StoragePool)
	d.Set(Attr_VolumeIDs, vgData.VolumeIDs)
	d.Set(Attr_VolumeGroupName, vgData.Name)

	return nil
}

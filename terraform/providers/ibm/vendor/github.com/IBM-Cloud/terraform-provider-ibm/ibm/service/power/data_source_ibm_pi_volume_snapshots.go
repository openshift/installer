// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIVolumeSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeSnapshotsRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_VolumesSnapshots: {
				Computed:    true,
				Description: "The list of volume snapshots.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_CreationDate: {
							Computed:    true,
							Description: "The date and time when the volume snapshot was created.",
							Type:        schema.TypeString,
						},
						Attr_CRN: {
							Computed:    true,
							Description: "The CRN of the volume snapshot.",
							Type:        schema.TypeString,
						},
						Attr_ID: {
							Computed:    true,
							Description: "The snapshot UUID.",
							Type:        schema.TypeString,
						},
						Attr_Name: {
							Computed:    true,
							Description: "The volume snapshot name.",
							Type:        schema.TypeString,
						},
						Attr_Size: {
							Computed:    true,
							Description: "The size of the volume snapshot, in gibibytes (GiB).",
							Type:        schema.TypeFloat,
						},
						Attr_Status: {
							Computed:    true,
							Description: "The status for the volume snapshot.",
							Type:        schema.TypeString,
						},
						Attr_UpdatedDate: {
							Computed:    true,
							Description: "The date and time when the volume snapshot was last updated.",
							Type:        schema.TypeString,
						},
						Attr_VolumeID: {
							Computed:    true,
							Description: "The volume UUID associated with the snapshot.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeSet,
			},
		},
	}
}

func dataSourceIBMPIVolumeSnapshotsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	client := instance.NewIBMPISnapshotClient(ctx, sess, cloudInstanceID)
	snapshots, err := client.V1VolumeSnapshotsGetall()
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set(Attr_VolumesSnapshots, flattenSnapshotsV1(snapshots.VolumeSnapshots))
	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)

	return nil
}

func flattenSnapshotsV1(snapshotList []*models.SnapshotV1) []map[string]interface{} {
	snapshots := make([]map[string]interface{}, 0, len(snapshotList))
	for _, snap := range snapshotList {
		snapshot := map[string]interface{}{
			Attr_CreationDate: snap.CreationDate.String(),
			Attr_CRN:          snap.Crn,
			Attr_ID:           *snap.ID,
			Attr_Name:         *snap.Name,
			Attr_Size:         *snap.Size,
			Attr_Status:       *snap.Status,
			Attr_UpdatedDate:  snap.UpdatedDate.String(),
			Attr_VolumeID:     *snap.VolumeID,
		}
		snapshots = append(snapshots, snapshot)
	}
	return snapshots
}

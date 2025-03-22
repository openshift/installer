// Copyright IBM Corp. 2024 All Rights Reserved.
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

func DataSourceIBMPIVolumeSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeSnapshotRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_VolumeSnapshotID: {
				Description:  "The volume snapshot id.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
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
	}
}

func dataSourceIBMPIVolumeSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	snapshotID := d.Get(Arg_VolumeSnapshotID).(string)

	client := instance.NewIBMPISnapshotClient(ctx, sess, cloudInstanceID)
	snapshot, err := client.V1VolumeSnapshotsGet(snapshotID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*snapshot.ID)
	d.Set(Attr_CreationDate, snapshot.CreationDate.String())
	d.Set(Attr_CRN, snapshot.Crn)
	d.Set(Attr_Name, *snapshot.Name)
	d.Set(Attr_Size, *snapshot.Size)
	d.Set(Attr_Status, snapshot.Status)
	d.Set(Attr_UpdatedDate, snapshot.UpdatedDate.String())
	d.Set(Attr_VolumeID, *snapshot.VolumeID)
	return nil
}

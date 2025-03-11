// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIPVMSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPISnapshotRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_InstanceName: {
				Description:  "The unique identifier or name of the instance.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_PVMSnapshots: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Action: {
							Computed:    true,
							Description: "Action performed on the instance snapshot.",
							Type:        schema.TypeString,
						},
						Attr_CreationDate: {
							Computed:    true,
							Description: "Date of snapshot creation.",
							Type:        schema.TypeString,
						},
						Attr_CRN: {
							Computed:    true,
							Description: "The CRN of this resource.",
							Type:        schema.TypeString,
						},
						Attr_Description: {
							Computed:    true,
							Description: "The description of the snapshot.",
							Type:        schema.TypeString,
						},
						Attr_ID: {
							Computed:    true,
							Description: "The unique identifier of the Power Virtual Machine instance snapshot.",
							Type:        schema.TypeString,
						},
						Attr_LastUpdatedDate: {
							Computed:    true,
							Description: "Date of last update.",
							Type:        schema.TypeString,
						},
						Attr_Name: {
							Computed:    true,
							Description: "The name of the Power Virtual Machine instance snapshot.",
							Type:        schema.TypeString,
						},
						Attr_PercentComplete: {
							Computed:    true,
							Description: "The snapshot completion percentage.",
							Type:        schema.TypeInt,
						},
						Attr_Status: {
							Computed:    true,
							Description: "The status of the Power Virtual Machine instance snapshot.",
							Type:        schema.TypeString,
						},
						Attr_UserTags: {
							Computed:    true,
							Description: "List of user tags attached to the resource.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Type:        schema.TypeSet,
						},
						Attr_VolumeSnapshots: {
							Computed:    true,
							Description: "A map of volume snapshots included in the Power Virtual Machine instance snapshot.",
							Type:        schema.TypeMap,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPISnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	powerinstancename := d.Get(Arg_InstanceName).(string)
	snapshot := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
	snapshotData, err := snapshot.GetSnapShotVM(powerinstancename)

	if err != nil {
		return diag.FromErr(err)
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Attr_PVMSnapshots, flattenPVMSnapshotInstances(snapshotData.Snapshots, meta))

	return nil
}

func flattenPVMSnapshotInstances(list []*models.Snapshot, meta interface{}) []map[string]interface{} {
	log.Printf("Calling the flattenPVMSnapshotInstances call with list %d", len(list))
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			Attr_Action:          i.Action,
			Attr_CreationDate:    i.CreationDate.String(),
			Attr_Description:     i.Description,
			Attr_ID:              *i.SnapshotID,
			Attr_LastUpdatedDate: i.LastUpdateDate.String(),
			Attr_Name:            *i.Name,
			Attr_PercentComplete: i.PercentComplete,
			Attr_Status:          i.Status,
			Attr_VolumeSnapshots: i.VolumeSnapshots,
		}
		if i.Crn != "" {
			l[Attr_CRN] = i.Crn
			tags, err := flex.GetGlobalTagsUsingCRN(meta, string(i.Crn), "", UserTagType)
			if err != nil {
				log.Printf("Error on get of pi pvm snapshot (%s) user_tags: %s", *i.SnapshotID, err)
			}
			l[Attr_UserTags] = tags
		}
		result = append(result, l)
	}
	return result
}

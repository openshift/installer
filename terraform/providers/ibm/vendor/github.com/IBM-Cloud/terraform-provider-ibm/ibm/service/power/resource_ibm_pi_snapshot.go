// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMPISnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPISnapshotCreate,
		ReadContext:   resourceIBMPISnapshotRead,
		UpdateContext: resourceIBMPISnapshotUpdate,
		DeleteContext: resourceIBMPISnapshotDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourcePowerUserTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{

			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_Description: {
				Description: "Description of the PVM instance snapshot.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_InstanceName: {
				Description:  "The name of the instance you want to take a snapshot of.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_SnapShotName: {
				Description:  "The unique name of the snapshot.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_UserTags: {
				Computed:    true,
				Description: "The user tags attached to this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
			Arg_VolumeIDs: {
				Description:      "A list of volume IDs of the instance that will be part of the snapshot. If none are provided, then all the volumes of the instance will be part of the snapshot.",
				DiffSuppressFunc: flex.ApplyOnce,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Optional:         true,
				Set:              schema.HashString,
				Type:             schema.TypeSet,
			},

			// Attributes
			Attr_CreationDate: {
				Computed:    true,
				Description: "Creation date of the snapshot.",
				Type:        schema.TypeString,
			},
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			Attr_LastUpdateDate: {
				Computed:    true,
				Description: "The last updated date of the snapshot.",
				Type:        schema.TypeString,
			},
			Attr_SnapshotID: {
				Computed:    true,
				Description: "ID of the PVM instance snapshot.",
				Type:        schema.TypeString,
			},
			Attr_Status: {
				Computed:    true,
				Description: "Status of the PVM instance snapshot.",
				Type:        schema.TypeString,
			},
			Attr_VolumeSnapshots: {
				Computed:    true,
				Description: "A map of volume snapshots included in the PVM instance snapshot.",
				Type:        schema.TypeMap,
			},
		},
		DeprecationMessage: "Resource ibm_pi_snapshot is deprecated. Use `ibm_pi_instance_snapshot` resource instead.",
	}
}

func resourceIBMPISnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	instanceid := d.Get(Arg_InstanceName).(string)
	name := d.Get(Arg_SnapShotName).(string)
	volumeIDs := flex.ExpandStringList((d.Get(Arg_VolumeIDs).(*schema.Set)).List())

	var description string
	if v, ok := d.GetOk(Arg_Description); ok {
		description = v.(string)
	}

	client := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)

	snapshotBody := &models.SnapshotCreate{Name: &name, Description: description}

	if len(volumeIDs) > 0 {
		snapshotBody.VolumeIDs = volumeIDs
	} else {
		log.Printf("no volumeids provided. Will snapshot the entire instance")
	}

	if v, ok := d.GetOk(Arg_UserTags); ok {
		snapshotBody.UserTags = flex.FlattenSet(v.(*schema.Set))
	}

	snapshotResponse, err := client.CreatePvmSnapShot(instanceid, snapshotBody)
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *snapshotResponse.SnapshotID))

	piSnapClient := instance.NewIBMPISnapshotClient(ctx, sess, cloudInstanceID)
	_, err = isWaitForPIInstanceSnapshotAvailable(ctx, piSnapClient, *snapshotResponse.SnapshotID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if _, ok := d.GetOk(Arg_UserTags); ok {
		if snapshotResponse.Crn != "" {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, string(snapshotResponse.Crn), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi snapshot (%s) pi_user_tags during creation: %s", *snapshotResponse.SnapshotID, err)
			}
		}
	}

	return resourceIBMPISnapshotRead(ctx, d, meta)
}

func resourceIBMPISnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("Calling the Snapshot Read function post create")
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, snapshotID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	snapshot := instance.NewIBMPISnapshotClient(ctx, sess, cloudInstanceID)
	snapshotdata, err := snapshot.Get(snapshotID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(Arg_SnapShotName, snapshotdata.Name)
	d.Set(Attr_CreationDate, snapshotdata.CreationDate.String())
	if snapshotdata.Crn != "" {
		d.Set(Attr_CRN, snapshotdata.Crn)
		tags, err := flex.GetGlobalTagsUsingCRN(meta, string(snapshotdata.Crn), "", UserTagType)
		if err != nil {
			log.Printf("Error on get of pi snapshot (%s) pi_user_tags: %s", *snapshotdata.SnapshotID, err)
		}
		d.Set(Arg_UserTags, tags)
	}
	d.Set(Attr_LastUpdateDate, snapshotdata.LastUpdateDate.String())
	d.Set(Attr_SnapshotID, *snapshotdata.SnapshotID)
	d.Set(Attr_Status, snapshotdata.Status)
	d.Set(Attr_VolumeSnapshots, snapshotdata.VolumeSnapshots)

	return nil
}

func resourceIBMPISnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("Calling the IBM Power Snapshot update call")
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, snapshotID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := instance.NewIBMPISnapshotClient(ctx, sess, cloudInstanceID)

	if d.HasChange(Arg_SnapShotName) || d.HasChange(Arg_Description) {
		name := d.Get(Arg_SnapShotName).(string)
		description := d.Get(Arg_Description).(string)
		snapshotBody := &models.SnapshotUpdate{Name: name, Description: description}

		_, err := client.Update(snapshotID, snapshotBody)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = isWaitForPIInstanceSnapshotAvailable(ctx, client, snapshotID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange(Arg_UserTags) {
		if crn, ok := d.GetOk(Attr_CRN); ok {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, crn.(string), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi snapshot (%s) pi_user_tags: %s", snapshotID, err)
			}
		}
	}

	return resourceIBMPISnapshotRead(ctx, d, meta)
}

func resourceIBMPISnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, snapshotID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := instance.NewIBMPISnapshotClient(ctx, sess, cloudInstanceID)
	snapshot, err := client.Get(snapshotID)
	if err != nil {
		// snapshot does not exist
		d.SetId("")
		return nil
	}

	log.Printf("The snapshot  to be deleted is in the following state .. %s", snapshot.Status)

	err = client.Delete(snapshotID)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = isWaitForPIInstanceSnapshotDeleted(ctx, client, snapshotID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

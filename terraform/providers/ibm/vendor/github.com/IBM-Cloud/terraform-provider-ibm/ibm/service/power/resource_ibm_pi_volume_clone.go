// Copyright IBM Corp. 2024 All Rights Reserved.
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMPIVolumeClone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIVolumeCloneCreate,
		ReadContext:   resourceIBMPIVolumeCloneRead,
		DeleteContext: resourceIBMPIVolumeCloneDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_TargetStorageTier: {
				Description: "The storage tier for the cloned volume(s).",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_ReplicationEnabled: {
				Description: "Indicates whether the cloned volume should have replication enabled. If no value is provided, it will default to the replication status of the source volume(s).",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeBool,
			},
			Arg_UserTags: {
				Description: "The user tags attached to this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				ForceNew:    true,
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
			Arg_VolumeCloneName: {
				Description:  "The base name of the newly cloned volume(s).",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_VolumeIDs: {
				Description: "List of volumes to be cloned.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				ForceNew:    true,
				Required:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},

			// Attributes
			Attr_ClonedVolumes: {
				Computed:    true,
				Description: "The List of cloned volumes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_CloneVolumeID: {
							Computed:    true,
							Description: "The ID of the newly cloned volume.",
							Type:        schema.TypeString,
						},
						Attr_SourceVolumeID: {
							Computed:    true,
							Description: "The ID of the source volume.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_FailureReason: {
				Computed:    true,
				Description: "The reason the clone volumes task has failed.",
				Type:        schema.TypeString,
			},
			Attr_PercentComplete: {
				Computed:    true,
				Description: "The completion percentage of the volume clone task.",
				Type:        schema.TypeInt,
			},
			Attr_Status: {
				Computed:    true,
				Description: "The status of the volume clone task.",
				Type:        schema.TypeString,
			},
			Attr_TaskID: {
				Computed:    true,
				Description: "The ID of the volume clone task.",
				Type:        schema.TypeString,
			},
		},
	}
}

func resourceIBMPIVolumeCloneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	vcName := d.Get(Arg_VolumeCloneName).(string)
	volids := flex.ExpandStringList((d.Get(Arg_VolumeIDs).(*schema.Set)).List())

	body := &models.VolumesCloneAsyncRequest{
		Name:      &vcName,
		VolumeIDs: volids,
	}

	if v, ok := d.GetOk(Arg_TargetStorageTier); ok {
		body.TargetStorageTier = v.(string)
	}

	if !d.GetRawConfig().GetAttr(Arg_ReplicationEnabled).IsNull() {
		body.TargetReplicationEnabled = flex.PtrToBool(d.Get(Arg_ReplicationEnabled).(bool))
	}

	if v, ok := d.GetOk(Arg_UserTags); ok {
		body.UserTags = flex.FlattenSet(v.(*schema.Set))
	}

	client := instance.NewIBMPICloneVolumeClient(ctx, sess, cloudInstanceID)
	volClone, err := client.Create(body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *volClone.CloneTaskID))

	_, err = isWaitForIBMPIVolumeCloneCompletion(ctx, client, *volClone.CloneTaskID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMPIVolumeCloneRead(ctx, d, meta)
}

func resourceIBMPIVolumeCloneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, vcTaskID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := instance.NewIBMPICloneVolumeClient(ctx, sess, cloudInstanceID)
	volCloneTask, err := client.Get(vcTaskID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(Attr_FailureReason, volCloneTask.FailedReason)
	if volCloneTask.PercentComplete != nil {
		d.Set(Attr_PercentComplete, *volCloneTask.PercentComplete)
	}
	if volCloneTask.Status != nil {
		d.Set(Attr_Status, *volCloneTask.Status)
	}
	d.Set(Attr_TaskID, vcTaskID)
	d.Set(Attr_ClonedVolumes, flattenClonedVolumes(volCloneTask.ClonedVolumes))

	return nil
}

func resourceIBMPIVolumeCloneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// There is no delete or unset concept for volume clone
	d.SetId("")
	return nil
}

func flattenClonedVolumes(list []*models.ClonedVolume) (cloneVolumes []map[string]interface{}) {
	if list != nil {
		cloneVolumes := make([]map[string]interface{}, len(list))
		for i, data := range list {
			l := map[string]interface{}{
				Attr_CloneVolumeID:  data.ClonedVolumeID,
				Attr_SourceVolumeID: data.SourceVolumeID,
			}
			cloneVolumes[i] = l
		}
		return cloneVolumes
	}
	return
}

func isWaitForIBMPIVolumeCloneCompletion(ctx context.Context, client *instance.IBMPICloneVolumeClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume clone (%s) to be completed.", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Creating},
		Target:     []string{State_Completed},
		Refresh:    isIBMPIVolumeCloneRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIVolumeCloneRefreshFunc(client *instance.IBMPICloneVolumeClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		volClone, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if *volClone.Status == State_Completed {
			return volClone, State_Completed, nil
		}

		return volClone, State_Creating, nil
	}
}

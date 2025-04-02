// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_volume_groups"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMPIVolumeGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIVolumeGroupCreate,
		ReadContext:   resourceIBMPIVolumeGroupRead,
		UpdateContext: resourceIBMPIVolumeGroupUpdate,
		DeleteContext: resourceIBMPIVolumeGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_ConsistencyGroupName: {
				ConflictsWith: []string{Arg_VolumeGroupName},
				Description:   "The name of consistency group at storage controller level",
				Optional:      true,
				Type:          schema.TypeString,
			},
			Arg_VolumeGroupName: {
				ConflictsWith: []string{Arg_ConsistencyGroupName},
				Description:   "Volume Group Name to create",
				Optional:      true,
				Type:          schema.TypeString,
			},
			Arg_VolumeIDs: {
				Description: "List of volumes to add in volume group",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},

			// Attributes
			Attr_ConsistencyGroupName: {
				Computed:    true,
				Description: "Consistency Group Name if volume is a part of volume group",
				Type:        schema.TypeString,
			},
			Attr_ReplicationSites: {
				Computed:    true,
				Description: "Indicates the replication sites of the volume group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
			Attr_ReplicationStatus: {
				Computed:    true,
				Description: "Volume Group Replication Status",
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
							Description: "List of volume IDs, which failed to be added to or removed from the volume group, with the given error.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Type:        schema.TypeList,
						},
					},
				},
				Type: schema.TypeSet,
			},
			Attr_VolumeGroupID: {
				Computed:    true,
				Description: "Volume Group ID",
				Type:        schema.TypeString,
			},
			Attr_VolumeGroupStatus: {
				Computed:    true,
				Description: "Volume Group Status",
				Type:        schema.TypeString,
			},
		},
	}
}

func resourceIBMPIVolumeGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	vgName := d.Get(Arg_VolumeGroupName).(string)
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	body := &models.VolumeGroupCreate{
		Name: vgName,
	}

	volids := flex.ExpandStringList((d.Get(Arg_VolumeIDs).(*schema.Set)).List())
	body.VolumeIDs = volids

	if v, ok := d.GetOk(Arg_ConsistencyGroupName); ok {
		body.ConsistencyGroupName = v.(string)
	}

	client := instance.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)
	vg, err := client.CreateVolumeGroup(body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *vg.ID))

	_, err = isWaitForIBMPIVolumeGroupAvailable(ctx, client, *vg.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMPIVolumeGroupRead(ctx, d, meta)
}

func resourceIBMPIVolumeGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, vgID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := instance.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)

	vg, err := client.GetDetails(vgID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(Arg_VolumeGroupName, vg.Name)
	d.Set(Arg_VolumeIDs, vg.VolumeIDs)
	d.Set(Attr_ConsistencyGroupName, vg.ConsistencyGroupName)
	d.Set(Attr_ReplicationSites, vg.ReplicationSites)
	d.Set(Attr_ReplicationStatus, vg.ReplicationStatus)
	if vg.StatusDescription != nil {
		d.Set(Attr_StatusDescriptionErrors, flattenVolumeGroupStatusDescription(vg.StatusDescription.Errors))
	}
	d.Set(Attr_VolumeGroupID, vg.ID)
	d.Set(Attr_VolumeGroupStatus, vg.Status)

	return nil
}

func resourceIBMPIVolumeGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, vgID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := instance.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)
	if d.HasChanges(Arg_VolumeIDs) {
		old, new := d.GetChange(Arg_VolumeIDs)
		oldList := old.(*schema.Set)
		newList := new.(*schema.Set)
		body := &models.VolumeGroupUpdate{
			AddVolumes:    flex.ExpandStringList(newList.Difference(oldList).List()),
			RemoveVolumes: flex.ExpandStringList(oldList.Difference(newList).List()),
		}
		err := client.UpdateVolumeGroup(vgID, body)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForIBMPIVolumeGroupAvailable(ctx, client, vgID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMPIVolumeGroupRead(ctx, d, meta)
}
func resourceIBMPIVolumeGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, vgID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := instance.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)

	volids := flex.ExpandStringList((d.Get(Arg_VolumeIDs).(*schema.Set)).List())
	if len(volids) > 0 {
		body := &models.VolumeGroupUpdate{
			RemoveVolumes: volids,
		}
		err = client.UpdateVolumeGroup(vgID, body)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForIBMPIVolumeGroupAvailable(ctx, client, vgID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = client.DeleteVolumeGroup(vgID)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = isWaitForIBMPIVolumeGroupDeleted(ctx, client, vgID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
func isWaitForIBMPIVolumeGroupAvailable(ctx context.Context, client *instance.IBMPIVolumeGroupClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume Group (%s) to be available.", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Retry, State_Creating},
		Target:     []string{State_Available},
		Refresh:    isIBMPIVolumeGroupRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIVolumeGroupRefreshFunc(client *instance.IBMPIVolumeGroupClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vg, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if vg.Status == State_Available {
			return vg, State_Available, nil
		}

		return vg, State_Creating, nil
	}
}

func isWaitForIBMPIVolumeGroupDeleted(ctx context.Context, client *instance.IBMPIVolumeGroupClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Deleting, State_Updating},
		Target:     []string{State_Deleted},
		Refresh:    isIBMPIVolumeGroupDeleteRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}
	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIVolumeGroupDeleteRefreshFunc(client *instance.IBMPIVolumeGroupClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vg, err := client.Get(id)
		if err != nil {
			uErr := errors.Unwrap(err)
			switch uErr.(type) {
			case *p_cloud_volume_groups.PcloudVolumegroupsGetNotFound:
				log.Printf("[DEBUG] volume-group does not exist while deleteing %v", err)
				return vg, State_Deleted, nil
			}
			return nil, "", err
		}
		if vg == nil {
			return vg, State_Deleted, nil
		}
		return vg, State_Deleting, nil
	}
}

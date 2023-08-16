// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_volume_groups"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud Instance ID - This is the service_instance_id.",
			},
			PIVolumeGroupName: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Volume Group Name to create",
				ConflictsWith: []string{PIVolumeGroupConsistencyGroupName},
			},
			PIVolumeGroupConsistencyGroupName: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The name of consistency group at storage controller level",
				ConflictsWith: []string{PIVolumeGroupName},
			},
			PIVolumeGroupsVolumeIds: {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of volumes to add in volume group",
			},

			// Computed Attributes
			"volume_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume Group ID",
			},
			"volume_group_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume Group Status",
			},
			"replication_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume Group Replication Status",
			},
			"status_description_errors": vgStatusDescriptionErrors(),
			"consistency_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Consistency Group Name if volume is a part of volume group",
			},
		},
	}
}

func resourceIBMPIVolumeGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	vgName := d.Get(PIVolumeGroupName).(string)
	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	body := &models.VolumeGroupCreate{
		Name: vgName,
	}

	volids := flex.ExpandStringList((d.Get(PIVolumeGroupsVolumeIds).(*schema.Set)).List())
	body.VolumeIDs = volids

	if v, ok := d.GetOk(PIVolumeGroupConsistencyGroupName); ok {
		body.ConsistencyGroupName = v.(string)
	}

	client := st.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)
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

	client := st.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)

	vg, err := client.GetDetails(vgID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("volume_group_id", vg.ID)
	d.Set("volume_group_status", vg.Status)
	d.Set("consistency_group_name", vg.ConsistencyGroupName)
	d.Set("replication_status", vg.ReplicationStatus)
	d.Set(PIVolumeGroupName, vg.Name)
	d.Set(PIVolumeGroupsVolumeIds, vg.VolumeIDs)
	d.Set("status_description_errors", flattenVolumeGroupStatusDescription(vg.StatusDescription.Errors))

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

	client := st.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)
	if d.HasChanges(PIVolumeGroupsVolumeIds) {
		old, new := d.GetChange(PIVolumeGroupsVolumeIds)
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

	client := st.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)

	volids := flex.ExpandStringList((d.Get(PIVolumeGroupsVolumeIds).(*schema.Set)).List())
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
func isWaitForIBMPIVolumeGroupAvailable(ctx context.Context, client *st.IBMPIVolumeGroupClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume Group (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PIVolumeProvisioning},
		Target:     []string{helpers.PIVolumeProvisioningDone},
		Refresh:    isIBMPIVolumeGroupRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIVolumeGroupRefreshFunc(client *st.IBMPIVolumeGroupClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vg, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if vg.Status == "available" {
			return vg, helpers.PIVolumeProvisioningDone, nil
		}

		return vg, helpers.PIVolumeProvisioning, nil
	}
}

func isWaitForIBMPIVolumeGroupDeleted(ctx context.Context, client *st.IBMPIVolumeGroupClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting", "updating"},
		Target:     []string{"deleted"},
		Refresh:    isIBMPIVolumeGroupDeleteRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}
	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIVolumeGroupDeleteRefreshFunc(client *st.IBMPIVolumeGroupClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vg, err := client.Get(id)
		if err != nil {
			uErr := errors.Unwrap(err)
			switch uErr.(type) {
			case *p_cloud_volume_groups.PcloudVolumegroupsGetNotFound:
				log.Printf("[DEBUG] volume-group does not exist while deleteing %v", err)
				return vg, "deleted", nil
			}
			return nil, "", err
		}
		if vg == nil {
			return vg, "deleted", nil
		}
		return vg, "deleting", nil
	}
}

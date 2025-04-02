// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_volumes"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMPIVolumeAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIVolumeAttachCreate,
		ReadContext:   resourceIBMPIVolumeAttachRead,
		DeleteContext: resourceIBMPIVolumeAttachDelete,
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
			Arg_InstanceID: {
				Description: "PI Instance Id",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			Arg_VolumeID: {
				Description: "Id of the volume to attach. Note these volumes should have been created",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},

			// Attribute
			Attr_Status: {
				Computed:    true,
				Description: "The status of the volume.",
				Type:        schema.TypeString,
			},
		},
	}
}

func resourceIBMPIVolumeAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	pvmInstanceID := d.Get(Arg_InstanceID).(string)
	volumeID := d.Get(Arg_VolumeID).(string)

	volClient := instance.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
	volinfo, err := volClient.Get(volumeID)
	if err != nil {
		return diag.FromErr(err)
	}

	if volinfo.State == State_Available || *volinfo.Shareable {
		log.Printf(" In the current state the volume can be attached to the instance ")
	}

	if volinfo.State == State_InUse && *volinfo.Shareable {

		log.Printf("Volume State /Status is  permitted and hence attaching the volume to the instance")
	}

	if volinfo.State == State_InUse && !*volinfo.Shareable {
		return diag.Errorf("the volume cannot be attached in the current state. The volume must be in the *available* state. No other states are permissible")
	}

	err = volClient.Attach(pvmInstanceID, volumeID)
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cloudInstanceID, pvmInstanceID, *volinfo.VolumeID))

	_, err = isWaitForIBMPIVolumeAttachAvailable(ctx, volClient, *volinfo.VolumeID, pvmInstanceID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMPIVolumeAttachRead(ctx, d, meta)
}

func resourceIBMPIVolumeAttachRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	ids, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, pvmInstanceID, volumeID := ids[0], ids[1], ids[2]

	client := instance.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)

	vol, err := client.CheckVolumeAttach(pvmInstanceID, volumeID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(Attr_Status, vol.State)
	return nil
}

func resourceIBMPIVolumeAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	ids, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, pvmInstanceID, volumeID := ids[0], ids[1], ids[2]
	client := instance.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)

	log.Printf("the id of the volume to detach is %s ", volumeID)

	err = client.Detach(pvmInstanceID, volumeID)
	if err != nil {
		uErr := errors.Unwrap(err)
		switch uErr.(type) {
		case *p_cloud_volumes.PcloudCloudinstancesVolumesGetNotFound:
			log.Printf("[DEBUG] volume does not exist while detaching %v", err)
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] volume detach failed %v", err)
		return diag.FromErr(err)
	}

	_, err = isWaitForIBMPIVolumeDetach(ctx, client, volumeID, pvmInstanceID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	// wait for power volume states to be back as available. if it's attached it will be in-use
	d.SetId("")
	return nil
}

func isWaitForIBMPIVolumeAttachAvailable(ctx context.Context, client *instance.IBMPIVolumeClient, id, pvmInstanceID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume (%s) to be available for attachment", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Retry, State_Creating},
		Target:     []string{State_InUse},
		Refresh:    isIBMPIVolumeAttachRefreshFunc(client, id, pvmInstanceID),
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIVolumeAttachRefreshFunc(client *instance.IBMPIVolumeClient, id, pvmInstanceID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vol, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if vol.State == State_InUse && flex.StringContains(vol.PvmInstanceIDs, pvmInstanceID) {
			return vol, State_InUse, nil
		}

		return vol, State_Creating, nil
	}
}

func isWaitForIBMPIVolumeDetach(ctx context.Context, client *instance.IBMPIVolumeClient, id, pvmInstanceID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume (%s) to be available after detachment", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Detaching, State_Deleting},
		Target:     []string{State_Available},
		Refresh:    isIBMPIVolumeDetachRefreshFunc(client, id, pvmInstanceID),
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIVolumeDetachRefreshFunc(client *instance.IBMPIVolumeClient, id, pvmInstanceID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vol, err := client.Get(id)
		if err != nil {
			uErr := errors.Unwrap(err)
			switch uErr.(type) {
			case *p_cloud_volumes.PcloudCloudinstancesVolumesGetNotFound:
				log.Printf("[DEBUG] volume does not exist while detaching %v", err)
				return vol, State_Available, nil
			}
			return nil, "", err
		}

		// Check if Instance ID is in the Volume's Instance list
		// Also validate the Volume state is 'available' when it is not Sharable
		// In case of Sharable Volume it can be `in-use` state
		if !flex.StringContains(vol.PvmInstanceIDs, pvmInstanceID) &&
			(*vol.Shareable || (!*vol.Shareable && vol.State == State_Available)) {
			return vol, State_Available, nil
		}

		return vol, State_Detaching, nil
	}
}

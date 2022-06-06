// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
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
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_volumes"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: " Cloud Instance ID - This is the service_instance_id.",
			},

			helpers.PIVolumeId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the volume to attach. Note these volumes should have been created",
			},

			helpers.PIInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "PI Instance Id",
			},

			// Computed Attribute
			helpers.PIVolumeAttachStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIBMPIVolumeAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	volumeID := d.Get(helpers.PIVolumeId).(string)
	pvmInstanceID := d.Get(helpers.PIInstanceId).(string)
	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	volClient := st.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
	volinfo, err := volClient.Get(volumeID)
	if err != nil {
		return diag.FromErr(err)
	}

	if volinfo.State == "available" || *volinfo.Shareable {
		log.Printf(" In the current state the volume can be attached to the instance ")
	}

	if volinfo.State == "in-use" && *volinfo.Shareable {

		log.Printf("Volume State /Status is  permitted and hence attaching the volume to the instance")
	}

	if volinfo.State == helpers.PIVolumeAllowableAttachStatus && !*volinfo.Shareable {
		return diag.Errorf("the volume cannot be attached in the current state. The volume must be in the *available* state. No other states are permissible")
	}

	err = volClient.Attach(pvmInstanceID, volumeID)
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cloudInstanceID, pvmInstanceID, *volinfo.VolumeID))

	_, err = isWaitForIBMPIVolumeAttachAvailable(ctx, volClient, *volinfo.VolumeID, cloudInstanceID, pvmInstanceID, d.Timeout(schema.TimeoutCreate))
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

	client := st.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)

	vol, err := client.CheckVolumeAttach(pvmInstanceID, volumeID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(helpers.PIVolumeAttachStatus, vol.State)
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
	client := st.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)

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

	_, err = isWaitForIBMPIVolumeDetach(ctx, client, volumeID, cloudInstanceID, pvmInstanceID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	// wait for power volume states to be back as available. if it's attached it will be in-use
	d.SetId("")
	return nil
}

func isWaitForIBMPIVolumeAttachAvailable(ctx context.Context, client *st.IBMPIVolumeClient, id, cloudInstanceID, pvmInstanceID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume (%s) to be available for attachment", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PIVolumeProvisioning},
		Target:     []string{helpers.PIVolumeAllowableAttachStatus},
		Refresh:    isIBMPIVolumeAttachRefreshFunc(client, id, cloudInstanceID, pvmInstanceID),
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIVolumeAttachRefreshFunc(client *st.IBMPIVolumeClient, id, cloudInstanceID, pvmInstanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vol, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if vol.State == "in-use" && flex.StringContains(vol.PvmInstanceIDs, pvmInstanceID) {
			return vol, helpers.PIVolumeAllowableAttachStatus, nil
		}

		return vol, helpers.PIVolumeProvisioning, nil
	}
}

func isWaitForIBMPIVolumeDetach(ctx context.Context, client *st.IBMPIVolumeClient, id, cloudInstanceID, pvmInstanceID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume (%s) to be available after detachment", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"detaching", helpers.PowerVolumeAttachDeleting},
		Target:     []string{helpers.PIVolumeProvisioningDone},
		Refresh:    isIBMPIVolumeDetachRefreshFunc(client, id, cloudInstanceID, pvmInstanceID),
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIVolumeDetachRefreshFunc(client *st.IBMPIVolumeClient, id, cloudInstanceID, pvmInstanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vol, err := client.Get(id)
		if err != nil {
			uErr := errors.Unwrap(err)
			switch uErr.(type) {
			case *p_cloud_volumes.PcloudCloudinstancesVolumesGetNotFound:
				log.Printf("[DEBUG] volume does not exist while detaching %v", err)
				return vol, helpers.PIVolumeProvisioningDone, nil
			}
			return nil, "", err
		}

		// Check if Instance ID is in the Volume's Instance list
		// Also validate the Volume state is 'available' when it is not Sharable
		// In case of Sharable Volume it can be `in-use` state
		if !flex.StringContains(vol.PvmInstanceIDs, pvmInstanceID) &&
			(*vol.Shareable || (!*vol.Shareable && vol.State == "available")) {
			return vol, helpers.PIVolumeProvisioningDone, nil
		}

		return vol, "detaching", nil
	}
}

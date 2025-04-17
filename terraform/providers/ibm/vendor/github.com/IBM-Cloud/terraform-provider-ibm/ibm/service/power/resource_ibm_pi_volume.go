// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_volumes"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func ResourceIBMPIVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIVolumeCreate,
		ReadContext:   resourceIBMPIVolumeRead,
		UpdateContext: resourceIBMPIVolumeUpdate,
		DeleteContext: resourceIBMPIVolumeDelete,
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
			helpers.PIVolumeName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Volume Name to create",
			},
			helpers.PIVolumeShareable: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to indicate if the volume can be shared across multiple instances?",
			},
			helpers.PIVolumeSize: {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "Size of the volume in GB",
			},
			helpers.PIVolumeType: {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validate.ValidateAllowedStringValues([]string{"ssd", "standard", "tier1", "tier3"}),
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Type of Disk, required if pi_affinity_policy and pi_volume_pool not provided, otherwise ignored",
			},
			helpers.PIVolumePool: {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Volume pool where the volume will be created; if provided then pi_volume_type and pi_affinity_policy values will be ignored",
			},
			PIAffinityPolicy: {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Affinity policy for data volume being created; ignored if pi_volume_pool provided; for policy affinity requires one of pi_affinity_instance or pi_affinity_volume to be specified; for policy anti-affinity requires one of pi_anti_affinity_instances or pi_anti_affinity_volumes to be specified",
				ValidateFunc:     validate.InvokeValidator("ibm_pi_volume", PIAffinityPolicy),
			},
			PIAffinityVolume: {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Volume (ID or Name) to base volume affinity policy against; required if requesting affinity and pi_affinity_instance is not provided",
				ConflictsWith:    []string{PIAffinityInstance},
			},
			PIAffinityInstance: {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "PVM Instance (ID or Name) to base volume affinity policy against; required if requesting affinity and pi_affinity_volume is not provided",
				ConflictsWith:    []string{PIAffinityVolume},
			},
			PIAntiAffinityVolumes: {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "List of volumes to base volume anti-affinity policy against; required if requesting anti-affinity and pi_anti_affinity_instances is not provided",
				ConflictsWith:    []string{PIAntiAffinityInstances},
			},
			PIAntiAffinityInstances: {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "List of pvmInstances to base volume anti-affinity policy against; required if requesting anti-affinity and pi_anti_affinity_volumes is not provided",
				ConflictsWith:    []string{PIAntiAffinityVolumes},
			},
			helpers.PIReplicationEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates if the volume should be replication enabled or not",
			},

			// Computed Attributes
			"volume_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume ID",
			},
			"volume_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume status",
			},

			"delete_on_termination": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Should the volume be deleted during termination",
			},
			"wwn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "WWN Of the volume",
			},
			"auxiliary": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "true if volume is auxiliary otherwise false",
			},
			"consistency_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Consistency Group Name if volume is a part of volume group",
			},
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume Group ID",
			},
			"replication_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Replication type(metro,global)",
			},
			"replication_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Replication status of a volume",
			},
			"mirroring_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mirroring state for replication enabled volume",
			},
			"primary_role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates whether master/aux volume is playing the primary role",
			},
			"auxiliary_volume_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates auxiliary volume name",
			},
			"master_volume_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates master volume name",
			},
		},
	}
}
func ResourceIBMPIVolumeValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "pi_affinity",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "affinity, anti-affinity"})
	ibmPIVolumeResourceValidator := validate.ResourceValidator{
		ResourceName: "ibm_pi_volume",
		Schema:       validateSchema}
	return &ibmPIVolumeResourceValidator
}

func resourceIBMPIVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get(helpers.PIVolumeName).(string)
	size := float64(d.Get(helpers.PIVolumeSize).(float64))
	var shared bool
	if v, ok := d.GetOk(helpers.PIVolumeShareable); ok {
		shared = v.(bool)
	}
	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	body := &models.CreateDataVolume{
		Name:      &name,
		Shareable: &shared,
		Size:      &size,
	}
	if v, ok := d.GetOk(helpers.PIVolumeType); ok {
		volType := v.(string)
		body.DiskType = volType
	}
	if v, ok := d.GetOk(helpers.PIVolumePool); ok {
		volumePool := v.(string)
		body.VolumePool = volumePool
	}
	if v, ok := d.GetOk(helpers.PIReplicationEnabled); ok {
		replicationEnabled := v.(bool)
		body.ReplicationEnabled = &replicationEnabled
	}
	if ap, ok := d.GetOk(PIAffinityPolicy); ok {
		policy := ap.(string)
		body.AffinityPolicy = &policy

		if policy == "affinity" {
			if av, ok := d.GetOk(PIAffinityVolume); ok {
				afvol := av.(string)
				body.AffinityVolume = &afvol
			}
			if ai, ok := d.GetOk(PIAffinityInstance); ok {
				afins := ai.(string)
				body.AffinityPVMInstance = &afins
			}
		} else {
			if avs, ok := d.GetOk(PIAntiAffinityVolumes); ok {
				afvols := flex.ExpandStringList(avs.([]interface{}))
				body.AntiAffinityVolumes = afvols
			}
			if ais, ok := d.GetOk(PIAntiAffinityInstances); ok {
				afinss := flex.ExpandStringList(ais.([]interface{}))
				body.AntiAffinityPVMInstances = afinss
			}
		}

	}

	client := st.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
	vol, err := client.CreateVolume(body)
	if err != nil {
		return diag.FromErr(err)
	}

	volumeid := *vol.VolumeID
	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, volumeid))

	_, err = isWaitForIBMPIVolumeAvailable(ctx, client, volumeid, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMPIVolumeRead(ctx, d, meta)
}

func resourceIBMPIVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, volumeID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := st.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)

	vol, err := client.Get(volumeID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set(helpers.PIVolumeName, vol.Name)
	d.Set(helpers.PIVolumeSize, vol.Size)
	if vol.Shareable != nil {
		d.Set(helpers.PIVolumeShareable, vol.Shareable)
	}
	d.Set(helpers.PIVolumeType, vol.DiskType)
	d.Set(helpers.PIVolumePool, vol.VolumePool)
	d.Set("volume_status", vol.State)
	if vol.VolumeID != nil {
		d.Set("volume_id", vol.VolumeID)
	}
	d.Set(helpers.PIReplicationEnabled, vol.ReplicationEnabled)
	d.Set("auxiliary", vol.Auxiliary)
	d.Set("consistency_group_name", vol.ConsistencyGroupName)
	d.Set("group_id", vol.GroupID)
	d.Set("replication_type", vol.ReplicationType)
	d.Set("replication_status", vol.ReplicationStatus)
	d.Set("mirroring_state", vol.MirroringState)
	d.Set("primary_role", vol.PrimaryRole)
	d.Set("master_volume_name", vol.MasterVolumeName)
	d.Set("auxiliary_volume_name", vol.AuxVolumeName)
	if vol.DeleteOnTermination != nil {
		d.Set("delete_on_termination", vol.DeleteOnTermination)
	}
	d.Set("wwn", vol.Wwn)
	d.Set(helpers.PICloudInstanceId, cloudInstanceID)

	return nil
}

func resourceIBMPIVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, volumeID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := st.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
	name := d.Get(helpers.PIVolumeName).(string)
	size := float64(d.Get(helpers.PIVolumeSize).(float64))
	var shareable bool
	if v, ok := d.GetOk(helpers.PIVolumeShareable); ok {
		shareable = v.(bool)
	}

	body := &models.UpdateVolume{
		Name:      &name,
		Shareable: &shareable,
		Size:      size,
	}
	volrequest, err := client.UpdateVolume(volumeID, body)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = isWaitForIBMPIVolumeAvailable(ctx, client, *volrequest.VolumeID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange(helpers.PIReplicationEnabled) {
		replicationEnabled := d.Get(helpers.PIReplicationEnabled).(bool)
		volActionBody := models.VolumeAction{
			ReplicationEnabled: &replicationEnabled,
		}

		err = client.VolumeAction(volumeID, &volActionBody)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForIBMPIVolumeAvailable(ctx, client, volumeID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMPIVolumeRead(ctx, d, meta)
}

func resourceIBMPIVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, volumeID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := st.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
	err = client.DeleteVolume(volumeID)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = isWaitForIBMPIVolumeDeleted(ctx, client, volumeID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func isWaitForIBMPIVolumeAvailable(ctx context.Context, client *st.IBMPIVolumeClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PIVolumeProvisioning},
		Target:     []string{helpers.PIVolumeProvisioningDone},
		Refresh:    isIBMPIVolumeRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIVolumeRefreshFunc(client *st.IBMPIVolumeClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vol, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if vol.State == "available" || vol.State == "in-use" {
			return vol, helpers.PIVolumeProvisioningDone, nil
		}

		return vol, helpers.PIVolumeProvisioning, nil
	}
}

func isWaitForIBMPIVolumeDeleted(ctx context.Context, client *st.IBMPIVolumeClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting", helpers.PIVolumeProvisioning},
		Target:     []string{"deleted"},
		Refresh:    isIBMPIVolumeDeleteRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}
	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIVolumeDeleteRefreshFunc(client *st.IBMPIVolumeClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vol, err := client.Get(id)
		if err != nil {
			uErr := errors.Unwrap(err)
			switch uErr.(type) {
			case *p_cloud_volumes.PcloudCloudinstancesVolumesGetNotFound:
				log.Printf("[DEBUG] volume does not exist %v", err)
				return vol, "deleted", nil
			}
			return nil, "", err
		}
		if vol == nil {
			return vol, "deleted", nil
		}
		return vol, "deleting", nil
	}
}

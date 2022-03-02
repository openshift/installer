// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

const (
	PIAffinityPolicy        = "pi_affinity_policy"
	PIAffinityVolume        = "pi_affinity_volume"
	PIAffinityInstance      = "pi_affinity_instance"
	PIAntiAffinityInstances = "pi_anti_affinity_instances"
	PIAntiAffinityVolumes   = "pi_anti_affinity_volumes"
)

func resourceIBMPIVolume() *schema.Resource {
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
				ValidateFunc:     validateAllowedStringValue([]string{"ssd", "standard", "tier1", "tier3"}),
				DiffSuppressFunc: applyOnce,
				Description:      "Type of Disk, required if pi_affinity_policy and pi_volume_pool not provided, otherwise ignored",
			},
			helpers.PIVolumePool: {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "Volume pool where the volume will be created; if provided then pi_volume_type and pi_affinity_policy values will be ignored",
			},
			PIAffinityPolicy: {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "Affinity policy for data volume being created; ignored if pi_volume_pool provided; for policy affinity requires one of pi_affinity_instance or pi_affinity_volume to be specified; for policy anti-affinity requires one of pi_anti_affinity_instances or pi_anti_affinity_volumes to be specified",
				ValidateFunc:     InvokeValidator("ibm_pi_volume", "pi_affinity"),
			},
			PIAffinityVolume: {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "Volume (ID or Name) to base volume affinity policy against; required if requesting affinity and pi_affinity_instance is not provided",
				ConflictsWith:    []string{PIAffinityInstance},
			},
			PIAffinityInstance: {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "PVM Instance (ID or Name) to base volume affinity policy against; required if requesting affinity and pi_affinity_volume is not provided",
				ConflictsWith:    []string{PIAffinityVolume},
			},
			PIAntiAffinityVolumes: {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: applyOnce,
				Description:      "List of volumes to base volume anti-affinity policy against; required if requesting anti-affinity and pi_anti_affinity_instances is not provided",
				ConflictsWith:    []string{PIAntiAffinityInstances},
			},
			PIAntiAffinityInstances: {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: applyOnce,
				Description:      "List of pvmInstances to base volume anti-affinity policy against; required if requesting anti-affinity and pi_anti_affinity_volumes is not provided",
				ConflictsWith:    []string{PIAntiAffinityVolumes},
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
		},
	}
}
func resourceIBMPIVolumeValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 0)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "pi_affinity",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "affinity, anti-affinity"})
	ibmPIVolumeResourceValidator := ResourceValidator{
		ResourceName: "ibm_pi_volume",
		Schema:       validateSchema}
	return &ibmPIVolumeResourceValidator
}

func resourceIBMPIVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
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
				afvols := expandStringList(avs.([]interface{}))
				body.AntiAffinityVolumes = afvols
			}
			if ais, ok := d.GetOk(PIAntiAffinityInstances); ok {
				afinss := expandStringList(ais.([]interface{}))
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
	sess, err := meta.(ClientSession).IBMPISession()
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
	if vol.DeleteOnTermination != nil {
		d.Set("delete_on_termination", vol.DeleteOnTermination)
	}
	d.Set("wwn", vol.Wwn)
	d.Set(helpers.PICloudInstanceId, cloudInstanceID)

	return nil
}

func resourceIBMPIVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
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

	return resourceIBMPIVolumeRead(ctx, d, meta)
}

func resourceIBMPIVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
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
			if strings.Contains(err.Error(), "Resource not found") {
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

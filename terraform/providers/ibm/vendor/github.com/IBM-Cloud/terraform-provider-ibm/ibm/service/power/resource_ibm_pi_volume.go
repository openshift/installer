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
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourcePowerUserTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_AffinityInstance: {
				ConflictsWith:    []string{Arg_AffinityVolume},
				Description:      "PVM Instance (ID or Name) to base volume affinity policy against; required if requesting 'affinity' and 'pi_affinity_volume' is not provided.",
				DiffSuppressFunc: flex.ApplyOnce,
				Optional:         true,
				Type:             schema.TypeString,
			},
			Arg_AffinityPolicy: {
				Description:      "Affinity policy for data volume being created; ignored if 'pi_volume_pool' provided; for policy 'affinity' requires one of 'pi_affinity_instance' or 'pi_affinity_volume' to be specified; for policy 'anti-affinity' requires one of 'pi_anti_affinity_instances' or 'pi_anti_affinity_volumes' to be specified; Allowable values: 'affinity', 'anti-affinity'.",
				DiffSuppressFunc: flex.ApplyOnce,
				Optional:         true,
				Type:             schema.TypeString,
				ValidateFunc:     validate.InvokeValidator("ibm_pi_volume", Arg_AffinityPolicy),
			},
			Arg_AffinityVolume: {
				ConflictsWith:    []string{Arg_AffinityInstance},
				Description:      "Volume (ID or Name) to base volume affinity policy against; required if requesting 'affinity' and 'pi_affinity_instance' is not provided.",
				DiffSuppressFunc: flex.ApplyOnce,
				Optional:         true,
				Type:             schema.TypeString,
			},
			Arg_AntiAffinityInstances: {
				ConflictsWith:    []string{Arg_AntiAffinityVolumes},
				Description:      "List of pvmInstances to base volume anti-affinity policy against; required if requesting 'anti-affinity' and 'pi_anti_affinity_volumes' is not provided.",
				DiffSuppressFunc: flex.ApplyOnce,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Optional:         true,
				Type:             schema.TypeList,
			},
			Arg_AntiAffinityVolumes: {
				ConflictsWith:    []string{Arg_AntiAffinityInstances},
				Description:      "List of volumes to base volume anti-affinity policy against; required if requesting 'anti-affinity' and 'pi_anti_affinity_instances' is not provided.",
				DiffSuppressFunc: flex.ApplyOnce,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Optional:         true,
				Type:             schema.TypeList,
			},
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_ReplicationEnabled: {
				Computed:    true,
				Description: "Indicates if the volume should be replication enabled or not.",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			Arg_ReplicationSites: {
				Description: "List of replication sites for volume replication.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				ForceNew:    true,
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
			Arg_UserTags: {
				Computed:    true,
				Description: "The user tags attached to this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
			Arg_VolumeName: {
				Description:  "The name of the volume.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_VolumePool: {
				Computed:         true,
				Description:      "Volume pool where the volume will be created; if provided then 'pi_affinity_policy' values will be ignored.",
				DiffSuppressFunc: flex.ApplyOnce,
				Optional:         true,
				Type:             schema.TypeString,
			},
			Arg_VolumeShareable: {
				Description: "If set to true, the volume can be shared across Power Systems Virtual Server instances. If set to false, you can attach it only to one instance.",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			Arg_VolumeSize: {
				Description:  "The size of the volume in GB.",
				Required:     true,
				Type:         schema.TypeFloat,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_VolumeType: {
				Computed:         true,
				Description:      "Type of disk, if diskType is not provided the disk type will default to 'tier3'",
				DiffSuppressFunc: flex.ApplyOnce,
				Optional:         true,
				Type:             schema.TypeString,
				ValidateFunc:     validate.ValidateAllowedStringValues([]string{"tier0", "tier1", "tier3", "tier5k"}),
			},

			// Attributes
			Attr_Auxiliary: {
				Computed:    true,
				Description: "Indicates if the volume is auxiliary or not.",
				Type:        schema.TypeBool,
			},
			Attr_AuxiliaryVolumeName: {
				Computed:    true,
				Description: "The auxiliary volume name.",
				Type:        schema.TypeString,
			},
			Attr_ConsistencyGroupName: {
				Computed:    true,
				Description: "The consistency group name if volume is a part of volume group.",
				Type:        schema.TypeString,
			},
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			Attr_DeleteOnTermination: {
				Computed:    true,
				Description: "Indicates if the volume should be deleted when the server terminates.",
				Type:        schema.TypeBool,
			},
			Attr_GroupID: {
				Computed:    true,
				Description: "The volume group id to which volume belongs.",
				Type:        schema.TypeString,
			},
			Attr_IOThrottleRate: {
				Computed:    true,
				Description: "Amount of iops assigned to the volume.",
				Type:        schema.TypeString,
			},
			Attr_MasterVolumeName: {
				Computed:    true,
				Description: "Indicates master volume name",
				Type:        schema.TypeString,
			},
			Attr_MirroringState: {
				Computed:    true,
				Description: "Mirroring state for replication enabled volume",
				Type:        schema.TypeString,
			},
			Attr_PrimaryRole: {
				Computed:    true,
				Description: "Indicates whether 'master'/'auxiliary' volume is playing the primary role.",
				Type:        schema.TypeString,
			},
			Attr_ReplicationStatus: {
				Computed:    true,
				Description: "The replication status of the volume.",
				Type:        schema.TypeString,
			},
			Attr_ReplicationSites: {
				Computed:    true,
				Description: "List of replication sites for volume replication.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
			Attr_ReplicationType: {
				Computed:    true,
				Description: "The replication type of the volume 'metro' or 'global'.",
				Type:        schema.TypeString,
			},
			Attr_VolumeID: {
				Computed:    true,
				Description: "The unique identifier of the volume.",
				Type:        schema.TypeString,
			},
			Attr_VolumeStatus: {
				Computed:    true,
				Description: "The status of the volume.",
				Type:        schema.TypeString,
			},
			Attr_WWN: {
				Computed:    true,
				Description: "The world wide name of the volume.",
				Type:        schema.TypeString,
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

	name := d.Get(Arg_VolumeName).(string)
	size := float64(d.Get(Arg_VolumeSize).(float64))
	var shared bool
	if v, ok := d.GetOk(Arg_VolumeShareable); ok {
		shared = v.(bool)
	}
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	body := &models.CreateDataVolume{
		Name:      &name,
		Shareable: &shared,
		Size:      &size,
	}
	if v, ok := d.GetOk(Arg_VolumeType); ok {
		volType := v.(string)
		body.DiskType = volType
	}
	if v, ok := d.GetOk(Arg_VolumePool); ok {
		volumePool := v.(string)
		body.VolumePool = volumePool
	}
	if v, ok := d.GetOk(Arg_ReplicationEnabled); ok {
		replicationEnabled := v.(bool)
		body.ReplicationEnabled = &replicationEnabled
	}
	if v, ok := d.GetOk(Arg_ReplicationSites); ok {
		if d.Get(Arg_ReplicationEnabled).(bool) {
			body.ReplicationSites = flex.FlattenSet(v.(*schema.Set))
		} else {
			return diag.Errorf("Replication (%s) must be enabled if replication sites are specified.", Arg_ReplicationEnabled)
		}
	}
	if ap, ok := d.GetOk(Arg_AffinityPolicy); ok {
		policy := ap.(string)
		body.AffinityPolicy = &policy

		if policy == "affinity" {
			if av, ok := d.GetOk(Arg_AffinityVolume); ok {
				afvol := av.(string)
				body.AffinityVolume = &afvol
			}
			if ai, ok := d.GetOk(Arg_AffinityInstance); ok {
				afins := ai.(string)
				body.AffinityPVMInstance = &afins
			}
		} else {
			if avs, ok := d.GetOk(Arg_AntiAffinityVolumes); ok {
				afvols := flex.ExpandStringList(avs.([]interface{}))
				body.AntiAffinityVolumes = afvols
			}
			if ais, ok := d.GetOk(Arg_AntiAffinityInstances); ok {
				afinss := flex.ExpandStringList(ais.([]interface{}))
				body.AntiAffinityPVMInstances = afinss
			}
		}

	}
	if v, ok := d.GetOk(Arg_UserTags); ok {
		body.UserTags = flex.FlattenSet(v.(*schema.Set))
	}

	client := instance.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
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

	if _, ok := d.GetOk(Arg_UserTags); ok {
		oldList, newList := d.GetChange(Arg_UserTags)
		err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, string(vol.Crn), "", UserTagType)
		if err != nil {
			log.Printf("Error on update of volume (%s) pi_user_tags during creation: %s", volumeid, err)
		}
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

	client := instance.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)

	vol, err := client.Get(volumeID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set(Arg_CloudInstanceID, cloudInstanceID)
	if vol.VolumeID != nil {
		d.Set(Attr_VolumeID, vol.VolumeID)
	}
	if vol.Crn != "" {
		d.Set(Attr_CRN, vol.Crn)
		tags, err := flex.GetGlobalTagsUsingCRN(meta, string(vol.Crn), "", UserTagType)
		if err != nil {
			log.Printf("Error on get of volume (%s) pi_user_tags: %s", *vol.VolumeID, err)
		}
		d.Set(Arg_UserTags, tags)
	}
	d.Set(Arg_VolumeName, vol.Name)
	d.Set(Arg_VolumePool, vol.VolumePool)
	if vol.Shareable != nil {
		d.Set(Arg_VolumeShareable, vol.Shareable)
	}
	d.Set(Arg_VolumeSize, vol.Size)
	d.Set(Arg_VolumeType, vol.DiskType)

	d.Set(Attr_Auxiliary, vol.Auxiliary)
	d.Set(Attr_AuxiliaryVolumeName, vol.AuxVolumeName)
	d.Set(Attr_ConsistencyGroupName, vol.ConsistencyGroupName)
	if vol.DeleteOnTermination != nil {
		d.Set(Attr_DeleteOnTermination, vol.DeleteOnTermination)
	}
	d.Set(Attr_GroupID, vol.GroupID)
	d.Set(Attr_IOThrottleRate, vol.IoThrottleRate)
	d.Set(Attr_MasterVolumeName, vol.MasterVolumeName)
	d.Set(Attr_MirroringState, vol.MirroringState)
	d.Set(Attr_PrimaryRole, vol.PrimaryRole)
	d.Set(Arg_ReplicationEnabled, vol.ReplicationEnabled)
	d.Set(Attr_ReplicationSites, vol.ReplicationSites)
	d.Set(Attr_ReplicationStatus, vol.ReplicationStatus)
	d.Set(Attr_ReplicationType, vol.ReplicationType)
	d.Set(Attr_VolumeStatus, vol.State)
	d.Set(Attr_WWN, vol.Wwn)

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

	client := instance.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
	name := d.Get(Arg_VolumeName).(string)
	size := float64(d.Get(Arg_VolumeSize).(float64))
	var shareable bool
	if v, ok := d.GetOk(Arg_VolumeShareable); ok {
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

	if d.HasChanges(Arg_ReplicationEnabled, Arg_VolumeType) {
		volActionBody := models.VolumeAction{}
		if d.HasChange(Arg_ReplicationEnabled) {
			volActionBody.ReplicationEnabled = flex.PtrToBool(d.Get(Arg_ReplicationEnabled).(bool))
		}
		if d.HasChange(Arg_VolumeType) {
			volActionBody.TargetStorageTier = flex.PtrToString(d.Get(Arg_VolumeType).(string))
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

	if d.HasChange(Arg_UserTags) {
		crn := d.Get(Attr_CRN)
		if crn != nil && crn != "" {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, crn.(string), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi volume (%s) pi_user_tags: %s", volumeID, err)
			}
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

	client := instance.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
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

func isWaitForIBMPIVolumeAvailable(ctx context.Context, client *instance.IBMPIVolumeClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume (%s) to be available.", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Retry, State_Creating},
		Target:     []string{State_Available},
		Refresh:    isIBMPIVolumeRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIVolumeRefreshFunc(client *instance.IBMPIVolumeClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vol, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if vol.State == State_Available || vol.State == State_InUse {
			return vol, State_Available, nil
		}

		return vol, State_Creating, nil
	}
}

func isWaitForIBMPIVolumeDeleted(ctx context.Context, client *instance.IBMPIVolumeClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Deleting, State_Creating},
		Target:     []string{State_Deleted},
		Refresh:    isIBMPIVolumeDeleteRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}
	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIVolumeDeleteRefreshFunc(client *instance.IBMPIVolumeClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vol, err := client.Get(id)
		if err != nil {
			uErr := errors.Unwrap(err)
			switch uErr.(type) {
			case *p_cloud_volumes.PcloudCloudinstancesVolumesGetNotFound:
				log.Printf("[DEBUG] volume does not exist %v", err)
				return vol, State_Deleted, nil
			}
			return nil, "", err
		}
		if vol == nil {
			return vol, State_Deleted, nil
		}
		return vol, State_Deleting, nil
	}
}

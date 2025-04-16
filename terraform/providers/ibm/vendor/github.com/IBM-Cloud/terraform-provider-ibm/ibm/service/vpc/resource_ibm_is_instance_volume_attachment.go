// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isInstanceId                             = "instance"
	isInstanceVolAttVol                      = "volume"
	isInstanceVolAttTags                     = "tags"
	isInstanceVolAttId                       = "volume_attachment_id"
	isInstanceVolAttIops                     = "volume_iops"
	isInstanceExistingVolume                 = "existing"
	isInstanceVolAttName                     = "name"
	isInstanceVolAttVolume                   = "volume"
	isInstanceVolumeDeleteOnInstanceDelete   = "delete_volume_on_instance_delete"
	isInstanceVolumeDeleteOnAttachmentDelete = "delete_volume_on_attachment_delete"
	isInstanceVolCapacity                    = "capacity"
	isInstanceVolIops                        = "iops"
	isInstanceVolEncryptionKey               = "encryption_key"
	isInstanceVolProfile                     = "profile"
)

func ResourceIBMISInstanceVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMisInstanceVolumeAttachmentCreate,
		Read:     resourceIBMisInstanceVolumeAttachmentRead,
		Update:   resourceIBMisInstanceVolumeAttachmentUpdate,
		Delete:   resourceIBMisInstanceVolumeAttachmentDelete,
		Exists:   resourceIBMisInstanceVolumeAttachmentExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceVolumeValidate(diff)
				}),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				}),
		),
		Schema: map[string]*schema.Schema{
			isInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance_volume_attachment", isInstanceId),
				Description:  "Instance id",
			},
			isInstanceVolAttId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this volume attachment",
			},

			isInstanceVolAttName: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance_volume_attachment", isInstanceVolAttName),
				Description:  "The user-defined name for this volume attachment.",
			},

			isInstanceVolumeDeleteOnInstanceDelete: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If set to true, when deleting the instance the volume will also be deleted.",
			},
			isInstanceVolumeDeleteOnAttachmentDelete: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to true, when deleting the attachment, the volume will also be deleted. Default value for this true.",
			},
			isInstanceVolAttVol: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{isInstanceVolIops, isInstanceVolumeAttVolumeReferenceName, isInstanceVolProfile, isInstanceVolCapacity, isInstanceVolumeSnapshot, isInstanceVolAttTags},
				ValidateFunc:  validate.InvokeValidator("ibm_is_instance_volume_attachment", isInstanceName),
				Description:   "Instance id",
			},

			isInstanceVolIops: {
				Type:          schema.TypeInt,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{isInstanceVolAttVol},
				Description:   "The maximum I/O operations per second (IOPS) for the volume.",
			},

			isInstanceVolumeAttVolumeReferenceName: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance_volume_attachment", isInstanceVolumeAttVolumeReferenceName),
				Description:  "The unique user-defined name for this volume",
			},

			isInstanceVolAttTags: {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{isInstanceVolAttVol},
				Elem:          &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_instance_volume_attachment", "tags")},
				Set:           flex.ResourceIBMVPCHash,
				Description:   "UserTags for the volume instance",
			},

			isInstanceVolProfile: {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{isInstanceVolAttVol},
				Computed:      true,
				ValidateFunc:  validate.InvokeValidator("ibm_is_instance_volume_attachment", isInstanceVolProfile),
				Description:   "The  globally unique name for the volume profile to use for this volume.",
			},

			isInstanceVolCapacity: {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  []string{isInstanceVolAttVol, isInstanceVolCapacity, isInstanceVolumeSnapshot},
				ConflictsWith: []string{isInstanceVolAttVol},
				ValidateFunc:  validate.InvokeValidator("ibm_is_instance_volume_attachment", isInstanceVolCapacity),
				Description:   "The capacity of the volume in gigabytes. The specified minimum and maximum capacity values for creating or updating volumes may expand in the future.",
			},
			isInstanceVolEncryptionKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.",
			},
			isInstanceVolumeSnapshot: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				AtLeastOneOf:  []string{isInstanceVolAttVol, isInstanceVolCapacity, isInstanceVolumeSnapshot},
				ConflictsWith: []string{isInstanceVolAttVol},
				Description:   "The snapshot of the volume to be attached",
			},
			isInstanceVolumeAttVolumeReferenceCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this volume",
			},
			isInstanceVolumeAttVolumeReferenceDeleted: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Link to documentation about deleted resources",
			},
			isInstanceVolumeAttVolumeReferenceHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this volume",
			},

			isInstanceVolumeAttDevice: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique identifier for the device which is exposed to the instance operating system",
			},

			isInstanceVolumeAttHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this volume attachment",
			},

			isInstanceVolumeAttStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of this volume attachment, one of [ attached, attaching, deleting, detaching ]",
			},

			isInstanceVolumeAttType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of volume attachment one of [ boot, data ]",
			},

			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMISInstanceVolumeAttachmentValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceId,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceVolAttName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceVolCapacity,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "10",
			MaxValue:                   "16000"})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceVolumeAttVolumeReferenceName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceVolProfile,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "general-purpose, 5iops-tier, 10iops-tier, custom",
		})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmISInstanceVolumeAttachmentValidator := validate.ResourceValidator{ResourceName: "ibm_is_instance_volume_attachment", Schema: validateSchema}
	return &ibmISInstanceVolumeAttachmentValidator
}

func instanceVolAttachmentCreate(d *schema.ResourceData, meta interface{}, instanceId string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	instanceVolAttproto := &vpcv1.CreateInstanceVolumeAttachmentOptions{
		InstanceID: &instanceId,
	}
	volumeIdStr := ""
	if volumeId, ok := d.GetOk(isInstanceVolAttVol); ok {
		volumeIdStr = volumeId.(string)
	}
	if volumeIdStr != "" {
		var volProtoVol = &vpcv1.VolumeAttachmentPrototypeVolumeVolumeIdentity{}
		volProtoVol.ID = &volumeIdStr
		instanceVolAttproto.Volume = volProtoVol
	} else {
		var volProtoVol = &vpcv1.VolumeAttachmentPrototypeVolumeVolumePrototypeInstanceContext{}
		if volname, ok := d.GetOk(isInstanceVolumeAttVolumeReferenceName); ok {
			volnamestr := volname.(string)
			volProtoVol.Name = &volnamestr
		}
		var userTags *schema.Set
		if v, ok := d.GetOk(isInstanceVolAttTags); ok {
			userTags = v.(*schema.Set)
			if userTags != nil && userTags.Len() != 0 {
				userTagsArray := make([]string, userTags.Len())
				for i, userTag := range userTags.List() {
					userTagStr := userTag.(string)
					userTagsArray[i] = userTagStr
				}
				schematicTags := os.Getenv("IC_ENV_TAGS")
				var envTags []string
				if schematicTags != "" {
					envTags = strings.Split(schematicTags, ",")
					userTagsArray = append(userTagsArray, envTags...)
				}
				volProtoVol.UserTags = userTagsArray
			}
		}
		volSnapshotStr := ""
		if volSnapshot, ok := d.GetOk(isInstanceVolumeSnapshot); ok {
			volSnapshotStr = volSnapshot.(string)
			volProtoVol.SourceSnapshot = &vpcv1.SnapshotIdentity{
				ID: &volSnapshotStr,
			}
		}
		encryptionCRNStr := ""
		if encryptionCRN, ok := d.GetOk(isInstanceVolEncryptionKey); ok {
			encryptionCRNStr = encryptionCRN.(string)
			volProtoVol.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
				CRN: &encryptionCRNStr,
			}
		}
		var snapCapacity int64
		if volSnapshotStr != "" {
			snapshotGet, _, err := sess.GetSnapshot(&vpcv1.GetSnapshotOptions{
				ID: &volSnapshotStr,
			})
			if err != nil {
				return fmt.Errorf("[ERROR] Error while getting snapshot details %q for instance %s: %q", volSnapshotStr, d.Id(), err)
			}
			snapCapacity = int64(int(*snapshotGet.MinimumCapacity))
		}
		var volCapacityInt int64
		if volCapacity, ok := d.GetOk(isInstanceVolCapacity); ok {
			volCapacityInt = int64(volCapacity.(int))
			if volCapacityInt != 0 && volCapacityInt > snapCapacity {
				volProtoVol.Capacity = &volCapacityInt
			}
		}
		var iops int64
		if volIops, ok := d.GetOk(isInstanceVolIops); ok {
			iops = int64(volIops.(int))
			if iops != 0 {
				volProtoVol.Iops = &iops
			}
			volProfileStr := "custom"
			volProtoVol.Profile = &vpcv1.VolumeProfileIdentity{
				Name: &volProfileStr,
			}
		} else {
			volProfileStr := "general-purpose"
			if volProfile, ok := d.GetOk(isInstanceVolProfile); ok {
				volProfileStr = volProfile.(string)
				volProtoVol.Profile = &vpcv1.VolumeProfileIdentity{
					Name: &volProfileStr,
				}
			} else {
				volProtoVol.Profile = &vpcv1.VolumeProfileIdentity{
					Name: &volProfileStr,
				}
			}
		}

		instanceVolAttproto.Volume = volProtoVol
	}

	if autoDelete, ok := d.GetOk(isInstanceVolumeDeleteOnInstanceDelete); ok {
		autoDeleteBool := autoDelete.(bool)
		instanceVolAttproto.DeleteVolumeOnInstanceDelete = &autoDeleteBool
	}
	if name, ok := d.GetOk(isInstanceVolAttName); ok {
		namestr := name.(string)
		instanceVolAttproto.Name = &namestr
	}

	isInstanceKey := "instance_key_" + instanceId
	conns.IbmMutexKV.Lock(isInstanceKey)
	defer conns.IbmMutexKV.Unlock(isInstanceKey)

	instanceVolAtt, response, err := sess.CreateInstanceVolumeAttachment(instanceVolAttproto)
	if err != nil {
		log.Printf("[DEBUG] Instance volume attachment create err %s\n%s", err, response)
		return fmt.Errorf("[ERROR] Error while attaching volume for instance %s: %q", instanceId, err)
	}
	d.SetId(makeTerraformVolAttID(instanceId, *instanceVolAtt.ID))
	_, err = isWaitForInstanceVolumeAttached(sess, d, instanceId, *instanceVolAtt.ID)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Instance (%s) volume attachment : %s", instanceId, *instanceVolAtt.ID)
	return nil
}

func resourceIBMisInstanceVolumeAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	instanceId := d.Get(isInstanceId).(string)
	err := instanceVolAttachmentCreate(d, meta, instanceId)
	if err != nil {
		return err
	}
	return resourceIBMisInstanceVolumeAttachmentRead(d, meta)
}

func resourceIBMisInstanceVolumeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	instanceID, id, err := parseVolAttTerraformID(d.Id())
	if err != nil {
		return err
	}
	err = instanceVolumeAttachmentGet(d, meta, instanceID, id)
	if err != nil {
		return err
	}
	return nil
}

func instanceVolumeAttachmentGet(d *schema.ResourceData, meta interface{}, instanceId, id string) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getinsVolAttOptions := &vpcv1.GetInstanceVolumeAttachmentOptions{
		InstanceID: &instanceId,
		ID:         &id,
	}
	volumeAtt, response, err := instanceC.GetInstanceVolumeAttachment(getinsVolAttOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Instance volume attachment : %s\n%s", err, response)
	}
	d.Set(isInstanceId, instanceId)

	if volumeAtt.Volume != nil {
		d.Set(isInstanceVolumeAttVolumeReferenceName, *volumeAtt.Volume.Name)
		d.Set(isInstanceVolumeAttVolumeReferenceCrn, *volumeAtt.Volume.CRN)
		if volumeAtt.Volume.Deleted != nil {
			d.Set(isInstanceVolumeAttVolumeReferenceDeleted, *volumeAtt.Volume.Deleted.MoreInfo)
		}
		d.Set(isInstanceVolumeAttVolumeReferenceHref, *volumeAtt.Volume.Href)
	}
	d.Set(isInstanceVolumeDeleteOnInstanceDelete, *volumeAtt.DeleteVolumeOnInstanceDelete)
	d.Set(isInstanceVolAttName, *volumeAtt.Name)
	if volumeAtt.Device != nil {
		d.Set(isInstanceVolumeAttDevice, *volumeAtt.Device.ID)
	}
	d.Set(isInstanceVolumeAttHref, *volumeAtt.Href)
	d.Set(isInstanceVolAttId, *volumeAtt.ID)
	d.Set(isInstanceVolumeAttStatus, *volumeAtt.Status)
	d.Set(isInstanceVolumeAttType, *volumeAtt.Type)
	if err = d.Set("version", response.Headers.Get("Etag")); err != nil {
		return fmt.Errorf("Error setting version: %s", err)
	}
	volId := *volumeAtt.Volume.ID
	getVolOptions := &vpcv1.GetVolumeOptions{
		ID: &volId,
	}
	volumeDetail, _, err := instanceC.GetVolume(getVolOptions)
	if err != nil || volumeDetail == nil {
		return fmt.Errorf("[ERROR] Error while getting volume details of volume %s ", id)
	}

	d.Set(isInstanceVolAttVol, *volumeDetail.ID)
	d.Set(isInstanceVolIops, *volumeDetail.Iops)
	d.Set(isInstanceVolProfile, *volumeDetail.Profile.Name)
	d.Set(isInstanceVolCapacity, *volumeDetail.Capacity)
	if volumeDetail.EncryptionKey != nil {
		d.Set(isInstanceVolEncryptionKey, *volumeDetail.EncryptionKey.CRN)
	}
	if volumeDetail.SourceSnapshot != nil {
		d.Set(isInstanceVolumeSnapshot, *volumeDetail.SourceSnapshot.ID)
	}
	return nil
}

func instanceVolAttUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}
	instanceId, id, err := parseVolAttTerraformID(d.Id())
	if err != nil {
		return err
	}
	updateInstanceVolAttOptions := &vpcv1.UpdateInstanceVolumeAttachmentOptions{
		InstanceID: &instanceId,
		ID:         &id,
	}
	flag := false

	// name && auto delete change
	volAttNamePatchModel := &vpcv1.VolumeAttachmentPatch{}
	if d.HasChange(isInstanceVolumeDeleteOnInstanceDelete) {
		autoDelete := d.Get(isInstanceVolumeDeleteOnInstanceDelete).(bool)
		volAttNamePatchModel.DeleteVolumeOnInstanceDelete = &autoDelete
		flag = true
	}

	if d.HasChange(isInstanceVolAttName) {
		name := d.Get(isInstanceVolAttName).(string)
		volAttNamePatchModel.Name = &name
		flag = true
	}
	if flag {
		volAttNamePatchModelAsPatch, err := volAttNamePatchModel.AsPatch()
		if err != nil || volAttNamePatchModelAsPatch == nil {
			return fmt.Errorf("[ERROR] Error Instance volume attachment (%s) as patch : %s", id, err)
		}
		updateInstanceVolAttOptions.VolumeAttachmentPatch = volAttNamePatchModelAsPatch

		instanceVolAttUpdate, response, err := instanceC.UpdateInstanceVolumeAttachment(updateInstanceVolAttOptions)
		if err != nil || instanceVolAttUpdate == nil {
			log.Printf("[DEBUG] Instance volume attachment updation err %s\n%s", err, response)
			return err
		}
	}

	if d.HasChange(isInstanceVolumeAttVolumeReferenceName) {
		newname := d.Get(isInstanceVolumeAttVolumeReferenceName).(string)
		volid := d.Get(isInstanceVolAttVol).(string)
		voloptions := &vpcv1.UpdateVolumeOptions{
			ID: &volid,
		}
		volumePatchModel := &vpcv1.VolumePatch{
			Name: &newname,
		}
		volumePatch, err := volumePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for VolumePatch: %s", err)
		}
		voloptions.VolumePatch = volumePatch
		_, response, err := instanceC.UpdateVolume(voloptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating volume name : %s\n%s", err, response)
		}
	}

	// profile/iops update

	volId := ""
	if volIdOk, ok := d.GetOk(isInstanceVolAttVol); ok {
		volId = volIdOk.(string)
	}

	if volId != "" && (d.HasChange(isInstanceVolIops) || d.HasChange(isInstanceVolProfile) || d.HasChange(isInstanceVolAttTags)) {
		insId := d.Get(isInstanceId).(string)
		getinsOptions := &vpcv1.GetInstanceOptions{
			ID: &insId,
		}
		instance, response, err := instanceC.GetInstance(getinsOptions)
		if err != nil || instance == nil {
			return fmt.Errorf("[ERROR] Error retrieving Instance (%s) : %s\n%s", insId, err, response)
		}

		if instance != nil && *instance.Status != "running" {
			actiontype := "start"
			createinsactoptions := &vpcv1.CreateInstanceActionOptions{
				InstanceID: &insId,
				Type:       &actiontype,
			}
			_, response, err = instanceC.CreateInstanceAction(createinsactoptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error starting Instance (%s) : %s\n%s", insId, err, response)
			}
			_, err = isWaitForInstanceAvailable(instanceC, insId, d.Timeout(schema.TimeoutCreate), d)
			if err != nil {
				return err
			}
		}
		updateVolumeProfileOptions := &vpcv1.UpdateVolumeOptions{
			ID: &volId,
		}
		volumeProfilePatchModel := &vpcv1.VolumePatch{}
		if d.HasChange(isInstanceVolProfile) {
			profile := d.Get(isInstanceVolProfile).(string)
			volumeProfilePatchModel.Profile = &vpcv1.VolumeProfileIdentity{
				Name: &profile,
			}
		} else if d.HasChange(isVolumeIops) {
			profile := d.Get(isInstanceVolProfile).(string)
			volumeProfilePatchModel.Profile = &vpcv1.VolumeProfileIdentity{
				Name: &profile,
			}
			iops := int64(d.Get(isVolumeIops).(int))
			volumeProfilePatchModel.Iops = &iops
		}
		if d.HasChange(isInstanceVolAttTags) && !d.IsNewResource() {
			if v, ok := d.GetOk(isInstanceVolAttTags); ok {
				userTags := v.(*schema.Set)
				if userTags != nil && userTags.Len() != 0 {
					userTagsArray := make([]string, userTags.Len())
					for i, userTag := range userTags.List() {
						userTagStr := userTag.(string)
						userTagsArray[i] = userTagStr
					}
					schematicTags := os.Getenv("IC_ENV_TAGS")
					var envTags []string
					if schematicTags != "" {
						envTags = strings.Split(schematicTags, ",")
						userTagsArray = append(userTagsArray, envTags...)
					}
					volumeProfilePatchModel.UserTags = userTagsArray
				}
			}

		}

		volumeProfilePatch, err := volumeProfilePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for volumeProfilePatch: %s", err)
		}
		optionsget := &vpcv1.GetVolumeOptions{
			ID: &volId,
		}
		_, response, err = instanceC.GetVolume(optionsget)
		if err != nil {
			return fmt.Errorf("[ERROR] Error getting Boot Volume (%s): %s\n%s", id, err, response)
		}
		eTag := response.Headers.Get("ETag")
		updateVolumeProfileOptions.IfMatch = &eTag
		updateVolumeProfileOptions.VolumePatch = volumeProfilePatch
		_, response, err = instanceC.UpdateVolume(updateVolumeProfileOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating volume profile/iops/userTags: %s\n%s", err, response)
		}
		isWaitForVolumeAvailable(instanceC, volId, d.Timeout(schema.TimeoutCreate))
	}

	// capacity update

	if volId != "" && d.HasChange(isInstanceVolCapacity) {

		getvolumeoptions := &vpcv1.GetVolumeOptions{
			ID: &volId,
		}
		vol, response, err := instanceC.GetVolume(getvolumeoptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Getting Volume (%s): %s\n%s", id, err, response)
		}

		if vol.VolumeAttachments == nil || len(vol.VolumeAttachments) == 0 || *vol.VolumeAttachments[0].Name == "" {
			return fmt.Errorf("[ERROR] Error volume capacity can't be updated since volume %s is not attached to any instance for VolumePatch", id)
		}

		getinsOptions := &vpcv1.GetInstanceOptions{
			ID: &instanceId,
		}
		instance, response, err := instanceC.GetInstance(getinsOptions)
		if err != nil || instance == nil {
			return fmt.Errorf("[ERROR] Error retrieving Instance (%s) : %s\n%s", instanceId, err, response)
		}
		if instance != nil && *instance.Status != "running" {
			actiontype := "start"
			createinsactoptions := &vpcv1.CreateInstanceActionOptions{
				InstanceID: &instanceId,
				Type:       &actiontype,
			}
			_, response, err = instanceC.CreateInstanceAction(createinsactoptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error starting Instance (%s) : %s\n%s", instanceId, err, response)
			}
			_, err = isWaitForInstanceAvailable(instanceC, instanceId, d.Timeout(schema.TimeoutCreate), d)
			return fmt.Errorf("[ERROR] Error starting Instance (%s) : %s\n%s", instanceId, err, response)
		}
		capacity := int64(d.Get(isVolumeCapacity).(int))
		updateVolumeOptions := &vpcv1.UpdateVolumeOptions{
			ID: &volId,
		}
		volumeCapacityPatchModel := &vpcv1.VolumePatch{}
		volumeCapacityPatchModel.Capacity = &capacity
		volumeCapacityPatch, err := volumeCapacityPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for volumeCapacityPatchModel: %s", err)
		}
		updateVolumeOptions.VolumePatch = volumeCapacityPatch
		_, response, err = instanceC.UpdateVolume(updateVolumeOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating volume capacity: %s\n%s", err, response)
		}
		_, err = isWaitForVolumeAvailable(instanceC, volId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceIBMisInstanceVolumeAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {

	err := instanceVolAttUpdate(d, meta)
	if err != nil {
		return err
	}
	return resourceIBMisInstanceVolumeAttachmentRead(d, meta)
}

func instanceVolAttDelete(d *schema.ResourceData, meta interface{}, instanceId, id, volId string, volDelete bool) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}

	deleteInstanceVolAttOptions := &vpcv1.DeleteInstanceVolumeAttachmentOptions{
		InstanceID: &instanceId,
		ID:         &id,
	}

	isInstanceKey := "instance_key_" + instanceId
	conns.IbmMutexKV.Lock(isInstanceKey)
	defer conns.IbmMutexKV.Unlock(isInstanceKey)

	_, err = instanceC.DeleteInstanceVolumeAttachment(deleteInstanceVolAttOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error while deleting volume attachment (%s) from instance (%s) : %q", id, instanceId, err)
	}
	_, err = isWaitForInstanceVolumeDetached(instanceC, d, instanceId, id)

	if err != nil {
		return fmt.Errorf("[ERROR] Error while deleting volume attachment (%s) from instance (%s) on wait : %q", id, instanceId, err)
	}
	if volDelete {
		deleteVolumeOptions := &vpcv1.DeleteVolumeOptions{
			ID: &volId,
		}
		response, err := instanceC.DeleteVolume(deleteVolumeOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error while deleting volume : %s\n%s", err, response)
		}
		_, err = isWaitForVolumeDeleted(instanceC, volId, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceIBMisInstanceVolumeAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	instanceId, id, err := parseVolAttTerraformID(d.Id())
	if err != nil {
		return err
	}

	volDelete := false
	if volDeleteOk, ok := d.GetOk(isInstanceVolumeDeleteOnAttachmentDelete); ok {
		volDelete = volDeleteOk.(bool)
	}
	volId := ""
	if volIdOk, ok := d.GetOk(isInstanceVolAttVol); ok {
		volId = volIdOk.(string)
	}

	err = instanceVolAttDelete(d, meta, instanceId, id, volId, volDelete)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceIBMisInstanceVolumeAttachmentExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	instanceId, id, err := parseVolAttTerraformID(d.Id())
	if err != nil {
		return false, err
	}
	exists, err := instanceVolAttExists(d, meta, instanceId, id)
	return exists, err
}

func instanceVolAttExists(d *schema.ResourceData, meta interface{}, instanceId, id string) (bool, error) {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getinsvolattOptions := &vpcv1.GetInstanceVolumeAttachmentOptions{
		InstanceID: &instanceId,
		ID:         &id,
	}
	_, response, err := instanceC.GetInstanceVolumeAttachment(getinsvolattOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting Instance volume attachment: %s\n%s", err, response)
	}
	return true, nil
}

func makeTerraformVolAttID(id1, id2 string) string {
	// Include both instance id and volume attachment to create a unique Terraform id.  As a bonus,
	// we can extract the instance id as needed for API calls such as READ.
	return fmt.Sprintf("%s/%s", id1, id2)
}

func parseVolAttTerraformID(s string) (string, string, error) {
	segments := strings.Split(s, "/")
	if len(segments) != 2 {
		return "", "", fmt.Errorf("invalid terraform Id %s (incorrect number of segments)", s)
	}
	if segments[0] == "" || segments[1] == "" {
		return "", "", fmt.Errorf("invalid terraform Id %s (one or more empty segments)", s)
	}
	return segments[0], segments[1], nil
}

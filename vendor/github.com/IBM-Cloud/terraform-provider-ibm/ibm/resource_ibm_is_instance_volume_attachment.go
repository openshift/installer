// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isInstanceId                             = "instance"
	isInstanceVolAttVol                      = "volume"
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

func resourceIBMISInstanceVolumeAttachment() *schema.Resource {
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

		Schema: map[string]*schema.Schema{
			isInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: InvokeValidator("ibm_is_instance_volume_attachment", isInstanceId),
				Description:  "Instance id",
			},
			isInstanceVolAttId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this volume attachment",
			},

			isInstanceVolAttName: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The user-defined name for this volume attachment.",
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
				ConflictsWith: []string{isInstanceVolIops, isInstanceVolumeAttVolumeReferenceName, isInstanceVolProfile, isInstanceVolCapacity, isInstanceVolumeSnapshot},
				ValidateFunc:  InvokeValidator("ibm_is_instance_volume_attachment", isInstanceName),
				Description:   "Instance id",
			},

			isInstanceVolIops: {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The maximum I/O operations per second (IOPS) for the volume.",
			},

			isInstanceVolumeAttVolumeReferenceName: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique user-defined name for this volume",
			},

			isInstanceVolProfile: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The  globally unique name for the volume profile to use for this volume.",
			},
			isInstanceVolCapacity: {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				AtLeastOneOf:  []string{isInstanceVolAttVol, isInstanceVolCapacity, isInstanceVolumeSnapshot},
				ConflictsWith: []string{isInstanceVolAttVol},
				ValidateFunc:  InvokeValidator("ibm_is_instance_volume_attachment", isInstanceVolCapacity),
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
		},
	}
}

func resourceIBMISInstanceVolumeAttachmentValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isInstanceId,
			ValidateFunctionIdentifier: ValidateNoZeroValues,
			Type:                       TypeString})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isInstanceVolAttName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isInstanceVolCapacity,
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			MinValue:                   "10",
			MaxValue:                   "2000"})

	ibmISInstanceVolumeAttachmentValidator := ResourceValidator{ResourceName: "ibm_is_instance_volume_attachment", Schema: validateSchema}
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
		volSnapshotStr := ""
		if volSnapshot, ok := d.GetOk(isInstanceVolumeSnapshot); ok {
			volSnapshotStr = volSnapshot.(string)
			volProtoVol.SourceSnapshot = &vpcv1.SnapshotIdentity{
				ID: &volSnapshotStr,
			}
		}
		var snapCapacity int64
		if volSnapshotStr != "" {
			snapshotGet, _, err := sess.GetSnapshot(&vpcv1.GetSnapshotOptions{
				ID: &volSnapshotStr,
			})
			if err != nil {
				return fmt.Errorf("Error while getting snapshot details %q for instance %s: %q", volSnapshotStr, d.Id(), err)
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

	instanceVolAtt, response, err := sess.CreateInstanceVolumeAttachment(instanceVolAttproto)
	if err != nil {
		log.Printf("[DEBUG] Instance volume attachment create err %s\n%s", err, response)
		return fmt.Errorf("Error while attaching volume for instance %s: %q", instanceId, err)
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
		return fmt.Errorf("Error getting Instance volume attachment : %s\n%s", err, response)
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

	volId := *volumeAtt.Volume.ID
	getVolOptions := &vpcv1.GetVolumeOptions{
		ID: &volId,
	}
	volumeDetail, _, err := instanceC.GetVolume(getVolOptions)
	if err != nil || volumeDetail == nil {
		return fmt.Errorf("Error while getting volume details of volume %s ", id)
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
	volAttPatchModel := &vpcv1.VolumeAttachmentPatch{}
	if d.HasChange(isInstanceVolumeDeleteOnInstanceDelete) {
		autoDelete := d.Get(isInstanceVolumeDeleteOnInstanceDelete).(bool)
		volAttPatchModel.DeleteVolumeOnInstanceDelete = &autoDelete
		flag = true
	}

	if d.HasChange(isInstanceVolAttName) {
		name := d.Get(isInstanceVolAttName).(string)
		volAttPatchModel.Name = &name
		flag = true
	}
	if flag {
		volAttPatchModelAsPatch, err := volAttPatchModel.AsPatch()
		if err != nil || volAttPatchModelAsPatch == nil {
			return fmt.Errorf("Error Instance volume attachment (%s) as patch : %s", id, err)
		}
		updateInstanceVolAttOptions.VolumeAttachmentPatch = volAttPatchModelAsPatch

		instanceVolAttUpdate, response, err := instanceC.UpdateInstanceVolumeAttachment(updateInstanceVolAttOptions)
		if err != nil || instanceVolAttUpdate == nil {
			log.Printf("[DEBUG] Instance volume attachment creation err %s\n%s", err, response)
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
	_, err = instanceC.DeleteInstanceVolumeAttachment(deleteInstanceVolAttOptions)
	_, err = isWaitForInstanceVolumeDetached(instanceC, d, instanceId, id)
	if err != nil {
		return fmt.Errorf("Error while deleting volume attachment (%s) from instance (%s) : %q", id, instanceId, err)
	}
	if volDelete {
		deleteVolumeOptions := &vpcv1.DeleteVolumeOptions{
			ID: &volId,
		}
		response, err := instanceC.DeleteVolume(deleteVolumeOptions)
		if err != nil {
			return fmt.Errorf("Error while deleting volume : %s\n%s", err, response)
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
		return false, fmt.Errorf("Error getting Instance volume attachment: %s\n%s", err, response)
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

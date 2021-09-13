// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isImageHref                   = "href"
	isImageName                   = "name"
	isImageTags                   = "tags"
	isImageOperatingSystem        = "operating_system"
	isImageStatus                 = "status"
	isImageVisibility             = "visibility"
	isImageFile                   = "file"
	isImageVolume                 = "source_volume"
	isImageMinimumProvisionedSize = "size"

	isImageResourceGroup    = "resource_group"
	isImageEncryptedDataKey = "encrypted_data_key"
	isImageEncryptionKey    = "encryption_key"
	isImageEncryption       = "encryption"
	isImageCheckSum         = "checksum"

	isImageProvisioning     = "provisioning"
	isImageProvisioningDone = "done"
	isImageDeleting         = "deleting"
	isImageDeleted          = "done"
)

func resourceIBMISImage() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISImageCreate,
		Read:     resourceIBMISImageRead,
		Update:   resourceIBMISImageUpdate,
		Delete:   resourceIBMISImageDelete,
		Exists:   resourceIBMISImageExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			isImageHref: {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: applyOnce,
				RequiredWith:     []string{isImageOperatingSystem},
				ExactlyOneOf:     []string{isImageHref, isImageVolume},
				Description:      "Image Href value",
			},

			isImageName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: InvokeValidator("ibm_is_image", isImageName),
				Description:  "Image name",
			},

			isImageEncryptedDataKey: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A base64-encoded, encrypted representation of the key that was used to encrypt the data for this image",
			},
			isImageEncryptionKey: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The CRN of the Key Protect Root Key or Hyper Protect Crypto Service Root Key for this resource",
			},
			isImageTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_is_image", "tag")},
				Set:         resourceIBMVPCHash,
				Description: "Tags for the image",
			},

			isImageOperatingSystem: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{isImageHref},
				Computed:     true,
				Description:  "Image Operating system",
			},

			isImageEncryption: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of encryption used on the image",
			},
			isImageStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of this image",
			},

			isImageMinimumProvisionedSize: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum size (in gigabytes) of a volume onto which this image may be provisioned",
			},

			isImageVisibility: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the image is publicly visible or private to the account",
			},

			isImageFile: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Details for the stored image file",
			},

			isImageVolume: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{isImageHref, isImageVolume},
				Description:  "Image volume id",
			},

			isImageResourceGroup: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "The resource group for this image",
			},

			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			isImageCheckSum: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SHA256 checksum of this image",
			},

			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func resourceIBMISImageValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isImageName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "tag",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})
	ibmISImageResourceValidator := ResourceValidator{ResourceName: "ibm_is_image", Schema: validateSchema}
	return &ibmISImageResourceValidator
}

func resourceIBMISImageCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] Image create")
	href := d.Get(isImageHref).(string)
	name := d.Get(isImageName).(string)
	operatingSystem := d.Get(isImageOperatingSystem).(string)
	volume := d.Get(isImageVolume).(string)

	if volume != "" {
		err := imgCreateByVolume(d, meta, name, volume)
		if err != nil {
			return err
		}
	} else {
		err := imgCreateByFile(d, meta, href, name, operatingSystem)
		if err != nil {
			return err
		}
	}

	return resourceIBMISImageRead(d, meta)
}

func imgCreateByFile(d *schema.ResourceData, meta interface{}, href, name, operatingSystem string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	imagePrototype := &vpcv1.ImagePrototypeImageByFile{
		Name: &name,
		File: &vpcv1.ImageFilePrototype{
			Href: &href,
		},
		OperatingSystem: &vpcv1.OperatingSystemIdentity{
			Name: &operatingSystem,
		},
	}
	if encryptionKey, ok := d.GetOk(isImageEncryptionKey); ok {
		encryptionKeyStr := encryptionKey.(string)
		// Construct an instance of the EncryptionKeyReference model
		encryptionKeyReferenceModel := new(vpcv1.EncryptionKeyIdentity)
		encryptionKeyReferenceModel.CRN = &encryptionKeyStr
		imagePrototype.EncryptionKey = encryptionKeyReferenceModel
	}
	if encDataKey, ok := d.GetOk(isImageEncryptedDataKey); ok {
		encDataKeyStr := encDataKey.(string)
		imagePrototype.EncryptedDataKey = &encDataKeyStr
	}
	if rgrp, ok := d.GetOk(isImageResourceGroup); ok {
		rg := rgrp.(string)
		imagePrototype.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}
	options := &vpcv1.CreateImageOptions{
		ImagePrototype: imagePrototype,
	}
	image, response, err := sess.CreateImage(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] Image creation err %s\n%s", err, response)
	}
	d.SetId(*image.ID)
	log.Printf("[INFO] Image ID : %s", *image.ID)
	_, err = isWaitForImageAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isImageTags); ok || v != "" {
		oldList, newList := d.GetChange(isImageTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *image.CRN)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc Image (%s) tags: %s", d.Id(), err)
		}
	}
	return nil
}
func imgCreateByVolume(d *schema.ResourceData, meta interface{}, name, volume string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	imagePrototype := &vpcv1.ImagePrototypeImageBySourceVolume{
		Name: &name,
	}
	var insId string
	imagePrototype.SourceVolume = &vpcv1.VolumeIdentity{
		ID: &volume,
	}
	options := &vpcv1.GetVolumeOptions{
		ID: &volume,
	}
	vol, response, err := sess.GetVolume(options)
	if err != nil || vol == nil {
		return fmt.Errorf("Error retrieving Volume (%s) details: %s\n%s", volume, err, response)
	}
	if vol.VolumeAttachments == nil {
		return fmt.Errorf("Error creating Image because the specified source_volume %s is not attached to a virtual server instance ", volume)
	}
	volAtt := &vol.VolumeAttachments[0]
	insId = *volAtt.Instance.ID
	getinsOptions := &vpcv1.GetInstanceOptions{
		ID: &insId,
	}
	instance, response, err := sess.GetInstance(getinsOptions)
	if err != nil || instance == nil {
		return fmt.Errorf("Error retrieving Instance (%s) to which the source_volume (%s) is attached : %s\n%s", insId, volume, err, response)
	}
	if instance != nil && *instance.Status == "running" {
		actiontype := "stop"
		createinsactoptions := &vpcv1.CreateInstanceActionOptions{
			InstanceID: &insId,
			Type:       &actiontype,
		}
		_, response, err = sess.CreateInstanceAction(createinsactoptions)
		if err != nil {
			return fmt.Errorf("Error stopping Instance (%s) to which the source_volume (%s) is attached  : %s\n%s", insId, volume, err, response)
		}
		_, err = isWaitForInstanceActionStop(sess, d.Timeout(schema.TimeoutCreate), insId, d)
		if err != nil {
			return err
		}
	} else if *instance.Status != "stopped" {
		_, err = isWaitForInstanceActionStop(sess, d.Timeout(schema.TimeoutCreate), insId, d)
		if err != nil {
			return err
		}
	}

	if encryptionKey, ok := d.GetOk(isImageEncryptionKey); ok {
		encryptionKeyStr := encryptionKey.(string)
		// Construct an instance of the EncryptionKeyReference model
		encryptionKeyReferenceModel := new(vpcv1.EncryptionKeyIdentity)
		encryptionKeyReferenceModel.CRN = &encryptionKeyStr
		imagePrototype.EncryptionKey = encryptionKeyReferenceModel
	}
	if rgrp, ok := d.GetOk(isImageResourceGroup); ok {
		rg := rgrp.(string)
		imagePrototype.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}
	imagOptions := &vpcv1.CreateImageOptions{
		ImagePrototype: imagePrototype,
	}
	image, response, err := sess.CreateImage(imagOptions)
	if err != nil {
		return fmt.Errorf("[DEBUG] Image creation err %s\n%s", err, response)
	}
	d.SetId(*image.ID)
	log.Printf("[INFO] Image ID : %s", *image.ID)
	_, err = isWaitForImageAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isImageTags); ok || v != "" {
		oldList, newList := d.GetChange(isImageTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *image.CRN)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc Image (%s) tags: %s", d.Id(), err)
		}
	}
	return nil
}

func isWaitForImageAvailable(imageC *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for image (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isImageProvisioning},
		Target:     []string{isImageProvisioningDone, ""},
		Refresh:    isImageRefreshFunc(imageC, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
func isImageRefreshFunc(imageC *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getimgoptions := &vpcv1.GetImageOptions{
			ID: &id,
		}
		image, response, err := imageC.GetImage(getimgoptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting Image: %s\n%s", err, response)
		}

		if *image.Status == "available" || *image.Status == "failed" {
			return image, isImageProvisioningDone, nil
		}

		return image, isImageProvisioning, nil
	}
}

func resourceIBMISImageUpdate(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	name := ""
	hasChanged := false

	if d.HasChange(isImageName) {
		name = d.Get(isImageName).(string)
		hasChanged = true
	}
	err := imgUpdate(d, meta, id, name, hasChanged)
	if err != nil {
		return err
	}

	return resourceIBMISImageRead(d, meta)
}

func imgUpdate(d *schema.ResourceData, meta interface{}, id, name string, hasChanged bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if d.HasChange(isImageTags) {
		options := &vpcv1.GetImageOptions{
			ID: &id,
		}
		image, response, err := sess.GetImage(options)
		if err != nil {
			return fmt.Errorf("Error getting Image IP: %s\n%s", err, response)
		}
		oldList, newList := d.GetChange(isImageTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *image.CRN)
		if err != nil {
			log.Printf(
				"Error on update of resource vpc Image (%s) tags: %s", id, err)
		}
	}
	if hasChanged {
		options := &vpcv1.UpdateImageOptions{
			ID: &id,
		}
		imagePatchModel := &vpcv1.ImagePatch{
			Name: &name,
		}
		imagePatch, err := imagePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for ImagePatch: %s", err)
		}
		options.ImagePatch = imagePatch
		_, response, err := sess.UpdateImage(options)
		if err != nil {
			return fmt.Errorf("Error on update of resource vpc Image: %s\n%s", err, response)
		}
	}
	return nil
}

func resourceIBMISImageRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	err := imgGet(d, meta, id)
	if err != nil {
		return err
	}
	return nil
}

func imgGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.GetImageOptions{
		ID: &id,
	}
	image, response, err := sess.GetImage(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Image (%s): %s\n%s", id, err, response)
	}
	// d.Set(isImageArchitecure, image.Architecture)
	if image.MinimumProvisionedSize != nil {
		d.Set(isImageMinimumProvisionedSize, *image.MinimumProvisionedSize)
	}
	d.Set(isImageName, *image.Name)
	d.Set(isImageOperatingSystem, *image.OperatingSystem.Name)
	// d.Set(isImageFormat, image.Format)
	if image.Encryption != nil {
		d.Set("encryption", *image.Encryption)
	}
	if image.EncryptionKey != nil {
		d.Set("encryption_key", *image.EncryptionKey.CRN)
	}
	if image.File != nil && image.File.Size != nil {
		d.Set(isImageFile, *image.File.Size)
	}
	if image.SourceVolume != nil {
		d.Set(isImageVolume, *image.SourceVolume.ID)
	}

	d.Set(isImageHref, *image.Href)
	d.Set(isImageStatus, *image.Status)
	d.Set(isImageVisibility, *image.Visibility)
	if image.Encryption != nil {
		d.Set(isImageEncryption, *image.Encryption)
	}
	if image.EncryptionKey != nil {
		d.Set(isImageEncryptionKey, *image.EncryptionKey.CRN)
	}
	if image.File != nil && image.File.Checksums != nil {
		d.Set(isImageCheckSum, *image.File.Checksums.Sha256)
	}
	tags, err := GetTagsUsingCRN(meta, *image.CRN)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc Image (%s) tags: %s", d.Id(), err)
	}
	d.Set(isImageTags, tags)
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/vpc-ext/compute/image")
	d.Set(ResourceName, *image.Name)
	d.Set(ResourceStatus, *image.Status)
	d.Set(ResourceCRN, *image.CRN)
	if image.ResourceGroup != nil {
		d.Set(isImageResourceGroup, *image.ResourceGroup.ID)
	}
	return nil
}

func resourceIBMISImageDelete(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	err := imgDelete(d, meta, id)
	if err != nil {
		return err
	}
	return nil
}

func imgDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getImageOptions := &vpcv1.GetImageOptions{
		ID: &id,
	}
	_, response, err := sess.GetImage(getImageOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error Getting Image (%s): %s\n%s", id, err, response)
	}

	options := &vpcv1.DeleteImageOptions{
		ID: &id,
	}
	response, err = sess.DeleteImage(options)
	if err != nil {
		return fmt.Errorf("Error Deleting Image : %s\n%s", err, response)
	}
	_, err = isWaitForImageDeleted(sess, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForImageDeleted(imageC *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for image (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isImageDeleting},
		Target:     []string{"", isImageDeleted},
		Refresh:    isImageDeleteRefreshFunc(imageC, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isImageDeleteRefreshFunc(imageC *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] delete function here")
		getimgoptions := &vpcv1.GetImageOptions{
			ID: &id,
		}
		image, response, err := imageC.GetImage(getimgoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return image, isImageDeleted, nil
			}
			return image, "", fmt.Errorf("Error Getting Image: %s\n%s", err, response)
		}
		return image, isImageDeleting, err
	}
}
func resourceIBMISImageExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id := d.Id()
	exists, err := imgExists(d, meta, id)
	return exists, err
}

func imgExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	options := &vpcv1.GetImageOptions{
		ID: &id,
	}
	_, response, err := sess.GetImage(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Image: %s\n%s", err, response)
	}
	return true, nil
}

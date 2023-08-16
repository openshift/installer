// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isImageId             = "image"
	isImageStatuses       = "statuses"
	isImageCreatedAt      = "created_at"
	isImageDeprecationAt  = "deprecation_at"
	isImageObsolescenceAt = "obsolescence_at"
)

func ResourceIBMISImageDeprecate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISImageDeprecateCreate,
		ReadContext:   resourceIBMISImageDeprecateRead,
		DeleteContext: resourceIBMISImageDeprecateDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set(isImageId, d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				},
			),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceValidateAccessTags(diff, v)
				}),
		),

		Schema: map[string]*schema.Schema{
			isImageHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image Href value",
			},

			isImageName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image name",
			},
			isImageId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Image identifier",
			},

			isImageEncryptedDataKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A base64-encoded, encrypted representation of the key that was used to encrypt the data for this image",
			},

			isImageCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image was created",
			},
			isImageDeprecationAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deprecation date and time (UTC) for this image. If absent, no deprecation date and time has been set.",
			},
			isImageObsolescenceAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The obsolescence date and time (UTC) for this image. If absent, no obsolescence date and time has been set.",
			},

			isImageEncryptionKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN of the Key Protect Root Key or Hyper Protect Crypto Service Root Key for this resource",
			},
			isImageTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "Tags for the image",
			},

			isImageOperatingSystem: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image Operating system",
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image volume id",
			},

			isImageResourceGroup: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group for this image",
			},

			IsImageCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			isImageCheckSum: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SHA256 checksum of this image",
			},

			isImageAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},
		},
	}
}

func resourceIBMISImageDeprecateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log.Printf("[DEBUG] Image deprecate create")
	id := d.Get(isImageId).(string)

	err := imgDeprecateCreate(context, d, meta, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMISImageDeprecateRead(context, d, meta)
}

func imgDeprecateCreate(context context.Context, d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	imageDeprecatePrototype := &vpcv1.DeprecateImageOptions{
		ID: &id,
	}
	response, err := sess.DeprecateImageWithContext(context, imageDeprecatePrototype)
	if err != nil {
		return fmt.Errorf("[ERROR] Image deprecate err %s\n%s", err, response)
	}
	d.SetId(id)
	log.Printf("[INFO] Image ID : %s", id)
	_, err = isWaitForImageDeprecate(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return nil
}

func isWaitForImageDeprecate(imageC *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for image (%s) to be deprecate.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isImageProvisioning},
		Target:     []string{isImageProvisioningDone, ""},
		Refresh:    isImageDeprecateRefreshFunc(imageC, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isImageDeprecateRefreshFunc(imageC *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getimgoptions := &vpcv1.GetImageOptions{
			ID: &id,
		}
		image, response, err := imageC.GetImage(getimgoptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Image: %s\n%s", err, response)
		}

		if *image.Status == "deprecated" {
			return image, isImageProvisioningDone, nil
		}
		if *image.Status == "failed" {
			return image, "", fmt.Errorf("[ERROR] Error Image(%s) went to failed state while deprecating ", *image.ID)
		}

		return image, isImageProvisioning, nil
	}
}

func resourceIBMISImageDeprecateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	id := d.Id()
	err := imgDeprecateGet(d, meta, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func imgDeprecateGet(d *schema.ResourceData, meta interface{}, id string) error {
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
		return fmt.Errorf("[ERROR] Error Getting Image (%s): %s\n%s", id, err, response)
	}
	if image.MinimumProvisionedSize != nil {
		d.Set(isImageMinimumProvisionedSize, *image.MinimumProvisionedSize)
	}

	if image.CreatedAt != nil {
		d.Set(isImageCreatedAt, image.CreatedAt.String())
	}
	if image.DeprecationAt != nil {
		d.Set(isImageDeprecationAt, image.DeprecationAt.String())
	}
	if image.ObsolescenceAt != nil {
		d.Set(isImageObsolescenceAt, image.ObsolescenceAt.String())
	}
	d.Set(isImageName, *image.Name)
	d.Set(isImageOperatingSystem, *image.OperatingSystem.Name)
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
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *image.CRN, "", isImageUserTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc Image (%s) tags: %s", d.Id(), err)
	}
	d.Set(isImageTags, tags)
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *image.CRN, "", isImageAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc Image (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isImageAccessTags, accesstags)
	d.Set(IsImageCRN, *image.CRN)
	if image.ResourceGroup != nil {
		d.Set(isImageResourceGroup, *image.ResourceGroup.ID)
	}
	return nil
}

func resourceIBMISImageDeprecateDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

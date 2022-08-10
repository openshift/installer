// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMISImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISImageRead,

		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"identifier", "name"},
				Description:  "Image name",
			},

			"identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"identifier", "name"},
				Description:  "Image id",
			},

			"visibility": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private"}),
				Description:  "Whether the image is publicly visible or private to the account",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of this image",
			},

			"os": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image Operating system",
			},
			"architecture": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The operating system architecture",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this image",
			},
			isImageCheckSum: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SHA256 Checksum for this image",
			},
			isImageEncryptionKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN of the Key Protect Root Key or Hyper Protect Crypto Service Root Key for this resource",
			},
			isImageEncryption: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of encryption used on the image",
			},
			"source_volume": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source volume id of the image",
			},
		},
	}
}

func dataSourceIBMISImageRead(d *schema.ResourceData, meta interface{}) error {

	name := d.Get("name").(string)
	identifier := d.Get("identifier").(string)
	var visibility string
	if v, ok := d.GetOk("visibility"); ok {
		visibility = v.(string)
	}
	if name != "" {
		err := imageGetByName(d, meta, name, visibility)
		if err != nil {
			return err
		}
	} else if identifier != "" {
		err := imageGetById(d, meta, identifier)
		if err != nil {
			return err
		}
	}

	return nil
}

func imageGetByName(d *schema.ResourceData, meta interface{}, name, visibility string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	listImagesOptions := &vpcv1.ListImagesOptions{
		Name: &name,
	}

	if visibility != "" {
		listImagesOptions.Visibility = &visibility
	}
	availableImages, response, err := sess.ListImages(listImagesOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Fetching Images %s\n%s", err, response)
	}
	allrecs := availableImages.Images

	if len(allrecs) == 0 {
		return fmt.Errorf("[ERROR] No image found with name  %s", name)
	}
	image := allrecs[0]
	d.SetId(*image.ID)
	d.Set("status", *image.Status)
	if *image.Status == "deprecated" {
		fmt.Printf("[WARN] Given image %s is deprecated and soon will be obsolete.", name)
	}
	d.Set("name", *image.Name)
	d.Set("visibility", *image.Visibility)
	d.Set("os", *image.OperatingSystem.Name)
	d.Set("architecture", *image.OperatingSystem.Architecture)
	d.Set("crn", *image.CRN)
	if image.Encryption != nil {
		d.Set("encryption", *image.Encryption)
	}
	if image.EncryptionKey != nil {
		d.Set("encryption_key", *image.EncryptionKey.CRN)
	}
	if image.File != nil && image.File.Checksums != nil {
		d.Set(isImageCheckSum, *image.File.Checksums.Sha256)
	}
	if image.SourceVolume != nil {
		d.Set("source_volume", *image.SourceVolume.ID)
	}
	return nil

}
func imageGetById(d *schema.ResourceData, meta interface{}, identifier string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getImageOptions := &vpcv1.GetImageOptions{
		ID: &identifier,
	}

	image, response, err := sess.GetImage(getImageOptions)
	if err != nil {
		if response.StatusCode == 404 {
			return fmt.Errorf("[ERROR] No image found with id  %s", identifier)
		}
		return fmt.Errorf("[ERROR] Error Fetching Images %s\n%s", err, response)
	}

	d.SetId(*image.ID)
	d.Set("status", *image.Status)
	if *image.Status == "deprecated" {
		fmt.Printf("[WARN] Given image %s is deprecated and soon will be obsolete.", name)
	}
	d.Set("name", *image.Name)
	d.Set("visibility", *image.Visibility)
	d.Set("os", *image.OperatingSystem.Name)
	d.Set("architecture", *image.OperatingSystem.Architecture)
	d.Set("crn", *image.CRN)
	if image.Encryption != nil {
		d.Set("encryption", *image.Encryption)
	}
	if image.EncryptionKey != nil {
		d.Set("encryption_key", *image.EncryptionKey.CRN)
	}
	if image.File != nil && image.File.Checksums != nil {
		d.Set(isImageCheckSum, *image.File.Checksums.Sha256)
	}
	if image.SourceVolume != nil {
		d.Set("source_volume", *image.SourceVolume.ID)
	}
	return nil
}

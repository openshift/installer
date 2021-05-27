// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMISImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISImageRead,

		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Image name",
			},

			"visibility": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"public", "private"}),
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
		},
	}
}

func dataSourceIBMISImageRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	var visibility string
	if v, ok := d.GetOk("visibility"); ok {
		visibility = v.(string)
	}
	if userDetails.generation == 1 {
		err := classicImageGet(d, meta, name, visibility)
		if err != nil {
			return err
		}
	} else {
		err := imageGet(d, meta, name, visibility)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicImageGet(d *schema.ResourceData, meta interface{}, name, visibility string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcclassicv1.Image{}
	for {
		listImagesOptions := &vpcclassicv1.ListImagesOptions{}
		if start != "" {
			listImagesOptions.Start = &start
		}
		if visibility != "" {
			listImagesOptions.Visibility = &visibility
		}
		availableImages, response, err := sess.ListImages(listImagesOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching Images %s\n%s", err, response)
		}
		start = GetNext(availableImages.Next)
		allrecs = append(allrecs, availableImages.Images...)
		if start == "" {
			break
		}
	}
	for _, image := range allrecs {
		if *image.Name == name {
			d.SetId(*image.ID)
			d.Set("status", *image.Status)
			d.Set("name", *image.Name)
			d.Set("visibility", *image.Visibility)
			d.Set("os", *image.OperatingSystem.Name)
			d.Set("architecture", *image.OperatingSystem.Architecture)
			d.Set("crn", *image.CRN)
			return nil
		}
	}
	return fmt.Errorf("No Image found with name %s", name)
}

func imageGet(d *schema.ResourceData, meta interface{}, name, visibility string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcv1.Image{}
	for {
		listImagesOptions := &vpcv1.ListImagesOptions{}
		if start != "" {
			listImagesOptions.Start = &start
		}
		if visibility != "" {
			listImagesOptions.Visibility = &visibility
		}
		availableImages, response, err := sess.ListImages(listImagesOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching Images %s\n%s", err, response)
		}
		start = GetNext(availableImages.Next)
		allrecs = append(allrecs, availableImages.Images...)
		if start == "" {
			break
		}
	}

	for _, image := range allrecs {
		if *image.Name == name {
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
			return nil
		}
	}

	return fmt.Errorf("No image found with name  %s", name)
}

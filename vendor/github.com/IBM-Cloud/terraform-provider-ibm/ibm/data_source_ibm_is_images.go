// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isImages                = "images"
	isImagesResourceGroupID = "resource_group"
)

func dataSourceIBMISImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISImagesRead,

		Schema: map[string]*schema.Schema{
			isImagesResourceGroupID: {
				Type:        schema.TypeString,
				Description: "The id of the resource group",
				Optional:    true,
			},
			isImageName: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_is_image", isImageName),
				Description:  "The name of the image",
			},
			isImageVisibility: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether the image is publicly visible or private to the account",
			},

			isImages: {
				Type:        schema.TypeList,
				Description: "List of images",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image name",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this image",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of this image",
						},
						"visibility": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the image is publicly visible or private to the account",
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
				},
			},
		},
	}
}

func dataSourceIBMISImagesRead(d *schema.ResourceData, meta interface{}) error {

	err := imageList(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func imageList(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcv1.Image{}

	var resourceGroupID string
	if v, ok := d.GetOk(isImagesResourceGroupID); ok {
		resourceGroupID = v.(string)
	}

	var imageName string
	if v, ok := d.GetOk(isImageName); ok {
		imageName = v.(string)
	}

	var visibility string
	if v, ok := d.GetOk(isImageVisibility); ok {
		visibility = v.(string)
	}

	listImagesOptions := &vpcv1.ListImagesOptions{}
	if resourceGroupID != "" {
		listImagesOptions.SetResourceGroupID(resourceGroupID)
	}
	if imageName != "" {
		listImagesOptions.SetName(imageName)
	}
	if visibility != "" {
		listImagesOptions.SetVisibility(visibility)
	}

	for {
		if start != "" {
			listImagesOptions.Start = &start
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
	imagesInfo := make([]map[string]interface{}, 0)
	for _, image := range allrecs {

		l := map[string]interface{}{
			"name":         *image.Name,
			"id":           *image.ID,
			"status":       *image.Status,
			"crn":          *image.CRN,
			"visibility":   *image.Visibility,
			"os":           *image.OperatingSystem.Name,
			"architecture": *image.OperatingSystem.Architecture,
		}
		if image.File != nil && image.File.Checksums != nil {
			l[isImageCheckSum] = *image.File.Checksums.Sha256
		}
		if image.Encryption != nil {
			l["encryption"] = *image.Encryption
		}
		if image.EncryptionKey != nil {
			l["encryption_key"] = *image.EncryptionKey.CRN
		}
		if image.SourceVolume != nil {
			l["source_volume"] = *image.SourceVolume.ID
		}
		imagesInfo = append(imagesInfo, l)
	}
	d.SetId(dataSourceIBMISSubnetsID(d))
	d.Set(isImages, imagesInfo)
	return nil
}

// dataSourceIBMISImagesId returns a reasonable ID for a image list.
func dataSourceIBMISImagesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

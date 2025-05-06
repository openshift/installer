// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isImageCatalogOffering         = "catalog_offering"
	isImageCatalogOfferingManaged  = "managed"
	isImageCatalogOfferingVersion  = "version"
	isImageCatalogOfferingCrn      = "crn"
	isImageCatalogOfferingDeleted  = "deleted"
	isImageCatalogOfferingMoreInfo = "more_info"
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
			isImageCatalogOffering: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isImageCatalogOfferingManaged: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this image is managed as part of a catalog offering. A managed image can be provisioned using its catalog offering CRN or catalog offering version CRN.",
						},
						isImageCatalogOfferingVersion: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The catalog offering version associated with this image.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// isImageCatalogOfferingDeleted: {
									// 	Type:        schema.TypeList,
									// 	Computed:    true,
									// 	Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
									// 	Elem: &schema.Resource{
									// 		Schema: map[string]*schema.Schema{
									// 			isImageCatalogOfferingMoreInfo: {
									// 				Type:        schema.TypeString,
									// 				Computed:    true,
									// 				Description: "Link to documentation about deleted resources.",
									// 			},
									// 		},
									// 	},
									// },
									isImageCatalogOfferingCrn: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this version of the IBM Cloud catalog offering.",
									},
								},
							},
						},
					},
				},
			},
			isImageAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access tags",
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
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *image.CRN, "", isImageAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource image (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isImageAccessTags, accesstags)
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
	if image.CreatedAt != nil {
		d.Set(isImageCreatedAt, image.CreatedAt.String())
	}
	if image.DeprecationAt != nil {
		d.Set(isImageDeprecationAt, image.DeprecationAt.String())
	}
	if image.ObsolescenceAt != nil {
		d.Set(isImageObsolescenceAt, image.ObsolescenceAt.String())
	}
	if image.CatalogOffering != nil {
		catalogOfferingList := []map[string]interface{}{}
		catalogOfferingMap := dataSourceImageCollectionCatalogOfferingToMap(*image.CatalogOffering)
		catalogOfferingList = append(catalogOfferingList, catalogOfferingMap)
		d.Set(isImageCatalogOffering, catalogOfferingList)
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
	if image.CatalogOffering != nil {
		catalogOfferingList := []map[string]interface{}{}
		catalogOfferingMap := dataSourceImageCollectionCatalogOfferingToMap(*image.CatalogOffering)
		catalogOfferingList = append(catalogOfferingList, catalogOfferingMap)
		d.Set(isImageCatalogOffering, catalogOfferingList)
	}
	return nil
}

func dataSourceImageCollectionCatalogOfferingToMap(imageCatalogOfferingItem vpcv1.ImageCatalogOffering) (imageCatalogOfferingMap map[string]interface{}) {
	imageCatalogOfferingMap = map[string]interface{}{}
	if imageCatalogOfferingItem.Managed != nil {
		imageCatalogOfferingMap[isImageCatalogOfferingManaged] = imageCatalogOfferingItem.Managed
	}
	if imageCatalogOfferingItem.Version != nil {
		imageCatalogOfferingVersionList := []map[string]interface{}{}
		imageCatalogOfferingVersionMap := map[string]interface{}{}
		imageCatalogOfferingVersionMap[isImageCatalogOfferingCrn] = imageCatalogOfferingItem.Version.CRN

		// if imageCatalogOfferingItem.Version.Deleted != nil {
		// 	imageCatalogOfferingVersionDeletedList := []map[string]interface{}{}
		// 	imageCatalogOfferingVersionDeletedMap := map[string]interface{}{}
		// 	imageCatalogOfferingVersionDeletedMap[isImageCatalogOfferingMoreInfo] = imageCatalogOfferingItem.Version.Deleted.MoreInfo
		// 	imageCatalogOfferingVersionDeletedList = append(imageCatalogOfferingVersionDeletedList, imageCatalogOfferingVersionDeletedMap)
		// 	imageCatalogOfferingVersionMap[isImageCatalogOfferingDeleted] = imageCatalogOfferingVersionDeletedList
		// }
		imageCatalogOfferingVersionList = append(imageCatalogOfferingVersionList, imageCatalogOfferingVersionMap)
		imageCatalogOfferingMap[isImageCatalogOfferingVersion] = imageCatalogOfferingVersionList
	}

	return imageCatalogOfferingMap
}

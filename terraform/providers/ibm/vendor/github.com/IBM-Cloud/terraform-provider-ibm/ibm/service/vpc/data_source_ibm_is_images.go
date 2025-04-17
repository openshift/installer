// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isImages                = "images"
	isImagesResourceGroupID = "resource_group"
	isImageCatalogManaged   = "catalog_managed"
)

func DataSourceIBMISImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISImagesRead,

		Schema: map[string]*schema.Schema{
			isImagesResourceGroupID: {
				Type:        schema.TypeString,
				Description: "The id of the resource group",
				Optional:    true,
			},
			isImageCatalogManaged: {
				Type:        schema.TypeBool,
				Description: "Lists images managed as part of a catalog offering. If an image is managed, accounts in the same enterprise with access to that catalog can specify the image's catalog offering version CRN to provision virtual server instances using the image",
				Optional:    true,
			},
			isImageName: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_image", isImageName),
				Description:  "The name of the image",
			},
			isImageStatus: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_images", isImageStatus),
				Description:  "The status of the image",
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
												isImageCatalogOfferingDeleted: {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															isImageCatalogOfferingMoreInfo: {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Link to documentation about deleted resources.",
															},
														},
													},
												},
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
				},
			},
		},
	}
}

func DataSourceIBMISImagesValidator() *validate.ResourceValidator {

	status := "available, deleting, deprecated, failed, obsolete, pending, tentative, unusable"
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isImageStatus,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              status})
	ibmISImageResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_images", Schema: validateSchema}
	return &ibmISImageResourceValidator
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

	var status string
	if v, ok := d.GetOk(isImageStatus); ok {
		status = v.(string)
	}
	var catalogManaged bool
	if v, ok := d.GetOk(isImageCatalogManaged); ok {
		catalogManaged = v.(bool)
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
			return fmt.Errorf("[ERROR] Error Fetching Images %s\n%s", err, response)
		}
		start = flex.GetNext(availableImages.Next)
		allrecs = append(allrecs, availableImages.Images...)
		if start == "" {
			break
		}
	}

	if status != "" {
		allrecsTemp := []vpcv1.Image{}
		for _, image := range allrecs {
			if status == *image.Status {
				allrecsTemp = append(allrecsTemp, image)
			}
		}
		allrecs = allrecsTemp
	}

	if catalogManaged {
		allrecsTemp := []vpcv1.Image{}
		for _, image := range allrecs {
			if image.CatalogOffering != nil && catalogManaged == *image.CatalogOffering.Managed {
				allrecsTemp = append(allrecsTemp, image)
			}
		}
		allrecs = allrecsTemp
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
		if image.CatalogOffering != nil {
			catalogOfferingList := []map[string]interface{}{}
			catalogOfferingMap := dataSourceImageCollectionCatalogOfferingToMap(*image.CatalogOffering)
			catalogOfferingList = append(catalogOfferingList, catalogOfferingMap)
			l[isImageCatalogOffering] = catalogOfferingList
		}
		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *image.CRN, "", isImageAccessTagType)
		if err != nil {
			log.Printf(
				"Error on get of resource image (%s) access tags: %s", d.Id(), err)
		}
		l[isImageAccessTags] = accesstags
		imagesInfo = append(imagesInfo, l)
	}
	d.SetId(dataSourceIBMISImagesID(d))
	d.Set(isImages, imagesInfo)
	return nil
}

// dataSourceIBMISImagesId returns a reasonable ID for a image list.
func dataSourceIBMISImagesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

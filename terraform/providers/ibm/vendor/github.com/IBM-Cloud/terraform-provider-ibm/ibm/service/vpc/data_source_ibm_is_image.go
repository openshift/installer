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
			"resource_group": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group for this IPsec policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this resource group.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this resource group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this resource group.",
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of this image",
			},
			"status_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current status (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about this status reason.",
						},
					},
				},
			},
			"operating_system": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operating system architecture",
						},
						"dedicated_host_only": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Images with this operating system can only be used on dedicated hosts or dedicated host groups",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A unique, display-friendly name for the operating system",
						},
						"family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The software family for this operating system",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this operating system",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this operating system",
						},
						"vendor": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vendor of the operating system",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The major release version of this operating system",
						},
					},
				},
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
	if len(image.StatusReasons) > 0 {
		d.Set("status_reasons", dataSourceIBMIsImageFlattenStatusReasons(image.StatusReasons))
	}
	d.Set("name", *image.Name)
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *image.CRN, "", isImageAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource image (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isImageAccessTags, accesstags)
	d.Set("visibility", *image.Visibility)

	if image.OperatingSystem != nil {
		operatingSystemList := []map[string]interface{}{}
		operatingSystemMap := dataSourceIBMISImageOperatingSystemToMap(*image.OperatingSystem)
		operatingSystemList = append(operatingSystemList, operatingSystemMap)
		d.Set("operating_system", operatingSystemList)
	}
	if image.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourceImageResourceGroupToMap(*image.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		d.Set("resource_group", resourceGroupList)
	}
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
	if len(image.StatusReasons) > 0 {
		d.Set("status_reasons", dataSourceIBMIsImageFlattenStatusReasons(image.StatusReasons))
	}
	d.Set("name", *image.Name)
	d.Set("visibility", *image.Visibility)
	if image.OperatingSystem != nil {
		operatingSystemList := []map[string]interface{}{}
		operatingSystemMap := dataSourceIBMISImageOperatingSystemToMap(*image.OperatingSystem)
		operatingSystemList = append(operatingSystemList, operatingSystemMap)
		d.Set("operating_system", operatingSystemList)
	}
	if image.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourceImageResourceGroupToMap(*image.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		d.Set("resource_group", resourceGroupList)
	}
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

func dataSourceIBMISImageOperatingSystemToMap(operatingSystemItem vpcv1.OperatingSystem) (operatingSystemMap map[string]interface{}) {
	operatingSystemMap = map[string]interface{}{}

	if operatingSystemItem.Architecture != nil {
		operatingSystemMap["architecture"] = operatingSystemItem.Architecture
	}
	if operatingSystemItem.DedicatedHostOnly != nil {
		operatingSystemMap["dedicated_host_only"] = operatingSystemItem.DedicatedHostOnly
	}
	if operatingSystemItem.DisplayName != nil {
		operatingSystemMap["display_name"] = operatingSystemItem.DisplayName
	}
	if operatingSystemItem.Family != nil {
		operatingSystemMap["family"] = operatingSystemItem.Family
	}
	if operatingSystemItem.Href != nil {
		operatingSystemMap["href"] = operatingSystemItem.Href
	}
	if operatingSystemItem.Name != nil {
		operatingSystemMap["name"] = operatingSystemItem.Name
	}
	if operatingSystemItem.Vendor != nil {
		operatingSystemMap["vendor"] = operatingSystemItem.Vendor
	}
	if operatingSystemItem.Version != nil {
		operatingSystemMap["version"] = operatingSystemItem.Version
	}
	return operatingSystemMap
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

func dataSourceIBMIsImageFlattenStatusReasons(result []vpcv1.ImageStatusReason) (statusReasons []map[string]interface{}) {
	for _, statusReasonsItem := range result {
		statusReasons = append(statusReasons, dataSourceIBMIsImageStatusReasonToMap(&statusReasonsItem))
	}

	return statusReasons
}

func dataSourceIBMIsImageStatusReasonToMap(model *vpcv1.ImageStatusReason) map[string]interface{} {
	modelMap := make(map[string]interface{})
	if model.Code != nil {
		modelMap["code"] = *model.Code
	}
	if model.Message != nil {
		modelMap["message"] = *model.Message
	}
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap
}
func dataSourceImageResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
	resourceGroupMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resourceGroupMap["href"] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resourceGroupMap["id"] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resourceGroupMap["name"] = resourceGroupItem.Name
	}

	return resourceGroupMap
}

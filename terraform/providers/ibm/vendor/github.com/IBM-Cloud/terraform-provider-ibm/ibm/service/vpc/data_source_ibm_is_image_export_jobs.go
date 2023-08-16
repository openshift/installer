// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsImageExports() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMIsImageExportsRead,

		Schema: map[string]*schema.Schema{
			"image": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The image identifier.",
			},
			"export_jobs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of image export jobs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"completed_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the image export job was completed.If absent, the export job has not yet completed.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the image export job was created.",
						},
						"encrypted_data_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A base64-encoded, encrypted representation of the key that was used to encrypt the data for the exported image. This key can be unwrapped with the image's `encryption_key` root key using either Key Protect or Hyper Protect Crypto Service.If absent, the export job is for an unencrypted image.",
						},
						"format": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The format of the exported image.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this image export job.",
						},
						"image_export_job": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this image export job.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this image export job.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
						"started_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the image export job started running.If absent, the export job has not yet started.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of this image export job:- `deleting`: Export job is being deleted- `failed`: Export job could not be completed successfully- `queued`: Export job is queued- `running`: Export job is in progress- `succeeded`: Export job was completed successfullyThe exported image object is automatically deleted for `failed` jobs.",
						},
						"status_reasons": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current status (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
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
						"storage_bucket": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Cloud Object Storage bucket of the exported image object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN of this Cloud Object Storage bucket.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name of this Cloud Object Storage bucket.",
									},
								},
							},
						},
						"storage_href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Cloud Object Storage location of the exported image object. The object at this location may not exist until the job is started, and will be incomplete while the job is running.After the job completes, the exported image object is not managed by the IBM VPC service, and may be removed or replaced with a different object by any user or service with IAM authorization to the bucket.",
						},
						"storage_object": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Cloud Object Storage object for the exported image. This object may not exist untilthe job is started, and will not be complete until the job completes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of this Cloud Object Storage object. Names are unique within a Cloud Object Storage bucket.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMIsImageExportsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listImageExportJobsOptions := &vpcv1.ListImageExportJobsOptions{}

	listImageExportJobsOptions.SetImageID(d.Get("image").(string))

	imageExportJobUnpaginatedCollection, response, err := vpcClient.ListImageExportJobsWithContext(context, listImageExportJobsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListImageExportJobsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListImageExportJobsWithContext failed %s\n%s", err, response))
	}

	d.SetId(DataSourceIBMIsImageExportsID(d))

	exportJobs := []map[string]interface{}{}
	if imageExportJobUnpaginatedCollection.ExportJobs != nil {
		for _, modelItem := range imageExportJobUnpaginatedCollection.ExportJobs {
			modelMap, err := DataSourceIBMIsImageExportsImageExportJobToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			exportJobs = append(exportJobs, modelMap)
		}
	}
	if err = d.Set("export_jobs", exportJobs); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting export_jobs %s", err))
	}

	return nil
}

// DataSourceIBMIsImageExportsID returns a reasonable ID for the list.
func DataSourceIBMIsImageExportsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsImageExportsImageExportJobToMap(model *vpcv1.ImageExportJob) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CompletedAt != nil {
		modelMap["completed_at"] = model.CompletedAt.String()
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.EncryptedDataKey != nil {
	}
	if model.Format != nil {
		modelMap["format"] = *model.Format
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["image_export_job"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	if model.StartedAt != nil {
		modelMap["started_at"] = model.StartedAt.String()
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.StatusReasons != nil {
		statusReasons := []map[string]interface{}{}
		for _, statusReasonsItem := range model.StatusReasons {
			statusReasonsItemMap, err := DataSourceIBMIsImageExportsImageExportJobStatusReasonToMap(&statusReasonsItem)
			if err != nil {
				return modelMap, err
			}
			statusReasons = append(statusReasons, statusReasonsItemMap)
		}
		modelMap["status_reasons"] = statusReasons
	}
	if model.StorageBucket != nil {
		storageBucketMap, err := DataSourceIBMIsImageExportsCloudObjectStorageBucketReferenceToMap(model.StorageBucket)
		if err != nil {
			return modelMap, err
		}
		modelMap["storage_bucket"] = []map[string]interface{}{storageBucketMap}
	}
	if model.StorageHref != nil {
		modelMap["storage_href"] = *model.StorageHref
	}
	if model.StorageObject != nil {
		storageObjectMap, err := DataSourceIBMIsImageExportsCloudObjectStorageObjectReferenceToMap(model.StorageObject)
		if err != nil {
			return modelMap, err
		}
		modelMap["storage_object"] = []map[string]interface{}{storageObjectMap}
	}
	return modelMap, nil
}

func DataSourceIBMIsImageExportsImageExportJobStatusReasonToMap(model *vpcv1.ImageExportJobStatusReason) (map[string]interface{}, error) {
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
	return modelMap, nil
}

func DataSourceIBMIsImageExportsCloudObjectStorageBucketReferenceToMap(model *vpcv1.CloudObjectStorageBucketReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func DataSourceIBMIsImageExportsCloudObjectStorageObjectReferenceToMap(model *vpcv1.CloudObjectStorageObjectReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

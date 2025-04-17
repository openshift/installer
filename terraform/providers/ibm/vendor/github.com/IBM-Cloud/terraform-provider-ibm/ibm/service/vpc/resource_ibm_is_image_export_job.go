// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsImageExportJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIBMIsImageExportCreate,
		ReadContext:   ResourceIBMIsImageExportRead,
		UpdateContext: ResourceIBMIsImageExportUpdate,
		DeleteContext: ResourceIBMIsImageExportDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"image": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The image identifier.",
			},
			"storage_bucket": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "The name of the Cloud Object Storage bucket to export the image to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ExactlyOneOf: []string{"storage_bucket.0.name", "storage_bucket.0.crn"},
							Description:  "Name of this Cloud Object Storage bucket",
						},
						"crn": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ExactlyOneOf: []string{"storage_bucket.0.name", "storage_bucket.0.crn"},
							Description:  "CRN of this Cloud Object Storage bucket",
						},
					},
				},
			},
			"format": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "qcow2",
				ValidateFunc: validate.InvokeValidator("ibm_is_image_export_job", "format"),
				Description:  "The format to use for the exported image. If the image is encrypted, only `qcow2` is supported.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_image_export_job", "name"),
				Description:  "The user-defined name for this image export job. Names must be unique within the image this export job resides in. If unspecified, the name will be a hyphenated list of randomly-selected words prefixed with the first 16 characters of the parent image name.The exported image object name in Cloud Object Storage (`storage_object.name` in the response) will be based on this name. The object name will be unique within the bucket.",
			},
			"completed_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image export job was completed.If absent, the export job has not yet completed.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image export job was created.",
			},
			"encrypted_data_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A base64-encoded, encrypted representation of the key that was used to encrypt the data for the exported image. This key can be unwrapped with the image's `encryption_key` root key using either Key Protect or Hyper Protect Crypto Service.If absent, the export job is for an unencrypted image.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this image export job.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
			"started_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image export job started running.If absent, the export job has not yet started.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of this image export job:- `deleting`: Export job is being deleted- `failed`: Export job could not be completed successfully- `queued`: Export job is queued- `running`: Export job is in progress- `succeeded`: Export job was completed successfullyThe exported image object is automatically deleted for `failed` jobs.",
			},
			"status_reasons": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current status (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A snake case string succinctly identifying the status reason.",
						},
						"message": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "An explanation of the status reason.",
						},
						"more_info": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Link to documentation about this status reason.",
						},
					},
				},
			},
			"storage_href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Cloud Object Storage location of the exported image object. The object at this location may not exist until the job is started, and will be incomplete while the job is running.After the job completes, the exported image object is not managed by the IBM VPC service, and may be removed or replaced with a different object by any user or service with IAM authorization to the bucket.",
			},
			"storage_object": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Cloud Object Storage object for the exported image. This object may not exist untilthe job is started, and will not be complete until the job completes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of this Cloud Object Storage object. Names are unique within a Cloud Object Storage bucket.",
						},
					},
				},
			},
			"image_export_job": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this image export job.",
			},
		},
	}
}

func ResourceIBMIsImageExportValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 1)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "format",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "qcow2, vhd",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_image_export_job", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIBMIsImageExportCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createImageExportJobOptions := &vpcv1.CreateImageExportJobOptions{}

	createImageExportJobOptions.SetImageID(d.Get("image").(string))
	storageBucketMap := d.Get("storage_bucket").([]interface{})[0].(map[string]interface{})
	storage_bucket_name := storageBucketMap["name"].(string)
	storage_bucket_crn := storageBucketMap["crn"].(string)
	storageBucket := &vpcv1.CloudObjectStorageBucketIdentity{}
	if storage_bucket_crn != "" {
		storageBucket.CRN = &storage_bucket_crn
	} else {
		storageBucket.Name = &storage_bucket_name
	}
	createImageExportJobOptions.SetStorageBucket(storageBucket)
	if format, ok := d.GetOk("format"); ok {
		createImageExportJobOptions.SetFormat(format.(string))
	}
	if name, ok := d.GetOk("name"); ok {
		createImageExportJobOptions.SetName(name.(string))
	}

	imageExportJob, response, err := vpcClient.CreateImageExportJobWithContext(context, createImageExportJobOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateImageExportJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateImageExportJobWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createImageExportJobOptions.ImageID, *imageExportJob.ID))

	return ResourceIBMIsImageExportRead(context, d, meta)
}

func ResourceIBMIsImageExportRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getImageExportJobOptions := &vpcv1.GetImageExportJobOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getImageExportJobOptions.SetImageID(parts[0])
	getImageExportJobOptions.SetID(parts[1])

	imageExportJob, response, err := vpcClient.GetImageExportJobWithContext(context, getImageExportJobOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetImageExportJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetImageExportJobWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("format", imageExportJob.Format); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting format: %s", err))
	}
	if err = d.Set("name", imageExportJob.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("completed_at", flex.DateTimeToString(imageExportJob.CompletedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting completed_at: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(imageExportJob.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("encrypted_data_key", imageExportJob.EncryptedDataKey); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting encrypted_data_key: %s", err))
	}
	if err = d.Set("href", imageExportJob.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("resource_type", imageExportJob.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if err = d.Set("started_at", flex.DateTimeToString(imageExportJob.StartedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting started_at: %s", err))
	}
	if err = d.Set("status", imageExportJob.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}

	if imageExportJob.StorageBucket != nil {
		storageBucketList := []map[string]interface{}{}
		storageBucketMap := map[string]interface{}{
			"name": *imageExportJob.StorageBucket.Name,
			"crn":  *imageExportJob.StorageBucket.CRN,
		}

		storageBucketList = append(storageBucketList, storageBucketMap)
		d.Set("storage_bucket", storageBucketList)
	}

	statusReasons := []map[string]interface{}{}
	for _, statusReasonsItem := range imageExportJob.StatusReasons {
		statusReasonsItemMap, err := ResourceIBMIsImageExportImageExportJobStatusReasonToMap(&statusReasonsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		statusReasons = append(statusReasons, statusReasonsItemMap)
	}
	if err = d.Set("status_reasons", statusReasons); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status_reasons: %s", err))
	}
	if err = d.Set("storage_href", imageExportJob.StorageHref); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting storage_href: %s", err))
	}
	storageObjectMap, err := ResourceIBMIsImageExportCloudObjectStorageObjectReferenceToMap(imageExportJob.StorageObject)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("storage_object", []map[string]interface{}{storageObjectMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting storage_object: %s", err))
	}
	if err = d.Set("image_export_job", imageExportJob.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting image_export_job: %s", err))
	}

	return nil
}

func ResourceIBMIsImageExportUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateImageExportJobOptions := &vpcv1.UpdateImageExportJobOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateImageExportJobOptions.SetImageID(parts[0])
	updateImageExportJobOptions.SetID(parts[1])

	hasChange := false

	patchVals := &vpcv1.ImageExportJobPatch{}

	if d.HasChange("name") {
		patchVals.Name = core.StringPtr(d.Get("name").(string))
		hasChange = true
	}

	if hasChange {
		updateImageExportJobOptions.ImageExportJobPatch, _ = patchVals.AsPatch()
		_, response, err := vpcClient.UpdateImageExportJobWithContext(context, updateImageExportJobOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateImageExportJobWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateImageExportJobWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIBMIsImageExportRead(context, d, meta)
}

func ResourceIBMIsImageExportDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteImageExportJobOptions := &vpcv1.DeleteImageExportJobOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteImageExportJobOptions.SetImageID(parts[0])
	deleteImageExportJobOptions.SetID(parts[1])

	response, err := vpcClient.DeleteImageExportJobWithContext(context, deleteImageExportJobOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteImageExportJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteImageExportJobWithContext failed %s\n%s", err, response))
	}
	_, err = isWaitForImageExportJobDeleted(context, d, meta, vpcClient, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func ResourceIBMIsImageExportMapToCloudObjectStorageBucketIdentity(modelMap map[string]interface{}) (vpcv1.CloudObjectStorageBucketIdentityIntf, error) {
	model := &vpcv1.CloudObjectStorageBucketIdentity{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	return model, nil
}

func ResourceIBMIsImageExportMapToCloudObjectStorageBucketIdentityCloudObjectStorageBucketIdentityByName(modelMap map[string]interface{}) (*vpcv1.CloudObjectStorageBucketIdentityCloudObjectStorageBucketIdentityByName, error) {
	model := &vpcv1.CloudObjectStorageBucketIdentityCloudObjectStorageBucketIdentityByName{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	return model, nil
}

func ResourceIBMIsImageExportMapToCloudObjectStorageBucketIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.CloudObjectStorageBucketIdentityByCRN, error) {
	model := &vpcv1.CloudObjectStorageBucketIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func ResourceIBMIsImageExportCloudObjectStorageBucketIdentityToMap(model vpcv1.CloudObjectStorageBucketIdentityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.CloudObjectStorageBucketIdentityCloudObjectStorageBucketIdentityByName); ok {
		return ResourceIBMIsImageExportCloudObjectStorageBucketIdentityCloudObjectStorageBucketIdentityByNameToMap(model.(*vpcv1.CloudObjectStorageBucketIdentityCloudObjectStorageBucketIdentityByName))
	} else if _, ok := model.(*vpcv1.CloudObjectStorageBucketIdentityByCRN); ok {
		return ResourceIBMIsImageExportCloudObjectStorageBucketIdentityByCRNToMap(model.(*vpcv1.CloudObjectStorageBucketIdentityByCRN))
	} else if _, ok := model.(*vpcv1.CloudObjectStorageBucketIdentity); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.CloudObjectStorageBucketIdentity)
		if model.Name != nil {
			modelMap["name"] = model.Name
		}
		if model.CRN != nil {
			modelMap["crn"] = model.CRN
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.CloudObjectStorageBucketIdentityIntf subtype encountered")
	}
}

func ResourceIBMIsImageExportCloudObjectStorageBucketIdentityCloudObjectStorageBucketIdentityByNameToMap(model *vpcv1.CloudObjectStorageBucketIdentityCloudObjectStorageBucketIdentityByName) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	return modelMap, nil
}

func ResourceIBMIsImageExportCloudObjectStorageBucketIdentityByCRNToMap(model *vpcv1.CloudObjectStorageBucketIdentityByCRN) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	return modelMap, nil
}

func ResourceIBMIsImageExportImageExportJobStatusReasonToMap(model *vpcv1.ImageExportJobStatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = model.Code
	modelMap["message"] = model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}

func ResourceIBMIsImageExportCloudObjectStorageObjectReferenceToMap(model *vpcv1.CloudObjectStorageObjectReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	return modelMap, nil
}

func isWaitForImageExportJobDeleted(context context.Context, d *schema.ResourceData, meta interface{}, vpcClient *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for image export job (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", "deleting"},
		Target:     []string{"", "done"},
		Refresh:    isImageExportJobDeleteRefreshFunc(context, d, meta, vpcClient, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isImageExportJobDeleteRefreshFunc(context context.Context, d *schema.ResourceData, meta interface{}, vpcClient *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] is image export job delete refresh here")
		parts, err := flex.SepIdParts(d.Id(), "/")
		if err != nil {
			return nil, "", err
		}
		getImgExpJobOptions := &vpcv1.GetImageExportJobOptions{}

		getImgExpJobOptions.SetImageID(parts[0])
		getImgExpJobOptions.SetID(parts[1])

		imageExportJob, response, err := vpcClient.GetImageExportJob(getImgExpJobOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return imageExportJob, "done", nil
			}
			return imageExportJob, "", fmt.Errorf("[ERROR] Error Getting Image export job: %s\n%s", err, response)
		}
		return imageExportJob, *imageExportJob.Status, err
	}
}

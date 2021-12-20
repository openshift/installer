// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/errors"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_images"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

func resourceIBMPIImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIImageCreate,
		ReadContext:   resourceIBMPIImageRead,
		UpdateContext: resourceIBMPIImageUpdate,
		DeleteContext: resourceIBMPIImageDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI cloud instance ID",
			},
			helpers.PIImageName: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Image name",
				DiffSuppressFunc: applyOnce,
			},
			helpers.PIImageId: {
				Type:             schema.TypeString,
				Optional:         true,
				ExactlyOneOf:     []string{helpers.PIImageId, helpers.PIImageBucketName},
				Description:      "Instance image id",
				DiffSuppressFunc: applyOnce,
				ConflictsWith:    []string{helpers.PIImageBucketName},
				ForceNew:         true,
			},

			// COS import variables
			helpers.PIImageBucketName: {
				Type:          schema.TypeString,
				Optional:      true,
				ExactlyOneOf:  []string{helpers.PIImageId, helpers.PIImageBucketName},
				Description:   "Cloud Object Storage bucket name; bucket-name[/optional/folder]",
				ConflictsWith: []string{helpers.PIImageId},
				RequiredWith:  []string{helpers.PIImageBucketRegion, helpers.PIImageBucketFileName},
				ForceNew:      true,
			},
			helpers.PIImageBucketAccess: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Indicates if the bucket has public or private access",
				Default:       "public",
				ValidateFunc:  validateAllowedStringValue([]string{"public", "private"}),
				ConflictsWith: []string{helpers.PIImageId},
				ForceNew:      true,
			},
			helpers.PIImageAccessKey: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Cloud Object Storage access key; required for buckets with private access",
				ForceNew:     true,
				Sensitive:    true,
				RequiredWith: []string{helpers.PIImageSecretKey},
			},
			helpers.PIImageSecretKey: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Cloud Object Storage secret key; required for buckets with private access",
				ForceNew:     true,
				Sensitive:    true,
				RequiredWith: []string{helpers.PIImageAccessKey},
			},
			helpers.PIImageBucketRegion: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Cloud Object Storage region",
				ConflictsWith: []string{helpers.PIImageId},
				RequiredWith:  []string{helpers.PIImageBucketName},
				ForceNew:      true,
			},
			helpers.PIImageBucketFileName: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Cloud Object Storage image filename",
				ConflictsWith: []string{helpers.PIImageId},
				RequiredWith:  []string{helpers.PIImageBucketName},
				ForceNew:      true,
			},
			helpers.PIImageStorageType: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Type of storage",
				RequiredWith: []string{helpers.PIImageBucketName},
				ForceNew:     true,
			},

			// Computed Attribute
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image ID",
			},
		},
	}
}

func resourceIBMPIImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		log.Printf("Failed to get the session")
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	imageName := d.Get(helpers.PIImageName).(string)

	client := st.NewIBMPIImageClient(ctx, sess, cloudInstanceID)
	// image copy
	if v, ok := d.GetOk(helpers.PIImageId); ok {
		imageid := v.(string)
		source := "root-project"
		var body = &models.CreateImage{
			ImageName: imageName,
			ImageID:   imageid,
			Source:    &source,
		}
		imageResponse, err := client.Create(body)
		if err != nil {
			return diag.FromErr(err)
		}

		IBMPIImageID := imageResponse.ImageID
		d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *IBMPIImageID))

		_, err = isWaitForIBMPIImageAvailable(ctx, client, *IBMPIImageID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			log.Printf("[DEBUG]  err %s", err)
			return diag.FromErr(err)
		}
	}

	// COS image import
	if v, ok := d.GetOk(helpers.PIImageBucketName); ok {
		bucketName := v.(string)
		bucketImageFileName := d.Get(helpers.PIImageBucketFileName).(string)
		bucketRegion := d.Get(helpers.PIImageBucketRegion).(string)
		bucketAccess := d.Get(helpers.PIImageBucketAccess).(string)
		storageType := d.Get(helpers.PIImageStorageType).(string)

		body := &models.CreateCosImageImportJob{
			ImageName:     &imageName,
			BucketName:    &bucketName,
			BucketAccess:  &bucketAccess,
			ImageFilename: &bucketImageFileName,
			Region:        &bucketRegion,
			StorageType:   storageType,
		}

		if v, ok := d.GetOk(helpers.PIImageAccessKey); ok {
			body.AccessKey = v.(string)
		}
		if v, ok := d.GetOk(helpers.PIImageSecretKey); ok {
			body.SecretKey = v.(string)
		}

		imageResponse, err := client.CreateCosImage(body)
		if err != nil {
			return diag.FromErr(err)
		}

		jobClient := st.NewIBMPIJobClient(ctx, sess, cloudInstanceID)
		_, err = waitForIBMPIJobCompleted(ctx, jobClient, *imageResponse.ID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		// Once the job is completed find by name
		image, err := client.Get(imageName)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *image.ImageID))
	}

	return resourceIBMPIImageRead(ctx, d, meta)
}

func resourceIBMPIImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, imageID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	imageC := st.NewIBMPIImageClient(ctx, sess, cloudInstanceID)
	imagedata, err := imageC.Get(imageID)
	if err != nil {
		uErr := errors.Unwrap(err)
		switch uErr.(type) {
		case *p_cloud_images.PcloudCloudinstancesImagesGetNotFound:
			log.Printf("[DEBUG] image does not exist %v", err)
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] get image failed %v", err)
		return diag.FromErr(err)
	}

	imageid := *imagedata.ImageID
	d.Set("image_id", imageid)
	d.Set(helpers.PICloudInstanceId, cloudInstanceID)

	return nil
}

func resourceIBMPIImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceIBMPIImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, imageID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	imageC := st.NewIBMPIImageClient(ctx, sess, cloudInstanceID)
	err = imageC.Delete(imageID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func isWaitForIBMPIImageAvailable(ctx context.Context, client *st.IBMPIImageClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Power Image (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PIImageQueStatus},
		Target:     []string{helpers.PIImageActiveStatus},
		Refresh:    isIBMPIImageRefreshFunc(ctx, client, id),
		Timeout:    timeout,
		Delay:      20 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIImageRefreshFunc(ctx context.Context, client *st.IBMPIImageClient, id string) resource.StateRefreshFunc {

	log.Printf("Calling the isIBMPIImageRefreshFunc Refresh Function....")
	return func() (interface{}, string, error) {
		image, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if image.State == "active" {
			return image, helpers.PIImageActiveStatus, nil
		}

		return image, helpers.PIImageQueStatus, nil
	}
}

func waitForIBMPIJobCompleted(ctx context.Context, client *st.IBMPIJobClient, jobID string, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{helpers.JobStatusQueued, helpers.JobStatusReadyForProcessing, helpers.JobStatusInProgress, helpers.JobStatusRunning, helpers.JobStatusWaiting},
		Target:  []string{helpers.JobStatusCompleted, helpers.JobStatusFailed},
		Refresh: func() (interface{}, string, error) {
			job, err := client.Get(jobID)
			if err != nil {
				log.Printf("[DEBUG] get job failed %v", err)
				return nil, "", fmt.Errorf(errors.GetJobOperationFailed, jobID, err)
			}
			if job == nil || job.Status == nil {
				log.Printf("[DEBUG] get job failed with empty response")
				return nil, "", fmt.Errorf("failed to get job status for job id %s", jobID)
			}
			if *job.Status.State == helpers.JobStatusFailed {
				log.Printf("[DEBUG] job status failed with message: %v", job.Status.Message)
				return nil, helpers.JobStatusFailed, fmt.Errorf("job status failed for job id %s with message: %v", jobID, job.Status.Message)
			}
			return job, *job.Status.State, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

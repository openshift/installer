// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/errors"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_images"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func ResourceIBMPIImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIImageCreate,
		ReadContext:   resourceIBMPIImageRead,
		DeleteContext: resourceIBMPIImageDelete,
		UpdateContext: resourceIBMPIImageUpdate,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI cloud instance ID",
				ForceNew:    true,
			},
			helpers.PIImageName: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Image name",
				DiffSuppressFunc: flex.ApplyOnce,
				ForceNew:         true,
			},
			helpers.PIImageId: {
				Type:             schema.TypeString,
				Optional:         true,
				ExactlyOneOf:     []string{helpers.PIImageId, helpers.PIImageBucketName},
				Description:      "Instance image id",
				DiffSuppressFunc: flex.ApplyOnce,
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
				ValidateFunc:  validate.ValidateAllowedStringValues([]string{"public", "private"}),
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of storage; If not specified, default is tier3",
				ForceNew:    true,
			},
			helpers.PIImageStoragePool: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Storage pool where the image will be loaded, if provided then pi_affinity_policy will be ignored",
				ForceNew:    true,
			},
			Arg_AffinityPolicy: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Affinity policy for image; ignored if pi_image_storage_pool provided; for policy affinity requires one of pi_affinity_instance or pi_affinity_volume to be specified; for policy anti-affinity requires one of pi_anti_affinity_instances or pi_anti_affinity_volumes to be specified",
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"affinity", "anti-affinity"}),
				ForceNew:     true,
			},
			Arg_AffinityVolume: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Volume (ID or Name) to base storage affinity policy against; required if requesting affinity and pi_affinity_instance is not provided",
				ConflictsWith: []string{Arg_AffinityInstance},
				ForceNew:      true,
			},
			Arg_AffinityInstance: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "PVM Instance (ID or Name) to base storage affinity policy against; required if requesting storage affinity and pi_affinity_volume is not provided",
				ConflictsWith: []string{Arg_AffinityVolume},
				ForceNew:      true,
			},
			Arg_AntiAffinityVolumes: {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "List of volumes to base storage anti-affinity policy against; required if requesting anti-affinity and pi_anti_affinity_instances is not provided",
				ConflictsWith: []string{Arg_AntiAffinityInstances},
				ForceNew:      true,
			},
			Arg_AntiAffinityInstances: {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "List of pvmInstances to base storage anti-affinity policy against; required if requesting anti-affinity and pi_anti_affinity_volumes is not provided",
				ConflictsWith: []string{Arg_AntiAffinityVolumes},
				ForceNew:      true,
			},
			Arg_ImageImportDetails: {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_LicenseType: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.ValidateAllowedStringValues([]string{BYOL}),
							Description:  "Origin of the license of the product.",
						},
						Attr_Product: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.ValidateAllowedStringValues([]string{Hana, Netweaver}),
							Description:  "Product within the image.",
						},
						Attr_Vendor: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.ValidateAllowedStringValues([]string{SAP}),
							Description:  "Vendor supporting the product.",
						},
					},
				},
			},
			Arg_UserTags: {
				Description: "The user tags attached to this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},

			// Computed Attribute
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image ID",
			},
		},
	}
}

func resourceIBMPIImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
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
		if tags, ok := d.GetOk(Arg_UserTags); ok {
			body.UserTags = flex.FlattenSet(tags.(*schema.Set))
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

		if _, ok := d.GetOk(Arg_UserTags); ok {
			if imageResponse.Crn != "" {
				oldList, newList := d.GetChange(Arg_UserTags)
				err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, string(imageResponse.Crn), "", UserTagType)
				if err != nil {
					log.Printf("Error on update of pi image (%s) pi_user_tags during creation: %s", *IBMPIImageID, err)
				}
			}
		}
	}

	// COS image import
	if v, ok := d.GetOk(helpers.PIImageBucketName); ok {
		bucketName := v.(string)
		bucketImageFileName := d.Get(helpers.PIImageBucketFileName).(string)
		bucketRegion := d.Get(helpers.PIImageBucketRegion).(string)
		bucketAccess := d.Get(helpers.PIImageBucketAccess).(string)

		body := &models.CreateCosImageImportJob{
			ImageName:     &imageName,
			BucketName:    &bucketName,
			BucketAccess:  &bucketAccess,
			ImageFilename: &bucketImageFileName,
			Region:        &bucketRegion,
		}

		if v, ok := d.GetOk(helpers.PIImageAccessKey); ok {
			body.AccessKey = v.(string)
		}
		if v, ok := d.GetOk(helpers.PIImageSecretKey); ok {
			body.SecretKey = v.(string)
		}

		if v, ok := d.GetOk(helpers.PIImageStorageType); ok {
			body.StorageType = v.(string)
		}
		if v, ok := d.GetOk(helpers.PIImageStoragePool); ok {
			body.StoragePool = v.(string)
		}
		if ap, ok := d.GetOk(Arg_AffinityPolicy); ok {
			policy := ap.(string)
			affinity := &models.StorageAffinity{
				AffinityPolicy: &policy,
			}

			if policy == "affinity" {
				if av, ok := d.GetOk(Arg_AffinityVolume); ok {
					afvol := av.(string)
					affinity.AffinityVolume = &afvol
				}
				if ai, ok := d.GetOk(Arg_AffinityInstance); ok {
					afins := ai.(string)
					affinity.AffinityPVMInstance = &afins
				}
			} else {
				if avs, ok := d.GetOk(Arg_AntiAffinityVolumes); ok {
					afvols := flex.ExpandStringList(avs.([]interface{}))
					affinity.AntiAffinityVolumes = afvols
				}
				if ais, ok := d.GetOk(Arg_AntiAffinityInstances); ok {
					afinss := flex.ExpandStringList(ais.([]interface{}))
					affinity.AntiAffinityPVMInstances = afinss
				}
			}
			body.StorageAffinity = affinity
		}
		if _, ok := d.GetOk(Arg_ImageImportDetails); ok {
			details := d.Get(Arg_ImageImportDetails + ".0").(map[string]interface{})
			importDetailsModel := models.ImageImportDetails{
				LicenseType: core.StringPtr(details[Attr_LicenseType].(string)),
				Product:     core.StringPtr(details[Attr_Product].(string)),
				Vendor:      core.StringPtr(details[Attr_Vendor].(string)),
			}
			body.ImportDetails = &importDetailsModel
		}
		if tags, ok := d.GetOk(Arg_UserTags); ok {
			body.UserTags = flex.FlattenSet(tags.(*schema.Set))
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

		if _, ok := d.GetOk(Arg_UserTags); ok {
			if image.Crn != "" {
				oldList, newList := d.GetChange(Arg_UserTags)
				err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, string(image.Crn), "", UserTagType)
				if err != nil {
					log.Printf("Error on update of pi image (%s) pi_user_tags during creation: %s", *image.ImageID, err)
				}
			}
		}
		d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *image.ImageID))
	}

	return resourceIBMPIImageRead(ctx, d, meta)
}

func resourceIBMPIImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
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
	if imagedata.Crn != "" {
		d.Set(Attr_CRN, imagedata.Crn)
		tags, err := flex.GetGlobalTagsUsingCRN(meta, string(imagedata.Crn), "", UserTagType)
		if err != nil {
			log.Printf("Error on get of image (%s) pi_user_tags: %s", *imagedata.ImageID, err)
		}
		d.Set(Arg_UserTags, tags)
	}
	d.Set("image_id", imageid)
	d.Set(helpers.PICloudInstanceId, cloudInstanceID)

	return nil
}

func resourceIBMPIImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, imageID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange(Arg_UserTags) {
		if crn, ok := d.GetOk(Attr_CRN); ok {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, crn.(string), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi image (%s) pi_user_tags: %s", imageID, err)
			}
		}
	}

	return resourceIBMPIImageRead(ctx, d, meta)
}

func resourceIBMPIImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
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

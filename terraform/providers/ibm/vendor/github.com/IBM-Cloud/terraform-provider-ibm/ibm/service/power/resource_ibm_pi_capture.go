// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_images"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const cloudStorageDestination string = "cloud-storage"
const imageCatalogDestination string = "image-catalog"

func ResourceIBMPICapture() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPICaptureCreate,
		ReadContext:   resourceIBMPICaptureRead,
		DeleteContext: resourceIBMPICaptureDelete,
		UpdateContext: resourceIBMPICaptureUpdate,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(75 * time.Minute),
			Delete: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: " Cloud Instance ID - This is the service_instance_id.",
			},

			helpers.PIInstanceName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance Name of the Power VM",
			},

			helpers.PIInstanceCaptureName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the capture to create. Note : this must be unique",
			},

			helpers.PIInstanceCaptureDestination: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Destination for the deployable image",
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"image-catalog", "cloud-storage", "both"}),
			},

			helpers.PIInstanceCaptureVolumeIds: {
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              schema.HashString,
				ForceNew:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "List of Data volume IDs",
			},

			helpers.PIInstanceCaptureCloudStorageRegion: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "List of Regions to use",
			},

			helpers.PIInstanceCaptureCloudStorageAccessKey: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Name of Cloud Storage Access Key",
			},
			helpers.PIInstanceCaptureCloudStorageSecretKey: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Name of the Cloud Storage Secret Key",
			},
			helpers.PIInstanceCaptureCloudStorageImagePath: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Cloud Storage Image Path (bucket-name [/folder/../..])",
			},
			Arg_UserTags: {
				Description: "List of user tags attached to the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
			// Computed Attribute
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of the resource.",
				Type:        schema.TypeString,
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image ID of Capture Instance",
			},
		},
	}
}

func resourceIBMPICaptureCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get(helpers.PIInstanceName).(string)
	capturename := d.Get(helpers.PIInstanceCaptureName).(string)
	capturedestination := d.Get(helpers.PIInstanceCaptureDestination).(string)
	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	client := st.NewIBMPIInstanceClient(context.Background(), sess, cloudInstanceID)

	captureBody := &models.PVMInstanceCapture{
		CaptureDestination: &capturedestination,
		CaptureName:        &capturename,
	}
	if capturedestination != imageCatalogDestination {
		if v, ok := d.GetOk(helpers.PIInstanceCaptureCloudStorageRegion); ok {
			captureBody.CloudStorageRegion = v.(string)
		} else {
			return diag.Errorf("%s is required when capture destination is %s", helpers.PIInstanceCaptureCloudStorageRegion, capturedestination)
		}
		if v, ok := d.GetOk(helpers.PIInstanceCaptureCloudStorageAccessKey); ok {
			captureBody.CloudStorageAccessKey = v.(string)
		} else {
			return diag.Errorf("%s is required when capture destination is %s ", helpers.PIInstanceCaptureCloudStorageAccessKey, capturedestination)
		}
		if v, ok := d.GetOk(helpers.PIInstanceCaptureCloudStorageImagePath); ok {
			captureBody.CloudStorageImagePath = v.(string)
		} else {
			return diag.Errorf("%s is required when capture destination is %s ", helpers.PIInstanceCaptureCloudStorageImagePath, capturedestination)
		}
		if v, ok := d.GetOk(helpers.PIInstanceCaptureCloudStorageSecretKey); ok {
			captureBody.CloudStorageSecretKey = v.(string)
		} else {
			return diag.Errorf("%s is required when capture destination is %s ", helpers.PIInstanceCaptureCloudStorageSecretKey, capturedestination)
		}
	}

	if v, ok := d.GetOk(helpers.PIInstanceCaptureVolumeIds); ok {
		volids := flex.ExpandStringList((v.(*schema.Set)).List())
		if len(volids) > 0 {
			captureBody.CaptureVolumeIDs = volids
		}
	}

	if v, ok := d.GetOk(Arg_UserTags); ok {
		captureBody.UserTags = flex.FlattenSet(v.(*schema.Set))
	}

	captureResponse, err := client.CaptureInstanceToImageCatalogV2(name, captureBody)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cloudInstanceID, capturename, capturedestination))
	jobClient := st.NewIBMPIJobClient(ctx, sess, cloudInstanceID)
	_, err = waitForIBMPIJobCompleted(ctx, jobClient, *captureResponse.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if _, ok := d.GetOk(Arg_UserTags); ok && capturedestination != cloudStorageDestination {
		imageClient := st.NewIBMPIImageClient(ctx, sess, cloudInstanceID)
		imagedata, err := imageClient.Get(capturename)
		if err != nil {
			if strings.Contains(err.Error(), NotFound) {
				d.SetId("")
			}
			return diag.Errorf("Error on get of ibm pi capture (%s) while applying pi_user_tags: %s", capturename, err)
		}
		if imagedata.Crn != "" {
			oldList, newList := d.GetChange(Arg_UserTags)
			err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, string(imagedata.Crn), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi capture (%s) pi_user_tags during creation: %s", *imagedata.ImageID, err)
			}
		}
	}

	return resourceIBMPICaptureRead(ctx, d, meta)
}

func resourceIBMPICaptureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := parts[0]
	captureID := parts[1]
	capturedestination := parts[2]
	if capturedestination != cloudStorageDestination {
		imageClient := st.NewIBMPIImageClient(ctx, sess, cloudInstanceID)
		imagedata, err := imageClient.Get(captureID)
		if err != nil {
			uErr := errors.Unwrap(err)
			switch uErr.(type) {
			case *p_cloud_images.PcloudCloudinstancesImagesGetNotFound:
				log.Printf("[DEBUG] image does not exist %v", err)
				d.SetId("")
				return diag.Errorf("image does not exist %v", err)
			}
			log.Printf("[DEBUG] get image failed %v", err)
			return diag.FromErr(err)
		}
		imageid := *imagedata.ImageID
		if imagedata.Crn != "" {
			d.Set(Attr_CRN, imagedata.Crn)
			tags, err := flex.GetGlobalTagsUsingCRN(meta, string(imagedata.Crn), "", UserTagType)
			if err != nil {
				log.Printf("Error on get of ibm pi capture (%s) pi_user_tags: %s", *imagedata.ImageID, err)
			}
			d.Set(Arg_UserTags, tags)
		}
		d.Set("image_id", imageid)
	}
	d.Set(helpers.PICloudInstanceId, cloudInstanceID)
	return nil
}

func resourceIBMPICaptureDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := parts[0]
	captureID := parts[1]
	capturedestination := parts[2]
	if capturedestination != cloudStorageDestination {
		imageClient := st.NewIBMPIImageClient(ctx, sess, cloudInstanceID)
		err = imageClient.Delete(captureID)
		if err != nil {
			uErr := errors.Unwrap(err)
			switch uErr.(type) {
			case *p_cloud_images.PcloudCloudinstancesImagesGetNotFound:
				log.Printf("[DEBUG] image does not exist while deleting %v", err)
				d.SetId("")
				return nil
			}
			log.Printf("[DEBUG] delete image failed %v", err)
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return nil
}

func resourceIBMPICaptureUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	captureID := parts[1]
	capturedestination := parts[2]

	if capturedestination != cloudStorageDestination && d.HasChange(Arg_UserTags) {
		if crn, ok := d.GetOk(Attr_CRN); ok {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, crn.(string), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi capture (%s) pi_user_tags: %s", captureID, err)
			}
		}
	}

	return resourceIBMPICaptureRead(ctx, d, meta)
}

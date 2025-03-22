// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMPIImageExport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIImageExportCreate,
		ReadContext:   resourceIBMPIImageExportRead,
		DeleteContext: resourceIBMPIImageExportDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_ImageAccessKey: {
				Description:  "Cloud Object Storage access key; required for buckets with private access",
				ForceNew:     true,
				Required:     true,
				Sensitive:    true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_ImageBucketName: {
				Description:  "Cloud Object Storage bucket name; bucket-name[/optional/folder]",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_ImageBucketRegion: {
				Description:  "Cloud Object Storage region",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_ImageID: {
				Description:      "Instance image id",
				DiffSuppressFunc: flex.ApplyOnce,
				ForceNew:         true,
				Required:         true,
				Type:             schema.TypeString,
				ValidateFunc:     validation.NoZeroValues,
			},
			Arg_ImageSecretKey: {
				Description:  "Cloud Object Storage secret key; required for buckets with private access",
				ForceNew:     true,
				Required:     true,
				Sensitive:    true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceIBMPIImageExportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		log.Printf("Failed to get the session")
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	imageid := d.Get(Arg_ImageID).(string)
	bucketName := d.Get(Arg_ImageBucketName).(string)
	accessKey := d.Get(Arg_ImageAccessKey).(string)

	client := instance.NewIBMPIImageClient(ctx, sess, cloudInstanceID)

	// image export
	var body = &models.ExportImage{
		BucketName: &bucketName,
		AccessKey:  &accessKey,
		Region:     d.Get(Arg_ImageBucketRegion).(string),
		SecretKey:  d.Get(Arg_ImageSecretKey).(string),
	}

	imageResponse, err := client.ExportImage(imageid, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", imageid, bucketName, d.Get(Arg_ImageBucketRegion).(string)))

	jobClient := instance.NewIBMPIJobClient(ctx, sess, cloudInstanceID)
	_, err = waitForIBMPIJobCompleted(ctx, jobClient, *imageResponse.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceIBMPIImageExportRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceIBMPIImageExportDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

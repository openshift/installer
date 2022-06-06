// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
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
			//required attributes
			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI cloud instance ID",
				ForceNew:    true,
			},
			helpers.PIImageId: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Instance image id",
				DiffSuppressFunc: flex.ApplyOnce,
				ForceNew:         true,
			},
			helpers.PIImageBucketName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud Object Storage bucket name; bucket-name[/optional/folder]",
				ForceNew:    true,
			},
			helpers.PIImageAccessKey: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud Object Storage access key; required for buckets with private access",
				Sensitive:   true,
				ForceNew:    true,
			},

			helpers.PIImageSecretKey: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud Object Storage secret key; required for buckets with private access",
				Sensitive:   true,
				ForceNew:    true,
			},
			helpers.PIImageBucketRegion: {
				Type:        schema.TypeString,
				Description: "Cloud Object Storage region",
				ForceNew:    true,
				Required:    true,
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

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	imageid := d.Get(helpers.PIImageId).(string)
	bucketName := d.Get(helpers.PIImageBucketName).(string)
	accessKey := d.Get(helpers.PIImageAccessKey).(string)

	client := st.NewIBMPIImageClient(ctx, sess, cloudInstanceID)

	// image export
	var body = &models.ExportImage{
		BucketName: &bucketName,
		AccessKey:  &accessKey,
		Region:     d.Get(helpers.PIImageBucketRegion).(string),
		SecretKey:  d.Get(helpers.PIImageSecretKey).(string),
	}

	imageResponse, err := client.ExportImage(imageid, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", imageid, bucketName, d.Get(helpers.PIImageBucketRegion).(string)))

	jobClient := st.NewIBMPIJobClient(ctx, sess, cloudInstanceID)
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

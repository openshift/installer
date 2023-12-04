// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cos

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMCosBucketObject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCosBucketObjectRead,

		Schema: map[string]*schema.Schema{
			"body": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "COS object body",
			},
			"bucket_crn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "COS bucket CRN",
			},
			"bucket_location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "COS bucket location",
			},
			"content_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "COS object content length",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "COS object content type",
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private", "direct"}),
				Description:  "COS endpoint type: public, private, direct",
				Default:      "public",
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "COS object MD5 hexdigest",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "COS object key",
			},
			"last_modified": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "COS object last modified date",
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_sql_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_lock_legal_hold_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_lock_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_lock_retain_until_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"website_redirect": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMCosBucketObjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	bucketCRN := d.Get("bucket_crn").(string)
	bucketName := strings.Split(bucketCRN, ":bucket:")[1]
	instanceCRN := fmt.Sprintf("%s::", strings.Split(bucketCRN, ":bucket:")[0])

	bucketLocation := d.Get("bucket_location").(string)
	endpointType := d.Get("endpoint_type").(string)

	bxSession, err := m.(conns.ClientSession).BluemixSession()
	if err != nil {
		return diag.FromErr(err)
	}

	s3Client, err := getS3Client(bxSession, bucketLocation, endpointType, instanceCRN)
	if err != nil {
		return diag.FromErr(err)
	}

	objectKey := d.Get("key").(string)
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	out, err := s3Client.HeadObject(headInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed getting COS bucket (%s) object (%s): %w", bucketName, objectKey, err))
	}

	log.Printf("[DEBUG] Received COS object: %s", out)

	d.Set("content_length", out.ContentLength)
	d.Set("content_type", out.ContentType)
	d.Set("etag", strings.Trim(aws.StringValue(out.ETag), `"`))
	if out.LastModified != nil {
		d.Set("last_modified", out.LastModified.Format(time.RFC1123))
	} else {
		d.Set("last_modified", "")
	}

	if isContentTypeAllowed(out.ContentType) {
		getInput := s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		}
		out, err := s3Client.GetObject(&getInput)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed getting COS object: %w", err))
		}

		buf := new(bytes.Buffer)
		bytesRead, err := buf.ReadFrom(out.Body)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed reading content of COS bucket (%s) object (%s): %w", bucketName, objectKey, err))
		}
		log.Printf("[INFO] Saving %d bytes from COS bucket (%s) object (%s)", bytesRead, bucketName, objectKey)
		d.Set("body", buf.String())
	} else {
		contentType := ""
		if out.ContentType == nil {
			contentType = "<EMPTY>"
		} else {
			contentType = aws.StringValue(out.ContentType)
		}

		log.Printf("[INFO] Ignoring body of COS bucket (%s) object (%s) with Content-Type %q", bucketName, objectKey, contentType)
	}
	if out.ObjectLockMode != nil {
		d.Set("object_lock_mode", out.ObjectLockMode)
	}
	if out.ObjectLockRetainUntilDate != nil {
		d.Set("object_lock_retain_until_date", out.ObjectLockRetainUntilDate.Format(time.RFC1123))
	}
	if out.ObjectLockLegalHoldStatus != nil {
		d.Set("object_lock_legal_hold_status", out.ObjectLockLegalHoldStatus)
	}
	if out.WebsiteRedirectLocation != nil {
		d.Set("website_redirect", out.WebsiteRedirectLocation)
	}

	objectID := getObjectId(bucketCRN, objectKey, bucketLocation)
	d.SetId(objectID)
	d.Set("version_id", out.VersionId)
	d.Set("object_sql_url", "cos://"+bucketLocation+"/"+bucketName+"/"+objectKey)
	return nil
}

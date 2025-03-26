// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cos

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	bxsession "github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/awserr"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	token "github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam/token"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validation "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMCOSBucketObject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCOSBucketObjectCreate,
		ReadContext:   resourceIBMCOSBucketObjectRead,
		UpdateContext: resourceIBMCOSBucketObjectUpdate,
		DeleteContext: resourceIBMCOSBucketObjectDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"body": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "COS object body",
			},
			"bucket_crn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "COS bucket CRN",
			},
			"bucket_location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "COS bucket location",
			},
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content_base64", "content_file"},
				Description:   "COS object content",
			},
			"content_base64": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content", "content_file"},
				Description:   "COS object content in base64 encoding",
			},
			"content_file": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content", "content_base64"},
				Description:   "COS object content file path",
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
				Optional:    true,
				Description: "COS object MD5 hexdigest",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "COS buckets need to be empty before they can be deleted. force_delete option empty the bucket and delete it.",
			},
			"object_sql_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access the object using an SQL Query instance.The reference url is used to perform queries against objects storing structured data.",
			},
			"object_lock_legal_hold_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(s3.ObjectLockLegalHoldStatus_Values(), false),
				Description:  "An object lock configuration on the object, the valid states are ON/OFF. When ON prevents deletion of the object version.",
			},
			"object_lock_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"object_lock_retain_until_date"},
				Description:  "Retention modes apply different levels of protection to the objects.",
			},
			"object_lock_retain_until_date": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"object_lock_mode"},
				ValidateFunc: validation.IsRFC3339Time,
				Description:  "An object cannot be deleted when the current time is earlier than the retainUntilDate. After this date, the object can be deleted.",
			},
			"website_redirect": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Redirect a request to another object or an URL",
			},
		},
	}
}

func resourceIBMCOSBucketObjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	var body io.ReadSeeker

	if v, ok := d.GetOk("content"); ok {
		content := v.(string)
		body = bytes.NewReader([]byte(content))
	} else if v, ok := d.GetOk("content_base64"); ok {
		content := v.(string)
		contentRaw, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error decoding content_base64: %s", err))
		}
		body = bytes.NewReader(contentRaw)
	} else if v, ok := d.GetOk("content_file"); ok {
		path := v.(string)
		file, err := os.Open(path)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error opening COS object file (%s): %s", path, err))
		}

		body = file
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Failed closing COS object file (%s): %s", path, err)
			}
		}()
	}

	putInput := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   body,
	}
	//if website redirect location if given for a an object
	if v, ok := d.GetOk("website_redirect"); ok {
		putInput.WebsiteRedirectLocation = aws.String(v.(string))
	}

	if _, err := s3Client.PutObject(putInput); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error putting object (%s) in COS bucket (%s): %s", objectKey, bucketName, err))
	}
	if v, ok := d.GetOk("object_lock_mode"); ok {
		if d, ok := d.GetOk("object_lock_retain_until_date"); ok {
			retainUntildate := parseDate(d.(string))
			putObjectRetentionInput := &s3.PutObjectRetentionInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
				Retention: &s3.ObjectLockRetention{
					Mode:            aws.String(v.(string)),
					RetainUntilDate: aws.Time(*retainUntildate),
				},
			}
			_, err = s3Client.PutObjectRetention(putObjectRetentionInput)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error putting objectlock retention on (%s) in COS bucket (%s): %s", objectKey, bucketName, err))
			}

		}
	}

	if v, ok := d.GetOk("object_lock_legal_hold_status"); ok {
		putObjectLegalHoldInput := &s3.PutObjectLegalHoldInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			LegalHold: &s3.ObjectLockLegalHold{
				Status: aws.String(v.(string)),
			},
		}
		_, err = s3Client.PutObjectLegalHold(putObjectLegalHoldInput)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error putting objectlock legalhold on (%s) in COS bucket (%s): %s", objectKey, bucketName, err))
		}
	}
	objectID := getObjectId(bucketCRN, objectKey, bucketLocation)
	d.SetId(objectID)

	return resourceIBMCOSBucketObjectRead(ctx, d, m)
}

func resourceIBMCOSBucketObjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	objectID := d.Id()

	bucketCRN := parseObjectId(objectID, "bucketCRN")
	bucketName := parseObjectId(objectID, "bucketName")
	bucketLocation := parseObjectId(objectID, "bucketLocation")
	instanceCRN := parseObjectId(objectID, "instanceCRN")
	endpointType := d.Get("endpoint_type").(string)

	d.Set("bucket_crn", bucketCRN)
	d.Set("bucket_location", bucketLocation)

	bxSession, err := m.(conns.ClientSession).BluemixSession()
	if err != nil {
		return diag.FromErr(err)
	}

	s3Client, err := getS3Client(bxSession, bucketLocation, endpointType, instanceCRN)
	if err != nil {
		return diag.FromErr(err)
	}

	objectKey := parseObjectId(objectID, "objectKey")
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	out, err := s3Client.HeadObject(headInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			d.SetId("") // Set state back to empty for terraform refresh
		}
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
	d.Set("key", objectKey)
	d.Set("version_id", out.VersionId)
	d.Set("object_sql_url", "cos://"+bucketLocation+"/"+bucketName+"/"+objectKey)
	return nil
}

func resourceIBMCOSBucketObjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	bucketCRN := d.Get("bucket_crn").(string)
	bucketName := strings.Split(bucketCRN, ":bucket:")[1]
	instanceCRN := fmt.Sprintf("%s::", strings.Split(bucketCRN, ":bucket:")[0])
	objectKey := d.Get("key").(string)
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
	if d.HasChanges("content", "content_base64", "content_file", "etag") {

		var body io.ReadSeeker

		if v, ok := d.GetOk("content"); ok {
			content := v.(string)
			body = bytes.NewReader([]byte(content))
		} else if v, ok := d.GetOk("content_base64"); ok {
			content := v.(string)
			contentRaw, err := base64.StdEncoding.DecodeString(content)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error decoding content_base64: %s", err))
			}
			body = bytes.NewReader(contentRaw)
		} else if v, ok := d.GetOk("content_file"); ok {
			path := v.(string)
			file, err := os.Open(path)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error opening COS object file (%s): %s", path, err))
			}

			body = file
			defer func() {
				err := file.Close()
				if err != nil {
					log.Printf("[WARN] Failed closing COS object file (%s): %s", path, err)
				}
			}()
		}

		objectKey := d.Get("key").(string)

		putInput := &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   body,
		}
		if d.HasChange("website_redirect") {
			if v, ok := d.GetOk("website_redirect"); ok {
				putInput.WebsiteRedirectLocation = aws.String(v.(string))
			}
		}

		if _, err := s3Client.PutObject(putInput); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error putting object (%s) in COS bucket (%s): %s", objectKey, bucketName, err))
		}

	}
	if d.HasChange("object_lock_legal_hold_status") {
		putObjectLegalHoldInput := &s3.PutObjectLegalHoldInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			LegalHold: &s3.ObjectLockLegalHold{
				Status: aws.String((d.Get("object_lock_legal_hold_status").(string))),
			},
		}
		_, err = s3Client.PutObjectLegalHold(putObjectLegalHoldInput)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error putting objectlock on legalhold (%s) in COS bucket (%s): %s", objectKey, bucketName, err))
		}
	}

	if d.HasChanges("object_lock_mode", "object_lock_retain_until_date") {
		putObjectRetentionInput := &s3.PutObjectRetentionInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Retention: &s3.ObjectLockRetention{
				Mode:            aws.String(d.Get("object_lock_mode").(string)),
				RetainUntilDate: parseDate(d.Get("object_lock_retain_until_date").(string)),
			},
		}
		_, err = s3Client.PutObjectRetention(putObjectRetentionInput)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error putting objectlock retention on (%s) in COS bucket (%s): %s", objectKey, bucketName, err))
		}
	}

	objectID := getObjectId(bucketCRN, objectKey, bucketLocation)
	d.SetId(objectID)

	return resourceIBMCOSBucketObjectRead(ctx, d, m)
}

func resourceIBMCOSBucketObjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	if _, ok := d.GetOk("version_id"); ok {
		err = deleteAllCOSObjectVersions(s3Client, bucketName, objectKey, d.Get("force_delete").(bool), false)
	} else {
		err = deleteCOSObjectVersion(s3Client, bucketName, objectKey, "", false)
	}

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getCosEndpoint(bucketLocation string, endpointType string) string {
	if bucketLocation != "" {
		hostUrl := "cloud-object-storage.appdomain.cloud"
		switch endpointType {
		case "public":
			return fmt.Sprintf("s3.%s.%s", bucketLocation, hostUrl)
		case "private":
			return fmt.Sprintf("s3.private.%s.%s", bucketLocation, hostUrl)
		case "direct":
			return fmt.Sprintf("s3.direct.%s.%s", bucketLocation, hostUrl)
		default:
			return fmt.Sprintf("s3.%s.%s", bucketLocation, hostUrl)
		}
	}
	return ""
}

func getS3Client(bxSession *bxsession.Session, bucketLocation string, endpointType string, instanceCRN string) (*s3.S3, error) {
	var s3Conf *aws.Config
	visibility := endpointType
	if endpointType == "direct" {
		visibility = "private"
	}
	apiEndpoint := getCosEndpoint(bucketLocation, endpointType)
	apiEndpoint = conns.FileFallBack(bxSession.Config.EndpointsFile, visibility, "IBMCLOUD_COS_ENDPOINT", bucketLocation, apiEndpoint)
	apiEndpoint = conns.EnvFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)
	if apiEndpoint == "" {
		return nil, fmt.Errorf("the endpoint doesn't exists for given location %s and endpoint type %s", bucketLocation, endpointType)
	}

	authEndpoint, err := bxSession.Config.EndpointLocator.IAMEndpoint()
	if err != nil {
		return nil, err
	}
	authEndpointPath := fmt.Sprintf("%s%s", authEndpoint, "/identity/token")
	apiKey := bxSession.Config.BluemixAPIKey
	if apiKey != "" {
		s3Conf = aws.NewConfig().WithEndpoint(apiEndpoint).WithCredentials(ibmiam.NewStaticCredentials(aws.NewConfig(), authEndpointPath, apiKey, instanceCRN)).WithS3ForcePathStyle(true)
	}
	iamAccessToken := bxSession.Config.IAMAccessToken
	if iamAccessToken != "" {
		initFunc := func() (*token.Token, error) {
			return &token.Token{
				AccessToken:  bxSession.Config.IAMAccessToken,
				RefreshToken: bxSession.Config.IAMRefreshToken,
				TokenType:    "Bearer",
				ExpiresIn:    int64((time.Hour * 248).Seconds()) * -1,
				Expiration:   time.Now().Add(-1 * time.Hour).Unix(),
			}, nil
		}
		s3Conf = aws.NewConfig().WithEndpoint(apiEndpoint).WithCredentials(ibmiam.NewCustomInitFuncCredentials(aws.NewConfig(), initFunc, authEndpointPath, instanceCRN)).WithS3ForcePathStyle(true)
	}
	s3Sess := session.Must(session.NewSession())
	return s3.New(s3Sess, s3Conf), nil
}

// This is to prevent potential issues w/ binary files
// and generally unprintable characters
// See https://github.com/hashicorp/terraform/pull/3858#issuecomment-156856738
func isContentTypeAllowed(contentType *string) bool {
	if contentType == nil {
		return false
	}

	allowedContentTypes := []*regexp.Regexp{
		regexp.MustCompile("^text/.+"),
		regexp.MustCompile("^application/json$"),
	}

	for _, r := range allowedContentTypes {
		if r.MatchString(*contentType) {
			return true
		}
	}

	return false
}

func objectExists(s3Client *s3.S3, bucketName string, objectKey string) (bool, error) {
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}
	_, err := s3Client.HeadObject(headInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func getObjectId(bucketCRN string, objectKey string, bucketLocation string) string {
	return fmt.Sprintf("%s:object:%s:location:%s", bucketCRN, objectKey, bucketLocation)
}

func parseObjectId(id string, info string) string {
	splitID := strings.Split(id, ":object:")
	bucketCRN := splitID[0]

	if info == "instanceCRN" {
		return fmt.Sprintf("%s::", strings.Split(bucketCRN, ":bucket:")[0])
	}

	if info == "bucketCRN" {
		return bucketCRN
	}

	if info == "bucketName" {
		return strings.Split(bucketCRN, ":bucket:")[1]
	}

	if info == "objectKey" {
		return strings.Split(splitID[1], ":location:")[0]
	}

	if info == "bucketLocation" {
		return strings.Split(splitID[1], ":location:")[1]
	}

	return parseBucketId(splitID[0], info)
}

// deleteAllCOSObjectVersions deletes all versions of a specified key from an COS bucket.
// If key is empty then all versions of all objects are deleted.
func deleteAllCOSObjectVersions(conn *s3.S3, bucketName, key string, force, ignoreObjectErrors bool) error {
	input := &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucketName),
	}
	if key != "" {
		input.Prefix = aws.String(key)
	}

	var lastErr error
	err := conn.ListObjectVersionsPages(input, func(page *s3.ListObjectVersionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, objectVersion := range page.Versions {
			objectKey := aws.StringValue(objectVersion.Key)
			objectVersionID := aws.StringValue(objectVersion.VersionId)
			log.Printf("[INFO] Object {%s} Version Id {%s}: ", objectKey, objectVersionID)

			if key != "" && key != objectKey {
				continue
			}

			err := deleteCOSObjectVersion(conn, bucketName, objectKey, objectVersionID, force)

			if err != nil {
				if strings.Contains(err.Error(), "AccessDenied") && force {
					// Remove any legal hold.
					_, err := conn.HeadObject(&s3.HeadObjectInput{
						Bucket:    aws.String(bucketName),
						Key:       objectVersion.Key,
						VersionId: objectVersion.VersionId,
					})

					if err != nil {
						log.Printf("[ERROR] Error getting COS Bucket (%s) Object (%s) Version (%s) metadata: %s", bucketName, objectKey, objectVersionID, err)
						lastErr = err
						continue
					}

					// AccessDenied for another reason.
					lastErr = fmt.Errorf("AccessDenied deleting COS Bucket (%s) Object (%s) Version: %s", bucketName, objectKey, objectVersionID)
					continue
				}
			}
		}

		return !lastPage
	})

	if err != nil {
		return err
	}

	if lastErr != nil {
		if !ignoreObjectErrors {
			return fmt.Errorf("[ERROR] Error deleting at least one object version, last error: %s", lastErr)
		}

		lastErr = nil
	}

	err = conn.ListObjectVersionsPages(input, func(page *s3.ListObjectVersionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, deleteMarker := range page.DeleteMarkers {
			deleteMarkerKey := aws.StringValue(deleteMarker.Key)
			deleteMarkerVersionID := aws.StringValue(deleteMarker.VersionId)

			if key != "" && key != deleteMarkerKey {
				continue
			}

			// Delete markers have no object lock protections.
			err := deleteCOSObjectVersion(conn, bucketName, deleteMarkerKey, deleteMarkerVersionID, false)

			if err != nil {
				lastErr = err
			}
		}

		return !lastPage
	})

	if err != nil {
		return err
	}

	if lastErr != nil {
		if !ignoreObjectErrors {
			return fmt.Errorf("[ERROR] Error deleting at least one object delete marker, last error: %s", lastErr)
		}

		lastErr = nil
	}

	return nil
}

// deleteCOSObjectVersion deletes a specific bucket object version.
func deleteCOSObjectVersion(conn *s3.S3, b, k, v string, force bool) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(b),
		Key:    aws.String(k),
	}

	if v != "" {
		input.VersionId = aws.String(v)
	}

	log.Printf("[INFO] Deleting COS Bucket (%s) Object (%s) Version: %s", b, k, v)
	_, err := conn.DeleteObject(input)

	if err != nil {
		log.Printf("[WARN] Error deleting S3 Bucket (%s) Object (%s) Version (%s): %s", b, k, v, err)
	}

	return err
}
func parseDate(v string) *time.Time {
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil
	}

	return aws.Time(t)
}

func objectDatetoString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

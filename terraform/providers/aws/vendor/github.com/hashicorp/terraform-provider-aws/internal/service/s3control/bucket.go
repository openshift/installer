package s3control

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

const (
	// Maximum amount of time to wait for s3control Bucket state to propagate
	bucketStatePropagationTimeout = 5 * time.Minute
)

func ResourceBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceBucketCreate,
		Read:   resourceBucketRead,
		Update: resourceBucketUpdate,
		Delete: resourceBucketDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 63),
					validation.StringMatch(regexp.MustCompile(`^[a-z0-9.-]+$`), "must contain only lowercase letters, numbers, periods, and hyphens"),
					validation.StringMatch(regexp.MustCompile(`^[a-z0-9]`), "must begin with lowercase letter or number"),
					validation.StringMatch(regexp.MustCompile(`[a-z0-9]$`), "must end with lowercase letter or number"),
				),
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"outpost_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"public_access_block_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceBucketCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	bucket := d.Get("bucket").(string)

	input := &s3control.CreateBucketInput{
		Bucket:    aws.String(bucket),
		OutpostId: aws.String(d.Get("outpost_id").(string)),
	}

	output, err := conn.CreateBucket(input)

	if err != nil {
		return fmt.Errorf("error creating S3 Control Bucket (%s): %w", bucket, err)
	}

	if output == nil {
		return fmt.Errorf("error creating S3 Control Bucket (%s): empty response", bucket)
	}

	d.SetId(aws.StringValue(output.BucketArn))

	if len(tags) > 0 {
		if err := bucketUpdateTags(conn, d.Id(), nil, tags); err != nil {
			return fmt.Errorf("error adding S3 Control Bucket (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceBucketRead(d, meta)
}

func resourceBucketRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	parsedArn, err := arn.Parse(d.Id())

	if err != nil {
		return fmt.Errorf("error parsing S3 Control Bucket ARN (%s): %w", d.Id(), err)
	}

	// ARN resource format: outpost/<outpost-id>/bucket/<my-bucket-name>
	arnResourceParts := strings.Split(parsedArn.Resource, "/")

	if parsedArn.AccountID == "" || len(arnResourceParts) != 4 {
		return fmt.Errorf("error parsing S3 Control Bucket ARN (%s): unknown format", d.Id())
	}

	input := &s3control.GetBucketInput{
		AccountId: aws.String(parsedArn.AccountID),
		Bucket:    aws.String(d.Id()),
	}

	output, err := conn.GetBucket(input)

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, "NoSuchBucket") {
		log.Printf("[WARN] S3 Control Bucket (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, "NoSuchOutpost") {
		log.Printf("[WARN] S3 Control Bucket (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading S3 Control Bucket (%s): %w", d.Id(), err)
	}

	if output == nil {
		return fmt.Errorf("error reading S3 Control Bucket (%s): empty response", d.Id())
	}

	d.Set("arn", d.Id())
	d.Set("bucket", output.Bucket)

	if output.CreationDate != nil {
		d.Set("creation_date", aws.TimeValue(output.CreationDate).Format(time.RFC3339))
	}

	d.Set("outpost_id", arnResourceParts[1])
	d.Set("public_access_block_enabled", output.PublicAccessBlockEnabled)

	tags, err := bucketListTags(conn, d.Id())

	if err != nil {
		return fmt.Errorf("error listing tags for S3 Control Bucket (%s): %w", d.Id(), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := bucketUpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating S3 Control Bucket (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceBucketRead(d, meta)
}

func resourceBucketDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).S3ControlConn

	parsedArn, err := arn.Parse(d.Id())

	if err != nil {
		return fmt.Errorf("error parsing S3 Control Bucket ARN (%s): %w", d.Id(), err)
	}

	input := &s3control.DeleteBucketInput{
		AccountId: aws.String(parsedArn.AccountID),
		Bucket:    aws.String(d.Id()),
	}

	// S3 Control Bucket have a backend state which cannot be checked so this error
	// can occur on deletion:
	//   InvalidBucketState: Bucket is in an invalid state
	err = resource.Retry(bucketStatePropagationTimeout, func() *resource.RetryError {
		_, err := conn.DeleteBucket(input)

		if tfawserr.ErrCodeEquals(err, "InvalidBucketState") {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.DeleteBucket(input)
	}

	if tfawserr.ErrCodeEquals(err, "NoSuchBucket") {
		return nil
	}

	if tfawserr.ErrCodeEquals(err, "NoSuchOutpost") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting S3 Control Bucket (%s): %w", d.Id(), err)
	}

	return nil
}

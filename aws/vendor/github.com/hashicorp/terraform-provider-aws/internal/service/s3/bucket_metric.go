package s3

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_s3_bucket_metric")
func ResourceBucketMetric() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceBucketMetricPut,
		ReadWithoutTimeout:   resourceBucketMetricRead,
		UpdateWithoutTimeout: resourceBucketMetricPut,
		DeleteWithoutTimeout: resourceBucketMetricDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: filterAtLeastOneOfKeys,
						},
						"tags": {
							Type:         schema.TypeMap,
							Optional:     true,
							Elem:         &schema.Schema{Type: schema.TypeString},
							AtLeastOneOf: filterAtLeastOneOfKeys,
						},
					},
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
		},
	}
}

func resourceBucketMetricPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).S3Conn(ctx)
	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)

	metricsConfiguration := &s3.MetricsConfiguration{
		Id: aws.String(name),
	}

	if v, ok := d.GetOk("filter"); ok {
		filterList := v.([]interface{})
		if filterMap, ok := filterList[0].(map[string]interface{}); ok {
			metricsConfiguration.Filter = ExpandMetricsFilter(ctx, filterMap)
		}
	}

	input := &s3.PutBucketMetricsConfigurationInput{
		Bucket:               aws.String(bucket),
		Id:                   aws.String(name),
		MetricsConfiguration: metricsConfiguration,
	}

	log.Printf("[DEBUG] Putting S3 Bucket Metrics Configuration: %s", input)
	err := retry.RetryContext(ctx, propagationTimeout, func() *retry.RetryError {
		_, err := conn.PutBucketMetricsConfigurationWithContext(ctx, input)

		if tfawserr.ErrCodeEquals(err, s3.ErrCodeNoSuchBucket) {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.PutBucketMetricsConfigurationWithContext(ctx, input)
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "putting S3 Bucket Metrics Configuration: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", bucket, name))

	return append(diags, resourceBucketMetricRead(ctx, d, meta)...)
}

func resourceBucketMetricDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).S3Conn(ctx)

	bucket, name, err := BucketMetricParseID(d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting S3 Bucket Metrics Configuration (%s): %s", d.Id(), err)
	}

	input := &s3.DeleteBucketMetricsConfigurationInput{
		Bucket: aws.String(bucket),
		Id:     aws.String(name),
	}

	log.Printf("[DEBUG] Deleting S3 Bucket Metrics Configuration: %s", input)
	_, err = conn.DeleteBucketMetricsConfigurationWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, s3.ErrCodeNoSuchBucket) {
		return diags
	}

	if tfawserr.ErrCodeEquals(err, errCodeNoSuchConfiguration) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting S3 Bucket Metrics Configuration (%s): %s", d.Id(), err)
	}

	return diags
}

func resourceBucketMetricRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).S3Conn(ctx)

	bucket, name, err := BucketMetricParseID(d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading S3 Bucket Metrics Configuration (%s): %s", d.Id(), err)
	}

	d.Set("bucket", bucket)
	d.Set("name", name)

	input := &s3.GetBucketMetricsConfigurationInput{
		Bucket: aws.String(bucket),
		Id:     aws.String(name),
	}

	log.Printf("[DEBUG] Reading S3 Bucket Metrics Configuration: %s", input)
	output, err := conn.GetBucketMetricsConfigurationWithContext(ctx, input)

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, s3.ErrCodeNoSuchBucket) {
		log.Printf("[WARN] S3 Bucket Metrics Configuration (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, errCodeNoSuchConfiguration) {
		log.Printf("[WARN] S3 Bucket Metrics Configuration (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading S3 Bucket Metrics Configuration (%s): %s", d.Id(), err)
	}

	if output == nil || output.MetricsConfiguration == nil {
		return sdkdiag.AppendErrorf(diags, "reading S3 Bucket Metrics Configuration (%s): empty response", d.Id())
	}

	if output.MetricsConfiguration.Filter != nil {
		if err := d.Set("filter", []interface{}{FlattenMetricsFilter(ctx, output.MetricsConfiguration.Filter)}); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting filter")
		}
	}

	return diags
}

func ExpandMetricsFilter(ctx context.Context, m map[string]interface{}) *s3.MetricsFilter {
	var prefix string
	if v, ok := m["prefix"]; ok {
		prefix = v.(string)
	}

	var tags []*s3.Tag
	if v, ok := m["tags"]; ok {
		tags = Tags(tftags.New(ctx, v).IgnoreAWS())
	}

	metricsFilter := &s3.MetricsFilter{}
	if prefix != "" && len(tags) > 0 {
		metricsFilter.And = &s3.MetricsAndOperator{
			Prefix: aws.String(prefix),
			Tags:   tags,
		}
	} else if len(tags) > 1 {
		metricsFilter.And = &s3.MetricsAndOperator{
			Tags: tags,
		}
	} else if len(tags) == 1 {
		metricsFilter.Tag = tags[0]
	} else {
		metricsFilter.Prefix = aws.String(prefix)
	}
	return metricsFilter
}

func FlattenMetricsFilter(ctx context.Context, metricsFilter *s3.MetricsFilter) map[string]interface{} {
	m := make(map[string]interface{})

	if and := metricsFilter.And; and != nil {
		if and.Prefix != nil {
			m["prefix"] = aws.StringValue(and.Prefix)
		}
		if and.Tags != nil {
			m["tags"] = KeyValueTags(ctx, and.Tags).IgnoreAWS().Map()
		}
	} else if metricsFilter.Prefix != nil {
		m["prefix"] = aws.StringValue(metricsFilter.Prefix)
	} else if metricsFilter.Tag != nil {
		tags := []*s3.Tag{
			metricsFilter.Tag,
		}
		m["tags"] = KeyValueTags(ctx, tags).IgnoreAWS().Map()
	}
	return m
}

func BucketMetricParseID(id string) (string, string, error) {
	idParts := strings.Split(id, ":")
	if len(idParts) != 2 {
		return "", "", fmt.Errorf("please make sure the ID is in the form BUCKET:NAME (i.e. my-bucket:EntireBucket")
	}
	bucket := idParts[0]
	name := idParts[1]
	return bucket, name, nil
}

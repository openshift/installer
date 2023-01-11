package macie2

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/macie2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceClassificationExportConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceClassificationExportConfigurationCreate,
		UpdateWithoutTimeout: resourceClassificationExportConfigurationUpdate,
		DeleteWithoutTimeout: resourceClassificationExportConfigurationDelete,
		ReadWithoutTimeout:   resourceClassificationExportConfigurationRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"s3_destination": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				AtLeastOneOf: []string{"s3_destination"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kms_key_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},
					},
				},
			},
		},
	}
}

func resourceClassificationExportConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	conn := meta.(*conns.AWSClient).Macie2Conn

	if d.IsNewResource() {
		output, err := conn.GetClassificationExportConfiguration(&macie2.GetClassificationExportConfigurationInput{})
		if err != nil {
			return diag.FromErr(fmt.Errorf("reading Macie classification export configuration failed: %w", err))
		}

		if (macie2.ClassificationExportConfiguration{}) != *output.Configuration { // nosemgrep: ci.prefer-aws-go-sdk-pointer-conversion-conditional
			return diag.FromErr(fmt.Errorf("creating Macie classification export configuration: a configuration already exists"))
		}
	}

	input := macie2.PutClassificationExportConfigurationInput{
		Configuration: &macie2.ClassificationExportConfiguration{},
	}

	if v, ok := d.GetOk("s3_destination"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.Configuration.S3Destination = expandClassificationExportConfiguration(v.([]interface{})[0].(map[string]interface{}))
	}

	log.Printf("[DEBUG] Creating Macie classification export configuration: %s", input)

	_, err := conn.PutClassificationExportConfiguration(&input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("creating Macie classification export configuration failed: %w", err))
	}

	return resourceClassificationExportConfigurationRead(ctx, d, meta)
}

func resourceClassificationExportConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).Macie2Conn

	input := macie2.PutClassificationExportConfigurationInput{
		Configuration: &macie2.ClassificationExportConfiguration{},
	}

	if v, ok := d.GetOk("s3_destination"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.Configuration.S3Destination = expandClassificationExportConfiguration(v.([]interface{})[0].(map[string]interface{}))
	} else {
		input.Configuration.S3Destination = nil
	}

	log.Printf("[DEBUG] Creating Macie classification export configuration: %s", input)

	_, err := conn.PutClassificationExportConfiguration(&input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("creating Macie classification export configuration failed: %w", err))
	}

	return resourceClassificationExportConfigurationRead(ctx, d, meta)
}

func resourceClassificationExportConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).Macie2Conn

	input := macie2.GetClassificationExportConfigurationInput{} // api does not have a getById() like endpoint.
	output, err := conn.GetClassificationExportConfiguration(&input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("reading Macie classification export configuration failed: %w", err))
	}

	if (macie2.ClassificationExportConfiguration{}) != *output.Configuration { // nosemgrep: ci.prefer-aws-go-sdk-pointer-conversion-conditional
		if (macie2.S3Destination{}) != *output.Configuration.S3Destination { // nosemgrep: ci.prefer-aws-go-sdk-pointer-conversion-conditional
			var flattenedS3Destination = flattenClassificationExportConfigurationS3DestinationResult(output.Configuration.S3Destination)
			if err := d.Set("s3_destination", []interface{}{flattenedS3Destination}); err != nil {
				return diag.FromErr(fmt.Errorf("error setting Macie classification export configuration s3_destination: %w", err))
			}
		}
		d.SetId(fmt.Sprintf("%s:%s:%s", "macie:classification_export_configuration", meta.(*conns.AWSClient).AccountID, meta.(*conns.AWSClient).Region))
	}

	return nil
}

func resourceClassificationExportConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).Macie2Conn

	input := macie2.PutClassificationExportConfigurationInput{
		Configuration: &macie2.ClassificationExportConfiguration{},
	}

	log.Printf("[DEBUG] deleting Macie classification export configuration: %s", input)

	_, err := conn.PutClassificationExportConfiguration(&input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("deleting Macie classification export configuration failed: %w", err))
	}

	return nil
}

func expandClassificationExportConfiguration(tfMap map[string]interface{}) *macie2.S3Destination {
	if tfMap == nil {
		return nil
	}

	apiObject := &macie2.S3Destination{}

	if v, ok := tfMap["bucket_name"].(string); ok {
		apiObject.BucketName = aws.String(v)
	}

	if v, ok := tfMap["key_prefix"].(string); ok {
		apiObject.KeyPrefix = aws.String(v)
	}

	if v, ok := tfMap["kms_key_arn"].(string); ok {
		apiObject.KmsKeyArn = aws.String(v)
	}

	return apiObject
}

func flattenClassificationExportConfigurationS3DestinationResult(apiObject *macie2.S3Destination) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.BucketName; v != nil {
		tfMap["bucket_name"] = aws.StringValue(v)
	}

	if v := apiObject.KeyPrefix; v != nil {
		tfMap["key_prefix"] = aws.StringValue(v)
	}

	if v := apiObject.KmsKeyArn; v != nil {
		tfMap["kms_key_arn"] = aws.StringValue(v)
	}

	return tfMap
}

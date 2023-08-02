package secretsmanager

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKDataSource("aws_secretsmanager_secret")
func DataSourceSecret() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceSecretRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: verify.ValidARN,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SecretsManagerConn(ctx)
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	var secretID string
	if v, ok := d.GetOk("arn"); ok {
		secretID = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		if secretID != "" {
			return sdkdiag.AppendErrorf(diags, "specify only arn or name")
		}
		secretID = v.(string)
	}

	if secretID == "" {
		return sdkdiag.AppendErrorf(diags, "must specify either arn or name")
	}

	input := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretID),
	}

	log.Printf("[DEBUG] Reading Secrets Manager Secret: %s", input)
	output, err := conn.DescribeSecretWithContext(ctx, input)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, secretsmanager.ErrCodeResourceNotFoundException) {
			return sdkdiag.AppendErrorf(diags, "Secrets Manager Secret %q not found", secretID)
		}
		return sdkdiag.AppendErrorf(diags, "reading Secrets Manager Secret: %s", err)
	}

	if output.ARN == nil {
		return sdkdiag.AppendErrorf(diags, "Secrets Manager Secret %q not found", secretID)
	}

	d.SetId(aws.StringValue(output.ARN))
	d.Set("arn", output.ARN)
	d.Set("description", output.Description)
	d.Set("kms_key_id", output.KmsKeyId)
	d.Set("name", output.Name)
	d.Set("policy", "")

	pIn := &secretsmanager.GetResourcePolicyInput{
		SecretId: aws.String(d.Id()),
	}
	log.Printf("[DEBUG] Reading Secrets Manager Secret policy: %s", pIn)
	pOut, err := conn.GetResourcePolicyWithContext(ctx, pIn)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Secrets Manager Secret policy: %s", err)
	}

	if pOut != nil && pOut.ResourcePolicy != nil {
		policy, err := structure.NormalizeJsonString(aws.StringValue(pOut.ResourcePolicy))
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "policy contains an invalid JSON: %s", err)
		}
		d.Set("policy", policy)
	}

	if err := d.Set("tags", KeyValueTags(ctx, output.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	return diags
}

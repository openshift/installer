package signer

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/signer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

// @SDKDataSource("aws_signer_signing_profile")
func DataSourceSigningProfile() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceSigningProfileRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform_display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"revocation_record": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"revocation_effective_from": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"revoked_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"revoked_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"signature_validity_period": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSigningProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SignerConn(ctx)
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	profileName := d.Get("name").(string)
	signingProfileOutput, err := conn.GetSigningProfileWithContext(ctx, &signer.GetSigningProfileInput{
		ProfileName: aws.String(profileName),
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Signer signing profile (%s): %s", d.Id(), err)
	}

	if err := d.Set("platform_id", signingProfileOutput.PlatformId); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing profile platform id: %s", err)
	}

	if err := d.Set("signature_validity_period", []interface{}{
		map[string]interface{}{
			"value": signingProfileOutput.SignatureValidityPeriod.Value,
			"type":  signingProfileOutput.SignatureValidityPeriod.Type,
		},
	}); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing profile signature validity period: %s", err)
	}

	if err := d.Set("platform_display_name", signingProfileOutput.PlatformDisplayName); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing profile platform display name: %s", err)
	}

	if err := d.Set("arn", signingProfileOutput.Arn); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing profile arn: %s", err)
	}

	if err := d.Set("version", signingProfileOutput.ProfileVersion); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing profile version: %s", err)
	}

	if err := d.Set("version_arn", signingProfileOutput.ProfileVersionArn); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing profile version arn: %s", err)
	}

	if err := d.Set("status", signingProfileOutput.Status); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing profile status: %s", err)
	}

	if err := d.Set("tags", KeyValueTags(ctx, signingProfileOutput.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing profile tags: %s", err)
	}

	if err := d.Set("revocation_record", flattenSigningProfileRevocationRecord(signingProfileOutput.RevocationRecord)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing profile revocation record: %s", err)
	}

	d.SetId(aws.StringValue(signingProfileOutput.ProfileName))

	return diags
}

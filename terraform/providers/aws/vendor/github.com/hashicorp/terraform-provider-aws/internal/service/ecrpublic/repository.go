package ecrpublic

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecrpublic"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_ecrpublic_repository", name="Repository")
// @Tags(identifierAttribute="arn")
func ResourceRepository() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceRepositoryCreate,
		ReadWithoutTimeout:   resourceRepositoryRead,
		UpdateWithoutTimeout: resourceRepositoryUpdate,
		DeleteWithoutTimeout: resourceRepositoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"repository_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(2, 205),
					validation.StringMatch(regexp.MustCompile(`(?:[a-z0-9]+(?:[._-][a-z0-9]+)*/)*[a-z0-9]+(?:[._-][a-z0-9]+)*`), "see: https://docs.aws.amazon.com/AmazonECRPublic/latest/APIReference/API_CreateRepository.html#API_CreateRepository_RequestSyntax"),
				),
			},
			"catalog_data": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"about_text": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 10240),
						},
						"architectures": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 50,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 1024),
						},
						"logo_image_blob": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"operating_systems": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 50,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"usage_text": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 10240),
						},
					},
				},
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
			},
			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"registry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repository_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},
	}
}

func resourceRepositoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ECRPublicConn(ctx)

	input := ecrpublic.CreateRepositoryInput{
		RepositoryName: aws.String(d.Get("repository_name").(string)),
		Tags:           GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("catalog_data"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.CatalogData = expandRepositoryCatalogData(v.([]interface{})[0].(map[string]interface{}))
	}

	out, err := conn.CreateRepositoryWithContext(ctx, &input)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating ECR Public repository: %s", err)
	}

	if out == nil {
		return sdkdiag.AppendErrorf(diags, "creating ECR Public Repository: empty response")
	}

	repository := out.Repository

	log.Printf("[DEBUG] ECR Public repository created: %q", aws.StringValue(repository.RepositoryArn))

	d.SetId(aws.StringValue(repository.RepositoryName))

	return append(diags, resourceRepositoryRead(ctx, d, meta)...)
}

func resourceRepositoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ECRPublicConn(ctx)

	log.Printf("[DEBUG] Reading ECR Public repository %s", d.Id())
	var out *ecrpublic.DescribeRepositoriesOutput
	input := &ecrpublic.DescribeRepositoriesInput{
		RepositoryNames: aws.StringSlice([]string{d.Id()}),
	}

	var err error
	err = retry.RetryContext(ctx, 1*time.Minute, func() *retry.RetryError {
		out, err = conn.DescribeRepositoriesWithContext(ctx, input)
		if d.IsNewResource() && tfawserr.ErrCodeEquals(err, ecrpublic.ErrCodeRepositoryNotFoundException) {
			return retry.RetryableError(err)
		}
		if err != nil {
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if tfresource.TimedOut(err) {
		out, err = conn.DescribeRepositoriesWithContext(ctx, input)
	}

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, ecrpublic.ErrCodeRepositoryNotFoundException) {
		log.Printf("[WARN] ECR Public Repository (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading ECR Public repository: %s", err)
	}

	if out == nil || len(out.Repositories) == 0 || out.Repositories[0] == nil {
		return sdkdiag.AppendErrorf(diags, "reading ECR Public Repository (%s): empty response", d.Id())
	}

	repository := out.Repositories[0]

	d.Set("repository_name", d.Id())
	d.Set("registry_id", repository.RegistryId)
	d.Set("arn", repository.RepositoryArn)
	d.Set("repository_uri", repository.RepositoryUri)

	if v, ok := d.GetOk("force_destroy"); ok {
		d.Set("force_destroy", v.(bool))
	} else {
		d.Set("force_destroy", false)
	}

	var catalogOut *ecrpublic.GetRepositoryCatalogDataOutput
	catalogInput := &ecrpublic.GetRepositoryCatalogDataInput{
		RepositoryName: aws.String(d.Id()),
		RegistryId:     repository.RegistryId,
	}

	catalogOut, err = conn.GetRepositoryCatalogDataWithContext(ctx, catalogInput)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading catalog data for ECR Public repository: %s", err)
	}

	if catalogOut != nil {
		flatCatalogData := flattenRepositoryCatalogData(catalogOut)
		if catalogData, ok := d.GetOk("catalog_data"); ok && len(catalogData.([]interface{})) > 0 && catalogData.([]interface{})[0] != nil {
			catalogDataMap := catalogData.([]interface{})[0].(map[string]interface{})
			if v, ok := catalogDataMap["logo_image_blob"].(string); ok && len(v) > 0 {
				flatCatalogData["logo_image_blob"] = v
			}
		}
		d.Set("catalog_data", []interface{}{flatCatalogData})
	} else {
		d.Set("catalog_data", nil)
	}

	return diags
}

func resourceRepositoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ECRPublicConn(ctx)

	deleteInput := &ecrpublic.DeleteRepositoryInput{
		RepositoryName: aws.String(d.Id()),
		RegistryId:     aws.String(d.Get("registry_id").(string)),
	}

	if v, ok := d.GetOk("force_destroy"); ok {
		deleteInput.Force = aws.Bool(v.(bool))
	}

	log.Printf("[DEBUG] Deleting ECR Public Repository: (%s)", d.Id())
	_, err := conn.DeleteRepositoryWithContext(ctx, deleteInput)

	if err != nil {
		if tfawserr.ErrCodeEquals(err, ecrpublic.ErrCodeRepositoryNotFoundException) {
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "deleting ECR Public repository: %s", err)
	}

	log.Printf("[DEBUG] Waiting for ECR Public Repository %q to be deleted", d.Id())
	input := &ecrpublic.DescribeRepositoriesInput{
		RepositoryNames: aws.StringSlice([]string{d.Id()}),
	}
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		_, err = conn.DescribeRepositoriesWithContext(ctx, input)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, ecrpublic.ErrCodeRepositoryNotFoundException) {
				return nil
			}
			return retry.NonRetryableError(err)
		}

		return retry.RetryableError(fmt.Errorf("%q: Timeout while waiting for the ECR Public Repository to be deleted", d.Id()))
	})
	if tfresource.TimedOut(err) {
		_, err = conn.DescribeRepositoriesWithContext(ctx, input)
	}

	if tfawserr.ErrCodeEquals(err, ecrpublic.ErrCodeRepositoryNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting ECR Public repository: %s", err)
	}

	log.Printf("[DEBUG] repository %q deleted.", d.Get("repository_name").(string))

	return diags
}

func resourceRepositoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ECRPublicConn(ctx)

	if d.HasChange("catalog_data") {
		if err := resourceRepositoryUpdateCatalogData(ctx, conn, d); err != nil {
			return sdkdiag.AppendErrorf(diags, "updating ECR Public Repository (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceRepositoryRead(ctx, d, meta)...)
}

func flattenRepositoryCatalogData(apiObject *ecrpublic.GetRepositoryCatalogDataOutput) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	catalogData := apiObject.CatalogData

	tfMap := map[string]interface{}{}

	if v := catalogData.AboutText; v != nil {
		tfMap["about_text"] = aws.StringValue(v)
	}

	if v := catalogData.Architectures; v != nil {
		tfMap["architectures"] = aws.StringValueSlice(v)
	}

	if v := catalogData.Description; v != nil {
		tfMap["description"] = aws.StringValue(v)
	}

	if v := catalogData.OperatingSystems; v != nil {
		tfMap["operating_systems"] = aws.StringValueSlice(v)
	}

	if v := catalogData.UsageText; v != nil {
		tfMap["usage_text"] = aws.StringValue(v)
	}

	return tfMap
}

func expandRepositoryCatalogData(tfMap map[string]interface{}) *ecrpublic.RepositoryCatalogDataInput {
	if tfMap == nil {
		return nil
	}

	repositoryCatalogDataInput := &ecrpublic.RepositoryCatalogDataInput{}

	if v, ok := tfMap["about_text"].(string); ok && v != "" {
		repositoryCatalogDataInput.AboutText = aws.String(v)
	}

	if v, ok := tfMap["architectures"].(*schema.Set); ok {
		repositoryCatalogDataInput.Architectures = flex.ExpandStringSet(v)
	}

	if v, ok := tfMap["description"].(string); ok && v != "" {
		repositoryCatalogDataInput.Description = aws.String(v)
	}

	if v, ok := tfMap["logo_image_blob"].(string); ok && len(v) > 0 {
		data, _ := base64.StdEncoding.DecodeString(v)
		repositoryCatalogDataInput.LogoImageBlob = data
	}

	if v, ok := tfMap["operating_systems"].(*schema.Set); ok {
		repositoryCatalogDataInput.OperatingSystems = flex.ExpandStringSet(v)
	}

	if v, ok := tfMap["usage_text"].(string); ok && v != "" {
		repositoryCatalogDataInput.UsageText = aws.String(v)
	}

	return repositoryCatalogDataInput
}

func resourceRepositoryUpdateCatalogData(ctx context.Context, conn *ecrpublic.ECRPublic, d *schema.ResourceData) error {
	if d.HasChange("catalog_data") {
		if v, ok := d.GetOk("catalog_data"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			input := ecrpublic.PutRepositoryCatalogDataInput{
				RepositoryName: aws.String(d.Id()),
				RegistryId:     aws.String(d.Get("registry_id").(string)),
				CatalogData:    expandRepositoryCatalogData(v.([]interface{})[0].(map[string]interface{})),
			}

			_, err := conn.PutRepositoryCatalogDataWithContext(ctx, &input)

			if err != nil {
				return fmt.Errorf("updating catalog data for repository(%s): %s", d.Id(), err)
			}
		}
	}

	return nil
}

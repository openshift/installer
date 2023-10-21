package servicecatalog

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_servicecatalog_product", name="Product")
// @Tags
func ResourceProduct() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceProductCreate,
		ReadWithoutTimeout:   resourceProductRead,
		UpdateWithoutTimeout: resourceProductUpdate,
		DeleteWithoutTimeout: resourceProductDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(ProductReadyTimeout),
			Read:   schema.DefaultTimeout(ProductReadTimeout),
			Update: schema.DefaultTimeout(ProductUpdateTimeout),
			Delete: schema.DefaultTimeout(ProductDeleteTimeout),
		},

		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      AcceptLanguageEnglish,
				ValidateFunc: validation.StringInSlice(AcceptLanguage_Values(), false),
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"distributor": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"has_default_path": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provisioning_artifact_parameters": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disable_template_validation": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"template_physical_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ExactlyOneOf: []string{
								"provisioning_artifact_parameters.0.template_url",
								"provisioning_artifact_parameters.0.template_physical_id",
							},
						},
						"template_url": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ExactlyOneOf: []string{
								"provisioning_artifact_parameters.0.template_url",
								"provisioning_artifact_parameters.0.template_physical_id",
							},
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(servicecatalog.ProvisioningArtifactType_Values(), false),
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"support_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"support_email": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"support_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(servicecatalog.ProductType_Values(), false),
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceProductCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ServiceCatalogConn(ctx)

	input := &servicecatalog.CreateProductInput{
		IdempotencyToken: aws.String(id.UniqueId()),
		Name:             aws.String(d.Get("name").(string)),
		Owner:            aws.String(d.Get("owner").(string)),
		ProductType:      aws.String(d.Get("type").(string)),
		ProvisioningArtifactParameters: expandProvisioningArtifactParameters(
			d.Get("provisioning_artifact_parameters").([]interface{})[0].(map[string]interface{}),
		),
		Tags: GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("accept_language"); ok {
		input.AcceptLanguage = aws.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("distributor"); ok {
		input.Distributor = aws.String(v.(string))
	}

	if v, ok := d.GetOk("support_description"); ok {
		input.SupportDescription = aws.String(v.(string))
	}

	if v, ok := d.GetOk("support_email"); ok {
		input.SupportEmail = aws.String(v.(string))
	}

	if v, ok := d.GetOk("support_url"); ok {
		input.SupportUrl = aws.String(v.(string))
	}

	var output *servicecatalog.CreateProductOutput
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var err error

		output, err = conn.CreateProductWithContext(ctx, input)

		if tfawserr.ErrMessageContains(err, servicecatalog.ErrCodeInvalidParametersException, "profile does not exist") {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		output, err = conn.CreateProductWithContext(ctx, input)
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Service Catalog Product: %s", err)
	}

	if output == nil {
		return sdkdiag.AppendErrorf(diags, "creating Service Catalog Product: empty response")
	}

	if output.ProductViewDetail == nil || output.ProductViewDetail.ProductViewSummary == nil {
		return sdkdiag.AppendErrorf(diags, "creating Service Catalog Product: no product view detail or summary")
	}

	if output.ProvisioningArtifactDetail == nil {
		return sdkdiag.AppendErrorf(diags, "creating Service Catalog Product: no provisioning artifact detail")
	}

	d.SetId(aws.StringValue(output.ProductViewDetail.ProductViewSummary.ProductId))

	if _, err := WaitProductReady(ctx, conn, aws.StringValue(input.AcceptLanguage),
		aws.StringValue(output.ProductViewDetail.ProductViewSummary.ProductId), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for Service Catalog Product (%s) to be ready: %s", d.Id(), err)
	}

	return append(diags, resourceProductRead(ctx, d, meta)...)
}

func resourceProductRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ServiceCatalogConn(ctx)

	output, err := WaitProductReady(ctx, conn, d.Get("accept_language").(string), d.Id(), d.Timeout(schema.TimeoutRead))

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Service Catalog Product (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "describing Service Catalog Product (%s): %s", d.Id(), err)
	}

	if output == nil || output.ProductViewDetail == nil || output.ProductViewDetail.ProductViewSummary == nil {
		return sdkdiag.AppendErrorf(diags, "getting Service Catalog Product (%s): empty response", d.Id())
	}

	pvs := output.ProductViewDetail.ProductViewSummary

	d.Set("arn", output.ProductViewDetail.ProductARN)
	if output.ProductViewDetail.CreatedTime != nil {
		d.Set("created_time", output.ProductViewDetail.CreatedTime.Format(time.RFC3339))
	}
	d.Set("description", pvs.ShortDescription)
	d.Set("distributor", pvs.Distributor)
	d.Set("has_default_path", pvs.HasDefaultPath)
	d.Set("name", pvs.Name)
	d.Set("owner", pvs.Owner)
	d.Set("status", output.ProductViewDetail.Status)
	d.Set("support_description", pvs.SupportDescription)
	d.Set("support_email", pvs.SupportEmail)
	d.Set("support_url", pvs.SupportUrl)
	d.Set("type", pvs.Type)

	SetTagsOut(ctx, output.Tags)

	return diags
}

func resourceProductUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ServiceCatalogConn(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		input := &servicecatalog.UpdateProductInput{
			Id: aws.String(d.Id()),
		}

		if v, ok := d.GetOk("accept_language"); ok {
			input.AcceptLanguage = aws.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			input.Description = aws.String(v.(string))
		}

		if v, ok := d.GetOk("distributor"); ok {
			input.Distributor = aws.String(v.(string))
		}

		if v, ok := d.GetOk("name"); ok {
			input.Name = aws.String(v.(string))
		}

		if v, ok := d.GetOk("owner"); ok {
			input.Owner = aws.String(v.(string))
		}

		if v, ok := d.GetOk("support_description"); ok {
			input.SupportDescription = aws.String(v.(string))
		}

		if v, ok := d.GetOk("support_email"); ok {
			input.SupportEmail = aws.String(v.(string))
		}

		if v, ok := d.GetOk("support_url"); ok {
			input.SupportUrl = aws.String(v.(string))
		}

		err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			_, err := conn.UpdateProductWithContext(ctx, input)

			if tfawserr.ErrMessageContains(err, servicecatalog.ErrCodeInvalidParametersException, "profile does not exist") {
				return retry.RetryableError(err)
			}

			if err != nil {
				return retry.NonRetryableError(err)
			}

			return nil
		})

		if tfresource.TimedOut(err) {
			_, err = conn.UpdateProductWithContext(ctx, input)
		}

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating Service Catalog Product (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := productUpdateTags(ctx, conn, d.Id(), o, n); err != nil {
			return sdkdiag.AppendErrorf(diags, "updating tags for Service Catalog Product (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceProductRead(ctx, d, meta)...)
}

func resourceProductDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ServiceCatalogConn(ctx)

	input := &servicecatalog.DeleteProductInput{
		Id: aws.String(d.Id()),
	}

	if v, ok := d.GetOk("accept_language"); ok {
		input.AcceptLanguage = aws.String(v.(string))
	}

	_, err := conn.DeleteProductWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Service Catalog Product (%s): %s", d.Id(), err)
	}

	if _, err := WaitProductDeleted(ctx, conn, d.Get("accept_language").(string), d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for Service Catalog Product (%s) to be deleted: %s", d.Id(), err)
	}

	return diags
}

func expandProvisioningArtifactParameters(tfMap map[string]interface{}) *servicecatalog.ProvisioningArtifactProperties {
	if tfMap == nil {
		return nil
	}

	apiObject := &servicecatalog.ProvisioningArtifactProperties{}

	if v, ok := tfMap["description"].(string); ok && v != "" {
		apiObject.Description = aws.String(v)
	}

	if v, ok := tfMap["disable_template_validation"].(bool); ok {
		apiObject.DisableTemplateValidation = aws.Bool(v)
	}

	info := make(map[string]*string)

	// schema will enforce that one of these is present
	if v, ok := tfMap["template_physical_id"].(string); ok && v != "" {
		info["ImportFromPhysicalId"] = aws.String(v)
	}

	if v, ok := tfMap["template_url"].(string); ok && v != "" {
		info["LoadTemplateFromURL"] = aws.String(v)
	}

	apiObject.Info = info

	if v, ok := tfMap["name"].(string); ok && v != "" {
		apiObject.Name = aws.String(v)
	}

	if v, ok := tfMap["type"].(string); ok && v != "" {
		apiObject.Type = aws.String(v)
	}

	return apiObject
}

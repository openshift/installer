package imagebuilder

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKDataSource("aws_imagebuilder_image_pipeline")
func DataSourceImagePipeline() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceImagePipelineRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},
			"container_recipe_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_last_run": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_next_run": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"distribution_configuration_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enhanced_image_metadata_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"image_recipe_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_tests_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_tests_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"timeout_minutes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"infrastructure_configuration_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pipeline_execution_start_condition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schedule_expression": {
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
		},
	}
}

func dataSourceImagePipelineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ImageBuilderConn(ctx)

	input := &imagebuilder.GetImagePipelineInput{}

	if v, ok := d.GetOk("arn"); ok {
		input.ImagePipelineArn = aws.String(v.(string))
	}

	output, err := conn.GetImagePipelineWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "getting Image Builder Image Pipeline: %s", err)
	}

	if output == nil || output.ImagePipeline == nil {
		return sdkdiag.AppendErrorf(diags, "getting Image Builder Image Pipeline: empty response")
	}

	imagePipeline := output.ImagePipeline

	d.SetId(aws.StringValue(imagePipeline.Arn))
	d.Set("arn", imagePipeline.Arn)
	d.Set("container_recipe_arn", imagePipeline.ContainerRecipeArn)
	d.Set("date_created", imagePipeline.DateCreated)
	d.Set("date_last_run", imagePipeline.DateLastRun)
	d.Set("date_next_run", imagePipeline.DateNextRun)
	d.Set("date_updated", imagePipeline.DateUpdated)
	d.Set("description", imagePipeline.Description)
	d.Set("distribution_configuration_arn", imagePipeline.DistributionConfigurationArn)
	d.Set("enhanced_image_metadata_enabled", imagePipeline.EnhancedImageMetadataEnabled)
	d.Set("image_recipe_arn", imagePipeline.ImageRecipeArn)

	if imagePipeline.ImageTestsConfiguration != nil {
		d.Set("image_tests_configuration", []interface{}{flattenImageTestsConfiguration(imagePipeline.ImageTestsConfiguration)})
	} else {
		d.Set("image_tests_configuration", nil)
	}

	d.Set("infrastructure_configuration_arn", imagePipeline.InfrastructureConfigurationArn)
	d.Set("name", imagePipeline.Name)
	d.Set("platform", imagePipeline.Platform)

	if imagePipeline.Schedule != nil {
		d.Set("schedule", []interface{}{flattenSchedule(imagePipeline.Schedule)})
	} else {
		d.Set("schedule", nil)
	}

	d.Set("status", imagePipeline.Status)
	d.Set("tags", KeyValueTags(ctx, imagePipeline.Tags).IgnoreAWS().IgnoreConfig(meta.(*conns.AWSClient).IgnoreTagsConfig).Map())

	return diags
}

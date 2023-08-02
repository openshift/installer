package apigatewayv2

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKDataSource("aws_apigatewayv2_export")
func DataSourceExport() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceExportRead,

		Schema: map[string]*schema.Schema{
			"api_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"body": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"1.0"}, false),
			},
			"include_extensions": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"specification": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"OAS30"}, false),
			},
			"stage_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"JSON", "YAML"}, false),
			},
		},
	}
}

func dataSourceExportRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayV2Conn(ctx)

	apiId := d.Get("api_id").(string)

	input := &apigatewayv2.ExportApiInput{
		ApiId:             aws.String(apiId),
		Specification:     aws.String(d.Get("specification").(string)),
		OutputType:        aws.String(d.Get("output_type").(string)),
		IncludeExtensions: aws.Bool(d.Get("include_extensions").(bool)),
	}

	if v, ok := d.GetOk("stage_name"); ok {
		input.StageName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("export_version"); ok {
		input.ExportVersion = aws.String(v.(string))
	}

	export, err := conn.ExportApiWithContext(ctx, input)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "exporting Gateway v2 API (%s): %s", apiId, err)
	}

	d.SetId(apiId)

	d.Set("body", string(export.Body))

	return diags
}

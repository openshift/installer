package cur

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKDataSource("aws_cur_report_definition")
func DataSourceReportDefinition() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceReportDefinitionRead,

		Schema: map[string]*schema.Schema{
			"report_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"time_unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compression": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"additional_schema_elements": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
			"s3_bucket": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"s3_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"s3_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"additional_artifacts": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
			"refresh_closed_reports": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"report_versioning": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceReportDefinitionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CURConn(ctx)

	reportName := d.Get("report_name").(string)

	reportDefinition, err := FindReportDefinitionByName(ctx, conn, reportName)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Report Definition (%s): %s", reportName, err)
	}

	if reportDefinition == nil {
		return sdkdiag.AppendErrorf(diags, "reading Report Definition (%s): not found", reportName)
	}

	d.SetId(aws.StringValue(reportDefinition.ReportName))
	d.Set("report_name", reportDefinition.ReportName)
	d.Set("time_unit", reportDefinition.TimeUnit)
	d.Set("format", reportDefinition.Format)
	d.Set("compression", reportDefinition.Compression)
	d.Set("additional_schema_elements", aws.StringValueSlice(reportDefinition.AdditionalSchemaElements))
	d.Set("s3_bucket", reportDefinition.S3Bucket)
	d.Set("s3_prefix", reportDefinition.S3Prefix)
	d.Set("s3_region", reportDefinition.S3Region)
	d.Set("additional_artifacts", aws.StringValueSlice(reportDefinition.AdditionalArtifacts))
	d.Set("refresh_closed_reports", reportDefinition.RefreshClosedReports)
	d.Set("report_versioning", reportDefinition.ReportVersioning)

	return diags
}

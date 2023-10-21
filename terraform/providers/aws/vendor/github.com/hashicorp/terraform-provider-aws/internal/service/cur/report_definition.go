package cur

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	cur "github.com/aws/aws-sdk-go/service/costandusagereportservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

// @SDKResource("aws_cur_report_definition")
func ResourceReportDefinition() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceReportDefinitionCreate,
		ReadWithoutTimeout:   resourceReportDefinitionRead,
		UpdateWithoutTimeout: resourceReportDefinitionUpdate,
		DeleteWithoutTimeout: resourceReportDefinitionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"report_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 256),
					validation.StringMatch(regexp.MustCompile(`[0-9A-Za-z!\-_.*\'()]+`), "The name must be unique, is case sensitive, and can't include spaces."),
				),
			},
			"time_unit": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(cur.TimeUnit_Values(), false),
			},
			"format": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(cur.ReportFormat_Values(), false),
			},
			"compression": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(cur.CompressionFormat_Values(), false),
			},
			"additional_schema_elements": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(cur.SchemaElement_Values(), false),
				},
				Required: true,
				ForceNew: true,
			},
			"s3_bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"s3_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"s3_region": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(cur.AWSRegion_Values(), false),
			},
			"additional_artifacts": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice(cur.AdditionalArtifact_Values(), false),
				},
				Optional: true,
			},
			"refresh_closed_reports": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"report_versioning": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      cur.ReportVersioningCreateNewReport,
				ValidateFunc: validation.StringInSlice(cur.ReportVersioning_Values(), false),
			},
		},
	}
}

func resourceReportDefinitionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CURConn(ctx)

	reportName := d.Get("report_name").(string)
	additionalArtifacts := flex.ExpandStringSet(d.Get("additional_artifacts").(*schema.Set))
	compression := d.Get("compression").(string)
	format := d.Get("format").(string)
	prefix := d.Get("s3_prefix").(string)
	reportVersioning := d.Get("report_versioning").(string)

	additionalArtifactsList := make([]string, 0)
	for i := 0; i < len(additionalArtifacts); i++ {
		additionalArtifactsList = append(additionalArtifactsList, *additionalArtifacts[i])
	}

	err := CheckReportDefinitionPropertyCombination(
		additionalArtifactsList,
		compression,
		format,
		prefix,
		reportVersioning,
	)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Cost And Usage Report Definition (%s): %s", reportName, err)
	}

	reportDefinition := &cur.ReportDefinition{
		ReportName:               aws.String(reportName),
		TimeUnit:                 aws.String(d.Get("time_unit").(string)),
		Format:                   aws.String(format),
		Compression:              aws.String(compression),
		AdditionalSchemaElements: flex.ExpandStringSet(d.Get("additional_schema_elements").(*schema.Set)),
		S3Bucket:                 aws.String(d.Get("s3_bucket").(string)),
		S3Prefix:                 aws.String(prefix),
		S3Region:                 aws.String(d.Get("s3_region").(string)),
		AdditionalArtifacts:      additionalArtifacts,
		RefreshClosedReports:     aws.Bool(d.Get("refresh_closed_reports").(bool)),
		ReportVersioning:         aws.String(reportVersioning),
	}

	reportDefinitionInput := &cur.PutReportDefinitionInput{
		ReportDefinition: reportDefinition,
	}

	_, err = conn.PutReportDefinitionWithContext(ctx, reportDefinitionInput)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Cost And Usage Report Definition (%s): %s", reportName, err)
	}

	d.SetId(reportName)

	return append(diags, resourceReportDefinitionRead(ctx, d, meta)...)
}

func resourceReportDefinitionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CURConn(ctx)

	reportDefinition, err := FindReportDefinitionByName(ctx, conn, d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Cost And Usage Report Definition (%s): %s", d.Id(), err)
	}

	if reportDefinition == nil {
		if d.IsNewResource() {
			return sdkdiag.AppendErrorf(diags, "reading Cost And Usage Report Definition (%s): not found after creation", d.Id())
		}
		log.Printf("[WARN] Cost And Usage Report Definition (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	reportName := aws.StringValue(reportDefinition.ReportName)
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "cur",
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("definition/%s", reportName),
	}.String()

	d.Set("arn", arn)

	d.SetId(reportName)
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

func resourceReportDefinitionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CURConn(ctx)

	additionalArtifacts := flex.ExpandStringSet(d.Get("additional_artifacts").(*schema.Set))
	compression := d.Get("compression").(string)
	format := d.Get("format").(string)
	prefix := d.Get("s3_prefix").(string)
	reportVersioning := d.Get("report_versioning").(string)

	additionalArtifactsList := make([]string, 0)
	for i := 0; i < len(additionalArtifacts); i++ {
		additionalArtifactsList = append(additionalArtifactsList, *additionalArtifacts[i])
	}

	err := CheckReportDefinitionPropertyCombination(
		additionalArtifactsList,
		compression,
		format,
		prefix,
		reportVersioning,
	)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating Cost And Usage Report Definition (%s): %s", d.Id(), err)
	}

	reportName := d.Get("report_name").(string)

	reportDefinition := &cur.ReportDefinition{
		ReportName:               aws.String(reportName),
		TimeUnit:                 aws.String(d.Get("time_unit").(string)),
		Format:                   aws.String(format),
		Compression:              aws.String(compression),
		AdditionalSchemaElements: flex.ExpandStringSet(d.Get("additional_schema_elements").(*schema.Set)),
		S3Bucket:                 aws.String(d.Get("s3_bucket").(string)),
		S3Prefix:                 aws.String(prefix),
		S3Region:                 aws.String(d.Get("s3_region").(string)),
		AdditionalArtifacts:      additionalArtifacts,
		RefreshClosedReports:     aws.Bool(d.Get("refresh_closed_reports").(bool)),
		ReportVersioning:         aws.String(reportVersioning),
	}

	reportDefinitionInput := &cur.ModifyReportDefinitionInput{
		ReportDefinition: reportDefinition,
		ReportName:       aws.String(reportName),
	}

	_, err = conn.ModifyReportDefinitionWithContext(ctx, reportDefinitionInput)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating Cost And Usage Report Definition (%s): %s", d.Id(), err)
	}

	return append(diags, resourceReportDefinitionRead(ctx, d, meta)...)
}

func resourceReportDefinitionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CURConn(ctx)

	log.Printf("[DEBUG] Deleting Cost And Usage Report Definition: %s", d.Id())
	_, err := conn.DeleteReportDefinitionWithContext(ctx, &cur.DeleteReportDefinitionInput{
		ReportName: aws.String(d.Id()),
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Cost And Usage Report Definition (%s): %s", d.Id(), err)
	}

	return diags
}

func CheckReportDefinitionPropertyCombination(additionalArtifacts []string, compression string, format string, prefix string, reportVersioning string) error {
	// perform various combination checks, AWS API unhelpfully just returns an empty ValidationException
	// these combinations have been determined from the Create Report AWS Console Web Form

	hasAthena := false

	for _, artifact := range additionalArtifacts {
		if artifact == cur.AdditionalArtifactAthena {
			hasAthena = true
			break
		}
	}

	if hasAthena {
		if len(additionalArtifacts) > 1 {
			return fmt.Errorf(
				"When %s exists within additional_artifacts, no other artifact type can be declared",
				cur.AdditionalArtifactAthena,
			)
		}

		if len(prefix) == 0 {
			return fmt.Errorf(
				"When %s exists within additional_artifacts, prefix cannot be empty",
				cur.AdditionalArtifactAthena,
			)
		}

		if reportVersioning != cur.ReportVersioningOverwriteReport {
			return fmt.Errorf(
				"When %s exists within additional_artifacts, report_versioning must be %s",
				cur.AdditionalArtifactAthena,
				cur.ReportVersioningOverwriteReport,
			)
		}

		if format != cur.ReportFormatParquet {
			return fmt.Errorf(
				"When %s exists within additional_artifacts, both format and compression must be %s",
				cur.AdditionalArtifactAthena,
				cur.ReportFormatParquet,
			)
		}
	} else if len(additionalArtifacts) > 0 && (format == cur.ReportFormatParquet) {
		return fmt.Errorf(
			"When additional_artifacts includes %s and/or %s, format must not be %s",
			cur.AdditionalArtifactQuicksight,
			cur.AdditionalArtifactRedshift,
			cur.ReportFormatParquet,
		)
	}

	if format == cur.ReportFormatParquet {
		if compression != cur.CompressionFormatParquet {
			return fmt.Errorf(
				"When format is %s, compression must also be %s",
				cur.ReportFormatParquet,
				cur.CompressionFormatParquet,
			)
		}
	} else {
		if compression == cur.CompressionFormatParquet {
			return fmt.Errorf(
				"When format is %s, compression must not be %s",
				format,
				cur.CompressionFormatParquet,
			)
		}
	}
	// end checks

	return nil
}

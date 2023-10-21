package ce

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_ce_anomaly_monitor", name="Anomaly Monitor")
// @Tags(identifierAttribute="id")
func ResourceAnomalyMonitor() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAnomalyMonitorCreate,
		ReadWithoutTimeout:   resourceAnomalyMonitorRead,
		UpdateWithoutTimeout: resourceAnomalyMonitorUpdate,
		DeleteWithoutTimeout: resourceAnomalyMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"monitor_dimension": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"monitor_specification"},
				ValidateFunc:  validation.StringInSlice(costexplorer.MonitorDimension_Values(), false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 1024),
					validation.StringMatch(regexp.MustCompile(`[\\S\\s]*`), "Must be a valid Anomaly Monitor Name matching expression: [\\S\\s]*")),
			},
			"monitor_specification": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: verify.SuppressEquivalentJSONDiffs,
				ConflictsWith:    []string{"monitor_dimension"},
			},
			"monitor_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(costexplorer.MonitorType_Values(), false),
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceAnomalyMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).CEConn(ctx)

	input := &costexplorer.CreateAnomalyMonitorInput{
		AnomalyMonitor: &costexplorer.AnomalyMonitor{
			MonitorName: aws.String(d.Get("name").(string)),
			MonitorType: aws.String(d.Get("monitor_type").(string)),
		},
		ResourceTags: GetTagsIn(ctx),
	}
	switch d.Get("monitor_type").(string) {
	case costexplorer.MonitorTypeDimensional:
		if v, ok := d.GetOk("monitor_dimension"); ok {
			input.AnomalyMonitor.MonitorDimension = aws.String(v.(string))
		} else {
			return diag.Errorf("If Monitor Type is %s, dimension attrribute is required", costexplorer.MonitorTypeDimensional)
		}
	case costexplorer.MonitorTypeCustom:
		if v, ok := d.GetOk("monitor_specification"); ok {
			expression := costexplorer.Expression{}

			if err := json.Unmarshal([]byte(v.(string)), &expression); err != nil {
				return diag.Errorf("parsing specification: %s", err)
			}

			input.AnomalyMonitor.MonitorSpecification = &expression
		} else {
			return diag.Errorf("If Monitor Type is %s, dimension attrribute is required", costexplorer.MonitorTypeCustom)
		}
	}

	resp, err := conn.CreateAnomalyMonitorWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("creating Anomaly Monitor: %s", err)
	}

	if resp == nil || resp.MonitorArn == nil {
		return diag.Errorf("creating Cost Explorer Anomaly Monitor resource (%s): empty output", d.Get("name").(string))
	}

	d.SetId(aws.StringValue(resp.MonitorArn))

	return resourceAnomalyMonitorRead(ctx, d, meta)
}

func resourceAnomalyMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).CEConn(ctx)

	monitor, err := FindAnomalyMonitorByARN(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		create.LogNotFoundRemoveState(names.CE, create.ErrActionReading, ResNameAnomalyMonitor, d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.CE, create.ErrActionReading, ResNameAnomalyMonitor, d.Id(), err)
	}

	if monitor.MonitorSpecification != nil {
		specificationToJson, err := json.Marshal(monitor.MonitorSpecification)
		if err != nil {
			return diag.Errorf("parsing specification response: %s", err)
		}
		specificationToSet, err := structure.NormalizeJsonString(string(specificationToJson))

		if err != nil {
			return diag.Errorf("Specification (%s) is invalid JSON: %s", specificationToSet, err)
		}

		d.Set("monitor_specification", specificationToSet)
	}

	d.Set("arn", monitor.MonitorArn)
	d.Set("monitor_dimension", monitor.MonitorDimension)
	d.Set("name", monitor.MonitorName)
	d.Set("monitor_type", monitor.MonitorType)

	return nil
}

func resourceAnomalyMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).CEConn(ctx)
	requestUpdate := false

	input := &costexplorer.UpdateAnomalyMonitorInput{
		MonitorArn: aws.String(d.Id()),
	}

	if d.HasChange("name") {
		input.MonitorName = aws.String(d.Get("name").(string))
		requestUpdate = true
	}

	if requestUpdate {
		_, err := conn.UpdateAnomalyMonitorWithContext(ctx, input)

		if err != nil {
			return create.DiagError(names.CE, create.ErrActionUpdating, ResNameAnomalyMonitor, d.Id(), err)
		}
	}

	return resourceAnomalyMonitorRead(ctx, d, meta)
}

func resourceAnomalyMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).CEConn(ctx)

	_, err := conn.DeleteAnomalyMonitorWithContext(ctx, &costexplorer.DeleteAnomalyMonitorInput{MonitorArn: aws.String(d.Id())})

	if err != nil && tfawserr.ErrCodeEquals(err, costexplorer.ErrCodeUnknownMonitorException) {
		return nil
	}

	if err != nil {
		return create.DiagError(names.CE, create.ErrActionDeleting, ResNameAnomalyMonitor, d.Id(), err)
	}

	return nil
}

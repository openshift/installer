package logs

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceMetricFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetricFilterUpdate,
		Read:   resourceMetricFilterRead,
		Update: resourceMetricFilterUpdate,
		Delete: resourceMetricFilterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMetricFilterImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validLogMetricFilterName,
			},

			"pattern": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
				StateFunc: func(v interface{}) string {
					s, ok := v.(string)
					if !ok {
						return ""
					}
					return strings.TrimSpace(s)
				},
			},

			"log_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validLogGroupName,
			},

			"metric_transformation": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validLogMetricFilterTransformationName,
						},
						"namespace": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validLogMetricFilterTransformationName,
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(0, 100),
						},
						"default_value": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: verify.ValidTypeStringNullableFloat,
						},
						"dimensions": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"unit": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      cloudwatchlogs.StandardUnitNone,
							ValidateFunc: validation.StringInSlice(cloudwatchlogs.StandardUnit_Values(), false),
						},
					},
				},
			},
		},
	}
}

func resourceMetricFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).LogsConn

	name := d.Get("name").(string)
	logGroupName := d.Get("log_group_name").(string)

	input := cloudwatchlogs.PutMetricFilterInput{
		FilterName:    aws.String(name),
		FilterPattern: aws.String(strings.TrimSpace(d.Get("pattern").(string))),
		LogGroupName:  aws.String(logGroupName),
	}

	transformations := d.Get("metric_transformation").([]interface{})
	o := transformations[0].(map[string]interface{})
	input.MetricTransformations = expandMetricTransformations(o)

	// Creating multiple filters on the same log group can sometimes cause
	// clashes, so use a mutex here (and on deletion) to serialise actions on
	// log groups.
	mutex_key := fmt.Sprintf(`log-group-%s`, d.Get(`log_group_name`))
	conns.GlobalMutexKV.Lock(mutex_key)
	defer conns.GlobalMutexKV.Unlock(mutex_key)
	log.Printf("[DEBUG] Creating/Updating CloudWatch Log Metric Filter: %s", input)
	_, err := conn.PutMetricFilter(&input)
	if err != nil {
		return fmt.Errorf("Creating/Updating CloudWatch Log Metric Filter failed: %w", err)
	}

	d.SetId(d.Get("name").(string))

	log.Println("[INFO] CloudWatch Log Metric Filter created/updated")

	return resourceMetricFilterRead(d, meta)
}

func resourceMetricFilterRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).LogsConn

	mf, err := LookupMetricFilter(conn, d.Get("name").(string),
		d.Get("log_group_name").(string), nil)
	if err != nil {
		if tfresource.NotFound(err) {
			log.Printf("[WARN] Removing CloudWatch Log Metric Filter as it is gone")
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Failed reading CloudWatch Log Metric Filter: %w", err)
	}

	log.Printf("[DEBUG] Found CloudWatch Log Metric Filter: %s", mf)

	d.Set("name", mf.FilterName)
	d.Set("pattern", mf.FilterPattern)
	if err := d.Set("metric_transformation", flattenMetricTransformations(mf.MetricTransformations)); err != nil {
		return fmt.Errorf("setting metric_transformation: %w", err)
	}

	return nil
}

func LookupMetricFilter(conn *cloudwatchlogs.CloudWatchLogs,
	name, logGroupName string, nextToken *string) (*cloudwatchlogs.MetricFilter, error) {

	input := cloudwatchlogs.DescribeMetricFiltersInput{
		FilterNamePrefix: aws.String(name),
		LogGroupName:     aws.String(logGroupName),
		NextToken:        nextToken,
	}
	log.Printf("[DEBUG] Reading CloudWatch Log Metric Filter: %s", input)
	resp, err := conn.DescribeMetricFilters(&input)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, cloudwatchlogs.ErrCodeResourceNotFoundException) {
			return nil, &resource.NotFoundError{
				Message: fmt.Sprintf("CloudWatch Log Metric Filter %q / %q not found via"+
					" initial DescribeMetricFilters call", name, logGroupName),
				LastError:   err,
				LastRequest: input,
			}
		}

		return nil, fmt.Errorf("Failed describing CloudWatch Log Metric Filter: %w", err)
	}

	for _, mf := range resp.MetricFilters {
		if aws.StringValue(mf.FilterName) == name {
			return mf, nil
		}
	}

	if resp.NextToken != nil {
		return LookupMetricFilter(conn, name, logGroupName, resp.NextToken)
	}

	return nil, &resource.NotFoundError{
		Message: fmt.Sprintf("CloudWatch Log Metric Filter %q / %q not found "+
			"in given results from DescribeMetricFilters", name, logGroupName),
		LastResponse: resp,
		LastRequest:  input,
	}
}

func resourceMetricFilterDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).LogsConn

	input := cloudwatchlogs.DeleteMetricFilterInput{
		FilterName:   aws.String(d.Get("name").(string)),
		LogGroupName: aws.String(d.Get("log_group_name").(string)),
	}
	// Creating multiple filters on the same log group can sometimes cause
	// clashes, so use a mutex here (and on creation) to serialise actions on
	// log groups.
	mutex_key := fmt.Sprintf(`log-group-%s`, d.Get(`log_group_name`))
	conns.GlobalMutexKV.Lock(mutex_key)
	defer conns.GlobalMutexKV.Unlock(mutex_key)
	log.Printf("[INFO] Deleting CloudWatch Log Metric Filter: %s", d.Id())
	_, err := conn.DeleteMetricFilter(&input)
	if err != nil {
		return fmt.Errorf("Error deleting CloudWatch Log Metric Filter: %w", err)
	}
	log.Println("[INFO] CloudWatch Log Metric Filter deleted")

	return nil
}

func resourceMetricFilterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.Split(d.Id(), ":")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		return nil, fmt.Errorf("Unexpected format of ID (%q), expected <log_group_name>:<name>", d.Id())
	}
	logGroupName := idParts[0]
	name := idParts[1]
	d.Set("log_group_name", logGroupName)
	d.Set("name", name)
	d.SetId(name)
	return []*schema.ResourceData{d}, nil
}

func expandMetricTransformations(m map[string]interface{}) []*cloudwatchlogs.MetricTransformation {
	transformation := cloudwatchlogs.MetricTransformation{
		MetricName:      aws.String(m["name"].(string)),
		MetricNamespace: aws.String(m["namespace"].(string)),
		MetricValue:     aws.String(m["value"].(string)),
	}

	if m["default_value"].(string) != "" {
		value, _ := strconv.ParseFloat(m["default_value"].(string), 64)
		transformation.DefaultValue = aws.Float64(value)
	}

	if dims := m["dimensions"].(map[string]interface{}); len(dims) > 0 {
		transformation.Dimensions = flex.ExpandStringMap(dims)
	}

	if v, ok := m["unit"].(string); ok && v != "" {
		transformation.Unit = aws.String(v)
	}

	return []*cloudwatchlogs.MetricTransformation{&transformation}
}

func flattenMetricTransformations(ts []*cloudwatchlogs.MetricTransformation) []interface{} {
	mts := make([]interface{}, 0)
	m := make(map[string]interface{})

	transform := ts[0]
	m["name"] = aws.StringValue(transform.MetricName)
	m["namespace"] = aws.StringValue(transform.MetricNamespace)
	m["value"] = aws.StringValue(transform.MetricValue)

	if transform.DefaultValue == nil {
		m["default_value"] = ""
	} else {
		m["default_value"] = strconv.FormatFloat(aws.Float64Value(transform.DefaultValue), 'f', -1, 64)
	}

	if dims := transform.Dimensions; len(dims) > 0 {
		m["dimensions"] = flex.PointersMapToStringList(dims)
	} else {
		m["dimensions"] = nil
	}

	if transform.Unit != nil {
		m["unit"] = aws.StringValue(transform.Unit)
	}

	mts = append(mts, m)

	return mts
}

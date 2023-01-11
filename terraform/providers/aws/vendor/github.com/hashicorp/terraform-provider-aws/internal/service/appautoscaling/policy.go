package appautoscaling

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

const (
	ResNamePolicy = "Policy"
)

func ResourcePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyCreate,
		Read:   resourcePolicyRead,
		Update: resourcePolicyUpdate,
		Delete: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePolicyImport,
		},

		Schema: map[string]*schema.Schema{
			"alarm_arns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// https://github.com/boto/botocore/blob/9f322b1/botocore/data/autoscaling/2011-01-01/service-2.json#L1862-L1873
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"policy_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "StepScaling",
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scalable_dimension": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"step_scaling_policy_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"adjustment_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cooldown": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"metric_aggregation_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"min_adjustment_magnitude": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"step_adjustment": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_interval_lower_bound": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"metric_interval_upper_bound": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"scaling_adjustment": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"target_tracking_scaling_policy_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"customized_metric_specification": {
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							ConflictsWith: []string{"target_tracking_scaling_policy_configuration.0.predefined_metric_specification"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dimensions": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"metric_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"namespace": {
										Type:     schema.TypeString,
										Required: true,
									},
									"statistic": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(applicationautoscaling.MetricStatistic_Values(), false),
									},
									"unit": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"disable_scale_in": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"predefined_metric_specification": {
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							ConflictsWith: []string{"target_tracking_scaling_policy_configuration.0.customized_metric_specification"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"predefined_metric_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"resource_label": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringLenBetween(0, 1023),
									},
								},
							},
						},
						"scale_in_cooldown": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"scale_out_cooldown": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"target_value": {
							Type:     schema.TypeFloat,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourcePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AppAutoScalingConn

	params, err := getPutScalingPolicyInput(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] ApplicationAutoScaling PutScalingPolicy: %#v", params)
	var resp *applicationautoscaling.PutScalingPolicyOutput
	err = resource.Retry(propagationTimeout, func() *resource.RetryError {
		var err error
		resp, err = conn.PutScalingPolicy(&params)
		if err != nil {
			if tfawserr.ErrMessageContains(err, applicationautoscaling.ErrCodeFailedResourceAccessException, "Rate exceeded") {
				return resource.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, applicationautoscaling.ErrCodeFailedResourceAccessException, "is not authorized to perform") {
				return resource.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, applicationautoscaling.ErrCodeFailedResourceAccessException, "token included in the request is invalid") {
				return resource.RetryableError(err)
			}
			if tfawserr.ErrCodeEquals(err, applicationautoscaling.ErrCodeObjectNotFoundException) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		resp, err = conn.PutScalingPolicy(&params)
	}
	if err != nil {
		return create.Error(names.AppAutoScaling, create.ErrActionCreating, ResNamePolicy, d.Get("name").(string), err)
	}

	d.Set("arn", resp.PolicyARN)
	d.SetId(d.Get("name").(string))
	log.Printf("[INFO] ApplicationAutoScaling scaling PolicyARN: %s", d.Get("arn").(string))

	return resourcePolicyRead(d, meta)
}

func resourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
	var p *applicationautoscaling.ScalingPolicy

	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		var err error
		p, err = getPolicy(d, meta)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, applicationautoscaling.ErrCodeFailedResourceAccessException) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if d.IsNewResource() && p == nil {
			return resource.RetryableError(&resource.NotFoundError{})
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		p, err = getPolicy(d, meta)
	}
	if err != nil {
		return create.Error(names.AppAutoScaling, create.ErrActionReading, ResNamePolicy, d.Id(), err)
	}

	if p == nil && !d.IsNewResource() {
		log.Printf("[WARN] Application AutoScaling Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read ApplicationAutoScaling policy: %s, SP: %s, Obj: %s", d.Get("name"), d.Get("name"), p)

	var alarmARNs = make([]string, 0, len(p.Alarms))
	for _, alarm := range p.Alarms {
		alarmARNs = append(alarmARNs, aws.StringValue(alarm.AlarmARN))
	}
	d.Set("alarm_arns", alarmARNs)
	d.Set("arn", p.PolicyARN)
	d.Set("name", p.PolicyName)
	d.Set("policy_type", p.PolicyType)
	d.Set("resource_id", p.ResourceId)
	d.Set("scalable_dimension", p.ScalableDimension)
	d.Set("service_namespace", p.ServiceNamespace)

	if err := d.Set("step_scaling_policy_configuration", flattenStepScalingPolicyConfiguration(p.StepScalingPolicyConfiguration)); err != nil {
		return create.SettingError(names.AppAutoScaling, ResNamePolicy, d.Id(), "step_scaling_policy_configuration", err)
	}
	if err := d.Set("target_tracking_scaling_policy_configuration", flattenTargetTrackingScalingPolicyConfiguration(p.TargetTrackingScalingPolicyConfiguration)); err != nil {
		return create.SettingError(names.AppAutoScaling, ResNamePolicy, d.Id(), "target_tracking_scaling_policy_configuration", err)
	}

	return nil
}

func resourcePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AppAutoScalingConn

	params, inputErr := getPutScalingPolicyInput(d)
	if inputErr != nil {
		return inputErr
	}

	log.Printf("[DEBUG] Application Autoscaling Update Scaling Policy: %#v", params)
	err := resource.Retry(propagationTimeout, func() *resource.RetryError {
		_, err := conn.PutScalingPolicy(&params)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, applicationautoscaling.ErrCodeFailedResourceAccessException) {
				return resource.RetryableError(err)
			}
			if tfawserr.ErrCodeEquals(err, applicationautoscaling.ErrCodeObjectNotFoundException) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		_, err = conn.PutScalingPolicy(&params)
	}
	if err != nil {
		return create.Error(names.AppAutoScaling, create.ErrActionUpdating, ResNamePolicy, d.Id(), err)
	}

	return resourcePolicyRead(d, meta)
}

func resourcePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AppAutoScalingConn
	p, err := getPolicy(d, meta)
	if err != nil {
		return create.Error(names.AppAutoScaling, create.ErrActionDeleting, ResNamePolicy, d.Id(), err)
	}
	if p == nil {
		return nil
	}

	params := applicationautoscaling.DeleteScalingPolicyInput{
		PolicyName:        aws.String(d.Get("name").(string)),
		ResourceId:        aws.String(d.Get("resource_id").(string)),
		ScalableDimension: aws.String(d.Get("scalable_dimension").(string)),
		ServiceNamespace:  aws.String(d.Get("service_namespace").(string)),
	}
	log.Printf("[DEBUG] Deleting Application AutoScaling Policy opts: %#v", params)
	err = resource.Retry(propagationTimeout, func() *resource.RetryError {
		_, err = conn.DeleteScalingPolicy(&params)

		if tfawserr.ErrCodeEquals(err, applicationautoscaling.ErrCodeFailedResourceAccessException) {
			return resource.RetryableError(err)
		}

		if tfawserr.ErrCodeEquals(err, applicationautoscaling.ErrCodeObjectNotFoundException) {
			return nil
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.DeleteScalingPolicy(&params)
	}

	if err != nil {
		return create.Error(names.AppAutoScaling, create.ErrActionDeleting, ResNamePolicy, d.Id(), err)
	}

	return nil
}

func resourcePolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	idParts, err := ValidPolicyImportInput(d.Id())
	if err != nil {
		return nil, create.Error(names.AppAutoScaling, create.ErrActionImporting, ResNamePolicy, d.Id(), err)
	}

	serviceNamespace := idParts[0]
	resourceId := idParts[1]
	scalableDimension := idParts[2]
	policyName := idParts[3]

	d.Set("service_namespace", serviceNamespace)
	d.Set("resource_id", resourceId)
	d.Set("scalable_dimension", scalableDimension)
	d.Set("name", policyName)
	d.SetId(policyName)
	return []*schema.ResourceData{d}, nil
}

func ValidPolicyImportInput(id string) ([]string, error) {

	idParts := strings.Split(id, "/")
	if len(idParts) < 4 {
		return nil, fmt.Errorf("unexpected format (%q), expected <service-namespace>/<resource-id>/<scalable-dimension>/<policy-name>", id)
	}

	var serviceNamespace, resourceId, scalableDimension, policyName string
	switch idParts[0] {
	case "dynamodb":
		serviceNamespace = idParts[0]

		dimensionIx := 3
		// DynamoDB resource ID can be "/table/tableName" or "/table/tableName/index/indexName"
		if idParts[dimensionIx] == "index" {
			dimensionIx = 5
		}

		resourceId = strings.Join(idParts[1:dimensionIx], "/")
		scalableDimension = idParts[dimensionIx]
		policyName = strings.Join(idParts[dimensionIx+1:], "/")
	default:
		serviceNamespace = idParts[0]
		resourceId = strings.Join(idParts[1:len(idParts)-2], "/")
		scalableDimension = idParts[len(idParts)-2]
		policyName = idParts[len(idParts)-1]

	}

	if serviceNamespace == "" || resourceId == "" || scalableDimension == "" || policyName == "" {
		return nil, fmt.Errorf("unexpected format (%q), expected <service-namespace>/<resource-id>/<scalable-dimension>/<policy-name>", id)
	}

	return []string{serviceNamespace, resourceId, scalableDimension, policyName}, nil
}

// Takes the result of flatmap.Expand for an array of step adjustments and
// returns a []*applicationautoscaling.StepAdjustment.
func expandStepAdjustments(configured []interface{}) ([]*applicationautoscaling.StepAdjustment, error) {
	var adjustments []*applicationautoscaling.StepAdjustment

	// Loop over our configured step adjustments and create an array
	// of aws-sdk-go compatible objects. We're forced to convert strings
	// to floats here because there's no way to detect whether or not
	// an uninitialized, optional schema element is "0.0" deliberately.
	// With strings, we can test for "", which is definitely an empty
	// struct value.
	for _, raw := range configured {
		data := raw.(map[string]interface{})
		a := &applicationautoscaling.StepAdjustment{
			ScalingAdjustment: aws.Int64(int64(data["scaling_adjustment"].(int))),
		}
		if data["metric_interval_lower_bound"] != "" {
			bound := data["metric_interval_lower_bound"]
			switch bound := bound.(type) {
			case string:
				f, err := strconv.ParseFloat(bound, 64)
				if err != nil {
					return nil, fmt.Errorf("metric_interval_lower_bound must be a float value represented as a string")
				}
				a.MetricIntervalLowerBound = aws.Float64(f)
			default:
				return nil, fmt.Errorf("metric_interval_lower_bound isn't a string")
			}
		}
		if data["metric_interval_upper_bound"] != "" {
			bound := data["metric_interval_upper_bound"]
			switch bound := bound.(type) {
			case string:
				f, err := strconv.ParseFloat(bound, 64)
				if err != nil {
					return nil, fmt.Errorf("metric_interval_upper_bound must be a float value represented as a string")
				}
				a.MetricIntervalUpperBound = aws.Float64(f)
			default:
				return nil, fmt.Errorf("metric_interval_upper_bound isn't a string")
			}
		}
		adjustments = append(adjustments, a)
	}

	return adjustments, nil
}

func expandCustomizedMetricSpecification(configured []interface{}) *applicationautoscaling.CustomizedMetricSpecification {
	spec := &applicationautoscaling.CustomizedMetricSpecification{}

	for _, raw := range configured {
		data := raw.(map[string]interface{})
		if v, ok := data["metric_name"]; ok {
			spec.MetricName = aws.String(v.(string))
		}

		if v, ok := data["namespace"]; ok {
			spec.Namespace = aws.String(v.(string))
		}

		if v, ok := data["unit"].(string); ok && v != "" {
			spec.Unit = aws.String(v)
		}

		if v, ok := data["statistic"]; ok {
			spec.Statistic = aws.String(v.(string))
		}

		if s, ok := data["dimensions"].(*schema.Set); ok && s.Len() > 0 {
			dimensions := make([]*applicationautoscaling.MetricDimension, s.Len())
			for i, d := range s.List() {
				dimension := d.(map[string]interface{})
				dimensions[i] = &applicationautoscaling.MetricDimension{
					Name:  aws.String(dimension["name"].(string)),
					Value: aws.String(dimension["value"].(string)),
				}
			}
			spec.Dimensions = dimensions
		}
	}
	return spec
}

func expandPredefinedMetricSpecification(configured []interface{}) *applicationautoscaling.PredefinedMetricSpecification {
	spec := &applicationautoscaling.PredefinedMetricSpecification{}

	for _, raw := range configured {
		data := raw.(map[string]interface{})

		if v, ok := data["predefined_metric_type"]; ok {
			spec.PredefinedMetricType = aws.String(v.(string))
		}

		if v, ok := data["resource_label"].(string); ok && v != "" {
			spec.ResourceLabel = aws.String(v)
		}
	}
	return spec
}

func getPutScalingPolicyInput(d *schema.ResourceData) (applicationautoscaling.PutScalingPolicyInput, error) {
	var params = applicationautoscaling.PutScalingPolicyInput{
		PolicyName: aws.String(d.Get("name").(string)),
		ResourceId: aws.String(d.Get("resource_id").(string)),
	}

	if v, ok := d.GetOk("policy_type"); ok {
		params.PolicyType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("service_namespace"); ok {
		params.ServiceNamespace = aws.String(v.(string))
	}

	if v, ok := d.GetOk("scalable_dimension"); ok {
		params.ScalableDimension = aws.String(v.(string))
	}

	if v, ok := d.GetOk("step_scaling_policy_configuration"); ok {
		params.StepScalingPolicyConfiguration = expandStepScalingPolicyConfiguration(v.([]interface{}))
	}

	if l, ok := d.GetOk("target_tracking_scaling_policy_configuration"); ok {
		v := l.([]interface{})
		if len(v) < 1 {
			return params, fmt.Errorf("empty target_tracking_scaling_policy_configuration block")
		}
		ttspCfg := v[0].(map[string]interface{})
		cfg := &applicationautoscaling.TargetTrackingScalingPolicyConfiguration{
			TargetValue: aws.Float64(ttspCfg["target_value"].(float64)),
		}

		if v, ok := ttspCfg["scale_in_cooldown"]; ok {
			cfg.ScaleInCooldown = aws.Int64(int64(v.(int)))
		}

		if v, ok := ttspCfg["scale_out_cooldown"]; ok {
			cfg.ScaleOutCooldown = aws.Int64(int64(v.(int)))
		}

		if v, ok := ttspCfg["disable_scale_in"]; ok {
			cfg.DisableScaleIn = aws.Bool(v.(bool))
		}

		if v, ok := ttspCfg["customized_metric_specification"].([]interface{}); ok && len(v) > 0 {
			cfg.CustomizedMetricSpecification = expandCustomizedMetricSpecification(v)
		}

		if v, ok := ttspCfg["predefined_metric_specification"].([]interface{}); ok && len(v) > 0 {
			cfg.PredefinedMetricSpecification = expandPredefinedMetricSpecification(v)
		}

		params.TargetTrackingScalingPolicyConfiguration = cfg
	}

	return params, nil
}

func getPolicy(d *schema.ResourceData, meta interface{}) (*applicationautoscaling.ScalingPolicy, error) {
	conn := meta.(*conns.AWSClient).AppAutoScalingConn

	params := applicationautoscaling.DescribeScalingPoliciesInput{
		PolicyNames:       []*string{aws.String(d.Get("name").(string))},
		ResourceId:        aws.String(d.Get("resource_id").(string)),
		ScalableDimension: aws.String(d.Get("scalable_dimension").(string)),
		ServiceNamespace:  aws.String(d.Get("service_namespace").(string)),
	}

	log.Printf("[DEBUG] Application AutoScaling Policy Describe Params: %#v", params)
	resp, err := conn.DescribeScalingPolicies(&params)
	if err != nil {
		return nil, fmt.Errorf("describing scaling policies: %s", err)
	}
	if len(resp.ScalingPolicies) == 0 {
		return nil, nil
	}

	return resp.ScalingPolicies[0], nil
}

func expandStepScalingPolicyConfiguration(cfg []interface{}) *applicationautoscaling.StepScalingPolicyConfiguration {
	if len(cfg) < 1 {
		return nil
	}

	out := &applicationautoscaling.StepScalingPolicyConfiguration{}

	m := cfg[0].(map[string]interface{})
	if v, ok := m["adjustment_type"]; ok {
		out.AdjustmentType = aws.String(v.(string))
	}
	if v, ok := m["cooldown"]; ok {
		out.Cooldown = aws.Int64(int64(v.(int)))
	}
	if v, ok := m["metric_aggregation_type"]; ok {
		out.MetricAggregationType = aws.String(v.(string))
	}
	if v, ok := m["min_adjustment_magnitude"].(int); ok && v > 0 {
		out.MinAdjustmentMagnitude = aws.Int64(int64(v))
	}
	if v, ok := m["step_adjustment"].(*schema.Set); ok && v.Len() > 0 {
		out.StepAdjustments, _ = expandStepAdjustments(v.List())
	}

	return out
}

func flattenStepScalingPolicyConfiguration(cfg *applicationautoscaling.StepScalingPolicyConfiguration) []interface{} {
	if cfg == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if cfg.AdjustmentType != nil {
		m["adjustment_type"] = aws.StringValue(cfg.AdjustmentType)
	}
	if cfg.Cooldown != nil {
		m["cooldown"] = aws.Int64Value(cfg.Cooldown)
	}
	if cfg.MetricAggregationType != nil {
		m["metric_aggregation_type"] = aws.StringValue(cfg.MetricAggregationType)
	}
	if cfg.MinAdjustmentMagnitude != nil {
		m["min_adjustment_magnitude"] = aws.Int64Value(cfg.MinAdjustmentMagnitude)
	}
	if cfg.StepAdjustments != nil {
		stepAdjustmentsResource := &schema.Resource{
			Schema: map[string]*schema.Schema{
				"metric_interval_lower_bound": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"metric_interval_upper_bound": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"scaling_adjustment": {
					Type:     schema.TypeInt,
					Required: true,
				},
			},
		}
		m["step_adjustment"] = schema.NewSet(schema.HashResource(stepAdjustmentsResource), flattenStepAdjustments(cfg.StepAdjustments))
	}

	return []interface{}{m}
}

func flattenStepAdjustments(adjs []*applicationautoscaling.StepAdjustment) []interface{} {
	out := make([]interface{}, len(adjs))

	for i, adj := range adjs {
		m := make(map[string]interface{})

		m["scaling_adjustment"] = int(aws.Int64Value(adj.ScalingAdjustment))

		if adj.MetricIntervalLowerBound != nil {
			m["metric_interval_lower_bound"] = fmt.Sprintf("%g", aws.Float64Value(adj.MetricIntervalLowerBound))
		}
		if adj.MetricIntervalUpperBound != nil {
			m["metric_interval_upper_bound"] = fmt.Sprintf("%g", aws.Float64Value(adj.MetricIntervalUpperBound))
		}

		out[i] = m
	}

	return out
}

func flattenTargetTrackingScalingPolicyConfiguration(cfg *applicationautoscaling.TargetTrackingScalingPolicyConfiguration) []interface{} {
	if cfg == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if v := cfg.CustomizedMetricSpecification; v != nil {
		m["customized_metric_specification"] = flattenCustomizedMetricSpecification(v)
	}

	if v := cfg.DisableScaleIn; v != nil {
		m["disable_scale_in"] = aws.BoolValue(v)
	}

	if v := cfg.PredefinedMetricSpecification; v != nil {
		m["predefined_metric_specification"] = flattenPredefinedMetricSpecification(v)
	}

	if v := cfg.ScaleInCooldown; v != nil {
		m["scale_in_cooldown"] = aws.Int64Value(v)
	}

	if v := cfg.ScaleOutCooldown; v != nil {
		m["scale_out_cooldown"] = aws.Int64Value(v)
	}

	if v := cfg.TargetValue; v != nil {
		m["target_value"] = aws.Float64Value(v)
	}

	return []interface{}{m}
}

func flattenCustomizedMetricSpecification(cfg *applicationautoscaling.CustomizedMetricSpecification) []interface{} {
	if cfg == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{}

	if v := cfg.Dimensions; len(v) > 0 {
		m["dimensions"] = flattenMetricDimensions(cfg.Dimensions)
	}

	if v := cfg.MetricName; v != nil {
		m["metric_name"] = aws.StringValue(v)
	}

	if v := cfg.Namespace; v != nil {
		m["namespace"] = aws.StringValue(v)
	}

	if v := cfg.Statistic; v != nil {
		m["statistic"] = aws.StringValue(v)
	}

	if v := cfg.Unit; v != nil {
		m["unit"] = aws.StringValue(v)
	}

	return []interface{}{m}
}

func flattenMetricDimensions(ds []*applicationautoscaling.MetricDimension) []interface{} {
	l := make([]interface{}, len(ds))
	for i, d := range ds {
		if ds == nil {
			continue
		}

		m := map[string]interface{}{}

		if v := d.Name; v != nil {
			m["name"] = aws.StringValue(v)
		}

		if v := d.Value; v != nil {
			m["value"] = aws.StringValue(v)
		}

		l[i] = m
	}
	return l
}

func flattenPredefinedMetricSpecification(cfg *applicationautoscaling.PredefinedMetricSpecification) []interface{} {
	if cfg == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{}

	if v := cfg.PredefinedMetricType; v != nil {
		m["predefined_metric_type"] = aws.StringValue(v)
	}

	if v := cfg.ResourceLabel; v != nil {
		m["resource_label"] = aws.StringValue(v)
	}

	return []interface{}{m}
}

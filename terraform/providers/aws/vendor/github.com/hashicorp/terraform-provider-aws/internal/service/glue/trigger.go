package glue

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/glue"
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

// @SDKResource("aws_glue_trigger", name="Trigger")
// @Tags(identifierAttribute="arn")
func ResourceTrigger() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceTriggerCreate,
		ReadWithoutTimeout:   resourceTriggerRead,
		UpdateWithoutTimeout: resourceTriggerUpdate,
		DeleteWithoutTimeout: resourceTriggerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"actions": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arguments": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"crawler_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"job_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"security_configuration": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"notification_property": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notify_delay_after": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
					},
				},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 2048),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"event_batching_condition": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_size": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
						"batch_window": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      900,
							ValidateFunc: validation.IntBetween(1, 900),
						},
					},
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"predicate": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"job_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"crawler_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"logical_operator": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      glue.LogicalOperatorEquals,
										ValidateFunc: validation.StringInSlice(glue.LogicalOperator_Values(), false),
									},
									"state": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(glue.JobRunState_Values(), false),
									},
									"crawl_state": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(glue.CrawlState_Values(), false),
									},
								},
							},
						},
						"logical": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      glue.LogicalAnd,
							ValidateFunc: validation.StringInSlice(glue.Logical_Values(), false),
						},
					},
				},
			},
			"schedule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_on_creation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(glue.TriggerType_Values(), false),
			},
			"workflow_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTriggerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlueConn(ctx)

	name := d.Get("name").(string)
	triggerType := d.Get("type").(string)
	input := &glue.CreateTriggerInput{
		Actions:         expandActions(d.Get("actions").([]interface{})),
		Name:            aws.String(name),
		Tags:            GetTagsIn(ctx),
		Type:            aws.String(triggerType),
		StartOnCreation: aws.Bool(d.Get("start_on_creation").(bool)),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("event_batching_condition"); ok {
		input.EventBatchingCondition = expandEventBatchingCondition(v.([]interface{}))
	}

	if v, ok := d.GetOk("predicate"); ok {
		input.Predicate = expandPredicate(v.([]interface{}))
	}

	if v, ok := d.GetOk("schedule"); ok {
		input.Schedule = aws.String(v.(string))
	}

	if v, ok := d.GetOk("workflow_name"); ok {
		input.WorkflowName = aws.String(v.(string))
	}

	if d.Get("enabled").(bool) && triggerType != glue.TriggerTypeOnDemand {
		start := true

		if triggerType == glue.TriggerTypeEvent {
			start = false
		}

		input.StartOnCreation = aws.Bool(start)
	}

	if v, ok := d.GetOk("workflow_name"); ok {
		input.WorkflowName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("start_on_creation"); ok {
		input.StartOnCreation = aws.Bool(v.(bool))
	}
	log.Printf("[DEBUG] Creating Glue Trigger: %s", input)
	err := retry.RetryContext(ctx, propagationTimeout, func() *retry.RetryError {
		_, err := conn.CreateTriggerWithContext(ctx, input)
		if err != nil {
			if tfawserr.ErrMessageContains(err, glue.ErrCodeInvalidInputException, "Service is unable to assume provided role") {
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		_, err = conn.CreateTriggerWithContext(ctx, input)
	}
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Glue Trigger (%s): %s", name, err)
	}

	d.SetId(name)

	log.Printf("[DEBUG] Waiting for Glue Trigger (%s) to create", d.Id())
	if _, err := waitTriggerCreated(ctx, conn, d.Id()); err != nil {
		if tfawserr.ErrCodeEquals(err, glue.ErrCodeEntityNotFoundException) {
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "waiting for Glue Trigger (%s) to be Created: %s", d.Id(), err)
	}

	if d.Get("enabled").(bool) && triggerType == glue.TriggerTypeOnDemand {
		input := &glue.StartTriggerInput{
			Name: aws.String(d.Id()),
		}

		log.Printf("[DEBUG] Starting Glue Trigger: %s", input)
		_, err := conn.StartTriggerWithContext(ctx, input)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "starting Glue Trigger (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceTriggerRead(ctx, d, meta)...)
}

func resourceTriggerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlueConn(ctx)

	output, err := FindTriggerByName(ctx, conn, d.Id())
	if err != nil {
		if tfawserr.ErrCodeEquals(err, glue.ErrCodeEntityNotFoundException) {
			log.Printf("[WARN] Glue Trigger (%s) not found, removing from state", d.Id())
			d.SetId("")
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "reading Glue Trigger (%s): %s", d.Id(), err)
	}

	trigger := output.Trigger
	if trigger == nil {
		log.Printf("[WARN] Glue Trigger (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err := d.Set("actions", flattenActions(trigger.Actions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting actions: %s", err)
	}

	triggerARN := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "glue",
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("trigger/%s", d.Id()),
	}.String()
	d.Set("arn", triggerARN)

	d.Set("description", trigger.Description)

	var enabled bool
	state := aws.StringValue(trigger.State)
	d.Set("state", state)

	if aws.StringValue(trigger.Type) == glue.TriggerTypeOnDemand || aws.StringValue(trigger.Type) == glue.TriggerTypeEvent {
		enabled = (state == glue.TriggerStateCreated || state == glue.TriggerStateCreating) && d.Get("enabled").(bool)
	} else {
		enabled = (state == glue.TriggerStateActivated || state == glue.TriggerStateActivating)
	}
	d.Set("enabled", enabled)

	if err := d.Set("predicate", flattenPredicate(trigger.Predicate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting predicate: %s", err)
	}

	if err := d.Set("event_batching_condition", flattenEventBatchingCondition(trigger.EventBatchingCondition)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting event_batching_condition: %s", err)
	}

	d.Set("name", trigger.Name)
	d.Set("schedule", trigger.Schedule)
	d.Set("type", trigger.Type)
	d.Set("workflow_name", trigger.WorkflowName)

	return diags
}

func resourceTriggerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlueConn(ctx)

	if d.HasChanges("actions", "description", "predicate", "schedule", "event_batching_condition") {
		triggerUpdate := &glue.TriggerUpdate{
			Actions: expandActions(d.Get("actions").([]interface{})),
		}

		if v, ok := d.GetOk("description"); ok {
			triggerUpdate.Description = aws.String(v.(string))
		}

		if v, ok := d.GetOk("predicate"); ok {
			triggerUpdate.Predicate = expandPredicate(v.([]interface{}))
		}

		if v, ok := d.GetOk("schedule"); ok {
			triggerUpdate.Schedule = aws.String(v.(string))
		}

		if v, ok := d.GetOk("event_batching_condition"); ok {
			triggerUpdate.EventBatchingCondition = expandEventBatchingCondition(v.([]interface{}))
		}

		input := &glue.UpdateTriggerInput{
			Name:          aws.String(d.Id()),
			TriggerUpdate: triggerUpdate,
		}

		log.Printf("[DEBUG] Updating Glue Trigger: %s", input)
		_, err := conn.UpdateTriggerWithContext(ctx, input)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating Glue Trigger (%s): %s", d.Id(), err)
		}

		if _, err := waitTriggerCreated(ctx, conn, d.Id()); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for Glue Trigger (%s) to be Update: %s", d.Id(), err)
		}
	}

	if d.HasChange("enabled") {
		if d.Get("enabled").(bool) {
			input := &glue.StartTriggerInput{
				Name: aws.String(d.Id()),
			}

			log.Printf("[DEBUG] Starting Glue Trigger: %s", input)
			_, err := conn.StartTriggerWithContext(ctx, input)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "starting Glue Trigger (%s): %s", d.Id(), err)
			}
		} else {
			//Skip if Trigger is type is ON_DEMAND and is in CREATED state as this means the trigger is not running or has ran already.
			if !(d.Get("type").(string) == glue.TriggerTypeOnDemand && d.Get("state").(string) == glue.TriggerStateCreated) {
				input := &glue.StopTriggerInput{
					Name: aws.String(d.Id()),
				}

				log.Printf("[DEBUG] Stopping Glue Trigger: %s", input)
				_, err := conn.StopTriggerWithContext(ctx, input)
				if err != nil {
					return sdkdiag.AppendErrorf(diags, "stopping Glue Trigger (%s): %s", d.Id(), err)
				}
			}
		}
	}

	return diags
}

func resourceTriggerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GlueConn(ctx)

	log.Printf("[DEBUG] Deleting Glue Trigger: %s", d.Id())
	err := deleteTrigger(ctx, conn, d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Glue Trigger (%s): %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Waiting for Glue Trigger (%s) to delete", d.Id())
	if _, err := waitTriggerDeleted(ctx, conn, d.Id()); err != nil {
		if tfawserr.ErrCodeEquals(err, glue.ErrCodeEntityNotFoundException) {
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "waiting for Glue Trigger (%s) to be Deleted: %s", d.Id(), err)
	}

	return diags
}

func deleteTrigger(ctx context.Context, conn *glue.Glue, Name string) error {
	input := &glue.DeleteTriggerInput{
		Name: aws.String(Name),
	}

	_, err := conn.DeleteTriggerWithContext(ctx, input)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, glue.ErrCodeEntityNotFoundException) {
			return nil
		}
		return err
	}

	return nil
}

func expandActions(l []interface{}) []*glue.Action {
	actions := []*glue.Action{}

	for _, mRaw := range l {
		m := mRaw.(map[string]interface{})

		action := &glue.Action{}

		if v, ok := m["crawler_name"].(string); ok && v != "" {
			action.CrawlerName = aws.String(v)
		}

		if v, ok := m["job_name"].(string); ok && v != "" {
			action.JobName = aws.String(v)
		}

		if v, ok := m["arguments"].(map[string]interface{}); ok && len(v) > 0 {
			action.Arguments = flex.ExpandStringMap(v)
		}

		if v, ok := m["timeout"].(int); ok && v > 0 {
			action.Timeout = aws.Int64(int64(v))
		}

		if v, ok := m["security_configuration"].(string); ok && v != "" {
			action.SecurityConfiguration = aws.String(v)
		}

		if v, ok := m["notification_property"].([]interface{}); ok && len(v) > 0 {
			action.NotificationProperty = expandTriggerNotificationProperty(v)
		}

		actions = append(actions, action)
	}

	return actions
}

func expandTriggerNotificationProperty(l []interface{}) *glue.NotificationProperty {
	m := l[0].(map[string]interface{})

	property := &glue.NotificationProperty{}

	if v, ok := m["notify_delay_after"].(int); ok && v > 0 {
		property.NotifyDelayAfter = aws.Int64(int64(v))
	}

	return property
}

func expandConditions(l []interface{}) []*glue.Condition {
	conditions := []*glue.Condition{}

	for _, mRaw := range l {
		m := mRaw.(map[string]interface{})

		condition := &glue.Condition{
			LogicalOperator: aws.String(m["logical_operator"].(string)),
		}

		if v, ok := m["crawler_name"].(string); ok && v != "" {
			condition.CrawlerName = aws.String(v)
		}

		if v, ok := m["crawl_state"].(string); ok && v != "" {
			condition.CrawlState = aws.String(v)
		}

		if v, ok := m["job_name"].(string); ok && v != "" {
			condition.JobName = aws.String(v)
		}

		if v, ok := m["state"].(string); ok && v != "" {
			condition.State = aws.String(v)
		}

		conditions = append(conditions, condition)
	}

	return conditions
}

func expandPredicate(l []interface{}) *glue.Predicate {
	m := l[0].(map[string]interface{})

	predicate := &glue.Predicate{
		Conditions: expandConditions(m["conditions"].([]interface{})),
	}

	if v, ok := m["logical"].(string); ok && v != "" {
		predicate.Logical = aws.String(v)
	}

	return predicate
}

func flattenActions(actions []*glue.Action) []interface{} {
	l := []interface{}{}

	for _, action := range actions {
		m := map[string]interface{}{
			"arguments": aws.StringValueMap(action.Arguments),
			"timeout":   int(aws.Int64Value(action.Timeout)),
		}

		if v := aws.StringValue(action.CrawlerName); v != "" {
			m["crawler_name"] = v
		}

		if v := aws.StringValue(action.JobName); v != "" {
			m["job_name"] = v
		}

		if v := aws.StringValue(action.SecurityConfiguration); v != "" {
			m["security_configuration"] = v
		}

		if v := action.NotificationProperty; v != nil {
			m["notification_property"] = flattenTriggerNotificationProperty(v)
		}

		l = append(l, m)
	}

	return l
}

func flattenConditions(conditions []*glue.Condition) []interface{} {
	l := []interface{}{}

	for _, condition := range conditions {
		m := map[string]interface{}{
			"logical_operator": aws.StringValue(condition.LogicalOperator),
		}

		if v := aws.StringValue(condition.CrawlerName); v != "" {
			m["crawler_name"] = v
		}

		if v := aws.StringValue(condition.CrawlState); v != "" {
			m["crawl_state"] = v
		}

		if v := aws.StringValue(condition.JobName); v != "" {
			m["job_name"] = v
		}

		if v := aws.StringValue(condition.State); v != "" {
			m["state"] = v
		}

		l = append(l, m)
	}

	return l
}

func flattenPredicate(predicate *glue.Predicate) []map[string]interface{} {
	if predicate == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"conditions": flattenConditions(predicate.Conditions),
		"logical":    aws.StringValue(predicate.Logical),
	}

	return []map[string]interface{}{m}
}

func flattenTriggerNotificationProperty(property *glue.NotificationProperty) []map[string]interface{} {
	if property == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"notify_delay_after": aws.Int64Value(property.NotifyDelayAfter),
	}

	return []map[string]interface{}{m}
}

func expandEventBatchingCondition(l []interface{}) *glue.EventBatchingCondition {
	m := l[0].(map[string]interface{})

	ebc := &glue.EventBatchingCondition{
		BatchSize: aws.Int64(int64(m["batch_size"].(int))),
	}

	if v, ok := m["batch_window"].(int); ok && v > 0 {
		ebc.BatchWindow = aws.Int64(int64(v))
	}

	return ebc
}

func flattenEventBatchingCondition(ebc *glue.EventBatchingCondition) []map[string]interface{} {
	if ebc == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"batch_size":   aws.Int64Value(ebc.BatchSize),
		"batch_window": aws.Int64Value(ebc.BatchWindow),
	}

	return []map[string]interface{}{m}
}

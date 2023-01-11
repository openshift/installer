package sfn

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceStateMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceStateMachineCreate,
		Read:   resourceStateMachineRead,
		Update: resourceStateMachineUpdate,
		Delete: resourceStateMachineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"definition": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024*1024), // 1048576
			},

			"logging_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include_execution_data": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"level": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(sfn.LogLevel_Values(), false),
						},
						"log_destination": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validStateMachineName,
			},

			"role_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),

			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      sfn.StateMachineTypeStandard,
				ValidateFunc: validation.StringInSlice(sfn.StateMachineType_Values(), false),
			},

			"tracing_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceStateMachineCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SFNConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	name := d.Get("name").(string)
	input := &sfn.CreateStateMachineInput{
		Definition: aws.String(d.Get("definition").(string)),
		Name:       aws.String(name),
		RoleArn:    aws.String(d.Get("role_arn").(string)),
		Tags:       Tags(tags.IgnoreAWS()),
		Type:       aws.String(d.Get("type").(string)),
	}

	if v, ok := d.GetOk("logging_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.LoggingConfiguration = expandLoggingConfiguration(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("tracing_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.TracingConfiguration = expandTracingConfiguration(v.([]interface{})[0].(map[string]interface{}))
	}

	var output *sfn.CreateStateMachineOutput

	log.Printf("[DEBUG] Creating Step Function State Machine: %s", input)
	err := resource.Retry(stateMachineCreatedTimeout, func() *resource.RetryError {
		var err error

		output, err = conn.CreateStateMachine(input)

		// Note: the instance may be in a deleting mode, hence the retry
		// when creating the step function. This can happen when we are
		// updating the resource (since there is no update API call).
		if tfawserr.ErrCodeEquals(err, sfn.ErrCodeStateMachineDeleting) {
			return resource.RetryableError(err)
		}

		// This is done to deal with IAM eventual consistency
		if tfawserr.ErrCodeEquals(err, "AccessDeniedException") {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		output, err = conn.CreateStateMachine(input)
	}

	if err != nil {
		return fmt.Errorf("error creating Step Function State Machine (%s): %w", name, err)
	}

	d.SetId(aws.StringValue(output.StateMachineArn))

	return resourceStateMachineRead(d, meta)
}

func resourceStateMachineRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SFNConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	output, err := FindStateMachineByARN(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Step Function State Machine (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Step Function State Machine (%s): %w", d.Id(), err)
	}

	d.Set("arn", output.StateMachineArn)
	if output.CreationDate != nil {
		d.Set("creation_date", aws.TimeValue(output.CreationDate).Format(time.RFC3339))
	} else {
		d.Set("creation_date", nil)
	}
	d.Set("definition", output.Definition)
	d.Set("name", output.Name)
	d.Set("role_arn", output.RoleArn)
	d.Set("type", output.Type)
	d.Set("status", output.Status)

	if output.LoggingConfiguration != nil {
		if err := d.Set("logging_configuration", []interface{}{flattenLoggingConfiguration(output.LoggingConfiguration)}); err != nil {
			return fmt.Errorf("error setting logging_configuration: %w", err)
		}
	} else {
		d.Set("logging_configuration", nil)
	}

	if output.TracingConfiguration != nil {
		if err := d.Set("tracing_configuration", []interface{}{flattenTracingConfiguration(output.TracingConfiguration)}); err != nil {
			return fmt.Errorf("error setting tracing_configuration: %w", err)
		}
	} else {
		d.Set("tracing_configuration", nil)
	}

	tags, err := ListTags(conn, d.Id())

	if err != nil {
		if tfawserr.ErrCodeEquals(err, "UnknownOperationException") {
			return nil
		}

		return fmt.Errorf("error listing tags for Step Function State Machine (%s): %w", d.Id(), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceStateMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SFNConn

	if d.HasChangesExcept("tags", "tags_all") {
		// "You must include at least one of definition or roleArn or you will receive a MissingRequiredParameter error"
		input := &sfn.UpdateStateMachineInput{
			StateMachineArn: aws.String(d.Id()),
			Definition:      aws.String(d.Get("definition").(string)),
			RoleArn:         aws.String(d.Get("role_arn").(string)),
		}

		if d.HasChange("logging_configuration") {
			if v, ok := d.GetOk("logging_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
				input.LoggingConfiguration = expandLoggingConfiguration(v.([]interface{})[0].(map[string]interface{}))
			}
		}

		if d.HasChange("tracing_configuration") {
			if v, ok := d.GetOk("tracing_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
				input.TracingConfiguration = expandTracingConfiguration(v.([]interface{})[0].(map[string]interface{}))
			}
		}

		log.Printf("[DEBUG] Updating Step Function State Machine: %s", input)
		_, err := conn.UpdateStateMachine(input)

		if err != nil {
			return fmt.Errorf("error updating Step Function State Machine (%s): %w", d.Id(), err)
		}

		// Handle eventual consistency after update.
		err = resource.Retry(stateMachineUpdatedTimeout, func() *resource.RetryError {
			output, err := FindStateMachineByARN(conn, d.Id())

			if err != nil {
				return resource.NonRetryableError(err)
			}

			if d.HasChange("definition") && !verify.JSONBytesEqual([]byte(aws.StringValue(output.Definition)), []byte(d.Get("definition").(string))) ||
				d.HasChange("role_arn") && aws.StringValue(output.RoleArn) != d.Get("role_arn").(string) ||
				d.HasChange("tracing_configuration.0.enabled") && output.TracingConfiguration != nil && aws.BoolValue(output.TracingConfiguration.Enabled) != d.Get("tracing_configuration.0.enabled").(bool) ||
				d.HasChange("logging_configuration.0.include_execution_data") && output.LoggingConfiguration != nil && aws.BoolValue(output.LoggingConfiguration.IncludeExecutionData) != d.Get("logging_configuration.0.include_execution_data").(bool) ||
				d.HasChange("logging_configuration.0.level") && output.LoggingConfiguration != nil && aws.StringValue(output.LoggingConfiguration.Level) != d.Get("logging_configuration.0.level").(string) {
				return resource.RetryableError(fmt.Errorf("Step Function State Machine (%s) eventual consistency", d.Id()))
			}

			return nil
		})

		if tfresource.TimedOut(err) {
			return fmt.Errorf("timed out waiting for Step Function State Machine (%s) update", d.Id())
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")
		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating tags: %w", err)
		}
	}

	return resourceStateMachineRead(d, meta)
}

func resourceStateMachineDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SFNConn

	_, err := conn.DeleteStateMachine(&sfn.DeleteStateMachineInput{
		StateMachineArn: aws.String(d.Id()),
	})

	if err != nil {
		return fmt.Errorf("error deleting Step Function State Machine (%s): %s", d.Id(), err)
	}

	if _, err := waitStateMachineDeleted(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for Step Function State Machine (%s) deletion: %w", d.Id(), err)
	}

	return nil
}

func expandLoggingConfiguration(tfMap map[string]interface{}) *sfn.LoggingConfiguration {
	if tfMap == nil {
		return nil
	}

	apiObject := &sfn.LoggingConfiguration{}

	if v, ok := tfMap["include_execution_data"].(bool); ok {
		apiObject.IncludeExecutionData = aws.Bool(v)
	}

	if v, ok := tfMap["level"].(string); ok && v != "" {
		apiObject.Level = aws.String(v)
	}

	if v, ok := tfMap["log_destination"].(string); ok && v != "" {
		apiObject.Destinations = []*sfn.LogDestination{{
			CloudWatchLogsLogGroup: &sfn.CloudWatchLogsLogGroup{
				LogGroupArn: aws.String(v),
			},
		}}
	}

	return apiObject
}

func flattenLoggingConfiguration(apiObject *sfn.LoggingConfiguration) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.IncludeExecutionData; v != nil {
		tfMap["include_execution_data"] = aws.BoolValue(v)
	}

	if v := apiObject.Level; v != nil {
		tfMap["level"] = aws.StringValue(v)
	}

	if v := apiObject.Destinations; len(v) > 0 {
		tfMap["log_destination"] = aws.StringValue(v[0].CloudWatchLogsLogGroup.LogGroupArn)
	}

	return tfMap
}

func expandTracingConfiguration(tfMap map[string]interface{}) *sfn.TracingConfiguration {
	if tfMap == nil {
		return nil
	}

	apiObject := &sfn.TracingConfiguration{}

	if v, ok := tfMap["enabled"].(bool); ok {
		apiObject.Enabled = aws.Bool(v)
	}

	return apiObject
}

func flattenTracingConfiguration(apiObject *sfn.TracingConfiguration) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Enabled; v != nil {
		tfMap["enabled"] = aws.BoolValue(v)
	}

	return tfMap
}

package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudConfigAggregateConfigRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudConfigAggregateConfigRuleCreate,
		Read:   resourceAlicloudConfigAggregateConfigRuleRead,
		Update: resourceAlicloudConfigAggregateConfigRuleUpdate,
		Delete: resourceAlicloudConfigAggregateConfigRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"aggregate_config_rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aggregator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_rule_trigger_types": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ConfigurationItemChangeNotification", "ScheduledNotification"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclude_resource_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"input_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"maximum_execution_frequency": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"One_Hour", "Six_Hours", "Three_Hours", "Twelve_Hours", "TwentyFour_Hours"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if fmt.Sprint(d.Get("config_rule_trigger_types")) == "ConfigurationItemChangeNotification" {
						return true
					}
					return false
				},
			},
			"region_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_types_scope": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"risk_level": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3}),
			},
			"source_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_owner": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALIYUN", "CUSTOM_FC"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
				Computed:     true,
			},
			"tag_key_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_value_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudConfigAggregateConfigRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAggregateConfigRule"
	request := make(map[string]interface{})
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	request["ConfigRuleName"] = d.Get("aggregate_config_rule_name")
	request["AggregatorId"] = d.Get("aggregator_id")
	request["ConfigRuleTriggerTypes"] = d.Get("config_rule_trigger_types")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("exclude_resource_ids_scope"); ok {
		request["ExcludeResourceIdsScope"] = v
	}
	if v, ok := d.GetOk("input_parameters"); ok {
		if v, err := convertMaptoJsonString(v.(map[string]interface{})); err == nil {
			request["InputParameters"] = v
		} else {
			return WrapError(err)
		}
	}
	if v, ok := d.GetOk("maximum_execution_frequency"); ok {
		request["MaximumExecutionFrequency"] = v
	}
	if v, ok := d.GetOk("region_ids_scope"); ok {
		request["RegionIdsScope"] = v
	}
	if v, ok := d.GetOk("resource_group_ids_scope"); ok {
		request["ResourceGroupIdsScope"] = v
	}
	request["ResourceTypesScope"] = convertListToCommaSeparate(d.Get("resource_types_scope").([]interface{}))
	request["RiskLevel"] = d.Get("risk_level")
	request["SourceIdentifier"] = d.Get("source_identifier")
	request["SourceOwner"] = d.Get("source_owner")
	if v, ok := d.GetOk("tag_key_scope"); ok {
		request["TagKeyScope"] = v
	}
	if v, ok := d.GetOk("tag_value_scope"); ok {
		request["TagValueScope"] = v
	}
	request["ClientToken"] = buildClientToken("CreateAggregateConfigRule")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_aggregate_config_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AggregatorId"], ":", response["ConfigRuleId"]))

	return resourceAlicloudConfigAggregateConfigRuleUpdate(d, meta)
}
func resourceAlicloudConfigAggregateConfigRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	object, err := configService.DescribeConfigAggregateConfigRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_aggregate_config_rule configService.DescribeConfigAggregateConfigRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("aggregator_id", parts[0])
	d.Set("config_rule_id", parts[1])
	d.Set("aggregate_config_rule_name", object["ConfigRuleName"])
	d.Set("config_rule_trigger_types", object["ConfigRuleTriggerTypes"])
	d.Set("description", object["Description"])
	d.Set("exclude_resource_ids_scope", object["ExcludeResourceIdsScope"])
	d.Set("input_parameters", object["InputParameters"])
	d.Set("maximum_execution_frequency", object["MaximumExecutionFrequency"])
	d.Set("region_ids_scope", object["RegionIdsScope"])
	d.Set("resource_group_ids_scope", object["ResourceGroupIdsScope"])
	d.Set("resource_types_scope", object["Scope"].(map[string]interface{})["ComplianceResourceTypes"])
	d.Set("risk_level", formatInt(object["RiskLevel"]))
	d.Set("source_identifier", object["Source"].(map[string]interface{})["Identifier"])
	d.Set("source_owner", object["Source"].(map[string]interface{})["Owner"])
	d.Set("status", object["ConfigRuleState"])
	d.Set("tag_key_scope", object["TagKeyScope"])
	d.Set("tag_value_scope", object["TagValueScope"])
	return nil
}
func resourceAlicloudConfigAggregateConfigRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	d.Partial(true)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"AggregatorId": parts[0],
		"ConfigRuleId": parts[1],
	}
	if !d.IsNewResource() && d.HasChange("config_rule_trigger_types") {
		update = true
	}
	request["ConfigRuleTriggerTypes"] = d.Get("config_rule_trigger_types")
	if !d.IsNewResource() && d.HasChange("resource_types_scope") {
		update = true
	}
	request["ResourceTypesScope"] = convertListToCommaSeparate(d.Get("resource_types_scope").([]interface{}))
	if !d.IsNewResource() && d.HasChange("risk_level") {
		update = true
	}
	request["RiskLevel"] = d.Get("risk_level")
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("exclude_resource_ids_scope") {
		update = true
		request["ExcludeResourceIdsScope"] = d.Get("exclude_resource_ids_scope")
	}
	if !d.IsNewResource() && d.HasChange("input_parameters") {
		update = true
		if v, err := convertMaptoJsonString(d.Get("input_parameters").(map[string]interface{})); err == nil {
			request["InputParameters"] = v
		} else {
			return WrapError(err)
		}
	}
	if !d.IsNewResource() && d.HasChange("maximum_execution_frequency") {
		update = true
		request["MaximumExecutionFrequency"] = d.Get("maximum_execution_frequency")
	}
	if !d.IsNewResource() && d.HasChange("region_ids_scope") {
		update = true
		request["RegionIdsScope"] = d.Get("region_ids_scope")
	}
	if !d.IsNewResource() && d.HasChange("resource_group_ids_scope") {
		update = true
		request["ResourceGroupIdsScope"] = d.Get("resource_group_ids_scope")
	}
	if !d.IsNewResource() && d.HasChange("tag_key_scope") {
		update = true
		request["TagKeyScope"] = d.Get("tag_key_scope")
	}
	if !d.IsNewResource() && d.HasChange("tag_value_scope") {
		update = true
		request["TagValueScope"] = d.Get("tag_value_scope")
	}
	if update {
		action := "UpdateAggregateConfigRule"
		conn, err := client.NewConfigClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateAggregateConfigRule")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	if d.HasChange("status") {
		update = true
	}
	if update {
		configService := ConfigService{client}
		object, err := configService.DescribeConfigAggregateConfigRule(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if value, exist := object["ConfigRuleState"]; exist && value.(string) != target {
			if target == "ACTIVE" {
				err := configService.ActiveAggregateConfigRules(parts[1], parts[0])
				if err != nil {
					return WrapError(err)
				}
			}
			if target == "INACTIVE" {
				err := configService.DeactiveAggregateConfigRules(parts[1], parts[0])
				if err != nil {
					return WrapError(err)
				}
			}
		}
		d.SetPartial("status")
	}
	d.Partial(false)
	stateConf := BuildStateConf([]string{}, []string{"ACTIVE", "INACTIVE"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, configService.ConfigAggregateConfigRuleStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudConfigAggregateConfigRuleRead(d, meta)
}
func resourceAlicloudConfigAggregateConfigRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteAggregateConfigRules"
	var response map[string]interface{}
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AggregatorId":  parts[0],
		"ConfigRuleIds": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ConfigRuleCanNotDelete", "ConfigRuleNotExists", "Invalid.AggregatorId.Value", "Invalid.ConfigRuleId.Value"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

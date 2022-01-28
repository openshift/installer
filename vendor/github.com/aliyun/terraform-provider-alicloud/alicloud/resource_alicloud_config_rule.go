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

func resourceAlicloudConfigRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudConfigRuleCreate,
		Read:   resourceAlicloudConfigRuleRead,
		Update: resourceAlicloudConfigRuleUpdate,
		Delete: resourceAlicloudConfigRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"member_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Removed:  "Field 'member_id' has been removed from provider version 1.124.1.",
			},
			"multi_account": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Removed:  "Field 'multi_account' has been removed from provider version 1.124.1.",
			},
			"scope_compliance_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'scope_compliance_resource_id' has been removed from provider version 1.124.1.",
			},
			"config_rule_trigger_types": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"source_detail_message_type"},
			},
			"source_detail_message_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'source_detail_message_type' has been deprecated from provider version 1.124.1. New field 'config_rule_trigger_types' instead.",
				ConflictsWith: []string{"config_rule_trigger_types"},
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
					if fmt.Sprint(d.Get("config_rule_trigger_types")) == "ConfigurationItemChangeNotification" || fmt.Sprint(d.Get("source_detail_message_type")) == "ConfigurationItemChangeNotification" {
						return true
					}
					return false
				},
				ConflictsWith: []string{"source_maximum_execution_frequency"},
			},
			"source_maximum_execution_frequency": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if fmt.Sprint(d.Get("config_rule_trigger_types")) == "ConfigurationItemChangeNotification" || fmt.Sprint(d.Get("source_detail_message_type")) == "ConfigurationItemChangeNotification" {
						return true
					}
					return false
				},
				ValidateFunc:  validation.StringInSlice([]string{"One_Hour", "Six_Hours", "Three_Hours", "Twelve_Hours", "TwentyFour_Hours"}, false),
				Deprecated:    "Field 'source_maximum_execution_frequency' has been deprecated from provider version 1.124.1. New field 'maximum_execution_frequency' instead.",
				ConflictsWith: []string{"maximum_execution_frequency"},
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
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"scope_compliance_resource_types"},
			},
			"scope_compliance_resource_types": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Deprecated:    "Field 'scope_compliance_resource_types' has been deprecated from provider version 1.124.1. New field 'resource_types_scope' instead.",
				ConflictsWith: []string{"resource_types_scope"},
			},
			"risk_level": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3}),
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func resourceAlicloudConfigRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateConfigRule"
	request := make(map[string]interface{})
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("config_rule_trigger_types"); ok {
		request["ConfigRuleTriggerTypes"] = v
	} else if v, ok := d.GetOk("source_detail_message_type"); ok {
		request["ConfigRuleTriggerTypes"] = v
	} else {
		return WrapError(Error("[ERROR] Argument 'source_detail_message_type' or 'config_rule_trigger_types' must be set one!"))
	}
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
	} else if v, ok := d.GetOk("source_maximum_execution_frequency"); ok {
		request["MaximumExecutionFrequency"] = v
	}
	if v, ok := d.GetOk("region_ids_scope"); ok {
		request["RegionIdsScope"] = v
	}
	if v, ok := d.GetOk("resource_group_ids_scope"); ok {
		request["ResourceGroupIdsScope"] = v
	}
	if v, ok := d.GetOk("resource_types_scope"); ok && v != nil {
		request["ResourceTypesScope"] = convertListToCommaSeparate(v.([]interface{}))
	} else if v, ok := d.GetOk("scope_compliance_resource_types"); ok && v != nil {
		request["ResourceTypesScope"] = convertListToCommaSeparate(v.([]interface{}))
	} else {
		return WrapError(Error("[ERROR] Argument 'scope_compliance_resource_types' or 'resource_types_scope' must be set one!"))
	}
	request["RiskLevel"] = d.Get("risk_level")
	request["ConfigRuleName"] = d.Get("rule_name")
	request["SourceIdentifier"] = d.Get("source_identifier")
	request["SourceOwner"] = d.Get("source_owner")
	if v, ok := d.GetOk("tag_key_scope"); ok {
		request["TagKeyScope"] = v
	}
	if v, ok := d.GetOk("tag_value_scope"); ok {
		request["TagValueScope"] = v
	}
	request["ClientToken"] = buildClientToken("CreateConfigRule")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ConfigRuleId"]))

	return resourceAlicloudConfigRuleUpdate(d, meta)
}
func resourceAlicloudConfigRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	object, err := configService.DescribeConfigRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_rule configService.DescribeConfigRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("exclude_resource_ids_scope", object["ExcludeResourceIdsScope"])
	d.Set("input_parameters", object["InputParameters"])
	d.Set("maximum_execution_frequency", object["MaximumExecutionFrequency"])
	d.Set("source_maximum_execution_frequency", object["MaximumExecutionFrequency"])
	d.Set("region_ids_scope", object["RegionIdsScope"])
	d.Set("resource_group_ids_scope", object["ResourceGroupIdsScope"])
	d.Set("resource_types_scope", object["Scope"].(map[string]interface{})["ComplianceResourceTypes"])
	d.Set("scope_compliance_resource_types", object["Scope"].(map[string]interface{})["ComplianceResourceTypes"])
	d.Set("risk_level", formatInt(object["RiskLevel"]))
	d.Set("rule_name", object["ConfigRuleName"])
	d.Set("source_identifier", object["Source"].(map[string]interface{})["Identifier"])
	d.Set("source_owner", object["Source"].(map[string]interface{})["Owner"])
	d.Set("status", object["ConfigRuleState"])
	d.Set("tag_key_scope", object["TagKeyScope"])
	d.Set("tag_value_scope", object["TagValueScope"])
	if v := object["Source"].(map[string]interface{})["SourceDetails"].([]interface{}); len(v) > 0 {
		d.Set("config_rule_trigger_types", v[0].(map[string]interface{})["MessageType"])
		d.Set("source_detail_message_type", v[0].(map[string]interface{})["MessageType"])
	}
	return nil
}
func resourceAlicloudConfigRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	d.Partial(true)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ConfigRuleId": d.Id(),
	}
	if !d.IsNewResource() && (d.HasChange("config_rule_trigger_types") || d.HasChange("source_detail_message_type")) {
		update = true
	}
	if !d.IsNewResource() && (d.HasChange("config_rule_trigger_types") || d.HasChange("source_detail_message_type")) {
		update = true
	}
	if _, ok := d.GetOk("config_rule_trigger_types"); ok {
		request["ConfigRuleTriggerTypes"] = d.Get("config_rule_trigger_types")
	} else {
		request["ConfigRuleTriggerTypes"] = d.Get("source_detail_message_type")
	}
	if !d.IsNewResource() && (d.HasChange("resource_types_scope") || d.HasChange("scope_compliance_resource_types")) {
		update = true
	}
	if _, ok := d.GetOk("resource_types_scope"); ok {
		request["ResourceTypesScope"] = convertListToCommaSeparate(d.Get("resource_types_scope").([]interface{}))
	} else {
		request["ResourceTypesScope"] = convertListToCommaSeparate(d.Get("scope_compliance_resource_types").([]interface{}))
	}
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
	if !d.IsNewResource() && d.HasChange("maximum_execution_frequency") || d.HasChange("source_maximum_execution_frequency") {
		update = true
		if _, ok := d.GetOk("maximum_execution_frequency"); ok {
			request["MaximumExecutionFrequency"] = d.Get("maximum_execution_frequency")
		} else {
			request["MaximumExecutionFrequency"] = d.Get("source_maximum_execution_frequency")
		}
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
		action := "UpdateConfigRule"
		conn, err := client.NewConfigClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateConfigRule")
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
		stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, configService.ConfigRuleStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	update = false
	if d.HasChange("status") {
		update = true
	}
	if update {
		configService := ConfigService{client}
		object, err := configService.DescribeConfigRule(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if value, exist := object["ConfigRuleState"]; exist && value.(string) != target {
			if target == "ACTIVE" {
				err := configService.ActiveConfigRule(d.Id())
				if err != nil {
					return WrapError(err)
				}
			}
			if target == "INACTIVE" {
				err := configService.StopConfigRule(d.Id())
				if err != nil {
					return WrapError(err)
				}
			}
		}
		d.SetPartial("status")
	}
	d.Partial(false)
	stateConf := BuildStateConf([]string{}, []string{"ACTIVE", "INACTIVE"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, configService.ConfigRuleStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudConfigRuleRead(d, meta)
}
func resourceAlicloudConfigRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteConfigRules"
	var response map[string]interface{}
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ConfigRuleIds": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"AccountNotExisted", "ConfigRuleNotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

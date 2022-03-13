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

func resourceAlicloudConfigAggregateCompliancePack() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudConfigAggregateCompliancePackCreate,
		Read:   resourceAlicloudConfigAggregateCompliancePackRead,
		Update: resourceAlicloudConfigAggregateCompliancePackUpdate,
		Delete: resourceAlicloudConfigAggregateCompliancePackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"aggregate_compliance_pack_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aggregator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"compliance_pack_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"config_rules": {
				Type:          schema.TypeSet,
				Optional:      true,
				Deprecated:    "Field 'config_rules' has been deprecated from provider version 1.141.0. New field 'config_rule_ids' instead.",
				ConflictsWith: []string{"config_rule_ids"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_rule_parameters": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameter_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"parameter_value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"managed_rule_identifier": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"config_rule_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"config_rules"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_rule_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"risk_level": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudConfigAggregateCompliancePackCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	var response map[string]interface{}
	action := "CreateAggregateCompliancePack"
	request := make(map[string]interface{})
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	request["CompliancePackName"] = d.Get("aggregate_compliance_pack_name")
	request["AggregatorId"] = d.Get("aggregator_id")
	if v, ok := d.GetOk("compliance_pack_template_id"); ok {
		request["CompliancePackTemplateId"] = v
	}
	configRulesMaps := make([]map[string]interface{}, 0)
	if _, ok := d.GetOk("config_rules"); ok {
		for _, configRules := range d.Get("config_rules").(*schema.Set).List() {
			configRulesArg := configRules.(map[string]interface{})
			configRulesMap := map[string]interface{}{
				"ManagedRuleIdentifier": configRulesArg["managed_rule_identifier"],
			}
			configRuleParametersMaps := make([]map[string]interface{}, 0)
			for _, configRuleParameters := range configRulesArg["config_rule_parameters"].(*schema.Set).List() {
				configRuleParametersArg := configRuleParameters.(map[string]interface{})
				configRuleParametersMap := map[string]interface{}{
					"ParameterName":  configRuleParametersArg["parameter_name"],
					"ParameterValue": configRuleParametersArg["parameter_value"],
				}
				configRuleParametersMaps = append(configRuleParametersMaps, configRuleParametersMap)
			}
			configRulesMap["ConfigRuleParameters"] = configRuleParametersMaps
			configRulesMaps = append(configRulesMaps, configRulesMap)
		}
	} else {
		for _, configRuleIds := range d.Get("config_rule_ids").(*schema.Set).List() {
			configRulesIdArg := configRuleIds.(map[string]interface{})
			configRuleIdMap := map[string]interface{}{
				"ConfigRuleId": configRulesIdArg["config_rule_id"],
			}
			configRulesMaps = append(configRulesMaps, configRuleIdMap)
		}
	}

	if v, err := convertArrayObjectToJsonString(configRulesMaps); err == nil {
		request["ConfigRules"] = v
	} else {
		return WrapError(err)
	}
	request["Description"] = d.Get("description")
	request["RiskLevel"] = d.Get("risk_level")
	request["ClientToken"] = buildClientToken("CreateAggregateCompliancePack")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_aggregate_compliance_pack", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AggregatorId"], ":", response["CompliancePackId"]))
	stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, configService.ConfigAggregateCompliancePackStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudConfigAggregateCompliancePackRead(d, meta)
}
func resourceAlicloudConfigAggregateCompliancePackRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	object, err := configService.DescribeConfigAggregateCompliancePack(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_aggregate_compliance_pack configService.DescribeConfigAggregateCompliancePack Failed!!! %s", err)
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
	d.Set("aggregate_compliance_pack_name", object["CompliancePackName"])
	d.Set("compliance_pack_template_id", object["CompliancePackTemplateId"])

	if _, ok := d.GetOk("config_rules"); ok {
		configRules := make([]map[string]interface{}, 0)
		if configRulesList, ok := object["ConfigRules"].([]interface{}); ok {
			for _, v := range configRulesList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"managed_rule_identifier": m1["ManagedRuleIdentifier"],
					}
					if m1["ConfigRuleParameters"] != nil {
						configRuleParametersMaps := make([]map[string]interface{}, 0)
						for _, configRuleParametersValue := range m1["ConfigRuleParameters"].([]interface{}) {
							configRuleParameters := configRuleParametersValue.(map[string]interface{})
							configRuleParametersMap := map[string]interface{}{
								"parameter_name":  configRuleParameters["ParameterName"],
								"parameter_value": configRuleParameters["ParameterValue"],
							}
							configRuleParametersMaps = append(configRuleParametersMaps, configRuleParametersMap)
						}
						temp1["config_rule_parameters"] = configRuleParametersMaps
					}
					configRules = append(configRules, temp1)

				}
			}
		}
		d.Set("config_rules", configRules)
	} else {

		configRuleIds := make([]map[string]interface{}, 0)
		if configRuleIdsList, ok := object["ConfigRules"].([]interface{}); ok {
			for _, v := range configRuleIdsList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"config_rule_id": m1["ConfigRuleId"],
					}
					configRuleIds = append(configRuleIds, temp1)
				}
			}
		}
		d.Set("config_rule_ids", configRuleIds)
	}

	d.Set("description", object["Description"])
	d.Set("risk_level", formatInt(object["RiskLevel"]))
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudConfigAggregateCompliancePackUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	var response map[string]interface{}
	d.Partial(true)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"CompliancePackId": parts[1],
		"AggregatorId":     parts[0],
	}

	if _, ok := d.GetOk("config_rules"); ok {
		configRulesMaps := make([]map[string]interface{}, 0)
		if d.HasChange("config_rules") {
			update = true
			for _, configRules := range d.Get("config_rules").(*schema.Set).List() {
				configRulesArg := configRules.(map[string]interface{})
				configRulesMap := map[string]interface{}{
					"ConfigRuleName":        configRulesArg["config_rule_name"],
					"ManagedRuleIdentifier": configRulesArg["managed_rule_identifier"],
				}
				configRuleParametersMaps := make([]map[string]interface{}, 0)
				for _, configRuleParameters := range configRulesArg["config_rule_parameters"].(*schema.Set).List() {
					configRuleParametersArg := configRuleParameters.(map[string]interface{})
					configRuleParametersMap := map[string]interface{}{
						"ParameterName":  configRuleParametersArg["parameter_name"],
						"ParameterValue": configRuleParametersArg["parameter_value"],
					}
					configRuleParametersMaps = append(configRuleParametersMaps, configRuleParametersMap)
				}
				configRulesMap["ConfigRuleParameters"] = configRuleParametersMaps
				configRulesMaps = append(configRulesMaps, configRulesMap)
			}
		}

		if v, err := convertArrayObjectToJsonString(configRulesMaps); err == nil {
			request["ConfigRules"] = v
		} else {
			return WrapError(err)
		}
	}

	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if d.HasChange("risk_level") {
		update = true
	}
	request["RiskLevel"] = d.Get("risk_level")

	if d.HasChange("aggregate_compliance_pack_name") {
		update = true
	}
	request["CompliancePackName"] = d.Get("aggregate_compliance_pack_name")
	if update {
		action := "UpdateAggregateCompliancePack"
		conn, err := client.NewConfigClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateAggregateCompliancePack")
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
		stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, configService.ConfigAggregateCompliancePackStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("risk_level")
		d.SetPartial("description")
		d.SetPartial("config_rules")
	}

	if d.HasChange("config_rule_ids") {

		oraw, nraw := d.GetChange("config_rule_ids")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}

		if len(remove) > 0 {
			removeRulesReq := map[string]interface{}{
				"AggregatorId":     parts[0],
				"CompliancePackId": parts[1],
			}

			ruleMaps := make([]interface{}, 0)
			for _, rule := range remove {
				ruleArg := rule.(map[string]interface{})
				ruleMaps = append(ruleMaps, ruleArg["config_rule_id"].(string))
			}
			removeRulesReq["ConfigRuleIds"] = convertListToCommaSeparate(ruleMaps)

			action := "DetachAggregateConfigRuleToCompliancePack"
			conn, err := client.NewConfigClient()
			if err != nil {
				return WrapError(err)
			}
			removeRulesReq["ClientToken"] = buildClientToken("DetachAggregateConfigRuleToCompliancePack")
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, removeRulesReq, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, removeRulesReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		if len(create) > 0 {
			addRulesReq := map[string]interface{}{
				"AggregatorId":     parts[0],
				"CompliancePackId": parts[1],
			}

			ruleMaps := make([]interface{}, 0)
			for _, rule := range create {
				ruleArg := rule.(map[string]interface{})
				ruleMaps = append(ruleMaps, ruleArg["config_rule_id"].(string))
			}
			addRulesReq["ConfigRuleIds"] = convertListToCommaSeparate(ruleMaps)

			action := "AttachAggregateConfigRuleToCompliancePack"
			conn, err := client.NewConfigClient()
			if err != nil {
				return WrapError(err)
			}
			addRulesReq["ClientToken"] = buildClientToken("AttachAggregateConfigRuleToCompliancePack")
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, addRulesReq, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, addRulesReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		d.SetPartial("config_rule_ids")
	}

	d.Partial(false)
	return resourceAlicloudConfigAggregateCompliancePackRead(d, meta)
}
func resourceAlicloudConfigAggregateCompliancePackDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteAggregateCompliancePacks"
	var response map[string]interface{}
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"CompliancePackIds": parts[1],
		"AggregatorId":      parts[0],
	}

	request["ClientToken"] = buildClientToken("DeleteAggregateCompliancePacks")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
		if IsExpectedErrors(err, []string{"Invalid.AggregatorId.Value", "Invalid.CompliancePackId.Value"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

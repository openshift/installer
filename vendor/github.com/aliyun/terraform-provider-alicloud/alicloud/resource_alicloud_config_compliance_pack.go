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

func resourceAlicloudConfigCompliancePack() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudConfigCompliancePackCreate,
		Read:   resourceAlicloudConfigCompliancePackRead,
		Update: resourceAlicloudConfigCompliancePackUpdate,
		Delete: resourceAlicloudConfigCompliancePackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"compliance_pack_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"compliance_pack_template_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config_rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_rule_parameters": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameter_name": {
										Type:     schema.TypeString,
										Required: true,
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

func resourceAlicloudConfigCompliancePackCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	var response map[string]interface{}
	action := "CreateCompliancePack"
	request := make(map[string]interface{})
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	request["CompliancePackName"] = d.Get("compliance_pack_name")
	request["CompliancePackTemplateId"] = d.Get("compliance_pack_template_id")
	configRulesMaps := make([]map[string]interface{}, 0)
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
	if v, err := convertArrayObjectToJsonString(configRulesMaps); err == nil {
		request["ConfigRules"] = v
	} else {
		return WrapError(err)
	}
	request["Description"] = d.Get("description")
	request["RiskLevel"] = d.Get("risk_level")
	request["ClientToken"] = buildClientToken("CreateCompliancePack")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 30*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_compliance_pack", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["CompliancePackId"]))
	stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, configService.ConfigCompliancePackStateRefreshFunc(d.Id(), []string{"CREATING"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudConfigCompliancePackRead(d, meta)
}
func resourceAlicloudConfigCompliancePackRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	object, err := configService.DescribeConfigCompliancePack(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_compliance_pack configService.DescribeConfigCompliancePack Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("compliance_pack_name", object["CompliancePackName"])
	d.Set("compliance_pack_template_id", object["CompliancePackTemplateId"])

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
	if err := d.Set("config_rules", configRules); err != nil {
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("risk_level", formatInt(object["RiskLevel"]))
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudConfigCompliancePackUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"CompliancePackId": d.Id(),
	}
	if d.HasChange("config_rules") {
		update = true
	}
	configRulesMaps := make([]map[string]interface{}, 0)
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
	if v, err := convertArrayObjectToJsonString(configRulesMaps); err == nil {
		request["ConfigRules"] = v
	} else {
		return WrapError(err)
	}
	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if d.HasChange("risk_level") {
		update = true
	}
	request["RiskLevel"] = d.Get("risk_level")
	if update {
		action := "UpdateCompliancePack"
		conn, err := client.NewConfigClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateCompliancePack")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 30*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"CompliancePackAlreadyPending"}) || NeedRetry(err) {
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
		stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, configService.ConfigCompliancePackStateRefreshFunc(d.Id(), []string{"CREATING"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudConfigCompliancePackRead(d, meta)
}
func resourceAlicloudConfigCompliancePackDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCompliancePacks"
	var response map[string]interface{}
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"CompliancePackIds": d.Id(),
	}

	request["ClientToken"] = buildClientToken("DeleteCompliancePacks")
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

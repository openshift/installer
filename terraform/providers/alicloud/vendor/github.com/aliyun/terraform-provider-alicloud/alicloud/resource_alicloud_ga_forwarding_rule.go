package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudGaForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaForwardingRuleCreate,
		Read:   resourceAlicloudGaForwardingRuleRead,
		Update: resourceAlicloudGaForwardingRuleUpdate,
		Delete: resourceAlicloudGaForwardingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"forwarding_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"forwarding_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"forwarding_rule_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_conditions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_condition_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Host", "Path"}, false),
						},
						"path_config": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"host_config": {
							Type:     schema.TypeSet,
							MinItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"rule_actions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"order": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"rule_action_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"forward_group_config": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_group_tuples": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"endpoint_group_id": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"accelerator_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceAlicloudGaForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateForwardingRules"
	request := make(map[string]interface{})
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	request["AcceleratorId"] = d.Get("accelerator_id")
	request["ListenerId"] = d.Get("listener_id")
	forwardingRule := make(map[string]interface{})
	if val, ok := d.GetOk("priority"); ok {
		forwardingRule["Priority"] = val
	}
	ruleConditions := d.Get("rule_conditions").(*schema.Set).List()
	ruleConditionsMap := make([]map[string]interface{}, 0)
	for _, ruleCondition := range ruleConditions {
		ruleCondition := ruleCondition.(map[string]interface{})
		ruleConditionMap := map[string]interface{}{}
		ruleConditionMap["RuleConditionType"] = ruleCondition["rule_condition_type"]
		if len(ruleCondition["path_config"].(*schema.Set).List()) > 0 {
			ruleConditionMap["PathConfig"] = map[string]interface{}{
				"Values": ruleCondition["path_config"].(*schema.Set).List()[0].(map[string]interface{})["values"],
			}
		}
		if len(ruleCondition["host_config"].(*schema.Set).List()) > 0 {
			ruleConditionMap["HostConfig"] = map[string]interface{}{
				"Values": ruleCondition["host_config"].(*schema.Set).List()[0].(map[string]interface{})["values"],
			}
		}
		ruleConditionsMap = append(ruleConditionsMap, ruleConditionMap)
	}
	forwardingRule["RuleConditions"] = ruleConditionsMap
	ruleActions := d.Get("rule_actions").(*schema.Set).List()
	ruleActionsMap := make([]map[string]interface{}, 0)
	for _, ruleAction := range ruleActions {
		ruleAction := ruleAction.(map[string]interface{})
		ruleActionMap := map[string]interface{}{}
		ruleActionMap["Order"] = ruleAction["order"]
		ruleActionMap["RuleActionType"] = ruleAction["rule_action_type"]
		forwardGroupConfigMap := map[string]interface{}{}
		serverGroupTuplesMap := make([]map[string]interface{}, 0)
		for _, serverGroupTuple := range ruleAction["forward_group_config"].(*schema.Set).List()[0].(map[string]interface{})["server_group_tuples"].(*schema.Set).List() {
			serverGroupTuplesMap = append(serverGroupTuplesMap, map[string]interface{}{
				"EndpointGroupId": serverGroupTuple.(map[string]interface{})["endpoint_group_id"],
			})
		}
		forwardGroupConfigMap["ServerGroupTuples"] = serverGroupTuplesMap
		ruleActionMap["ForwardGroupConfig"] = forwardGroupConfigMap
		ruleActionsMap = append(ruleActionsMap, ruleActionMap)
	}
	forwardingRule["RuleActions"] = ruleActionsMap
	if val, ok := d.GetOk("forwarding_rule_name"); ok {
		forwardingRule["ForwardingRuleName"] = val
	}
	request["ForwardingRules"] = []interface{}{forwardingRule}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken(action)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_forwarding_rule", action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ForwardingRules", response)
	if err != nil || len(v.([]interface{})) < 1 {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	response = v.([]interface{})[0].(map[string]interface{})
	d.SetId(fmt.Sprintf("%s:%s:%s", request["AcceleratorId"].(string), request["ListenerId"].(string), fmt.Sprint(response["ForwardingRuleId"])))
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaForwardingRuleStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudGaForwardingRuleRead(d, meta)
}
func resourceAlicloudGaForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	object, err := gaService.DescribeGaForwardingRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_ip_set gaService.DescribeGaForwardingRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("accelerator_id", parts[0])
	d.Set("listener_id", parts[1])
	d.Set("priority", object["Priority"])
	d.Set("forwarding_rule_id", object["ForwardingRuleId"])
	d.Set("forwarding_rule_name", object["ForwardingRuleName"])
	d.Set("forwarding_rule_status", object["ForwardingRuleStatus"])
	ruleConditionsMap := make([]map[string]interface{}, 0)
	for _, ruleCondition := range object["RuleConditions"].([]interface{}) {
		ruleCondition := ruleCondition.(map[string]interface{})
		ruleConditionMap := map[string]interface{}{}
		ruleConditionMap["rule_condition_type"] = ruleCondition["RuleConditionType"]
		if ruleCondition["PathConfig"].(map[string]interface{})["Values"] != nil {
			ruleConditionMap["path_config"] = []map[string]interface{}{
				{
					"values": ruleCondition["PathConfig"].(map[string]interface{})["Values"],
				},
			}
		}
		if ruleCondition["HostConfig"].(map[string]interface{})["Values"] != nil {
			ruleConditionMap["host_config"] = []map[string]interface{}{
				{
					"values": ruleCondition["HostConfig"].(map[string]interface{})["Values"],
				},
			}
		}
		ruleConditionsMap = append(ruleConditionsMap, ruleConditionMap)
	}
	d.Set("rule_conditions", ruleConditionsMap)
	ruleActionsMap := make([]map[string]interface{}, 0)
	for _, ruleAction := range object["RuleActions"].([]interface{}) {
		ruleAction := ruleAction.(map[string]interface{})
		ruleActionMap := map[string]interface{}{}
		ruleActionMap["order"] = ruleAction["Order"]
		ruleActionMap["rule_action_type"] = ruleAction["RuleActionType"]
		serverGroupTuplesMap := make([]map[string]interface{}, 0)
		for _, serverGroupTuple := range ruleAction["ForwardGroupConfig"].(map[string]interface{})["ServerGroupTuples"].([]interface{}) {
			serverGroupTuplesMap = append(serverGroupTuplesMap, map[string]interface{}{
				"endpoint_group_id": serverGroupTuple.(map[string]interface{})["EndpointGroupId"],
			})
		}
		ruleActionMap["forward_group_config"] = []map[string]interface{}{
			{
				"server_group_tuples": serverGroupTuplesMap,
			},
		}
		ruleActionsMap = append(ruleActionsMap, ruleActionMap)
	}
	d.Set("rule_actions", ruleActionsMap)
	return nil
}
func resourceAlicloudGaForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AcceleratorId": parts[0],
		"ListenerId":    parts[1],
	}
	forwardingRule := make(map[string]interface{})
	forwardingRule["ForwardingRuleId"] = parts[2]
	forwardingRule["Priority"] = d.Get("priority")
	ruleConditions := d.Get("rule_conditions").(*schema.Set).List()
	ruleConditionsMap := make([]map[string]interface{}, 0)
	for _, ruleCondition := range ruleConditions {
		ruleCondition := ruleCondition.(map[string]interface{})
		ruleConditionMap := map[string]interface{}{}
		ruleConditionMap["RuleConditionType"] = ruleCondition["rule_condition_type"]
		if len(ruleCondition["path_config"].(*schema.Set).List()) > 0 {
			ruleConditionMap["PathConfig"] = map[string]interface{}{
				"Values": ruleCondition["path_config"].(*schema.Set).List()[0].(map[string]interface{})["values"],
			}
		}
		if len(ruleCondition["host_config"].(*schema.Set).List()) > 0 {
			ruleConditionMap["HostConfig"] = map[string]interface{}{
				"Values": ruleCondition["host_config"].(*schema.Set).List()[0].(map[string]interface{})["values"],
			}
		}
		ruleConditionsMap = append(ruleConditionsMap, ruleConditionMap)
	}
	forwardingRule["RuleConditions"] = ruleConditionsMap
	ruleActions := d.Get("rule_actions").(*schema.Set).List()
	ruleActionsMap := make([]map[string]interface{}, 0)
	for _, ruleAction := range ruleActions {
		ruleAction := ruleAction.(map[string]interface{})
		ruleActionMap := map[string]interface{}{}
		ruleActionMap["Order"] = ruleAction["order"]
		ruleActionMap["RuleActionType"] = ruleAction["rule_action_type"]
		forwardGroupConfigMap := map[string]interface{}{}
		serverGroupTuplesMap := make([]map[string]interface{}, 0)
		for _, serverGroupTuple := range ruleAction["forward_group_config"].(*schema.Set).List()[0].(map[string]interface{})["server_group_tuples"].(*schema.Set).List() {
			serverGroupTuplesMap = append(serverGroupTuplesMap, map[string]interface{}{
				"EndpointGroupId": serverGroupTuple.(map[string]interface{})["endpoint_group_id"],
			})
		}
		forwardGroupConfigMap["ServerGroupTuples"] = serverGroupTuplesMap
		ruleActionMap["ForwardGroupConfig"] = forwardGroupConfigMap
		ruleActionsMap = append(ruleActionsMap, ruleActionMap)
	}
	forwardingRule["RuleActions"] = ruleActionsMap
	if val, ok := d.GetOk("forwarding_rule_name"); ok {
		forwardingRule["ForwardingRuleName"] = val
	}
	request["ForwardingRules"] = []interface{}{forwardingRule}
	request["RegionId"] = client.RegionId
	action := "UpdateForwardingRules"
	request["ClientToken"] = buildClientToken(action)
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken(action)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.ForwardingRule"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaForwardingRuleStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudGaForwardingRuleRead(d, meta)
}
func resourceAlicloudGaForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteForwardingRules"
	var response map[string]interface{}
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AcceleratorId": parts[0],
		"ListenerId":    parts[1],
	}
	request["ForwardingRuleIds"] = []string{parts[2]}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken(action)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.ForwardingRule"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

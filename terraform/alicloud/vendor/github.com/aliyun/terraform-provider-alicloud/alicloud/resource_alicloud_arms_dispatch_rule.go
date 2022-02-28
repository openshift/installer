package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudArmsDispatchRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudArmsDispatchRuleCreate,
		Read:   resourceAlicloudArmsDispatchRuleRead,
		Update: resourceAlicloudArmsDispatchRuleUpdate,
		Delete: resourceAlicloudArmsDispatchRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"is_recover": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"group_rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"group_wait_time": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"group_interval": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"grouping_fields": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"repeat_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("dispatch_type").(string) == "DISCARD_ALERT"
				},
			},
			"dispatch_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CREATE_ALERT", "DISCARD_ALERT"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label_match_expression_grid": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label_match_expression_groups": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"label_match_expressions": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"operator": {
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
			"notify_rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notify_objects": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notify_object_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"notify_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"ARMS_CONTACT", "ARMS_CONTACT_GROUP"}, false),
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"notify_channels": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("dispatch_type").(string) == "DISCARD_ALERT"
				},
			},
			"dispatch_rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudArmsDispatchRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDispatchRule"
	request := make(map[string]interface{})
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	dispatchRuleMap := make(map[string]interface{}, 0)
	if v, ok := d.GetOk("is_recover"); ok {
		dispatchRuleMap["isRecover"] = v
	}

	if v, ok := d.GetOk("group_rules"); ok {
		groupRulesMaps := make([]map[string]interface{}, 0)
		for _, groupRules := range v.(*schema.Set).List() {
			groupRulesArg := groupRules.(map[string]interface{})
			groupRulesMap := map[string]interface{}{
				"groupWait":      groupRulesArg["group_wait_time"],
				"groupInterval":  groupRulesArg["group_interval"],
				"groupingFields": groupRulesArg["grouping_fields"],
				"repeatInterval": groupRulesArg["repeat_interval"],
			}
			groupRulesMaps = append(groupRulesMaps, groupRulesMap)
		}
		dispatchRuleMap["groupRules"] = groupRulesMaps
	}
	if v, ok := d.GetOk("dispatch_type"); ok {
		dispatchRuleMap["dispatchType"] = v
	}
	if v, ok := d.GetOk("label_match_expression_grid"); ok {
		labelMatchExpressionGrid := v.(*schema.Set).List()[0]
		labelMatchExpressionGridArg := labelMatchExpressionGrid.(map[string]interface{})
		labelMatchExpressionGroupsMaps := make([]map[string]interface{}, 0)
		for _, labelMatchExpressionGroups := range labelMatchExpressionGridArg["label_match_expression_groups"].(*schema.Set).List() {
			labelMatchExpressionGroupsArg := labelMatchExpressionGroups.(map[string]interface{})
			labelMatchExpressionsMaps := make([]map[string]interface{}, 0)
			for _, labelMatchExpressions := range labelMatchExpressionGroupsArg["label_match_expressions"].(*schema.Set).List() {
				labelMatchExpressionsArg := labelMatchExpressions.(map[string]interface{})
				labelMatchExpressionsMap := map[string]interface{}{
					"key":      labelMatchExpressionsArg["key"],
					"value":    labelMatchExpressionsArg["value"],
					"operator": labelMatchExpressionsArg["operator"],
				}
				labelMatchExpressionsMaps = append(labelMatchExpressionsMaps, labelMatchExpressionsMap)
			}
			labelMatchExpressionGroupsMaps = append(labelMatchExpressionGroupsMaps, map[string]interface{}{
				"labelMatchExpressions": labelMatchExpressionsMaps,
			})
		}
		dispatchRuleMap["labelMatchExpressionGrid"] = map[string]interface{}{
			"labelMatchExpressionGroups": labelMatchExpressionGroupsMaps,
		}
	}

	if v, ok := d.GetOk("notify_rules"); ok {
		notifyRulesMaps := make([]map[string]interface{}, 0)
		for _, notifyRules := range v.(*schema.Set).List() {
			notifyRulesArg := notifyRules.(map[string]interface{})
			notifyObjectsMaps := make([]map[string]interface{}, 0)
			for _, notifyObjects := range notifyRulesArg["notify_objects"].(*schema.Set).List() {
				notifyObjectsArg := notifyObjects.(map[string]interface{})
				notifyObjectsMap := map[string]interface{}{
					"notifyType":     notifyObjectsArg["notify_type"],
					"name":           notifyObjectsArg["name"],
					"notifyObjectId": notifyObjectsArg["notify_object_id"],
				}
				notifyObjectsMaps = append(notifyObjectsMaps, notifyObjectsMap)
			}
			notifyRulesMap := map[string]interface{}{
				"notifyObjects":  notifyObjectsMaps,
				"notifyChannels": notifyRulesArg["notify_channels"].([]interface{}),
			}
			notifyRulesMaps = append(notifyRulesMaps, notifyRulesMap)
		}
		dispatchRuleMap["notifyRules"] = notifyRulesMaps
	}

	if v, ok := d.GetOk("dispatch_rule_name"); ok {
		dispatchRuleMap["name"] = v
	}

	request["RegionId"] = client.RegionId
	if v, err := convertMaptoJsonString(dispatchRuleMap); err != nil {
		return WrapError(err)
	} else {
		request["DispatchRule"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_dispatch_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DispatchRuleId"]))

	return resourceAlicloudArmsDispatchRuleRead(d, meta)
}
func resourceAlicloudArmsDispatchRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := ArmsService{client}
	object, err := armsService.DescribeArmsDispatchRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_dispatch_rule armsService.DescribeArmsDispatchRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if groupRulesList, ok := object["GroupRules"]; ok && groupRulesList != nil {
		groupRulesMaps := make([]map[string]interface{}, 0)
		for _, groupRulesListItem := range groupRulesList.([]interface{}) {
			if groupRulesItemMap, ok := groupRulesListItem.(map[string]interface{}); ok {
				groupRulesMap := make(map[string]interface{}, 0)
				groupRulesMap["group_interval"] = groupRulesItemMap["GroupInterval"]
				groupRulesMap["group_wait_time"] = groupRulesItemMap["GroupWaitTime"]
				groupRulesMap["group_id"] = groupRulesItemMap["GroupId"]
				groupRulesMap["grouping_fields"] = groupRulesItemMap["GroupingFields"]
				groupRulesMap["repeat_interval"] = groupRulesItemMap["RepeatInterval"]
				groupRulesMaps = append(groupRulesMaps, groupRulesMap)
			}
		}
		d.Set("group_rules", groupRulesMaps)
	}
	d.Set("status", object["State"])
	if labelMatchExpressionGrid, ok := object["LabelMatchExpressionGrid"]; ok && labelMatchExpressionGrid != nil {
		labelMatchExpressionGridMaps := make([]map[string]interface{}, 0)
		labelMatchExpressionGridMap := make(map[string]interface{})

		labelMatchExpressionGroupsMaps := make([]map[string]interface{}, 0)
		if v, ok := labelMatchExpressionGrid.(map[string]interface{})["LabelMatchExpressionGroups"]; ok && v != nil {
			for _, labelMatchExpressionGroups := range v.([]interface{}) {
				labelMatchExpressionGroupsMap := make(map[string]interface{})
				labelMatchExpressionsMaps := make([]map[string]interface{}, 0)
				if v, ok := labelMatchExpressionGroups.(map[string]interface{})["LabelMatchExpressions"]; ok && v != nil {
					for _, labelMatchExpressions := range v.([]interface{}) {
						labelMatchExpressionsArg := labelMatchExpressions.(map[string]interface{})
						labelMatchExpressionsMap := make(map[string]interface{}, 0)
						labelMatchExpressionsMap["operator"] = labelMatchExpressionsArg["Operator"]
						labelMatchExpressionsMap["key"] = labelMatchExpressionsArg["Key"]
						labelMatchExpressionsMap["value"] = labelMatchExpressionsArg["Value"]
						labelMatchExpressionsMaps = append(labelMatchExpressionsMaps, labelMatchExpressionsMap)
					}
				}
				labelMatchExpressionGroupsMap["label_match_expressions"] = labelMatchExpressionsMaps
				labelMatchExpressionGroupsMaps = append(labelMatchExpressionGroupsMaps, labelMatchExpressionGroupsMap)
			}
		}
		labelMatchExpressionGridMap["label_match_expression_groups"] = labelMatchExpressionGroupsMaps
		labelMatchExpressionGridMaps = append(labelMatchExpressionGridMaps, labelMatchExpressionGridMap)
		if err := d.Set("label_match_expression_grid", labelMatchExpressionGridMaps); err != nil {
			return WrapError(err)
		}
	}

	if notifyRulesList, ok := object["NotifyRules"]; ok && notifyRulesList != nil {
		notifyRulesMaps := make([]map[string]interface{}, 0)
		for _, notifyRulesListItem := range notifyRulesList.([]interface{}) {
			if notifyRulesItemMap, ok := notifyRulesListItem.(map[string]interface{}); ok {
				notifyRulesMap := make(map[string]interface{}, 0)
				notifyObjectsMaps := make([]map[string]interface{}, 0)
				for _, notifyObjects := range notifyRulesItemMap["NotifyObjects"].([]interface{}) {
					notifyObjectsArg := notifyObjects.(map[string]interface{})
					notifyObjectsMap := make(map[string]interface{}, 0)
					notifyObjectsMap["notify_type"] = convertArmsDispatchRuleNotifyTypeResponse(notifyObjectsArg["NotifyType"])
					notifyObjectsMap["notify_object_id"] = notifyObjectsArg["NotifyObjectId"]
					notifyObjectsMap["name"] = notifyObjectsArg["Name"]
					notifyObjectsMaps = append(notifyObjectsMaps, notifyObjectsMap)
				}
				notifyRulesMap["notify_objects"] = notifyObjectsMaps
				notifyRulesMap["notify_channels"] = notifyRulesItemMap["NotifyChannels"]
				notifyRulesMaps = append(notifyRulesMaps, notifyRulesMap)
			}
		}
		d.Set("notify_rules", notifyRulesMaps)
	}

	d.Set("dispatch_rule_name", object["Name"])
	return nil
}
func resourceAlicloudArmsDispatchRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "UpdateDispatchRule"
	request := make(map[string]interface{})
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	dispatchRuleMap := make(map[string]interface{}, 0)
	dispatchRuleMap["isRecover"] = d.Get("is_recover")

	dispatchRuleMap["ruleid"] = d.Id()
	dispatchRuleMap["id"] = d.Id()
	if _, ok := d.GetOk("group_rules"); ok {
		oraw, nraw := d.GetChange("group_rules")
		groupRulesMaps := make([]map[string]interface{}, 0)
		groupId := 0
		for _, groupRules := range oraw.(*schema.Set).List() {
			groupRulesArg := groupRules.(map[string]interface{})
			groupId = groupRulesArg["group_id"].(int)
		}
		for _, groupRules := range nraw.(*schema.Set).List() {
			groupRulesArg := groupRules.(map[string]interface{})
			groupRulesMap := map[string]interface{}{
				"groupId":        groupId,
				"groupWait":      groupRulesArg["group_wait_time"],
				"groupInterval":  groupRulesArg["group_interval"],
				"groupingFields": groupRulesArg["grouping_fields"],
				"repeatInterval": groupRulesArg["repeat_interval"],
			}
			groupRulesMaps = append(groupRulesMaps, groupRulesMap)
		}
		dispatchRuleMap["groupRules"] = groupRulesMaps
	}
	if v, ok := d.GetOk("dispatch_type"); ok {
		dispatchRuleMap["dispatchType"] = v
	}
	if v, ok := d.GetOk("label_match_expression_grid"); ok {
		labelMatchExpressionGrid := v.(*schema.Set).List()[0]
		labelMatchExpressionGridArg := labelMatchExpressionGrid.(map[string]interface{})
		labelMatchExpressionGroupsMaps := make([]map[string]interface{}, 0)
		for _, labelMatchExpressionGroups := range labelMatchExpressionGridArg["label_match_expression_groups"].(*schema.Set).List() {
			labelMatchExpressionGroupsArg := labelMatchExpressionGroups.(map[string]interface{})
			labelMatchExpressionsMaps := make([]map[string]interface{}, 0)
			for _, labelMatchExpressions := range labelMatchExpressionGroupsArg["label_match_expressions"].(*schema.Set).List() {
				labelMatchExpressionsArg := labelMatchExpressions.(map[string]interface{})
				labelMatchExpressionsMap := map[string]interface{}{
					"key":      labelMatchExpressionsArg["key"],
					"value":    labelMatchExpressionsArg["value"],
					"operator": labelMatchExpressionsArg["operator"],
				}
				labelMatchExpressionsMaps = append(labelMatchExpressionsMaps, labelMatchExpressionsMap)
			}
			labelMatchExpressionGroupsMaps = append(labelMatchExpressionGroupsMaps, map[string]interface{}{
				"labelMatchExpressions": labelMatchExpressionsMaps,
			})
		}
		dispatchRuleMap["labelMatchExpressionGrid"] = map[string]interface{}{
			"labelMatchExpressionGroups": labelMatchExpressionGroupsMaps,
		}
	}

	if v, ok := d.GetOk("notify_rules"); ok {
		notifyRulesMaps := make([]map[string]interface{}, 0)
		for _, notifyRules := range v.(*schema.Set).List() {
			notifyRulesArg := notifyRules.(map[string]interface{})
			notifyObjectsMaps := make([]map[string]interface{}, 0)
			for _, notifyObjects := range notifyRulesArg["notify_objects"].(*schema.Set).List() {
				notifyObjectsArg := notifyObjects.(map[string]interface{})
				notifyObjectsMap := map[string]interface{}{
					"notifyType":     notifyObjectsArg["notify_type"],
					"name":           notifyObjectsArg["name"],
					"notifyObjectId": notifyObjectsArg["notify_object_id"],
				}
				notifyObjectsMaps = append(notifyObjectsMaps, notifyObjectsMap)
			}
			notifyRulesMap := map[string]interface{}{
				"notifyObjects":  notifyObjectsMaps,
				"notifyChannels": notifyRulesArg["notify_channels"].([]interface{}),
			}
			notifyRulesMaps = append(notifyRulesMaps, notifyRulesMap)
		}
		dispatchRuleMap["notifyRules"] = notifyRulesMaps
	}

	if v, ok := d.GetOk("dispatch_rule_name"); ok {
		dispatchRuleMap["name"] = v
	}

	request["RegionId"] = client.RegionId
	if v, err := convertMaptoJsonString(dispatchRuleMap); err != nil {
		return WrapError(err)
	} else {
		request["DispatchRule"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_dispatch_rule", action, AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudArmsDispatchRuleRead(d, meta)
}
func resourceAlicloudArmsDispatchRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDispatchRule"
	var response map[string]interface{}
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Id":       d.Id(),
		"RegionId": client.RegionId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
func convertArmsDispatchRuleNotifyTypeResponse(source interface{}) interface{} {
	switch source {
	case "CONTACT":
		return "ARMS_CONTACT"
	case "CONTACT_GROUP":
		return "ARMS_CONTACT_GROUP"
	}
	return source
}

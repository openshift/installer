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

func resourceAlicloudPvtzRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPvtzRuleCreate,
		Read:   resourceAlicloudPvtzRuleRead,
		Update: resourceAlicloudPvtzRuleUpdate,
		Delete: resourceAlicloudPvtzRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"forward_ips": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"OUTBOUND"}, false),
			},
			"zone_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudPvtzRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddResolverRule"
	request := make(map[string]interface{})
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	request["EndpointId"] = d.Get("endpoint_id")
	request["Lang"] = "en"
	request["Name"] = d.Get("rule_name")
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	for k, forwardConfig := range d.Get("forward_ips").(*schema.Set).List() {
		forwardConfigArg := forwardConfig.(map[string]interface{})
		request[fmt.Sprintf("ForwardIp.%d.Port", k+1)] = forwardConfigArg["port"]
		request[fmt.Sprintf("ForwardIp.%d.Ip", k+1)] = forwardConfigArg["ip"]
	}

	request["ZoneName"] = d.Get("zone_name")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pvtz_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RuleId"]))

	return resourceAlicloudPvtzRuleRead(d, meta)
}
func resourceAlicloudPvtzRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	object, err := pvtzService.DescribePvtzRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pvtz_rule pvtzService.DescribePvtzRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("endpoint_id", object["EndpointId"])
	d.Set("rule_name", object["Name"])
	forwardConfigsSli := make([]map[string]interface{}, 0)
	if forwardConfigs, ok := object["ForwardIps"].([]interface{}); ok {
		for _, forwardConfigArgs := range forwardConfigs {
			forwardConfigArg := forwardConfigArgs.(map[string]interface{})
			forwardConfigsMap := make(map[string]interface{})
			forwardConfigsMap["ip"] = forwardConfigArg["Ip"]
			forwardConfigsMap["port"] = formatInt(forwardConfigArg["Port"])
			forwardConfigsSli = append(forwardConfigsSli, forwardConfigsMap)
		}
	}
	d.Set("forward_ips", forwardConfigsSli)
	d.Set("type", object["Type"])
	d.Set("zone_name", object["ZoneName"])
	return nil
}
func resourceAlicloudPvtzRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"RuleId": d.Id(),
	}

	if d.HasChange("rule_name") {
		update = true
		request["Name"] = d.Get("rule_name")
	}

	if d.HasChange("forward_ips") {
		update = true
		for k, forwardConfig := range d.Get("forward_ips").(*schema.Set).List() {
			forwardConfigArg := forwardConfig.(map[string]interface{})
			request[fmt.Sprintf("ForwardIp.%d.Port", k+1)] = forwardConfigArg["port"]
			request[fmt.Sprintf("ForwardIp.%d.Ip", k+1)] = forwardConfigArg["ip"]
		}
	}

	if update {
		request["Lang"] = "en"
		action := "UpdateResolverRule"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudPvtzRuleRead(d, meta)
}
func resourceAlicloudPvtzRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteResolverRule"
	var response map[string]interface{}
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RuleId": d.Id(),
	}

	request["Lang"] = "en"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"ResolverRule.NotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

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

func resourceAlicloudAlidnsMonitorConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsMonitorConfigCreate,
		Read:   resourceAlicloudAlidnsMonitorConfigRead,
		Update: resourceAlicloudAlidnsMonitorConfigUpdate,
		Delete: resourceAlicloudAlidnsMonitorConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"addr_pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"evaluation_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3}),
			},
			"interval": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{60}),
			},
			"isp_city_node": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"isp_code": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"monitor_extend_info": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"protocol_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS", "PING", "TCP"}, false),
			},
			"timeout": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{2000, 3000, 5000, 10000}),
			},
		},
	}
}

func resourceAlicloudAlidnsMonitorConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddDnsGtmMonitor"
	request := make(map[string]interface{})
	conn, err := client.NewAlidnsClient()
	if err != nil {
		return WrapError(err)
	}
	request["AddrPoolId"] = d.Get("addr_pool_id")
	request["EvaluationCount"] = d.Get("evaluation_count")
	request["Interval"] = d.Get("interval")
	if v, ok := d.GetOk("isp_city_node"); ok {
		for i, ispCityNode := range v.(*schema.Set).List() {
			ispCityNodeArg := ispCityNode.(map[string]interface{})
			request[fmt.Sprintf("IspCityNode.%d.CityCode", i+1)] = ispCityNodeArg["city_code"]
			request[fmt.Sprintf("IspCityNode.%d.IspCode", i+1)] = ispCityNodeArg["isp_code"]
		}
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	request["MonitorExtendInfo"] = d.Get("monitor_extend_info")
	request["ProtocolType"] = d.Get("protocol_type")
	request["Timeout"] = d.Get("timeout")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_monitor_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["MonitorConfigId"]))
	return resourceAlicloudAlidnsMonitorConfigRead(d, meta)
}
func resourceAlicloudAlidnsMonitorConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	object, err := alidnsService.DescribeAlidnsMonitorConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_monitor_config alidnsService.DescribeAlidnsMonitorConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := object["EvaluationCount"]; ok {
		d.Set("evaluation_count", formatInt(v))
	}
	if v, ok := object["Interval"]; ok {
		d.Set("interval", formatInt(v))
	}
	d.Set("monitor_extend_info", object["MonitorExtendInfo"])
	d.Set("protocol_type", object["ProtocolType"])
	if v, ok := object["Timeout"]; ok {
		d.Set("timeout", formatInt(v))
	}
	if ispCityNodesList, ok := object["IspCityNodes"]; ok {
		ispCityNodesArg := ispCityNodesList.(map[string]interface{})
		if ispCityNodeConfig, ok := ispCityNodesArg["IspCityNode"]; ok {
			ispCityNodeConfigArgs := ispCityNodeConfig.([]interface{})
			ispCityNodesMaps := make([]map[string]interface{}, 0)
			for _, ispCityNodeMapArgitem := range ispCityNodeConfigArgs {
				ispCityNodeMapArg := ispCityNodeMapArgitem.(map[string]interface{})
				ispCityNodesMap := map[string]interface{}{}
				ispCityNodesMap["city_code"] = ispCityNodeMapArg["CityCode"]
				ispCityNodesMap["isp_code"] = ispCityNodeMapArg["IspCode"]
				ispCityNodesMaps = append(ispCityNodesMaps, ispCityNodesMap)
			}
			d.Set("isp_city_node", ispCityNodesMaps)
		}
	}
	return nil
}
func resourceAlicloudAlidnsMonitorConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"MonitorConfigId": d.Id(),
	}
	if d.HasChange("isp_city_node") {
		update = true
	}
	if v, ok := d.GetOk("isp_city_node"); ok {
		for i, ispCityNode := range v.(*schema.Set).List() {
			ispCityNodeArg := ispCityNode.(map[string]interface{})
			request[fmt.Sprintf("IspCityNode.%d.CityCode", i+1)] = ispCityNodeArg["city_code"]
			request[fmt.Sprintf("IspCityNode.%d.IspCode", i+1)] = ispCityNodeArg["isp_code"]
		}
	}
	if d.HasChange("monitor_extend_info") {
		update = true
	}
	request["MonitorExtendInfo"] = d.Get("monitor_extend_info")
	if d.HasChange("protocol_type") {
		update = true
	}
	request["ProtocolType"] = d.Get("protocol_type")
	if d.HasChange("evaluation_count") {
		update = true
		request["EvaluationCount"] = d.Get("evaluation_count")
	}
	if d.HasChange("interval") {
		update = true
		request["Interval"] = d.Get("interval")
	}
	if d.HasChange("timeout") {
		update = true
		request["Timeout"] = d.Get("timeout")
	}
	if update {
		action := "UpdateDnsGtmMonitor"
		conn, err := client.NewAlidnsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudAlidnsMonitorConfigRead(d, meta)
}
func resourceAlicloudAlidnsMonitorConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudAlidnsMonitorConfig. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

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

func resourceAlicloudAlidnsAccessStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsAccessStrategyCreate,
		Read:   resourceAlicloudAlidnsAccessStrategyRead,
		Update: resourceAlicloudAlidnsAccessStrategyUpdate,
		Delete: resourceAlicloudAlidnsAccessStrategyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"access_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AUTO", "DEFAULT", "FAILOVER"}, false),
			},
			"default_addr_pool_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPV4", "IPV6", "DOMAIN"}, false),
			},
			"lines": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"line_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"default_addr_pools": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addr_pool_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"lba_weight": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"default_latency_optimization": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"OPEN", "CLOSE"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("strategy_mode"); ok && v.(string) == "LATENCY" {
						return false
					}
					return true
				},
			},
			"default_lba_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL_RR", "RATIO"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("strategy_mode"); ok && v.(string) == "GEO" {
						return false
					}
					return true
				},
			},
			"default_max_return_addr_num": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("strategy_mode"); ok && v.(string) == "LATENCY" {
						return false
					}
					return true
				},
			},
			"default_min_available_addr_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"failover_addr_pool_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"failover_addr_pools": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addr_pool_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"lba_weight": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"failover_latency_optimization": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"OPEN", "CLOSE"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("strategy_mode"); ok && v.(string) == "LATENCY" {
						return false
					}
					return true
				},
			},
			"failover_lba_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL_RR", "RATIO"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("strategy_mode"); ok && v.(string) == "GEO" {
						return false
					}
					return true
				},
			},
			"failover_max_return_addr_num": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("strategy_mode"); ok && v.(string) == "LATENCY" {
						return false
					}
					return true
				},
			},
			"failover_min_available_addr_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"strategy_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"GEO", "LATENCY"}, false),
			},
			"strategy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudAlidnsAccessStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddDnsGtmAccessStrategy"
	request := make(map[string]interface{})
	conn, err := client.NewAlidnsClient()
	if err != nil {
		return WrapError(err)
	}

	request["DefaultAddrPoolType"] = d.Get("default_addr_pool_type")
	if v, ok := d.GetOk("default_latency_optimization"); ok {
		request["DefaultLatencyOptimization"] = v
	}
	if v, ok := d.GetOk("default_lba_strategy"); ok {
		request["DefaultLbaStrategy"] = v
	}
	if v, ok := d.GetOk("default_max_return_addr_num"); ok {
		request["DefaultMaxReturnAddrNum"] = v
	}
	request["DefaultMinAvailableAddrNum"] = d.Get("default_min_available_addr_num")
	if v, ok := d.GetOk("lines"); ok {
		lines := make([]interface{}, 0)
		for _, line := range v.(*schema.Set).List() {
			lineArg := line.(map[string]interface{})
			if v, ok := lineArg["line_code"]; ok && fmt.Sprint(v) != "" {
				lines = append(lines, v)
			}
		}
		request["Lines"] = convertListToJsonString(lines)
	}
	if v, ok := d.GetOk("failover_addr_pools"); ok {
		failoverAddrPoolMaps := make([]map[string]interface{}, 0)
		for _, failoverAddrPool := range v.(*schema.Set).List() {
			failoverAddrPoolArg := failoverAddrPool.(map[string]interface{})
			failoverAddrPoolMap := map[string]interface{}{}
			failoverAddrPoolMap["Id"] = failoverAddrPoolArg["addr_pool_id"]
			if v, ok := failoverAddrPoolArg["lba_weight"]; ok && fmt.Sprint(v) != "0" {
				failoverAddrPoolMap["LbaWeight"] = v
			}
			failoverAddrPoolMaps = append(failoverAddrPoolMaps, failoverAddrPoolMap)
		}
		request["FailoverAddrPool"] = failoverAddrPoolMaps
	}
	if v, ok := d.GetOk("default_addr_pools"); ok {
		defaultAddrPoolMaps := make([]map[string]interface{}, 0)
		for _, defaultAddrPool := range v.(*schema.Set).List() {
			defaultAddrPoolArg := defaultAddrPool.(map[string]interface{})
			defaultAddrPoolMap := map[string]interface{}{}
			defaultAddrPoolMap["Id"] = defaultAddrPoolArg["addr_pool_id"]
			if v, ok := defaultAddrPoolArg["lba_weight"]; ok && fmt.Sprint(v) != "0" {
				defaultAddrPoolMap["LbaWeight"] = v
			}
			defaultAddrPoolMaps = append(defaultAddrPoolMaps, defaultAddrPoolMap)
		}
		request["DefaultAddrPool"] = defaultAddrPoolMaps
	}
	if v, ok := d.GetOk("failover_addr_pool_type"); ok {
		request["FailoverAddrPoolType"] = v
	}
	if v, ok := d.GetOk("failover_latency_optimization"); ok {
		request["FailoverLatencyOptimization"] = v
	}
	if v, ok := d.GetOk("failover_lba_strategy"); ok {
		request["FailoverLbaStrategy"] = v
	}
	if v, ok := d.GetOk("failover_max_return_addr_num"); ok {
		request["FailoverMaxReturnAddrNum"] = v
	}
	if v, ok := d.GetOk("failover_min_available_addr_num"); ok {
		request["FailoverMinAvailableAddrNum"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	request["StrategyMode"] = d.Get("strategy_mode")
	request["StrategyName"] = d.Get("strategy_name")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_access_strategy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["StrategyId"]))

	return resourceAlicloudAlidnsAccessStrategyRead(d, meta)
}
func resourceAlicloudAlidnsAccessStrategyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	object, err := alidnsService.DescribeAlidnsAccessStrategy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_access_strategy alidnsService.DescribeAlidnsAccessStrategy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("default_addr_pool_type", object["DefaultAddrPoolType"])
	if defaultAddrPoolsMap, ok := object["DefaultAddrPools"].(map[string]interface{}); ok && defaultAddrPoolsMap != nil {
		if defaultAddrPoolList, ok := defaultAddrPoolsMap["DefaultAddrPool"]; ok && defaultAddrPoolList != nil {
			defaultAddrPoolsMaps := make([]map[string]interface{}, 0)
			for _, defaultAddrPoolListItem := range defaultAddrPoolList.([]interface{}) {
				if v, ok := defaultAddrPoolListItem.(map[string]interface{}); ok {
					defaultAddrPoolListItemMap := make(map[string]interface{}, 0)
					defaultAddrPoolListItemMap["addr_pool_id"] = v["Id"]
					defaultAddrPoolListItemMap["lba_weight"] = v["LbaWeight"]
					defaultAddrPoolsMaps = append(defaultAddrPoolsMaps, defaultAddrPoolListItemMap)
				}
			}
			d.Set("default_addr_pools", defaultAddrPoolsMaps)
		}
	}

	d.Set("default_latency_optimization", object["DefaultLatencyOptimization"])
	d.Set("default_lba_strategy", object["DefaultLbaStrategy"])
	if v, ok := object["DefaultMaxReturnAddrNum"]; ok {
		d.Set("default_max_return_addr_num", formatInt(v))
	}
	if v, ok := object["DefaultMinAvailableAddrNum"]; ok {
		d.Set("default_min_available_addr_num", formatInt(v))
	}
	d.Set("failover_addr_pool_type", object["FailoverAddrPoolType"])
	if failoverAddrPoolsMap, ok := object["FailoverAddrPools"].(map[string]interface{}); ok && failoverAddrPoolsMap != nil {
		if failoverAddrPoolList, ok := failoverAddrPoolsMap["FailoverAddrPool"]; ok && failoverAddrPoolList != nil {
			failoverAddrPoolsMaps := make([]map[string]interface{}, 0)
			for _, failoverAddrPoolListItem := range failoverAddrPoolList.([]interface{}) {
				if v, ok := failoverAddrPoolListItem.(map[string]interface{}); ok {
					failoverAddrPoolListItemMap := make(map[string]interface{}, 0)
					failoverAddrPoolListItemMap["addr_pool_id"] = v["Id"]
					failoverAddrPoolListItemMap["lba_weight"] = v["LbaWeight"]
					failoverAddrPoolsMaps = append(failoverAddrPoolsMaps, failoverAddrPoolListItemMap)
				}
			}
			d.Set("failover_addr_pools", failoverAddrPoolsMaps)
		}
	}
	if LinesMap, ok := object["Lines"].(map[string]interface{}); ok && LinesMap != nil {
		if LineList, ok := LinesMap["Line"]; ok && LineList != nil {
			LinesMaps := make([]map[string]interface{}, 0)
			for _, lineListItem := range LineList.([]interface{}) {
				if v, ok := lineListItem.(map[string]interface{}); ok {
					lineListItemMap := make(map[string]interface{}, 0)
					lineListItemMap["line_code"] = v["LineCode"]
					LinesMaps = append(LinesMaps, lineListItemMap)
				}
			}
			d.Set("lines", LinesMaps)
		}
	}

	d.Set("failover_latency_optimization", object["FailoverLatencyOptimization"])
	d.Set("failover_lba_strategy", object["FailoverLbaStrategy"])
	if v, ok := object["FailoverMaxReturnAddrNum"]; ok {
		d.Set("failover_max_return_addr_num", formatInt(v))
	}
	if v, ok := object["FailoverMinAvailableAddrNum"]; ok {
		d.Set("failover_min_available_addr_num", formatInt(v))
	}
	d.Set("instance_id", object["InstanceId"])
	d.Set("strategy_mode", object["StrategyMode"])
	d.Set("strategy_name", object["StrategyName"])
	d.Set("access_mode", object["AccessMode"])
	return nil
}
func resourceAlicloudAlidnsAccessStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}

	update := false
	request := map[string]interface{}{
		"StrategyId": d.Id(),
	}
	if d.HasChange("access_mode") {
		update = true
		if v, ok := d.GetOk("access_mode"); ok {
			request["AccessMode"] = v
		}
	}
	if d.HasChange("default_addr_pool_type") {
		update = true
	}
	request["DefaultAddrPoolType"] = d.Get("default_addr_pool_type")
	if d.HasChange("lines") {
		update = true
	}
	if v, ok := d.GetOk("lines"); ok {
		lines := make([]interface{}, 0)
		for _, line := range v.(*schema.Set).List() {
			lineArg := line.(map[string]interface{})
			if v, ok := lineArg["line_code"]; ok && fmt.Sprint(v) != "" {
				lines = append(lines, v)
			}
		}
		request["Lines"] = convertListToJsonString(lines)
	}

	if d.HasChange("default_min_available_addr_num") {
		update = true
	}
	request["DefaultMinAvailableAddrNum"] = d.Get("default_min_available_addr_num")

	if d.HasChange("strategy_name") {
		update = true
	}
	request["StrategyName"] = d.Get("strategy_name")
	if d.HasChange("default_latency_optimization") {
		update = true
	}
	if v, ok := d.GetOk("default_latency_optimization"); ok {
		request["DefaultLatencyOptimization"] = v
	}
	if d.HasChange("default_lba_strategy") {
		update = true
	}
	if v, ok := d.GetOk("default_lba_strategy"); ok {
		request["DefaultLbaStrategy"] = v
	}
	if d.HasChange("default_max_return_addr_num") {
		update = true
	}
	if v, ok := d.GetOk("default_max_return_addr_num"); ok {
		request["DefaultMaxReturnAddrNum"] = v
	}
	if d.HasChange("failover_addr_pool_type") {
		update = true
	}
	if v, ok := d.GetOk("failover_addr_pool_type"); ok {
		request["FailoverAddrPoolType"] = v
	}
	if d.HasChange("failover_latency_optimization") {
		update = true
	}
	if v, ok := d.GetOk("failover_latency_optimization"); ok {
		request["FailoverLatencyOptimization"] = v
	}
	if d.HasChange("failover_lba_strategy") {
		update = true
	}
	if v, ok := d.GetOk("failover_lba_strategy"); ok {
		request["FailoverLbaStrategy"] = v
	}
	if d.HasChange("failover_max_return_addr_num") {
		update = true
	}
	if v, ok := d.GetOk("failover_max_return_addr_num"); ok {
		request["FailoverMaxReturnAddrNum"] = v
	}
	if d.HasChange("failover_min_available_addr_num") {
		update = true
	}
	if v, ok := d.GetOk("failover_min_available_addr_num"); ok {
		request["FailoverMinAvailableAddrNum"] = v
	}
	if d.HasChange("default_addr_pools") {
		update = true
	}
	if v, ok := d.GetOk("default_addr_pools"); ok {
		defaultAddrPoolMaps := make([]map[string]interface{}, 0)
		for _, defaultAddrPool := range v.(*schema.Set).List() {
			defaultAddrPoolArg := defaultAddrPool.(map[string]interface{})
			defaultAddrPoolMap := map[string]interface{}{}
			defaultAddrPoolMap["Id"] = defaultAddrPoolArg["addr_pool_id"]
			if v, ok := defaultAddrPoolArg["lba_weight"]; ok && fmt.Sprint(v) != "0" {
				defaultAddrPoolMap["LbaWeight"] = v
			}
			defaultAddrPoolMaps = append(defaultAddrPoolMaps, defaultAddrPoolMap)
		}
		request["DefaultAddrPool"] = defaultAddrPoolMaps
	}
	if d.HasChange("failover_addr_pools") {
		update = true
	}
	if v, ok := d.GetOk("failover_addr_pools"); ok {
		failoverAddrPoolMaps := make([]map[string]interface{}, 0)
		for _, failoverAddrPool := range v.(*schema.Set).List() {
			failoverAddrPoolArg := failoverAddrPool.(map[string]interface{})
			failoverAddrPoolMap := map[string]interface{}{}
			failoverAddrPoolMap["Id"] = failoverAddrPoolArg["addr_pool_id"]
			if v, ok := failoverAddrPoolArg["lba_weight"]; ok && fmt.Sprint(v) != "0" {
				failoverAddrPoolMap["LbaWeight"] = v
			}
			failoverAddrPoolMaps = append(failoverAddrPoolMaps, failoverAddrPoolMap)
		}
		request["FailoverAddrPool"] = failoverAddrPoolMaps
	}
	if update {
		if v, ok := d.GetOk("lang"); ok {
			request["Lang"] = v
		}
		action := "UpdateDnsGtmAccessStrategy"
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
	return resourceAlicloudAlidnsAccessStrategyRead(d, meta)
}
func resourceAlicloudAlidnsAccessStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDnsGtmAccessStrategy"
	var response map[string]interface{}
	conn, err := client.NewAlidnsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"StrategyId": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
	return nil
}

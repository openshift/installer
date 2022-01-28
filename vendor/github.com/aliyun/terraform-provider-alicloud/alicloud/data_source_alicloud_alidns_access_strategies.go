package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudAlidnsAccessStrategies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlidnsAccessStrategiesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"strategy_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"GEO", "LATENCY"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"strategies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_strategy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_addr_pool_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_addr_pools": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addr_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"addr_pool_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"lba_weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"default_available_addr_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_latency_optimization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_lba_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_max_return_addr_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_min_available_addr_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"effective_addr_pool_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failover_addr_pool_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failover_addr_pools": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addr_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"addr_pool_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"lba_weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"failover_available_addr_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"failover_latency_optimization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failover_lba_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failover_max_return_addr_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"failover_min_available_addr_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lines": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"line_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"line_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"strategy_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"strategy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudAlidnsAccessStrategiesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDnsGtmAccessStrategies"
	request := make(map[string]interface{})
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	request["StrategyMode"] = d.Get("strategy_mode")
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	conn, err := client.NewAlidnsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-09"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alidns_access_strategies", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Strategies.Strategy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Strategies.Strategy", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil {
				if !nameRegex.MatchString(fmt.Sprint(item["StrategyName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["StrategyId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                             fmt.Sprint(object["StrategyId"]),
			"access_strategy_id":             fmt.Sprint(object["StrategyId"]),
			"create_time":                    object["CreateTime"],
			"create_timestamp":               fmt.Sprint(object["CreateTimestamp"]),
			"effective_addr_pool_group_type": object["EffectiveAddrPoolGroupType"],
			"strategy_name":                  object["StrategyName"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["StrategyName"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["StrategyId"])
		alidnsService := AlidnsService{client}
		getResp, err := alidnsService.DescribeAlidnsAccessStrategy(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["default_addr_pool_type"] = getResp["DefaultAddrPoolType"]

		defaultAddrPool := make([]map[string]interface{}, 0)
		if v, ok := getResp["DefaultAddrPools"]; ok {
			defaultAddrPools := v.(map[string]interface{})
			if v, ok := defaultAddrPools["DefaultAddrPool"]; ok {
				defaultAddrPoolList := v.([]interface{})
				for _, v := range defaultAddrPoolList {
					if m1, ok := v.(map[string]interface{}); ok {
						temp1 := map[string]interface{}{
							"addr_count":   formatInt(m1["AddrCount"]),
							"addr_pool_id": m1["Id"],
							"lba_weight":   formatInt(m1["LbaWeight"]),
							"name":         m1["Name"],
						}
						defaultAddrPool = append(defaultAddrPool, temp1)
					}
				}
			}
		}

		mapping["default_addr_pools"] = defaultAddrPool
		if v, ok := getResp["DefaultAvailableAddrNum"]; ok && fmt.Sprint(v) != "0" {
			mapping["default_available_addr_num"] = formatInt(v)
		}
		mapping["default_latency_optimization"] = getResp["DefaultLatencyOptimization"]
		mapping["default_lba_strategy"] = getResp["DefaultLbaStrategy"]
		if v, ok := getResp["DefaultMaxReturnAddrNum"]; ok && fmt.Sprint(v) != "0" {
			mapping["default_max_return_addr_num"] = formatInt(v)
		}
		if v, ok := getResp["DefaultMinAvailableAddrNum"]; ok && fmt.Sprint(v) != "0" {
			mapping["default_min_available_addr_num"] = formatInt(v)
		}
		mapping["failover_addr_pool_type"] = getResp["FailoverAddrPoolType"]
		mapping["access_mode"] = getResp["AccessMode"]

		failoverAddrPool := make([]map[string]interface{}, 0)
		if v, ok := getResp["FailoverAddrPools"]; ok {
			failoverAddrPools := v.(map[string]interface{})
			if v, ok := failoverAddrPools["FailoverAddrPool"]; ok {
				failoverAddrPoolList := v.([]interface{})
				for _, v := range failoverAddrPoolList {
					if m1, ok := v.(map[string]interface{}); ok {
						temp1 := map[string]interface{}{
							"addr_count":   formatInt(m1["AddrCount"]),
							"addr_pool_id": m1["Id"],
							"lba_weight":   formatInt(m1["LbaWeight"]),
							"name":         m1["Name"],
						}
						failoverAddrPool = append(failoverAddrPool, temp1)
					}
				}
			}
		}
		mapping["failover_addr_pools"] = failoverAddrPool
		if v, ok := getResp["FailoverAvailableAddrNum"]; ok && fmt.Sprint(v) != "0" {
			mapping["failover_available_addr_num"] = formatInt(v)
		}
		mapping["failover_latency_optimization"] = getResp["FailoverLatencyOptimization"]
		mapping["failover_lba_strategy"] = getResp["FailoverLbaStrategy"]
		if v, ok := getResp["FailoverMaxReturnAddrNum"]; ok && fmt.Sprint(v) != "0" {
			mapping["failover_max_return_addr_num"] = formatInt(v)
		}
		if v, ok := getResp["FailoverMinAvailableAddrNum"]; ok && fmt.Sprint(v) != "0" {
			mapping["failover_min_available_addr_num"] = formatInt(v)
		}
		mapping["instance_id"] = getResp["InstanceId"]

		lineMaps := make([]map[string]interface{}, 0)
		if v, ok := getResp["Lines"]; ok {
			lines := v.(map[string]interface{})
			if v, ok := lines["Line"]; ok {
				lineList := v.([]interface{})
				for _, v := range lineList {
					if m1, ok := v.(map[string]interface{}); ok {
						temp1 := map[string]interface{}{
							"group_code": m1["GroupCode"],
							"group_name": m1["GroupName"],
							"line_code":  m1["LineCode"],
							"line_name":  m1["LineName"],
						}
						lineMaps = append(lineMaps, temp1)
					}
				}
			}
		}
		mapping["lines"] = lineMaps
		mapping["strategy_mode"] = getResp["StrategyMode"]
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("strategies", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

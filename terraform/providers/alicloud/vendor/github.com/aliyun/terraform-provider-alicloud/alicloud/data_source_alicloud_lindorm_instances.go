package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudLindormInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudLindormInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
			"query_str": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVATION", "CLASS_CHANGING", "CREATING", "DELETED", "INSTANCE_LEVEL_MODIFY", "LOCKED", "MINOR_VERSION_TRANSING", "NET_MODIFYING", "RESIZING", "RESTARTING"}, false),
			},
			"support_engine": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 7),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_renew": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cold_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deletion_proection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_usage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_warning_threshold": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_engine_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"file_engine_specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_storage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_white_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"lts_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lts_node_specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phoenix_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"phoenix_node_specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"search_engine_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"search_engine_specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"table_engine_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"table_engine_specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_series_engine_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"time_serires_engine_specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
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

func dataSourceAlicloudLindormInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "GetLindormInstanceList"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("query_str"); ok {
		request["QueryStr"] = v
	}
	if v, ok := d.GetOk("support_engine"); ok {
		request["SupportEngine"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var instanceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		instanceNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, val := range v.([]interface{}) {
			if val == nil {
				continue
			}
			idsMap[val.(string)] = val.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewHitsdbClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-15"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_lindorm_instances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.InstanceList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.InstanceList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if instanceNameRegex != nil && !instanceNameRegex.MatchString(fmt.Sprint(item["InstanceAlias"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["InstanceStatus"].(string) {
				continue
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
			"create_time":      fmt.Sprint(object["CreateMilliseconds"]),
			"engine_type":      formatInt(object["EngineType"]),
			"expired_time":     fmt.Sprint(object["ExpiredMilliseconds"]),
			"id":               fmt.Sprint(object["InstanceId"]),
			"instance_id":      fmt.Sprint(object["InstanceId"]),
			"instance_name":    object["InstanceAlias"],
			"instance_storage": object["InstanceStorage"],
			"network_type":     object["NetworkType"],
			"payment_type":     convertLindormInstancePaymentTypeResponse(object["PayType"]),
			"service_type":     object["ServiceType"],
			"status":           object["InstanceStatus"],
			"vpc_id":           object["VpcId"],
			"zone_id":          object["ZoneId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["InstanceAlias"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["InstanceId"])
		hitsdbService := HitsdbService{client}
		ipWhiteList, err := hitsdbService.GetInstanceIpWhiteList(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["ip_white_list"] = convertListToJsonString(ipWhiteList)
		getResp1, err := hitsdbService.DescribeLindormInstance(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["auto_renew"] = getResp1["AutoRenew"]
		if v, ok := getResp1["ColdStorage"]; ok {
			mapping["cold_storage"] = formatInt(v)
		}
		mapping["deletion_proection"] = getResp1["DeletionProtection"]
		mapping["disk_category"] = getResp1["DiskCategory"]
		mapping["disk_usage"] = getResp1["DiskUsage"]
		mapping["disk_warning_threshold"] = getResp1["DiskThreshold"]
		mapping["resource_owner_id"] = getResp1["AliUid"]
		mapping["vswitch_id"] = getResp1["VswitchId"]

		getResp2, err := hitsdbService.GetLindormInstanceEngineInfo(id)
		if err != nil {
			return WrapError(err)
		}
		if v, ok := getResp2["FileEngineNodeCount"]; ok {
			mapping["file_engine_node_count"] = formatInt(v)
		}
		mapping["file_engine_specification"] = getResp2["FileEngineSpecification"]
		if v, ok := getResp2["LtsNodeCount"]; ok {
			mapping["lts_node_count"] = formatInt(v)
		}
		mapping["lts_node_specification"] = getResp2["LtsNodeSpecification"]
		if v, ok := getResp2["PhoenixNodeCount"]; ok {
			mapping["phoenix_node_count"] = formatInt(v)
		}
		mapping["phoenix_node_specification"] = getResp2["PhoenixNodeSpecification"]
		if v, ok := getResp2["SearchEngineNodeCount"]; ok {
			mapping["search_engine_node_count"] = formatInt(v)
		}
		mapping["search_engine_specification"] = getResp2["SearchEngineSpecification"]
		if v, ok := getResp2["TableEngineNodeCount"]; ok {
			mapping["table_engine_node_count"] = formatInt(v)
		}
		mapping["table_engine_specification"] = getResp2["TableEngineSpecification"]
		if v, ok := getResp2["TimeSeriesNodeCount"]; ok {
			mapping["time_series_engine_node_count"] = formatInt(v)
		}
		mapping["time_serires_engine_specification"] = getResp2["TimeSeriesSpecification"]

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

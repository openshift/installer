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

func dataSourceAlicloudPvtzRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPvtzRulesRead,
		Schema: map[string]*schema.Schema{
			"endpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forward_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"endpoint_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bind_vpcs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudPvtzRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeResolverRules"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("endpoint_id"); ok {
		request["EndpointId"] = v
	}

	request["Lang"] = "en"
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["NeedDetailAttributes"] = true
	var objects []map[string]interface{}
	var ruleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		ruleNameRegex = r
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
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_pvtz_rules", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Rules", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Rules", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if ruleNameRegex != nil && !ruleNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
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
			"create_time":   object["CreateTime"],
			"endpoint_id":   object["EndpointId"],
			"endpoint_name": object["EndpointName"],
			"id":            fmt.Sprint(object["Id"]),
			"rule_id":       fmt.Sprint(object["Id"]),
			"rule_name":     object["Name"],
			"type":          object["Type"],
			"zone_name":     object["ZoneName"],
		}
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
		mapping["forward_ips"] = forwardConfigsSli

		bindVpcsSlice := make([]map[string]interface{}, 0)
		if bindVpcs, ok := object["BindVpcs"].([]interface{}); ok {
			for _, bindVpcsArgs := range bindVpcs {
				bindVpcsArg := bindVpcsArgs.(map[string]interface{})
				bindVpcsMap := make(map[string]interface{})
				bindVpcsMap["vpc_id"] = bindVpcsArg["VpcId"]
				bindVpcsMap["region_id"] = bindVpcsArg["RegionId"]
				bindVpcsMap["vpc_name"] = bindVpcsArg["VpcName"]
				bindVpcsMap["region_name"] = bindVpcsArg["RegionName"]
				bindVpcsSlice = append(bindVpcsSlice, bindVpcsMap)
			}
		}
		mapping["bind_vpcs"] = bindVpcsSlice
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("rules", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

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

func dataSourceAlicloudAlidnsAddressPools() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlidnsAddressPoolsRead,
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
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address_pool_name": {
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
						"lba_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"monitor_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"monitor_config_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"attribute_info": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"lba_weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"remark": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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

func dataSourceAlicloudAlidnsAddressPoolsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDnsGtmInstanceAddressPools"
	request := make(map[string]interface{})
	request["InstanceId"] = d.Get("instance_id")
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var addressPoolNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		addressPoolNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alidns_address_pools", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AddrPools.AddrPool", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AddrPools.AddrPool", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if addressPoolNameRegex != nil && !addressPoolNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AddrPoolId"])]; !ok {
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
			"id":                fmt.Sprint(object["AddrPoolId"]),
			"instance_id":       request["InstanceId"],
			"address_pool_id":   fmt.Sprint(object["AddrPoolId"]),
			"address_pool_name": object["Name"],
			"create_time":       object["CreateTime"],
			"create_timestamp":  fmt.Sprint(object["CreateTimestamp"]),
			"lba_strategy":      object["LbaStrategy"],
			"type":              object["Type"],
			"update_timestamp":  fmt.Sprint(object["UpdateTimestamp"]),
			"update_time":       fmt.Sprint(object["UpdateTime"]),
			"monitor_status":    object["MonitorStatus"],
			"monitor_config_id": object["MonitorConfigId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["AddrPoolId"])
		alidnsService := AlidnsService{client}
		getResp, err := alidnsService.DescribeAlidnsAddressPool(id)
		if err != nil {
			return WrapError(err)
		}
		addressMaps := make([]map[string]interface{}, 0)
		if v, ok := getResp["Addrs"]; ok {
			addressList := v.(map[string]interface{})
			if v, ok := addressList["Addr"]; ok {
				addressSli := v.([]interface{})
				for _, v := range addressSli {
					if m1, ok := v.(map[string]interface{}); ok {
						addressMap := make(map[string]interface{})
						addressMap["attribute_info"] = m1["AttributeInfo"]
						addressMap["remark"] = m1["Remark"]
						if v, ok := m1["LbaWeight"]; ok {
							addressMap["lba_weight"] = v
						}
						addressMap["address"] = m1["Addr"]
						addressMap["mode"] = m1["Mode"]
						addressMaps = append(addressMaps, addressMap)
					}
				}
			}
		}
		mapping["address"] = addressMaps

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("pools", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudAlidnsGtmInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlidnsGtmInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dingtalk_notice": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"email_notice": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"notice_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sms_notice": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"alert_group": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cname_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
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
						"strategy_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_cname_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_rr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_user_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_zone_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"package_edition": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlidnsGtmInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDnsGtmInstances"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}

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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alidns_gtm_instances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.GtmInstances", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.GtmInstances", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
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
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"payment_type":      object["PaymentType"],
			"create_time":       object["CreateTime"],
			"expire_time":       object["ExpireTime"],
			"id":                fmt.Sprint(object["InstanceId"]),
			"instance_id":       fmt.Sprint(object["InstanceId"]),
			"resource_group_id": object["ResourceGroupId"],
			"package_edition":   object["VersionCode"],
		}

		if config, ok := object["Config"].(map[string]interface{}); ok {
			mapping["public_zone_name"] = config["PublicZoneName"]
			mapping["cname_type"] = config["CnameType"]
			mapping["instance_name"] = config["InstanceName"]
			mapping["strategy_mode"] = config["StrategyMode"]
			mapping["public_cname_mode"] = config["PublicCnameMode"]
			mapping["public_rr"] = config["PublicRr"]
			mapping["public_user_domain_name"] = config["PublicUserDomainName"]
			if v, ok := config["Ttl"]; ok {
				mapping["ttl"] = formatInt(v)
			}
			if v, ok := config["AlertGroup"].(string); ok {
				vv, err := convertJsonStringToList(v)
				if err != nil {
					return WrapError(err)
				} else {
					mapping["alert_group"] = vv
				}
			}

			if alertConfigsList, ok := config["AlertConfig"]; ok {
				alertConfigConfigArgs := alertConfigsList.([]interface{})
				alertConfigsMaps := make([]map[string]interface{}, 0)
				for _, alertConfigMapArgitem := range alertConfigConfigArgs {
					alertConfigMapArg := alertConfigMapArgitem.(map[string]interface{})
					alertConfigsMap := map[string]interface{}{}
					alertConfigsMap["sms_notice"] = alertConfigMapArg["SmsNotice"]
					alertConfigsMap["notice_type"] = alertConfigMapArg["NoticeType"]
					alertConfigsMap["email_notice"] = alertConfigMapArg["EmailNotice"]
					alertConfigsMap["dingtalk_notice"] = alertConfigMapArg["DingtalkNotice"]
					alertConfigsMaps = append(alertConfigsMaps, alertConfigsMap)
				}
				mapping["alert_config"] = alertConfigsMaps
			}
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
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

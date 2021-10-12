package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEipanycastAnycastEipAddresses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEipanycastAnycastEipAddressesRead,
		Schema: map[string]*schema.Schema{
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
			"anycast_eip_address_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"bind_instance_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"business_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"internet_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "PayByTraffic",
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo"}, false),
				Default:      "PayAsYouGo",
			},
			"service_location": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Allocated", "Associated", "Associating", "Unassociating"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ali_uid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"anycast_eip_address_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"anycast_eip_bind_info_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bind_instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"bind_instance_region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"bind_instance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"bind_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"anycast_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEipanycastAnycastEipAddressesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListAnycastEipAddresses"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("anycast_eip_address_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("bind_instance_ids"); ok {
		request["BindInstanceIds"] = convertListToJsonString(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("business_status"); ok {
		request["BusinessStatus"] = v
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}
	if v, ok := d.GetOk("ip_address"); ok {
		request["AnycastEipAddress"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["InstanceChargeType"] = convertEipanycastAnycastEipAddressPaymentTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("service_location"); ok {
		request["ServiceLocation"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var anycastEipAddressNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		anycastEipAddressNameRegex = r
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
	conn, err := client.NewEipanycastClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-09"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_eipanycast_anycast_eip_addresses", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.AnycastList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AnycastList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if anycastEipAddressNameRegex != nil {
				if !anycastEipAddressNameRegex.MatchString(fmt.Sprint(item["Name"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AnycastId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"ali_uid":                  object["AliUid"],
			"anycast_eip_address_name": object["Name"],
			"id":                       fmt.Sprint(object["AnycastId"]),
			"anycast_id":               fmt.Sprint(object["AnycastId"]),
			"bandwidth":                formatInt(object["Bandwidth"]),
			"bid":                      object["Bid"],
			"business_status":          object["BusinessStatus"],
			"description":              object["Description"],
			"internet_charge_type":     object["InternetChargeType"],
			"ip_address":               object["IpAddress"],
			"payment_type":             convertEipanycastAnycastEipAddressPaymentTypeResponse(object["InstanceChargeType"].(string)),
			"service_location":         object["ServiceLocation"],
			"status":                   object["Status"],
		}

		anycastEipBindInfoList := make([]map[string]interface{}, 0)
		if anycastEipBindInfoListList, ok := object["AnycastEipBindInfoList"].([]interface{}); ok {
			for _, v := range anycastEipBindInfoListList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"bind_instance_id":        m1["BindInstanceId"],
						"bind_instance_region_id": m1["BindInstanceRegionId"],
						"bind_instance_type":      m1["BindInstanceType"],
						"bind_time":               m1["BindTime"],
					}
					anycastEipBindInfoList = append(anycastEipBindInfoList, temp1)
				}
			}
		}
		mapping["anycast_eip_bind_info_list"] = anycastEipBindInfoList
		ids = append(ids, fmt.Sprint(object["AnycastId"]))
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

	if err := d.Set("addresses", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

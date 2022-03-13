package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpcIpv6InternetBandwidths() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpcIpv6InternetBandwidthsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ipv6_internet_bandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipv6_address_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"FinacialLocked", "Normal", "SecurityLocked"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidths": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_address_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_internet_bandwidth_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
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

func dataSourceAlicloudVpcIpv6InternetBandwidthsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeIpv6Addresses"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("ipv6_internet_bandwidth_id"); ok {
		request["Ipv6InternetBandwidthId"] = v
	}
	if v, ok := d.GetOk("ipv6_address_id"); ok {
		request["Ipv6AddressId"] = v
	}
	request["RegionId"] = client.RegionId
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_ipv6_internet_bandwidths", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Ipv6Addresses.Ipv6Address", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Ipv6Addresses.Ipv6Address", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Ipv6InternetBandwidth"].(map[string]interface{})["Ipv6InternetBandwidthId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Ipv6InternetBandwidth"].(map[string]interface{})["BusinessStatus"].(string) {
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
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"ipv6_address_id": object["Ipv6AddressId"],
			"ipv6_gateway_id": object["Ipv6GatewayId"],
		}
		if ipv6InternetBandwidth, ok := object["Ipv6InternetBandwidth"]; ok {
			if v, ok := ipv6InternetBandwidth.(map[string]interface{}); ok {
				if v, ok := v["Bandwidth"]; ok {
					mapping["bandwidth"] = formatInt(v)
				}
				if v, ok := v["BusinessStatus"]; ok {
					mapping["status"] = v
				}
				if v, ok := v["InstanceChargeType"]; ok {
					mapping["payment_type"] = convertVpcIpv6InternetBandwidthInstanceChargeTypeResponse(v.(string))
				}
				if v, ok := v["InternetChargeType"]; ok {
					mapping["internet_charge_type"] = v
				}
				if v, ok := v["Ipv6InternetBandwidthId"]; ok {
					mapping["ipv6_internet_bandwidth_id"] = fmt.Sprint(v)
					mapping["id"] = fmt.Sprint(v)
				}
			}
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("bandwidths", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

func convertVpcIpv6InternetBandwidthInstanceChargeTypeResponse(source string) string {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	}
	return source
}

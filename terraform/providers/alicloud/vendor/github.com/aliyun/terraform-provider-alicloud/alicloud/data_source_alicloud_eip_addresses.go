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

func dataSourceAlicloudEipAddresses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEipAddressesRead,
		Schema: map[string]*schema.Schema{
			"ip_addresses": {
				Type:       schema.TypeList,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				ForceNew:   true,
				MinItems:   1,
				Deprecated: "Field 'ip_addresses' has been deprecated from provider version 1.126.0 and it will be removed in the future version. Please use the new attribute 'ip_address' instead.",
			},
			"in_use": {
				Type:     schema.TypeBool,
				Optional: true,
				Removed:  "Field 'in_use' has been removed from provider version 1.126.0.",
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
			"address_name": {
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
			"associated_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"associated_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"include_reservation_data": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"isp": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"BGP", "BGP_PRO"}, false),
			},
			"lock_reason": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"segment_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Associating", "Available", "InUse", "Unassociating"}, false),
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allocation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_regions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_package_bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_package_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_package_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deletion_protection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"has_reservation_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hd_monitor_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
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
						"isp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_locks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reservation_active_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reservation_bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reservation_internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reservation_order_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"second_limited": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"segment_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
			"eips": {
				Type:       schema.TypeList,
				Computed:   true,
				Deprecated: "Field 'eips' has been deprecated from provider version 1.126.0 and it will be removed in the future version. Please use the new attribute 'addresses' instead.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deletion_protection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func dataSourceAlicloudEipAddressesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeEipAddresses"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("address_name"); ok {
		request["EipName"] = v
	}
	if v, ok := d.GetOk("associated_instance_id"); ok {
		request["AssociatedInstanceId"] = v
	}
	if v, ok := d.GetOk("associated_instance_type"); ok {
		request["AssociatedInstanceType"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOkExists("include_reservation_data"); ok {
		request["IncludeReservationData"] = v
	}
	if v, ok := d.GetOk("ip_address"); ok {
		request["EipAddress"] = v
	}
	if v, ok := d.GetOk("isp"); ok {
		request["ISP"] = v
	}
	if v, ok := d.GetOk("lock_reason"); ok {
		request["LockReason"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertEipAddressPaymentTypeRequest(v)
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("segment_instance_id"); ok {
		request["SegmentInstanceId"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var addressNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		addressNameRegex = r
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
	ipsMap := make(map[string]string)
	if v, ok := d.GetOk("ip_addresses"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	tagsMap := make(map[string]interface{})
	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		tagsMap = v.(map[string]interface{})
	}
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_eip_addresses", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.EipAddresses.EipAddress", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.EipAddresses.EipAddress", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if addressNameRegex != nil {
				if !addressNameRegex.MatchString(fmt.Sprint(item["Name"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AllocationId"])]; !ok {
					continue
				}
			}
			if len(ipsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["IpAddress"])]; !ok {
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
	eips := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"address_name":                     object["Name"],
			"id":                               fmt.Sprint(object["AllocationId"]),
			"allocation_id":                    fmt.Sprint(object["AllocationId"]),
			"available_regions":                object["AvailableRegions"].(map[string]interface{})["AvailableRegion"],
			"bandwidth":                        object["Bandwidth"],
			"bandwidth_package_bandwidth":      object["BandwidthPackageBandwidth"],
			"bandwidth_package_id":             object["BandwidthPackageId"],
			"bandwidth_package_type":           object["BandwidthPackageType"],
			"create_time":                      object["AllocationTime"],
			"deletion_protection":              object["DeletionProtection"],
			"description":                      object["Descritpion"],
			"expired_time":                     object["ExpiredTime"],
			"has_reservation_data":             fmt.Sprint(object["HasReservationData"]),
			"hd_monitor_status":                object["HDMonitorStatus"],
			"instance_id":                      object["InstanceId"],
			"instance_region_id":               object["InstanceRegionId"],
			"instance_type":                    object["InstanceType"],
			"internet_charge_type":             object["InternetChargeType"],
			"ip_address":                       object["IpAddress"],
			"isp":                              object["ISP"],
			"payment_type":                     convertEipAddressPaymentTypeResponse(object["ChargeType"]),
			"reservation_active_time":          object["ReservationActiveTime"],
			"reservation_bandwidth":            object["ReservationBandwidth"],
			"reservation_internet_charge_type": object["ReservationInternetChargeType"],
			"reservation_order_type":           object["ReservationOrderType"],
			"resource_group_id":                object["ResourceGroupId"],
			"second_limited":                   object["SecondLimited"],
			"segment_instance_id":              object["SegmentInstanceId"],
			"status":                           object["Status"],
		}
		eipMapping := map[string]interface{}{
			"id":                   object["AllocationId"],
			"status":               object["Status"],
			"ip_address":           object["IpAddress"],
			"bandwidth":            object["Bandwidth"],
			"instance_id":          object["InstanceId"],
			"instance_type":        object["InstanceType"],
			"internet_charge_type": object["InternetChargeType"],
			"creation_time":        object["AllocationTime"],
			"deletion_protection":  object["DeletionProtection"],
		}
		if v := object["OperationLocks"].(map[string]interface{})["LockReason"].([]interface{}); len(v) > 0 {
			mapping["operation_locks"] = v[0].(map[string]interface{})["LockReason"]
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["AllocationId"]))
			names = append(names, object["Name"])
			s = append(s, mapping)
			eips = append(eips, eipMapping)
			continue
		}

		vpcService := VpcService{client}
		id := fmt.Sprint(object["AllocationId"])
		getResp, err := vpcService.ListTagResources(id, "EIP")
		if err != nil {
			return WrapError(err)
		}
		mapping["tags"] = tagsToMap(getResp)
		march := false
		for tagKey, tagValue := range tagsToMap(getResp) {
			if v, ok := tagsMap[tagKey]; ok && v == tagValue {
				march = true
				break
			}
		}
		if !march && len(tagsMap) > 0 {
			continue
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		s = append(s, mapping)
		eips = append(eips, eipMapping)
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

	if err := d.Set("eips", eips); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

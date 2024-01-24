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

func dataSourceAlicloudCommonBandwidthPackages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCommonBandwidthPackagesRead,
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
			"bandwidth_package_name": {
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
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Available", "Pending"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"packages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_package_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_package_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"business_status": {
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
							Type:     schema.TypeBool,
							Computed: true,
						},
						"isp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allocation_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"bandwidth_package_ip_relation_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ratio": {
							Type:     schema.TypeInt,
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
						"service_managed": {
							Type:     schema.TypeInt,
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

func dataSourceAlicloudCommonBandwidthPackagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeCommonBandwidthPackages"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("bandwidth_package_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOkExists("include_reservation_data"); ok {
		request["IncludeReservationData"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var bandwidthPackageNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		bandwidthPackageNameRegex = r
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_common_bandwidth_packages", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.CommonBandwidthPackages.CommonBandwidthPackage", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.CommonBandwidthPackages.CommonBandwidthPackage", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if bandwidthPackageNameRegex != nil {
				if !bandwidthPackageNameRegex.MatchString(fmt.Sprint(item["Name"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["BandwidthPackageId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
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
			"bandwidth":                        object["Bandwidth"],
			"id":                               fmt.Sprint(object["BandwidthPackageId"]),
			"bandwidth_package_id":             fmt.Sprint(object["BandwidthPackageId"]),
			"bandwidth_package_name":           object["Name"],
			"name":                             object["Name"],
			"business_status":                  object["BusinessStatus"],
			"deletion_protection":              object["DeletionProtection"],
			"description":                      object["Description"],
			"expired_time":                     object["ExpiredTime"],
			"has_reservation_data":             object["HasReservationData"],
			"isp":                              object["ISP"],
			"internet_charge_type":             object["InternetChargeType"],
			"payment_type":                     object["InstanceChargeType"],
			"ratio":                            formatInt(object["Ratio"]),
			"reservation_active_time":          object["ReservationActiveTime"],
			"reservation_bandwidth":            object["ReservationBandwidth"],
			"reservation_internet_charge_type": object["ReservationInternetChargeType"],
			"reservation_order_type":           object["ReservationOrderType"],
			"resource_group_id":                object["ResourceGroupId"],
			"service_managed":                  formatInt(object["ServiceManaged"]),
			"status":                           object["Status"],
		}

		publicIpAddresse := make([]map[string]interface{}, 0)
		if publicIpAddresseList, ok := object["PublicIpAddresses"].(map[string]interface{})["PublicIpAddresse"].([]interface{}); ok {
			for _, v := range publicIpAddresseList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"allocation_id":                        m1["AllocationId"],
						"bandwidth_package_ip_relation_status": m1["BandwidthPackageIpRelationStatus"],
						"ip_address":                           m1["IpAddress"],
					}
					publicIpAddresse = append(publicIpAddresse, temp1)
				}
			}
		}
		mapping["public_ip_addresses"] = publicIpAddresse
		ids = append(ids, fmt.Sprint(object["BandwidthPackageId"]))
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

	if err := d.Set("packages", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

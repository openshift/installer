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

func dataSourceAlicloudGaAccelerators() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGaAcceleratorsRead,
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"active", "binding", "configuring", "deleting", "finacialLocked", "init", "unbinding"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"accelerators": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accelerator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accelerator_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"basic_bandwidth_package": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bandwidth": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"bandwidth_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cross_domain_bandwidth_package": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bandwidth": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ddos_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"second_dns_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec": {
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

func dataSourceAlicloudGaAcceleratorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListAccelerators"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var acceleratorNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		acceleratorNameRegex = r
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
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ga_accelerators", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Accelerators", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Accelerators", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if acceleratorNameRegex != nil {
				if !acceleratorNameRegex.MatchString(fmt.Sprint(item["Name"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AcceleratorId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["State"].(string) {
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
			"id":               fmt.Sprint(object["AcceleratorId"]),
			"accelerator_id":   fmt.Sprint(object["AcceleratorId"]),
			"accelerator_name": object["Name"],
			"cen_id":           object["CenId"],
			"ddos_id":          object["DdosId"],
			"description":      object["Description"],
			"dns_name":         object["DnsName"],
			"expired_time":     formatInt(object["ExpiredTime"]),
			"payment_type":     object["InstanceChargeType"],
			"second_dns_name":  object["SecondDnsName"],
			"spec":             object["Spec"],
			"status":           object["State"],
		}

		basicBandwidthPackageSli := make([]map[string]interface{}, 0)
		if v, ok := object["BasicBandwidthPackage"].(map[string]interface{}); ok && len(v) > 0 {
			basicBandwidthPackage := object["BasicBandwidthPackage"]
			basicBandwidthPackageMap := make(map[string]interface{})
			basicBandwidthPackageMap["bandwidth"] = basicBandwidthPackage.(map[string]interface{})["Bandwidth"]
			basicBandwidthPackageMap["bandwidth_type"] = basicBandwidthPackage.(map[string]interface{})["BandwidthType"]
			basicBandwidthPackageMap["instance_id"] = basicBandwidthPackage.(map[string]interface{})["InstanceId"]
			basicBandwidthPackageSli = append(basicBandwidthPackageSli, basicBandwidthPackageMap)
		}
		mapping["basic_bandwidth_package"] = basicBandwidthPackageSli

		crossDomainBandwidthPackageSli := make([]map[string]interface{}, 0)
		if v, ok := object["CrossDomainBandwidthPackage"].(map[string]interface{}); ok && len(v) > 0 {
			crossDomainBandwidthPackage := object["CrossDomainBandwidthPackage"]
			crossDomainBandwidthPackageMap := make(map[string]interface{})
			crossDomainBandwidthPackageMap["bandwidth"] = crossDomainBandwidthPackage.(map[string]interface{})["Bandwidth"]
			crossDomainBandwidthPackageMap["instance_id"] = crossDomainBandwidthPackage.(map[string]interface{})["InstanceId"]
			crossDomainBandwidthPackageSli = append(crossDomainBandwidthPackageSli, crossDomainBandwidthPackageMap)
		}
		mapping["cross_domain_bandwidth_package"] = crossDomainBandwidthPackageSli
		ids = append(ids, fmt.Sprint(object["AcceleratorId"]))
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

	if err := d.Set("accelerators", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

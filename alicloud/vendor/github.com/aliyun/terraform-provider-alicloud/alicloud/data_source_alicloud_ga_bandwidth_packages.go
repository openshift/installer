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

func dataSourceAlicloudGaBandwidthPackages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGaBandwidthPackagesRead,
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
				ValidateFunc: validation.StringInSlice([]string{"active", "binded", "binding", "finacialLocked", "init", "unbinding", "updating"}, false),
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "CrossDomain"}, false),
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
							Type:     schema.TypeInt,
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
						"bandwidth_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cbn_geographic_region_ida": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cbn_geographic_region_idb": {
							Type:     schema.TypeString,
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
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
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

func dataSourceAlicloudGaBandwidthPackagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListBandwidthPackages"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("status"); ok {
		request["State"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
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
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ga_bandwidth_packages", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.BandwidthPackages", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.BandwidthPackages", response)
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
			"bandwidth":                 formatInt(object["Bandwidth"]),
			"id":                        fmt.Sprint(object["BandwidthPackageId"]),
			"bandwidth_package_id":      fmt.Sprint(object["BandwidthPackageId"]),
			"bandwidth_package_name":    object["Name"],
			"bandwidth_type":            object["BandwidthType"],
			"cbn_geographic_region_ida": object["CbnGeographicRegionIdA"],
			"cbn_geographic_region_idb": object["CbnGeographicRegionIdB"],
			"description":               object["Description"],
			"expired_time":              object["ExpiredTime"],
			"payment_type":              convertGaBandwidthPackagePaymentTypeResponse(object["ChargeType"].(string)),
			"type":                      object["Type"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["BandwidthPackageId"]))
			names = append(names, object["Name"])
			s = append(s, mapping)
			continue
		}

		gaService := GaService{client}
		id := fmt.Sprint(object["BandwidthPackageId"])
		getResp, err := gaService.DescribeGaBandwidthPackage(id)
		if err != nil {
			return WrapError(err)
		}
		if statusOk && status != "" && status != getResp["State"].(string) {
			continue
		}
		mapping["status"] = getResp["State"]
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

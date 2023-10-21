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

func dataSourceAlicloudPvtzZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPvtzZonesRead,
		Schema: map[string]*schema.Schema{
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"en", "jp", "zh"}, false),
			},
			"query_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"query_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"search_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"LIKE", "EXACT"}, false),
			},
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_vpcs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"create_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_ptr": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"proxy_pattern": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slave_dns": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"update_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
							Removed:  "Field 'creation_time' has been removed from provider version 1.107.0",
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
							Removed:  "Field 'update_time' has been removed from provider version 1.107.0",
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

func dataSourceAlicloudPvtzZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeZones"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("query_region_id"); ok {
		request["QueryRegionId"] = v
	}
	if v, ok := d.GetOk("query_vpc_id"); ok {
		request["QueryVpcId"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("search_mode"); ok {
		request["SearchMode"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var zoneNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		zoneNameRegex = r
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
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_pvtz_zones", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Zones.Zone", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Zones.Zone", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if zoneNameRegex != nil {
				if !zoneNameRegex.MatchString(item["ZoneName"].(string)) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ZoneId"])]; !ok {
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
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"create_timestamp":  formatInt(object["CreateTimestamp"]),
			"is_ptr":            object["IsPtr"],
			"proxy_pattern":     object["ProxyPattern"],
			"record_count":      formatInt(object["RecordCount"]),
			"remark":            object["Remark"],
			"resource_group_id": object["ResourceGroupId"],
			"update_timestamp":  formatInt(object["UpdateTimestamp"]),
			"id":                fmt.Sprint(object["ZoneId"]),
			"zone_id":           fmt.Sprint(object["ZoneId"]),
			"zone_name":         object["ZoneName"],
			"name":              object["ZoneName"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["ZoneId"]))
			names = append(names, object["ZoneName"].(string))
			s = append(s, mapping)
			continue
		}

		pvtzService := PvtzService{client}
		id := fmt.Sprint(object["ZoneId"])
		getResp, err := pvtzService.DescribePvtzZone(id)
		if err != nil {
			return WrapError(err)
		}

		vpc := make([]map[string]interface{}, 0)
		if vpcList, ok := getResp["BindVpcs"].(map[string]interface{})["Vpc"].([]interface{}); ok {
			for _, v := range vpcList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"region_id":   m1["RegionId"],
						"region_name": m1["RegionName"],
						"vpc_id":      m1["VpcId"],
						"vpc_name":    m1["VpcName"],
					}
					vpc = append(vpc, temp1)
				}
			}
		}
		mapping["bind_vpcs"] = vpc
		mapping["slave_dns"] = getResp["SlaveDns"]
		ids = append(ids, fmt.Sprint(object["ZoneId"]))
		names = append(names, object["ZoneName"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

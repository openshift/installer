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

func dataSourceAlicloudSnatEntries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSnatEntriesRead,
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
			"snat_entry_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snat_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snat_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Available", "Deleting", "Pending"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snat_entry_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snat_entry_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snat_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_vswitch_id": {
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

func dataSourceAlicloudSnatEntriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeSnatTableEntries"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("snat_entry_name"); ok {
		request["SnatEntryName"] = v
	}
	if v, ok := d.GetOk("snat_ip"); ok {
		request["SnatIp"] = v
	}
	request["SnatTableId"] = d.Get("snat_table_id")
	if v, ok := d.GetOk("source_cidr"); ok {
		request["SourceCIDR"] = v
	}
	if v, ok := d.GetOk("source_vswitch_id"); ok {
		request["SourceVSwitchId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var snatEntryNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		snatEntryNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_snat_entries", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.SnatTableEntries.SnatTableEntry", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.SnatTableEntries.SnatTableEntry", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if snatEntryNameRegex != nil {
				if !snatEntryNameRegex.MatchString(fmt.Sprint(item["SnatEntryName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["SnatEntryId"])]; !ok {
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
			"id":                fmt.Sprint(object["SnatEntryId"]),
			"snat_entry_id":     fmt.Sprint(object["SnatEntryId"]),
			"snat_entry_name":   object["SnatEntryName"],
			"snat_ip":           object["SnatIp"],
			"source_cidr":       object["SourceCIDR"],
			"source_vswitch_id": object["SourceVSwitchId"],
			"status":            object["Status"],
		}
		ids = append(ids, fmt.Sprint(object["SnatEntryId"]))
		names = append(names, object["SnatEntryName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("entries", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

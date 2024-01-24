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

func dataSourceAlicloudForwardEntries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudForwardEntriesRead,
		Schema: map[string]*schema.Schema{
			"external_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"external_port": {
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
			"forward_entry_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"forward_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"internal_port": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"any", "tcp", "udp"}, false),
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
						"external_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forward_entry_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forward_entry_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_protocol": {
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

func dataSourceAlicloudForwardEntriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeForwardTableEntries"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("external_ip"); ok {
		request["ExternalIp"] = v
	}
	if v, ok := d.GetOk("external_port"); ok {
		request["ExternalPort"] = v
	}
	if v, ok := d.GetOk("forward_entry_name"); ok {
		request["ForwardEntryName"] = v
	}
	request["ForwardTableId"] = d.Get("forward_table_id")
	if v, ok := d.GetOk("internal_ip"); ok {
		request["InternalIp"] = v
	}
	if v, ok := d.GetOk("internal_port"); ok {
		request["InternalPort"] = v
	}
	if v, ok := d.GetOk("ip_protocol"); ok {
		request["IpProtocol"] = v
	}
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var forwardEntryNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		forwardEntryNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_forward_entries", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.ForwardTableEntries.ForwardTableEntry", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ForwardTableEntries.ForwardTableEntry", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if forwardEntryNameRegex != nil {
				if !forwardEntryNameRegex.MatchString(fmt.Sprint(item["ForwardEntryName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ForwardEntryId"])]; !ok {
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
			"external_ip":        object["ExternalIp"],
			"external_port":      object["ExternalPort"],
			"id":                 fmt.Sprint(object["ForwardEntryId"]),
			"forward_entry_id":   fmt.Sprint(object["ForwardEntryId"]),
			"forward_entry_name": object["ForwardEntryName"],
			"name":               object["ForwardEntryName"],
			"internal_ip":        object["InternalIp"],
			"internal_port":      object["InternalPort"],
			"ip_protocol":        object["IpProtocol"],
			"status":             object["Status"],
		}
		ids = append(ids, fmt.Sprint(object["ForwardEntryId"]))
		names = append(names, object["ForwardEntryName"])
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

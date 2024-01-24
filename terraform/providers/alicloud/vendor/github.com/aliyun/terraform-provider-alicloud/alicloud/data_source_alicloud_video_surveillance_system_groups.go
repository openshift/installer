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

func dataSourceAlicloudVideoSurveillanceSystemGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVideoSurveillanceSystemGroupsRead,
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
			"in_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"gb28181", "rtmp"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"callback": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"gb_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gb_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"in_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"out_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"play_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"push_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stats": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_num": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ied_num": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipc_num": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"platform_num": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVideoSurveillanceSystemGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeGroups"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("in_protocol"); ok {
		request["InProtocol"] = v
	}
	request["IncludeStats"] = true
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNum"] = 1
	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
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
	conn, err := client.NewVsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-12"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_video_surveillance_system_groups", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Groups", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Groups", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"app":          object["App"],
			"callback":     object["Callback"],
			"create_time":  object["CreatedTime"],
			"description":  object["Description"],
			"enabled":      object["Enabled"],
			"gb_id":        object["GbId"],
			"gb_ip":        object["GbIp"],
			"id":           fmt.Sprint(object["Id"]),
			"group_id":     fmt.Sprint(object["Id"]),
			"group_name":   object["Name"],
			"in_protocol":  object["InProtocol"],
			"out_protocol": object["OutProtocol"],
			"play_domain":  object["PlayDomain"],
			"push_domain":  object["PushDomain"],
			"stats": []map[string]interface{}{
				{
					"ipc_num":      object["IpcNum"],
					"device_num":   object["DeviceNum"],
					"ied_num":      object["IedNum"],
					"platform_num": object["PlatformNum"],
				}},
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
		names = append(names, object["Name"])
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

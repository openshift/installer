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

func dataSourceAlicloudVpcTrafficMirrorSessions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpcTrafficMirrorSessionsRead,
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Created", "Creating", "Deleting", "Modifying"}, false),
			},
			"traffic_mirror_filter_id": {
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
			"traffic_mirror_session_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"traffic_mirror_source_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"traffic_mirror_target_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sessions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"packet_length": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_mirror_filter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_mirror_session_business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_mirror_session_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_mirror_session_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_mirror_session_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_mirror_source_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"traffic_mirror_target_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_mirror_target_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_network_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpcTrafficMirrorSessionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTrafficMirrorSessions"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("enabled"); ok {
		request["Enabled"] = v
	}
	if v, ok := d.GetOk("priority"); ok {
		request["Priority"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("traffic_mirror_filter_id"); ok {
		request["TrafficMirrorFilterId"] = v
	}
	if v, ok := d.GetOk("traffic_mirror_session_name"); ok {
		request["TrafficMirrorSessionName"] = v
	}
	if v, ok := d.GetOk("traffic_mirror_source_id"); ok {
		request["TrafficMirrorSourceId"] = v
	}
	if v, ok := d.GetOk("traffic_mirror_target_id"); ok {
		request["TrafficMirrorTargetId"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var trafficMirrorSessionNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		trafficMirrorSessionNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_traffic_mirror_sessions", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.TrafficMirrorSessions", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TrafficMirrorSessions", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if trafficMirrorSessionNameRegex != nil && !trafficMirrorSessionNameRegex.MatchString(fmt.Sprint(item["TrafficMirrorSessionName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TrafficMirrorSessionId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["TrafficMirrorSessionStatus"].(string) {
				continue
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
			"enabled":                                object["Enabled"],
			"packet_length":                          formatInt(object["PacketLength"]),
			"priority":                               formatInt(object["Priority"]),
			"status":                                 object["TrafficMirrorSessionStatus"],
			"traffic_mirror_filter_id":               object["TrafficMirrorFilterId"],
			"traffic_mirror_session_business_status": object["TrafficMirrorSessionBusinessStatus"],
			"traffic_mirror_session_description":     object["TrafficMirrorSessionDescription"],
			"id":                                     fmt.Sprint(object["TrafficMirrorSessionId"]),
			"traffic_mirror_session_id":              fmt.Sprint(object["TrafficMirrorSessionId"]),
			"traffic_mirror_session_name":            object["TrafficMirrorSessionName"],
			"traffic_mirror_source_ids":              object["TrafficMirrorSourceIds"],
			"traffic_mirror_target_id":               object["TrafficMirrorTargetId"],
			"traffic_mirror_target_type":             object["TrafficMirrorTargetType"],
			"virtual_network_id":                     formatInt(object["VirtualNetworkId"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["TrafficMirrorSessionName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("sessions", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

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

func dataSourceAlicloudEciVirtualNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEciVirtualNodesRead,
		Schema: map[string]*schema.Schema{
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Cleaned", "Failed", "Pending", "Ready"}, false),
			},
			"tags": tagsSchema(),
			"vswitch_id": {
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
			"virtual_node_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eni_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"events": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"first_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_timestamp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"reason": {
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
						"internet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intranet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ram_role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
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
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_node_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEciVirtualNodesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeVirtualNodes"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("virtual_node_name"); ok {
		request["VirtualNodeName"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	var objects []map[string]interface{}
	var virtualNodeNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		virtualNodeNameRegex = r
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
	if v, ok := d.GetOk("tags"); ok {
		index := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", index)] = key
			request[fmt.Sprintf("Tag.%d.Value", index)] = value.(string)
			index = index + 1
		}
	}
	var response map[string]interface{}
	conn, err := client.NewEciClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-08"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_eci_virtual_nodes", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.VirtualNodes", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.VirtualNodes", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if virtualNodeNameRegex != nil && !virtualNodeNameRegex.MatchString(fmt.Sprint(item["VirtualNodeName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["VirtualNodeId"])]; !ok {
					continue
				}
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
			"cpu":               formatInt(object["Cpu"]),
			"create_time":       object["CreationTime"],
			"eni_instance_id":   object["EniInstanceId"],
			"internet_ip":       object["InternetIp"],
			"intranet_ip":       object["IntranetIp"],
			"memory":            formatInt(object["Memory"]),
			"ram_role_name":     object["RamRoleName"],
			"resource_group_id": object["ResourceGroupId"],
			"security_group_id": object["SecurityGroupId"],
			"status":            object["Status"],
			"vswitch_id":        object["VSwitchId"],
			"id":                fmt.Sprint(object["VirtualNodeId"]),
			"virtual_node_id":   fmt.Sprint(object["VirtualNodeId"]),
			"virtual_node_name": object["VirtualNodeName"],
			"vpc_id":            object["VpcId"],
			"zone_id":           object["ZoneId"],
		}

		events := make([]map[string]interface{}, 0)
		if eventsList, ok := object["Events"].([]interface{}); ok {
			for _, v := range eventsList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"count":           formatInt(m1["Count"]),
						"first_timestamp": m1["FirstTimestamp"],
						"last_timestamp":  m1["LastTimestamp"],
						"message":         m1["Message"],
						"name":            m1["Name"],
						"reason":          m1["Reason"],
						"type":            m1["Type"],
					}
					events = append(events, temp1)
				}
			}
		}
		mapping["events"] = events

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["Key"].(string)
				value := t.(map[string]interface{})["Value"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["VirtualNodeName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("nodes", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

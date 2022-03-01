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

func dataSourceAlicloudEcsNetworkInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsNetworkInterfacesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
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
			"network_interface_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'name' has been deprecated from provider version 1.123.1. New field 'network_interface_name' instead",
			},
			"primary_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"private_ip": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'private_ip' has been deprecated from provider version 1.123.1. New field 'primary_ip_address' instead",
			},
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
			"service_managed": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Attaching", "Available", "CreateFailed", "Creating", "Deleting", "Detaching", "InUse", "Linked", "Linking", "Unlinking"}, false),
			},
			"tags": tagsSchema(),
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Primary", "Secondary"}, false),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"private_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"queue_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"service_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_managed": {
							Type:     schema.TypeBool,
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
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
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcsNetworkInterfacesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeNetworkInterfaces"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("network_interface_name"); ok {
		request["NetworkInterfaceName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["NetworkInterfaceName"] = v
	}
	if v, ok := d.GetOk("primary_ip_address"); ok {
		request["PrimaryIpAddress"] = v
	} else if v, ok := d.GetOk("private_ip"); ok {
		request["PrimaryIpAddress"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}
	if v, ok := d.GetOkExists("service_managed"); ok {
		request["ServiceManaged"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag"] = tags
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var networkInterfaceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		networkInterfaceNameRegex = r
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
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_network_interfaces", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.NetworkInterfaceSets.NetworkInterfaceSet", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.NetworkInterfaceSets.NetworkInterfaceSet", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if networkInterfaceNameRegex != nil {
				if !networkInterfaceNameRegex.MatchString(fmt.Sprint(item["NetworkInterfaceName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["NetworkInterfaceId"])]; !ok {
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
			"description":            object["Description"],
			"instance_id":            object["InstanceId"],
			"mac":                    object["MacAddress"],
			"id":                     fmt.Sprint(object["NetworkInterfaceId"]),
			"network_interface_id":   fmt.Sprint(object["NetworkInterfaceId"]),
			"network_interface_name": object["NetworkInterfaceName"],
			"name":                   object["NetworkInterfaceName"],
			"primary_ip_address":     object["PrivateIpAddress"],
			"private_ip":             object["PrivateIpAddress"],
			"queue_number":           formatInt(object["QueueNumber"]),
			"resource_group_id":      object["ResourceGroupId"],
			"security_group_ids":     object["SecurityGroupIds"].(map[string]interface{})["SecurityGroupId"],
			"service_id":             formatInt(object["ServiceID"]),
			"service_managed":        object["ServiceManaged"],
			"status":                 object["Status"],
			"type":                   object["Type"],
			"vswitch_id":             object["VSwitchId"],
			"vpc_id":                 object["VpcId"],
			"zone_id":                object["ZoneId"],
			"creation_time":          object["CreationTime"],
			"security_groups":        object["SecurityGroupIds"].(map[string]interface{})["SecurityGroupId"],
		}
		if v := object["PrivateIpSets"].(map[string]interface{})["PrivateIpSet"].([]interface{}); len(v) > 0 {
			privateIps := make([]interface{}, 0, len(object["PrivateIpSets"].(map[string]interface{})["PrivateIpSet"].([]interface{})))
			for _, v := range object["PrivateIpSets"].(map[string]interface{})["PrivateIpSet"].([]interface{}) {
				if !v.(map[string]interface{})["Primary"].(bool) {
					privateIps = append(privateIps, v.(map[string]interface{})["PrivateIpAddress"])
				}
			}
			mapping["private_ip_addresses"] = privateIps
			mapping["private_ips"] = privateIps
		}

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.Tag", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["TagKey"].(string)
				value := t.(map[string]interface{})["TagValue"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags
		ids = append(ids, fmt.Sprint(object["NetworkInterfaceId"]))
		names = append(names, object["NetworkInterfaceName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("interfaces", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

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

func dataSourceAlicloudEcsKeyPairs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsKeyPairsRead,
		Schema: map[string]*schema.Schema{
			"finger_print": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pairs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"finger_print": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_pair_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Resource{Schema: outputInstancesSchema()},
						},
					},
				},
			},
			"key_pairs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"finger_print": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_pair_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Resource{Schema: outputInstancesSchema()},
						},
					},
				},
				Deprecated: "Field 'key_pairs' has been deprecated from provider version 1.121.0. New field 'pairs' instead.",
			},
		},
	}
}

func dataSourceAlicloudEcsKeyPairsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeKeyPairs"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("finger_print"); ok {
		request["KeyPairFingerPrint"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
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
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var keyPairNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		keyPairNameRegex = r
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
	objectMap := make(map[string][]map[string]interface{})
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_key_pairs", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.KeyPairs.KeyPair", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.KeyPairs.KeyPair", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if keyPairNameRegex != nil {
				if !keyPairNameRegex.MatchString(fmt.Sprint(item["KeyPairName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["KeyPairName"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
			objectMap[item["KeyPairName"].(string)] = make([]map[string]interface{}, 0)

		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	action = "DescribeInstances"
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_key_pairs", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Instances.Instance", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Instances.Instance", response)
		}
		result, _ := resp.([]interface{})

		for _, inst := range result {
			itemInst := inst.(map[string]interface{})
			if itemInst["KeyPairName"] != nil {
				if _, ok := objectMap[itemInst["KeyPairName"].(string)]; ok {
					publicIp := itemInst["EipAddress"].(map[string]interface{})["IpAddress"]
					if publicIp == "" && len(itemInst["PublicIpAddress"].(map[string]interface{})["IpAddress"].([]interface{})) > 0 {
						publicIp = itemInst["PublicIpAddress"].(map[string]interface{})["IpAddress"].([]interface{})[0]
					}
					var privateIp string
					if len(itemInst["InnerIpAddress"].(map[string]interface{})["IpAddress"].([]interface{})) > 0 {
						privateIp = itemInst["InnerIpAddress"].(map[string]interface{})["IpAddress"].([]interface{})[0].(string)

					} else if len(itemInst["VpcAttributes"].(map[string]interface{})["PrivateIpAddress"].(map[string]interface{})["IpAddress"].([]interface{})) > 0 {
						privateIp = itemInst["VpcAttributes"].(map[string]interface{})["PrivateIpAddress"].(map[string]interface{})["IpAddress"].([]interface{})[0].(string)
					}
					mapping := map[string]interface{}{
						"availability_zone": itemInst["ZoneId"],
						"instance_id":       itemInst["InstanceId"],
						"instance_name":     itemInst["InstanceName"],
						"vswitch_id":        itemInst["VpcAttributes"].(map[string]interface{})["VSwitchId"],
						"public_ip":         publicIp,
						"private_ip":        privateIp,
					}
					val := objectMap[itemInst["KeyPairName"].(string)]
					val = append(val, mapping)
					objectMap[itemInst["KeyPairName"].(string)] = val
				}
			}
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
			"finger_print":      object["KeyPairFingerPrint"],
			"id":                fmt.Sprint(object["KeyPairName"]),
			"key_pair_name":     fmt.Sprint(object["KeyPairName"]),
			"key_name":          object["KeyPairName"],
			"resource_group_id": object["ResourceGroupId"],
			"instances":         objectMap[object["KeyPairName"].(string)],
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
		ids = append(ids, fmt.Sprint(object["KeyPairName"]))
		names = append(names, object["KeyPairName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("key_pairs", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("pairs", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCddcDedicatedHostGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCddcDedicatedHostGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"engine": {
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocation_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bastion_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_allocate_ration": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"cpu_allocated_amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"cpu_allocation_ratio": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_count_group_by_host_type": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"place_holder": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"dedicated_host_group_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deploy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_allocate_ration": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"disk_allocated_amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"disk_allocation_ratio": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_used_amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"disk_utility": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_replace_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem_allocate_ration": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"mem_allocated_amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"mem_allocation_ratio": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem_used_amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"mem_utility": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"text": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_id_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceAlicloudCddcDedicatedHostGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDedicatedHostGroups"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("engine"); ok {
		request["Engine"] = v
	}
	request["RegionId"] = client.RegionId
	var objects []map[string]interface{}
	var userNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		userNameRegex = r
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
	conn, err := client.NewCddcClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-20"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cddc_dedicated_host_groups", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.DedicatedHostGroups.DedicatedHostGroups", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DedicatedHostGroups.DedicatedHostGroups", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if userNameRegex != nil && !userNameRegex.MatchString(fmt.Sprint(item["DedicatedHostGroupDesc"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["DedicatedHostGroupId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"allocation_policy":         object["AllocationPolicy"],
			"bastion_instance_id":       object["BastionInstanceId"],
			"cpu_allocate_ration":       object["CpuAllocateRation"],
			"cpu_allocated_amount":      object["CpuAllocatedAmount"],
			"cpu_allocation_ratio":      formatInt(object["CpuAllocationRatio"]),
			"create_time":               object["CreateTime"],
			"dedicated_host_group_desc": object["DedicatedHostGroupDesc"],
			"id":                        fmt.Sprint(object["DedicatedHostGroupId"]),
			"dedicated_host_group_id":   fmt.Sprint(object["DedicatedHostGroupId"]),
			"deploy_type":               object["DeployType"],
			"disk_allocate_ration":      object["DiskAllocateRation"],
			"disk_allocated_amount":     object["DiskAllocatedAmount"],
			"disk_allocation_ratio":     formatInt(object["DiskAllocationRatio"]),
			"disk_used_amount":          object["DiskUsedAmount"],
			"disk_utility":              object["DiskUtility"],
			"engine":                    object["Engine"],
			"host_number":               formatInt(object["HostNumber"]),
			"host_replace_policy":       object["HostReplacePolicy"],
			"instance_number":           formatInt(object["InstanceNumber"]),
			"mem_allocate_ration":       object["MemAllocateRation"],
			"mem_allocated_amount":      object["MemAllocatedAmount"],
			"mem_allocation_ratio":      formatInt(object["MemAllocationRatio"]),
			"mem_used_amount":           object["MemUsedAmount"],
			"mem_utility":               object["MemUtility"],
			"text":                      object["Text"],
			"vpc_id":                    object["VPCId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["DedicatedHostGroupDesc"])
		s = append(s, mapping)
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

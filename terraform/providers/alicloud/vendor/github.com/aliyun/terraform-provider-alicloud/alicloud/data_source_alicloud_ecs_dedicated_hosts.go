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

func dataSourceAlicloudEcsDedicatedHosts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsDedicatedHostsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"dedicated_host_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"dedicated_host_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dedicated_host_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"operation_locks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lock_reason": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Available", "Creating", "PermanentFailure", "Released", "UnderAssessment"}, false),
			},
			"tags": tagsSchema(),
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_on_maintenance": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_placement": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_release_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"capacity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"available_local_storage": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"available_memory": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"available_vcpus": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"available_vgpus": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"local_storage_category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"total_local_storage": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"total_memory": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"total_vcpus": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"total_vgpus": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"cores": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_over_commit_ratio": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_type": {
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
						"gpu_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"machine_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_attributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"slb_udp_timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"udp_timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"operation_locks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lock_reason": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"physical_gpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sale_cycle": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sockets": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"supported_custom_instance_type_families": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"supported_instance_type_families": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"supported_instance_types_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"tags": {
							Type:     schema.TypeMap,
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

func dataSourceAlicloudEcsDedicatedHostsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDedicatedHosts"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("dedicated_host_name"); ok {
		request["DedicatedHostName"] = v
	}
	if v, ok := d.GetOk("dedicated_host_type"); ok {
		request["DedicatedHostType"] = v
	}
	if v, ok := d.GetOk("dedicated_host_id"); ok {
		request["DedicatedHostId"] = v
	}
	if v, ok := d.GetOk("operation_locks"); ok {
		operationLocksMaps := make([]map[string]interface{}, 0)
		for _, operationLocks := range v.([]interface{}) {
			operationLocksMap := make(map[string]interface{})
			operationLocksArg := operationLocks.(map[string]interface{})
			operationLocksMap["LockReason"] = operationLocksArg["lock_reason"]
			operationLocksMaps = append(operationLocksMaps, operationLocksMap)
		}
		request["OperationLocks"] = operationLocksMaps

	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
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
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var dedicatedHostNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		dedicatedHostNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_dedicated_hosts", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DedicatedHosts.DedicatedHost", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DedicatedHosts.DedicatedHost", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if dedicatedHostNameRegex != nil {
				if !dedicatedHostNameRegex.MatchString(fmt.Sprint(item["DedicatedHostName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DedicatedHostId"])]; !ok {
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
			"action_on_maintenance": object["ActionOnMaintenance"],
			"auto_placement":        object["AutoPlacement"],
			"auto_release_time":     object["AutoReleaseTime"],
			"cores":                 formatInt(object["Cores"]),
			"cpu_over_commit_ratio": object["CpuOverCommitRatio"],
			"id":                    fmt.Sprint(object["DedicatedHostId"]),
			"dedicated_host_id":     fmt.Sprint(object["DedicatedHostId"]),
			"dedicated_host_name":   object["DedicatedHostName"],
			"dedicated_host_type":   object["DedicatedHostType"],
			"description":           object["Description"],
			"expired_time":          object["ExpiredTime"],
			"gpu_spec":              object["GPUSpec"],
			"machine_id":            object["MachineId"],
			"payment_type":          object["ChargeType"],
			"physical_gpus":         formatInt(object["PhysicalGpus"]),
			"resource_group_id":     object["ResourceGroupId"],
			"sale_cycle":            object["SaleCycle"],
			"sockets":               formatInt(object["Sockets"]),
			"status":                object["Status"],
			"supported_custom_instance_type_families": object["SupportedCustomInstanceTypeFamilies"].(map[string]interface{})["SupportedCustomInstanceTypeFamily"],
			"supported_instance_type_families":        object["SupportedInstanceTypeFamilies"].(map[string]interface{})["SupportedInstanceTypeFamily"],
			"supported_instance_types_list":           object["SupportedInstanceTypesList"].(map[string]interface{})["SupportedInstanceTypesList"],
			"zone_id":                                 object["ZoneId"],
		}
		operationLock := make([]map[string]interface{}, 0)
		if operationLockList, ok := object["OperationLocks"].(map[string]interface{})["OperationLock"].([]interface{}); ok {
			for _, v := range operationLockList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"lock_reason": m1["LockReason"],
					}
					operationLock = append(operationLock, temp1)
				}
			}
		}
		mapping["operation_locks"] = operationLock

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
		ids = append(ids, fmt.Sprint(object["DedicatedHostId"]))
		names = append(names, object["DedicatedHostName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("hosts", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

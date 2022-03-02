package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCddcDedicatedHosts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCddcDedicatedHostsRead,
		Schema: map[string]*schema.Schema{
			"allocation_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Allocatable", "Suspended"}, false),
			},
			"dedicated_host_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"host_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"dhg_local_ssd", "dhg_cloud_ssd"}, false),
			},
			"order_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1", "2", "3", "4", "5", "6", "7"}, false),
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
						"bastion_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_allocation_ratio": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_used": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_group_id": {
							Type:     schema.TypeString,
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
						"disk_allocation_ratio": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ecs_class_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_cpu": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_mem": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_storage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mem_allocation_ratio": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_used": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"open_permission": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allocation_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_used": {
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudCddcDedicatedHostsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDedicatedHosts"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("allocation_status"); ok {
		request["AllocationStatus"] = v
	}
	request["DedicatedHostGroupId"] = d.Get("dedicated_host_group_id")
	if v, ok := d.GetOk("status"); ok {
		request["HostStatus"] = v
	}
	if v, ok := d.GetOk("host_type"); ok {
		request["HostType"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Tags"] = respJson
	}
	if v, ok := d.GetOk("order_id"); ok {
		request["OrderId"] = v
	}
	request["PageNumbers"] = 1
	request["PageSize"] = PageSizeLarge
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	var objects []map[string]interface{}

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
	for {
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cddc_dedicated_hosts", action, AlibabaCloudSdkGoERROR)
		}
		if _, ok := response["DedicatedHosts"]; !ok {
			break
		}
		resp, err := jsonpath.Get("$.DedicatedHosts.DedicatedHosts", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DedicatedHosts.DedicatedHosts", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DedicatedHostGroupId"], ":", item["DedicatedHostId"])]; !ok {
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
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"bastion_instance_id":     object["BastionInstanceId"],
			"cpu_allocation_ratio":    object["CPUAllocationRatio"],
			"cpu_used":                object["CpuUsed"],
			"create_time":             object["CreatedTime"],
			"dedicated_host_group_id": object["DedicatedHostGroupId"],
			"id":                      fmt.Sprint(object["DedicatedHostGroupId"], ":", object["DedicatedHostId"]),
			"dedicated_host_id":       fmt.Sprint(object["DedicatedHostId"]),
			"disk_allocation_ratio":   object["DiskAllocationRatio"],
			"ecs_class_code":          object["EcsClassCode"],
			"end_time":                object["EndTime"],
			"engine":                  object["Engine"],
			"host_class":              object["HostClass"],
			"host_cpu":                object["HostCPU"],
			"host_mem":                object["HostMem"],
			"host_name":               object["HostName"],
			"host_storage":            object["HostStorage"],
			"host_type":               object["HostType"],
			"image_category":          object["ImageCategory"],
			"ip_address":              object["IPAddress"],
			"mem_allocation_ratio":    object["MemAllocationRatio"],
			"memory_used":             object["MemoryUsed"],
			"open_permission":         object["OpenPermission"],
			"allocation_status":       convertCddcAllocationStatusResponse(object["AllocationStatus"].(string)),
			"status":                  object["HostStatus"],
			"storage_used":            object["StorageUsed"],
			"vswitch_id":              object["VSwitchId"],
			"vpc_id":                  object["VPCId"],
			"zone_id":                 object["ZoneId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["DedicatedHostGroupId"], ":", object["DedicatedHostId"])
		cddcService := CddcService{client}
		getResp, err := cddcService.DescribeCddcDedicatedHost(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["expired_time"] = getResp["ExpiredTime"]
		getResp1, err := cddcService.ListTagResources(id, "DEDICATEDHOST")
		if err != nil {
			return WrapError(err)
		}
		mapping["tags"] = tagsToMap(getResp1)

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
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

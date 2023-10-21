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

func dataSourceAlicloudEcsDedicatedHostClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsDedicatedHostClustersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"dedicated_host_cluster_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"dedicated_host_cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
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
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedicated_host_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dedicated_host_cluster_capacity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"total_memory": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"available_memory": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"available_vcpus": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"total_vcpus": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"local_storage_capacities": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"available_disk": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"data_disk_category": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"total_disk": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
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

func dataSourceAlicloudEcsDedicatedHostClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDedicatedHostClusters"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("dedicated_host_cluster_ids"); ok {
		request["DedicatedHostClusterIds"] = convertListToJsonString(v.([]interface{}))
	}
	if v, ok := d.GetOk("dedicated_host_cluster_name"); ok {
		request["DedicatedHostClusterName"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var dedicatedHostClusterNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		dedicatedHostClusterNameRegex = r
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
	tagsMap := make(map[string]interface{})
	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		tagsMap = v.(map[string]interface{})
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
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_dedicated_host_clusters", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DedicatedHostClusters.DedicatedHostCluster", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DedicatedHostClusters.DedicatedHostCluster", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if dedicatedHostClusterNameRegex != nil && !dedicatedHostClusterNameRegex.MatchString(fmt.Sprint(item["DedicatedHostClusterName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DedicatedHostClusterId"])]; !ok {
					continue
				}
			}
			if len(tagsMap) > 0 {
				if _, ok := item["Tags"].(map[string]interface{}); !ok {
					continue
				}
				tagsItem := tagsToMap(item["Tags"].(map[string]interface{})["Tag"])
				if len(tagsItem) != len(tagsMap) {
					continue
				}
				match := true
				for key, value := range tagsItem {
					if v, ok := tagsMap[key]; !ok || v.(string) != value {
						match = false
						break
					}
				}
				if !match {
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
			"id":                          fmt.Sprint(object["DedicatedHostClusterId"]),
			"dedicated_host_cluster_id":   fmt.Sprint(object["DedicatedHostClusterId"]),
			"dedicated_host_cluster_name": object["DedicatedHostClusterName"],
			"description":                 object["Description"],
			"resource_group_id":           object["ResourceGroupId"],
			"zone_id":                     object["ZoneId"],
			"dedicated_host_ids":          object["DedicatedHostIds"].(map[string]interface{})["DedicatedHostId"],
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

		dedicatedHostClusterCapacity := make([]map[string]interface{}, 0)
		if dedicatedHostClusterCapacityArr, ok := object["DedicatedHostClusterCapacity"].(map[string]interface{}); ok {
			temp1 := map[string]interface{}{
				"total_memory":     formatInt(dedicatedHostClusterCapacityArr["TotalMemory"]),
				"available_memory": formatInt(dedicatedHostClusterCapacityArr["AvailableMemory"]),
				"available_vcpus":  formatInt(dedicatedHostClusterCapacityArr["AvailableVcpus"]),
				"total_vcpus":      formatInt(dedicatedHostClusterCapacityArr["TotalVcpus"]),
			}

			localStorageCapacity := make([]map[string]interface{}, 0)
			if localStorageCapacitiesArr, ok := dedicatedHostClusterCapacityArr["LocalStorageCapacities"].(map[string]interface{}); ok {
				if localStorageCapacityList, ok := localStorageCapacitiesArr["LocalStorageCapacity"].([]interface{}); ok {
					for _, v := range localStorageCapacityList {
						if m1, ok := v.(map[string]interface{}); ok {
							temp1 := map[string]interface{}{
								"available_disk":     formatInt(m1["AvailableDisk"]),
								"data_disk_category": m1["DataDiskCategory"],
								"total_disk":         formatInt(m1["TotalDisk"]),
							}
							localStorageCapacity = append(localStorageCapacity, temp1)
						}
					}
				}
			}
			temp1["local_storage_capacities"] = localStorageCapacity

			dedicatedHostClusterCapacity = append(dedicatedHostClusterCapacity, temp1)
		}
		mapping["dedicated_host_cluster_capacity"] = dedicatedHostClusterCapacity

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["DedicatedHostClusterName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

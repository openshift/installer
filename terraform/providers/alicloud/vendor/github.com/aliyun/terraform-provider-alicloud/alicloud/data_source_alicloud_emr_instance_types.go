package alicloud

import (
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/emr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEmrInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEmrInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"destination_resource": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Zone",
					"Network",
					"InstanceType",
					"SystemDisk",
					"DataDisk",
				}, false),
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"support_local_storage": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PostPaid",
					"PrePaid",
				}, false),
			},
			"support_node_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_storage_capacity": {
							Type:     schema.TypeInt,
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

func dataSourceAlicloudEmrInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := emr.CreateListEmrAvailableResourceRequest()
	request.RegionId = string(client.Region)
	if dstRes, ok := d.GetOk("destination_resource"); ok {
		request.DestinationResource = strings.TrimSpace(dstRes.(string))
	}
	if typ, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = strings.TrimSpace(typ.(string))
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = strings.TrimSpace(instanceType.(string))
	}
	if chargeType, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = strings.TrimSpace(chargeType.(string))
	}

	raw, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
		return emrClient.ListEmrAvailableResource(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_emr_instance_types",
			request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	supportedResources := make(map[string][]emr.SupportedResource)
	resourceResponse, _ := raw.(*emr.ListEmrAvailableResourceResponse)
	for _, zoneInfo := range resourceResponse.EmrZoneInfoList.EmrZoneInfo {
		resourceInfo := zoneInfo.EmrResourceInfoList.EmrResourceInfo
		if len(resourceInfo) == 1 {
			supportedResources[zoneInfo.ZoneId] = resourceInfo[0].SupportedResourceList.SupportedResource
		}
	}

	return emrClusterInstanceTypesAttributes(d, supportedResources)
}

func emrClusterInstanceTypesAttributes(d *schema.ResourceData,
	supportedResources map[string][]emr.SupportedResource) error {
	var ids []string
	var zoneIDs []string
	var s []map[string]interface{}

	optZoneID := d.Get("zone_id").(string)

	for k, v := range supportedResources {

		if k == "" || len(v) == 0 { // ignore empty zoneId or empty emr instance type of the specific zoneId
			continue
		} else if optZoneID != "" && k != optZoneID { // get supported resources of a specific zoneId
			continue
		} else if optZoneID != "" && k == optZoneID {
			zoneIDs = append(zoneIDs, k)
			break
		}

		zoneIDs = append(zoneIDs, k)
	}
	sort.Strings(zoneIDs)
	localStorage := d.Get("support_local_storage").(bool)
	supportNodeType := d.Get("support_node_type").([]interface{})
	nodeTypeFilter := func(filter []interface{}, source []string) bool {
		if len(source) == 0 {
			return false
		}
		sourceMapping := make(map[string]struct{})
		for _, s := range source {
			sourceMapping[s] = struct{}{}
		}
		for _, f := range filter {
			if _, ok := sourceMapping[f.(string)]; !ok {
				return false
			}
		}
		return true
	}

	for _, zoneID := range zoneIDs {
		mapping := map[string]interface{}{
			"zone_id": zoneID,
		}
		if v, ok := supportedResources[zoneID]; ok {
			selector := make(map[string]interface{})
			selectFunc := func(m map[string]interface{}, tpe emr.SupportedResource, localStorage bool) {
				prefix := "ecs.i2." // use those instance types first to localStorage
				if localStorage {
					prefix = "ecs.g5." // use those instance types first to cloudStorage
				}
				if preID, ok := selector["id"]; !ok {
					selector["id"] = tpe.EmrInstanceType.InstanceType
					selector["cpu_core_count"] = tpe.EmrInstanceType.CpuCoreCount
					if localStorage {
						selector["local_storage_capacity"] = tpe.EmrInstanceType.LocalStorageCapacity
					}
				} else if !strings.HasPrefix(tpe.EmrInstanceType.InstanceType, prefix) {
					if !strings.HasPrefix(preID.(string), prefix) &&
						selector["cpu_core_count"].(int) > tpe.EmrInstanceType.CpuCoreCount {
						selector["id"] = tpe.EmrInstanceType.InstanceType
						selector["cpu_core_count"] = tpe.EmrInstanceType.CpuCoreCount
						if localStorage {
							selector["local_storage_capacity"] = tpe.EmrInstanceType.LocalStorageCapacity
						}
					}
					return
				} else if !strings.HasPrefix(preID.(string), prefix) {
					selector["id"] = tpe.EmrInstanceType.InstanceType
					selector["cpu_core_count"] = tpe.EmrInstanceType.CpuCoreCount
					if localStorage {
						selector["local_storage_capacity"] = tpe.EmrInstanceType.LocalStorageCapacity
					}
				} else if selector["cpu_core_count"].(int) > tpe.EmrInstanceType.CpuCoreCount {
					selector["id"] = tpe.EmrInstanceType.InstanceType
					selector["cpu_core_count"] = tpe.EmrInstanceType.CpuCoreCount
					if localStorage {
						selector["local_storage_capacity"] = tpe.EmrInstanceType.LocalStorageCapacity
					}
				}
			}

			for _, tpe := range v {
				if nodeTypeFilter(supportNodeType, tpe.SupportNodeTypeList.SupportNodeType) == false {
					continue
				}

				if localStorage == true && tpe.EmrInstanceType.LocalStorageAmount > 0 {
					selectFunc(selector, tpe, true)
				} else if localStorage == false && tpe.EmrInstanceType.LocalStorageAmount == 0 {
					selectFunc(selector, tpe, false)
				}
			}
			if v, ok := selector["id"]; !ok {
				continue
			} else {
				mapping["id"] = v
			}
			if v, ok := selector["local_storage_capacity"]; ok {
				mapping["local_storage_capacity"] = v
			}
			ids = append(ids, selector["id"].(string))
			s = append(s, mapping)
		}
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("types", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

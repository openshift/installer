package alicloud

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDBZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDBZonesRead,

		Schema: map[string]*schema.Schema{
			"multi": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"multi_zone": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{"PrePaid", "PostPaid"}, false),
			},
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "SQLServer", "PostgreSQL", "PPAS", "MariaDB"}, false),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "HighAvailability", "AlwaysOn", "Finance"}, false),
			},
			"db_instance_storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_ssd", "local_ssd", "cloud_essd", "cloud_essd2", "cloud_essd3"}, false),
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
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"multi_zone_ids": {
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

func dataSourceAlicloudDBZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	engines := make([]string, 0)
	if v, ok := d.GetOk("engine"); ok && v.(string) != "" {
		engines = append(engines, v.(string))
	} else {
		engines = []string{"MySQL", "SQLServer", "PostgreSQL", "PPAS", "MariaDB"}
	}

	action := "DescribeAvailableZones"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"SourceIp": client.SourceIp,
	}
	if v, ok := d.GetOk("engine_version"); ok && v.(string) != "" {
		request["EngineVersion"] = v.(string)
	}
	if v, ok := d.GetOk("zone_id"); ok && v.(string) != "" {
		request["ZoneId"] = v.(string)
	}
	instanceChargeType := d.Get("instance_charge_type").(string)
	if instanceChargeType == string(PostPaid) {
		request["CommodityCode"] = "bards"
	} else {
		request["CommodityCode"] = "rds"
	}
	multiZone := false
	if v, ok := d.GetOkExists("multi_zone"); ok {
		multiZone = v.(bool)
	} else if v, ok := d.GetOkExists("multi"); ok {
		multiZone = v.(bool)
	}
	var targetCategory, targetStorageType string
	if v, ok := d.GetOk("category"); ok && v.(string) != "" {
		targetCategory = v.(string)
	}
	if v, ok := d.GetOk("db_instance_storage_type"); ok && v.(string) != "" {
		targetStorageType = v.(string)
	}
	var ids []string
	var s []map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	for _, engine := range engines {
		request["Engine"] = engine
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_zones", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AvailableZones", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AvailableZones", response)
		}

		for _, r := range resp.([]interface{}) {
			availableZoneItem := r.(map[string]interface{})

			zoneId := fmt.Sprint(availableZoneItem["ZoneId"])
			if (multiZone && !strings.Contains(zoneId, MULTI_IZ_SYMBOL)) || (!multiZone && strings.Contains(zoneId, MULTI_IZ_SYMBOL)) {
				continue
			}

			if targetCategory == "" && targetStorageType == "" {
				ids = append(ids, zoneId)
				continue
			}
			for _, r := range availableZoneItem["SupportedEngines"].([]interface{}) {
				supportedEngineItem := r.(map[string]interface{})
				for _, r := range supportedEngineItem["SupportedEngineVersions"].([]interface{}) {
					supportedEngineVersionItem := r.(map[string]interface{})
					for _, r := range supportedEngineVersionItem["SupportedCategorys"].([]interface{}) {
						supportedCategoryItem := r.(map[string]interface{})
						if targetCategory != "" && targetCategory != fmt.Sprint(supportedCategoryItem["Category"]) {
							continue
						}
						if targetStorageType == "" {
							ids = append(ids, zoneId)
							goto NEXT
						}
						for _, r := range supportedCategoryItem["SupportedStorageTypes"].([]interface{}) {
							supportedStorageTypeItem := r.(map[string]interface{})
							if targetStorageType != fmt.Sprint(supportedStorageTypeItem["StorageType"]) {
								continue
							}
							ids = append(ids, zoneId)
							goto NEXT
						}
					}
				}
			}
		NEXT:
			continue
		}
	}
	if len(ids) > 0 {
		sort.Strings(ids)
	}

	for _, zoneId := range ids {
		mapping := map[string]interface{}{
			"id":             zoneId,
			"multi_zone_ids": splitMultiZoneId(zoneId),
		}
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

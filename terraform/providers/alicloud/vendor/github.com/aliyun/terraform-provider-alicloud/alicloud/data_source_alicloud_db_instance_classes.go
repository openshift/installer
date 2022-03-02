package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDBInstanceClasses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDBInstanceClassesRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"sorted_by": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Price"}, false),
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
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
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_ssd", "local_ssd", "cloud_essd", "cloud_essd2", "cloud_essd3"}, false),
			},
			"db_instance_storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_ssd", "local_ssd", "cloud_essd", "cloud_essd2", "cloud_essd3"}, false),
			},
			"multi_zone": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			// Computed values.
			"instance_classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_ids": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sub_zone_ids": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Computed: true,
									},
								},
							},
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"price": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_range": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"max": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"step": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceAlicloudDBInstanceClassesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	instanceChargeType := d.Get("instance_charge_type").(string)
	if instanceChargeType == string(PostPaid) {
		instanceChargeType = string(Postpaid)
	} else {
		instanceChargeType = string(Prepaid)
	}
	zoneId, zoneIdOk := d.GetOk("zone_id")
	engine, engineOk := d.GetOk("engine")
	engineVersion, engineVersionOk := d.GetOk("engine_version")
	dbInstanceClass, dbInstanceClassOk := d.GetOk("db_instance_class")
	dbInstanceStorageType, dbInstanceStorageTypeOk := d.GetOk("db_instance_storage_type")
	if !dbInstanceStorageTypeOk || dbInstanceStorageType.(string) == "" {
		dbInstanceStorageType, dbInstanceStorageTypeOk = d.GetOk("storage_type")
	}
	category, categoryOk := d.GetOk("category")

	availableZones := make([]map[string]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	// if all filters can be got, there is no need to invoking DescribeAvailableZones to get them
	if zoneIdOk && zoneId.(string) != "" &&
		engineOk && engine.(string) != "" &&
		engineVersionOk && engineVersion.(string) != "" &&
		dbInstanceStorageTypeOk && dbInstanceStorageType.(string) != "" &&
		categoryOk && category.(string) != "" {

		action := "DescribeAvailableClasses"
		request := map[string]interface{}{
			"RegionId":              client.RegionId,
			"SourceIp":              client.SourceIp,
			"ZoneId":                zoneId,
			"InstanceChargeType":    instanceChargeType,
			"Engine":                engine,
			"EngineVersion":         engineVersion,
			"DBInstanceStorageType": dbInstanceStorageType,
			"Category":              category,
		}
		var response map[string]interface{}
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_instance_classes", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DBInstanceClasses", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DBInstanceClasses", response)
		}
		zoneIds := make([]map[string]interface{}, 0)
		zoneIds = append(zoneIds, map[string]interface{}{
			"id":           zoneId,
			"sub_zone_ids": splitMultiZoneId(zoneId.(string)),
		})
		for _, r := range resp.([]interface{}) {
			instanceClassItem := r.(map[string]interface{})
			if dbInstanceClassOk && dbInstanceClass != "" && dbInstanceClass != fmt.Sprint(instanceClassItem["DBInstanceClass"]) {
				continue
			}
			mapping := map[string]interface{}{
				"instance_class": fmt.Sprint(instanceClassItem["DBInstanceClass"]),
				"zone_ids":       zoneIds,
				"storage_range": map[string]interface{}{
					"min":  fmt.Sprint(instanceClassItem["DBInstanceStorageRange"].(map[string]interface{})["MinValue"]),
					"max":  fmt.Sprint(instanceClassItem["DBInstanceStorageRange"].(map[string]interface{})["MaxValue"]),
					"step": fmt.Sprint(instanceClassItem["DBInstanceStorageRange"].(map[string]interface{})["Step"]),
				},
			}
			s = append(s, mapping)
			ids = append(ids, fmt.Sprint(instanceClassItem["DBInstanceClass"]))
		}
	} else {
		// 1. Invkoing DescribeAvailableZones to get available zones and other filters
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
		if instanceChargeType == string(PostPaid) {
			request["CommodityCode"] = "bards"
		} else {
			request["CommodityCode"] = "rds"
		}
		multiZone := false
		if v, ok := d.GetOk("multi_zone"); ok {
			multiZone = v.(bool)
		}
		var targetCategory, targetStorageType string
		if v, ok := d.GetOk("category"); ok && v.(string) != "" {
			targetCategory = v.(string)
		}
		if v, ok := d.GetOk("db_instance_storage_type"); ok && v.(string) != "" {
			targetStorageType = v.(string)
		}
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
					availableZones = append(availableZones, availableZoneItem)
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
								availableZones = append(availableZones, availableZoneItem)
								goto NEXT
							}
							for _, r := range supportedCategoryItem["SupportedStorageTypes"].([]interface{}) {
								supportedStorageTypeItem := r.(map[string]interface{})
								if targetStorageType != fmt.Sprint(supportedStorageTypeItem["StorageType"]) {
									continue
								}
								availableZones = append(availableZones, availableZoneItem)
								goto NEXT
							}
						}
					}
				}
			NEXT:
				continue
			}
		}
		// 2. Iterating the availableZones and invoking DescribeAvailableClasses to get available classes
		for _, availableZone := range availableZones {
			action := "DescribeAvailableClasses"
			request := map[string]interface{}{
				"RegionId":           client.RegionId,
				"SourceIp":           client.SourceIp,
				"ZoneId":             fmt.Sprint(availableZone["ZoneId"]),
				"InstanceChargeType": instanceChargeType,
			}
			zoneIds := make([]map[string]interface{}, 0)
			zoneIds = append(zoneIds, map[string]interface{}{
				"id":           fmt.Sprint(availableZone["ZoneId"]),
				"sub_zone_ids": splitMultiZoneId(fmt.Sprint(availableZone["ZoneId"])),
			})
			for _, r := range availableZone["SupportedEngines"].([]interface{}) {
				supportedEngineItem := r.(map[string]interface{})
				request["Engine"] = fmt.Sprint(supportedEngineItem["Engine"])
				for _, r := range supportedEngineItem["SupportedEngineVersions"].([]interface{}) {
					supportedEngineVersionItem := r.(map[string]interface{})
					request["EngineVersion"] = fmt.Sprint(supportedEngineVersionItem["Version"])
					for _, r := range supportedEngineVersionItem["SupportedCategorys"].([]interface{}) {
						supportedCategoryItem := r.(map[string]interface{})
						request["Category"] = fmt.Sprint(supportedCategoryItem["Category"])
						for _, r := range supportedCategoryItem["SupportedStorageTypes"].([]interface{}) {
							storageTypeItem := r.(map[string]interface{})
							request["DBInstanceStorageType"] = fmt.Sprint(storageTypeItem["StorageType"])
							var response map[string]interface{}
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
								return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_instance_classes", action, AlibabaCloudSdkGoERROR)
							}
							resp, err := jsonpath.Get("$.DBInstanceClasses", response)
							if err != nil {
								return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DBInstanceClasses", response)
							}
							for _, r := range resp.([]interface{}) {
								instanceClassItem := r.(map[string]interface{})
								if dbInstanceClassOk && dbInstanceClass != "" && dbInstanceClass != fmt.Sprint(instanceClassItem["DBInstanceClass"]) {
									continue
								}
								mapping := map[string]interface{}{
									"instance_class": fmt.Sprint(instanceClassItem["DBInstanceClass"]),
									"zone_ids":       zoneIds,
									"storage_range": map[string]interface{}{
										"min":  fmt.Sprint(instanceClassItem["DBInstanceStorageRange"].(map[string]interface{})["MinValue"]),
										"max":  fmt.Sprint(instanceClassItem["DBInstanceStorageRange"].(map[string]interface{})["MaxValue"]),
										"step": fmt.Sprint(instanceClassItem["DBInstanceStorageRange"].(map[string]interface{})["Step"]),
									},
								}
								s = append(s, mapping)
								ids = append(ids, fmt.Sprint(instanceClassItem["DBInstanceClass"]))
							}
						}
					}
				}
			}
		}
	}

	d.SetId(dataResourceIdHash(ids))
	err := d.Set("instance_classes", s)
	if err != nil {
		return WrapError(err)
	}
	d.Set("ids", ids)
	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), s)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

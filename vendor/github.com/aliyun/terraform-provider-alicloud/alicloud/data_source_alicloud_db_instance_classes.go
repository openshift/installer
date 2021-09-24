package alicloud

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_ssd", "local_ssd"}, false),
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

	request := rds.CreateDescribeAvailableResourceRequest()
	request.RegionId = client.RegionId
	request.ZoneId = d.Get("zone_id").(string)
	instanceChargeType := d.Get("instance_charge_type").(string)
	multiZone := d.Get("multi_zone").(bool)
	if instanceChargeType == string(PostPaid) {
		instanceChargeType = string(Postpaid)
	} else {
		instanceChargeType = string(Prepaid)
	}
	request.InstanceChargeType = instanceChargeType
	var response = &rds.DescribeAvailableResourceResponse{}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeAvailableResource(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{Throttling}) {
				time.Sleep(time.Duration(5) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*rds.DescribeAvailableResourceResponse)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_instance_classes", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	type ClassInfosItem struct {
		Index        int
		StorageRange map[string]string
		ZoneIds      []map[string]interface{}
		Price        float64
	}

	classInfos := make(map[string]ClassInfosItem)
	indexMap := make(map[string]int)
	ids := []string{}

	engine, engineGot := d.GetOk("engine")
	engineVersion, engineVersionGot := d.GetOk("engine_version")
	dbInstanceClass, dbInstanceClassGot := d.GetOk("db_instance_class")
	storageType, storageTypeGot := d.GetOk("storage_type")
	category, categoryGot := d.GetOk("category")

	for _, AvailableZone := range response.AvailableZones.AvailableZone {
		id_item := []string{}
		if multiZone {
			if !strings.Contains(AvailableZone.ZoneId, "MAZ") {
				continue
			}
			for _, v := range splitMultiZoneId(AvailableZone.ZoneId) {
				id_item = append(id_item, v)
			}
		} else {
			if strings.Contains(AvailableZone.ZoneId, "MAZ") {
				continue
			}
			id_item = []string{AvailableZone.ZoneId}
		}

		zoneId := map[string]interface{}{
			"id":           AvailableZone.ZoneId,
			"sub_zone_ids": id_item,
		}

		ids = append(ids, AvailableZone.ZoneId)
		for _, SupportedEngine := range AvailableZone.SupportedEngines.SupportedEngine {
			if engineGot && engine.(string) != SupportedEngine.Engine {
				continue
			}
			ids = append(ids, SupportedEngine.Engine)
			for _, SupportedEngineVersion := range SupportedEngine.SupportedEngineVersions.SupportedEngineVersion {
				if engineVersionGot && engineVersion != SupportedEngineVersion.Version {
					continue
				}
				ids = append(ids, SupportedEngineVersion.Version)
				for _, SupportedCategory := range SupportedEngineVersion.SupportedCategorys.SupportedCategory {
					if categoryGot && category.(string) != SupportedCategory.Category {
						continue
					}
					for _, SupportedStorageType := range SupportedCategory.SupportedStorageTypes.SupportedStorageType {
						if storageTypeGot && storageType.(string) != SupportedStorageType.StorageType {
							continue
						}
						for _, AvailableResource := range SupportedStorageType.AvailableResources.AvailableResource {
							if dbInstanceClassGot && dbInstanceClass.(string) != AvailableResource.DBInstanceClass {
								continue
							}
							zoneIds := []map[string]interface{}{}
							if _, ok := classInfos[AvailableResource.DBInstanceClass]; ok {
								zoneIds = append(classInfos[AvailableResource.DBInstanceClass].ZoneIds, zoneId)
							} else {
								zoneIds = []map[string]interface{}{zoneId}
								indexMap[AvailableResource.DBInstanceClass] = len(classInfos)
							}
							classInfos[AvailableResource.DBInstanceClass] = ClassInfosItem{
								Index: indexMap[AvailableResource.DBInstanceClass],
								StorageRange: map[string]string{
									"min":  strconv.Itoa(AvailableResource.DBInstanceStorageRange.Min),
									"max":  strconv.Itoa(AvailableResource.DBInstanceStorageRange.Max),
									"step": strconv.Itoa(AvailableResource.DBInstanceStorageRange.Step),
								},
								ZoneIds: zoneIds,
							}
						}
					}
				}
			}
		}
	}

	infos := make([]map[string]interface{}, len(classInfos))
	var classIds []string
	for k, v := range classInfos {
		infos[v.Index] = map[string]interface{}{
			"zone_ids":       v.ZoneIds,
			"storage_range":  v.StorageRange,
			"instance_class": k,
		}
		classIds = append(classIds, k)
	}
	// The price will be query and assigned when the engine_version be specified and the sorted_by is "Price", Instead, it is assigned a null character.
	sortedBy := d.Get("sorted_by").(string)
	if engineVersionGot && len(infos) > 0 && sortedBy == "Price" {
		bssopenapiService := BssopenapiService{client}
		priceList, err := getDBInstanceClassPrice(bssopenapiService, instanceChargeType,
			infos, engineVersion.(string))
		if err != nil {
			return WrapError(err)
		}
		for i, info := range infos {
			info["price"] = fmt.Sprintf("%.4f", priceList[i])
		}
		sort.SliceStable(infos, func(i, j int) bool {
			iPrice, _ := strconv.ParseFloat(infos[i]["price"].(string), 64)
			jPrice, _ := strconv.ParseFloat(infos[j]["price"].(string), 64)
			return iPrice < jPrice
		})
	}

	d.SetId(dataResourceIdHash(ids))
	err = d.Set("instance_classes", infos)
	if err != nil {
		return WrapError(err)
	}
	d.Set("ids", classIds)
	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), infos)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func getDBInstanceClassPrice(bssopenapiService BssopenapiService, instanceChargeType string, infos []map[string]interface{}, version string) ([]float64, error) {
	client := bssopenapiService.client
	var modules interface{}
	var productType string
	moduleCode := "DBInstanceClass"
	var payAsYouGo []bssopenapi.GetPayAsYouGoPriceModuleList
	var subsciption []bssopenapi.GetSubscriptionPriceModuleList
	for _, info := range infos {
		config := fmt.Sprintf("DBInstanceClass:%s,EngineVersion:%s,Region:%s", info["instance_class"], version, client.Region)
		if instanceChargeType == string(Postpaid) {
			payAsYouGo = append(payAsYouGo, bssopenapi.GetPayAsYouGoPriceModuleList{
				ModuleCode: moduleCode,
				Config:     config,
				PriceType:  "Hour",
			})
		} else {
			subsciption = append(subsciption, bssopenapi.GetSubscriptionPriceModuleList{
				ModuleCode: moduleCode,
				Config:     config,
			})
		}
	}

	if len(payAsYouGo) != 0 {
		modules = payAsYouGo
		productType = "bards"
	} else {
		modules = subsciption
	}

	return bssopenapiService.GetInstanceTypePrice("rds", productType, modules)
}

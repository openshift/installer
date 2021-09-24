package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDBInstanceEngines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDBInstanceEnginesRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
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
			// Computed values.
			"instance_engines": {
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
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDBInstanceEnginesRead(d *schema.ResourceData, meta interface{}) error {
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_instance_engines", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	type engineInfosItem struct {
		Index         int
		Engine        string
		EngineVersion string
		ZoneIds       []map[string]interface{}
	}

	engineInfos := make(map[string]engineInfosItem)
	indexMap := make(map[string]int)
	ids := []string{}

	engine, engineGot := d.GetOk("engine")
	engineVersion, engineVersionGot := d.GetOk("engine_version")

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
			if engineGot && engine != SupportedEngine.Engine {
				continue
			}
			ids = append(ids, SupportedEngine.Engine)
			for _, SupportedEngineVersion := range SupportedEngine.SupportedEngineVersions.SupportedEngineVersion {
				if engineVersionGot && engineVersion != SupportedEngineVersion.Version {
					continue
				}
				ids = append(ids, SupportedEngineVersion.Version)

				for _, SupportedCategory := range SupportedEngineVersion.SupportedCategorys.SupportedCategory {
					zoneIds := []map[string]interface{}{}
					if _, ok := engineInfos[SupportedCategory.Category]; ok {
						zoneIds = append(engineInfos[SupportedCategory.Category].ZoneIds, zoneId)
					} else {
						zoneIds = []map[string]interface{}{zoneId}
						indexMap[SupportedCategory.Category] = len(engineInfos)
					}
					engineInfos[SupportedCategory.Category] = engineInfosItem{
						Index:         indexMap[SupportedCategory.Category],
						Engine:        SupportedEngine.Engine,
						EngineVersion: SupportedEngineVersion.Version,
						ZoneIds:       zoneIds,
					}
				}
			}
		}
	}

	infos := make([]map[string]interface{}, len(engineInfos))
	for k, v := range engineInfos {
		infos[v.Index] = map[string]interface{}{
			"zone_ids":       v.ZoneIds,
			"engine":         v.Engine,
			"engine_version": v.EngineVersion,
			"category":       k,
		}
	}

	d.SetId(dataResourceIdHash(ids))
	err = d.Set("instance_engines", infos)
	if err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		err = writeToFile(output.(string), infos)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

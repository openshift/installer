package alicloud

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"

	"github.com/denverdino/aliyungo/common"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudZonesRead,

		Schema: map[string]*schema.Schema{
			"available_instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
			},
			"available_resource_creation": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ResourceTypeInstance),
					string(ResourceTypeRds),
					string(ResourceTypePolarDB),
					string(ResourceTypeRkv),
					string(ResourceTypeVSwitch),
					string(ResourceTypeDisk),
					string(IoOptimized),
					string(ResourceTypeFC),
					string(ResourceTypeElasticsearch),
					string(ResourceTypeSlb),
					string(ResourceTypeMongoDB),
					string(ResourceTypeGpdb),
					string(ResourceTypeHBase),
					string(ResourceTypeAdb),
				}, false),
			},
			"available_slb_address_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(Vpc),
					string(ClassicIntranet),
					string(ClassicInternet),
				}, false),
			},
			"available_slb_address_ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(IPV4),
					string(IPV6),
				}, false),
			},
			"available_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"all",
					"cloud",
					"ephemeral_ssd",
					"cloud_essd",
					"cloud_efficiency",
					"cloud_ssd",
					"local_disk",
				}, false),
			},

			"multi": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  PostPaid,
				// %q must contain a valid InstanceChargeType, expected common.PrePaid, common.PostPaid
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Vpc", "Classic"}, false),
			},
			"spot_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      NoSpot,
				ValidateFunc: validation.StringInSlice([]string{"NoSpot", "SpotAsPriceGo", "SpotWithPriceLimit"}, false),
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			// Computed values.
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"available_resource_creation": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"available_disk_categories": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"multi_zone_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"slb_slave_zone_ids": {
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

func dataSourceAlicloudZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	resType := d.Get("available_resource_creation").(string)
	multi := d.Get("multi").(bool)
	var zoneIds []string
	rdsZones := make(map[string]string)
	polarDBZones := make(map[string]string)
	rkvZones := make(map[string]string)
	mongoDBZones := make(map[string]string)
	gpdbZones := make(map[string]string)
	hbaseZones := make(map[string]string)
	adbZones := make(map[string]string)
	instanceChargeType := d.Get("instance_charge_type").(string)

	if strings.ToLower(Trim(resType)) == strings.ToLower(string(ResourceTypePolarDB)) {
		request := polardb.CreateDescribeRegionsRequest()
		request.RegionId = client.RegionId
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeRegions(request)
		})
		if err != nil {
			return WrapError(fmt.Errorf("[ERROR] DescribeRegions got an error: %#v", err))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		regions, _ := raw.(*polardb.DescribeRegionsResponse)
		if len(regions.Regions.Region) <= 0 {
			return WrapError(fmt.Errorf("[ERROR] There is no available region for PolarDB."))
		}
		for _, r := range regions.Regions.Region {
			for _, zone := range r.Zones.Zone {
				if multi && strings.Contains(zone.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
					zoneIds = append(zoneIds, zone.ZoneId)
					continue
				}
				polarDBZones[zone.ZoneId] = r.RegionId
			}
		}
	}
	if strings.ToLower(Trim(resType)) == strings.ToLower(string(ResourceTypeRkv)) {
		request := r_kvstore.CreateDescribeAvailableResourceRequest()
		request.RegionId = client.RegionId
		request.InstanceChargeType = instanceChargeType
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DescribeAvailableResource(request)
		})
		if err != nil {
			return WrapError(fmt.Errorf("[ERROR] DescribeAvailableResource got an error: %#v", err))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		zones, _ := raw.(*r_kvstore.DescribeAvailableResourceResponse)
		if len(zones.AvailableZones.AvailableZone) <= 0 {
			return WrapError(fmt.Errorf("[ERROR] There is no available zones for KVStore"))
		}
		for _, zone := range zones.AvailableZones.AvailableZone {
			if multi && strings.Contains(zone.ZoneId, MULTI_IZ_SYMBOL) {
				zoneIds = append(zoneIds, zone.ZoneId)
				continue
			}
			rkvZones[zone.ZoneId] = zone.RegionId
		}
	}
	if strings.ToLower(Trim(resType)) == strings.ToLower(string(ResourceTypeMongoDB)) {
		request := dds.CreateDescribeRegionsRequest()
		request.RegionId = client.RegionId
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DescribeRegions(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		regions, _ := raw.(*dds.DescribeRegionsResponse)
		if len(regions.Regions.DdsRegion) <= 0 {
			return WrapError(fmt.Errorf("[ERROR] There is no available region for MongoDB."))
		}
		for _, r := range regions.Regions.DdsRegion {
			for _, zonid := range r.Zones.Zone {
				if multi && strings.Contains(zonid.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
					zoneIds = append(zoneIds, zonid.ZoneId)
					continue
				}
				mongoDBZones[zonid.ZoneId] = r.RegionId
			}
		}
	}
	if strings.ToLower(Trim(resType)) == strings.ToLower(string(ResourceTypeHBase)) {
		request := hbase.CreateDescribeRegionsRequest()
		request.RegionId = client.RegionId
		raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.DescribeRegions(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		regions, _ := raw.(*hbase.DescribeRegionsResponse)
		if len(regions.Regions.Region) <= 0 {
			return WrapError(fmt.Errorf("[ERROR] There is no available region for HBase."))
		}
		for _, r := range regions.Regions.Region {
			for _, zonid := range r.Zones.Zone {
				if r.RegionId == string(client.Region) {
					hbaseZones[zonid.Id] = r.RegionId
				}
			}
		}
	}
	if strings.ToLower(Trim(resType)) == strings.ToLower(string(ResourceTypeAdb)) {
		request := adb.CreateDescribeRegionsRequest()
		request.RegionId = client.RegionId
		raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.DescribeRegions(request)
		})
		if err != nil {
			return WrapError(fmt.Errorf("[ERROR] DescribeRegions got an error: %#v", err))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		regions, _ := raw.(*adb.DescribeRegionsResponse)
		if len(regions.Regions.Region) <= 0 {
			return WrapError(fmt.Errorf("[ERROR] There is no available region for adb."))
		}
		for _, r := range regions.Regions.Region {
			for _, zone := range r.Zones.Zone {
				if multi && strings.Contains(zone.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
					zoneIds = append(zoneIds, zone.ZoneId)
					continue
				}
				adbZones[zone.ZoneId] = r.RegionId
			}
		}
	}
	if strings.ToLower(Trim(resType)) == strings.ToLower(string(ResourceTypeGpdb)) {
		request := gpdb.CreateDescribeRegionsRequest()
		request.RegionId = client.RegionId
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.DescribeRegions(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*gpdb.DescribeRegionsResponse)
		if len(response.Regions.Region) <= 0 {
			return WrapError(fmt.Errorf("[ERROR] There is no available region for gpdb."))
		}
		for _, r := range response.Regions.Region {
			for _, zoneId := range r.Zones.Zone {
				if multi && strings.Contains(zoneId.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
					zoneIds = append(zoneIds, zoneId.ZoneId)
					continue
				}
				gpdbZones[zoneId.ZoneId] = r.RegionId
			}
		}
	}
	elasticsearchZones := make(map[string]string)
	if strings.ToLower(Trim(resType)) == strings.ToLower(string(ResourceTypeElasticsearch)) {
		request := elasticsearch.CreateGetRegionConfigurationRequest()
		request.RegionId = client.RegionId
		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.GetRegionConfiguration(request)
		})

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.GetActionName(), request)
		zones, _ := raw.(*elasticsearch.GetRegionConfigurationResponse)
		for _, zoneID := range zones.Result.Zones {
			if multi && strings.Contains(zoneID, MULTI_IZ_SYMBOL) {
				zoneIds = append(zoneIds, zoneID)
				continue
			}

			elasticsearchZones[zoneID] = string(client.Region)
		}
	}

	if len(zoneIds) > 0 {
		sort.Strings(zoneIds)
		return zoneIdsDescriptionAttributes(d, zoneIds)
	}

	// Retrieving available zones for VPC-FC
	if strings.ToLower(Trim(resType)) == strings.ToLower(string(ResourceTypeFC)) {
		var clientInfo *fc.Client
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			clientInfo = fcClient
			return fcClient.GetAccountSettings(fc.NewGetAccountSettingsInput())
		})
		if err != nil {
			return fmt.Errorf("[API ERROR] FC GetAccountSettings: %#v", err)
		}
		addDebug("GetAccountSettings", raw, clientInfo)
		out, _ := raw.(*fc.GetAccountSettingsOutput)
		if out != nil && len(out.AvailableAZs) > 0 {
			sort.Strings(out.AvailableAZs)
			return zoneIdsDescriptionAttributes(d, out.AvailableAZs)
		}
	}

	// Retrieving available zones for SLB
	slaveZones := make(map[string][]string)
	if strings.ToLower(Trim(resType)) == strings.ToLower(string(ResourceTypeSlb)) {
		request := slb.CreateDescribeAvailableResourceRequest()
		request.RegionId = client.RegionId
		if ipVersion, ok := d.GetOk("available_slb_address_ip_version"); ok {
			request.AddressIPVersion = ipVersion.(string)
		}
		if addressType, ok := d.GetOk("available_slb_address_type"); ok {
			request.AddressType = addressType.(string)
		}
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeAvailableResource(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*slb.DescribeAvailableResourceResponse)
		for _, resource := range response.AvailableResources.AvailableResource {
			slaveIds := slaveZones[resource.MasterZoneId]
			slaveIds = append(slaveIds, resource.SlaveZoneId)
			if len(slaveIds) > 0 {
				sort.Strings(slaveIds)
			}
			slaveZones[resource.MasterZoneId] = slaveIds
		}
	}

	_, validZones, _, err := ecsService.DescribeAvailableResources(d, meta, ZoneResource)
	if err != nil {
		return err
	}

	req := ecs.CreateDescribeZonesRequest()
	req.RegionId = client.RegionId
	req.InstanceChargeType = instanceChargeType
	if v, ok := d.GetOk("spot_strategy"); ok && v.(string) != "" {
		req.SpotStrategy = v.(string)
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeZones(req)
	})
	if err != nil {
		return fmt.Errorf("DescribeZones got an error: %#v", err)
	}
	addDebug(req.GetActionName(), raw, req.RpcRequest, req)
	resp, _ := raw.(*ecs.DescribeZonesResponse)
	if resp == nil || len(resp.Zones.Zone) < 1 {
		return fmt.Errorf("There are no availability zones in the region: %#v.", client.Region)
	}

	mapZones := make(map[string]ecs.Zone)
	insType, _ := d.Get("available_instance_type").(string)
	diskType, _ := d.Get("available_disk_category").(string)

	for _, zone := range resp.Zones.Zone {
		for _, v := range validZones {
			if zone.ZoneId != v.ZoneId {
				continue
			}
			if len(zone.AvailableInstanceTypes.InstanceTypes) <= 0 ||
				(insType != "" && !constraints(zone.AvailableInstanceTypes.InstanceTypes, insType)) {
				continue
			}
			if len(zone.AvailableDiskCategories.DiskCategories) <= 0 ||
				(diskType != "" && !constraints(zone.AvailableDiskCategories.DiskCategories, diskType)) {
				continue
			}
			if len(rdsZones) > 0 {
				if _, ok := rdsZones[zone.ZoneId]; !ok {
					continue
				}
			}
			if len(polarDBZones) > 0 {
				if _, ok := polarDBZones[zone.ZoneId]; !ok {
					continue
				}
			}
			if len(rkvZones) > 0 {
				if _, ok := rkvZones[zone.ZoneId]; !ok {
					continue
				}
			}
			if len(mongoDBZones) > 0 {
				if _, ok := mongoDBZones[zone.ZoneId]; !ok {
					continue
				}
			}
			if len(hbaseZones) > 0 {
				if _, ok := hbaseZones[zone.ZoneId]; !ok {
					continue
				}
			}
			if len(gpdbZones) > 0 {
				if _, ok := gpdbZones[zone.ZoneId]; !ok {
					continue
				}
			}
			if len(elasticsearchZones) > 0 {
				if _, ok := elasticsearchZones[zone.ZoneId]; !ok {
					continue
				}
			}
			if len(slaveZones) > 0 {
				if _, ok := slaveZones[zone.ZoneId]; !ok {
					continue
				}
			}
			if len(adbZones) > 0 {
				if _, ok := adbZones[zone.ZoneId]; !ok {
					continue
				}
			}

			zoneIds = append(zoneIds, zone.ZoneId)
			mapZones[zone.ZoneId] = zone
		}
	}

	if len(zoneIds) > 0 {
		// Sort zones before reading
		sort.Strings(zoneIds)
	}

	var s []map[string]interface{}
	for _, zoneId := range zoneIds {
		mapping := map[string]interface{}{"id": zoneId}
		if len(slaveZones) > 0 {
			mapping["slb_slave_zone_ids"] = slaveZones[zoneId]
		}
		if !d.Get("enable_details").(bool) {
			s = append(s, mapping)
			continue
		}
		mapping["local_name"] = mapZones[zoneId].LocalName
		mapping["available_instance_types"] = mapZones[zoneId].AvailableInstanceTypes.InstanceTypes
		mapping["available_resource_creation"] = mapZones[zoneId].AvailableResourceCreation.ResourceTypes
		mapping["available_disk_categories"] = mapZones[zoneId].AvailableDiskCategories.DiskCategories
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(zoneIds))
	if err := d.Set("zones", s); err != nil {
		return err
	}

	if err := d.Set("ids", zoneIds); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

// check array constraints str
func constraints(arr interface{}, v string) bool {
	arrs := reflect.ValueOf(arr)
	len := arrs.Len()
	for i := 0; i < len; i++ {
		if arrs.Index(i).String() == v {
			return true
		}
	}
	return false
}

func zoneIdsDescriptionAttributes(d *schema.ResourceData, zones []string) error {
	var s []map[string]interface{}
	var zoneIds []string
	for _, t := range zones {
		mapping := map[string]interface{}{
			"id":             t,
			"multi_zone_ids": splitMultiZoneId(t),
		}
		s = append(s, mapping)
		zoneIds = append(zoneIds, t)
	}

	d.SetId(dataResourceIdHash(zones))
	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", zoneIds); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

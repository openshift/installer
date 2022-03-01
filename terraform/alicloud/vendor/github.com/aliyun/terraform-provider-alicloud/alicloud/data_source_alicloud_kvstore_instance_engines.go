package alicloud

import (
	"strings"
	"time"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudKVStoreInstanceEngines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKVStoreInstanceEnginesRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(KVStoreMemcache),
					string(KVStoreRedis),
				}, false),
				Default: string(KVStoreRedis),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PrePaid,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
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
						"zone_id": {
							Type:     schema.TypeString,
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
					},
				},
			},
		},
	}
}

func dataSourceAlicloudKVStoreInstanceEnginesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := r_kvstore.CreateDescribeAvailableResourceRequest()
	request.RegionId = client.RegionId
	request.ZoneId = d.Get("zone_id").(string)
	instanceChargeType := d.Get("instance_charge_type").(string)
	request.InstanceChargeType = instanceChargeType
	request.Engine = d.Get("engine").(string)
	var response = &r_kvstore.DescribeAvailableResourceResponse{}
	err := resource.Retry(time.Minute*5, func() *resource.RetryError {
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DescribeAvailableResource(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{Throttling}) {
				time.Sleep(time.Duration(5) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*r_kvstore.DescribeAvailableResourceResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kvstore_instance_engines", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	var infos []map[string]interface{}
	var ids []string

	engine, engineGot := d.GetOk("engine")
	engine = strings.ToLower(engine.(string))
	engineVersion, engineVersionGot := d.GetOk("engine_version")

	for _, AvailableZone := range response.AvailableZones.AvailableZone {
		zondId := AvailableZone.ZoneId
		ids = append(ids, zondId)
		versions := make(map[string]interface{})
		for _, SupportedEngine := range AvailableZone.SupportedEngines.SupportedEngine {
			if engineGot && engine != SupportedEngine.Engine {
				continue
			}
			ids = append(ids, SupportedEngine.Engine)
			if strings.ToLower(engine.(string)) == "memcache" {
				info := make(map[string]interface{})
				info["zone_id"] = AvailableZone.ZoneId
				info["engine"] = SupportedEngine.Engine
				info["engine_version"] = "2.8"
				ids = append(ids, "2.8")
				infos = append(infos, info)
			} else {
				for _, editionType := range SupportedEngine.SupportedEditionTypes.SupportedEditionType {
					for _, seriesType := range editionType.SupportedSeriesTypes.SupportedSeriesType {
						for _, SupportedEngineVersion := range seriesType.SupportedEngineVersions.SupportedEngineVersion {
							if engineVersionGot && engineVersion.(string) != SupportedEngineVersion.Version {
								continue
							}
							versions[SupportedEngineVersion.Version] = nil
						}
					}
				}
				for version := range versions {
					info := make(map[string]interface{})
					info["zone_id"] = AvailableZone.ZoneId
					info["engine"] = SupportedEngine.Engine
					info["engine_version"] = version
					ids = append(ids, version)
					infos = append(infos, info)
				}
			}
		}
	}

	d.SetId(dataResourceIdHash(ids))
	err = d.Set("instance_engines", infos)
	if err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), infos)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

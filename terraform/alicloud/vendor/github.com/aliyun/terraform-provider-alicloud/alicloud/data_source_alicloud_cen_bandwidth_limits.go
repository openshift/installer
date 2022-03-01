package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCenBandwidthLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenBandwidthLimitsRead,

		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values

			"limits": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opposite_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenBandwidthLimitsRead(d *schema.ResourceData, meta interface{}) error {
	var allCenBwLimits []cbn.CenInterRegionBandwidthLimit

	instanceIds := make([]string, 0)
	if v, ok := d.GetOk("instance_ids"); ok {
		for _, vv := range v.([]interface{}) {
			instanceIds = append(instanceIds, Trim(vv.(string)))
		}
	} else {
		instanceIds = append(instanceIds, "")
	}

	for _, instanceId := range instanceIds {
		tmpAllCenBwLimits, err := getCenBandwidthLimits(instanceId, meta)
		if err != nil {
			return WrapError(err)
		}
		allCenBwLimits = append(allCenBwLimits, tmpAllCenBwLimits...)
	}

	return cenInterRegionBandwidthLimitsAttributes(d, allCenBwLimits)
}

func getCenBandwidthLimits(instanceId string, meta interface{}) ([]cbn.CenInterRegionBandwidthLimit, error) {
	client := meta.(*connectivity.AliyunClient)

	request := cbn.CreateDescribeCenInterRegionBandwidthLimitsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	if instanceId != "" {
		request.CenId = instanceId
	}

	var allCenBwLimits []cbn.CenInterRegionBandwidthLimit

	deadline := time.Now().Add(10 * time.Minute)
	for {
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenInterRegionBandwidthLimits(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.User"}) {
				if time.Now().After(deadline) {
					return nil, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_bandwidth_limits", request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				time.Sleep(10 * time.Second)
				continue
			}
			return allCenBwLimits, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_bandwidth_limits", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cbn.DescribeCenInterRegionBandwidthLimitsResponse)

		if len(response.CenInterRegionBandwidthLimits.CenInterRegionBandwidthLimit) < 1 {
			break
		}
		allCenBwLimits = append(allCenBwLimits, response.CenInterRegionBandwidthLimits.CenInterRegionBandwidthLimit...)

		if len(response.CenInterRegionBandwidthLimits.CenInterRegionBandwidthLimit) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return allCenBwLimits, WrapError(err)
		}
		request.PageNumber = page
	}

	return allCenBwLimits, nil
}

func cenInterRegionBandwidthLimitsAttributes(d *schema.ResourceData, allCenBwLimits []cbn.CenInterRegionBandwidthLimit) error {
	var ids []string
	var s []map[string]interface{}

	for _, cenBwLimit := range allCenBwLimits {
		mapping := map[string]interface{}{
			"instance_id":        cenBwLimit.CenId,
			"local_region_id":    cenBwLimit.LocalRegionId,
			"opposite_region_id": cenBwLimit.OppositeRegionId,
			"status":             cenBwLimit.Status,
			"bandwidth_limit":    cenBwLimit.BandwidthLimit,
		}

		id := cenBwLimit.CenId + COLON_SEPARATED + cenBwLimit.LocalRegionId + COLON_SEPARATED + cenBwLimit.OppositeRegionId
		ids = append(ids, id)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("limits", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

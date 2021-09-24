package alicloud

import (
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCenBandwidthPackages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenBandwidthPackagesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"include_reservation_data": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Idle", "InUse"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"packages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_bandwidth_package_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_bandwidth_package_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"geographic_region_a_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"geographic_region_b_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"geographic_span_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"has_reservation_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_cross_border": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_package_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reservation_active_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reservation_bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reservation_internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reservation_order_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenBandwidthPackagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cbn.CreateDescribeCenBandwidthPackagesRequest()
	filters := make([]cbn.DescribeCenBandwidthPackagesFilter, 0)
	if v, ok := d.GetOk("instance_id"); ok {
		filter := cbn.DescribeCenBandwidthPackagesFilter{
			Key:   "CenId",
			Value: &[]string{v.(string)},
		}
		filters = append(filters, filter)
		request.Filter = &filters
	}

	if v, ok := d.GetOkExists("include_reservation_data"); ok {
		request.IncludeReservationData = requests.NewBoolean(v.(bool))
	}

	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []cbn.CenBandwidthPackage
	var cenBandwidthPackageNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		cenBandwidthPackageNameRegex = r
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
	status, statusOk := d.GetOk("status")
	var response *cbn.DescribeCenBandwidthPackagesResponse
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.DescribeCenBandwidthPackages(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"AliyunGoClientFailure", "ServiceUnavailable", "Throttling", "Throttling.User"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			response, _ = raw.(*cbn.DescribeCenBandwidthPackagesResponse)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_bandwidth_packages", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		for _, item := range response.CenBandwidthPackages.CenBandwidthPackage {
			if cenBandwidthPackageNameRegex != nil {
				if !cenBandwidthPackageNameRegex.MatchString(item.Name) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.CenBandwidthPackageId]; !ok {
					continue
				}
			}
			if statusOk && status != "" && status != item.Status {
				continue
			}
			objects = append(objects, item)
		}
		if len(response.CenBandwidthPackages.CenBandwidthPackage) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"bandwidth":                        object.Bandwidth,
			"business_status":                  object.BusinessStatus,
			"id":                               object.CenBandwidthPackageId,
			"cen_bandwidth_package_id":         object.CenBandwidthPackageId,
			"cen_bandwidth_package_name":       object.Name,
			"name":                             object.Name,
			"cen_ids":                          object.CenIds.CenId,
			"description":                      object.Description,
			"expired_time":                     object.ExpiredTime,
			"geographic_region_a_id":           convertGeographicRegionAIdResponse(object.GeographicRegionAId),
			"geographic_region_b_id":           convertGeographicRegionBIdResponse(object.GeographicRegionBId),
			"geographic_span_id":               object.GeographicSpanId,
			"has_reservation_data":             object.HasReservationData,
			"is_cross_border":                  object.IsCrossBorder,
			"payment_type":                     object.BandwidthPackageChargeType,
			"bandwidth_package_charge_type":    object.BandwidthPackageChargeType,
			"reservation_active_time":          object.ReservationActiveTime,
			"reservation_bandwidth":            object.ReservationBandwidth,
			"reservation_internet_charge_type": object.ReservationInternetChargeType,
			"reservation_order_type":           object.ReservationOrderType,
			"status":                           object.Status,
		}
		if len(object.CenIds.CenId) > 0 {
			mapping["instance_id"] = object.CenIds.CenId[0]
		}
		ids = append(ids, object.CenBandwidthPackageId)
		names = append(names, object.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("packages", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

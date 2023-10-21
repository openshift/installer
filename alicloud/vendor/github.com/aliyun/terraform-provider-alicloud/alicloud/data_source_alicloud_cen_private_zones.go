package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCenPrivateZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenPrivateZonesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Creating", "Deleting"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_zone_dns_servers": {
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

func dataSourceAlicloudCenPrivateZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cbn.CreateDescribeCenPrivateZoneRoutesRequest()
	request.CenId = d.Get("cen_id").(string)
	if v, ok := d.GetOk("host_region_id"); ok {
		request.HostRegionId = v.(string)
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []cbn.PrivateZoneInfo

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
	var response *cbn.DescribeCenPrivateZoneRoutesResponse
	for {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenPrivateZoneRoutes(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_private_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ = raw.(*cbn.DescribeCenPrivateZoneRoutesResponse)

		for _, item := range response.PrivateZoneInfos.PrivateZoneInfo {
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.AccessRegionId]; !ok {
					continue
				}
			}
			if statusOk && status != "" && status != item.Status {
				continue
			}
			objects = append(objects, item)
		}
		if len(response.PrivateZoneInfos.PrivateZoneInfo) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))
	for i, object := range objects {
		mapping := map[string]interface{}{
			"access_region_id":         object.AccessRegionId,
			"cen_id":                   response.CenId,
			"host_region_id":           object.HostRegionId,
			"host_vpc_id":              object.HostVpcId,
			"private_zone_dns_servers": response.PrivateZoneDnsServers,
			"status":                   object.Status,
		}
		ids[i] = object.AccessRegionId
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCenRegionRouteEntries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenRegionDomainRouteEntriesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"entries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenRegionDomainRouteEntriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cbn.CreateDescribeCenRegionDomainRouteEntriesRequest()
	request.RegionId = client.RegionId
	request.CenId = d.Get("instance_id").(string)
	request.CenRegionId = d.Get("region_id").(string)

	request.PageSize = requests.NewInteger(PageSizeLarge)

	var allCenRouteEntries []cbn.CenRouteEntry
	for pageNumber := 1; ; pageNumber++ {
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenRegionDomainRouteEntries(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_region_route_entries", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cbn.DescribeCenRegionDomainRouteEntriesResponse)

		if len(response.CenRouteEntries.CenRouteEntry) < 1 {
			break
		}
		allCenRouteEntries = append(allCenRouteEntries, response.CenRouteEntries.CenRouteEntry...)

		if len(response.CenRouteEntries.CenRouteEntry) < PageSizeLarge {
			break
		}

		request.PageNumber = requests.NewInteger(pageNumber)
	}

	return cenRegionDomainRouteEntriesAttributes(d, allCenRouteEntries)
}

func cenRegionDomainRouteEntriesAttributes(d *schema.ResourceData, allCenRouteEntries []cbn.CenRouteEntry) error {
	var s []map[string]interface{}

	for _, cenRouteEntry := range allCenRouteEntries {
		mapping := map[string]interface{}{
			"cidr_block":         cenRouteEntry.DestinationCidrBlock,
			"type":               cenRouteEntry.Type,
			"next_hop_id":        cenRouteEntry.NextHopInstanceId,
			"next_hop_type":      cenRouteEntry.NextHopType,
			"next_hop_region_id": cenRouteEntry.NextHopRegionId,
		}

		s = append(s, mapping)
	}
	id := d.Get("instance_id").(string) + COLON_SEPARATED + d.Get("region_id").(string)
	d.SetId(id)
	if err := d.Set("entries", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

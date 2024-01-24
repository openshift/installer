package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCenRouteEntries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenPublishedRouteEntriesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
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
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
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
						"route_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operational_mode": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"publish_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// Complex computed value
						"conflicts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_block": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_type": {
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
				},
			},
		},
	}
}

func dataSourceAlicloudCenPublishedRouteEntriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}

	request := cbn.CreateDescribePublishedRouteEntriesRequest()
	request.RegionId = client.RegionId
	request.CenId = d.Get("instance_id").(string)
	request.ChildInstanceRouteTableId = d.Get("route_table_id").(string)
	if v, ok := d.GetOk("cidr_block"); ok {
		request.DestinationCidrBlock = v.(string)
	}

	childInstanceId, childInstanceType, err := cenService.CreateCenRouteEntryParas(request.ChildInstanceRouteTableId)
	if err != nil {
		return WrapError(err)
	}
	request.ChildInstanceId = childInstanceId
	request.ChildInstanceType = childInstanceType
	request.ChildInstanceRegionId = client.RegionId

	request.PageSize = requests.NewInteger(PageSizeLarge)

	var allPublishedRouteEntries []cbn.PublishedRouteEntry
	for pageNumber := 1; ; pageNumber++ {
		request.PageNumber = requests.NewInteger(pageNumber)
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribePublishedRouteEntries(request)
		})

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_route_entries", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cbn.DescribePublishedRouteEntriesResponse)

		if len(response.PublishedRouteEntries.PublishedRouteEntry) < 1 {
			break
		}
		allPublishedRouteEntries = append(allPublishedRouteEntries, response.PublishedRouteEntries.PublishedRouteEntry...)

		if len(response.PublishedRouteEntries.PublishedRouteEntry) < PageSizeLarge {
			break
		}
	}

	return cenPublishedRouteEntriesAttributes(d, allPublishedRouteEntries)
}

func cenPublishedRouteEntriesAttributes(d *schema.ResourceData, allPublishedRouteEntries []cbn.PublishedRouteEntry) error {
	var ids []string
	var s []map[string]interface{}

	for _, routeEntry := range allPublishedRouteEntries {
		mapping := map[string]interface{}{
			"route_table_id":   routeEntry.ChildInstanceRouteTableId,
			"cidr_block":       routeEntry.DestinationCidrBlock,
			"next_hop_type":    routeEntry.NextHopType,
			"next_hop_id":      routeEntry.NextHopId,
			"operational_mode": routeEntry.OperationalMode,
			"publish_status":   routeEntry.PublishStatus,
			"route_type":       routeEntry.RouteType,
			// Complex types get their own functions
			"conflicts": routeConflictsMappings(routeEntry.Conflicts.Conflict),
		}

		id := routeEntry.ChildInstanceRouteTableId + COLON_SEPARATED + routeEntry.DestinationCidrBlock
		ids = append(ids, id)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("entries", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

// Returns a set of route conflicts mappings.
func routeConflictsMappings(m []cbn.Conflict) []map[string]interface{} {
	var s []map[string]interface{}

	for _, v := range m {
		mapping := map[string]interface{}{
			"cidr_block":    v.DestinationCidrBlock,
			"region_id":     v.RegionId,
			"instance_id":   v.InstanceId,
			"instance_type": v.InstanceType,
			"status":        v.Status,
		}

		s = append(s, mapping)
	}

	return s
}

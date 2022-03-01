package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCenRouteMaps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenRouteMapsRead,
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cen_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Creating", "Deleting"}, false),
			},
			"transmit_direction": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"RegionIn", "RegionOut"}, false),
			},
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"maps": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"as_path_match_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_match_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"community_match_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"community_operate_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_child_instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"destination_cidr_blocks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"destination_instance_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"destination_instance_ids_reverse_match": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"destination_route_table_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"map_result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"match_asns": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_community_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"next_priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"operate_community_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"preference": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"prepend_as_path": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_map_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"source_child_instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"source_instance_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"source_instance_ids_reverse_match": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"source_region_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"source_route_table_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transmit_direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenRouteMapsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cbn.CreateDescribeCenRouteMapsRequest()
	request.CenId = d.Get("cen_id").(string)
	if v, ok := d.GetOk("cen_region_id"); ok {
		request.CenRegionId = v.(string)
	}
	if v, ok := d.GetOk("transmit_direction"); ok {
		request.TransmitDirection = v.(string)
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []cbn.RouteMap
	var descriptionRegex *regexp.Regexp
	if v, ok := d.GetOk("description_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		descriptionRegex = r
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
	for {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenRouteMaps(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_route_maps", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*cbn.DescribeCenRouteMapsResponse)

		for _, item := range response.RouteMaps.RouteMap {
			if descriptionRegex != nil {
				if !descriptionRegex.MatchString(item.Description) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.RouteMapId]; !ok {
					continue
				}
			}
			if statusOk && status != "" && status != item.Status {
				continue
			}
			objects = append(objects, item)
		}
		if len(response.RouteMaps.RouteMap) < PageSizeLarge {
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
			"as_path_match_mode":                     object.AsPathMatchMode,
			"cen_id":                                 object.CenId,
			"cen_region_id":                          object.CenRegionId,
			"cidr_match_mode":                        object.CidrMatchMode,
			"community_match_mode":                   object.CommunityMatchMode,
			"community_operate_mode":                 object.CommunityOperateMode,
			"description":                            object.Description,
			"destination_child_instance_types":       object.DestinationChildInstanceTypes.DestinationChildInstanceType,
			"destination_cidr_blocks":                object.DestinationCidrBlocks.DestinationCidrBlock,
			"destination_instance_ids":               object.DestinationInstanceIds.DestinationInstanceId,
			"destination_instance_ids_reverse_match": object.DestinationInstanceIdsReverseMatch,
			"destination_route_table_ids":            object.DestinationRouteTableIds.DestinationRouteTableId,
			"map_result":                             object.MapResult,
			"match_asns":                             object.MatchAsns.MatchAsn,
			"match_community_set":                    object.MatchCommunitySet.MatchCommunity,
			"next_priority":                          object.NextPriority,
			"operate_community_set":                  object.OperateCommunitySet.OperateCommunity,
			"preference":                             object.Preference,
			"prepend_as_path":                        object.PrependAsPath.AsPath,
			"priority":                               object.Priority,
			"id":                                     object.RouteMapId,
			"route_map_id":                           object.RouteMapId,
			"route_types":                            object.RouteTypes.RouteType,
			"source_child_instance_types":            object.SourceChildInstanceTypes.SourceChildInstanceType,
			"source_instance_ids":                    object.SourceInstanceIds.SourceInstanceId,
			"source_instance_ids_reverse_match":      object.SourceInstanceIdsReverseMatch,
			"source_region_ids":                      object.SourceRegionIds.SourceRegionId,
			"source_route_table_ids":                 object.SourceRouteTableIds.SourceRouteTableId,
			"status":                                 object.Status,
			"transmit_direction":                     object.TransmitDirection,
		}
		ids[i] = object.RouteMapId
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("maps", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

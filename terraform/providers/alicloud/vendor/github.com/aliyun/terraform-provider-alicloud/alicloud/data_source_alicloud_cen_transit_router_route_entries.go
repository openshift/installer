package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCenTransitRouterRouteEntries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenTransitRouterRouteEntriesRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Creating", "Deleting", "Updating"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"transit_router_route_entry_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"transit_router_route_entry_names": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"transit_router_route_entry_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Creating", "Deleting", "Updating"}, false),
			},
			"transit_router_route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_entry_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_entry_destination_cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_entry_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_entry_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_entry_next_hop_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_entry_next_hop_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_entry_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenTransitRouterRouteEntriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTransitRouterRouteEntries"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("transit_router_route_entry_ids"); ok {
		request["TransitRouterRouteEntryIds"] = convertListToJsonString(v.([]interface{}))
	}
	if v, ok := d.GetOk("transit_router_route_entry_names"); ok {
		request["TransitRouterRouteEntryNames"] = convertListToJsonString(v.([]interface{}))
	}
	if v, ok := d.GetOk("transit_router_route_entry_status"); ok {
		request["TransitRouterRouteEntryStatus"] = v
	}
	request["TransitRouterRouteTableId"] = d.Get("transit_router_route_table_id")
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var transitRouterRouteEntryNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		transitRouterRouteEntryNameRegex = r
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
	var response map[string]interface{}
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_transit_router_route_entries", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.TransitRouterRouteEntries", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouterRouteEntries", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if transitRouterRouteEntryNameRegex != nil {
				if !transitRouterRouteEntryNameRegex.MatchString(fmt.Sprint(item["TransitRouterRouteEntryName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TransitRouterRouteEntryId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["TransitRouterRouteEntryStatus"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"status":                                 object["TransitRouterRouteEntryStatus"],
			"transit_router_route_entry_description": object["TransitRouterRouteEntryDescription"],
			"transit_router_route_entry_destination_cidr_block": object["TransitRouterRouteEntryDestinationCidrBlock"],
			"id":                                       fmt.Sprint(object["TransitRouterRouteEntryId"]),
			"transit_router_route_entry_id":            fmt.Sprint(object["TransitRouterRouteEntryId"]),
			"transit_router_route_entry_name":          object["TransitRouterRouteEntryName"],
			"transit_router_route_entry_next_hop_id":   object["TransitRouterRouteEntryNextHopId"],
			"transit_router_route_entry_next_hop_type": object["TransitRouterRouteEntryNextHopType"],
			"transit_router_route_entry_type":          object["TransitRouterRouteEntryType"],
		}
		ids = append(ids, fmt.Sprint(object["TransitRouterRouteEntryId"]))
		names = append(names, object["TransitRouterRouteEntryName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("entries", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

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

func dataSourceAlicloudCenTransitRouterRouteTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenTransitRouterRouteTablesRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Creating", "Deleting", "Updating"}, false),
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"transit_router_route_table_ids": {
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
			"transit_router_route_table_names": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"transit_router_route_table_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Creating", "Deleting", "Updating"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_table_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_table_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenTransitRouterRouteTablesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTransitRouterRouteTables"
	request := make(map[string]interface{})
	request["TransitRouterId"] = d.Get("transit_router_id")
	if v, ok := d.GetOk("transit_router_route_table_ids"); ok {
		request["TransitRouterRouteTableIds"] = convertListToJsonString(v.([]interface{}))
	}
	if v, ok := d.GetOk("transit_router_route_table_names"); ok {
		request["TransitRouterRouteTableNames"] = convertListToJsonString(v.([]interface{}))
	}
	if v, ok := d.GetOk("transit_router_route_table_status"); ok {
		request["TransitRouterRouteTableStatus"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var transitRouterRouteTableNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		transitRouterRouteTableNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_transit_router_route_tables", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.TransitRouterRouteTables", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouterRouteTables", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if transitRouterRouteTableNameRegex != nil {
				if !transitRouterRouteTableNameRegex.MatchString(fmt.Sprint(item["TransitRouterRouteTableName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TransitRouterRouteTableId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["TransitRouterRouteTableStatus"].(string) {
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
			"status":                                 object["TransitRouterRouteTableStatus"],
			"transit_router_route_table_description": object["TransitRouterRouteTableDescription"],
			"id":                                     fmt.Sprint(object["TransitRouterRouteTableId"]),
			"transit_router_route_table_id":          fmt.Sprint(object["TransitRouterRouteTableId"]),
			"transit_router_route_table_name":        object["TransitRouterRouteTableName"],
			"transit_router_route_table_type":        object["TransitRouterRouteTableType"],
		}
		ids = append(ids, fmt.Sprint(object["TransitRouterRouteTableId"]))
		names = append(names, object["TransitRouterRouteTableName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("tables", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

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

func dataSourceAlicloudCenTransitRouterPeerAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenTransitRouterPeerAttachmentsRead,
		Schema: map[string]*schema.Schema{
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
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Attached", "Attaching", "Detaching"}, false),
			},
			"transit_router_attachment_id": {
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
			"transit_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_publish_route_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cen_bandwidth_package_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"geographic_span_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_transit_router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_transit_router_owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_transit_router_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_attachment_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_attachment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_attachment_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenTransitRouterPeerAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTransitRouterPeerAttachments"
	request := make(map[string]interface{})
	request["CenId"] = d.Get("cen_id")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}
	if v, ok := d.GetOk("transit_router_attachment_id"); ok {
		request["TransitRouterAttachmentId"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var transitRouterPeerAttachmentNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		transitRouterPeerAttachmentNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_transit_router_peer_attachments", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.TransitRouterAttachments", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouterAttachments", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if transitRouterPeerAttachmentNameRegex != nil {
				if !transitRouterPeerAttachmentNameRegex.MatchString(fmt.Sprint(item["TransitRouterAttachmentName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TransitRouterAttachmentId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
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
			"auto_publish_route_enabled":            object["AutoPublishRouteEnabled"],
			"bandwidth":                             formatInt(object["Bandwidth"]),
			"cen_bandwidth_package_id":              object["CenBandwidthPackageId"],
			"geographic_span_id":                    object["GeographicSpanId"],
			"peer_transit_router_id":                object["PeerTransitRouterId"],
			"peer_transit_router_owner_id":          fmt.Sprint(object["PeerTransitRouterOwnerId"]),
			"peer_transit_router_region_id":         object["PeerTransitRouterRegionId"],
			"resource_type":                         object["ResourceType"],
			"status":                                object["Status"],
			"transit_router_attachment_description": object["TransitRouterAttachmentDescription"],
			"id":                                    fmt.Sprint(object["TransitRouterAttachmentId"]),
			"transit_router_attachment_id":          fmt.Sprint(object["TransitRouterAttachmentId"]),
			"transit_router_attachment_name":        object["TransitRouterAttachmentName"],
			"transit_router_id":                     object["TransitRouterId"],
		}
		ids = append(ids, fmt.Sprint(object["TransitRouterAttachmentId"]))
		names = append(names, object["TransitRouterAttachmentName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("attachments", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

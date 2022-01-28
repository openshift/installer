package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpcBgpPeers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpcBgpPeersRead,
		Schema: map[string]*schema.Schema{
			"bgp_group_id": {
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
			"router_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Available", "Deleted", "Deleting", "Modifying", "Pending"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bfd_multi_hop": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bgp_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bgp_peer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bgp_peer_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bgp_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_bfd": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"hold": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_fake": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"keepalive": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_asn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_asn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_limit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"router_id": {
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

func dataSourceAlicloudVpcBgpPeersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeBgpPeers"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("bgp_group_id"); ok {
		request["BgpGroupId"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("router_id"); ok {
		request["RouterId"] = v
	}
	request["PageSize"] = PageSizeMedium
	request["PageNumber"] = 1
	var objects []map[string]interface{}
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
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_bgp_peers", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.BgpPeers.BgpPeer", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.BgpPeers.BgpPeer", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["BgpPeerId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"auth_key":        object["AuthKey"],
			"bfd_multi_hop":   formatInt(object["BfdMultiHop"]),
			"bgp_group_id":    object["BgpGroupId"],
			"id":              fmt.Sprint(object["BgpPeerId"]),
			"bgp_peer_id":     fmt.Sprint(object["BgpPeerId"]),
			"bgp_peer_name":   object["Name"],
			"bgp_status":      object["BgpStatus"],
			"description":     object["Description"],
			"enable_bfd":      object["EnableBfd"],
			"hold":            object["Hold"],
			"ip_version":      object["IpVersion"],
			"is_fake":         object["IsFake"],
			"keepalive":       object["Keepalive"],
			"local_asn":       object["LocalAsn"],
			"peer_asn":        object["PeerAsn"],
			"peer_ip_address": object["PeerIpAddress"],
			"route_limit":     object["RouteLimit"],
			"router_id":       object["RouterId"],
			"status":          object["Status"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("peers", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

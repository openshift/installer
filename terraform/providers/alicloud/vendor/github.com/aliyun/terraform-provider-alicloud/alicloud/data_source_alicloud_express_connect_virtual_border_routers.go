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

func dataSourceAlicloudExpressConnectVirtualBorderRouters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudExpressConnectVirtualBorderRoutersRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"active", "deleting", "recovering", "terminated", "terminating", "unconfirmed"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"routers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_point_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"activation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"circuit_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_box_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"detect_multiplier": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ecc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_ipv6": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"local_gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_ipv6_gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_rx_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_tx_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"payment_vbr_expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_ipv6_gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peering_ipv6_subnet_mask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peering_subnet_mask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"physical_connection_business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"physical_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"physical_connection_owner_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"physical_connection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"recovery_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"termination_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_border_router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_border_router_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vlan_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vlan_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudExpressConnectVirtualBorderRoutersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeVirtualBorderRouters"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("filter"); ok {
		for filterPtr, filter := range v.([]interface{}) {
			filterArg := filter.(map[string]interface{})
			request[fmt.Sprintf("Filter.%d.Key", filterPtr+1)] = filterArg["key"]
			for valuePtr, value := range filterArg["values"].([]interface{}) {
				request[fmt.Sprintf("Filter.%d.Value.%d", filterPtr+1, valuePtr+1)] = value
			}
		}
	}
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var virtualBorderRouterNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		virtualBorderRouterNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_express_connect_virtual_border_routers", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.VirtualBorderRouterSet.VirtualBorderRouterType", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.VirtualBorderRouterSet.VirtualBorderRouterType", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if virtualBorderRouterNameRegex != nil && !virtualBorderRouterNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["VbrId"])]; !ok {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"access_point_id":                     object["AccessPointId"],
			"activation_time":                     object["ActivationTime"],
			"circuit_code":                        object["CircuitCode"],
			"cloud_box_instance_id":               object["CloudBoxInstanceId"],
			"create_time":                         object["CreationTime"],
			"description":                         object["Description"],
			"detect_multiplier":                   formatInt(object["DetectMultiplier"]),
			"ecc_id":                              object["EccId"],
			"enable_ipv6":                         object["EnableIpv6"],
			"local_gateway_ip":                    object["LocalGatewayIp"],
			"local_ipv6_gateway_ip":               object["LocalIpv6GatewayIp"],
			"min_rx_interval":                     formatInt(object["MinRxInterval"]),
			"min_tx_interval":                     formatInt(object["MinTxInterval"]),
			"payment_vbr_expire_time":             object["PConnVbrExpireTime"],
			"peer_gateway_ip":                     object["PeerGatewayIp"],
			"peer_ipv6_gateway_ip":                object["PeerIpv6GatewayIp"],
			"peering_ipv6_subnet_mask":            object["PeeringIpv6SubnetMask"],
			"peering_subnet_mask":                 object["PeeringSubnetMask"],
			"physical_connection_business_status": object["PhysicalConnectionBusinessStatus"],
			"physical_connection_id":              object["PhysicalConnectionId"],
			"physical_connection_owner_uid":       object["PhysicalConnectionOwnerUid"],
			"physical_connection_status":          object["PhysicalConnectionStatus"],
			"recovery_time":                       object["RecoveryTime"],
			"route_table_id":                      object["RouteTableId"],
			"status":                              object["Status"],
			"termination_time":                    object["TerminationTime"],
			"type":                                object["Type"],
			"id":                                  fmt.Sprint(object["VbrId"]),
			"virtual_border_router_id":            fmt.Sprint(object["VbrId"]),
			"virtual_border_router_name":          object["Name"],
			"vlan_id":                             formatInt(object["VlanId"]),
			"vlan_interface_id":                   object["VlanInterfaceId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("routers", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

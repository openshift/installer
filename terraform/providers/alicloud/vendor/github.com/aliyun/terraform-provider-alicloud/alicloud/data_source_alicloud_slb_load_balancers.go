package alicloud

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudSlbLoadBalancers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbLoadBalancersRead,
		Schema: map[string]*schema.Schema{
			"master_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Removed:  "Field 'master_availability_zone' has removed from v1.123.1 and replace by 'master_zone_id'.",
			},
			"slave_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Removed:  "Field 'slave_availability_zone' has removed from v1.123.1 and replace by 'slave_zone_id'.",
			},
			"address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"address_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
			},
			"address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"internet", "intranet"}, false),
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, true),
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
			"load_balancer_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"master_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"vpc", "classic"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"server_intranet_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"slave_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"active", "inactive", "locked"}, false),
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
			"balancers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address_ip_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_release_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backend_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"server_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time_stamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"delete_protection": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time_stamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_ports_and_protocal": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"listener_protocal": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"listener_ports_and_protocol": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"forward_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"listener_forward": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"listener_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"listener_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modification_protection_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modification_protection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"renewal_cyc_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"renewal_duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"renewal_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slave_zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"slbs": {
				Type:       schema.TypeList,
				Computed:   true,
				Deprecated: "Field 'slbs' has deprecated from v1.123.1 and replace by 'balancers'.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slave_availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudSlbLoadBalancersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeLoadBalancers"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("address"); ok {
		request["Address"] = v
	}
	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIPVersion"] = v
	}
	if v, ok := d.GetOk("address_type"); ok {
		request["AddressType"] = v
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = convertSlbLoadBalancerInternetChargeTypeRequest(v)
	}
	if v, ok := d.GetOk("load_balancer_name"); ok {
		request["LoadBalancerName"] = v
	}
	if v, ok := d.GetOk("master_zone_id"); ok {
		request["MasterZoneId"] = v
	}
	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertSlbLoadBalancerPaymentTypeRequest(v)
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("server_id"); ok {
		request["ServerId"] = v
	}
	if v, ok := d.GetOk("server_intranet_address"); ok {
		request["ServerIntranetAddress"] = v
	}
	if v, ok := d.GetOk("slave_zone_id"); ok {
		request["SlaveZoneId"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["LoadBalancerStatus"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag"] = tags
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
	}
	var objects []interface{}
	var loadBalancerNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		loadBalancerNameRegex = r
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
	var response map[string]interface{}
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_load_balancers", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.LoadBalancers.LoadBalancer", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.LoadBalancers.LoadBalancer", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if loadBalancerNameRegex != nil {
				if !loadBalancerNameRegex.MatchString(fmt.Sprint(item["LoadBalancerName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["LoadBalancerId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	slbs := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"address":                        object["Address"],
			"address_ip_version":             object["AddressIPVersion"],
			"address_type":                   object["AddressType"],
			"create_time_stamp":              formatInt(object["CreateTimeStamp"]),
			"id":                             fmt.Sprint(object["LoadBalancerId"]),
			"load_balancer_id":               fmt.Sprint(object["LoadBalancerId"]),
			"load_balancer_name":             object["LoadBalancerName"],
			"master_zone_id":                 object["MasterZoneId"],
			"modification_protection_reason": object["ModificationProtectionReason"],
			"modification_protection_status": object["ModificationProtectionStatus"],
			"network_type":                   object["NetworkType"],
			"payment_type":                   convertSlbLoadBalancerPaymentTypeResponse(object["PayType"]),
			"region_id_alias":                object["RegionIdAlias"],
			"resource_group_id":              object["ResourceGroupId"],
			"slave_zone_id":                  object["SlaveZoneId"],
			"status":                         object["LoadBalancerStatus"],
			"vswitch_id":                     object["VSwitchId"],
			"vpc_id":                         object["VpcId"],
		}

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.Tag", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["TagKey"].(string)
				value := t.(map[string]interface{})["TagValue"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags

		slb := map[string]interface{}{
			"id":                       fmt.Sprint(object["LoadBalancerId"]),
			"region_id":                object["RegionIdAlias"],
			"master_availability_zone": object["MasterZoneId"],
			"slave_availability_zone":  object["SlaveZoneId"],
			"status":                   object["LoadBalancerStatus"],
			"name":                     object["LoadBalancerName"],
			"network_type":             object["NetworkType"],
			"vpc_id":                   object["VpcId"],
			"vswitch_id":               object["VSwitchId"],
			"address":                  object["Address"],
			"internet":                 fmt.Sprint(object["NetworkType"]) == strings.ToLower(string(Internet)),
			"creation_time":            object["CreationTime"],
			"tags":                     tags,
		}
		slbs = append(slbs, slb)

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["LoadBalancerId"]))
			names = append(names, object["LoadBalancerName"])
			s = append(s, mapping)
			continue
		}

		slbService := SlbService{client}
		id := fmt.Sprint(object["LoadBalancerId"])
		getResp, err := slbService.DescribeSlbLoadBalancer(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["auto_release_time"] = getResp["AutoReleaseTime"]

		backendServer := make([]map[string]interface{}, 0)
		if backendServerList, ok := getResp["BackendServers"].(map[string]interface{})["BackendServer"].([]interface{}); ok {
			for _, v := range backendServerList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"description": m1["Description"],
						"server_id":   m1["ServerId"],
						"type":        m1["Type"],
						"weight":      formatInt(m1["Weight"]),
					}
					backendServer = append(backendServer, temp1)
				}
			}
		}
		mapping["backend_servers"] = backendServer
		mapping["bandwidth"] = getResp["Bandwidth"]
		mapping["delete_protection"] = getResp["DeleteProtection"]
		mapping["end_time"] = getResp["EndTime"]
		mapping["end_time_stamp"] = getResp["EndTimeStamp"]
		mapping["internet_charge_type"] = convertSlbLoadBalancerInternetChargeTypeResponse(getResp["InternetChargeType"])

		listenerPortAndProtocal := make([]map[string]interface{}, 0)
		if listenerPortAndProtocalList, ok := getResp["ListenerPortsAndProtocal"].(map[string]interface{})["ListenerPortAndProtocal"].([]interface{}); ok {
			for _, v := range listenerPortAndProtocalList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"listener_port":     formatInt(m1["ListenerPort"]),
						"listener_protocal": m1["ListenerProtocal"],
					}
					listenerPortAndProtocal = append(listenerPortAndProtocal, temp1)
				}
			}
		}
		mapping["listener_ports_and_protocal"] = listenerPortAndProtocal

		listenerPortAndProtocol := make([]map[string]interface{}, 0)
		if listenerPortAndProtocolList, ok := getResp["ListenerPortsAndProtocol"].(map[string]interface{})["ListenerPortAndProtocol"].([]interface{}); ok {
			for _, v := range listenerPortAndProtocolList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"description":       m1["Description"],
						"forward_port":      formatInt(m1["ForwardPort"]),
						"listener_forward":  m1["ListenerForward"],
						"listener_port":     formatInt(m1["ListenerPort"]),
						"listener_protocol": m1["ListenerProtocol"],
					}
					listenerPortAndProtocol = append(listenerPortAndProtocol, temp1)
				}
			}
		}
		mapping["listener_ports_and_protocol"] = listenerPortAndProtocol
		mapping["load_balancer_spec"] = getResp["LoadBalancerSpec"]
		mapping["renewal_cyc_unit"] = getResp["RenewalCycUnit"]
		mapping["renewal_duration"] = getResp["RenewalDuration"]
		mapping["renewal_status"] = getResp["RenewalStatus"]
		ids = append(ids, fmt.Sprint(object["LoadBalancerId"]))
		names = append(names, object["LoadBalancerName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("balancers", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("slbs", slbs); err != nil {
		return WrapError(err)
	}
	if err := d.Set("total_count", formatInt(response["TotalCount"])); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

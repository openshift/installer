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

func dataSourceAlicloudExpressConnectPhysicalConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudExpressConnectPhysicalConnectionsRead,
		Schema: map[string]*schema.Schema{
			"include_reservation_data": {
				Type:     schema.TypeBool,
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Allocated", "Allocating", "Allocation Failed", "Approved", "Canceled", "Confirmed", "Enabled", "Initial", "Rejected", "Terminated"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_point_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ad_location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"circuit_code": {
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
						"enabled_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"has_reservation_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"loa_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"physical_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"physical_connection_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"redundant_physical_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reservation_active_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reservation_internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reservation_order_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudExpressConnectPhysicalConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribePhysicalConnections"
	request := make(map[string]interface{})
	if v, ok := d.GetOkExists("include_reservation_data"); ok {
		request["IncludeReservationData"] = v
	}
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var physicalConnectionNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		physicalConnectionNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_express_connect_physical_connections", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.PhysicalConnectionSet.PhysicalConnectionType", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PhysicalConnectionSet.PhysicalConnectionType", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if physicalConnectionNameRegex != nil && !physicalConnectionNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PhysicalConnectionId"])]; !ok {
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
			"access_point_id":                  object["AccessPointId"],
			"ad_location":                      object["AdLocation"],
			"bandwidth":                        fmt.Sprint(object["Bandwidth"]),
			"business_status":                  object["BusinessStatus"],
			"circuit_code":                     object["CircuitCode"],
			"create_time":                      object["CreationTime"],
			"description":                      object["Description"],
			"enabled_time":                     object["EnabledTime"],
			"end_time":                         object["EndTime"],
			"has_reservation_data":             fmt.Sprint(object["HasReservationData"]),
			"line_operator":                    object["LineOperator"],
			"loa_status":                       object["LoaStatus"],
			"payment_type":                     object["ChargeType"],
			"peer_location":                    object["PeerLocation"],
			"id":                               fmt.Sprint(object["PhysicalConnectionId"]),
			"physical_connection_id":           fmt.Sprint(object["PhysicalConnectionId"]),
			"physical_connection_name":         object["Name"],
			"port_number":                      object["PortNumber"],
			"port_type":                        object["PortType"],
			"redundant_physical_connection_id": object["RedundantPhysicalConnectionId"],
			"reservation_active_time":          object["ReservationActiveTime"],
			"reservation_internet_charge_type": object["ReservationInternetChargeType"],
			"reservation_order_type":           object["ReservationOrderType"],
			"spec":                             object["Spec"],
			"status":                           object["Status"],
			"type":                             object["Type"],
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

	if err := d.Set("connections", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

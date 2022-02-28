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

func dataSourceAlicloudDatabaseGatewayGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDatabaseGatewayGatewaysRead,
		Schema: map[string]*schema.Schema{
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
			"search_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"EXCEPTION", "NEW", "RUNNING", "STOPPED"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connect_endpoint_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"current_daemon_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"current_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_point": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gateway_instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gateway_instance_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_update_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"local_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"output_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"gateway_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hosts": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
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
		},
	}
}

func dataSourceAlicloudDatabaseGatewayGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "GetUserGateways"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("search_key"); ok {
		request["SearchKey"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var gatewayNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		gatewayNameRegex = r
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
	conn, err := client.NewDgClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-27"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_database_gateway_gateways", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
		}
		result, _ := convertJsonStringToList(resp.(string))
		for _, v := range result {
			item := v.(map[string]interface{})
			if gatewayNameRegex != nil && !gatewayNameRegex.MatchString(fmt.Sprint(item["gatewayName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["gatewayId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["status"].(string) {
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
			"create_time":   object["gmtCreateString"],
			"gateway_desc":  object["gatewayDesc"],
			"gateway_name":  object["gatewayName"],
			"hosts":         object["hosts"],
			"id":            object["gatewayId"],
			"modified_time": object["gmtModifiedString"],
			"parent_id":     object["parentId"],
			"status":        object["status"],
			"user_id":       object["userId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["gatewayName"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["gatewayId"])
		dgService := DgService{client}
		getResp, err := dgService.GetUserGatewayInstances(id)
		if err != nil {
			return WrapError(err)
		}

		gatewayInstances := make([]map[string]interface{}, 0)
		for _, v := range getResp {
			if m, ok := v.(map[string]interface{}); ok {
				temp := map[string]interface{}{
					"connect_endpoint_type":   m["connectEndpointType"],
					"current_daemon_version":  m["currentDaemonVersion"],
					"current_version":         m["currentVersion"],
					"end_point":               m["endPoint"],
					"gateway_instance_id":     m["gatewayInstanceId"],
					"gateway_instance_status": m["gatewayInstanceStatus"],
					"last_update_time":        m["lastUpdateTime"],
					"local_ip":                m["localIP"],
					"message":                 m["message"],
					"output_ip":               m["outputIP"],
				}
				gatewayInstances = append(gatewayInstances, temp)
			}
		}
		mapping["gateway_instances"] = gatewayInstances
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("gateways", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

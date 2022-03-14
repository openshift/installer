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

func dataSourceAlicloudMseGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMseGatewaysRead,
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
			"gateway_name": {
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
				ValidateFunc: validation.StringInSlice([]string{"0", "1", "2", "3", "4", "6", "8", "9", "10", "11", "12", "13"}, false),
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
			"gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_unique_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replica": {
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
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slb_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"associate_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"slb_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"slb_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"slb_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gmt_create": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gateway_slb_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gateway_slb_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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

func dataSourceAlicloudMseGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListGateway"
	request := make(map[string]interface{})
	filterParams := make(map[string]interface{})

	if v, ok := d.GetOk("gateway_name"); ok {
		filterParams["Name"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		filterParams["Vpc"] = v
	}

	if len(filterParams) > 0 {
		if v, err := convertMaptoJsonString(filterParams); err != nil {
			return WrapError(err)
		} else {
			request["FilterParams"] = v
		}
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
	conn, err := client.NewMseClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-31"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_mse_gateways", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.Result", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Result", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if gatewayNameRegex != nil && !gatewayNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["GatewayUniqueId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != fmt.Sprint(item["Status"]) {
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
			"gateway_name":      object["Name"],
			"id":                fmt.Sprint(object["GatewayUniqueId"]),
			"gateway_unique_id": fmt.Sprint(object["GatewayUniqueId"]),
			"payment_type":      object["ChargeType"],
			"replica":           fmt.Sprint(object["Replica"]),
			"spec":              object["Spec"],
			"status":            fmt.Sprint(object["Status"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["GatewayUniqueId"])
		mseService := MseService{client}
		getResp, err := mseService.DescribeMseGateway(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["backup_vswitch_id"] = getResp["Vswitch2"]
		mapping["vswitch_id"] = getResp["Vswitch"]
		mapping["vpc_id"] = getResp["Vpc"]

		slbListObject, err := mseService.ListGatewaySlb(id)
		if err != nil {
			return WrapError(err)
		}
		slbList := slbListObject["SlbList"]

		slbMapList := make([]map[string]interface{}, 0)
		for _, v := range slbList.([]interface{}) {
			slbMap := v.(map[string]interface{})
			slb := map[string]interface{}{
				"associate_id":       slbMap["Id"],
				"slb_id":             slbMap["SlbId"],
				"slb_ip":             slbMap["SlbIp"],
				"slb_port":           slbMap["SlbPort"],
				"type":               slbMap["Type"],
				"gmt_create":         slbMap["GmtCreate"],
				"gateway_slb_mode":   slbMap["GatewaySlbMode"],
				"gateway_slb_status": slbMap["GatewaySlbStatus"],
			}
			slbMapList = append(slbMapList, slb)
		}
		mapping["slb_list"] = slbMapList
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

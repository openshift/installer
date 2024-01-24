package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudPrivatelinkVpcEndpoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPrivatelinkVpcEndpointsRead,
		Schema: map[string]*schema.Schema{
			"connection_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Pending", "Connecting", "Connected", "Disconnecting", "Disconnected", "Deleting", "ServiceDeleted"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Creating", "Deleted", "Deleting", "Inactive", "Pending"}, false),
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
			"vpc_endpoint_name": {
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
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_endpoint_name": {
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudPrivatelinkVpcEndpointsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListVpcEndpoints"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("connection_status"); ok {
		request["ConnectionStatus"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("service_name"); ok {
		request["ServiceName"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["EndpointStatus"] = v
	}
	if v, ok := d.GetOk("vpc_endpoint_name"); ok {
		request["EndpointName"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var vpcEndpointNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		vpcEndpointNameRegex = r
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
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_privatelink_vpc_endpoints", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Endpoints", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Endpoints", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if vpcEndpointNameRegex != nil {
				if !vpcEndpointNameRegex.MatchString(fmt.Sprint(item["EndpointName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["EndpointId"])]; !ok {
					continue
				}
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
			"bandwidth":                object["Bandwidth"],
			"connection_status":        object["ConnectionStatus"],
			"endpoint_business_status": object["EndpointBusinessStatus"],
			"endpoint_description":     object["EndpointDescription"],
			"endpoint_domain":          object["EndpointDomain"],
			"id":                       fmt.Sprint(object["EndpointId"]),
			"endpoint_id":              fmt.Sprint(object["EndpointId"]),
			"service_id":               object["ServiceId"],
			"service_name":             object["ServiceName"],
			"status":                   object["EndpointStatus"],
			"vpc_endpoint_name":        object["EndpointName"],
			"vpc_id":                   object["VpcId"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["EndpointId"]))
			names = append(names, object["EndpointName"])
			s = append(s, mapping)
			continue
		}

		privatelinkService := PrivatelinkService{client}
		id := fmt.Sprint(object["EndpointId"])
		getResp, err := privatelinkService.ListVpcEndpointSecurityGroups(id)
		if err != nil {
			return WrapError(err)
		}
		if v := getResp["SecurityGroups"].([]interface{}); len(v) > 0 {
			mapping["security_group_ids"] = convertSecurityGroupIdToStringList(v)
		}
		ids = append(ids, fmt.Sprint(object["EndpointId"]))
		names = append(names, object["EndpointName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("endpoints", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

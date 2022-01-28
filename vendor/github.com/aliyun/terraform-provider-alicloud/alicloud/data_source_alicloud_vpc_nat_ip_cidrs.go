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

func dataSourceAlicloudVpcNatIpCidrs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpcNatIpCidrsRead,
		Schema: map[string]*schema.Schema{
			"nat_gateway_id": {
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
			"nat_ip_cidr_name": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nat_ip_cidrs": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Available"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"nat_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_ip_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_ip_cidr_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_ip_cidr_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_ip_cidr_name": {
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

func dataSourceAlicloudVpcNatIpCidrsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListNatIpCidrs"
	request := make(map[string]interface{})
	request["NatGatewayId"] = d.Get("nat_gateway_id")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("nat_ip_cidrs"); ok {
		request["NatIpCidrs"] = v
	}
	if v, ok := d.GetOk("nat_ip_cidr_name"); ok {
		request["NatIpCidrName"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["NatIpCidrStatus"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var natIpCidrNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		natIpCidrNameRegex = r
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
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("ListNatIpCidrs")
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_nat_ip_cidrs", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.NatIpCidrs", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.NatIpCidrs", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if natIpCidrNameRegex != nil && !natIpCidrNameRegex.MatchString(fmt.Sprint(item["NatIpCidrName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%s:%s", item["NatGatewayId"], item["NatIpCidr"])]; !ok {
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
			"create_time":             object["CreationTime"],
			"is_default":              object["IsDefault"],
			"nat_gateway_id":          object["NatGatewayId"],
			"id":                      fmt.Sprintf("%s:%s", object["NatGatewayId"], object["NatIpCidr"]),
			"nat_ip_cidr":             fmt.Sprint(object["NatIpCidr"]),
			"nat_ip_cidr_description": object["NatIpCidrDescription"],
			"nat_ip_cidr_id":          object["NatIpCidrId"],
			"status":                  object["NatIpCidrStatus"],
			"nat_ip_cidr_name":        object["NatIpCidrName"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, mapping["nat_ip_cidr_name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("cidrs", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

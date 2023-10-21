package alicloud

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpnCustomerGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnCgwsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},

			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpnCgwsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateDescribeCustomerGatewaysRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var allCgws []vpc.CustomerGateway

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeCustomerGateways(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "aliclioud_vpn_customer_gateways", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeCustomerGatewaysResponse)
		if len(response.CustomerGateways.CustomerGateway) < 1 {
			break
		}
		allCgws = append(allCgws, response.CustomerGateways.CustomerGateway...)
		if len(response.CustomerGateways.CustomerGateway) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredCgws []vpc.CustomerGateway
	var reg *regexp.Regexp
	var ids []string
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, item := range v.([]interface{}) {
			if item == nil {
				continue
			}
			ids = append(ids, strings.Trim(item.(string), " "))
		}
	}

	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		if r, err := regexp.Compile(nameRegex.(string)); err == nil {
			reg = r
		} else {
			return WrapError(err)
		}
	}

	for _, cgw := range allCgws {
		if reg != nil {
			if !reg.MatchString(cgw.Name) {
				continue
			}
		}

		if ids != nil && len(ids) != 0 {
			for _, id := range ids {
				if cgw.CustomerGatewayId == id {
					filteredCgws = append(filteredCgws, cgw)
				}
			}
		} else {
			filteredCgws = append(filteredCgws, cgw)
		}
	}

	return vpnCgwsDecriptionAttributes(d, filteredCgws, meta)
}

func vpnCgwsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.CustomerGateway, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"id":          vpn.CustomerGatewayId,
			"name":        vpn.Name,
			"ip_address":  vpn.IpAddress,
			"description": vpn.Description,
			"create_time": TimestampToStr(vpn.CreateTime),
		}
		ids = append(ids, vpn.CustomerGatewayId)
		names = append(names, vpn.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("gateways", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

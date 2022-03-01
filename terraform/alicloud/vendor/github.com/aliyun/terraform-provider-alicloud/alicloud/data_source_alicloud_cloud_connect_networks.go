package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCloudConnectNetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudConnectNetworkRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"networks": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCloudConnectNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateDescribeCloudConnectNetworksRequest()

	var allCcnInstances []smartag.CloudConnectNetwork
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithSagClient(func(ccnClient *smartag.Client) (interface{}, error) {
			return ccnClient.DescribeCloudConnectNetworks(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_connect_networks", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*smartag.DescribeCloudConnectNetworksResponse)
		cloudConnectNetworks := response.CloudConnectNetworks.CloudConnectNetwork
		for _, cloudConnectNetwork := range cloudConnectNetworks {
			allCcnInstances = append(allCcnInstances, cloudConnectNetwork)
		}

		if len(cloudConnectNetworks) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredCcnInstancesTemp []smartag.CloudConnectNetwork

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	for _, ccnInstance := range allCcnInstances {
		if v, ok := d.GetOk("id"); ok && v.(string) != "" && ccnInstance.CcnId != v.(string) {
			continue
		}
		if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
			r := regexp.MustCompile(v.(string))
			if !r.MatchString(ccnInstance.Name) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[ccnInstance.CcnId]; !ok {
				continue
			}
		}
		filteredCcnInstancesTemp = append(filteredCcnInstancesTemp, ccnInstance)

	}

	return cloudConnectNetworkAttributes(d, filteredCcnInstancesTemp, meta)
}

func cloudConnectNetworkAttributes(d *schema.ResourceData, ccnInstances []smartag.CloudConnectNetwork, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, ccnInstance := range ccnInstances {
		mapping := map[string]interface{}{
			"id":          ccnInstance.CcnId,
			"name":        ccnInstance.Name,
			"description": ccnInstance.Description,
			"cidr_block":  ccnInstance.CidrBlock,
			"is_default":  ccnInstance.IsDefault,
		}
		names = append(names, ccnInstance.Name)
		ids = append(ids, ccnInstance.CcnId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("networks", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudApiGatewayApis() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudApigatewayApisRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"api_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'api_id' has been deprecated from provider version 1.52.2. New field 'ids' replaces it.",
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
			"apis": {
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
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudApigatewayApisRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cloudapi.CreateDescribeApisRequest()
	request.RegionId = client.RegionId

	if groupId, ok := d.GetOk("group_id"); ok {
		request.GroupId = groupId.(string)
	}
	if apiId, ok := d.GetOk("api_id"); ok {
		request.ApiId = apiId.(string)
	}

	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var allapis []cloudapi.ApiSummary

	for {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeApis(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_api_gateway_apis", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.DescribeApisResponse)

		allapis = append(allapis, response.ApiSummarys.ApiSummary...)

		if len(allapis) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredApisTemp []cloudapi.ApiSummary

	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	for _, api := range allapis {
		if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
			r := regexp.MustCompile(v.(string))
			if !r.MatchString(api.ApiName) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[api.ApiId]; !ok {
				continue
			}
		}
		filteredApisTemp = append(filteredApisTemp, api)
	}

	return apiGatewayApisDescribeSummarys(d, filteredApisTemp, meta)
}

func apiGatewayApisDescribeSummarys(d *schema.ResourceData, apis []cloudapi.ApiSummary, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, api := range apis {
		mapping := map[string]interface{}{
			"id":          api.ApiId,
			"name":        api.ApiName,
			"region_id":   api.RegionId,
			"group_id":    api.GroupId,
			"group_name":  api.GroupName,
			"description": api.Description,
		}
		ids = append(ids, api.ApiId)
		s = append(s, mapping)
		names = append(names, api.ApiName)

	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("apis", s); err != nil {
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

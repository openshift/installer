package alicloud

import (
	"regexp"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDdoscooInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDdoscooInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
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
						"base_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"domain_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDdoscooInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ddoscoo.CreateDescribeInstancesRequest()
	request.RegionId = client.RegionId
	request.PageSize = strconv.Itoa(PageSizeSmall)
	request.PageNumber = "1"
	var instances []ddoscoo.Instance

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		ids := expandStringList(v.([]interface{}))
		request.InstanceIds = &ids
	}
	// describe ddoscoo instance filtered by name_regex
	for {
		raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.DescribeInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ddoscoo_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ddoscoo.DescribeInstancesResponse)
		if len(response.Instances) < 1 {
			break
		}

		for _, item := range response.Instances {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.Remark) {
					continue
				}
			}
			instances = append(instances, item)
		}

		if len(response.Instances) < PageSizeLarge {
			break
		}

		currentPageNo, err := strconv.Atoi(request.PageNumber)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ddoscoo_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		if page, err := getNextpageNumber(requests.NewInteger(currentPageNo)); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = string(page)
		}
	}

	// describe instance spec filtered by instanceids
	nameMap := make(map[string]string)
	var instanceIds []string
	for _, instance := range instances {
		instanceIds = append(instanceIds, instance.InstanceId)
		nameMap[instance.InstanceId] = instance.Remark
	}

	if len(instanceIds) < 1 {
		return WrapError(extractDdoscooInstance(d, nameMap, []ddoscoo.InstanceSpec{}))
	}

	specReq := ddoscoo.CreateDescribeInstanceSpecsRequest()
	specReq.InstanceIds = &instanceIds

	raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.DescribeInstanceSpecs(specReq)
	})

	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ddoscoo_instances", specReq.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(specReq.GetActionName(), raw, specReq.RpcRequest, specReq)
	response, _ := raw.(*ddoscoo.DescribeInstanceSpecsResponse)

	return WrapError(extractDdoscooInstance(d, nameMap, response.InstanceSpecs))
}

func extractDdoscooInstance(d *schema.ResourceData, nameMap map[string]string, instanceSpecs []ddoscoo.InstanceSpec) error {
	var instanceIds []string
	var names []string
	var s []map[string]interface{}

	for _, item := range instanceSpecs {
		mapping := map[string]interface{}{
			"id":                item.InstanceId,
			"name":              nameMap[item.InstanceId],
			"bandwidth":         item.ElasticBandwidth,
			"base_bandwidth":    item.BaseBandwidth,
			"service_bandwidth": item.BandwidthMbps,
			"port_count":        item.PortLimit,
			"domain_count":      item.DomainLimit,
		}
		instanceIds = append(instanceIds, item.InstanceId)
		names = append(names, nameMap[item.InstanceId])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(instanceIds))
	if err := d.Set("ids", instanceIds); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

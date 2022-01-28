package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_bastionhost"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudBastionhostInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudBastionhostInstancesRead,

		Schema: map[string]*schema.Schema{
			"description_regex": {
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
				ForceNew: true,
			},
			"descriptions": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_domain": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"instance_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"license_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_network_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"tags": tagsSchema(),
					},
				},
			},
			"tags": tagsSchema(),
		},
	}
}

func dataSourceAlicloudBastionhostInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := yundun_bastionhost.CreateDescribeInstanceBastionhostRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.CurrentPage = requests.NewInteger(1)
	var instances []yundun_bastionhost.Instance

	// get name Regex
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("description_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	if v, ok := d.GetOk("ids"); ok {
		ids, _ := v.([]interface{})
		var ids_str []string
		for _, v_instance_id := range ids {
			if v_instance_id == nil {
				continue
			}
			ids_str = append(ids_str, v_instance_id.(string))
		}
		request.InstanceId = &ids_str
	}

	if v, ok := d.GetOk("tags"); ok {
		var tags []yundun_bastionhost.DescribeInstanceBastionhostTag
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, yundun_bastionhost.DescribeInstanceBastionhostTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &tags
	}
	for {
		raw, err := client.WithBastionhostClient(func(bastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
			return bastionhostClient.DescribeInstanceBastionhost(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_bastionhost_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*yundun_bastionhost.DescribeInstanceBastionhostResponse)
		if len(response.Instances) < 1 {
			break
		}

		for _, e := range response.Instances {
			if nameRegex != nil && !nameRegex.MatchString(e.Description) {
				continue
			}
			instances = append(instances, e)
		}

		if len(response.Instances) < PageSizeSmall {
			break
		}

		currentPageNo := request.CurrentPage
		if page, err := getNextpageNumber(currentPageNo); err != nil {
			return WrapError(err)
		} else {
			request.CurrentPage = page
		}
	}

	var instanceIds []string
	for _, instance := range instances {
		instanceIds = append(instanceIds, instance.InstanceId)
	}
	if len(instanceIds) < 1 {
		return WrapError(extractBastionhostInstance(d, nil, nil))
	}
	specs := make([]map[string]interface{}, 0)
	var tags []yundun_bastionhost.TagResources
	BastionhostService := YundunBastionhostService{client}

	for _, instanceId := range instanceIds {
		object, err := BastionhostService.DescribeBastionhostInstance(instanceId)
		if err != nil {
			return WrapError(err)
		}
		specs = append(specs, object)

		{
			request := yundun_bastionhost.CreateListTagResourcesRequest()
			request.RegionId = client.RegionId
			request.ResourceType = strings.ToUpper(string(TagResourceInstance))
			request.ResourceId = &[]string{instanceId}
			raw, err := client.WithBastionhostClient(func(client *yundun_bastionhost.Client) (interface{}, error) {
				return client.ListTagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_yundun_bastionhost_tags ", request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			res, _ := raw.(*yundun_bastionhost.ListTagResourcesResponse)
			tags = append(tags, res.TagResources)
		}
	}
	return WrapError(extractBastionhostInstance(d, specs, tags))
}

func extractBastionhostInstance(d *schema.ResourceData, specs []map[string]interface{}, tags []yundun_bastionhost.TagResources) error {
	var instanceIds []string
	var descriptions []string
	var instances []map[string]interface{}

	for i := 0; i < len(specs); i++ {
		instanceMap := map[string]interface{}{
			"id":                    specs[i]["InstanceId"],
			"description":           specs[i]["Description"],
			"user_vswitch_id":       specs[i]["VswitchId"],
			"private_domain":        specs[i]["IntranetEndpoint"],
			"public_domain":         specs[i]["InternetEndpoint"],
			"instance_status":       specs[i]["InstanceStatus"],
			"license_code":          specs[i]["LicenseCode"],
			"public_network_access": specs[i]["PublicNetworkAccess"],
			"security_group_ids":    specs[i]["AuthorizedSecurityGroups"],
			"tags":                  bastionhostTagsToMap(tags[i].TagResource),
		}
		instanceIds = append(instanceIds, fmt.Sprint(specs[i]["InstanceId"]))
		descriptions = append(descriptions, fmt.Sprint(specs[i]["Description"]))
		instances = append(instances, instanceMap)
	}

	d.SetId(dataResourceIdHash(instanceIds))
	if err := d.Set("ids", instanceIds); err != nil {
		return WrapError(err)
	}

	if err := d.Set("descriptions", descriptions); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", instances); err != nil {
		return WrapError(err)
	}
	// storage locally
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), instances)
	}
	return nil
}

func bastionhostTagsToMap(tags []yundun_bastionhost.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !bastionhostTagIgnored(t) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func bastionhostTagIgnored(t yundun_bastionhost.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

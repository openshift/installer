package alicloud

import (
	"log"
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_dbaudit"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDbauditInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDbauditInstancesRead,

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
						"tags": tagsSchema(),
					},
				},
			},
			"tags": tagsSchema(),
		},
	}
}

func dataSourceAlicloudDbauditInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := yundun_dbaudit.CreateDescribeInstancesRequest()
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.CurrentPage = requests.NewInteger(1)
	var instances []yundun_dbaudit.Instance

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
		var tags []yundun_dbaudit.DescribeInstancesTag
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, yundun_dbaudit.DescribeInstancesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &tags
	}
	for {
		raw, err := client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
			return dbauditClient.DescribeInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_yundun_dbaudit", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*yundun_dbaudit.DescribeInstancesResponse)
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

	var instanceTags []yundun_dbaudit.TagResources
	for _, inst := range instances {
		request := yundun_dbaudit.CreateListTagResourcesRequest()
		request.RegionId = client.RegionId
		request.ResourceType = strings.ToUpper(string(TagResourceInstance))
		request.ResourceId = &[]string{inst.InstanceId}
		raw, err := client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
			return dbauditClient.ListTagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_yundun_dbaudit_tags", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*yundun_dbaudit.ListTagResourcesResponse)
		instanceTags = append(instanceTags, yundun_dbaudit.TagResources{TagResource: response.TagResources})
	}
	return WrapError(extractDbauditInstance(d, instances, instanceTags))
}

func extractDbauditInstance(d *schema.ResourceData, specs []yundun_dbaudit.Instance, tags []yundun_dbaudit.TagResources) error {

	var instanceIds []string
	var descriptions []string
	var instances []map[string]interface{}
	for i := 0; i < len(specs); i++ {
		instanceMap := map[string]interface{}{
			"id":                    specs[i].InstanceId,
			"description":           specs[i].Description,
			"user_vswitch_id":       specs[i].VswitchId,
			"private_domain":        specs[i].IntranetEndpoint,
			"public_domain":         specs[i].InternetEndpoint,
			"instance_status":       specs[i].InstanceStatus,
			"license_code":          specs[i].LicenseCode,
			"public_network_access": specs[i].PublicNetworkAccess,
			"tags":                  dbauditTagsToMap(tags[i].TagResource),
		}
		instanceIds = append(instanceIds, specs[i].InstanceId)
		descriptions = append(descriptions, specs[i].Description)
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
	log.Printf("DEBUF data source finnished")
	return nil
}

func dbauditTagsToMap(tags []yundun_dbaudit.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !dbauditTagIgnored(t) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func dbauditTagIgnored(t yundun_dbaudit.TagResource) bool {
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

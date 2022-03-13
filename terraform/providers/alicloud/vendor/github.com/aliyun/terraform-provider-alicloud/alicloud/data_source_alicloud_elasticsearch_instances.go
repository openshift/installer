package alicloud

import (
	"encoding/json"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudElasticsearch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudElasticsearchRead,

		Schema: map[string]*schema.Schema{
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"5.5.3_with_X-Pack",
					"6.3.2_with_X-Pack",
					"6.7.0_with_X-Pack",
				}, false),
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
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
						"instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_node_amount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_node_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_node_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_node_disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudElasticsearchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := elasticsearch.CreateListInstanceRequest()
	request.RegionId = client.RegionId
	request.EsVersion = d.Get("version").(string)
	request.Size = requests.NewInteger(PageSizeLarge)
	request.Page = requests.NewInteger(1)

	if v, ok := d.GetOk("tags"); ok {
		var reqTags []map[string]string

		for k, v := range v.(map[string]interface{}) {
			reqTags = append(reqTags, map[string]string{
				"tagKey":   k,
				"tagValue": v.(string),
			})
		}

		reqTagsStr, err := json.Marshal(reqTags)
		if err != nil {
			return WrapError(err)
		}
		request.Tags = string(reqTagsStr)
	}

	var instances []elasticsearch.Instance

	for {
		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.ListInstance(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_elasticsearch_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		response, _ := raw.(*elasticsearch.ListInstanceResponse)
		if len(response.Result) < 1 {
			break
		}

		for _, item := range response.Result {
			instances = append(instances, item)
		}

		if len(response.Result) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.Page)
		if err != nil {
			return WrapError(err)
		}
		request.Page = page
	}

	var filteredInstances []elasticsearch.Instance

	var descriptionRegex *regexp.Regexp
	if v, ok := d.GetOk("description_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		descriptionRegex = r
	}

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

	for _, instance := range instances {
		if descriptionRegex != nil && !descriptionRegex.MatchString(instance.Description) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[instance.InstanceId]; !ok {
				continue
			}
		}
		filteredInstances = append(filteredInstances, instance)
	}

	return WrapError(extractInstance(d, filteredInstances))
}

func extractInstance(d *schema.ResourceData, instances []elasticsearch.Instance) error {
	var ids []string
	var descriptions []string
	var s []map[string]interface{}

	for _, item := range instances {
		mapping := map[string]interface{}{
			"id":                   item.InstanceId,
			"description":          item.Description,
			"instance_charge_type": getChargeType(item.PaymentType),
			"data_node_amount":     item.NodeAmount,
			"data_node_spec":       item.NodeSpec.Spec,
			"data_node_disk_size":  item.NodeSpec.Disk,
			"data_node_disk_type":  item.NodeSpec.DiskType,
			"status":               item.Status,
			"version":              item.EsVersion,
			"created_at":           item.CreatedAt,
			"updated_at":           item.UpdatedAt,
			"vswitch_id":           item.NetworkConfig.VswitchId,
			"tags":                 elasticsearchTagsToMap(item.Tags),
		}

		ids = append(ids, item.InstanceId)
		descriptions = append(descriptions, item.Description)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("descriptions", descriptions); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

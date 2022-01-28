package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAlikafkaConsumerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlikafkaConsumerGroupsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"consumer_id_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			// Computed values
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"consumer_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlikafkaConsumerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	request := alikafka.CreateGetConsumerListRequest()
	request.InstanceId = d.Get("instance_id").(string)
	request.RegionId = client.RegionId

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.GetConsumerList(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alikafka_consumer_groups", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alikafka.GetConsumerListResponse)

	var filteredConsumerGroups []alikafka.ConsumerVO
	nameRegex, ok := d.GetOk("consumer_id_regex")
	for _, consumer := range response.ConsumerList.ConsumerVO {
		var r *regexp.Regexp

		if ok && nameRegex.(string) != "" {
			if nameRegex != "" {
				r, err = regexp.Compile(nameRegex.(string))
				if err != nil {
					return WrapError(err)
				}
			}
		}

		if r != nil && !r.MatchString(consumer.ConsumerId) {
			continue
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(consumer.InstanceId, ":", consumer.ConsumerId)]; !ok {
				continue
			}
		}

		filteredConsumerGroups = append(filteredConsumerGroups, consumer)
	}

	return alikafkaConsumerGroupsDecriptionAttributes(d, filteredConsumerGroups, meta)
}

func alikafkaConsumerGroupsDecriptionAttributes(d *schema.ResourceData, consumerGroupsInfo []alikafka.ConsumerVO, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, item := range consumerGroupsInfo {
		mapping := map[string]interface{}{
			"id":          fmt.Sprint(item.InstanceId, ":", item.ConsumerId),
			"instance_id": item.InstanceId,
			"consumer_id": item.ConsumerId,
			"remark":      item.Remark,
			"tags":        alikafkaService.tagVOTagsToMap(item.Tags.TagVO),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, item.ConsumerId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), ids)
	}
	return nil

}

package alicloud

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAlikafkaTopics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlikafkaTopicsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				ForceNew: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"topics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_topic": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"compact_topic": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"partition_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlikafkaTopicsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	request := alikafka.CreateGetTopicListRequest()
	request.InstanceId = d.Get("instance_id").(string)
	request.RegionId = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.GetTopicList(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alikafka_topics", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*alikafka.GetTopicListResponse)

	var filteredTopics []alikafka.TopicVO
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r, err = regexp.Compile(nameRegex.(string))
			if err != nil {
				return WrapError(err)
			}
		}
		for _, topic := range response.TopicList.TopicVO {
			if r != nil && !r.MatchString(topic.Topic) {
				continue
			}

			filteredTopics = append(filteredTopics, topic)
		}
	} else {
		filteredTopics = response.TopicList.TopicVO
	}
	return alikafkaTopicsDecriptionAttributes(d, filteredTopics, meta)
}

func alikafkaTopicsDecriptionAttributes(d *schema.ResourceData, topicsInfo []alikafka.TopicVO, meta interface{}) error {
	var names []string
	var s []map[string]interface{}

	for _, item := range topicsInfo {
		mapping := map[string]interface{}{
			"topic":         item.Topic,
			"create_time":   time.Unix(int64(item.CreateTime)/1000, 0).Format("2006-01-02 03:04:05"),
			"local_topic":   item.LocalTopic,
			"compact_topic": item.CompactTopic,
			"partition_num": item.PartitionNum,
			"remark":        item.Remark,
			"status":        item.Status,
		}

		names = append(names, item.Topic)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(names))

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("topics", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil

}

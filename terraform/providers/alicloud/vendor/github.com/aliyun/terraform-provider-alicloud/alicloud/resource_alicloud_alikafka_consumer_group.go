package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAlikafkaConsumerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlikafkaConsumerGroupCreate,
		Update: resourceAlicloudAlikafkaConsumerGroupUpdate,
		Read:   resourceAlicloudAlikafkaConsumerGroupRead,
		Delete: resourceAlicloudAlikafkaConsumerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"consumer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudAlikafkaConsumerGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	instanceId := d.Get("instance_id").(string)
	regionId := client.RegionId
	consumerId := d.Get("consumer_id").(string)

	request := alikafka.CreateCreateConsumerGroupRequest()
	request.InstanceId = instanceId
	request.RegionId = regionId
	request.ConsumerId = consumerId

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.CreateConsumerGroup(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_consumer_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(instanceId + ":" + consumerId)
	return resourceAlicloudAlikafkaConsumerGroupUpdate(d, meta)
}

func resourceAlicloudAlikafkaConsumerGroupRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	object, err := alikafkaService.DescribeAlikafkaConsumerGroup(d.Id())
	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("consumer_id", object.ConsumerId)

	tags, err := alikafkaService.DescribeTags(d.Id(), nil, TagResourceConsumerGroup)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", alikafkaService.tagsToMap(tags))

	return nil
}

func resourceAlicloudAlikafkaConsumerGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	if err := alikafkaService.setInstanceTags(d, TagResourceConsumerGroup); err != nil {
		return WrapError(err)
	}
	return resourceAlicloudAlikafkaConsumerGroupRead(d, meta)
}

func resourceAlicloudAlikafkaConsumerGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	consumerId := parts[1]

	request := alikafka.CreateDeleteConsumerGroupRequest()
	request.ConsumerId = consumerId
	request.InstanceId = instanceId
	request.RegionId = client.RegionId

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DeleteConsumerGroup(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(alikafkaService.WaitForAlikafkaConsumerGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}

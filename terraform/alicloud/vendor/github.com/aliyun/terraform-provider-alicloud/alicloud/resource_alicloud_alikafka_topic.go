package alicloud

import (
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudAlikafkaTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlikafkaTopicCreate,
		Update: resourceAlicloudAlikafkaTopicUpdate,
		Read:   resourceAlicloudAlikafkaTopicRead,
		Delete: resourceAlicloudAlikafkaTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"local_topic": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"compact_topic": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"partition_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      12,
				ValidateFunc: validation.IntBetween(0, 360),
			},
			"remark": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudAlikafkaTopicCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	instanceId := d.Get("instance_id").(string)
	regionId := client.RegionId
	topic := d.Get("topic").(string)

	request := alikafka.CreateCreateTopicRequest()
	request.InstanceId = instanceId
	request.RegionId = regionId
	request.Topic = topic
	if v, ok := d.GetOk("local_topic"); ok {
		request.LocalTopic = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("compact_topic"); ok {
		request.CompactTopic = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("partition_num"); ok {
		request.PartitionNum = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = v.(string)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.CreateTopic(request)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_topic", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(instanceId + ":" + topic)

	// wait topic status change from Creating to running
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, alikafkaService.KafkaTopicStatusRefreshFunc(d.Id()))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudAlikafkaTopicUpdate(d, meta)
}

func resourceAlicloudAlikafkaTopicUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	d.Partial(true)
	if err := alikafkaService.setInstanceTags(d, TagResourceTopic); err != nil {
		return WrapError(err)
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudAlikafkaTopicRead(d, meta)
	}

	instanceId := d.Get("instance_id").(string)
	if d.HasChange("remark") {
		remark := d.Get("remark").(string)
		topic := d.Get("topic").(string)
		modifyRemarkRequest := alikafka.CreateModifyTopicRemarkRequest()
		modifyRemarkRequest.InstanceId = instanceId
		modifyRemarkRequest.RegionId = client.RegionId
		modifyRemarkRequest.Topic = topic
		modifyRemarkRequest.Remark = remark

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.ModifyTopicRemark(modifyRemarkRequest)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(modifyRemarkRequest.GetActionName(), raw, modifyRemarkRequest.RpcRequest, modifyRemarkRequest)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyRemarkRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("remark")
	}

	if d.HasChange("partition_num") {
		o, n := d.GetChange("partition_num")
		oldPartitionNum := o.(int)
		newPartitionNum := n.(int)

		if newPartitionNum < oldPartitionNum {
			return WrapError(errors.New("partition_num only support adjust to a greater value."))
		} else {
			topic := d.Get("topic").(string)
			modifyPartitionReq := alikafka.CreateModifyPartitionNumRequest()
			modifyPartitionReq.InstanceId = instanceId
			modifyPartitionReq.RegionId = client.RegionId
			modifyPartitionReq.Topic = topic
			modifyPartitionReq.AddPartitionNum = requests.NewInteger(newPartitionNum - oldPartitionNum)

			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
					return alikafkaClient.ModifyPartitionNum(modifyPartitionReq)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
						time.Sleep(10 * time.Second)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(modifyPartitionReq.GetActionName(), raw, modifyPartitionReq.RpcRequest, modifyPartitionReq)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyPartitionReq.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			d.SetPartial("partition_num")
		}
	}

	d.Partial(false)
	return resourceAlicloudAlikafkaTopicRead(d, meta)
}

func resourceAlicloudAlikafkaTopicRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	object, err := alikafkaService.DescribeAlikafkaTopic(d.Id())
	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("topic", object.Topic)
	d.Set("local_topic", object.LocalTopic)
	d.Set("compact_topic", object.CompactTopic)
	d.Set("partition_num", object.PartitionNum)
	d.Set("remark", object.Remark)

	tags, err := alikafkaService.DescribeTags(d.Id(), nil, TagResourceTopic)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", alikafkaService.tagsToMap(tags))

	return nil
}

func resourceAlicloudAlikafkaTopicDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := alikafka.CreateDeleteTopicRequest()
	request.Topic = topic
	request.InstanceId = instanceId
	request.RegionId = client.RegionId

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DeleteTopic(request)
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

	return WrapError(alikafkaService.WaitForAlikafkaTopic(d.Id(), Deleted, DefaultTimeoutMedium))
}

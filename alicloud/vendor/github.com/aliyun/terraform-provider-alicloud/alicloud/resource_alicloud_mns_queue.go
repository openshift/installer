package alicloud

import (
	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudMNSQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMNSQueueCreate,
		Read:   resourceAlicloudMNSQueueRead,
		Update: resourceAlicloudMNSQueueUpdate,
		Delete: resourceAlicloudMNSQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 256),
			},
			"delay_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 604800),
			},
			"maximum_message_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      65536,
				ValidateFunc: validation.IntBetween(1024, 65536),
			},
			"message_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      345600,
				ValidateFunc: validation.IntBetween(60, 604800),
			},
			"visibility_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validation.IntBetween(1, 43200),
			},
			"polling_wait_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 1800),
			},
		},
	}
}

func resourceAlicloudMNSQueueCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	name := d.Get("name").(string)
	var delaySeconds, maximumMessageSize, messageRetentionPeriod, visibilityTimeout, pollingWaitSeconds int
	if v, ok := d.GetOk("delay_seconds"); ok {
		delaySeconds = v.(int)
	}
	if v, ok := d.GetOk("maximum_message_size"); ok {
		maximumMessageSize = v.(int)
	}
	if v, ok := d.GetOk("message_retention_period"); ok {
		messageRetentionPeriod = v.(int)
	}
	if v, ok := d.GetOk("visibility_timeout"); ok {
		visibilityTimeout = v.(int)
	}
	if v, ok := d.GetOk("polling_wait_seconds"); ok {
		pollingWaitSeconds = v.(int)
	}

	raw, err := client.WithMnsQueueManager(func(queueManager ali_mns.AliQueueManager) (interface{}, error) {
		return nil, queueManager.CreateQueue(name, int32(delaySeconds), int32(maximumMessageSize), int32(messageRetentionPeriod), int32(visibilityTimeout), int32(pollingWaitSeconds), 3)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mns_queue", "CreateQueue", AliMnsERROR)
	}
	addDebug("CreateQueue", raw)
	d.SetId(name)
	return resourceAlicloudMNSQueueRead(d, meta)
}

func resourceAlicloudMNSQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mnsService := MnsService{client}

	object, err := mnsService.DescribeMnsQueue(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.QueueName)
	d.Set("delay_seconds", object.DelaySeconds)
	d.Set("maximum_message_size", object.MaxMessageSize)
	d.Set("message_retention_period", object.MessageRetentionPeriod)
	d.Set("visibility_timeout", object.VisibilityTimeout)
	d.Set("polling_wait_seconds", object.PollingWaitSeconds)

	return nil
}

func resourceAlicloudMNSQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	attributeUpdate := false
	var delaySeconds, maximumMessageSize, messageRetentionPeriod, visibilityTimeouts, pollingWaitSeconds int
	delaySeconds = d.Get("delay_seconds").(int)
	maximumMessageSize = d.Get("maximum_message_size").(int)
	messageRetentionPeriod = d.Get("message_retention_period").(int)
	visibilityTimeouts = d.Get("visibility_timeout").(int)
	pollingWaitSeconds = d.Get("polling_wait_seconds").(int)
	name := d.Id()
	if d.HasChange("delay_seconds") {
		attributeUpdate = true
	}

	if d.HasChange("maximum_message_size") {
		attributeUpdate = true
	}

	if d.HasChange("message_retention_period") {
		attributeUpdate = true
	}
	if d.HasChange("visibility_timeout") {
		attributeUpdate = true
	}
	if d.HasChange("polling_wait_seconds") {
		attributeUpdate = true
	}

	if attributeUpdate {
		raw, err := client.WithMnsQueueManager(func(queueManager ali_mns.AliQueueManager) (interface{}, error) {
			return nil, queueManager.SetQueueAttributes(name, int32(delaySeconds), int32(maximumMessageSize), int32(messageRetentionPeriod), int32(visibilityTimeouts), int32(pollingWaitSeconds), 3)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetQueueAttributes", AliMnsERROR)
		}
		addDebug("SetQueueAttributes", raw)
	}
	return resourceAlicloudMNSQueueRead(d, meta)
}

func resourceAlicloudMNSQueueDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mnsService := MnsService{client}
	name := d.Id()
	raw, err := client.WithMnsQueueManager(func(queueManager ali_mns.AliQueueManager) (interface{}, error) {
		return nil, queueManager.DeleteQueue(name)
	})
	if err != nil {
		if mnsService.QueueNotExistFunc(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteQueue", AliMnsERROR)
	}
	addDebug("DeleteQueue", raw)
	return WrapError(mnsService.WaitForMnsQueue(d.Id(), Deleted, DefaultTimeout))
}

package alicloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudMNSSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMNSSubscriptionCreate,
		Read:   resourceAlicloudMNSSubscriptionRead,
		Update: resourceAlicloudMNSSubscriptionUpdate,
		Delete: resourceAlicloudMNSSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 256),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 256),
			},

			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"filter_tag": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 16),
			},

			"notify_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ali_mns.BACKOFF_RETRY),
				ValidateFunc: validation.StringInSlice([]string{
					string(ali_mns.BACKOFF_RETRY), string(ali_mns.EXPONENTIAL_DECAY_RETRY),
				}, false),
			},

			"notify_content_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ali_mns.SIMPLIFIED),
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudMNSSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	topicName := d.Get("topic_name").(string)
	name := d.Get("name").(string)
	endpoint := d.Get("endpoint").(string)
	notifyStrategyStr := d.Get("notify_strategy").(string)
	notifyContentFormatStr := d.Get("notify_content_format").(string)
	var filterTag string
	if v, ok := d.GetOk("filter_tag"); ok {
		filterTag = v.(string)
	}
	notifyStrategy := ali_mns.NotifyStrategyType(notifyStrategyStr)
	notifyContentFormat := ali_mns.NotifyContentFormatType(notifyContentFormatStr)
	subRequest := ali_mns.MessageSubsribeRequest{
		Endpoint:            endpoint,
		FilterTag:           filterTag,
		NotifyStrategy:      notifyStrategy,
		NotifyContentFormat: notifyContentFormat,
	}
	raw, err := client.WithMnsSubscriptionManagerByTopicName(topicName, func(subscriptionManager ali_mns.AliMNSTopic) (interface{}, error) {
		return nil, subscriptionManager.Subscribe(name, subRequest)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mns_topic_subscription", "Subscribe", AliMnsERROR)
	}
	addDebug("Subscribe", raw)
	d.SetId(fmt.Sprintf("%s%s%s", topicName, COLON_SEPARATED, name))
	return resourceAlicloudMNSSubscriptionRead(d, meta)
}

func resourceAlicloudMNSSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mnsService := MnsService{client}

	object, err := mnsService.DescribeMnsTopicSubscription(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("topic_name", object.TopicName)
	d.Set("name", object.SubscriptionName)
	d.Set("endpoint", object.Endpoint)
	d.Set("filter_tag", object.FilterTag)
	d.Set("notify_strategy", object.NotifyStrategy)
	d.Set("notify_content_format", object.NotifyContentFormat)
	return nil
}

func resourceAlicloudMNSSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("notify_strategy") {
		client := meta.(*connectivity.AliyunClient)
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
		topicName, name := parts[0], parts[1]
		notifyStrategy := ali_mns.NotifyStrategyType(d.Get("notify_strategy").(string))
		raw, err := client.WithMnsSubscriptionManagerByTopicName(topicName, func(subscriptionManager ali_mns.AliMNSTopic) (interface{}, error) {
			return nil, subscriptionManager.SetSubscriptionAttributes(name, notifyStrategy)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetSubscriptionAttributes", AliMnsERROR)
		}
		addDebug("SetSubscriptionAttributes", raw)
	}
	return resourceAlicloudMNSSubscriptionRead(d, meta)
}

func resourceAlicloudMNSSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mnsService := MnsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	topicName, name := parts[0], parts[1]

	raw, err := client.WithMnsSubscriptionManagerByTopicName(topicName, func(subscriptionManager ali_mns.AliMNSTopic) (interface{}, error) {
		return nil, subscriptionManager.Unsubscribe(name)
	})
	if err != nil {
		if mnsService.TopicNotExistFunc(err) || mnsService.SubscriptionNotExistFunc(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "Unsubscribe", AliMnsERROR)
	}
	addDebug("Unsubscribe", raw)
	return WrapError(mnsService.WaitForMnsTopicSubscription(d.Id(), Deleted, DefaultTimeout))
}

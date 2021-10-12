package alicloud

import (
	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudMNSTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMNSTopicCreate,
		Read:   resourceAlicloudMNSTopicRead,
		Update: resourceAlicloudMNSTopicUpdate,
		Delete: resourceAlicloudMNSTopicDelete,
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

			"maximum_message_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      65536,
				ValidateFunc: validation.IntBetween(1024, 65536),
			},

			"logging_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAlicloudMNSTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	name := d.Get("name").(string)
	maximumMessageSize := d.Get("maximum_message_size").(int)
	loggingEnabled := d.Get("logging_enabled").(bool)
	raw, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
		return nil, topicManager.CreateTopic(name, int32(maximumMessageSize), loggingEnabled)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mns_topic", "CreateTopic", AliMnsERROR)
	}
	addDebug("CreateTopic", raw)
	d.SetId(name)
	return resourceAlicloudMNSTopicRead(d, meta)
}

func resourceAlicloudMNSTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mnsService := MnsService{client}
	object, err := mnsService.DescribeMnsTopic(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.TopicName)
	d.Set("maximum_message_size", object.MaxMessageSize)
	d.Set("logging_enabled", object.LoggingEnabled)
	return nil
}

func resourceAlicloudMNSTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	attributeUpdate := false
	maximumMessageSize := d.Get("maximum_message_size").(int)
	loggingEnabled := d.Get("logging_enabled").(bool)
	if d.HasChange("maximum_message_size") {
		attributeUpdate = true
	}

	if d.HasChange("logging_enabled") {
		attributeUpdate = true
	}

	if attributeUpdate {
		raw, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
			return nil, topicManager.SetTopicAttributes(d.Id(), int32(maximumMessageSize), loggingEnabled)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetTopicAttributes", AliMnsERROR)
		}
		addDebug("SetTopicAttributes", raw)
	}
	return resourceAlicloudMNSTopicRead(d, meta)
}

func resourceAlicloudMNSTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mnsService := MnsService{client}
	raw, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
		return nil, topicManager.DeleteTopic(d.Id())
	})

	if err != nil {
		if mnsService.TopicNotExistFunc(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteTopic", AliMnsERROR)
	}
	addDebug("DeleteTopic", raw)

	return WrapError(mnsService.WaitForMnsTopic(d.Id(), Deleted, DefaultTimeout))
}

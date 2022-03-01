package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudAmqpQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAmqpQueueCreate,
		Read:   resourceAlicloudAmqpQueueRead,
		Update: resourceAlicloudAmqpQueueUpdate,
		Delete: resourceAlicloudAmqpQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_delete_state": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"auto_expire_state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dead_letter_exchange": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dead_letter_routing_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclusive_state": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"max_length": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"maximum_priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"message_ttl": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"queue_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"virtual_host_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudAmqpQueueCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateQueue"
	request := make(map[string]interface{})
	conn, err := client.NewOnsproxyClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("auto_delete_state"); ok {
		request["AutoDeleteState"] = v
	}
	if v, ok := d.GetOk("auto_expire_state"); ok {
		request["AutoExpireState"] = v
	}
	if v, ok := d.GetOk("dead_letter_exchange"); ok {
		request["DeadLetterExchange"] = v
	}
	if v, ok := d.GetOk("dead_letter_routing_key"); ok {
		request["DeadLetterRoutingKey"] = v
	}
	if v, ok := d.GetOkExists("exclusive_state"); ok {
		request["ExclusiveState"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("max_length"); ok {
		request["MaxLength"] = v
	}
	if v, ok := d.GetOk("maximum_priority"); ok {
		request["MaximumPriority"] = v
	}
	if v, ok := d.GetOk("message_ttl"); ok {
		request["MessageTTL"] = v
	}
	request["QueueName"] = d.Get("queue_name")
	request["VirtualHost"] = d.Get("virtual_host_name")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_amqp_queue", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["VirtualHost"], ":", request["QueueName"]))

	return resourceAlicloudAmqpQueueRead(d, meta)
}
func resourceAlicloudAmqpQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	amqpOpenService := AmqpOpenService{client}
	object, err := amqpOpenService.DescribeAmqpQueue(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_amqp_queue amqpOpenService.DescribeAmqpQueue Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("instance_id", parts[0])
	d.Set("queue_name", parts[2])
	d.Set("virtual_host_name", parts[1])
	d.Set("auto_delete_state", object["AutoDeleteState"])
	d.Set("exclusive_state", object["ExclusiveState"])
	return nil
}
func resourceAlicloudAmqpQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudAmqpQueueRead(d, meta)
}
func resourceAlicloudAmqpQueueDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteQueue"
	var response map[string]interface{}
	conn, err := client.NewOnsproxyClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId":  parts[0],
		"QueueName":   parts[2],
		"VirtualHost": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

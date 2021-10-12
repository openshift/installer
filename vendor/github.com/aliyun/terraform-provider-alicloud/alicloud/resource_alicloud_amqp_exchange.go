package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAmqpExchange() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAmqpExchangeCreate,
		Read:   resourceAlicloudAmqpExchangeRead,
		Update: resourceAlicloudAmqpExchangeUpdate,
		Delete: resourceAlicloudAmqpExchangeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alternate_exchange": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_delete_state": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"exchange_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"exchange_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DIRECT", "FANOUT", "HEADERS", "TOPIC"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internal": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"virtual_host_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudAmqpExchangeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateExchange"
	request := make(map[string]interface{})
	conn, err := client.NewOnsproxyClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("alternate_exchange"); ok {
		request["AlternateExchange"] = v
	}
	request["AutoDeleteState"] = d.Get("auto_delete_state")
	request["ExchangeName"] = d.Get("exchange_name")
	request["ExchangeType"] = d.Get("exchange_type")
	request["InstanceId"] = d.Get("instance_id")
	request["Internal"] = d.Get("internal")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_amqp_exchange", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["VirtualHost"], ":", request["ExchangeName"]))

	return resourceAlicloudAmqpExchangeRead(d, meta)
}
func resourceAlicloudAmqpExchangeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	amqpOpenService := AmqpOpenService{client}
	object, err := amqpOpenService.DescribeAmqpExchange(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_amqp_exchange amqpOpenService.DescribeAmqpExchange Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("exchange_name", parts[2])
	d.Set("instance_id", parts[0])
	d.Set("virtual_host_name", parts[1])
	d.Set("auto_delete_state", object["AutoDeleteState"])
	d.Set("exchange_type", object["ExchangeType"])
	return nil
}
func resourceAlicloudAmqpExchangeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudAmqpExchangeRead(d, meta)
}
func resourceAlicloudAmqpExchangeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteExchange"
	var response map[string]interface{}
	conn, err := client.NewOnsproxyClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ExchangeName": parts[2],
		"InstanceId":   parts[0],
		"VirtualHost":  parts[1],
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

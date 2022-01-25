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

func resourceAlicloudAmqpBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAmqpBindingCreate,
		Read:   resourceAlicloudAmqpBindingRead,
		Delete: resourceAlicloudAmqpBindingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"argument": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"x-match:all", "x-match:any"}, false),
			},
			"binding_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"binding_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"EXCHANGE", "QUEUE"}, false),
			},
			"destination_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_exchange": {
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

func resourceAlicloudAmqpBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBinding"
	request := make(map[string]interface{})
	conn, err := client.NewOnsproxyClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("argument"); ok {
		request["Argument"] = v
	}
	request["BindingKey"] = d.Get("binding_key")
	request["BindingType"] = d.Get("binding_type")
	request["DestinationName"] = d.Get("destination_name")
	request["InstanceId"] = d.Get("instance_id")
	request["SourceExchange"] = d.Get("source_exchange")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_amqp_binding", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["VirtualHost"], ":", request["SourceExchange"], ":", request["DestinationName"]))

	return resourceAlicloudAmqpBindingRead(d, meta)
}
func resourceAlicloudAmqpBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	amqpOpenService := AmqpOpenService{client}
	object, err := amqpOpenService.DescribeAmqpBinding(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_amqp_binding amqpOpenService.DescribeAmqpBinding Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	d.Set("destination_name", parts[3])
	d.Set("instance_id", parts[0])
	d.Set("source_exchange", parts[2])
	d.Set("virtual_host_name", parts[1])
	d.Set("argument", object["Argument"])
	d.Set("binding_key", object["BindingKey"])
	d.Set("binding_type", object["BindingType"])
	return nil
}
func resourceAlicloudAmqpBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteBinding"
	var response map[string]interface{}
	conn, err := client.NewOnsproxyClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DestinationName": parts[3],
		"InstanceId":      parts[0],
		"SourceExchange":  parts[2],
		"VirtualHost":     parts[1],
	}

	request["BindingKey"] = d.Get("binding_key")
	request["BindingType"] = d.Get("binding_type")
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

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

func resourceAlicloudConfigDeliveryChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudConfigDeliveryChannelCreate,
		Read:   resourceAlicloudConfigDeliveryChannelRead,
		Update: resourceAlicloudConfigDeliveryChannelUpdate,
		Delete: resourceAlicloudConfigDeliveryChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"delivery_channel_assume_role_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delivery_channel_condition": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delivery_channel_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delivery_channel_target_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delivery_channel_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MNS", "OSS", "SLS"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
		},
	}
}

func resourceAlicloudConfigDeliveryChannelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "PutDeliveryChannel"
	request := make(map[string]interface{})
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	request["DeliveryChannelAssumeRoleArn"] = d.Get("delivery_channel_assume_role_arn")
	if v, ok := d.GetOk("delivery_channel_condition"); ok {
		request["DeliveryChannelCondition"] = v
	}

	if v, ok := d.GetOk("delivery_channel_name"); ok {
		request["DeliveryChannelName"] = v
	}

	request["DeliveryChannelTargetArn"] = d.Get("delivery_channel_target_arn")
	request["DeliveryChannelType"] = d.Get("delivery_channel_type")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-08"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"DeliveryChannelSlsUnreachableError"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)

		d.SetId(fmt.Sprint(response["DeliveryChannelId"]))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_delivery_channel", action, AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudConfigDeliveryChannelRead(d, meta)
}
func resourceAlicloudConfigDeliveryChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	object, err := configService.DescribeConfigDeliveryChannel(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_delivery_channel configService.DescribeConfigDeliveryChannel Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("delivery_channel_assume_role_arn", object["DeliveryChannelAssumeRoleArn"])
	d.Set("delivery_channel_condition", object["DeliveryChannelCondition"])
	d.Set("delivery_channel_name", object["DeliveryChannelName"])
	d.Set("delivery_channel_target_arn", object["DeliveryChannelTargetArn"])
	d.Set("delivery_channel_type", object["DeliveryChannelType"])
	d.Set("description", object["Description"])
	d.Set("status", formatInt(object["Status"]))
	return nil
}
func resourceAlicloudConfigDeliveryChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"DeliveryChannelId": d.Id(),
	}
	if d.HasChange("delivery_channel_assume_role_arn") {
		update = true
	}
	request["DeliveryChannelAssumeRoleArn"] = d.Get("delivery_channel_assume_role_arn")
	if d.HasChange("delivery_channel_target_arn") {
		update = true
	}
	request["DeliveryChannelTargetArn"] = d.Get("delivery_channel_target_arn")
	request["DeliveryChannelType"] = d.Get("delivery_channel_type")
	if d.HasChange("delivery_channel_condition") {
		update = true
		request["DeliveryChannelCondition"] = d.Get("delivery_channel_condition")
	}
	if d.HasChange("delivery_channel_name") {
		update = true
		request["DeliveryChannelName"] = d.Get("delivery_channel_name")
	}
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("status") {
		update = true
		request["Status"] = d.Get("status")
	}
	if update {
		action := "PutDeliveryChannel"
		conn, err := client.NewConfigClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-08"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"DeliveryChannelSlsUnreachableError"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudConfigDeliveryChannelRead(d, meta)
}
func resourceAlicloudConfigDeliveryChannelDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudConfigDeliveryChannel. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

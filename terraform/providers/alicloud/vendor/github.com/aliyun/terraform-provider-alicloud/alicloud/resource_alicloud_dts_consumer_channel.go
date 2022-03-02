package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDtsConsumerChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDtsConsumerChannelCreate,
		Read:   resourceAlicloudDtsConsumerChannelRead,
		Update: resourceAlicloudDtsConsumerChannelUpdate,
		Delete: resourceAlicloudDtsConsumerChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"consumer_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"consumer_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,

				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"consumer_group_password": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(8, 32),
			},
			"consumer_group_user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,

				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-zA-Z0-9_]{1,16}$"), "The length of the name is limited to `1` to `16` characters. It can contain uppercase and lowercase letters, numbers, underscores (_)"),
			},
			"dts_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudDtsConsumerChannelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateConsumerChannel"
	request := make(map[string]interface{})
	conn, err := client.NewDtsClient()
	if err != nil {
		return WrapError(err)
	}
	request["ConsumerGroupName"] = d.Get("consumer_group_name")
	request["ConsumerGroupPassword"] = d.Get("consumer_group_password")
	request["ConsumerGroupUserName"] = d.Get("consumer_group_user_name")
	request["DtsInstanceId"] = d.Get("dts_instance_id")

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.JobStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dts_consumer_channel", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DtsInstanceId"], ":", response["ConsumerGroupID"]))

	return resourceAlicloudDtsConsumerChannelRead(d, meta)
}
func resourceAlicloudDtsConsumerChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := DtsService{client}
	object, err := dtsService.DescribeDtsConsumerChannel(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dts_consumer_channel dtsService.DescribeDtsConsumerChannel Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("consumer_group_id", parts[1])
	d.Set("dts_instance_id", parts[0])
	d.Set("consumer_group_name", object["ConsumerGroupName"])
	d.Set("consumer_group_user_name", object["ConsumerGroupUserName"])
	return nil
}
func resourceAlicloudDtsConsumerChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"ConsumerGroupId": parts[1],
		"DtsInstanceId":   parts[0],
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("consumer_group_password") {
		update = true
		request["ConsumerGroupPassword"] = d.Get("consumer_group_password")
	}
	if update {
		action := "ModifyConsumerChannel"
		conn, err := client.NewDtsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.JobStatus"}) || NeedRetry(err) {
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
	}
	return resourceAlicloudDtsConsumerChannelRead(d, meta)
}
func resourceAlicloudDtsConsumerChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteConsumerChannel"
	var response map[string]interface{}
	conn, err := client.NewDtsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ConsumerGroupId": parts[1],
		"DtsInstanceId":   parts[0],
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.JobStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	time.Sleep(10 * time.Second)
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

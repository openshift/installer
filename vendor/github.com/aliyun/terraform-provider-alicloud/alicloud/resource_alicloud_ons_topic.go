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

func resourceAlicloudOnsTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOnsTopicCreate,
		Read:   resourceAlicloudOnsTopicRead,
		Update: resourceAlicloudOnsTopicUpdate,
		Delete: resourceAlicloudOnsTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"message_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 4, 5}),
			},
			"perm": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{2, 4, 6}),
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"tags": tagsSchema(),
			"topic_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"topic"},
				ValidateFunc:  validation.StringLenBetween(1, 64),
			},
			"topic": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'topic' has been deprecated from version 1.97.0. Use 'topic_name' instead.",
				ConflictsWith: []string{"topic_name"},
				ValidateFunc:  validation.StringLenBetween(1, 64),
			},
		},
	}
}

func resourceAlicloudOnsTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "OnsTopicCreate"
	request := make(map[string]interface{})
	conn, err := client.NewOnsClient()
	if err != nil {
		return WrapError(err)
	}
	request["InstanceId"] = d.Get("instance_id")
	request["MessageType"] = d.Get("message_type")
	if v, ok := d.GetOk("remark"); ok {
		request["Remark"] = v
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request["Topic"] = v
	} else if v, ok := d.GetOk("topic"); ok {
		request["Topic"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "topic" or "topic_name" must be set one!`))
	}

	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.User"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ons_topic", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["Topic"]))

	return resourceAlicloudOnsTopicUpdate(d, meta)
}
func resourceAlicloudOnsTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}
	object, err := onsService.DescribeOnsTopic(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ons_topic onsService.DescribeOnsTopic Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("instance_id", parts[0])
	d.Set("topic_name", parts[1])
	d.Set("topic", parts[1])
	d.Set("message_type", formatInt(object["MessageType"]))
	d.Set("remark", object["Remark"])
	d.Set("tags", tagsToMap(object["Tags"].(map[string]interface{})["Tag"]))

	onsTopicStatusObject, err := onsService.OnsTopicStatus(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("perm", formatInt(onsTopicStatusObject["Perm"]))
	return nil
}
func resourceAlicloudOnsTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := onsService.SetResourceTags(d, "TOPIC"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	if d.HasChange("perm") {
		request := map[string]interface{}{
			"InstanceId": parts[0],
			"Topic":      parts[1],
		}
		request["Perm"] = d.Get("perm")
		action := "OnsTopicUpdate"
		conn, err := client.NewOnsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
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
		d.SetPartial("perm")
	}
	d.Partial(false)
	return resourceAlicloudOnsTopicRead(d, meta)
}
func resourceAlicloudOnsTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "OnsTopicDelete"
	var response map[string]interface{}
	conn, err := client.NewOnsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId": parts[0],
		"Topic":      parts[1],
	}

	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.User"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

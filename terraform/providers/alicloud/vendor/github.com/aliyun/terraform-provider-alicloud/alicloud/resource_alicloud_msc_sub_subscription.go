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

func resourceAlicloudMscSubSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMscSubSubscriptionCreate,
		Read:   resourceAlicloudMscSubSubscriptionRead,
		Update: resourceAlicloudMscSubSubscriptionUpdate,
		Delete: resourceAlicloudMscSubSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"channel": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"contact_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{-1, -2, 0, 1}),
			},
			"item_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pmsg_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{-1, -2, 0, 1}),
			},
			"sms_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{-1, -2, 0, 1}),
			},
			"tts_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{-1, -2, 0, 1}),
			},
			"webhook_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"webhook_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{-1, -2, 0, 1}),
			},
		},
	}
}

func resourceAlicloudMscSubSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSubscriptionItem"
	request := make(map[string]interface{})
	conn, err := client.NewMscopensubscriptionClient()
	if err != nil {
		return WrapError(err)
	}
	request["ItemName"] = d.Get("item_name")
	request["Locale"] = "en"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-07-13"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_msc_sub_subscription", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	responseSubscriptionItem := response["SubscriptionItem"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseSubscriptionItem["ItemId"]))

	return resourceAlicloudMscSubSubscriptionUpdate(d, meta)
}
func resourceAlicloudMscSubSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mscOpenSubscriptionService := MscOpenSubscriptionService{client}
	object, err := mscOpenSubscriptionService.DescribeMscSubSubscription(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_msc_sub_subscription mscOpenSubscriptionService.DescribeMscSubSubscription Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("channel", object["Channel"])

	if contactIds, ok := object["ContactIds"]; ok && contactIds != nil {
		d.Set("contact_ids", convertJsonStringToStringList(contactIds))
	}
	d.Set("description", object["Description"])
	if v, ok := object["EmailStatus"]; ok {
		d.Set("email_status", formatInt(v))
	}
	d.Set("item_name", object["ItemName"])
	if v, ok := object["PmsgStatus"]; ok {
		d.Set("pmsg_status", formatInt(v))
	}
	if v, ok := object["SmsStatus"]; ok {
		d.Set("sms_status", formatInt(v))
	}
	if v, ok := object["TtsStatus"]; ok {
		d.Set("tts_status", formatInt(v))
	}

	if webhookIds, ok := object["WebhookIds"]; ok && webhookIds != nil {
		d.Set("webhook_ids", convertJsonStringToStringList(webhookIds))
	}
	if v, ok := object["WebhookStatus"]; ok {
		d.Set("webhook_status", formatInt(v))
	}
	return nil
}
func resourceAlicloudMscSubSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ItemId": d.Id(),
	}
	if d.HasChange("email_status") {
		update = true
	}
	if v, ok := d.GetOkExists("email_status"); ok {
		request["EmailStatus"] = v
	}
	if d.HasChange("pmsg_status") {
		update = true
	}
	if v, ok := d.GetOkExists("pmsg_status"); ok {
		request["PmsgStatus"] = v
	}
	if d.HasChange("sms_status") {
		update = true
	}
	if v, ok := d.GetOkExists("sms_status"); ok {
		request["SmsStatus"] = v
	}
	if d.HasChange("tts_status") {
		update = true
	}
	if v, ok := d.GetOkExists("tts_status"); ok {
		request["TtsStatus"] = v
	}
	if d.HasChange("webhook_status") {
		update = true
	}
	if v, ok := d.GetOkExists("webhook_status"); ok {
		request["WebhookStatus"] = v
	}
	if d.HasChange("contact_ids") {
		update = true
		if v, ok := d.GetOk("contact_ids"); ok {
			request["ContactIds"] = convertListToJsonString(v.([]interface{}))
		}
	}
	request["Locale"] = "en"
	request["RegionId"] = client.RegionId
	if d.HasChange("webhook_ids") {
		update = true
		if v, ok := d.GetOk("webhook_ids"); ok {
			request["WebhookIds"] = convertListToJsonString(v.([]interface{}))
		}
	}

	if update {
		action := "UpdateSubscriptionItem"
		conn, err := client.NewMscopensubscriptionClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateSubscriptionItem")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-07-13"), StringPointer("AK"), nil, request, &runtime)
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}

	return resourceAlicloudMscSubSubscriptionRead(d, meta)
}
func resourceAlicloudMscSubSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudMscSubSubscription. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

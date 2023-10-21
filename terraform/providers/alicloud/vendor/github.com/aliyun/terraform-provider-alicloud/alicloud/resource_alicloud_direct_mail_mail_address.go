package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudDirectMailMailAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDirectMailMailAddressCreate,
		Read:   resourceAlicloudDirectMailMailAddressRead,
		Update: resourceAlicloudDirectMailMailAddressUpdate,
		Delete: resourceAlicloudDirectMailMailAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z0-9_-]+@[a-z0-9_-]+(\.[a-z0-9_-]+)+$`), "The email address must be filled in the format of account@domain, and only lowercase letters or numbers can be used."),
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9A-Za-z]{10,20}$`), "The password must be length 10-20 string, contains numbers, uppercase letters, lowercase letters at the same time."),
			},
			"reply_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sendtype": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"batch", "trigger"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDirectMailMailAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateMailAddress"
	request := make(map[string]interface{})
	conn, err := client.NewDmClient()
	if err != nil {
		return WrapError(err)
	}
	request["AccountName"] = d.Get("account_name")
	if v, ok := d.GetOk("reply_address"); ok {
		request["ReplyAddress"] = v
	}
	request["Sendtype"] = d.Get("sendtype")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-11-23"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_direct_mail_mail_address", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["MailAddressId"]))

	return resourceAlicloudDirectMailMailAddressUpdate(d, meta)
}
func resourceAlicloudDirectMailMailAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dmService := DmService{client}
	object, err := dmService.DescribeDirectMailMailAddress(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_direct_mail_mail_address dmService.DescribeDirectMailMailAddress Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("account_name", object["AccountName"])
	d.Set("reply_address", object["ReplyAddress"])
	d.Set("sendtype", object["Sendtype"])
	d.Set("status", object["AccountStatus"])
	return nil
}
func resourceAlicloudDirectMailMailAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"MailAddressId": d.Id(),
	}
	if d.HasChange("password") {
		update = true
		if v, ok := d.GetOk("password"); ok {
			request["Password"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("reply_address") {
		update = true
		if v, ok := d.GetOk("reply_address"); ok {
			request["ReplyAddress"] = v
		}
	}
	if update {
		action := "ModifyMailAddress"
		conn, err := client.NewDmClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-11-23"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	}
	return resourceAlicloudDirectMailMailAddressRead(d, meta)
}
func resourceAlicloudDirectMailMailAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMailAddress"
	var response map[string]interface{}
	conn, err := client.NewDmClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"MailAddressId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-11-23"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidMailAddressId.Malformed"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

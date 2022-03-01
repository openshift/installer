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

func resourceAlicloudMscSubContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMscSubContactCreate,
		Read:   resourceAlicloudMscSubContactRead,
		Update: resourceAlicloudMscSubContactUpdate,
		Delete: resourceAlicloudMscSubContactDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"contact_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[\u4E00-\u9FA5a-zA-Z]{2,12}$"), "The name must be 2 to 12 characters in length."),
			},
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`), "The email must has correct format."),
			},
			"mobile": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9]+$`), "The mobile only has digit."),
			},
			"position": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CEO", "Finance Director", "Maintenance Director", "Other", "Project Director", "Technical Director"}, false),
			},
		},
	}
}

func resourceAlicloudMscSubContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateContact"
	request := make(map[string]interface{})
	conn, err := client.NewMscopensubscriptionClient()
	if err != nil {
		return WrapError(err)
	}
	request["ContactName"] = d.Get("contact_name")
	request["Email"] = d.Get("email")
	request["Mobile"] = d.Get("mobile")
	request["Position"] = d.Get("position")
	request["ClientToken"] = buildClientToken("CreateContact")
	request["Locale"] = "en"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_msc_sub_contact", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(formatInt(response["ContactId"])))

	return resourceAlicloudMscSubContactRead(d, meta)
}
func resourceAlicloudMscSubContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mscOpenSubscriptionService := MscOpenSubscriptionService{client}
	object, err := mscOpenSubscriptionService.DescribeMscSubContact(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_msc_sub_contact mscOpenSubscriptionService.DescribeMscSubContact Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("contact_name", object["ContactName"])
	d.Set("email", object["Email"])
	d.Set("mobile", object["Mobile"])
	d.Set("position", convertMscSubContactPositionResponse(fmt.Sprint(object["Position"])))
	return nil
}
func resourceAlicloudMscSubContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ContactId": d.Id(),
	}
	if d.HasChange("contact_name") {
		update = true
	}
	request["ContactName"] = d.Get("contact_name")
	if d.HasChange("email") {
		update = true
	}
	request["Email"] = d.Get("email")
	if d.HasChange("mobile") {
		update = true
	}
	request["Mobile"] = d.Get("mobile")
	if update {
		action := "UpdateContact"
		conn, err := client.NewMscopensubscriptionClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateContact")
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
	return resourceAlicloudMscSubContactRead(d, meta)
}
func resourceAlicloudMscSubContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteContact"
	var response map[string]interface{}
	conn, err := client.NewMscopensubscriptionClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ContactId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if IsExpectedErrors(err, []string{"ResourceNotFound"}) {
		return nil
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}

func convertMscSubContactPositionResponse(source string) string {
	switch source {
	case "Others":
		return "Other"
	}
	return source
}

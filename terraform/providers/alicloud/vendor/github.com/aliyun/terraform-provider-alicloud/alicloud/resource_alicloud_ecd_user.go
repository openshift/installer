package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEcdUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcdUserCreate,
		Read:   resourceAlicloudEcdUserRead,
		Update: resourceAlicloudEcdUserUpdate,
		Delete: resourceAlicloudEcdUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"end_user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"phone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Unlocked", "Locked"}, false),
			},
		},
	}
}

func resourceAlicloudEcdUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateUsers"
	request := make(map[string]interface{})
	conn, err := client.NewEdsuserClient()
	if err != nil {
		return WrapError(err)
	}

	requestUsers := make(map[string]interface{})
	requestUsersMap := make([]interface{}, 0)

	requestUsers["Email"] = d.Get("email")
	requestUsers["EndUserId"] = d.Get("end_user_id")
	if v, ok := d.GetOk("password"); ok {
		requestUsers["Password"] = v
	}
	if v, ok := d.GetOk("phone"); ok {
		requestUsers["Phone"] = v
	}

	requestUsersMap = append(requestUsersMap, requestUsers)
	request["Users"] = requestUsersMap

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"Forbidden"}) {
				conn.Endpoint = String(connectivity.EcdOpenAPIEndpointUser)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_user", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(requestUsers["EndUserId"]))
	edsUserService := EdsUserService{client}

	stateConf := BuildStateConf([]string{}, []string{"Unlocked", "Locked"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edsUserService.EcdUserStateRefreshFunc(d.Id()))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEcdUserUpdate(d, meta)
}
func resourceAlicloudEcdUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edsUserService := EdsUserService{client}
	object, err := edsUserService.DescribeEcdUser(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_user edsUserService.DescribeEcdUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("end_user_id", d.Id())
	d.Set("email", object["Email"])
	d.Set("phone", object["Phone"])
	d.Set("status", convertEcdUserStatusResponse(fmt.Sprint(formatInt(object["Status"]))))
	return nil
}
func resourceAlicloudEcdUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edsUserService := EdsUserService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("status") {
		object, err := edsUserService.DescribeEcdUser(d.Id())
		if err != nil {
			return WrapError(err)
		}
		if target, ok := d.GetOk("status"); ok {
			if convertEcdUserStatusRequest(strconv.Itoa(formatInt(object["Status"]))) != target {
				if target == "Unlocked" {
					request := map[string]interface{}{
						"Users": []string{d.Id()},
					}

					action := "UnlockUsers"
					conn, err := client.NewEdsuserClient()
					if err != nil {
						return WrapError(err)
					}
					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
						response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
						if err != nil {
							if NeedRetry(err) {
								wait()
								return resource.RetryableError(err)
							}
							if IsExpectedErrors(err, []string{"Forbidden"}) {
								conn.Endpoint = String(connectivity.EcdOpenAPIEndpointUser)
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
				if target == "Locked" {
					request := map[string]interface{}{
						"Users": []string{d.Id()},
					}

					action := "LockUsers"
					conn, err := client.NewEdsuserClient()
					if err != nil {
						return WrapError(err)
					}
					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
						response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
						if err != nil {
							if NeedRetry(err) {
								wait()
								return resource.RetryableError(err)
							}
							if IsExpectedErrors(err, []string{"Forbidden"}) {
								conn.Endpoint = String(connectivity.EcdOpenAPIEndpointUser)
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
				d.SetPartial("status")
			}
		}

	}
	d.Partial(false)
	return resourceAlicloudEcdUserRead(d, meta)
}
func resourceAlicloudEcdUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RemoveUsers"
	var response map[string]interface{}
	conn, err := client.NewEdsuserClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Users": []string{d.Id()},
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"Forbidden"}) {
				conn.Endpoint = String(connectivity.EcdOpenAPIEndpointUser)
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
func convertEcdUserStatusResponse(source string) string {
	switch source {
	case "0":
		return "Unlocked"
	case "9":
		return "Locked"
	}
	return source
}
func convertEcdUserStatusRequest(source string) string {
	switch source {
	case "Unlocked":
		return "0"
	case "Locked":
		return "9"
	}
	return source
}

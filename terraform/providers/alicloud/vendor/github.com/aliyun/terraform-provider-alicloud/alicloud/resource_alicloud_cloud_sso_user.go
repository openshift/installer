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

func resourceAlicloudCloudSsoUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudSsoUserCreate,
		Read:   resourceAlicloudCloudSsoUserRead,
		Update: resourceAlicloudCloudSsoUserUpdate,
		Delete: resourceAlicloudCloudSsoUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"first_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"last_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled"}, false),
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w-.@]{1,64}$`), "The name of the resource. The name must be 1 to 64 characters in length and  can contain letters, digits, underscores (_), and hyphens (-)."),
			},
		},
	}
}

func resourceAlicloudCloudSsoUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateUser"
	request := make(map[string]interface{})
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["DirectoryId"] = d.Get("directory_id")
	if v, ok := d.GetOk("display_name"); ok {
		request["DisplayName"] = v
	}
	if v, ok := d.GetOk("email"); ok {
		request["Email"] = v
	}
	if v, ok := d.GetOk("first_name"); ok {
		request["FirstName"] = v
	}
	if v, ok := d.GetOk("last_name"); ok {
		request["LastName"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("user_name"); ok {
		request["UserName"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_sso_user", action, AlibabaCloudSdkGoERROR)
	}
	responseUser := response["User"].(map[string]interface{})
	d.SetId(fmt.Sprint(request["DirectoryId"], ":", responseUser["UserId"]))

	return resourceAlicloudCloudSsoUserRead(d, meta)
}
func resourceAlicloudCloudSsoUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}
	object, err := cloudssoService.DescribeCloudSsoUser(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_user cloudssoService.DescribeCloudSsoUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("directory_id", parts[0])
	d.Set("user_id", parts[1])
	d.Set("description", object["Description"])
	d.Set("display_name", object["DisplayName"])
	d.Set("email", object["Email"])
	d.Set("first_name", object["FirstName"])
	d.Set("last_name", object["LastName"])
	d.Set("status", object["Status"])
	d.Set("user_name", object["UserName"])
	return nil
}
func resourceAlicloudCloudSsoUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"UserId":      parts[1],
	}
	if d.HasChange("status") {
		update = true
		if v, ok := d.GetOk("status"); ok {
			request["NewStatus"] = v
		}
	}
	if update {
		action := "UpdateUserStatus"
		conn, err := client.NewCloudssoClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("status")
	}
	update = false
	updateUserReq := map[string]interface{}{
		"DirectoryId": parts[0],
		"UserId":      parts[1],
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			updateUserReq["NewDescription"] = v
		}
	}
	if d.HasChange("display_name") {
		update = true
		if v, ok := d.GetOk("display_name"); ok {
			updateUserReq["NewDisplayName"] = v
		}
	}
	if d.HasChange("email") {
		update = true
		if v, ok := d.GetOk("email"); ok {
			updateUserReq["NewEmail"] = v
		}
	}
	if d.HasChange("first_name") {
		update = true
		if v, ok := d.GetOk("first_name"); ok {
			updateUserReq["NewFirstName"] = v
		}
	}
	if d.HasChange("last_name") {
		update = true
		if v, ok := d.GetOk("last_name"); ok {
			updateUserReq["NewLastName"] = v
		}
	}
	if update {
		action := "UpdateUser"
		conn, err := client.NewCloudssoClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, updateUserReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateUserReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("description")
		d.SetPartial("display_name")
		d.SetPartial("email")
		d.SetPartial("first_name")
		d.SetPartial("last_name")
	}
	d.Partial(false)
	return resourceAlicloudCloudSsoUserRead(d, meta)
}
func resourceAlicloudCloudSsoUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteUser"
	var response map[string]interface{}
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"UserId":      parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"DeletionConflict.User.AccessAssigment"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.User"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

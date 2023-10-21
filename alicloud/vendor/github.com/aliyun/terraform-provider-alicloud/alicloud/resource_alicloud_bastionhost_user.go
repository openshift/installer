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

func resourceAlicloudBastionhostUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBastionhostUserCreate,
		Read:   resourceAlicloudBastionhostUserRead,
		Update: resourceAlicloudBastionhostUserUpdate,
		Delete: resourceAlicloudBastionhostUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mobile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mobile_country_code": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AE", "AU", "CH", "CN", "DE", "GB", "HK", "ID", "IN", "JP", "KR", "MO", "MY", "PH", "RU", "SE", "SG", "TW", "US"}, false),
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"source": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Local", "Ram"}, false),
			},
			"source_user_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Frozen", "Normal"}, false),
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudBastionhostUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateUser"
	request := make(map[string]interface{})
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("comment"); ok {
		request["Comment"] = v
	}
	if v, ok := d.GetOk("display_name"); ok {
		request["DisplayName"] = v
	}
	if v, ok := d.GetOk("email"); ok {
		request["Email"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("mobile"); ok {
		request["Mobile"] = v
	} else if v, ok := d.GetOk("mobile_country_code"); ok && v.(string) == "" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "mobile", "mobile_country_code", d.Get("mobile_country_code")))
	}
	if v, ok := d.GetOk("mobile_country_code"); ok {
		request["MobileCountryCode"] = v
	}
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	} else if v, ok := d.GetOk("source"); ok && v.(string) == "Local" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "password", "source", d.Get("source")))
	}
	request["RegionId"] = client.RegionId
	request["Source"] = d.Get("source")
	if v, ok := d.GetOk("source_user_id"); ok {
		request["SourceUserId"] = v
	} else if v, ok := d.GetOk("source"); ok && v.(string) == "Ram" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "source_user_id", "source", d.Get("source")))
	}
	request["UserName"] = d.Get("user_name")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_user", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", response["UserId"]))

	return resourceAlicloudBastionhostUserUpdate(d, meta)
}
func resourceAlicloudBastionhostUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	yundunBastionhostService := YundunBastionhostService{client}
	object, err := yundunBastionhostService.DescribeBastionhostUser(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_bastionhost_user yundunBastionhostService.DescribeBastionhostUser Failed!!! %s", err)
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
	d.Set("user_id", parts[1])
	d.Set("comment", object["Comment"])
	d.Set("display_name", object["DisplayName"])
	d.Set("email", object["Email"])
	d.Set("mobile", object["Mobile"])
	d.Set("mobile_country_code", object["MobileCountryCode"])
	d.Set("source", object["Source"])
	d.Set("source_user_id", object["SourceUserId"])
	d.Set("status", convertArrayToString(object["UserState"], ""))
	d.Set("user_name", object["UserName"])
	return nil
}
func resourceAlicloudBastionhostUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	yundunBastionhostService := YundunBastionhostService{client}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"InstanceId": parts[0],
		"UserId":     parts[1],
	}
	if !d.IsNewResource() && d.HasChange("comment") {
		update = true
		if v, ok := d.GetOk("comment"); ok {
			request["Comment"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("display_name") {
		update = true
		if v, ok := d.GetOk("display_name"); ok {
			request["DisplayName"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("email") {
		update = true
		if v, ok := d.GetOk("email"); ok {
			request["Email"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("mobile") {
		update = true
		if v, ok := d.GetOk("mobile"); ok {
			request["Mobile"] = v
		} else if v, ok := d.GetOk("mobile_country_code"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "mobile", "mobile_country_code", d.Get("mobile_country_code")))
		}
		request["MobileCountryCode"] = d.Get("mobile_country_code")
	}
	if !d.IsNewResource() && d.HasChange("mobile_country_code") {
		update = true
		if v, ok := d.GetOk("mobile_country_code"); ok {
			request["MobileCountryCode"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("password") {
		update = true
		if v, ok := d.GetOk("password"); ok {
			request["Password"] = v
		} else if v, ok := d.GetOk("source"); ok && v.(string) == "Local" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "password", "source", d.Get("source")))
		}
	}
	request["RegionId"] = client.RegionId
	if update {
		action := "ModifyUser"
		conn, err := client.NewBastionhostClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("comment")
		d.SetPartial("display_name")
		d.SetPartial("email")
		d.SetPartial("mobile")
		d.SetPartial("mobile_country_code")
		d.SetPartial("password")
	}
	if d.HasChange("status") {
		object, err := yundunBastionhostService.DescribeBastionhostUser(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if convertArrayToString(object["UserState"], "") != target {
			if target == "Frozen" {
				request := map[string]interface{}{
					"InstanceId": parts[0],
					"UserIds":    convertListToJsonString([]interface{}{parts[1]}),
				}
				request["RegionId"] = client.RegionId
				action := "LockUsers"
				conn, err := client.NewBastionhostClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			if target == "Normal" {
				request := map[string]interface{}{
					"InstanceId": parts[0],
					"UserIds":    convertListToJsonString([]interface{}{parts[1]}),
				}
				request["RegionId"] = client.RegionId
				action := "UnlockUsers"
				conn, err := client.NewBastionhostClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudBastionhostUserRead(d, meta)
}
func resourceAlicloudBastionhostUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteUser"
	var response map[string]interface{}
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId": parts[0],
		"UserId":     parts[1],
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"OBJECT_NOT_FOUND"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

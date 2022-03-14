package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudArmsAlertContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudArmsAlertContactCreate,
		Read:   resourceAlicloudArmsAlertContactRead,
		Update: resourceAlicloudArmsAlertContactUpdate,
		Delete: resourceAlicloudArmsAlertContactDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alert_contact_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ding_robot_webhook_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"phone_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_noc": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudArmsAlertContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAlertContact"
	request := make(map[string]interface{})
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("alert_contact_name"); ok {
		request["ContactName"] = v
	}
	if v, ok := d.GetOk("ding_robot_webhook_url"); ok {
		request["DingRobotWebhookUrl"] = v
	} else if v, ok := d.GetOk("email"); ok && v.(string) == "" {
		if v, ok := d.GetOk("phone_num"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "ding_robot_webhook_url", "email", d.Get("email"), "phone_num", d.Get("phone_num")))
		}
	}
	if v, ok := d.GetOk("email"); ok {
		request["Email"] = v
	} else if v, ok := d.GetOk("ding_robot_webhook_url"); ok && v.(string) == "" {
		if v, ok := d.GetOk("phone_num"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "email", "ding_robot_webhook_url", d.Get("ding_robot_webhook_url"), "phone_num", d.Get("phone_num")))
		}
	}
	if v, ok := d.GetOk("phone_num"); ok {
		request["PhoneNum"] = v
	} else if v, ok := d.GetOk("ding_robot_webhook_url"); ok && v.(string) == "" {
		if v, ok := d.GetOk("email"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "phone_num", "ding_robot_webhook_url", d.Get("ding_robot_webhook_url"), "email", d.Get("email")))
		}
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOkExists("system_noc"); ok {
		request["SystemNoc"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_alert_contact", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ContactId"]))

	return resourceAlicloudArmsAlertContactRead(d, meta)
}
func resourceAlicloudArmsAlertContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := ArmsService{client}
	object, err := armsService.DescribeArmsAlertContact(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_alert_contact armsService.DescribeArmsAlertContact Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("alert_contact_name", object["ContactName"])
	d.Set("ding_robot_webhook_url", object["DingRobot"])
	d.Set("email", object["Email"])
	d.Set("phone_num", object["Phone"])
	d.Set("system_noc", object["SystemNoc"])
	return nil
}
func resourceAlicloudArmsAlertContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ContactId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("alert_contact_name") {
		update = true
	}
	if v, ok := d.GetOk("alert_contact_name"); ok {
		request["ContactName"] = v
	}
	if d.HasChange("ding_robot_webhook_url") {
		update = true
	}
	if v, ok := d.GetOk("ding_robot_webhook_url"); ok {
		request["DingRobotWebhookUrl"] = v
	} else if v, ok := d.GetOk("email"); ok && v.(string) == "" {
		if v, ok := d.GetOk("phone_num"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "ding_robot_webhook_url", "email", d.Get("email"), "phone_num", d.Get("phone_num")))
		}
	}
	if d.HasChange("email") {
		update = true
	}
	if v, ok := d.GetOk("email"); ok {
		request["Email"] = v
	} else if v, ok := d.GetOk("ding_robot_webhook_url"); ok && v.(string) == "" {
		if v, ok := d.GetOk("phone_num"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "email", "ding_robot_webhook_url", d.Get("ding_robot_webhook_url"), "phone_num", d.Get("phone_num")))
		}
	}
	if d.HasChange("phone_num") {
		update = true
	}
	if v, ok := d.GetOk("phone_num"); ok {
		request["PhoneNum"] = v
	} else if v, ok := d.GetOk("ding_robot_webhook_url"); ok && v.(string) == "" {
		if v, ok := d.GetOk("email"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "phone_num", "ding_robot_webhook_url", d.Get("ding_robot_webhook_url"), "email", d.Get("email")))
		}
	}
	if d.HasChange("system_noc") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("system_noc"); ok {
		request["SystemNoc"] = v
	}
	if update {
		action := "UpdateAlertContact"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudArmsAlertContactRead(d, meta)
}
func resourceAlicloudArmsAlertContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAlertContact"
	var response map[string]interface{}
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ContactId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

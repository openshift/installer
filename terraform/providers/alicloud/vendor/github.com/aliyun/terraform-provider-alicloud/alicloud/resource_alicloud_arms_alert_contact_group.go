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

func resourceAlicloudArmsAlertContactGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudArmsAlertContactGroupCreate,
		Read:   resourceAlicloudArmsAlertContactGroupRead,
		Update: resourceAlicloudArmsAlertContactGroupUpdate,
		Delete: resourceAlicloudArmsAlertContactGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alert_contact_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contact_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlicloudArmsAlertContactGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAlertContactGroup"
	request := make(map[string]interface{})
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	request["ContactGroupName"] = d.Get("alert_contact_group_name")
	if v, ok := d.GetOk("contact_ids"); ok {
		request["ContactIds"] = convertArrayToString(v.(*schema.Set).List(), " ")
	}
	request["RegionId"] = client.RegionId
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_alert_contact_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ContactGroupId"]))

	return resourceAlicloudArmsAlertContactGroupRead(d, meta)
}
func resourceAlicloudArmsAlertContactGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := ArmsService{client}
	object, err := armsService.DescribeArmsAlertContactGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_alert_contact_group armsService.DescribeArmsAlertContactGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("alert_contact_group_name", object["ContactGroupName"])
	contactIdsItems := make([]string, 0)
	if contacts, ok := object["Contacts"]; ok && contacts != nil {
		for _, contactsItem := range contacts.([]interface{}) {
			if contactId, ok := contactsItem.(map[string]interface{})["ContactId"]; ok && contactId != nil {
				contactIdsItems = append(contactIdsItems, fmt.Sprint(contactId))
			}
		}
	}
	d.Set("contact_ids", contactIdsItems)
	return nil
}
func resourceAlicloudArmsAlertContactGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ContactGroupId": d.Id(),
	}
	if d.HasChange("alert_contact_group_name") {
		update = true
	}
	request["ContactGroupName"] = d.Get("alert_contact_group_name")
	request["RegionId"] = client.RegionId
	if d.HasChange("contact_ids") {
		update = true
	}
	if v, ok := d.GetOk("contact_ids"); ok {
		request["ContactIds"] = convertArrayToString(v.(*schema.Set).List(), " ")
	}
	if update {
		action := "UpdateAlertContactGroup"
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
	return resourceAlicloudArmsAlertContactGroupRead(d, meta)
}
func resourceAlicloudArmsAlertContactGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAlertContactGroup"
	var response map[string]interface{}
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ContactGroupId": d.Id(),
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

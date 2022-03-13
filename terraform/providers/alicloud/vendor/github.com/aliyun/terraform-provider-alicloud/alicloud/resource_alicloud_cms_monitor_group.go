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

func resourceAlicloudCmsMonitorGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsMonitorGroupCreate,
		Read:   resourceAlicloudCmsMonitorGroupRead,
		Update: resourceAlicloudCmsMonitorGroupUpdate,
		Delete: resourceAlicloudCmsMonitorGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"contact_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"monitor_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudCmsMonitorGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewCmsClient()
	var response map[string]interface{}
	request := make(map[string]interface{})
	if v, exist := d.GetOk("resource_group_id"); exist {
		action := "CreateMonitorGroupByResourceGroupId"
		request["RegionId"] = client.RegionId
		request["ResourceGroupId"] = v.(string)
		request["ResourceGroupName"] = d.Get("resource_group_name")
		for k, v := range d.Get("contact_groups").([]interface{}) {
			request[fmt.Sprintf("ContactGroupList.%d", k+1)] = v.(string)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_monitor_group", action, AlibabaCloudSdkGoERROR)
		}
		d.SetId(fmt.Sprint(response["Id"]))
		return resourceAlicloudCmsMonitorGroupUpdate(d, meta)
	}

	action := "CreateMonitorGroup"
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("contact_groups"); ok && v != nil {
		request["ContactGroups"] = convertListToCommaSeparate(v.([]interface{}))
	}
	request["GroupName"] = d.Get("monitor_group_name")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_monitor_group", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		return WrapError(Error("CreateMonitorGroup failed for " + response["Message"].(string)))
	}

	d.SetId(fmt.Sprint(formatInt(response["GroupId"])))

	return resourceAlicloudCmsMonitorGroupUpdate(d, meta)
}
func resourceAlicloudCmsMonitorGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsMonitorGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_monitor_group cmsService.DescribeCmsMonitorGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	var contactGroups []interface{}
	if v, ok := object["ContactGroups"].(map[string]interface{})["ContactGroup"]; ok {
		for _, contactGroup := range v.([]interface{}) {
			contactGroups = append(contactGroups, contactGroup.(map[string]interface{})["Name"])
		}
	}
	d.Set("contact_groups", contactGroups)
	d.Set("monitor_group_name", object["GroupName"])
	d.Set("tags", tagsToMap(object["Tags"].(map[string]interface{})["Tag"]))
	return nil
}
func resourceAlicloudCmsMonitorGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"GroupId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("contact_groups") {
		update = true
		request["ContactGroups"] = convertListToCommaSeparate(d.Get("contact_groups").([]interface{}))
	}
	if d.HasChange("monitor_group_name") {
		update = true
		request["GroupName"] = d.Get("monitor_group_name")
	}
	if update {
		action := "ModifyMonitorGroup"
		conn, err := client.NewCmsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if !response["Success"].(bool) {
			return WrapError(Error("ModifyMonitorGroup failed for " + response["Message"].(string)))
		}
		d.SetPartial("contact_groups")
		d.SetPartial("monitor_group_name")
	}
	if d.HasChange("tags") {
		if err := cmsService.SetResourceTags(d, ""); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAlicloudCmsMonitorGroupRead(d, meta)
}
func resourceAlicloudCmsMonitorGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMonitorGroup"
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"GroupId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["Code"]), []string{"400", "403", "404", "ResourceNotFound"}) {
		return nil
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		return WrapError(Error("DeleteMonitorGroup failed for " + response["Message"].(string)))
	}
	return nil
}

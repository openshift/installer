package alicloud

import (
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSecurityCenterGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSecurityCenterGroupCreate,
		Read:   resourceAlicloudSecurityCenterGroupRead,
		Update: resourceAlicloudSecurityCenterGroupUpdate,
		Delete: resourceAlicloudSecurityCenterGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
func resourceAlicloudSecurityCenterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateOrUpdateAssetGroup"
	request := make(map[string]interface{})
	conn, err := client.NewSasClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("group_id"); ok {
		request["GroupId"] = v
	}
	if v, ok := d.GetOk("group_name"); ok {
		request["GroupName"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_security_center_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["GroupId"]))
	return resourceAlicloudSecurityCenterGroupRead(d, meta)
}
func resourceAlicloudSecurityCenterGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}
	describeAllGroupsObject, err := sasService.DescribeAllGroups(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("group_name", describeAllGroupsObject["GroupName"])
	d.Set("group_id", describeAllGroupsObject["GroupId"])
	return nil
}
func resourceAlicloudSecurityCenterGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var response map[string]interface{}
	request := map[string]interface{}{
		"GroupId": d.Id(),
	}
	if v, ok := d.GetOk("group_name"); ok {
		request["GroupName"] = v
	}

	action := "CreateOrUpdateAssetGroup"
	conn, err := client.NewSasClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudSecurityCenterGroupRead(d, meta)
}
func resourceAlicloudSecurityCenterGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteGroup"
	var response map[string]interface{}
	conn, err := client.NewSasClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"GroupId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

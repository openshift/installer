package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCloudSsoGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudSsoGroupCreate,
		Read:   resourceAlicloudCloudSsoGroupRead,
		Update: resourceAlicloudCloudSsoGroupUpdate,
		Delete: resourceAlicloudCloudSsoGroupDelete,
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
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w-.]{1,128}$`), "The name of the resource. The name must be 1 to 128 characters in length and  can contain letters, digits, underscores (_), and hyphens (-)."),
			},
		},
	}
}

func resourceAlicloudCloudSsoGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateGroup"
	request := make(map[string]interface{})
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["DirectoryId"] = d.Get("directory_id")
	request["GroupName"] = d.Get("group_name")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_sso_group", action, AlibabaCloudSdkGoERROR)
	}
	responseGroup := response["Group"].(map[string]interface{})
	d.SetId(fmt.Sprint(request["DirectoryId"], ":", responseGroup["GroupId"]))

	return resourceAlicloudCloudSsoGroupRead(d, meta)
}
func resourceAlicloudCloudSsoGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}
	object, err := cloudssoService.DescribeCloudSsoGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_group cloudssoService.DescribeCloudSsoGroup Failed!!! %s", err)
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
	d.Set("group_id", parts[1])
	d.Set("description", object["Description"])
	d.Set("group_name", object["GroupName"])
	return nil
}
func resourceAlicloudCloudSsoGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"GroupId":     parts[1],
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["NewDescription"] = v
		}
	}
	if d.HasChange("group_name") {
		update = true
		if v, ok := d.GetOk("group_name"); ok {
			request["NewGroupName"] = v
		}
	}
	if update {
		action := "UpdateGroup"
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
	}
	return resourceAlicloudCloudSsoGroupRead(d, meta)
}
func resourceAlicloudCloudSsoGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteGroup"
	var response map[string]interface{}
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"GroupId":     parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"DeletionConflict.Group.AccessAssigment"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Group"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

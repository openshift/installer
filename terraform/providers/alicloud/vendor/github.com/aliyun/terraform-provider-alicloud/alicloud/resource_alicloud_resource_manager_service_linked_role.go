package alicloud

import (
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudResourceManagerServiceLinkedRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerServiceLinkedRoleCreate,
		Read:   resourceAlicloudResourceManagerServiceLinkedRoleRead,
		Delete: resourceAlicloudResourceManagerServiceLinkedRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"custom_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerServiceLinkedRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateServiceLinkedRole"
	request := make(map[string]interface{})
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}

	request["ServiceName"] = d.Get("service_name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("custom_suffix"); ok {
		request["CustomSuffix"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_service_linked_role", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(d.Get("service_name"), ":", response["Role"].(map[string]interface{})["RoleName"]))

	return resourceAlicloudResourceManagerServiceLinkedRoleRead(d, meta)
}
func resourceAlicloudResourceManagerServiceLinkedRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	parts, _ := ParseResourceId(d.Id(), 2)

	object, err := ramService.DescribeRamServiceLinkedRole(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("role_id", object.Role.RoleId)
	d.Set("role_name", object.Role.RoleName)
	d.Set("arn", object.Role.Arn)
	d.Set("service_name", parts[0])
	return nil

}

func resourceAlicloudResourceManagerServiceLinkedRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "DeleteServiceLinkedRole"
	request := make(map[string]interface{})
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}
	parts, _ := ParseResourceId(d.Id(), 2)
	request["RoleName"] = parts[1]

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_service_linked_role", action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

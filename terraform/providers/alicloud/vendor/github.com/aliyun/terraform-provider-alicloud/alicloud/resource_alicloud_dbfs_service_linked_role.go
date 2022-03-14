package alicloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDbfsServiceLinkedRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDbfsServiceLinkedRoleCreate,
		Read:   resourceAlicloudDbfsServiceLinkedRoleRead,
		Delete: resourceAlicloudDbfsServiceLinkedRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"product_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AliyunServiceRoleForDbfs",
				}, false),
			},
			"status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}
func resourceAlicloudDbfsServiceLinkedRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateServiceLinkedRole"
	request := map[string]interface{}{}
	conn, err := client.NewDbfsClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	request["ClientToken"] = buildClientToken("CreateServiceLinkedRole")
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	d.SetId(d.Get("product_name").(string))
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityAlreadyExist.Role"}) {
			return resourceAlicloudDbfsServiceLinkedRoleRead(d, meta)
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dbfs_service_linked_role", action, AlibabaCloudSdkGoERROR)
	}
	return resourceAlicloudDbfsServiceLinkedRoleRead(d, meta)
}

func resourceAlicloudDbfsServiceLinkedRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbfsService := DbfsService{client}
	object, err := dbfsService.DescribeDbfsServiceLinkedRole(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dbfs_event_bus dbfsService.DescribeDbfsServiceLinkedRole Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("product_name", d.Id())
	d.Set("status", object["DbfsLinkedRole"])
	return nil
}

func resourceAlicloudDbfsServiceLinkedRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteServiceLinkedRole"
	var response map[string]interface{}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RoleName": d.Id(),
	}
	wait := incrementalWait(3*time.Second, 0*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

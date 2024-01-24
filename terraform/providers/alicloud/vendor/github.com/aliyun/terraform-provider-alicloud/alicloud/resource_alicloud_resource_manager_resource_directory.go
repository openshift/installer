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

func resourceAlicloudResourceManagerResourceDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerResourceDirectoryCreate,
		Read:   resourceAlicloudResourceManagerResourceDirectoryRead,
		Update: resourceAlicloudResourceManagerResourceDirectoryUpdate,
		Delete: resourceAlicloudResourceManagerResourceDirectoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"master_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_folder_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled"}, false),
			},
		},
	}
}

func resourceAlicloudResourceManagerResourceDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "InitResourceDirectory"
	request := make(map[string]interface{})
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
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
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_resource_directory", action, AlibabaCloudSdkGoERROR)
	}
	responseResourceDirectory := response["ResourceDirectory"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseResourceDirectory["ResourceDirectoryId"]))

	return resourceAlicloudResourceManagerResourceDirectoryUpdate(d, meta)
}
func resourceAlicloudResourceManagerResourceDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerResourceDirectory(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_resource_directory resourcemanagerService.DescribeResourceManagerResourceDirectory Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("master_account_id", object["MasterAccountId"])
	d.Set("master_account_name", object["MasterAccountName"])
	d.Set("root_folder_id", object["RootFolderId"])
	d.Set("status", object["ScpStatus"])
	return nil
}
func resourceAlicloudResourceManagerResourceDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("status") {
		object, err := resourcemanagerService.DescribeResourceManagerResourceDirectory(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["ScpStatus"].(string) != target {
			if target == "Disabled" {
				request := map[string]interface{}{}
				action := "DisableControlPolicy"
				conn, err := client.NewResourcemanagerClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				stateConf := BuildStateConf([]string{}, []string{"Disabled"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, resourcemanagerService.ResourceManagerResourceDirectoryStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "Enabled" {
				request := map[string]interface{}{}
				action := "EnableControlPolicy"
				conn, err := client.NewResourcemanagerClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				stateConf := BuildStateConf([]string{}, []string{"Enabled"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, resourcemanagerService.ResourceManagerResourceDirectoryStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudResourceManagerResourceDirectoryRead(d, meta)
}
func resourceAlicloudResourceManagerResourceDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DestroyResourceDirectory"
	var response map[string]interface{}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{}

	wait := incrementalWait(3*time.Second, 0*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2020-03-31"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
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
	return nil
}

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

func resourceAlicloudResourceManagerSharedTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerSharedTargetCreate,
		Read:   resourceAlicloudResourceManagerSharedTargetRead,
		Delete: resourceAlicloudResourceManagerSharedTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Delete: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"resource_share_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerSharedTargetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcesharingService := ResourcesharingService{client}
	var response map[string]interface{}
	action := "AssociateResourceShare"
	request := make(map[string]interface{})
	conn, err := client.NewRessharingClient()
	if err != nil {
		return WrapError(err)
	}
	request["ResourceShareId"] = d.Get("resource_share_id")
	request["Targets"] = []string{d.Get("target_id").(string)}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_shared_target", action, AlibabaCloudSdkGoERROR)
	}
	response = response["ResourceShareAssociations"].([]interface{})[0].(map[string]interface{})
	d.SetId(fmt.Sprint(response["ResourceShareId"], ":", response["EntityId"]))
	stateConf := BuildStateConf([]string{}, []string{"Associated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, resourcesharingService.ResourceManagerSharedTargetStateRefreshFunc(d.Id(), []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudResourceManagerSharedTargetRead(d, meta)
}
func resourceAlicloudResourceManagerSharedTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcesharingService := ResourcesharingService{client}
	object, err := resourcesharingService.DescribeResourceManagerSharedTarget(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_shared_target resourcesharingService.DescribeResourceManagerSharedTarget Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("resource_share_id", parts[0])
	d.Set("target_id", parts[1])
	d.Set("status", object["AssociationStatus"])
	return nil
}
func resourceAlicloudResourceManagerSharedTargetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	resourcesharingService := ResourcesharingService{client}
	action := "DisassociateResourceShare"
	var response map[string]interface{}
	conn, err := client.NewRessharingClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ResourceShareId": parts[0],
		"Targets":         []string{parts[1]},
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	stateConf := BuildStateConf([]string{}, []string{"Disassociated"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, resourcesharingService.ResourceManagerSharedTargetStateRefreshFunc(d.Id(), []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudResourceManagerPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerPolicyAttachmentCreate,
		Read:   resourceAlicloudResourceManagerPolicyAttachmentRead,
		Delete: resourceAlicloudResourceManagerPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Custom", "System"}, false),
			},
			"principal_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"principal_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IMSGroup", "IMSUser", "ServiceRole"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachPolicy"
	request := make(map[string]interface{})
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}
	request["PolicyName"] = d.Get("policy_name")
	request["PolicyType"] = d.Get("policy_type")
	request["PrincipalName"] = d.Get("principal_name")
	request["PrincipalType"] = d.Get("principal_type")
	request["ResourceGroupId"] = d.Get("resource_group_id")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_policy_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["PolicyName"], ":", request["PolicyType"], ":", request["PrincipalName"], ":", request["PrincipalType"], ":", request["ResourceGroupId"]))

	return resourceAlicloudResourceManagerPolicyAttachmentRead(d, meta)
}
func resourceAlicloudResourceManagerPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	if len(strings.Split(d.Id(), ":")) != 5 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), d.Get("resource_group_id").(string)))
	}
	_, err := resourcemanagerService.DescribeResourceManagerPolicyAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_policy_attachment resourcemanagerService.DescribeResourceManagerPolicyAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 5)
	if err != nil {
		return WrapError(err)
	}
	d.Set("policy_name", parts[0])
	d.Set("policy_type", parts[1])
	d.Set("principal_name", parts[2])
	d.Set("principal_type", parts[3])
	d.Set("resource_group_id", parts[4])
	return nil
}
func resourceAlicloudResourceManagerPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if len(strings.Split(d.Id(), ":")) != 5 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), d.Get("resource_group_id").(string)))
	}
	parts, err := ParseResourceId(d.Id(), 5)
	if err != nil {
		return WrapError(err)
	}
	action := "DetachPolicy"
	var response map[string]interface{}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PolicyName":      parts[0],
		"PolicyType":      parts[1],
		"PrincipalName":   parts[2],
		"PrincipalType":   parts[3],
		"ResourceGroupId": parts[4],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExists.ResourceGroup"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

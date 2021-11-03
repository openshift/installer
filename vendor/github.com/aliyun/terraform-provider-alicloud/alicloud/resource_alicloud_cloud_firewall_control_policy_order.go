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

func resourceAlicloudCloudFirewallControlPolicyOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudFirewallControlPolicyOrderCreate,
		Read:   resourceAlicloudCloudFirewallControlPolicyOrderRead,
		Update: resourceAlicloudCloudFirewallControlPolicyOrderUpdate,
		Delete: resourceAlicloudCloudFirewallControlPolicyOrderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"acl_uuid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"in", "out"}, false),
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudCloudFirewallControlPolicyOrderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyControlPolicyPriority"
	request := make(map[string]interface{})
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}
	request["Direction"] = d.Get("direction")
	request["Order"] = d.Get("order")
	request["AclUuid"] = d.Get("acl_uuid")

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_control_policy_order", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AclUuid"], ":", request["Direction"]))

	return resourceAlicloudCloudFirewallControlPolicyRead(d, meta)
}

func resourceAlicloudCloudFirewallControlPolicyOrderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "ModifyControlPolicyPriority"
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}

	update := false
	request := map[string]interface{}{
		"AclUuid":   parts[0],
		"Direction": parts[1],
	}
	if d.HasChange("order") {
		update = true
		request["Order"] = d.Get("order")
	}

	if update {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
	}
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_control_policy_order", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AclUuid"], ":", request["Direction"]))

	return resourceAlicloudCloudFirewallControlPolicyRead(d, meta)
}

func resourceAlicloudCloudFirewallControlPolicyOrderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}
	object, err := cloudfwService.DescribeCloudFirewallControlPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_control_policy_order cloudfwService.DescribeCloudFirewallControlPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	d.Set("acl_uuid", parts[0])
	d.Set("direction", parts[1])
	d.Set("order", formatInt(object["Order"]))

	return nil
}

func resourceAlicloudCloudFirewallControlPolicyOrderDelete(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] Resource alicloud_cloud_firewall_control_policy_order [%s]  will not be deleted", d.Id())
	return nil
}

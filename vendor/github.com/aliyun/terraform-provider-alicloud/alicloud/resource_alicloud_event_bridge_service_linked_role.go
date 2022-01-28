package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEventBridgeServiceLinkedRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEventBridgeServiceLinkedRoleCreate,
		Read:   resourceAlicloudEventBridgeServiceLinkedRoleRead,
		Delete: resourceAlicloudEventBridgeServiceLinkedRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"product_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AliyunServiceRoleForEventBridgeSourceRocketMQ",
					"AliyunServiceRoleForEventBridgeSourceActionTrail",
					"AliyunServiceRoleForEventBridgeSourceRabbitMQ",
					"AliyunServiceRoleForEventBridgeSendToMNS",
					"AliyunServiceRoleForEventBridgeConnectVPC",
					"AliyunServiceRoleForEventBridgeSendToFC",
					"AliyunServiceRoleForEventBridgeSendToSMS",
					"AliyunServiceRoleForEventBridgeSendToDirectMail",
					"AliyunServiceRoleForEventBridgeSendToRabbitMQ",
					"AliyunServiceRoleForEventBridgeSendToRocketMQ",
				}, false),
			},
		},
	}
}
func resourceAlicloudEventBridgeServiceLinkedRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventbridgeService := EventbridgeService{client}
	var response map[string]interface{}
	action := "CreateServiceLinkedRoleForProduct"
	request := map[string]interface{}{
		"ProductName": d.Get("product_name"),
	}
	conn, err := client.NewEventbridgeClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_event_bridge_service_linked_role", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("CreateServiceLinkedRoleForProduct failed, response: %v", response))
	}
	d.SetId(request["ProductName"].(string))

	stateConf := BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, eventbridgeService.CheckRoleForProductRefreshFunc(d.Id(), []string{"false"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEventBridgeServiceLinkedRoleRead(d, meta)
}

func resourceAlicloudEventBridgeServiceLinkedRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventbridgeService := EventbridgeService{client}
	object, err := eventbridgeService.DescribeEventBridgeServiceLinkedRole(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_event_bridge_event_bus eventbridgeService.DescribeEventBridgeServiceLinkedRole Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("product_name", object["StsRoleName"])
	return nil
}

func resourceAlicloudEventBridgeServiceLinkedRoleDelete(d *schema.ResourceData, meta interface{}) error {
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

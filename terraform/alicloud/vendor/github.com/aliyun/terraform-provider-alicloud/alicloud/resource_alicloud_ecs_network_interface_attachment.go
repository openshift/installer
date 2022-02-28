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

func resourceAlicloudEcsNetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsNetworkInterfaceAttachmentCreate,
		Read:   resourceAlicloudEcsNetworkInterfaceAttachmentRead,
		Update: resourceAlicloudEcsNetworkInterfaceAttachmentUpdate,
		Delete: resourceAlicloudEcsNetworkInterfaceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"trunk_network_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"wait_for_network_configuration_ready": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudEcsNetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "AttachNetworkInterface"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request["InstanceId"] = d.Get("instance_id")
	request["NetworkInterfaceId"] = d.Get("network_interface_id")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("trunk_network_instance_id"); ok {
		request["TrunkNetworkInstanceId"] = v
	}

	if v, ok := d.GetOkExists("wait_for_network_configuration_ready"); ok {
		request["WaitForNetworkConfigurationReady"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, NetworkInterfaceInvalidOperations) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_network_interface_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["NetworkInterfaceId"], ":", request["InstanceId"]))
	parts, err := ParseResourceId(d.Id(), 2)
	stateConf := BuildStateConf([]string{}, []string{"InUse"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsNetworkInterfaceStateRefreshFunc(parts[0], []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEcsNetworkInterfaceAttachmentRead(d, meta)
}
func resourceAlicloudEcsNetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeEcsNetworkInterfaceAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_netWork_interface_attachment ecsService.DescribeNetworkInterfaceAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object["InstanceId"])
	d.Set("network_interface_id", object["NetworkInterfaceId"])

	return nil
}
func resourceAlicloudEcsNetworkInterfaceAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudEcsNetworkInterfaceAttachmentRead(d, meta)
}
func resourceAlicloudEcsNetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DetachNetworkInterface"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId":         parts[1],
		"NetworkInterfaceId": parts[0],
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("trunk_network_instance_id"); ok {
		request["TrunkNetworkInstanceId"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, NetworkInterfaceInvalidOperations) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidEcsId.NotFound", "InvalidEniId.NotFound", "InvalidSecurityGroupId.NotFound", "InvalidVSwitchId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecsService.EcsNetworkInterfaceStateRefreshFunc(parts[0], []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

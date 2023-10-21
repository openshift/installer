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

func resourceAlicloudEcsAutoSnapshotPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsAutoSnapshotPolicyAttachmentCreate,
		Read:   resourceAlicloudEcsAutoSnapshotPolicyAttachmentRead,
		Delete: resourceAlicloudEcsAutoSnapshotPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_snapshot_policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEcsAutoSnapshotPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ApplyAutoSnapshotPolicy"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request["autoSnapshotPolicyId"] = d.Get("auto_snapshot_policy_id")
	request["diskIds"] = convertListToJsonString([]interface{}{d.Get("disk_id")})
	request["regionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_auto_snapshot_policy_attachment", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(request["autoSnapshotPolicyId"], ":", d.Get("disk_id")))
	return resourceAlicloudEcsAutoSnapshotPolicyAttachmentRead(d, meta)
}
func resourceAlicloudEcsAutoSnapshotPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	_, err := ecsService.DescribeEcsAutoSnapshotPolicyAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_auto_snapshot_policy_attachment ecsService.DescribeEcsAutoSnapshotPolicyAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)

	d.Set("auto_snapshot_policy_id", parts[0])
	d.Set("disk_id", parts[1])

	return nil
}
func resourceAlicloudEcsAutoSnapshotPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "CancelAutoSnapshotPolicy"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"diskIds": convertListToJsonString([]interface{}{parts[1]}),
	}

	request["regionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidDiskId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

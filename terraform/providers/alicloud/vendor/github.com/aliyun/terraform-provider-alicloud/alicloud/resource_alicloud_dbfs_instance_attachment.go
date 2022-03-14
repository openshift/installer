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

func resourceAlicloudDbfsInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDbfsInstanceAttachmentCreate,
		Read:   resourceAlicloudDbfsInstanceAttachmentRead,
		Delete: resourceAlicloudDbfsInstanceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ecs_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudDbfsInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachDbfs"
	request := make(map[string]interface{})
	conn, err := client.NewDbfsClient()
	if err != nil {
		return WrapError(err)
	}
	request["ECSInstanceId"] = d.Get("ecs_id")
	request["FsId"] = d.Get("instance_id")

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dbfs_instance_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["FsId"], ":", request["ECSInstanceId"]))
	dbfsService := DbfsService{client}
	stateConf := BuildStateConf([]string{}, []string{"attached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dbfsService.DbfsInstanceStateRefreshFunc(fmt.Sprint(request["FsId"]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudDbfsInstanceAttachmentRead(d, meta)
}
func resourceAlicloudDbfsInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbfsService := DbfsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := dbfsService.DescribeDbfsInstanceAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dbfs_instance_attachment dbfsService.DescribeDbfsInstanceAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("status", object["Status"])
	d.Set("instance_id", object["FsId"])
	if ecsListList, ok := object["EcsList"]; ok && ecsListList != nil {
		for _, ecsListListItem := range ecsListList.([]interface{}) {
			if ecsListListItemMap, ok := ecsListListItem.(map[string]interface{}); ok {
				if ok && ecsListListItemMap["EcsId"] == parts[1] {
					d.Set("ecs_id", ecsListListItemMap["EcsId"])
					break
				}
			}
		}
	}
	return nil
}

func resourceAlicloudDbfsInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbfsService := DbfsService{client}
	action := "DetachDbfs"
	var response map[string]interface{}
	conn, err := client.NewDbfsClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"FsId":          parts[0],
		"ECSInstanceId": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"unattached"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, dbfsService.DbfsInstanceStateRefreshFunc(parts[0], []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

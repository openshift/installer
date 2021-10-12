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

func resourceAlicloudEcsDiskAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsDiskAttachmentCreate,
		Read:   resourceAlicloudEcsDiskAttachmentRead,
		Update: resourceAlicloudEcsDiskAttachmentUpdate,
		Delete: resourceAlicloudEcsDiskAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bootable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"delete_with_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"device": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_pair_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEcsDiskAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachDisk"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("bootable"); ok {
		request["Bootable"] = v
	}

	if v, ok := d.GetOkExists("delete_with_instance"); ok {
		request["DeleteWithInstance"] = v
	}

	request["DiskId"] = d.Get("disk_id")
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	}

	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}

	ecsService := EcsService{client}
	oldDisk, err := ecsService.DescribeEcsDisk(d.Get("disk_id").(string))
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, DiskInvalidOperation) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_disk_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DiskId"], ":", request["InstanceId"]))
	parts, err := ParseResourceId(d.Id(), 2)
	stateConf := BuildStateConf([]string{}, []string{"In_use"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecsService.EcsDiskStateRefreshFunc(parts[0], []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	newDisk, err := ecsService.DescribeEcsDisk(d.Get("disk_id").(string))
	if err != nil {
		return WrapError(err)
	}
	if newDisk["DeleteAutoSnapshot"].(bool) != oldDisk["DeleteAutoSnapshot"].(bool) {
		action := "ModifyDiskAttribute"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}

		request := map[string]interface{}{
			"DiskId": d.Get("disk_id"),
		}
		request["RegionId"] = client.RegionId
		request["DeleteAutoSnapshot"] = oldDisk["DeleteAutoSnapshot"]
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudEcsDiskAttachmentRead(d, meta)
}
func resourceAlicloudEcsDiskAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	disk, err := ecsService.DescribeEcsDiskAttachment(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", disk["InstanceId"])
	d.Set("disk_id", disk["DiskId"])
	d.Set("device", disk["Device"])

	return nil
}
func resourceAlicloudEcsDiskAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudEcsDiskAttachmentRead(d, meta)
}
func resourceAlicloudEcsDiskAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DetachDisk"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DiskId":     parts[0],
		"InstanceId": parts[1],
	}

	if v, ok := d.GetOkExists("delete_with_instance"); ok {
		request["DeleteWithInstance"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, DiskInvalidOperation) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDisk.AlreadyDetached", "InvalidDiskId.NotFound", "InvalidDiskId.OperationNotSupported", "InvalidInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	ecsService := EcsService{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, ecsService.EcsDiskStateRefreshFunc(parts[0], []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

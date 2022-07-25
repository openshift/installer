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

func resourceAlicloudBastionhostHostGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBastionhostHostGroupCreate,
		Read:   resourceAlicloudBastionhostHostGroupRead,
		Update: resourceAlicloudBastionhostHostGroupUpdate,
		Delete: resourceAlicloudBastionhostHostGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudBastionhostHostGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHostGroup"
	request := make(map[string]interface{})
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("comment"); ok {
		request["Comment"] = v
	}
	request["HostGroupName"] = d.Get("host_group_name")
	request["InstanceId"] = d.Get("instance_id")
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_host_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", response["HostGroupId"]))

	return resourceAlicloudBastionhostHostGroupRead(d, meta)
}
func resourceAlicloudBastionhostHostGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	yundunBastionhostService := YundunBastionhostService{client}
	object, err := yundunBastionhostService.DescribeBastionhostHostGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_bastionhost_host_group yundunBastionhostService.DescribeBastionhostHostGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("host_group_id", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("comment", object["Comment"])
	d.Set("host_group_name", object["HostGroupName"])
	return nil
}
func resourceAlicloudBastionhostHostGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"HostGroupId": parts[1],
		"InstanceId":  parts[0],
	}
	if d.HasChange("comment") {
		update = true
		if v, ok := d.GetOk("comment"); ok {
			request["Comment"] = v
		}
	}
	if d.HasChange("host_group_name") {
		update = true
		request["HostGroupName"] = d.Get("host_group_name")
	}
	request["RegionId"] = client.RegionId
	if update {
		action := "ModifyHostGroup"
		conn, err := client.NewBastionhostClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	}
	return resourceAlicloudBastionhostHostGroupRead(d, meta)
}
func resourceAlicloudBastionhostHostGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteHostGroup"
	var response map[string]interface{}
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"HostGroupId": parts[1],
		"InstanceId":  parts[0],
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"OBJECT_NOT_FOUND"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

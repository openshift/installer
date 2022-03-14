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

func resourceAlicloudBastionhostHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBastionhostHostCreate,
		Read:   resourceAlicloudBastionhostHostRead,
		Update: resourceAlicloudBastionhostHostUpdate,
		Delete: resourceAlicloudBastionhostHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"active_address_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Private", "Public"}, false),
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_private_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_public_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Linux", "Windows"}, false),
			},
			"source": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Ecs", "Local", "Rds"}, false),
			},
			"source_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudBastionhostHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHost"
	request := make(map[string]interface{})
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	request["ActiveAddressType"] = d.Get("active_address_type")
	if v, ok := d.GetOk("comment"); ok {
		request["Comment"] = v
	}
	request["HostName"] = d.Get("host_name")
	if v, ok := d.GetOk("host_private_address"); ok {
		request["HostPrivateAddress"] = v
	}
	if v, ok := d.GetOk("host_public_address"); ok {
		request["HostPublicAddress"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("instance_region_id"); ok {
		request["InstanceRegionId"] = v
	}
	request["OSType"] = d.Get("os_type")
	request["RegionId"] = client.RegionId
	request["Source"] = d.Get("source")
	if v, ok := d.GetOk("source_instance_id"); ok {
		request["SourceInstanceId"] = v
	}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_host", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", response["HostId"]))

	return resourceAlicloudBastionhostHostRead(d, meta)
}
func resourceAlicloudBastionhostHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	yundunBastionhostService := YundunBastionhostService{client}
	object, err := yundunBastionhostService.DescribeBastionhostHost(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_bastionhost_host yundunBastionhostService.DescribeBastionhostHost Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("host_id", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("active_address_type", object["ActiveAddressType"])
	d.Set("comment", object["Comment"])
	d.Set("host_name", object["HostName"])
	d.Set("host_private_address", object["HostPrivateAddress"])
	d.Set("host_public_address", object["HostPublicAddress"])
	d.Set("os_type", object["OSType"])
	d.Set("source", object["Source"])
	d.Set("source_instance_id", object["SourceInstanceId"])
	return nil
}
func resourceAlicloudBastionhostHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"HostIds":    convertListToJsonString([]interface{}{parts[1]}),
		"InstanceId": parts[0],
	}
	if d.HasChange("active_address_type") {
		update = true
	}
	request["ActiveAddressType"] = d.Get("active_address_type")
	request["RegionId"] = client.RegionId
	if update {
		action := "ModifyHostsActiveAddressType"
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
		d.SetPartial("active_address_type")
	}
	update = false
	modifyHostReq := map[string]interface{}{
		"HostId":     parts[1],
		"InstanceId": parts[0],
	}
	if d.HasChange("comment") {
		update = true
		if v, ok := d.GetOk("comment"); ok {
			modifyHostReq["Comment"] = v
		}
	}
	if d.HasChange("host_name") {
		update = true
		modifyHostReq["HostName"] = d.Get("host_name")
	}
	if d.HasChange("host_private_address") {
		update = true
		if v, ok := d.GetOk("host_private_address"); ok {
			modifyHostReq["HostPrivateAddress"] = v
		}
	}
	if d.HasChange("host_public_address") {
		update = true
		if v, ok := d.GetOk("host_public_address"); ok {
			modifyHostReq["HostPublicAddress"] = v
		}
	}
	if d.HasChange("os_type") {
		update = true
		modifyHostReq["OSType"] = d.Get("os_type")
	}
	modifyHostReq["RegionId"] = client.RegionId
	if update {
		action := "ModifyHost"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, modifyHostReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyHostReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("comment")
		d.SetPartial("host_name")
		d.SetPartial("host_private_address")
		d.SetPartial("host_public_address")
		d.SetPartial("os_type")
	}
	d.Partial(false)
	return resourceAlicloudBastionhostHostRead(d, meta)
}
func resourceAlicloudBastionhostHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteHost"
	var response map[string]interface{}
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"HostId":     parts[1],
		"InstanceId": parts[0],
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

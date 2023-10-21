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

func resourceAlicloudVideoSurveillanceSystemGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVideoSurveillanceSystemGroupCreate,
		Read:   resourceAlicloudVideoSurveillanceSystemGroupRead,
		Update: resourceAlicloudVideoSurveillanceSystemGroupUpdate,
		Delete: resourceAlicloudVideoSurveillanceSystemGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"in_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"gb28181", "rtmp"}, false),
			},
			"out_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"flv", "rtmp", "hls"}, false),
			},
			"play_domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"push_domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"callback": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"capture_image": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"capture_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"capture_oss_bucket": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capture_oss_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capture_video": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lazy_pull": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudVideoSurveillanceSystemGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateGroup"
	request := make(map[string]interface{})
	conn, err := client.NewVsClient()
	if err != nil {
		return WrapError(err)
	}
	request["Region"] = client.RegionId
	request["InProtocol"] = d.Get("in_protocol")
	request["OutProtocol"] = d.Get("out_protocol")
	request["PlayDomain"] = d.Get("play_domain")
	request["PushDomain"] = d.Get("push_domain")
	request["Name"] = d.Get("group_name")

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_video_surveillance_system_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAlicloudVideoSurveillanceSystemGroupUpdate(d, meta)
}
func resourceAlicloudVideoSurveillanceSystemGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vsService := VsService{client}
	object, err := vsService.DescribeVideoSurveillanceSystemGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_video_surveillance_system_group vsService.DescribeVideoSurveillanceSystemGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("callback", object["Callback"])
	d.Set("description", object["Description"])
	d.Set("enabled", object["Enabled"])
	d.Set("group_name", object["Name"])
	d.Set("lazy_pull", object["LazyPull"])
	d.Set("in_protocol", object["InProtocol"])
	d.Set("out_protocol", object["OutProtocol"])
	d.Set("play_domain", object["PlayDomain"])
	d.Set("push_domain", object["PushDomain"])
	d.Set("capture_oss_bucket", object["CaptureOssBucket"])
	d.Set("capture_oss_path", object["CaptureOssPath"])
	d.Set("capture_video", object["CaptureVideo"])
	d.Set("capture_image", object["CaptureImage"])
	d.Set("capture_interval", object["CaptureInterval"])
	return nil
}
func resourceAlicloudVideoSurveillanceSystemGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"Id": d.Id(),
	}
	if d.HasChange("callback") {
		update = true
		if v, ok := d.GetOk("callback"); ok {
			request["Callback"] = v
		}
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("enabled") || d.IsNewResource() {
		update = true
		if v, ok := d.GetOkExists("enabled"); ok {
			request["Enabled"] = v
		}
	}
	if d.HasChange("group_name") {
		update = true
		if v, ok := d.GetOk("group_name"); ok {
			request["Name"] = v
		}
	}
	if d.HasChange("in_protocol") {
		update = true
		if v, ok := d.GetOk("in_protocol"); ok {
			request["InProtocol"] = v
		}
	}
	if d.HasChange("out_protocol") {
		update = true
		if v, ok := d.GetOk("out_protocol"); ok {
			request["OutProtocol"] = v
		}
	}
	if update {
		action := "ModifyGroup"
		conn, err := client.NewVsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"2001"}) || NeedRetry(err) {
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
	return resourceAlicloudVideoSurveillanceSystemGroupRead(d, meta)
}
func resourceAlicloudVideoSurveillanceSystemGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	conn, err := client.NewVsClient()
	action := "DeleteGroup"
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Id": d.Id(),
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}

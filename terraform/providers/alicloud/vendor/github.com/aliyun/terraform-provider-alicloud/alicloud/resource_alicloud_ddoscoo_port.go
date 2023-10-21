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

func resourceAlicloudDdoscooPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDdoscooPortCreate,
		Read:   resourceAlicloudDdoscooPortRead,
		Update: resourceAlicloudDdoscooPortUpdate,
		Delete: resourceAlicloudDdoscooPortDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"backend_port": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"frontend_port": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"frontend_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"real_servers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAlicloudDdoscooPortCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePort"
	request := make(map[string]interface{})
	conn, err := client.NewDdoscooClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("backend_port"); ok {
		request["BackendPort"] = v
	}

	request["FrontendPort"] = d.Get("frontend_port")
	request["FrontendProtocol"] = d.Get("frontend_protocol")
	request["InstanceId"] = d.Get("instance_id")
	request["RealServers"] = d.Get("real_servers")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_port", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["FrontendPort"], ":", request["FrontendProtocol"]))

	return resourceAlicloudDdoscooPortRead(d, meta)
}
func resourceAlicloudDdoscooPortRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}
	object, err := ddoscooService.DescribeDdoscooPort(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddoscoo_port ddoscooService.DescribeDdoscooPort Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("frontend_port", parts[1])
	d.Set("frontend_protocol", parts[2])
	d.Set("instance_id", parts[0])
	d.Set("backend_port", fmt.Sprint(formatInt(object["BackendPort"])))
	d.Set("real_servers", object["RealServers"])
	return nil
}
func resourceAlicloudDdoscooPortUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"FrontendPort":     parts[1],
		"FrontendProtocol": parts[2],
		"InstanceId":       parts[0],
	}
	request["BackendPort"] = d.Get("backend_port")
	if d.HasChange("real_servers") {
		update = true
	}
	request["RealServers"] = d.Get("real_servers")
	if update {
		action := "ModifyPort"
		conn, err := client.NewDdoscooClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudDdoscooPortRead(d, meta)
}
func resourceAlicloudDdoscooPortDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	action := "DeletePort"
	var response map[string]interface{}
	conn, err := client.NewDdoscooClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"FrontendPort":     parts[1],
		"FrontendProtocol": parts[2],
		"InstanceId":       parts[0],
	}
	if v, ok := d.GetOk("backend_port"); ok {
		request["BackendPort"] = v
	}
	request["RealServers"] = d.Get("real_servers")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}

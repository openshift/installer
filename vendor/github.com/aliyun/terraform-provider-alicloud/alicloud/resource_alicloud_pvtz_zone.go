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

func resourceAlicloudPvtzZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPvtzZoneCreate,
		Read:   resourceAlicloudPvtzZoneRead,
		Update: resourceAlicloudPvtzZoneUpdate,
		Delete: resourceAlicloudPvtzZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"is_ptr": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"en", "jp", "zh"}, false),
			},
			"proxy_pattern": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"RECORD", "ZONE"}, false),
				Default:      "ZONE",
			},
			"record_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_client_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'name' has been deprecated from version 1.107.0. Use 'zone_name' instead.",
				ConflictsWith: []string{"zone_name"},
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "Field 'creation_time' has been removed from provider version 1.107.0",
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "Field 'update_time' has been removed from provider version 1.107.0",
			},
		},
	}
}

func resourceAlicloudPvtzZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddZone"
	request := make(map[string]interface{})
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	if v, ok := d.GetOk("proxy_pattern"); ok {
		request["ProxyPattern"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("user_client_ip"); ok {
		request["UserClientIp"] = v
	}

	if v, ok := d.GetOk("zone_name"); ok {
		request["ZoneName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["ZoneName"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable", "System.Busy", "Throttling.User"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pvtz_zone", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ZoneId"]))

	return resourceAlicloudPvtzZoneUpdate(d, meta)
}
func resourceAlicloudPvtzZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	object, err := pvtzService.DescribePvtzZone(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_private_zone_zone pvtzService.DescribePvtzZone Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("is_ptr", object["IsPtr"])
	d.Set("proxy_pattern", object["ProxyPattern"])
	d.Set("record_count", formatInt(object["RecordCount"]))
	d.Set("remark", object["Remark"])
	d.Set("zone_name", object["ZoneName"])
	d.Set("name", object["ZoneName"])
	return nil
}
func resourceAlicloudPvtzZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"ZoneId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("proxy_pattern") {
		update = true
	}
	request["ProxyPattern"] = d.Get("proxy_pattern")
	if update {
		if _, ok := d.GetOk("lang"); ok {
			request["Lang"] = d.Get("lang")
		}
		if _, ok := d.GetOk("user_client_ip"); ok {
			request["UserClientIp"] = d.Get("user_client_ip")
		}
		action := "SetProxyPattern"
		conn, err := client.NewPvtzClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable", "System.Busy", "Throttling.User"}) || NeedRetry(err) {
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
		d.SetPartial("proxy_pattern")
	}
	update = false
	updateZoneRemarkReq := map[string]interface{}{
		"ZoneId": d.Id(),
	}
	if d.HasChange("remark") {
		update = true
		updateZoneRemarkReq["Remark"] = d.Get("remark")
	}
	if update {
		if _, ok := d.GetOk("lang"); ok {
			updateZoneRemarkReq["Lang"] = d.Get("lang")
		}
		if _, ok := d.GetOk("user_client_ip"); ok {
			updateZoneRemarkReq["UserClientIp"] = d.Get("user_client_ip")
		}
		action := "UpdateZoneRemark"
		conn, err := client.NewPvtzClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, updateZoneRemarkReq, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable", "System.Busy", "Throttling.User"}) || NeedRetry(err) {
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
		d.SetPartial("remark")
	}
	d.Partial(false)
	return resourceAlicloudPvtzZoneRead(d, meta)
}
func resourceAlicloudPvtzZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteZone"
	var response map[string]interface{}
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ZoneId": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("user_client_ip"); ok {
		request["UserClientIp"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"System.Busy", "Throttling.User", "Zone.VpcExists"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Zone.Invalid.Id", "Zone.Invalid.UserId", "Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

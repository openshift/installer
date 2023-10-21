package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudPvtzZoneRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPvtzZoneRecordCreate,
		Read:   resourceAlicloudPvtzZoneRecordRead,
		Update: resourceAlicloudPvtzZoneRecordUpdate,
		Delete: resourceAlicloudPvtzZoneRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 99),
				Default:      1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) != "MX"
				},
			},
			"record_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rr": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"resource_record"},
			},
			"resource_record": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'resource_record' has been deprecated from version 1.109.0. Use 'rr' instead.",
				ConflictsWith: []string{"rr"},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DISABLE", "ENABLE"}, false),
				Default:      "ENABLE",
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "CNAME", "MX", "PTR", "SRV", "TXT"}, false),
			},
			"user_client_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudPvtzZoneRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddZoneRecord"
	request := make(map[string]interface{})
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	if v, ok := d.GetOk("priority"); ok {
		request["Priority"] = v
	}

	if v, ok := d.GetOk("rr"); ok {
		request["Rr"] = v
	} else if v, ok := d.GetOk("resource_record"); ok {
		request["Rr"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "resource_record" or "rr" must be set one!`))
	}

	if v, ok := d.GetOk("ttl"); ok {
		request["Ttl"] = v
	}

	request["Type"] = d.Get("type")
	if v, ok := d.GetOk("user_client_ip"); ok {
		request["UserClientIp"] = v
	}

	request["Value"] = d.Get("value")
	request["ZoneId"] = d.Get("zone_id")
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
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pvtz_zone_record", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RecordId"], ":", request["ZoneId"]))

	return resourceAlicloudPvtzZoneRecordUpdate(d, meta)
}
func resourceAlicloudPvtzZoneRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	object, err := pvtzService.DescribePvtzZoneRecord(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_private_zone_zone_record pvtzService.DescribePvtzZoneRecord Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if strings.Split(d.Id(), ":")[0] != fmt.Sprint(object["RecordId"]) {
		d.SetId(fmt.Sprintf("%v:%v", fmt.Sprint(object["RecordId"]), fmt.Sprint(object["ZoneId"])))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("record_id", fmt.Sprint(parts[0]))
	d.Set("zone_id", parts[1])
	d.Set("priority", formatInt(object["Priority"]))
	d.Set("remark", object["Remark"])
	d.Set("rr", object["Rr"])
	d.Set("resource_record", object["Rr"])
	d.Set("status", object["Status"])
	d.Set("ttl", formatInt(object["Ttl"]))
	d.Set("type", object["Type"])
	d.Set("value", object["Value"])
	return nil
}
func resourceAlicloudPvtzZoneRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	pvtzService := PvtzService{client}
	object, err := pvtzService.DescribePvtzZoneRecord(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if strings.Split(d.Id(), ":")[0] != fmt.Sprint(object["RecordId"]) {
		d.SetId(fmt.Sprintf("%v:%v", fmt.Sprint(object["RecordId"]), fmt.Sprint(object["ZoneId"])))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"RecordId": parts[0],
	}
	if d.HasChange("remark") {
		update = true
		request["Remark"] = d.Get("remark")
	}
	if update {
		if _, ok := d.GetOk("lang"); ok {
			request["Lang"] = d.Get("lang")
		}
		action := "UpdateRecordRemark"
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
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("remark")
	}
	update = false
	setZoneRecordStatusReq := map[string]interface{}{
		"RecordId": parts[0],
	}
	if d.HasChange("status") {
		update = true
	}
	setZoneRecordStatusReq["Status"] = d.Get("status")
	if update {
		if _, ok := d.GetOk("lang"); ok {
			setZoneRecordStatusReq["Lang"] = d.Get("lang")
		}
		if _, ok := d.GetOk("user_client_ip"); ok {
			setZoneRecordStatusReq["UserClientIp"] = d.Get("user_client_ip")
		}
		action := "SetZoneRecordStatus"
		conn, err := client.NewPvtzClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, setZoneRecordStatusReq, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable", "System.Busy", "Throttling.User"}) || NeedRetry(err) {
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
		d.SetPartial("status")
	}
	update = false
	updateZoneRecordReq := map[string]interface{}{
		"RecordId": parts[0],
	}
	if !d.IsNewResource() && d.HasChange("rr") {
		update = true
		updateZoneRecordReq["Rr"] = d.Get("rr")
	}
	if !d.IsNewResource() && d.HasChange("resource_record") {
		update = true
		updateZoneRecordReq["Rr"] = d.Get("resource_record")
	}
	if updateZoneRecordReq["Rr"] == nil {
		updateZoneRecordReq["Rr"] = d.Get("rr")
	}
	if !d.IsNewResource() && d.HasChange("type") {
		update = true
	}
	updateZoneRecordReq["Type"] = d.Get("type")
	if !d.IsNewResource() && d.HasChange("value") {
		update = true
	}
	updateZoneRecordReq["Value"] = d.Get("value")
	if !d.IsNewResource() && d.HasChange("priority") {
		update = true
		updateZoneRecordReq["Priority"] = d.Get("priority")
	}
	if !d.IsNewResource() && d.HasChange("ttl") {
		update = true
		updateZoneRecordReq["Ttl"] = d.Get("ttl")
	}
	if update {
		if _, ok := d.GetOk("lang"); ok {
			updateZoneRecordReq["Lang"] = d.Get("lang")
		}
		if _, ok := d.GetOk("user_client_ip"); ok {
			updateZoneRecordReq["UserClientIp"] = d.Get("user_client_ip")
		}
		action := "UpdateZoneRecord"
		conn, err := client.NewPvtzClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, updateZoneRecordReq, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable", "System.Busy", "Throttling.User"}) || NeedRetry(err) {
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
		d.SetPartial("resource_record")
		d.SetPartial("rr")
		d.SetPartial("type")
		d.SetPartial("value")
		d.SetPartial("priority")
		d.SetPartial("ttl")
	}
	d.Partial(false)
	return resourceAlicloudPvtzZoneRecordRead(d, meta)
}
func resourceAlicloudPvtzZoneRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	object, err := pvtzService.DescribePvtzZoneRecord(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if strings.Split(d.Id(), ":")[0] != fmt.Sprint(object["RecordId"]) {
		d.SetId(fmt.Sprintf("%v:%v", fmt.Sprint(object["RecordId"]), fmt.Sprint(object["ZoneId"])))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteZoneRecord"
	var response map[string]interface{}
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RecordId": parts[0],
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
			if IsExpectedErrors(err, []string{"System.Busy", "Throttling.User"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"Zone.Invalid.Id", "Zone.Invalid.UserId", "Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

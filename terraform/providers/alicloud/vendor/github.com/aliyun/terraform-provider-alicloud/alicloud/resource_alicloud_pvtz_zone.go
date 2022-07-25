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
			"user_info": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"sync_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ON", "OFF"}, false),
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

	if v, ok := object["SyncHostTask"]; ok && v != nil {
		syncObject := v.(map[string]interface{})
		syncArray := make([]map[string]interface{}, 0)
		for _, raw := range syncObject["EcsRegions"].(map[string]interface{})["EcsRegion"].([]interface{}) {
			obj := raw.(map[string]interface{})
			user := make(map[string]interface{}, 0)
			user["user_id"] = obj["UserId"]
			regions := make([]interface{}, 0)
			for _, region := range obj["RegionIds"].(map[string]interface{})["RegionId"].([]interface{}) {
				regions = append(regions, region)
			}
			user["region_ids"] = regions
			syncArray = append(syncArray, user)
		}
		d.Set("user_info", syncArray)
		d.Set("sync_status", syncObject["Status"])
	}
	return nil
}
func resourceAlicloudPvtzZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
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

	update = false
	UpdateSyncEcsHostTaskReq := map[string]interface{}{
		"ZoneId": d.Id(),
	}
	if d.HasChange("sync_status") || d.HasChange("user_info") {
		update = true
	}
	if v, ok := d.GetOk("sync_status"); ok {
		UpdateSyncEcsHostTaskReq["Status"] = v
	}
	if v, ok := d.GetOk("user_info"); ok {
		for _, raw := range v.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			for index, val := range obj["region_ids"].(*schema.Set).List() {
				UpdateSyncEcsHostTaskReq[fmt.Sprintf("Region.%d.UserId", index+1)] = obj["user_id"]
				UpdateSyncEcsHostTaskReq[fmt.Sprintf("Region.%d.RegionId", index+1)] = val
			}
		}
	}
	if update {
		if _, ok := d.GetOk("lang"); ok {
			updateZoneRemarkReq["Lang"] = d.Get("lang")
		}
		if _, ok := d.GetOk("user_info"); !ok {
			object, err := pvtzService.DescribePvtzZone(d.Id())
			if err != nil {
				return WrapError(err)
			}
			if v, ok := object["SyncHostTask"]; ok && v != nil {
				syncObject := v.(map[string]interface{})
				for _, raw := range syncObject["EcsRegions"].(map[string]interface{})["EcsRegion"].([]interface{}) {
					obj := raw.(map[string]interface{})
					for index, region := range obj["RegionIds"].(map[string]interface{})["RegionId"].([]interface{}) {
						UpdateSyncEcsHostTaskReq[fmt.Sprintf("Region.%d.UserId", index+1)] = obj["user_id"]
						UpdateSyncEcsHostTaskReq[fmt.Sprintf("Region.%d.RegionId", index+1)] = region
					}
				}
			}
		}
		action := "UpdateSyncEcsHostTask"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, UpdateSyncEcsHostTaskReq, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"MissingRegion"}) {
					log.Printf("[DEBUG] Resource alicloud_private_zone_zone UpdateSyncEcsHostTask Missed Region!!! %s", err)
					return nil
				}
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
		d.SetPartial("sync_status")
		d.SetPartial("user_info")
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

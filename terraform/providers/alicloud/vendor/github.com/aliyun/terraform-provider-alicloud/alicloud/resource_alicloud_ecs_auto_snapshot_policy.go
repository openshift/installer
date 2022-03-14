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

func resourceAlicloudEcsAutoSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsAutoSnapshotPolicyCreate,
		Read:   resourceAlicloudEcsAutoSnapshotPolicyRead,
		Update: resourceAlicloudEcsAutoSnapshotPolicyUpdate,
		Delete: resourceAlicloudEcsAutoSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"copied_snapshots_retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"enable_cross_region_copy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repeat_weekdays": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"target_copy_regions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"time_points": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAlicloudEcsAutoSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "CreateAutoSnapshotPolicy"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("copied_snapshots_retention_days"); ok {
		request["CopiedSnapshotsRetentionDays"] = v
	}

	if v, ok := d.GetOkExists("enable_cross_region_copy"); ok {
		request["EnableCrossRegionCopy"] = v
	}

	if v, ok := d.GetOk("name"); ok {
		request["autoSnapshotPolicyName"] = v
	}

	request["regionId"] = client.RegionId
	request["repeatWeekdays"] = convertListToJsonString(d.Get("repeat_weekdays").(*schema.Set).List())
	request["retentionDays"] = d.Get("retention_days")
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if v, ok := d.GetOk("target_copy_regions"); ok {
		request["TargetCopyRegions"] = convertListToJsonString(v.(*schema.Set).List())
	}

	request["timePoints"] = convertListToJsonString(d.Get("time_points").(*schema.Set).List())
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_auto_snapshot_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AutoSnapshotPolicyId"]))
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecsService.EcsAutoSnapshotPolicyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEcsAutoSnapshotPolicyRead(d, meta)
}
func resourceAlicloudEcsAutoSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsAutoSnapshotPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_auto_snapshot_policy ecsService.DescribeEcsAutoSnapshotPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("copied_snapshots_retention_days", formatInt(object["CopiedSnapshotsRetentionDays"]))
	d.Set("enable_cross_region_copy", object["EnableCrossRegionCopy"])
	d.Set("name", object["AutoSnapshotPolicyName"])
	if object["RepeatWeekdays"] != nil {
		if repeatWeekdays, err := convertJsonStringToList(object["RepeatWeekdays"].(string)); err != nil {
			return WrapError(err)
		} else {
			d.Set("repeat_weekdays", repeatWeekdays)
		}
	}
	d.Set("retention_days", formatInt(object["RetentionDays"]))
	d.Set("status", object["Status"])
	d.Set("tags", tagsToMap(object["Tags"].(map[string]interface{})["Tag"]))
	if object["TargetCopyRegions"] != nil {
		if targetCopyRegions, err := convertJsonStringToList(object["TargetCopyRegions"].(string)); err != nil {
			return WrapError(err)
		} else {
			d.Set("target_copy_regions", targetCopyRegions)
		}
	}
	if object["TimePoints"] != nil {
		if timePoints, err := convertJsonStringToList(object["TimePoints"].(string)); err != nil {
			return WrapError(err)
		} else {
			d.Set("time_points", timePoints)
		}
	}
	return nil
}
func resourceAlicloudEcsAutoSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "snapshotpolicy"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := map[string]interface{}{
		"autoSnapshotPolicyId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("copied_snapshots_retention_days") {
		update = true
		request["CopiedSnapshotsRetentionDays"] = d.Get("copied_snapshots_retention_days")
	}
	if !d.IsNewResource() && d.HasChange("enable_cross_region_copy") {
		update = true
		request["EnableCrossRegionCopy"] = d.Get("enable_cross_region_copy")
	}
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["autoSnapshotPolicyName"] = d.Get("name")
	}
	if !d.IsNewResource() && d.HasChange("repeat_weekdays") {
		update = true
		request["repeatWeekdays"] = convertListToJsonString(d.Get("repeat_weekdays").(*schema.Set).List())
	}
	if !d.IsNewResource() && d.HasChange("retention_days") {
		update = true
		request["retentionDays"] = d.Get("retention_days")
	}
	if !d.IsNewResource() && d.HasChange("target_copy_regions") {
		update = true
		request["TargetCopyRegions"] = convertListToJsonString(d.Get("target_copy_regions").(*schema.Set).List())
	}
	if !d.IsNewResource() && d.HasChange("time_points") {
		update = true
		request["timePoints"] = convertListToJsonString(d.Get("time_points").(*schema.Set).List())
	}
	if update {
		action := "ModifyAutoSnapshotPolicyEx"
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
		d.SetPartial("copied_snapshots_retention_days")
		d.SetPartial("enable_cross_region_copy")
		d.SetPartial("name")
		d.SetPartial("repeat_weekdays")
		d.SetPartial("retention_days")
		d.SetPartial("target_copy_regions")
		d.SetPartial("time_points")
	}
	d.Partial(false)
	return resourceAlicloudEcsAutoSnapshotPolicyRead(d, meta)
}
func resourceAlicloudEcsAutoSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAutoSnapshotPolicy"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"autoSnapshotPolicyId": d.Id(),
	}

	request["regionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError", "OperationConflict", "ServiceUnavailable", "SnapshotCreatedDisk", "SnapshotCreatedImage"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterInvalid"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

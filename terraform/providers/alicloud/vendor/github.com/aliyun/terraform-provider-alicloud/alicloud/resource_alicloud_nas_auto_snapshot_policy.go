package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudNasAutoSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasAutoSnapshotPolicyCreate,
		Read:   resourceAlicloudNasAutoSnapshotPolicyRead,
		Update: resourceAlicloudNasAutoSnapshotPolicyUpdate,
		Delete: resourceAlicloudNasAutoSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_snapshot_policy_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with 'http://', 'https://'."), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9:_-]{1,127}$`), `The name must start with a letter. It must be 2 to 128 characters in length. It can contain digits, colons (:), underscores (_), and hyphens (-).`)),
			},
			"repeat_weekdays": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"retention_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.Any(validation.IntInSlice([]int{-1}), validation.IntBetween(1, 65536)),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceAlicloudNasAutoSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAutoSnapshotPolicy"
	request := make(map[string]interface{})
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("auto_snapshot_policy_name"); ok {
		request["AutoSnapshotPolicyName"] = v
	}
	request["FileSystemType"] = "extreme"
	request["RepeatWeekdays"] = convertListToCommaSeparate(d.Get("repeat_weekdays").(*schema.Set).List())
	if v, ok := d.GetOk("retention_days"); ok {
		request["RetentionDays"] = v
	}
	request["TimePoints"] = convertListToCommaSeparate(d.Get("time_points").(*schema.Set).List())
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_auto_snapshot_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AutoSnapshotPolicyId"]))
	nasService := NasService{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nasService.NasAutoSnapshotPolicyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudNasAutoSnapshotPolicyRead(d, meta)
}
func resourceAlicloudNasAutoSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	object, err := nasService.DescribeNasAutoSnapshotPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_auto_snapshot_policy nasService.DescribeNasAutoSnapshotPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("auto_snapshot_policy_name", object["AutoSnapshotPolicyName"])
	d.Set("repeat_weekdays", strings.Split(object["RepeatWeekdays"].(string), ","))
	d.Set("retention_days", formatInt(object["RetentionDays"]))
	d.Set("status", object["Status"])
	d.Set("time_points", strings.Split(object["TimePoints"].(string), ","))
	return nil
}
func resourceAlicloudNasAutoSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"AutoSnapshotPolicyId": d.Id(),
	}
	if d.HasChange("auto_snapshot_policy_name") {
		update = true
		if v, ok := d.GetOk("auto_snapshot_policy_name"); ok {
			request["AutoSnapshotPolicyName"] = v
		}
	}
	if d.HasChange("repeat_weekdays") {
		update = true
		request["RepeatWeekdays"] = convertListToCommaSeparate(d.Get("repeat_weekdays").(*schema.Set).List())
	}
	if d.HasChange("retention_days") {
		update = true
		if v, ok := d.GetOk("retention_days"); ok {
			request["RetentionDays"] = v
		}
	}
	if d.HasChange("time_points") {
		update = true
		request["TimePoints"] = convertListToCommaSeparate(d.Get("time_points").(*schema.Set).List())
	}
	if update {
		action := "ModifyAutoSnapshotPolicy"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nasService.NasAutoSnapshotPolicyStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudNasAutoSnapshotPolicyRead(d, meta)
}
func resourceAlicloudNasAutoSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	action := "DeleteAutoSnapshotPolicy"
	var response map[string]interface{}
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AutoSnapshotPolicyId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidLifecyclePolicy.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nasService.NasAutoSnapshotPolicyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

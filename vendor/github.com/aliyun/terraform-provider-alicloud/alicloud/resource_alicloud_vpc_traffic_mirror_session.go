package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudVpcTrafficMirrorSession() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcTrafficMirrorSessionCreate,
		Read:   resourceAlicloudVpcTrafficMirrorSessionRead,
		Update: resourceAlicloudVpcTrafficMirrorSessionUpdate,
		Delete: resourceAlicloudVpcTrafficMirrorSessionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 32766),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"traffic_mirror_filter_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_mirror_session_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"traffic_mirror_session_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-zA-Z\u4E00-\u9FA5][\u4E00-\u9FA5A-Za-z0-9_-]{2,128}$"), "The name must be `2` to `128` characters in length and can contain digits, underscores (_), and hyphens (-). It must start with a letter."),
			},
			"traffic_mirror_source_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"traffic_mirror_target_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_mirror_target_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"NetworkInterface", "SLB"}, false),
			},
			"virtual_network_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 16777215),
			},
		},
	}
}

func resourceAlicloudVpcTrafficMirrorSessionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateTrafficMirrorSession"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("enabled"); ok {
		request["Enabled"] = v
	}
	request["Priority"] = d.Get("priority")
	request["RegionId"] = client.RegionId
	request["TrafficMirrorFilterId"] = d.Get("traffic_mirror_filter_id")
	if v, ok := d.GetOk("traffic_mirror_session_description"); ok {
		request["TrafficMirrorSessionDescription"] = v
	}
	if v, ok := d.GetOk("traffic_mirror_session_name"); ok {
		request["TrafficMirrorSessionName"] = v
	}
	request["TrafficMirrorSourceIds"] = d.Get("traffic_mirror_source_ids")
	request["TrafficMirrorTargetId"] = d.Get("traffic_mirror_target_id")
	request["TrafficMirrorTargetType"] = d.Get("traffic_mirror_target_type")
	if v, ok := d.GetOk("virtual_network_id"); ok {
		request["VirtualNetworkId"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateTrafficMirrorSession")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_traffic_mirror_session", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TrafficMirrorSessionId"]))
	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcTrafficMirrorSessionStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcTrafficMirrorSessionRead(d, meta)
}
func resourceAlicloudVpcTrafficMirrorSessionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcTrafficMirrorSession(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_traffic_mirror_session vpcService.DescribeVpcTrafficMirrorSession Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("enabled", object["Enabled"])
	d.Set("priority", formatInt(object["Priority"]))
	d.Set("status", object["TrafficMirrorSessionStatus"])
	d.Set("traffic_mirror_filter_id", object["TrafficMirrorFilterId"])
	d.Set("traffic_mirror_session_description", object["TrafficMirrorSessionDescription"])
	d.Set("traffic_mirror_session_name", object["TrafficMirrorSessionName"])
	d.Set("traffic_mirror_source_ids", object["TrafficMirrorSourceIds"])
	d.Set("traffic_mirror_target_id", object["TrafficMirrorTargetId"])
	d.Set("traffic_mirror_target_type", object["TrafficMirrorTargetType"])
	d.Set("virtual_network_id", formatInt(object["VirtualNetworkId"]))
	return nil
}
func resourceAlicloudVpcTrafficMirrorSessionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"TrafficMirrorSessionId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("enabled") {
		update = true
		if v, ok := d.GetOkExists("enabled"); ok {
			request["Enabled"] = v
		}
	}
	if d.HasChange("priority") {
		update = true
		request["Priority"] = d.Get("priority")
	}
	if d.HasChange("traffic_mirror_filter_id") {
		update = true
		request["TrafficMirrorFilterId"] = d.Get("traffic_mirror_filter_id")
	}
	if d.HasChange("traffic_mirror_session_description") {
		update = true
		if v, ok := d.GetOk("traffic_mirror_session_description"); ok {
			request["TrafficMirrorSessionDescription"] = v
		}
	}
	if d.HasChange("traffic_mirror_session_name") {
		update = true
		if v, ok := d.GetOk("traffic_mirror_session_name"); ok {
			request["TrafficMirrorSessionName"] = v
		}
	}
	request["TrafficMirrorTargetId"] = d.Get("traffic_mirror_target_id")
	if d.HasChange("traffic_mirror_target_id") {
		update = true
	}

	request["TrafficMirrorTargetType"] = d.Get("traffic_mirror_target_type")

	if d.HasChange("traffic_mirror_target_type") {
		update = true
	}
	if d.HasChange("virtual_network_id") {
		update = true
		if v, ok := d.GetOk("virtual_network_id"); ok {
			request["VirtualNetworkId"] = v
		}
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
		action := "UpdateTrafficMirrorSessionAttribute"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("UpdateTrafficMirrorSessionAttribute")
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
		vpcService := VpcService{client}
		stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcTrafficMirrorSessionStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("enabled")
		d.SetPartial("priority")
		d.SetPartial("traffic_mirror_filter_id")
		d.SetPartial("traffic_mirror_session_description")
		d.SetPartial("traffic_mirror_session_name")
		d.SetPartial("traffic_mirror_target_id")
		d.SetPartial("traffic_mirror_target_type")
		d.SetPartial("virtual_network_id")
	}
	d.Partial(false)
	if d.HasChange("traffic_mirror_source_ids") {
		oldTrafficMirrorSourceIds, newTrafficMirrorSourceIds := d.GetChange("traffic_mirror_source_ids")
		removed := oldTrafficMirrorSourceIds.([]interface{})
		added := newTrafficMirrorSourceIds.([]interface{})
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		if len(removed) > 0 {
			removeSourcesFromTrafficMirrorSessionRequest := map[string]interface{}{
				"TrafficMirrorSessionId": d.Id(),
				"RegionId":               client.RegionId,
			}
			removeSourcesFromTrafficMirrorSessionRequest["TrafficMirrorSourceIds"] = removed
			if _, ok := d.GetOkExists("dry_run"); ok {
				removeSourcesFromTrafficMirrorSessionRequest["DryRun"] = d.Get("dry_run")
			}
			action := "RemoveSourcesFromTrafficMirrorSession"
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				request["ClientToken"] = buildClientToken("RemoveSourcesFromTrafficMirrorSession")
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, removeSourcesFromTrafficMirrorSessionRequest, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, removeSourcesFromTrafficMirrorSessionRequest)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, vpcService.VpcTrafficMirrorSessionStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}
		if len(added) > 0 {
			AddSourcesToTrafficMirrorSessionRequest := map[string]interface{}{
				"TrafficMirrorSessionId": d.Id(),
				"RegionId":               client.RegionId,
			}
			AddSourcesToTrafficMirrorSessionRequest["TrafficMirrorSourceIds"] = added
			if _, ok := d.GetOkExists("dry_run"); ok {
				AddSourcesToTrafficMirrorSessionRequest["DryRun"] = d.Get("dry_run")
			}
			action := "AddSourcesToTrafficMirrorSession"
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				request["ClientToken"] = buildClientToken("AddSourcesToTrafficMirrorSession")
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, AddSourcesToTrafficMirrorSessionRequest, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, AddSourcesToTrafficMirrorSessionRequest)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, vpcService.VpcTrafficMirrorSessionStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}
		d.SetPartial("TrafficMirrorSourceIds")
	}
	return resourceAlicloudVpcTrafficMirrorSessionRead(d, meta)
}
func resourceAlicloudVpcTrafficMirrorSessionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	action := "DeleteTrafficMirrorSession"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TrafficMirrorSessionId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteTrafficMirrorSession")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.TrafficMirrorSession"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.TrafficMirrorSession"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.VpcTrafficMirrorSessionStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

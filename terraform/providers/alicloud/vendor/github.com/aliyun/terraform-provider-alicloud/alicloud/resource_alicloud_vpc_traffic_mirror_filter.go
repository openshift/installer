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

func resourceAlicloudVpcTrafficMirrorFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcTrafficMirrorFilterCreate,
		Read:   resourceAlicloudVpcTrafficMirrorFilterRead,
		Update: resourceAlicloudVpcTrafficMirrorFilterUpdate,
		Delete: resourceAlicloudVpcTrafficMirrorFilterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"traffic_mirror_filter_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][\w\-]{1,255}$`), "The description must be `2` to `256` characters in length. It must start with a letter and cannot start with `http://` or `https://`"),
			},
			"traffic_mirror_filter_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][\w\-.]{1,127}$`), "The name must be `2` to `128` characters in length, and can contain digits, periods (.), underscores (_), and hyphens (-). It must start with a letter and cannot start with `http://` or `https://`."),
			},
		},
	}
}

func resourceAlicloudVpcTrafficMirrorFilterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateTrafficMirrorFilter"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("traffic_mirror_filter_description"); ok {
		request["TrafficMirrorFilterDescription"] = v
	}
	if v, ok := d.GetOk("traffic_mirror_filter_name"); ok {
		request["TrafficMirrorFilterName"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateTrafficMirrorFilter")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_traffic_mirror_filter", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TrafficMirrorFilterId"]))
	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcTrafficMirrorFilterStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcTrafficMirrorFilterRead(d, meta)
}
func resourceAlicloudVpcTrafficMirrorFilterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcTrafficMirrorFilter(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_traffic_mirror_filter vpcService.DescribeVpcTrafficMirrorFilter Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("status", object["TrafficMirrorFilterStatus"])
	d.Set("traffic_mirror_filter_description", object["TrafficMirrorFilterDescription"])
	d.Set("traffic_mirror_filter_name", object["TrafficMirrorFilterName"])
	return nil
}
func resourceAlicloudVpcTrafficMirrorFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"TrafficMirrorFilterId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("traffic_mirror_filter_description") {
		update = true
		if v, ok := d.GetOk("traffic_mirror_filter_description"); ok {
			request["TrafficMirrorFilterDescription"] = v
		}
	}
	if d.HasChange("traffic_mirror_filter_name") {
		update = true
		if v, ok := d.GetOk("traffic_mirror_filter_name"); ok {
			request["TrafficMirrorFilterName"] = v
		}
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
		action := "UpdateTrafficMirrorFilterAttribute"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("UpdateTrafficMirrorFilterAttribute")
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
	}
	return resourceAlicloudVpcTrafficMirrorFilterRead(d, meta)
}
func resourceAlicloudVpcTrafficMirrorFilterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTrafficMirrorFilter"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TrafficMirrorFilterId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteTrafficMirrorFilter")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.TrafficMirrorFilter", "IncorrectStatus.TrafficMirrorRule"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.TrafficMirrorFilter"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

package alicloud

import (
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpcTrafficMirrorService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpcTrafficMirrorServiceRead,
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
				Optional:     true,
				Default:      "Off",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAlicloudVpcTrafficMirrorServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("VpcTrafficMirrorServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	var response map[string]interface{}
	request := map[string]interface{}{}

	action := "GetTrafficMirrorServiceStatus"
	request["ClientToken"] = buildClientToken(action)
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
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

		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_traffic_mirror_service", action, AlibabaCloudSdkGoERROR)
	}

	if v, ok := response["Enabled"]; ok && v.(bool) {
		d.SetId("VpcTrafficMirrorServiceHasBeenOpened")
		d.Set("status", "Opened")
		return nil
	}

	action = "OpenTrafficMirrorService"
	request["ClientToken"] = buildClientToken(action)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
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
		if IsExpectedErrors(err, []string{"ORDER.OPEND"}) {
			d.SetId("VpcTrafficMirrorServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_traffic_mirror_service", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId("VpcTrafficMirrorServiceHasBeenOpened")
	d.Set("status", "Opened")
	return nil
}

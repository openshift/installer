package alicloud

import (
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCenTransitRouterService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenTransitRouterServiceRead,
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
func dataSourceAlicloudCenTransitRouterServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("CenTransitRouterServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}

	var response map[string]interface{}
	request := map[string]interface{}{}
	action := "CheckTransitRouterService"
	request["ClientToken"] = buildClientToken(action)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"USER_NOT_OPEN_TR_SERVICE"}) {
			action = "OpenTransitRouterService"
			request["ClientToken"] = buildClientToken(action)
			err = resource.Retry(3*time.Minute, func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
					d.SetId("CenTransitRouterServiceHasBeenOpened")
					d.Set("status", "Opened")
					return nil
				}
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_transit_router_service", action, AlibabaCloudSdkGoERROR)
			}

			d.SetId("CenTransitRouterServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_transit_router_service", action, AlibabaCloudSdkGoERROR)
	}

	if v, ok := response["Enabled"]; ok && v.(bool) {
		d.SetId("CenTransitRouterServiceHasBeenOpened")
		d.Set("status", "Opened")
		return nil
	}

	d.SetId("CenTransitRouterServiceHasNotBeenOpened")
	d.Set("status", "")
	return nil
}

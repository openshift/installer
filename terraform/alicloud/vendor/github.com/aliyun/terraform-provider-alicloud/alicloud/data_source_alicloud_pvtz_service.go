package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudPvtzService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPvtzServiceRead,

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
func dataSourceAlicloudPvtzServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("PvtzServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	var response map[string]interface{}
	action := "CreateInstance"
	request := map[string]interface{}{
		"ProductCode":       "pvtz",
		"ProductType":       "pvtzpost",
		"SubscriptionType":  "PayAsYouGo",
		"Parameter.1.Code":  "CommodityType",
		"Parameter.1.Value": "pvtz",
	}
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, nil)
	if err != nil {
		if IsExpectedErrors(err, []string{"SYSTEM.SALE_VALIDATE_NO_SPECIFIC_CODE_FAILED"}) {
			d.SetId("PvtzServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_pvtz_service", action, AlibabaCloudSdkGoERROR)
	}

	if response["Success"] != nil && !response["Success"].(bool) {
		if response["Code"] != nil && response["Code"].(string) == "SYSTEM.SALE_VALIDATE_NO_SPECIFIC_CODE_FAILED" {
			d.SetId("PvtzServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(fmt.Errorf("%s", response), DataDefaultErrorMsg, "alicloud_pvtz_service", action, AlibabaCloudSdkGoERROR)
	}

	if response["Data"] != nil {
		d.SetId(fmt.Sprintf("%v", response["Data"].(map[string]interface{})["OrderId"]))
	} else {
		log.Printf("[ERROR] When opening pvtz service, invoking %s got an nil data. Response: %s.", action, response)
		d.SetId("PvtzServiceHasBeenOpened")
	}
	d.Set("status", "Opened")

	return nil
}

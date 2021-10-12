package alicloud

import (
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCloudStorageGatewayService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudStorageGatewayServiceRead,

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
func dataSourceAlicloudCloudStorageGatewayServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("CloudStorageGatewayServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	action := "OpenSgwService"
	conn, err := meta.(*connectivity.AliyunClient).NewTeaCommonClient(connectivity.OpenCloudStorageGatewayService)
	if err != nil {
		return WrapError(err)
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, nil, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, response, nil)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"OrderOpend"}) {
			d.SetId("CloudStorageGatewayServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_storage_gateway_service", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId("CloudStorageGatewayServiceHasBeenOpened")
	d.Set("status", "Opened")

	return nil
}

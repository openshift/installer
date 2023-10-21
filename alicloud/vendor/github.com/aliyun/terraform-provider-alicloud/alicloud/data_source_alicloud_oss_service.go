package alicloud

import (
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOssService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOssServiceRead,

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
func dataSourceAlicloudOssServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("OssServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	conn, err := meta.(*connectivity.AliyunClient).NewTeaCommonClient(connectivity.OpenOssService)
	if err != nil {
		return WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer("OpenOssService"), nil, StringPointer("POST"), StringPointer("2019-04-22"), StringPointer("AK"), nil, nil, &util.RuntimeOptions{})

	addDebug("OpenOssService", response, nil)
	if err != nil {
		if IsExpectedErrors(err, []string{"SYSTEM.SALE_VALIDATE_NO_SPECIFIC_CODE_FAILEDError", "ORDER.OPEND"}) {
			d.SetId("OssServicHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_oss_service", "OpenOssService", AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v", response["OrderId"]))
	d.Set("status", "Opened")

	return nil
}

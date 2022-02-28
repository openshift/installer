package alicloud

import (
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudNasService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNasServiceRead,

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
func dataSourceAlicloudNasServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("NasServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	conn, err := meta.(*connectivity.AliyunClient).NewTeaCommonClient(connectivity.OpenNasService)
	if err != nil {
		return WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer("OpenNASService"), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, nil, &util.RuntimeOptions{})
	addDebug("OpenNASService", response, nil)
	if err != nil {
		if IsExpectedErrors(err, []string{"ORDER.OPENDError", "ORDER.OPEND"}) {
			d.SetId("NasServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nas_service", "OpenNASService", AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v", response["OrderId"]))
	d.Set("status", "Opened")

	return nil
}

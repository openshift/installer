package alicloud

import (
	"time"

	rpc "github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const maxWaitTime = 60

func dataSourceAlicloudLogService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudLogServiceRead,

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
func dataSourceAlicloudLogServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("LogServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}
	conn, err := meta.(*connectivity.AliyunClient).NewTeaCommonClient(fmt.Sprintf("%s.log.aliyuncs.com/open-api", meta.(*connectivity.AliyunClient).RegionId))
	if err != nil {
		return WrapError(err)
	}
	isNotOpened, err := waitServiceReady(conn, false)
	if err == nil {
		d.SetId("LogServiceHasBeenOpened")
		d.Set("status", "Opened")
		return nil
	}
	if isNotOpened {
		response, err := conn.DoRequest(StringPointer("OpenSlsService"), nil, StringPointer("POST"), StringPointer("2019-10-23"), StringPointer("AK"), nil, nil, &util.RuntimeOptions{})
		addDebug("OpenSlsService", response, nil)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_log_service", "OpenLogService", AlibabaCloudSdkGoERROR)
		}
		if response["Success"] != nil && !response["Success"].(bool) {
			return WrapErrorf(fmt.Errorf("%s", response), DataDefaultErrorMsg, "alicloud_log_service", "OpenLogService", AlibabaCloudSdkGoERROR)
		}
		_, err = waitServiceReady(conn, true)
		if err == nil {
			d.SetId("LogServiceHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return WrapError(err)
	}
	return WrapError(err)
}

func waitServiceReady(conn *rpc.Client, hasOpened bool) (bool, error) {
	beginTime := time.Now().Unix()
	for {
		response, err := conn.DoRequest(StringPointer("GetSlsService"), nil, StringPointer("POST"), StringPointer("2019-10-23"), StringPointer("AK"), nil, nil, &util.RuntimeOptions{})
		addDebug("GetSlsService", response, nil)
		if err != nil {
			return false, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_log_service", "GetLogService", AlibabaCloudSdkGoERROR)
		}
		if response["Success"] != nil && !response["Success"].(bool) {
			return false, WrapErrorf(fmt.Errorf("%s", response), DataDefaultErrorMsg, "alicloud_log_service", "GetLogService", AlibabaCloudSdkGoERROR)
		}
		status := response["Status"].(string)
		if "Opened" == status {
			return false, nil
		}
		if hasOpened || "Opening" == status {
			if time.Now().Unix()-beginTime >= maxWaitTime {
				return false, fmt.Errorf("wait until the maxWaitTime(60s) is still in the %s state", status)
			}
			time.Sleep(time.Second)
			continue
		}
		return true, fmt.Errorf("incorrect status: %s", status)
	}

}

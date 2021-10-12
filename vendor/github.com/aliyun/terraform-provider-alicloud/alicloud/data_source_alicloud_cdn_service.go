package alicloud

import (
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCdnService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCdnServiceRead,

		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
				Optional:     true,
				Default:      "Off",
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"PayByTraffic", "PayByBandwidth"}, false),
				Optional:     true,
				Default:      "PayByTraffic",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"opening_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"changing_charge_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"changing_affect_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAlicloudCdnServiceRead(d *schema.ResourceData, meta interface{}) error {
	opened := false
	enable := ""
	if v, ok := d.GetOk("enable"); ok {
		enable = v.(string)
	}

	conn, err := meta.(*connectivity.AliyunClient).NewTeaCommonClient(connectivity.OpenCdnService)
	if err != nil {
		return WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer("DescribeCdnService"), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, nil, &util.RuntimeOptions{})
	addDebug("DescribeCdnService", response, nil)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cdn_service", "DescribeCdnService", AlibabaCloudSdkGoERROR)
	}

	if response["OpeningTime"] != nil && response["OpeningTime"].(string) != "" {
		opened = true
	}
	if enable == "On" {
		chargeType := ""
		if v, ok := d.GetOk("internet_charge_type"); ok {
			chargeType = v.(string)
		}

		if chargeType == "" {
			return WrapError(fmt.Errorf("Field 'internet_charge_type' is required when 'enable' is 'On'."))
		}
		requestBody := map[string]interface{}{"InternetChargeType": chargeType}
		if opened && chargeType != response["ChangingChargeType"].(string) {
			resp, err := conn.DoRequest(StringPointer("ModifyCdnService"), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, requestBody, &util.RuntimeOptions{})

			addDebug("ModifyCdnService", resp, nil)
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cdn_service", "ModifyCdnService", AlibabaCloudSdkGoERROR)
			}
		}
		if !opened {
			resp, err := conn.DoRequest(StringPointer("OpenCdnService"), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, requestBody, &util.RuntimeOptions{})

			addDebug("OpenCdnService", resp, nil)
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cdn_service", "OpenCdnService", AlibabaCloudSdkGoERROR)
			}
		}
		response, err = conn.DoRequest(StringPointer("DescribeCdnService"), nil, StringPointer("POST"), StringPointer("2018-05-10"), StringPointer("AK"), nil, nil, &util.RuntimeOptions{})

		addDebug("DescribeCdnService", response, nil)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cdn_service", "DescribeCdnService", AlibabaCloudSdkGoERROR)
		}
	}

	if opened {
		d.SetId("CdnServiceHasBeenOpened")
		d.Set("status", "Opened")
	} else {
		d.SetId("CdnServiceHasNotBeenOpened")
		d.Set("status", "")
	}

	d.Set("internet_charge_type", response["InternetChargeType"])
	d.Set("opening_time", response["OpeningTime"])
	d.Set("changing_charge_type", response["ChangingChargeType"])
	d.Set("changing_affect_time", response["ChangingAffectTime"])

	return nil
}

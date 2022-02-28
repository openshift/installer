package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudQuotasQuotaApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudQuotasQuotaApplicationCreate,
		Read:   resourceAlicloudQuotasQuotaApplicationRead,
		Delete: resourceAlicloudQuotasQuotaApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"approve_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"audit_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Async", "Sync"}, false),
				Default:      "Async",
			},
			"audit_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"desire_value": {
				Type:     schema.TypeFloat,
				Required: true,
				ForceNew: true,
			},
			"dimensions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
				ForceNew: true,
			},
			"effective_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notice_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3}),
				Default:      0,
			},
			"product_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota_action_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CommonQuota", "FlowControl"}, false),
			},
			"quota_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quota_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quota_unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudQuotasQuotaApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateQuotaApplication"
	request := make(map[string]interface{})
	conn, err := client.NewQuotasClient()
	if err != nil {
		return WrapError(err)
	}
	request["SourceIp"] = client.SourceIp
	if v, ok := d.GetOk("audit_mode"); ok {
		request["AuditMode"] = v
	}

	request["DesireValue"] = d.Get("desire_value")
	if v, ok := d.GetOk("dimensions"); ok {
		dimensionsMaps := make([]map[string]interface{}, 0)
		for _, dimensions := range v.(*schema.Set).List() {
			dimensionsMap := make(map[string]interface{})
			dimensionsArg := dimensions.(map[string]interface{})
			dimensionsMap["Key"] = dimensionsArg["key"]
			dimensionsMap["Value"] = dimensionsArg["value"]
			dimensionsMaps = append(dimensionsMaps, dimensionsMap)
		}
		request["Dimensions"] = dimensionsMaps

	}

	if v, ok := d.GetOk("notice_type"); ok {
		request["NoticeType"] = v
	}

	request["ProductCode"] = d.Get("product_code")
	request["QuotaActionCode"] = d.Get("quota_action_code")
	if v, ok := d.GetOk("quota_category"); ok {
		request["QuotaCategory"] = v
	}

	request["Reason"] = d.Get("reason")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_quotas_quota_application", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ApplicationId"]))

	return resourceAlicloudQuotasQuotaApplicationRead(d, meta)
}
func resourceAlicloudQuotasQuotaApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	quotasService := QuotasService{client}
	object, err := quotasService.DescribeQuotasQuotaApplication(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_quotas_quota_application quotasService.DescribeQuotasQuotaApplication Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("approve_value", object["ApproveValue"])
	d.Set("audit_reason", object["AuditReason"])
	d.Set("desire_value", object["DesireValue"])

	dimensionList := make([]map[string]interface{}, 0)
	if dimension, ok := object["Dimension"]; ok {
		for k, v := range dimension.(map[string]interface{}) {
			dimensionMap := make(map[string]interface{})
			dimensionMap["key"] = k
			dimensionMap["value"] = v
			dimensionList = append(dimensionList, dimensionMap)
		}
	}

	if err := d.Set("dimensions", dimensionList); err != nil {
		return WrapError(err)
	}
	d.Set("effective_time", object["EffectiveTime"])
	d.Set("expire_time", object["ExpireTime"])
	d.Set("notice_type", object["NoticeType"])
	d.Set("product_code", object["ProductCode"])
	d.Set("quota_action_code", object["QuotaActionCode"])
	d.Set("quota_description", object["QuotaDescription"])
	d.Set("quota_name", object["QuotaName"])
	d.Set("quota_unit", object["QuotaUnit"])
	d.Set("reason", object["Reason"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudQuotasQuotaApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudQuotasQuotaApplication. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

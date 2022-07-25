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

func resourceAlicloudSddpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSddpInstanceCreate,
		Read:   resourceAlicloudSddpInstanceRead,
		Update: resourceAlicloudSddpInstanceUpdate,
		Delete: resourceAlicloudSddpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Subscription"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 12, 2, 24, 3, 6}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 12),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
							return false
						}
					}
					return true
				},
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"logistics": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dataphin_count": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dataphin": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"no", "yes"}, false),
			},
			"sddp_version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"version_audit", "version_company", "version_dlp"}, false),
			},
			"sdc": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ud_cbool": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"no", "yes"}, false),
			},
			"sd_cbool": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"no", "yes"}, false),
			},
			"udc": {
				Type:     schema.TypeString,
				Required: true,
			},
			"modify_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Upgrade", "Downgrade"}, false),
				Optional:     true,
			},
			"authed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"instance_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"odps_set": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"oss_bucket_set": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"oss_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rds_set": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remain_days": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudSddpInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateInstance"
	request := make(map[string]interface{})
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}

	request["ClientToken"] = buildClientToken("CreateInstance")
	request["ProductCode"] = "sddp"
	request["ProductType"] = "sddp_pre"
	request["SubscriptionType"] = d.Get("payment_type")

	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOk("renewal_duration"); ok {
		request["RenewPeriod"] = v
	} else if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "renewal_duration", "renewal_status", d.Get("renewal_status")))
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	} else if d.Get("payment_type").(string) == "Subscription" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v", "period", "payment_type", "Subscription"))
	}

	if v, ok := d.GetOk("logistics"); ok {
		request["Logistics"] = v
	}

	parameterMapList := make([]map[string]interface{}, 0)
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "sddp_version",
		"Value": d.Get("sddp_version"),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "SDCbool",
		"Value": d.Get("sd_cbool"),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "SDC",
		"Value": d.Get("sdc"),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "UDCbool",
		"Value": d.Get("ud_cbool"),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "UDC",
		"Value": d.Get("udc"),
	})
	if v, ok := d.GetOk("dataphin"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "dataphin",
			"Value": v,
		})
	} else if d.Get("sddp_version").(string) == "version_audit" || d.Get("sddp_version").(string) == "version_company" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v or %v", "dataphin", "sddp_version", "version_audit", "version_company"))
	}
	if v, ok := d.GetOk("dataphin_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "dataphin_count",
			"Value": v,
		})
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "region",
		"Value": client.RegionId,
	})
	request["Parameter"] = parameterMapList
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &runtime)
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
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sddp_instance", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["InstanceId"]))
	return resourceAlicloudSddpInstanceRead(d, meta)
}
func resourceAlicloudSddpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sddpService := SddpService{client}
	object, err := sddpService.DescribeSddpInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sddp_instance sddpService.DescribeSddpInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("payment_type", convertSddpInstancePaymentTypeToStandard(object["ChargeType"]))
	d.Set("sddp_version", object["Version"])
	d.Set("sdc", fmt.Sprint(formatInt(object["InstanceNum"])))

	d.Set("authed", object["Authed"])
	d.Set("renewal_status", object["RenewStatus"])
	d.Set("instance_num", fmt.Sprint(formatInt(object["InstanceNum"])))
	d.Set("odps_set", object["OdpsSet"])
	d.Set("oss_bucket_set", object["OssBucketSet"])
	d.Set("oss_size", fmt.Sprint(formatInt(object["OssSize"])))
	d.Set("rds_set", object["RdsSet"])
	d.Set("status", fmt.Sprint(formatInt(object["InstanceStatus"])))
	d.Set("remain_days", object["RemainDays"])
	return nil
}
func resourceAlicloudSddpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	request["SubscriptionType"] = d.Get("payment_type")
	request["ClientToken"] = buildClientToken("ModifyInstance")
	request["ProductCode"] = "sddp"
	request["ProductType"] = "sddp_pre"
	parameterMapList := make([]map[string]interface{}, 0)
	if d.HasChange("sdc") {
		update = true
	}
	if v, ok := d.GetOk("sdc"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "SDC",
			"Value": v,
		})
	}
	if d.HasChange("oss_size") {
		update = true
	}
	if v, ok := d.GetOk("oss_size"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "UDC",
			"Value": v,
		})
	}

	if d.HasChange("sd_cbool") {
		update = true
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "SDCbool",
		"Value": d.Get("sd_cbool"),
	})
	if d.HasChange("ud_cbool") {
		update = true
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "UDCbool",
		"Value": d.Get("ud_cbool"),
	})
	if d.HasChange("sddp_version") {
		update = true
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "sddp_version",
		"Value": d.Get("sddp_version"),
	})

	if d.HasChange("sdc") {
		update = true
	}
	if v, ok := d.GetOk("sdc"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "SDC",
			"Value": v,
		})
	}
	if d.HasChange("udc") {
		update = true
	}
	if v, ok := d.GetOk("udc"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "UDC",
			"Value": v,
		})
	}
	if d.HasChange("dataphin") {
		update = true
	}
	if v, ok := d.GetOk("dataphin"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "dataphin",
			"Value": v,
		})
	}
	if d.HasChange("dataphin_count") {
		update = true
	}
	if v, ok := d.GetOk("dataphin_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "dataphin_count",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList
	if update {
		if v, ok := d.GetOk("modify_type"); ok {
			request["ModifyType"] = v
		}
		action := "ModifyInstance"
		conn, err := client.NewBssopenapiClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &runtime)
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
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}
	return resourceAlicloudSddpInstanceRead(d, meta)
}
func resourceAlicloudSddpInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudSddpInstance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertSddpInstancePaymentTypeToStandard(source interface{}) interface{} {
	switch source {
	case "PREPAY":
		return "Subscription"
	}
	return source
}

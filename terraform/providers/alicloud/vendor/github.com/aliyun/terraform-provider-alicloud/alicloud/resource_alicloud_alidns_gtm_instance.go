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

func resourceAlicloudAlidnsGtmInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsGtmInstanceCreate,
		Read:   resourceAlicloudAlidnsGtmInstanceRead,
		Update: resourceAlicloudAlidnsGtmInstanceUpdate,
		Delete: resourceAlicloudAlidnsGtmInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alert_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dingtalk_notice": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"email_notice": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"notice_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"ADDR_RESUME", "ADDR_ALERT", "ADDR_POOL_GROUP_UNAVAILABLE", "ADDR_POOL_GROUP_AVAILABLE", "ACCESS_STRATEGY_POOL_GROUP_SWITCH", "MONITOR_NODE_IP_CHANGE"}, false),
						},
						"sms_notice": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"alert_group": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cname_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLIC"}, false),
			},
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Subscription"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return fmt.Sprint(d.Get("renewal_status")) == "ManualRenewal"
				},
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
			},
			"package_edition": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "ultimate"}, false),
			},
			"health_check_task_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 100000),
			},
			"sms_notification_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 100000),
			},
			"strategy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"GEO", "LATENCY"}, false),
			},
			"public_cname_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"CUSTOM", "SYSTEM_ASSIGN"}, false),
			},
			"public_rr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return fmt.Sprint(d.Get("public_cname_mode")) == "SYSTEM_ASSIGN"
				},
			},
			"public_user_domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_zone_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return fmt.Sprint(d.Get("public_cname_mode")) == "SYSTEM_ASSIGN"
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600}),
			},
		},
	}
}

func resourceAlicloudAlidnsGtmInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateInstance"
	request := make(map[string]interface{})
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	request["SubscriptionType"] = d.Get("payment_type")
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	request["ProductCode"] = "dns"
	request["ProductType"] = "dns_gtm_public_cn"
	if v, ok := d.GetOk("renew_period"); ok {
		request["RenewPeriod"] = v
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	parameters := []map[string]string{
		{
			"Code":  "PackageEdition",
			"Value": fmt.Sprint(d.Get("package_edition")),
		},
		{
			"Code":  "HealthcheckTaskCount",
			"Value": fmt.Sprint(d.Get("health_check_task_count")),
		},
		{
			"Code":  "SmsNotificationCount",
			"Value": fmt.Sprint(d.Get("sms_notification_count")),
		},
	}

	request["Parameter"] = parameters
	request["ClientToken"] = buildClientToken("CreateInstance")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = "dns_gtm_public_intl"
				conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_gtm_instance", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["InstanceId"]))
	return resourceAlicloudAlidnsGtmInstanceUpdate(d, meta)
}
func resourceAlicloudAlidnsGtmInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	object, err := alidnsService.DescribeAlidnsGtmInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_gtm_instance alidnsService.DescribeAlidnsGtmInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("payment_type", object["PaymentType"])
	d.Set("package_edition", object["VersionCode"])
	if config, ok := object["Config"].(map[string]interface{}); ok {
		d.Set("cname_type", config["CnameType"])
		d.Set("instance_name", config["InstanceName"])
		d.Set("strategy_mode", config["StrategyMode"])
		d.Set("public_cname_mode", config["PublicCnameMode"])
		d.Set("public_rr", config["PublicRr"])
		d.Set("public_user_domain_name", config["PublicUserDomainName"])
		d.Set("public_zone_name", config["PubicZoneName"])
		if v, ok := config["Ttl"]; ok {
			d.Set("ttl", formatInt(v))
		}

		v, err := convertJsonStringToList(config["AlertGroup"].(string))
		if err != nil {
			return WrapError(err)
		} else {
			d.Set("alert_group", v)
		}

		if alertConfigsList, ok := config["AlertConfig"]; ok {
			alertConfigsArg := alertConfigsList.(map[string]interface{})
			if alertConfigConfig, ok := alertConfigsArg["AlertConfig"]; ok {
				alertConfigConfigArgs := alertConfigConfig.([]interface{})
				alertConfigsMaps := make([]map[string]interface{}, 0)
				for _, alertConfigMapArgitem := range alertConfigConfigArgs {
					alertConfigMapArg := alertConfigMapArgitem.(map[string]interface{})
					alertConfigsMap := map[string]interface{}{}
					alertConfigsMap["sms_notice"] = alertConfigMapArg["SmsNotice"]
					alertConfigsMap["notice_type"] = alertConfigMapArg["NoticeType"]
					alertConfigsMap["email_notice"] = alertConfigMapArg["EmailNotice"]
					alertConfigsMap["dingtalk_notice"] = alertConfigMapArg["DingtalkNotice"]
					alertConfigsMaps = append(alertConfigsMaps, alertConfigsMap)
				}
				d.Set("alert_config", alertConfigsMaps)
			}
		}
	}
	return nil
}
func resourceAlicloudAlidnsGtmInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"ResourceId": d.Id(),
	}
	if d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["NewResourceGroupId"] = v
		}
	}
	if update {
		if v, ok := d.GetOk("lang"); ok {
			request["Lang"] = v
		}
		action := "MoveGtmResourceGroup"
		conn, err := client.NewAlidnsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}
	update = false
	switchDnsGtmInstanceStrategyModeRequest := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("strategy_mode") {
		update = true
		if v, ok := d.GetOk("strategy_mode"); ok {
			switchDnsGtmInstanceStrategyModeRequest["StrategyMode"] = v
		}
	}
	if update {
		if v, ok := d.GetOk("lang"); ok {
			request["Lang"] = v
		}
		action := "SwitchDnsGtmInstanceStrategyMode"
		conn, err := client.NewAlidnsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-09"), StringPointer("AK"), nil, switchDnsGtmInstanceStrategyModeRequest, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("strategy_mode")
	}
	update = false
	updateDnsGtmInstanceGlobalConfigReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("alert_config") {
		update = true
		alertConfigMaps := make([]map[string]interface{}, 0)
		if v, ok := d.GetOk("alert_config"); ok {
			for _, alertConfig := range v.(*schema.Set).List() {
				alertConfigArg := alertConfig.(map[string]interface{})
				alertConfigMap := map[string]interface{}{}
				alertConfigMap["SmsNotice"] = alertConfigArg["sms_notice"]
				alertConfigMap["NoticeType"] = alertConfigArg["notice_type"]
				alertConfigMap["EmailNotice"] = alertConfigArg["email_notice"]
				alertConfigMap["DingtalkNotice"] = alertConfigArg["dingtalk_notice"]
				alertConfigMaps = append(alertConfigMaps, alertConfigMap)
			}
			updateDnsGtmInstanceGlobalConfigReq["AlertConfig"] = alertConfigMaps
		}

	}
	if d.HasChange("alert_group") {
		update = true
		if v, ok := d.GetOk("alert_group"); ok {
			updateDnsGtmInstanceGlobalConfigReq["AlertGroup"] = convertListToJsonString(v.(*schema.Set).List())
		}
	}
	if d.HasChange("instance_name") {
		update = true
	}
	if v, ok := d.GetOk("instance_name"); ok {
		updateDnsGtmInstanceGlobalConfigReq["InstanceName"] = v
	}
	if d.HasChange("ttl") {
		update = true
		if v, ok := d.GetOk("ttl"); ok {
			updateDnsGtmInstanceGlobalConfigReq["Ttl"] = v
		}
	}
	if d.HasChange("public_cname_mode") {
		update = true
		if v, ok := d.GetOk("public_cname_mode"); ok {
			updateDnsGtmInstanceGlobalConfigReq["PublicCnameMode"] = v
		}
	}
	if d.HasChange("public_rr") {
		update = true
		if v, ok := d.GetOk("public_rr"); ok {
			updateDnsGtmInstanceGlobalConfigReq["PublicRr"] = v
		}
	}
	if d.HasChange("public_user_domain_name") {
		update = true
		if v, ok := d.GetOk("public_user_domain_name"); ok {
			updateDnsGtmInstanceGlobalConfigReq["PublicUserDomainName"] = v
		}
	}
	if d.HasChange("public_zone_name") {
		update = true
		if v, ok := d.GetOk("public_zone_name"); ok {
			updateDnsGtmInstanceGlobalConfigReq["PublicZoneName"] = v
		}
	}
	if d.HasChange("cname_type") {
		update = true
		if v, ok := d.GetOk("cname_type"); ok {
			updateDnsGtmInstanceGlobalConfigReq["CnameType"] = v
		}
	}
	if update {
		if v, ok := d.GetOkExists("force_update"); ok {
			updateDnsGtmInstanceGlobalConfigReq["ForceUpdate"] = v
		}
		if v, ok := d.GetOk("lang"); ok {
			updateDnsGtmInstanceGlobalConfigReq["Lang"] = v
		}
		action := "UpdateDnsGtmInstanceGlobalConfig"
		conn, err := client.NewAlidnsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-09"), StringPointer("AK"), nil, updateDnsGtmInstanceGlobalConfigReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateDnsGtmInstanceGlobalConfigReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("alert_config")
		d.SetPartial("alert_group")
		d.SetPartial("cname_type")
		d.SetPartial("instance_name")
		d.SetPartial("public_cname_mode")
		d.SetPartial("public_rr")
		d.SetPartial("public_user_domain_name")
		d.SetPartial("public_zone_name")
		d.SetPartial("ttl")
	}
	d.Partial(false)
	return resourceAlicloudAlidnsGtmInstanceRead(d, meta)
}
func resourceAlicloudAlidnsGtmInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudAlidnsGtmInstance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

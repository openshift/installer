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

func resourceAlicloudAlidnsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsInstanceCreate,
		Read:   resourceAlicloudAlidnsInstanceRead,
		Update: resourceAlicloudAlidnsInstanceUpdate,
		Delete: resourceAlicloudAlidnsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"dns_security": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"advanced", "basic", "no"}, false),
			},
			"domain_numbers": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Subscription"}, false),
				Default:      "Subscription",
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if fmt.Sprint(d.Get("renewal_status")) == "ManualRenewal" {
						return true
					}
					return false
				},
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
			},
			"version_code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"version_enterprise_advanced", "version_enterprise_basic", "version_personal"}, false),
			},
			"version_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudAlidnsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
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
	request["ProductType"] = "alidns_pre"
	if v, ok := d.GetOk("renew_period"); ok {
		request["RenewPeriod"] = v
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	request["Parameter"] = []map[string]string{
		{
			"Code":  "DNSSecurity",
			"Value": d.Get("dns_security").(string),
		},
		{
			"Code":  "DomainNumbers",
			"Value": d.Get("domain_numbers").(string),
		},
		{
			"Code":  "Version",
			"Value": d.Get("version_code").(string),
		},
	}
	request["ClientToken"] = buildClientToken("CreateInstance")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_instance", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, "alicloud_alidns_instance", action, AlibabaCloudSdkGoERROR)
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["InstanceId"]))

	return resourceAlicloudAlidnsInstanceRead(d, meta)
}
func resourceAlicloudAlidnsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	object, err := alidnsService.DescribeAlidnsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_instance alidnsService.DescribeAlidnsInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("dns_security", convertDnsSecurityResponse(object["DnsSecurity"]))
	d.Set("domain_numbers", object["BindDomainCount"])
	d.Set("version_code", object["VersionCode"])
	d.Set("version_name", object["VersionName"])

	res, err := alidnsService.QueryAvailableInstances(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("payment_type", res["SubscriptionType"])
	d.Set("renewal_status", res["RenewStatus"])
	if fmt.Sprint(res["RenewalDurationUnit"]) == "M" {
		d.Set("renew_period", formatInt(res["RenewalDuration"]))
	} else {
		d.Set("renew_period", formatInt(res["RenewalDuration"])*12)
	}

	return nil
}
func resourceAlicloudAlidnsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudAlidnsInstanceRead(d, meta)
}
func resourceAlicloudAlidnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudAlidnsInstance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
func convertDnsSecurityResponse(source interface{}) interface{} {
	switch source {
	case "DNS Anti-DDoS Advanced":
		return "advanced"
	case "DNS Anti-DDoS Basic":
		return "basic"
	case "Not Required":
		return "no"
	}
	return source
}

package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudWafCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudWafCertificateCreate,
		Read:   resourceAlicloudWafCertificateRead,
		Update: resourceAlicloudWafCertificateUpdate,
		Delete: resourceAlicloudWafCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"certificate_id"},
				ForceNew:      true,
			},
			"private_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"certificate_id"},
				ForceNew:      true,
			},
			"certificate_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"certificate", "private_key", "certificate_name"},
			},
			"certificate_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"certificate_id"},
			},
		},
	}
}

func resourceAlicloudWafCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	conn, err := client.NewWafClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	if sslId, ok := d.GetOk("certificate_id"); ok {
		action := "CreateCertificateByCertificateId"
		request["InstanceId"] = d.Get("instance_id")
		request["CertificateId"] = sslId.(string)
		if v, ok := d.GetOk("domain"); ok {
			request["Domain"] = v
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_waf_certificate", action, AlibabaCloudSdkGoERROR)
		}
		d.SetId(fmt.Sprint(request["InstanceId"], ":", request["Domain"], ":", formatInt(response["CertificateId"])))
		return resourceAlicloudWafCertificateRead(d, meta)
	}

	request = make(map[string]interface{})
	action := "CreateCertificate"
	if v, ok := d.GetOk("certificate"); ok {
		request["Certificate"] = v
	} else {
		return WrapErrorf(err, RequiredWhenMsg, "certificate", "certificate_id", "null")
	}
	if v, ok := d.GetOk("private_key"); ok {
		request["PrivateKey"] = v
	} else {
		return WrapErrorf(err, RequiredWhenMsg, "private_key", "certificate_id", "null")
	}
	if v, ok := d.GetOk("certificate_name"); ok {
		request["CertificateName"] = v
	} else {
		return WrapErrorf(err, RequiredWhenMsg, "certificate_name", "certificate_id", "null")
	}
	request["Domain"] = d.Get("domain")
	request["InstanceId"] = d.Get("instance_id")

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_waf_certificate", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["Domain"], ":", formatInt(response["CertificateId"])))

	return resourceAlicloudWafCertificateRead(d, meta)
}
func resourceAlicloudWafCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafOpenapiService := Waf_openapiService{client}
	object, err := wafOpenapiService.DescribeWafCertificate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_waf_certificate wafOpenapiService.DescribeWafCertificate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("certificate_id", object["CertificateId"])
	d.Set("domain", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("certificate_name", object["CertificateName"])
	return nil
}
func resourceAlicloudWafCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudWafCertificateRead(d, meta)
}
func resourceAlicloudWafCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudWafCertificate. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

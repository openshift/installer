package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSslCertificatesServiceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSslCertificatesServiceCertificateCreate,
		Read:   resourceAlicloudSslCertificatesServiceCertificateRead,
		Update: resourceAlicloudSslCertificatesServiceCertificateUpdate,
		Delete: resourceAlicloudSslCertificatesServiceCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cert": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validation.StringLenBetween(1, 64),
				ForceNew:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.129.0 and it will be remove in the future version. Please use the new attribute 'certificate_name' instead.",
				ConflictsWith: []string{"certificate_name"},
				ValidateFunc:  validation.StringLenBetween(1, 64),
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudSslCertificatesServiceCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateUserCertificate"
	request := make(map[string]interface{})
	conn, err := client.NewCasClient()
	if err != nil {
		return WrapError(err)
	}
	request["Cert"] = d.Get("cert")
	request["Key"] = d.Get("key")

	if v, ok := d.GetOk("certificate_name"); ok {
		request["Name"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	} else {
		return WrapError(fmt.Errorf("Field 'certificate_name' is required"))
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-07-13"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ssl_certificates_service_certificate", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["CertId"]))

	return resourceAlicloudSslCertificatesServiceCertificateRead(d, meta)
}
func resourceAlicloudSslCertificatesServiceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	casService := CasService{client}
	object, err := casService.DescribeSslCertificatesServiceCertificate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ssl_certificates_service_certificate casService.DescribeSslCertificatesServiceCertificate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("cert", object["Cert"])
	d.Set("certificate_name", object["Name"])
	d.Set("name", object["Name"])
	d.Set("key", object["Key"])
	return nil
}
func resourceAlicloudSslCertificatesServiceCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudSslCertificatesServiceCertificateRead(d, meta)
}
func resourceAlicloudSslCertificatesServiceCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteUserCertificate"
	var response map[string]interface{}
	conn, err := client.NewCasClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"CertId": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-07-13"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}

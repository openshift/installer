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

func resourceAlicloudRamSamlProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamSamlProviderCreate,
		Read:   resourceAlicloudRamSamlProviderRead,
		Update: resourceAlicloudRamSamlProviderUpdate,
		Delete: resourceAlicloudRamSamlProviderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encodedsaml_metadata_document": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return ramSAMLProviderDiffSuppressFunc(old, new)
				},
			},
			"saml_provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamSamlProviderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSAMLProvider"
	request := make(map[string]interface{})
	conn, err := client.NewImsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("encodedsaml_metadata_document"); ok {
		request["EncodedSAMLMetadataDocument"] = v
	}

	request["SAMLProviderName"] = d.Get("saml_provider_name")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_saml_provider", action, AlibabaCloudSdkGoERROR)
	}
	responseSAMLProvider := response["SAMLProvider"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseSAMLProvider["SAMLProviderName"]))

	return resourceAlicloudRamSamlProviderRead(d, meta)
}
func resourceAlicloudRamSamlProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	imsService := ImsService{client}
	object, err := imsService.DescribeRamSamlProvider(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_saml_provider imsService.DescribeRamSamlProvider Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("saml_provider_name", d.Id())
	d.Set("arn", object["Arn"])
	d.Set("description", object["Description"])
	d.Set("encodedsaml_metadata_document", object["EncodedSAMLMetadataDocument"])
	d.Set("update_date", object["UpdateDate"])
	return nil
}
func resourceAlicloudRamSamlProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewImsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"SAMLProviderName": d.Id(),
	}
	if d.HasChange("description") {
		update = true
		request["NewDescription"] = d.Get("description")
	}
	if d.HasChange("encodedsaml_metadata_document") {
		update = true
		request["NewEncodedSAMLMetadataDocument"] = d.Get("encodedsaml_metadata_document")
	}
	if update {
		action := "UpdateSAMLProvider"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudRamSamlProviderRead(d, meta)
}
func resourceAlicloudRamSamlProviderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSAMLProvider"
	var response map[string]interface{}
	conn, err := client.NewImsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"SAMLProviderName": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"EntityNotExist.SAMLProviderError"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

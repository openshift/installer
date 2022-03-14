package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOosServiceSetting() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOosServiceSettingCreate,
		Read:   resourceAlicloudOosServiceSettingRead,
		Update: resourceAlicloudOosServiceSettingUpdate,
		Delete: resourceAlicloudOosServiceSettingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"delivery_oss_bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOkExists("delivery_sls_enabled"); ok && v.(bool) {
						return false
					}
					return true
				},
			},
			"delivery_oss_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delivery_oss_key_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOkExists("delivery_sls_enabled"); ok && v.(bool) {
						return false
					}
					return true
				},
			},
			"delivery_sls_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delivery_sls_project_name": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOkExists("delivery_sls_enabled"); ok && v.(bool) {
						return false
					}
					return true
				},
			},
		},
	}
}

func resourceAlicloudOosServiceSettingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "SetServiceSettings"
	request := make(map[string]interface{})
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("delivery_oss_bucket_name"); ok {
		request["DeliveryOssBucketName"] = v
	}
	if v, ok := d.GetOkExists("delivery_oss_enabled"); ok {
		request["DeliveryOssEnabled"] = v
	}
	if v, ok := d.GetOk("delivery_oss_key_prefix"); ok {
		request["DeliveryOssKeyPrefix"] = v
	}
	if v, ok := d.GetOkExists("delivery_sls_enabled"); ok {
		request["DeliverySlsEnabled"] = v
	}
	if v, ok := d.GetOk("delivery_sls_project_name"); ok {
		request["DeliverySlsProjectName"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_service_setting", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(strings.Trim(uuid.New().String(), "-")))
	return resourceAlicloudOosServiceSettingRead(d, meta)
}
func resourceAlicloudOosServiceSettingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}
	object, err := oosService.DescribeOosServiceSetting(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oos_service_setting oosService.DescribeOosServiceSetting Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("delivery_oss_bucket_name", object["DeliveryOssBucketName"])
	d.Set("delivery_oss_enabled", object["DeliveryOssEnabled"])
	d.Set("delivery_oss_key_prefix", object["DeliveryOssKeyPrefix"])
	d.Set("delivery_sls_enabled", object["DeliverySlsEnabled"])
	d.Set("delivery_sls_project_name", object["DeliverySlsProjectName"])
	d.Set("rdc_enterprise_id", object["RdcEnterpriseId"])
	return nil
}
func resourceAlicloudOosServiceSettingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{}
	if d.HasChange("delivery_oss_enabled") {
		update = true
		if v, ok := d.GetOkExists("delivery_oss_enabled"); ok {
			request["DeliveryOssEnabled"] = v
		}
	}
	if d.HasChange("delivery_oss_bucket_name") {
		update = true
		if v, ok := d.GetOk("delivery_oss_bucket_name"); ok {
			request["DeliveryOssBucketName"] = v
		}
	}
	if d.HasChange("delivery_oss_key_prefix") {
		update = true
		if v, ok := d.GetOk("delivery_oss_key_prefix"); ok {
			request["DeliveryOssKeyPrefix"] = v
		}
	}
	if d.HasChange("delivery_sls_enabled") {
		update = true
		if v, ok := d.GetOkExists("delivery_sls_enabled"); ok {
			request["DeliverySlsEnabled"] = v
		}
	}
	if d.HasChange("delivery_sls_project_name") {
		update = true
		if v, ok := d.GetOk("delivery_sls_project_name"); ok {
			request["DeliverySlsProjectName"] = v
		}
	}
	if update {
		action := "SetServiceSettings"
		request["RegionId"] = client.RegionId
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	}
	return resourceAlicloudOosServiceSettingRead(d, meta)
}
func resourceAlicloudOosServiceSettingDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudOosServiceSetting. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

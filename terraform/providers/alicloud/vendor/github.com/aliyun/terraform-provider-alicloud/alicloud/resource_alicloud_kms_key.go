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

func resourceAlicloudKmsKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKmsKeyCreate,
		Read:   resourceAlicloudKmsKeyRead,
		Update: resourceAlicloudKmsKeyUpdate,
		Delete: resourceAlicloudKmsKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"automatic_rotation": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled"}, false),
				Default:      "Disabled",
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Aliyun_AES_256", "Aliyun_SM4", "RSA_2048", "EC_P256", "EC_P256K", "EC_SM2"}, false),
			},
			"is_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'is_enabled' has been deprecated from provider version 1.85.0. New field 'key_state' instead.",
			},
			"key_usage": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENCRYPT/DECRYPT", "SIGN/VERIFY"}, false),
				Default:      "ENCRYPT/DECRYPT",
			},
			"last_rotation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"material_expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_rotation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Aliyun_KMS", "EXTERNAL"}, false),
				Default:      "Aliyun_KMS",
			},
			"pending_window_in_days": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"deletion_window_in_days"},
				ValidateFunc:  validation.IntBetween(7, 30),
			},
			"deletion_window_in_days": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'deletion_window_in_days' has been deprecated from provider version 1.85.0. New field 'pending_window_in_days' instead.",
				ConflictsWith: []string{"pending_window_in_days"},
				ValidateFunc:  validation.IntBetween(7, 30),
			},
			"primary_key_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"SOFTWARE", "HSM"}, false),
				Default:      "SOFTWARE",
			},
			"rotation_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Disabled", "Enabled", "PendingDeletion"}, false),
				ConflictsWith: []string{"key_state"},
			},
			"key_state": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Disabled", "Enabled", "PendingDeletion"}, false),
				ConflictsWith: []string{"status"},
				Deprecated:    "Field 'key_state' has been deprecated from provider version 1.123.1. New field 'status' instead.",
			},
		},
	}
}

func resourceAlicloudKmsKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateKey"
	request := make(map[string]interface{})
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("automatic_rotation"); ok {
		request["EnableAutomaticRotation"] = convertKmsKeyAutomaticRotationRequest(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("key_spec"); ok {
		request["KeySpec"] = v
	}

	if v, ok := d.GetOk("key_usage"); ok {
		request["KeyUsage"] = v
	}

	if v, ok := d.GetOk("origin"); ok {
		request["Origin"] = v
	}

	if v, ok := d.GetOk("protection_level"); ok {
		request["ProtectionLevel"] = v
	}

	if v, ok := d.GetOk("rotation_interval"); ok {
		request["RotationInterval"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_key", action, AlibabaCloudSdkGoERROR)
	}
	responseKeyMetadata := response["KeyMetadata"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseKeyMetadata["KeyId"]))

	return resourceAlicloudKmsKeyUpdate(d, meta)
}
func resourceAlicloudKmsKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}
	object, err := kmsService.DescribeKmsKey(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_key kmsService.DescribeKmsKey Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("arn", object["Arn"])
	d.Set("automatic_rotation", object["AutomaticRotation"])
	d.Set("creator", object["Creator"])
	d.Set("creation_date", object["CreationDate"])
	d.Set("delete_date", object["DeleteDate"])
	d.Set("description", object["Description"])
	d.Set("key_spec", object["KeySpec"])
	d.Set("key_usage", object["KeyUsage"])
	d.Set("last_rotation_date", object["LastRotationDate"])
	d.Set("material_expire_time", object["MaterialExpireTime"])
	d.Set("next_rotation_date", object["NextRotationDate"])
	d.Set("origin", object["Origin"])
	d.Set("primary_key_version", object["PrimaryKeyVersion"])
	d.Set("protection_level", object["ProtectionLevel"])
	d.Set("rotation_interval", object["RotationInterval"])
	d.Set("status", object["KeyState"])
	d.Set("key_state", object["KeyState"])
	return nil
}
func resourceAlicloudKmsKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("description") {
		request := map[string]interface{}{
			"KeyId": d.Id(),
		}
		request["Description"] = d.Get("description")
		action := "UpdateKeyDescription"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("description")
	}
	update := false
	request := map[string]interface{}{
		"KeyId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("automatic_rotation") {
		update = true
	}
	request["EnableAutomaticRotation"] = convertKmsKeyAutomaticRotationRequest(d.Get("automatic_rotation").(string))
	if !d.IsNewResource() && d.HasChange("rotation_interval") {
		update = true
		request["RotationInterval"] = d.Get("rotation_interval")
	}
	if update {
		action := "UpdateRotationPolicy"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("automatic_rotation")
		d.SetPartial("rotation_interval")
	}
	if d.HasChange("status") || d.HasChange("key_state") || d.HasChange("is_enabled") {
		object, err := kmsService.DescribeKmsKey(d.Id())
		if err != nil {
			return WrapError(err)
		}
		var target = ""
		if k, ok := d.GetOk("status"); ok {
			target = k.(string)
		} else if k, ok := d.GetOk("key_state"); ok {
			target = k.(string)
		} else {
			if k, ok := d.GetOk("is_enabled"); ok {
				if k.(bool) {
					target = "Enable"
				} else {
					target = "Disabled"
				}
			}
		}
		if object["KeyState"].(string) != target {
			if target == "Disabled" {
				request := map[string]interface{}{
					"KeyId": d.Id(),
				}
				action := "DisableKey"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			if target == "Enabled" {
				request := map[string]interface{}{
					"KeyId": d.Id(),
				}
				action := "EnableKey"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			d.SetPartial("status")
			d.SetPartial("key_state")
			d.SetPartial("is_enabled")
		}
	}
	d.Partial(false)
	return resourceAlicloudKmsKeyRead(d, meta)
}
func resourceAlicloudKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ScheduleKeyDeletion"
	var response map[string]interface{}
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"KeyId": d.Id(),
	}

	if v, ok := d.GetOk("pending_window_in_days"); ok {
		request["PendingWindowInDays"] = v
	} else if v, ok := d.GetOk("deletion_window_in_days"); ok {
		request["PendingWindowInDays"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "pending_window_in_days" or "deletion_window_in_days" must be set one!`))
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
func convertKmsKeyAutomaticRotationRequest(source interface{}) interface{} {
	switch source {
	case "Disabled":
		return false
	case "Enabled":
		return true
	}
	return false
}

package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudKmsSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKmsSecretCreate,
		Read:   resourceAlicloudKmsSecretRead,
		Update: resourceAlicloudKmsSecretUpdate,
		Delete: resourceAlicloudKmsSecretDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
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
			"enable_automatic_rotation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"encryption_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"force_delete_without_recovery": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"planned_delete_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"recovery_window_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("force_delete_without_recovery").(bool)
				},
			},
			"rotation_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret_data": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"secret_data_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"binary", "text"}, false),
				Default:      "text",
			},
			"secret_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"version_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version_stages": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlicloudKmsSecretCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSecret"
	request := make(map[string]interface{})
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOkExists("enable_automatic_rotation"); ok {
		request["EnableAutomaticRotation"] = v
	}

	if v, ok := d.GetOk("encryption_key_id"); ok {
		request["EncryptionKeyId"] = v
	}

	if v, ok := d.GetOk("rotation_interval"); ok {
		request["RotationInterval"] = v
	}

	request["SecretData"] = d.Get("secret_data")
	if v, ok := d.GetOk("secret_data_type"); ok {
		request["SecretDataType"] = v
	}

	request["SecretName"] = d.Get("secret_name")
	if v, ok := d.GetOk("tags"); ok {
		addTags := make([]JsonTag, 0)
		for key, value := range v.(map[string]interface{}) {
			addTags = append(addTags, JsonTag{
				TagKey:   key,
				TagValue: value.(string),
			})
		}
		tags, err := json.Marshal(addTags)
		if err != nil {
			return WrapError(err)
		}
		request["Tags"] = string(tags)
	}
	request["VersionId"] = d.Get("version_id")
	wait := incrementalWait(3*time.Second, 1*time.Second)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_secret", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SecretName"]))

	return resourceAlicloudKmsSecretUpdate(d, meta)
}
func resourceAlicloudKmsSecretRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}
	object, err := kmsService.DescribeKmsSecret(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_secret kmsService.DescribeKmsSecret Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("secret_name", d.Id())
	d.Set("arn", object["Arn"])
	d.Set("description", object["Description"])
	d.Set("encryption_key_id", object["EncryptionKeyId"])
	d.Set("planned_delete_time", object["PlannedDeleteTime"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}

	getSecretValueObject, err := kmsService.GetSecretValue(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("secret_data", getSecretValueObject["SecretData"])
	d.Set("secret_data_type", getSecretValueObject["SecretDataType"])
	d.Set("version_id", getSecretValueObject["VersionId"])
	d.Set("version_stages", getSecretValueObject["VersionStages"].(map[string]interface{})["VersionStage"])
	return nil
}
func resourceAlicloudKmsSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}
	var response map[string]interface{}
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := kmsService.SetResourceTags(d, "secret"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	if !d.IsNewResource() && d.HasChange("description") {
		request := map[string]interface{}{
			"SecretName": d.Id(),
		}
		request["Description"] = d.Get("description")
		action := "UpdateSecret"
		conn, err := client.NewKmsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 1*time.Second)
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
		"SecretName": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("secret_data") {
		update = true
	}
	request["SecretData"] = d.Get("secret_data")
	if !d.IsNewResource() && d.HasChange("version_id") {
		update = true
	}
	request["VersionId"] = d.Get("version_id")
	if !d.IsNewResource() && d.HasChange("secret_data_type") {
		update = true
		request["SecretDataType"] = d.Get("secret_data_type")
	}
	if d.HasChange("version_stages") {
		update = true
		request["VersionStages"] = convertListToJsonString(d.Get("version_stages").(*schema.Set).List())
	}
	if update {
		action := "PutSecretValue"
		conn, err := client.NewKmsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 1*time.Second)
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
		d.SetPartial("secret_data")
		d.SetPartial("version_id")
		d.SetPartial("secret_data_type")
		d.SetPartial("version_stages")
	}
	d.Partial(false)
	return resourceAlicloudKmsSecretRead(d, meta)
}
func resourceAlicloudKmsSecretDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSecret"
	var response map[string]interface{}
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"SecretName": d.Id(),
	}

	if v, ok := d.GetOkExists("force_delete_without_recovery"); ok {
		request["ForceDeleteWithoutRecovery"] = fmt.Sprintf("%v", v.(bool))
	}
	if v, ok := d.GetOk("recovery_window_in_days"); ok {
		request["RecoveryWindowInDays"] = v
	}
	wait := incrementalWait(3*time.Second, 1*time.Second)
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
		if IsExpectedErrors(err, []string{"Forbidden.ResourceNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

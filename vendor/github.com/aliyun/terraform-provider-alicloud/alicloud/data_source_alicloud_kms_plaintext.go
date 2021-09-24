package alicloud

import (
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudKmsPlaintext() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKmsPlaintextRead,

		Schema: map[string]*schema.Schema{
			"plaintext": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"ciphertext_blob": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAlicloudKmsPlaintextRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	// Since a plaintext has no ID, we create an ID based on
	// current unix time.
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	action := "Decrypt"
	request := make(map[string]interface{})
	request["CiphertextBlob"] = d.Get("ciphertext_blob")
	request["RegionId"] = client.RegionId

	if context := d.Get("encryption_context"); context != nil {
		cm := context.(map[string]interface{})
		contextJson, err := convertMaptoJsonString(cm)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_plaintext", action, AlibabaCloudSdkGoERROR)
		}
		request["EncryptionContext"] = contextJson
	}

	var response map[string]interface{}
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_ciphertext", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Plaintext", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Plaintext", response)
	}
	d.Set("plaintext", resp)
	resp, err = jsonpath.Get("$.KeyId", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.KeyId", response)
	}
	d.Set("key_id", resp)

	return nil
}

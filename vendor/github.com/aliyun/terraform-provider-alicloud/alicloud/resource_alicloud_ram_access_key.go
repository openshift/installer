package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/encryption"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudRamAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamAccessKeyCreate,
		Read:   resourceAlicloudRamAccessKeyRead,
		Update: resourceAlicloudRamAccessKeyUpdate,
		Delete: resourceAlicloudRamAccessKeyDelete,

		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secret_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      Active,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive"}, false),
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"pgp_key": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"key_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypted_secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	request := ram.CreateCreateAccessKeyRequest()
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		request.UserName = v.(string)
	}

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateAccessKey(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_access_key", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ram.CreateAccessKeyResponse)

	if v, ok := d.GetOk("pgp_key"); ok {
		pgpKey := v.(string)
		encryptionKey, err := encryption.RetrieveGPGKey(pgpKey)
		if err != nil {
			return WrapError(err)
		}
		fingerprint, encrypted, err := encryption.EncryptValue(encryptionKey, response.AccessKey.AccessKeySecret, "Alicloud RAM Access Key Secret")
		if err != nil {
			return WrapError(err)
		}
		d.Set("key_fingerprint", fingerprint)
		d.Set("encrypted_secret", encrypted)
	} else {
		if err := d.Set("secret", response.AccessKey.AccessKeySecret); err != nil {
			return WrapError(err)
		}
	}
	if output, ok := d.GetOk("secret_file"); ok && output != nil {
		// create a secret_file and write access key to it.
		writeToFile(output.(string), response.AccessKey)
	}

	d.SetId(response.AccessKey.AccessKeyId)
	err = ramService.WaitForRamAccessKey(d.Id(), request.UserName, Active, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}
	return resourceAlicloudRamAccessKeyUpdate(d, meta)
}

func resourceAlicloudRamAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateUpdateAccessKeyRequest()
	request.RegionId = client.RegionId
	request.UserAccessKeyId = d.Id()
	request.Status = d.Get("status").(string)

	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		request.UserName = v.(string)
	}

	if d.HasChange("status") {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.UpdateAccessKey(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudRamAccessKeyRead(d, meta)
}

func resourceAlicloudRamAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramservice := RamService{client}
	userName := ""
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		userName = v.(string)
	}
	object, err := ramservice.DescribeRamAccessKey(d.Id(), userName)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("status", object.Status)
	return nil
}

func resourceAlicloudRamAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	request := ram.CreateDeleteAccessKeyRequest()
	request.RegionId = client.RegionId
	request.UserAccessKeyId = d.Id()

	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		request.UserName = v.(string)
	}

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.DeleteAccessKey(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, request.UserName, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(ramService.WaitForRamAccessKey(d.Id(), request.UserName, Deleted, DefaultTimeoutMedium))

}

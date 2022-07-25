package alicloud

import (
	"fmt"
	"log"
	"time"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudKvstoreAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKvstoreAccountCreate,
		Read:   resourceAlicloudKvstoreAccountRead,
		Update: resourceAlicloudKvstoreAccountUpdate,
		Delete: resourceAlicloudKvstoreAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"account_privilege": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"RoleReadOnly", "RoleReadWrite", "RoleRepl"}, false),
				Default:      "RoleReadWrite",
			},
			"account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Normal"}, false),
				Default:      "Normal",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
		},
	}
}

func resourceAlicloudKvstoreAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}

	request := r_kvstore.CreateCreateAccountRequest()
	request.AccountName = d.Get("account_name").(string)
	request.AccountPassword = d.Get("account_password").(string)
	if request.AccountPassword == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request.AccountPassword = decryptResp
		}
	}
	if request.AccountPassword == "" {
		return WrapError(Error("One of the 'account_password' and 'kms_encrypted_password' should be set."))
	}
	if v, ok := d.GetOk("account_privilege"); ok {
		request.AccountPrivilege = v.(string)
	}

	if v, ok := d.GetOk("account_type"); ok {
		request.AccountType = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		request.AccountDescription = v.(string)
	}

	request.InstanceId = d.Get("instance_id").(string)
	wait := incrementalWait(3*time.Second, 30*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.CreateAccount(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDeniedDBStatus"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*r_kvstore.CreateAccountResponse)
		d.SetId(fmt.Sprintf("%v:%v", response.InstanceId, response.AcountName))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kvstore_account", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, r_kvstoreService.KvstoreAccountStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudKvstoreAccountRead(d, meta)
}
func resourceAlicloudKvstoreAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}
	object, err := r_kvstoreService.DescribeKvstoreAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kvstore_account r_kvstoreService.DescribeKvstoreAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("account_name", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("account_type", object.AccountType)
	d.Set("description", object.AccountDescription)
	d.Set("status", object.AccountStatus)
	if len(object.DatabasePrivileges.DatabasePrivilege) > 0 {
		d.Set("account_privilege", object.DatabasePrivileges.DatabasePrivilege[0].AccountPrivilege)
	}
	return nil
}
func resourceAlicloudKvstoreAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	r_kvstoreService := R_kvstoreService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("account_privilege") {
		request := r_kvstore.CreateGrantAccountPrivilegeRequest()
		request.AccountName = parts[1]
		request.InstanceId = parts[0]
		request.AccountPrivilege = d.Get("account_privilege").(string)
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.GrantAccountPrivilege(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreAccountStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("account_privilege")
	}
	if d.HasChange("description") {
		request := r_kvstore.CreateModifyAccountDescriptionRequest()
		request.AccountName = parts[1]
		request.InstanceId = parts[0]
		request.AccountDescription = d.Get("description").(string)
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ModifyAccountDescription(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("description")
	}
	if d.HasChange("account_password") || d.HasChange("kms_encrypted_password") {
		request := r_kvstore.CreateResetAccountPasswordRequest()
		request.AccountName = parts[1]
		request.InstanceId = parts[0]
		password := d.Get("account_password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)

		if password == "" && kmsPassword == "" {
			return WrapError(Error("One of the 'account_password' and 'kms_encrypted_password' should be set."))
		}

		if password != "" {
			request.AccountPassword = password
		} else {
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request.AccountPassword = decryptResp
		}
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.ResetAccountPassword(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, r_kvstoreService.KvstoreAccountStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("kms_encrypted_password")
		d.SetPartial("kms_encryption_context")
		d.SetPartial("account_password")
	}
	d.Partial(false)
	return resourceAlicloudKvstoreAccountRead(d, meta)
}
func resourceAlicloudKvstoreAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := r_kvstore.CreateDeleteAccountRequest()
	request.AccountName = parts[1]
	request.InstanceId = parts[0]
	raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DeleteAccount(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidAccountName.NotFound", "InvalidInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

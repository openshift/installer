package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudPolarDBAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBAccountCreate,
		Read:   resourceAlicloudPolarDBAccountRead,
		Update: resourceAlicloudPolarDBAccountUpdate,
		Delete: resourceAlicloudPolarDBAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
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

			"account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string("Normal"), string("Super")}, false),
				Default:      "Normal",
				ForceNew:     true,
			},

			"account_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudPolarDBAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	request := polardb.CreateCreateAccountRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("db_cluster_id").(string)
	request.AccountName = d.Get("account_name").(string)

	password := d.Get("account_password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)

	if password == "" && kmsPassword == "" {
		return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}
	if password != "" {
		request.AccountPassword = password
	} else {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.AccountPassword = decryptResp
	}

	request.AccountType = d.Get("account_type").(string)

	// Description will not be set when account type is normal and it is a API bug
	if v, ok := d.GetOk("account_description"); ok && v.(string) != "" {
		request.AccountDescription = v.(string)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.CreateAccount(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_account", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBClusterId, COLON_SEPARATED, request.AccountName))

	if err := polarDBService.WaitForPolarDBAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudPolarDBAccountRead(d, meta)
}

func resourceAlicloudPolarDBAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	object, err := polarDBService.DescribePolarDBAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	d.Set("db_cluster_id", parts[0])
	d.Set("account_name", object.AccountName)
	d.Set("account_type", object.AccountType)
	d.Set("account_description", object.AccountDescription)

	return nil
}

func resourceAlicloudPolarDBAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	d.Partial(true)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	instanceId := parts[0]
	accountName := parts[1]

	if d.HasChange("account_description") {
		if err := polarDBService.WaitForPolarDBAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := polardb.CreateModifyAccountDescriptionRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = instanceId
		request.AccountName = accountName
		request.AccountDescription = d.Get("account_description").(string)

		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
				return polarDBClient.ModifyAccountDescription(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"ConcurrentTaskExceeded"}) {
					time.Sleep(DefaultIntervalShort * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("account_description")
	}

	if d.HasChange("account_password") || d.HasChange("kms_encrypted_password") {
		if err := polarDBService.WaitForPolarDBAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := polardb.CreateModifyAccountPasswordRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = instanceId
		request.AccountName = accountName

		password := d.Get("account_password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)
		if password == "" && kmsPassword == "" {
			return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
		}

		if password != "" {
			d.SetPartial("account_password")
			request.NewAccountPassword = password
		} else {
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request.NewAccountPassword = decryptResp
			d.SetPartial("kms_encrypted_password")
			d.SetPartial("kms_encryption_context")
		}

		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
				return polarDBClient.ModifyAccountPassword(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"ConcurrentTaskExceeded"}) {
					time.Sleep(DefaultIntervalShort * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("account_password")
	}

	d.Partial(false)
	return resourceAlicloudPolarDBAccountRead(d, meta)
}

func resourceAlicloudPolarDBAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := polardb.CreateDeleteAccountRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DeleteAccount(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentTaskExceeded"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		if IsExpectedErrors(err, []string{"InvalidAccountName.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return polarDBService.WaitForPolarDBAccount(d.Id(), Deleted, DefaultTimeoutMedium)
}

package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudRdsAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRdsAccountCreate,
		Read:   resourceAlicloudRdsAccountRead,
		Update: resourceAlicloudRdsAccountUpdate,
		Delete: resourceAlicloudRdsAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_description": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"description"},
			},
			"description": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'description' has been deprecated from provider version 1.120.0. New field 'account_description' instead.",
				ConflictsWith: []string{"account_description"},
			},
			"account_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9_]{0,61}[a-z0-9]$`), "The name can consist of lowercase letters, numbers, underscores, and must begin with letters and end with letters or numbers"),
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9_]{0,61}[a-z0-9]$`), "The name can consist of lowercase letters, numbers, underscores, and must begin with letters and end with letters or numbers"),
				Deprecated:    "Field 'name' has been deprecated from provider version 1.120.0. New field 'account_name' instead.",
				ConflictsWith: []string{"account_name"},
			},
			"account_password": {
				Type:          schema.TypeString,
				Sensitive:     true,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"password"},
			},
			"password": {
				Type:          schema.TypeString,
				Sensitive:     true,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'password' has been deprecated from provider version 1.120.0. New field 'account_password' instead.",
				ConflictsWith: []string{"account_password"},
			},
			"account_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Normal", "Super"}, false),
				ConflictsWith: []string{"type"},
			},
			"type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Normal", "Super"}, false),
				Deprecated:    "Field 'type' has been deprecated from provider version 1.120.0. New field 'account_type' instead.",
				ConflictsWith: []string{"account_type"},
			},
			"db_instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"instance_id"},
			},
			"instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'instance_id' has been deprecated from provider version 1.120.0. New field 'db_instance_id' instead.",
				ConflictsWith: []string{"db_instance_id"},
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

func resourceAlicloudRdsAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	var response map[string]interface{}
	action := "CreateAccount"
	request := make(map[string]interface{})
	request["SourceIp"] = client.SourceIp
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("account_description"); ok {
		request["AccountDescription"] = v
	} else if v, ok := d.GetOk("description"); ok {
		request["AccountDescription"] = v
	}

	if v, ok := d.GetOk("account_name"); ok {
		request["AccountName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["AccountName"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "name" or "account_name" must be set one!`))
	}

	request["AccountPassword"] = d.Get("account_password")
	if v, ok := d.GetOk("account_password"); ok {
		request["AccountPassword"] = v
	} else if v, ok := d.GetOk("password"); ok {
		request["AccountPassword"] = v
	} else if v, ok := d.GetOk("kms_encrypted_password"); ok {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(v.(string), d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["AccountPassword"] = decryptResp
	} else {
		return WrapError(Error("One of the 'account_password' and 'password' and 'kms_encrypted_password' should be set."))
	}
	if v, ok := d.GetOk("account_type"); ok {
		request["AccountType"] = v
	} else if v, ok := d.GetOk("type"); ok {
		request["AccountType"] = v
	}

	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	} else if v, ok := d.GetOk("instance_id"); ok {
		request["DBInstanceId"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "instance_id" or "db_instance_id" must be set one!`))
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError", "OperationDenied.DBClusterStatus", "OperationDenied.DBInstanceStatus", "OperationDenied.DBStatus", "OperationDenied.OutofUsage"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBInstanceId"], ":", request["AccountName"]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, rdsService.RdsAccountStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudRdsAccountRead(d, meta)
}
func resourceAlicloudRdsAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeRdsAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rds_account rdsService.DescribeRdsAccount Failed!!! %s", err)
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
	d.Set("name", parts[1])
	d.Set("db_instance_id", parts[0])
	d.Set("instance_id", parts[0])
	d.Set("account_description", object["AccountDescription"])
	d.Set("description", object["AccountDescription"])
	d.Set("account_type", object["AccountType"])
	d.Set("type", object["AccountType"])
	d.Set("status", object["AccountStatus"])
	return nil
}
func resourceAlicloudRdsAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"AccountName":  parts[1],
		"DBInstanceId": parts[0],
		"SourceIp":     client.SourceIp,
	}
	if d.HasChange("account_description") {
		update = true
		request["AccountDescription"] = d.Get("account_description")
	} else if d.HasChange("description") {
		update = true
		request["AccountDescription"] = d.Get("description")
	}

	if update {

		action := "ModifyAccountDescription"
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsAccountStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("description")
		d.SetPartial("account_description")
	}
	update = false
	resetAccountPasswordReq := map[string]interface{}{
		"AccountName":  parts[1],
		"DBInstanceId": parts[0],
		"SourceIp":     client.SourceIp,
	}
	if d.HasChange("account_password") {
		update = true
		resetAccountPasswordReq["AccountPassword"] = d.Get("account_password").(string)
	} else if d.HasChange("password") {
		update = true
		resetAccountPasswordReq["AccountPassword"] = d.Get("password").(string)
	} else if d.HasChange("kms_encrypted_password") {
		update = true
		kmsPassword := d.Get("kms_encrypted_password").(string)
		kmsService := KmsService{meta.(*connectivity.AliyunClient)}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		resetAccountPasswordReq["AccountPassword"] = decryptResp
	}
	if update {

		action := "ResetAccountPassword"
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, resetAccountPasswordReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, resetAccountPasswordReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsAccountStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("password")
		d.SetPartial("kms_encrypted_password")
		d.SetPartial("kms_encryption_context")
		d.SetPartial("account_password")
	}
	d.Partial(false)
	return resourceAlicloudRdsAccountRead(d, meta)
}
func resourceAlicloudRdsAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	rdsService := RdsService{client}
	action := "DeleteAccount"
	var response map[string]interface{}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AccountName":  parts[1],
		"DBInstanceId": parts[0],
		"SourceIp":     client.SourceIp,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		object, err := rdsService.DescribeRdsAccount(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		if fmt.Sprint(object["AccountStatus"]) == "Lock" {
			action = "UnlockAccount"
			wait = incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				return resource.NonRetryableError(err)
			}
			action = "DeleteAccount"
			return resource.RetryableError(fmt.Errorf("there need to delete account %s again after unlock it", d.Id()))
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, rdsService.RdsAccountStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

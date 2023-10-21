package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudGpdbAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGpdbAccountCreate,
		Read:   resourceAlicloudGpdbAccountRead,
		Update: resourceAlicloudGpdbAccountUpdate,
		Delete: resourceAlicloudGpdbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z][\w\\_]{2,255}$`), "The description of the account. The description must be 2 to 256 characters in length and can contain letters, digits, underscores (_)."),
			},
			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9_]{1,14}[a-z0-9]$`), "The name of the account. The name must be 2 to 16 characters in length and can contain lower letters, digits, underscores (_)."),
			},
			"account_password": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(8, 32),
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGpdbAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAccount"
	request := make(map[string]interface{})
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("account_description"); ok {
		request["AccountDescription"] = v
	}
	request["AccountName"] = d.Get("account_name")
	request["DBInstanceId"] = d.Get("db_instance_id")
	request["AccountPassword"] = d.Get("account_password")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBInstanceId"], ":", request["AccountName"]))
	gpdbService := GpdbService{client}
	stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, gpdbService.GpdbAccountStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGpdbAccountRead(d, meta)
}
func resourceAlicloudGpdbAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	object, err := gpdbService.DescribeGpdbAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_account gpdbService.DescribeGpdbAccount Failed!!! %s", err)
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
	d.Set("db_instance_id", parts[0])
	d.Set("account_description", object["AccountDescription"])
	d.Set("status", convertGpdbAccountStatusResponse(object["AccountStatus"]))
	return nil
}
func resourceAlicloudGpdbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"AccountName":  parts[1],
		"DBInstanceId": parts[0],
	}

	update := false
	if d.HasChange("account_password") {
		update = true
		if v, ok := d.GetOk("account_password"); ok {
			request["AccountPassword"] = v
		}
	}

	if update {
		action := "ResetAccountPassword"
		conn, err := client.NewGpdbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {

			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

	return resourceAlicloudGpdbAccountRead(d, meta)
}
func resourceAlicloudGpdbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudGpdbAccount. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
func convertGpdbAccountStatusResponse(source interface{}) interface{} {
	switch source {
	case "Creating":
		return "0"
	case "Active":
		return "1"
	case "Deleting":
		return "3"
	}
	return source
}

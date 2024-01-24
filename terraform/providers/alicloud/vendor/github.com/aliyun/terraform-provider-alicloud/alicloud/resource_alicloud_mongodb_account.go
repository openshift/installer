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

func resourceAlicloudMongodbAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongodbAccountCreate,
		Read:   resourceAlicloudMongodbAccountRead,
		Update: resourceAlicloudMongodbAccountUpdate,
		Delete: resourceAlicloudMongodbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\"."), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{1,255}$`), "It must be `2` to `256` characters in length, The description must start with a letter, and can contain letters, digits, underscores (_), and hyphens (-).")),
			},
			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"root"}, false),
			},
			"account_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`[a-zA-Z!#$%^&*()_+-=]{8,32}`), "account_password must consist of uppercase letters, lowercase letters, numbers, and special characters"),
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
		},
	}
}

func resourceAlicloudMongodbAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ResetAccountPassword"
	request := make(map[string]interface{})
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	request["AccountName"] = d.Get("account_name")
	request["AccountPassword"] = d.Get("account_password")
	request["DBInstanceId"] = d.Get("instance_id")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_account", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(request["DBInstanceId"], ":", request["AccountName"]))
	MongoDBService := MongoDBService{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, MongoDBService.MongodbAccountStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudMongodbAccountUpdate(d, meta)
}
func resourceAlicloudMongodbAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	MongoDBService := MongoDBService{client}
	object, err := MongoDBService.DescribeMongodbAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_account MongoDBService.DescribeMongodbAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("account_name", object["AccountName"])
	d.Set("instance_id", object["DBInstanceId"])
	d.Set("account_description", object["AccountDescription"])
	d.Set("status", object["AccountStatus"])
	return nil
}
func resourceAlicloudMongodbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	MongoDBService := MongoDBService{client}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("account_description") {
		request := map[string]interface{}{
			"AccountName":  parts[1],
			"DBInstanceId": parts[0],
		}
		if v, ok := d.GetOk("account_description"); ok {
			request["AccountDescription"] = v
		}
		action := "ModifyAccountDescription"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, MongoDBService.MongodbAccountStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("account_description")
	}

	if !d.IsNewResource() && d.HasChange("account_password") {
		request := map[string]interface{}{
			"AccountName":  parts[1],
			"DBInstanceId": parts[0],
		}
		request["AccountPassword"] = d.Get("account_password")
		action := "ResetAccountPassword"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, MongoDBService.MongodbAccountStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("account_password")
	}

	d.Partial(false)
	return resourceAlicloudMongodbAccountRead(d, meta)
}
func resourceAlicloudMongodbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource Alicloud Mongodb Account. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

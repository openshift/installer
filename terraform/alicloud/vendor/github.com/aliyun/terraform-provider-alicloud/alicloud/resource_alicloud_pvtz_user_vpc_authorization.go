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

func resourceAlicloudPvtzUserVpcAuthorization() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPvtzUserVpcAuthorizationCreate,
		Read:   resourceAlicloudPvtzUserVpcAuthorizationRead,
		Update: resourceAlicloudPvtzUserVpcAuthorizationUpdate,
		Delete: resourceAlicloudPvtzUserVpcAuthorizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auth_channel": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"RESOURCE_DIRECTORY"}, false),
			},
			"auth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"NORMAL", "CLOUD_PRODUCT"}, false),
			},
			"authorized_user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudPvtzUserVpcAuthorizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddUserVpcAuthorization"
	request := make(map[string]interface{})
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("auth_channel"); ok {
		request["AuthChannel"] = v
	}
	if v, ok := d.GetOk("auth_type"); ok {
		request["AuthType"] = v
	}
	request["AuthorizedUserId"] = d.Get("authorized_user_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"System.Busy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pvtz_user_vpc_authorization", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AuthorizedUserId"], ":", request["AuthType"]))

	return resourceAlicloudPvtzUserVpcAuthorizationRead(d, meta)
}
func resourceAlicloudPvtzUserVpcAuthorizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	_, err := pvtzService.DescribePvtzUserVpcAuthorization(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pvtz_user_vpc_authorization pvtzService.DescribePvtzUserVpcAuthorization Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("auth_type", parts[1])
	d.Set("authorized_user_id", parts[0])
	return nil
}
func resourceAlicloudPvtzUserVpcAuthorizationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudPvtzUserVpcAuthorizationRead(d, meta)
}
func resourceAlicloudPvtzUserVpcAuthorizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteUserVpcAuthorization"
	var response map[string]interface{}
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AuthType":         parts[1],
		"AuthorizedUserId": parts[0],
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"System.Busy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ZoneVpc.Auth.NotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

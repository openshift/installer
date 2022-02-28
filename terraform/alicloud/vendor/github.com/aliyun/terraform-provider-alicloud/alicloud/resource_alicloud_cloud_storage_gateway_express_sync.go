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

func resourceAlicloudCloudStorageGatewayExpressSync() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudStorageGatewayExpressSyncCreate,
		Read:   resourceAlicloudCloudStorageGatewayExpressSyncRead,
		Delete: resourceAlicloudCloudStorageGatewayExpressSyncDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bucket_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"bucket_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"express_sync_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5a-zA-Z][\u4e00-\u9fa5a-zA-Z0-9-_.]{0,127}$"), "The length of the name is limited to `1` to `128` characters. It can contain uppercase and lowercase letters, Chinese characters, numbers, English periods (.), underscores (_), or hyphens (-), and must start with  letters."),
			},
		},
	}
}

func resourceAlicloudCloudStorageGatewayExpressSyncCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateExpressSync"
	request := make(map[string]interface{})
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}
	request["BucketName"] = d.Get("bucket_name")
	if v, ok := d.GetOk("bucket_prefix"); ok {
		request["BucketPrefix"] = v
	}
	request["BucketRegion"] = d.Get("bucket_region")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["Name"] = d.Get("express_sync_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_storage_gateway_express_sync", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["ExpressSyncId"]))
	return resourceAlicloudCloudStorageGatewayExpressSyncRead(d, meta)
}
func resourceAlicloudCloudStorageGatewayExpressSyncRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	object, err := sgwService.DescribeExpressSyncs(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_express_sync sgwService.DescribeExpressSyncs Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bucket_name", object["BucketName"])
	d.Set("bucket_prefix", object["BucketPrefix"])
	d.Set("bucket_region", object["BucketRegion"])
	d.Set("description", object["Description"])
	d.Set("express_sync_name", object["Name"])
	return nil
}

func resourceAlicloudCloudStorageGatewayExpressSyncDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	action := "DeleteExpressSync"
	var response map[string]interface{}
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ExpressSyncId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"ExpressSyncNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutDelete), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(d.Id(), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

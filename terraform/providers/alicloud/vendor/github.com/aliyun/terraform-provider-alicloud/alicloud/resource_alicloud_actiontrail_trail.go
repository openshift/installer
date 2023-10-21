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

func resourceAlicloudActiontrailTrail() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudActiontrailTrailCreate,
		Read:   resourceAlicloudActiontrailTrailRead,
		Update: resourceAlicloudActiontrailTrailUpdate,
		Delete: resourceAlicloudActiontrailTrailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"event_rw": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"All", "Read", "Write"}, false),
				Default:      "Write",
			},
			"is_organization_trail": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"mns_topic_arn": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'mns_topic_arn' has been deprecated from version 1.118.0",
			},
			"oss_bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oss_key_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'role_name' has been deprecated from version 1.118.0",
			},
			"oss_write_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sls_project_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sls_write_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enable", "Disable"}, false),
				Default:      "Enable",
			},
			"trail_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'name' has been deprecated from version 1.95.0. Use 'trail_name' instead.",
				ConflictsWith: []string{"trail_name"},
			},
			"trail_region": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"All", "cn-beijing", "cn-hangzhou"}, false),
				Default:      "All",
			},
		},
	}
}

func resourceAlicloudActiontrailTrailCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actiontrailService := ActiontrailService{client}
	var response map[string]interface{}
	action := "CreateTrail"
	request := make(map[string]interface{})
	conn, err := client.NewActiontrailClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("event_rw"); ok {
		request["EventRW"] = v
	}

	if v, ok := d.GetOkExists("is_organization_trail"); ok {
		request["IsOrganizationTrail"] = v
	}

	if v, ok := d.GetOk("oss_bucket_name"); ok {
		request["OssBucketName"] = v
	}

	if v, ok := d.GetOk("oss_key_prefix"); ok {
		request["OssKeyPrefix"] = v
	}

	if v, ok := d.GetOk("oss_write_role_arn"); ok {
		request["OssWriteRoleArn"] = v
	}

	if v, ok := d.GetOk("sls_project_arn"); ok {
		request["SlsProjectArn"] = v
	}

	if v, ok := d.GetOk("sls_write_role_arn"); ok {
		request["SlsWriteRoleArn"] = v
	}

	if v, ok := d.GetOk("trail_name"); ok {
		request["Name"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "name" or "trail_name" must be set one!`))
	}

	if v, ok := d.GetOk("trail_region"); ok {
		request["TrailRegion"] = v
	}

	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_actiontrail_trail", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Name"]))
	stateConf := BuildStateConf([]string{}, []string{"Fresh"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, actiontrailService.ActiontrailTrailStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudActiontrailTrailUpdate(d, meta)
}
func resourceAlicloudActiontrailTrailRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actiontrailService := ActiontrailService{client}
	object, err := actiontrailService.DescribeActiontrailTrail(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_actiontrail_trail actiontrailService.DescribeActiontrailTrail Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("trail_name", d.Id())
	d.Set("name", d.Id())
	d.Set("event_rw", object["EventRW"])
	d.Set("is_organization_trail", object["IsOrganizationTrail"])
	d.Set("oss_bucket_name", object["OssBucketName"])
	d.Set("oss_key_prefix", object["OssKeyPrefix"])
	d.Set("oss_write_role_arn", object["OssWriteRoleArn"])
	d.Set("sls_project_arn", object["SlsProjectArn"])
	d.Set("sls_write_role_arn", object["SlsWriteRoleArn"])
	d.Set("status", object["Status"])
	d.Set("trail_region", object["TrailRegion"])
	return nil
}
func resourceAlicloudActiontrailTrailUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actiontrailService := ActiontrailService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"Name": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("event_rw") {
		update = true
		request["EventRW"] = d.Get("event_rw")
	}
	if !d.IsNewResource() && d.HasChange("oss_bucket_name") {
		update = true
		request["OssBucketName"] = d.Get("oss_bucket_name")
	}
	if !d.IsNewResource() && d.HasChange("oss_key_prefix") {
		update = true
		request["OssKeyPrefix"] = d.Get("oss_key_prefix")
	}
	if !d.IsNewResource() && d.HasChange("oss_write_role_arn") {
		update = true
		request["OssWriteRoleArn"] = d.Get("oss_write_role_arn")
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("sls_project_arn") {
		update = true
		request["SlsProjectArn"] = d.Get("sls_project_arn")
	}
	if !d.IsNewResource() && d.HasChange("sls_write_role_arn") {
		update = true
		request["SlsWriteRoleArn"] = d.Get("sls_write_role_arn")
	}
	if !d.IsNewResource() && d.HasChange("trail_region") {
		update = true
		request["TrailRegion"] = d.Get("trail_region")
	}
	if update {

		if v, ok := d.GetOk("sls_project_arn"); ok {
			request["SlsProjectArn"] = v
		}
		if v, ok := d.GetOk("sls_write_role_arn"); ok {
			request["SlsWriteRoleArn"] = v
		}
		if v, ok := d.GetOk("oss_bucket_name"); ok {
			request["OssBucketName"] = v
		}
		if v, ok := d.GetOk("oss_write_role_arn"); ok {
			request["OssWriteRoleArn"] = v
		}

		action := "UpdateTrail"
		conn, err := client.NewActiontrailClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) || NeedRetry(err) {
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
		d.SetPartial("event_rw")
		d.SetPartial("oss_bucket_name")
		d.SetPartial("oss_key_prefix")
		d.SetPartial("oss_write_role_arn")
		d.SetPartial("sls_project_arn")
		d.SetPartial("sls_write_role_arn")
		d.SetPartial("trail_region")
	}
	if d.HasChange("status") {
		object, err := actiontrailService.DescribeActiontrailTrail(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Disable" {
				request := map[string]interface{}{
					"Name": d.Id(),
				}
				action := "StopLogging"
				conn, err := client.NewActiontrailClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2020-07-06"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) || NeedRetry(err) {
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
				stateConf := BuildStateConf([]string{}, []string{"Disable"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, actiontrailService.ActiontrailTrailStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "Enable" {
				request := map[string]interface{}{
					"Name": d.Id(),
				}
				action := "StartLogging"
				conn, err := client.NewActiontrailClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) || NeedRetry(err) {
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
				stateConf := BuildStateConf([]string{}, []string{"Enable"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, actiontrailService.ActiontrailTrailStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudActiontrailTrailRead(d, meta)
}
func resourceAlicloudActiontrailTrailDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTrail"
	var response map[string]interface{}
	conn, err := client.NewActiontrailClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Name": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"TrailNotFoundException"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

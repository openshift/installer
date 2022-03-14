package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCrChartRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCrChartRepositoryCreate,
		Read:   resourceAlicloudCrChartRepositoryRead,
		Update: resourceAlicloudCrChartRepositoryUpdate,
		Delete: resourceAlicloudCrChartRepositoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repo_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repo_namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repo_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PRIVATE", "PUBLIC"}, false),
			},
			"summary": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudCrChartRepositoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateChartRepository"
	request := make(map[string]interface{})
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	request["InstanceId"] = d.Get("instance_id")
	request["RepoName"] = d.Get("repo_name")
	request["RepoNamespaceName"] = d.Get("repo_namespace_name")
	if v, ok := d.GetOk("repo_type"); ok {
		request["RepoType"] = v
	}
	if v, ok := d.GetOk("summary"); ok {
		request["Summary"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_chart_repository", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["RepoNamespaceName"], ":", request["RepoName"]))

	return resourceAlicloudCrChartRepositoryRead(d, meta)
}
func resourceAlicloudCrChartRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}
	object, err := crService.DescribeCrChartRepository(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_chart_repository crService.DescribeCrChartRepository Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("repo_name", parts[2])
	d.Set("repo_namespace_name", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("repo_type", object["RepoType"])
	d.Set("summary", object["Summary"])
	return nil
}
func resourceAlicloudCrChartRepositoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"RepoName":          parts[2],
		"RepoNamespaceName": parts[1],
		"InstanceId":        parts[0],
	}

	if d.HasChange("repo_type") {
		update = true
	}
	if v, ok := d.GetOk("repo_type"); ok {
		request["RepoType"] = v
	}
	if d.HasChange("summary") {
		update = true
		if v, ok := d.GetOk("summary"); ok {
			request["Summary"] = v
		}
	}
	if update {
		action := "UpdateChartRepository"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudCrChartRepositoryRead(d, meta)
}
func resourceAlicloudCrChartRepositoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteChartRepository"
	var response map[string]interface{}
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RepoName":          parts[2],
		"RepoNamespaceName": parts[1],
		"InstanceId":        parts[0],
	}

	request["InstanceId"] = d.Get("instance_id")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}

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

func resourceAlicloudCrChartNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCrChartNamespaceCreate,
		Read:   resourceAlicloudCrChartNamespaceRead,
		Update: resourceAlicloudCrChartNamespaceUpdate,
		Delete: resourceAlicloudCrChartNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_create_repo": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"default_repo_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PRIVATE", "PUBLIC"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCrChartNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateChartNamespace"
	request := make(map[string]interface{})
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("auto_create_repo"); ok {
		request["AutoCreateRepo"] = v
	}
	if v, ok := d.GetOk("default_repo_type"); ok {
		request["DefaultRepoType"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	request["NamespaceName"] = d.Get("namespace_name")
	request["RegionId"] = client.RegionId
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_chart_namespace", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["NamespaceName"]))
	return resourceAlicloudCrChartNamespaceRead(d, meta)
}
func resourceAlicloudCrChartNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}
	object, err := crService.DescribeCrChartNamespace(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_chart_namespace crService.DescribeCrChartNamespace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("instance_id", parts[0])
	d.Set("namespace_name", parts[1])
	d.Set("auto_create_repo", object["AutoCreateRepo"])
	d.Set("default_repo_type", object["DefaultRepoType"])
	return nil
}
func resourceAlicloudCrChartNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"InstanceId":    parts[0],
		"NamespaceName": parts[1],
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("auto_create_repo") {
		update = true
		if v, ok := d.GetOkExists("auto_create_repo"); ok {
			request["AutoCreateRepo"] = v
		}
	}
	if d.HasChange("default_repo_type") {
		update = true
		if v, ok := d.GetOk("default_repo_type"); ok {
			request["DefaultRepoType"] = v
		}
	}
	if update {
		action := "UpdateChartNamespace"
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
	return resourceAlicloudCrChartNamespaceRead(d, meta)
}
func resourceAlicloudCrChartNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteChartNamespace"
	var response map[string]interface{}
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId":    parts[0],
		"NamespaceName": parts[1],
	}

	request["RegionId"] = client.RegionId
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

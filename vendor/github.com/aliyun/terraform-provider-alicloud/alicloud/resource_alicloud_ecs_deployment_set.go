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

func resourceAlicloudEcsDeploymentSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsDeploymentSetCreate,
		Read:   resourceAlicloudEcsDeploymentSetRead,
		Update: resourceAlicloudEcsDeploymentSetUpdate,
		Delete: resourceAlicloudEcsDeploymentSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"deployment_set_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([\w\\:\-]){2,128}$`), "\t\nThe name of the deployment set.\n\nThe name must be 2 to 128 characters in length and can contain letters, digits, colons (:), underscores (_), and hyphens (-)."),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"domain": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Default"}, false),
			},
			"granularity": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Host"}, false),
			},
			"on_unable_to_redeploy_failed_instance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CancelMembershipAndStart", "KeepStopped"}, false),
			},
			"strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Availability"}, false),
			},
		},
	}
}

func resourceAlicloudEcsDeploymentSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDeploymentSet"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("deployment_set_name"); ok {
		request["DeploymentSetName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}
	if v, ok := d.GetOk("granularity"); ok {
		request["Granularity"] = v
	}
	if v, ok := d.GetOk("on_unable_to_redeploy_failed_instance"); ok {
		request["OnUnableToRedeployFailedInstance"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("strategy"); ok {
		request["Strategy"] = v
	}
	request["ClientToken"] = buildClientToken("CreateDeploymentSet")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_deployment_set", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DeploymentSetId"]))

	return resourceAlicloudEcsDeploymentSetRead(d, meta)
}
func resourceAlicloudEcsDeploymentSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsDeploymentSet(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_deployment_set ecsService.DescribeEcsDeploymentSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("domain", convertEcsDeploymentSetDomainResponse(object["Domain"]))
	d.Set("granularity", convertEcsDeploymentSetGranularityResponse(object["Granularity"]))
	d.Set("deployment_set_name", object["DeploymentSetName"])
	d.Set("description", object["DeploymentSetDescription"])
	d.Set("strategy", object["DeploymentStrategy"])

	return nil
}
func resourceAlicloudEcsDeploymentSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"DeploymentSetId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("deployment_set_name") {
		update = true
		if v, ok := d.GetOk("deployment_set_name"); ok {
			request["DeploymentSetName"] = v
		}
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if update {
		action := "ModifyDeploymentSetAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudEcsDeploymentSetRead(d, meta)
}
func resourceAlicloudEcsDeploymentSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDeploymentSet"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DeploymentSetId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
func convertEcsDeploymentSetDomainResponse(source interface{}) interface{} {
	switch source {
	case "default":
		return "Default"
	}
	return source
}
func convertEcsDeploymentSetGranularityResponse(source interface{}) interface{} {
	switch source {
	case "host":
		return "Host"
	}
	return source
}

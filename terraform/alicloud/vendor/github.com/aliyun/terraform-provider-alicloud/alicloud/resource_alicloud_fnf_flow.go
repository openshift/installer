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

func resourceAlicloudFnfFlow() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFnfFlowCreate,
		Read:   resourceAlicloudFnfFlowRead,
		Update: resourceAlicloudFnfFlowUpdate,
		Delete: resourceAlicloudFnfFlowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"definition": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flow_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DEFAULT", "FDL"}, false),
			},
		},
	}
}

func resourceAlicloudFnfFlowCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateFlow"
	request := make(map[string]interface{})
	conn, err := client.NewFnfClient()
	if err != nil {
		return WrapError(err)
	}
	request["Definition"] = d.Get("definition")
	request["Description"] = d.Get("description")
	request["Name"] = d.Get("name")
	if v, ok := d.GetOk("role_arn"); ok {
		request["RoleArn"] = v
	}

	request["Type"] = d.Get("type")
	wait := incrementalWait(3*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalServerError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fnf_flow", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Name"]))

	return resourceAlicloudFnfFlowRead(d, meta)
}
func resourceAlicloudFnfFlowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fnfService := FnfService{client}
	object, err := fnfService.DescribeFnfFlow(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fnf_flow fnfService.DescribeFnfFlow Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", d.Id())
	d.Set("definition", object["Definition"])
	d.Set("description", object["Description"])
	d.Set("flow_id", object["Id"])
	d.Set("last_modified_time", object["LastModifiedTime"])
	d.Set("role_arn", object["RoleArn"])
	d.Set("type", object["Type"])
	return nil
}
func resourceAlicloudFnfFlowUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"Name": d.Id(),
	}
	if d.HasChange("definition") {
		update = true
		request["Definition"] = d.Get("definition")
	}
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("role_arn") {
		update = true
		request["RoleArn"] = d.Get("role_arn")
	}
	if d.HasChange("type") {
		update = true
		request["Type"] = d.Get("type")
	}
	if update {
		action := "UpdateFlow"
		conn, err := client.NewFnfClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 1*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"ConcurrentUpdateError", "InternalServerError"}) || NeedRetry(err) {
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
	}
	return resourceAlicloudFnfFlowRead(d, meta)
}
func resourceAlicloudFnfFlowDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteFlow"
	var response map[string]interface{}
	conn, err := client.NewFnfClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Name": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-03-15"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalServerError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"FlowNotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

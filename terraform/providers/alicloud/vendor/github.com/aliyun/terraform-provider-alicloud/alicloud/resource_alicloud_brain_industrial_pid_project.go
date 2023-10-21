package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudBrainIndustrialPidProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBrainIndustrialPidProjectCreate,
		Read:   resourceAlicloudBrainIndustrialPidProjectRead,
		Update: resourceAlicloudBrainIndustrialPidProjectUpdate,
		Delete: resourceAlicloudBrainIndustrialPidProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"pid_organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pid_project_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pid_project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudBrainIndustrialPidProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePidProject"
	request := make(map[string]interface{})
	conn, err := client.NewAistudioClient()
	if err != nil {
		return WrapError(err)
	}
	request["PidOrganisationId"] = d.Get("pid_organization_id")
	if v, ok := d.GetOk("pid_project_desc"); ok {
		request["PidProjectDesc"] = v
	}

	request["PidProjectName"] = d.Get("pid_project_name")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("CreatePidProject")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_brain_industrial_pid_project", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		return WrapError(Error("CreatePidProject failed for " + response["Message"].(string)))
	}

	d.SetId(fmt.Sprint(response["PidProjectId"]))

	return resourceAlicloudBrainIndustrialPidProjectRead(d, meta)
}
func resourceAlicloudBrainIndustrialPidProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	brain_industrialService := Brain_industrialService{client}
	object, err := brain_industrialService.DescribeBrainIndustrialPidProject(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_brain_industrial_pid_project brain_industrialService.DescribeBrainIndustrialPidProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("pid_organization_id", object["PidOrganisationId"])
	d.Set("pid_project_desc", object["PidProjectDesc"])
	d.Set("pid_project_name", object["PidProjectName"])
	return nil
}
func resourceAlicloudBrainIndustrialPidProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"PidProjectId": d.Id(),
	}
	if d.HasChange("pid_organization_id") {
		update = true
	}
	request["PidOrganisationId"] = d.Get("pid_organization_id")
	if d.HasChange("pid_project_desc") {
		update = true
		request["PidProjectDesc"] = d.Get("pid_project_desc")
	}
	if d.HasChange("pid_project_name") {
		update = true
		request["PidProjectName"] = d.Get("pid_project_name")
	}
	if update {
		action := "UpdatePidProject"
		conn, err := client.NewAistudioClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if fmt.Sprintf(`%v`, response["Code"]) != "200" {
			return WrapError(Error("UpdatePidProject failed for " + response["Message"].(string)))
		}
	}
	return resourceAlicloudBrainIndustrialPidProjectRead(d, meta)
}
func resourceAlicloudBrainIndustrialPidProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePidProject"
	var response map[string]interface{}
	conn, err := client.NewAistudioClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PidProjectId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		return WrapError(Error("DeletePidProject failed for " + response["Message"].(string)))
	}
	return nil
}

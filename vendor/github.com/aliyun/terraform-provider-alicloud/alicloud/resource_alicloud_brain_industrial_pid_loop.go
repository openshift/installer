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

func resourceAlicloudBrainIndustrialPidLoop() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBrainIndustrialPidLoopCreate,
		Read:   resourceAlicloudBrainIndustrialPidLoopRead,
		Update: resourceAlicloudBrainIndustrialPidLoopUpdate,
		Delete: resourceAlicloudBrainIndustrialPidLoopDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"pid_loop_configuration": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"pid_loop_dcs_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard"}, false),
			},
			"pid_loop_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pid_loop_is_crucial": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"pid_loop_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pid_loop_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1", "2", "3", "4", "5"}, false),
			},
			"pid_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudBrainIndustrialPidLoopCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePidLoop"
	request := make(map[string]interface{})
	conn, err := client.NewAistudioClient()
	if err != nil {
		return WrapError(err)
	}
	request["PidLoopConfiguration"] = d.Get("pid_loop_configuration")
	request["PidLoopDcsType"] = d.Get("pid_loop_dcs_type")
	if v, ok := d.GetOk("pid_loop_desc"); ok {
		request["PidLoopDesc"] = v
	}

	request["IsCrucialPidLoop"] = d.Get("pid_loop_is_crucial")
	request["PidLoopName"] = d.Get("pid_loop_name")
	request["PidLoopType"] = d.Get("pid_loop_type")
	request["PidProjectId"] = d.Get("pid_project_id")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("CreatePidLoop")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_brain_industrial_pid_loop", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		return WrapError(Error("CreatePidLoop failed for " + response["Message"].(string)))
	}

	d.SetId(fmt.Sprint(response["PidLoopId"]))

	return resourceAlicloudBrainIndustrialPidLoopRead(d, meta)
}
func resourceAlicloudBrainIndustrialPidLoopRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	brain_industrialService := Brain_industrialService{client}
	object, err := brain_industrialService.DescribeBrainIndustrialPidLoop(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_brain_industrial_pid_loop brain_industrialService.DescribeBrainIndustrialPidLoop Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("pid_loop_configuration", object["PidLoopConfiguration"])
	d.Set("pid_loop_dcs_type", object["PidLoopDcsType"])
	d.Set("pid_loop_desc", object["PidLoopDesc"])
	d.Set("pid_loop_is_crucial", object["PidLoopIsCrucial"])
	d.Set("pid_loop_name", object["PidLoopName"])
	d.Set("pid_loop_type", object["PidLoopType"])
	d.Set("pid_project_id", object["PidProjectId"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudBrainIndustrialPidLoopUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"PidLoopId": d.Id(),
	}
	if d.HasChange("pid_loop_configuration") {
		update = true
	}
	request["PidLoopConfiguration"] = d.Get("pid_loop_configuration")
	if d.HasChange("pid_loop_dcs_type") {
		update = true
	}
	request["PidLoopDcsType"] = d.Get("pid_loop_dcs_type")
	if d.HasChange("pid_loop_is_crucial") {
		update = true
	}
	request["Crucial"] = d.Get("pid_loop_is_crucial")
	if d.HasChange("pid_loop_name") {
		update = true
	}
	request["PidLoopName"] = d.Get("pid_loop_name")
	if d.HasChange("pid_loop_type") {
		update = true
	}
	request["PidLoopType"] = d.Get("pid_loop_type")
	if d.HasChange("pid_project_id") {
		update = true
	}
	request["PidProjectId"] = d.Get("pid_project_id")
	if d.HasChange("pid_loop_desc") {
		update = true
		request["PidLoopDesc"] = d.Get("pid_loop_desc")
	}
	if update {
		action := "UpdatePidLoop"
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
			return WrapError(Error("UpdatePidLoop failed for " + response["Message"].(string)))
		}
	}
	return resourceAlicloudBrainIndustrialPidLoopRead(d, meta)
}
func resourceAlicloudBrainIndustrialPidLoopDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePidLoop"
	var response map[string]interface{}
	conn, err := client.NewAistudioClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PidLoopId": d.Id(),
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
		return WrapError(Error("DeletePidLoop failed for " + response["Message"].(string)))
	}
	return nil
}

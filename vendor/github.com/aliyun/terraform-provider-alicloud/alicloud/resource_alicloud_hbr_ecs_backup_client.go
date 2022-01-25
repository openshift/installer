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

func resourceAlicloudHbrEcsBackupClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHbrEcsBackupClientCreate,
		Read:   resourceAlicloudHbrEcsBackupClientRead,
		Update: resourceAlicloudHbrEcsBackupClientUpdate,
		Delete: resourceAlicloudHbrEcsBackupClientDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"CLASSIC", "PUBLIC", "VPC"}, false),
			},
			"data_proxy_setting": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"CUSTOM", "DISABLE", "USE_CONTROL_PROXY"}, false),
			},
			"max_cpu_core": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_worker": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"proxy_host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"proxy_password": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"proxy_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"proxy_user": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVATED", "STOPPED", "INSTALL_FAILED", "REGISTERED", "DEACTIVATED", "INSTALLING", "NOT_INSTALLED", "UPGRADING", "UPGRADE_FAILED", "UNINSTALLING", "UNINSTALL_FAILED", "UNKNOWN"}, false),
			},
			"use_https": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudHbrEcsBackupClientCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "InstallBackupClients"
	request := make(map[string]interface{})
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	request["InstanceIds"] = "[\"" + d.Get("instance_id").(string) + "\"]"
	wait := incrementalWait(30*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"CloudAssistNotRunningOnInstance"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_ecs_backup_client", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	hbrService := HbrService{client}
	taskId := fmt.Sprint(response["TaskId"])
	taskConf := BuildStateConf([]string{}, []string{"ACTIVATED"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, hbrService.HbrTaskRefreshFunc(taskId, []string{"INSTALL_FAILED", "DEACTIVATED", "UNKNOWN"}))
	if _, err := taskConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	taskResult, err := hbrService.DescribeHbrTask(taskId)
	if err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	d.SetId(fmt.Sprint(taskResult["ClientId"]))
	return resourceAlicloudHbrEcsBackupClientUpdate(d, meta)
}
func resourceAlicloudHbrEcsBackupClientRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	object, err := hbrService.DescribeHbrEcsBackupClient(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_ecs_backup_client hbrService.DescribeHbrEcsBackupClient Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_id", object["InstanceId"])
	d.Set("status", object["Status"])
	if object["Settings"] != nil {
		d.Set("use_https", object["Settings"].(map[string]interface{})["UseHttps"])
		d.Set("data_network_type", object["Settings"].(map[string]interface{})["DataNetworkType"])
		d.Set("data_proxy_setting", object["Settings"].(map[string]interface{})["DataProxySetting"])
		d.Set("max_cpu_core", object["Settings"].(map[string]interface{})["MaxCpuCore"])
		d.Set("max_worker", object["Settings"].(map[string]interface{})["MaxWorker"])
		d.Set("proxy_host", object["Settings"].(map[string]interface{})["ProxyHost"])
		d.Set("proxy_password", object["Settings"].(map[string]interface{})["ProxyPassword"])
		d.Set("proxy_user", object["Settings"].(map[string]interface{})["ProxyUser"])
		d.Set("proxy_port", fmt.Sprint(formatInt(object["Settings"].(map[string]interface{})["ProxyPort"])))
	}
	return nil
}
func resourceAlicloudHbrEcsBackupClientUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ClientId": d.Id(),
	}
	if d.HasChange("use_https") || d.IsNewResource() {
		update = true
		if v, ok := d.GetOkExists("use_https"); ok {
			request["UseHttps"] = v
		}
	}
	if d.HasChange("data_network_type") {
		update = true
		if v, ok := d.GetOk("data_network_type"); ok {
			request["DataNetworkType"] = v
		}
	}
	if d.HasChange("data_proxy_setting") {
		update = true
		if v, ok := d.GetOk("data_proxy_setting"); ok {
			request["DataProxySetting"] = v
		}
	}
	if d.HasChange("max_cpu_core") {
		update = true
		if v, ok := d.GetOk("max_cpu_core"); ok {
			request["MaxCpuCore"] = v
		}
	}
	if d.HasChange("max_worker") {
		update = true
		if v, ok := d.GetOk("max_worker"); ok {
			request["MaxWorker"] = v
		}
	}
	if d.HasChange("proxy_host") {
		update = true
		if v, ok := d.GetOk("proxy_host"); ok {
			request["ProxyHost"] = v
		}
	}
	if d.HasChange("proxy_password") {
		update = true
		if v, ok := d.GetOk("proxy_password"); ok {
			request["ProxyPassword"] = v
		}
	}
	if d.HasChange("proxy_port") {
		update = true
		if v, ok := d.GetOk("proxy_port"); ok {
			request["ProxyPort"] = v
		}
	}
	if d.HasChange("proxy_user") {
		update = true
		if v, ok := d.GetOk("proxy_user"); ok {
			request["ProxyUser"] = v
		}
	}
	if update {
		action := "UpdateClientSettings"
		conn, err := client.NewHbrClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}

	if d.HasChange("status") {
		client := meta.(*connectivity.AliyunClient)
		hbrService := HbrService{client}
		object, err := hbrService.DescribeHbrEcsBackupClient(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "ACTIVATED" {
				err = resourceAlicloudHbrEcsBackupClientCreate(d, meta)
			} else if target == "STOPPED" {
				err = resourceAlicloudHbrEcsBackupClientUninstall(d, meta)
			} else {
				return WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
			if err != nil {
				return WrapError(err)
			}
		}
	}
	d.SetPartial("status")

	return resourceAlicloudHbrEcsBackupClientRead(d, meta)
}
func resourceAlicloudHbrEcsBackupClientUninstall(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "UninstallBackupClients"
	var response map[string]interface{}
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ClientIds": "[\"" + d.Id() + "\"]",
	}

	request["InstanceIds"] = "[\"" + d.Get("instance_id").(string) + "\"]"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	hbrService := HbrService{client}
	taskConf := BuildStateConf([]string{}, []string{"NOT_INSTALLED"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, hbrService.HbrTaskRefreshFunc(fmt.Sprint(response["TaskId"]), []string{"UNINSTALL_FAILED", "UNKNOWN"}))
	if _, err := taskConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
func resourceAlicloudHbrEcsBackupClientDelete(d *schema.ResourceData, meta interface{}) error {
	err := resourceAlicloudHbrEcsBackupClientUninstall(d, meta)
	if err != nil {
		return WrapError(err)
	}
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBackupClient"
	var response map[string]interface{}
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ClientId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}

package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEhpcJobTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEhpcJobTemplateCreate,
		Read:   resourceAlicloudEhpcJobTemplateRead,
		Update: resourceAlicloudEhpcJobTemplateUpdate,
		Delete: resourceAlicloudEhpcJobTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"array_request": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clock_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"command_line": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"job_template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mem": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"package_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"queue": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"re_runable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"runas_user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stderr_redirect_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stdout_redirect_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"thread": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"variables": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudEhpcJobTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateJobTemplate"
	request := make(map[string]interface{})
	conn, err := client.NewEhpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("array_request"); ok {
		request["ArrayRequest"] = v
	}
	if v, ok := d.GetOk("clock_time"); ok {
		request["ClockTime"] = v
	}
	request["CommandLine"] = d.Get("command_line")
	if v, ok := d.GetOk("gpu"); ok {
		request["Gpu"] = v
	}
	request["Name"] = d.Get("job_template_name")
	if v, ok := d.GetOk("mem"); ok {
		request["Mem"] = v
	}
	if v, ok := d.GetOk("node"); ok {
		request["Node"] = v
	}
	if v, ok := d.GetOk("package_path"); ok {
		request["PackagePath"] = v
	}
	if v, ok := d.GetOk("priority"); ok {
		request["Priority"] = v
	}
	if v, ok := d.GetOk("queue"); ok {
		request["Queue"] = v
	}
	if v, ok := d.GetOkExists("re_runable"); ok {
		request["ReRunable"] = v
	}
	if v, ok := d.GetOk("runas_user"); ok {
		request["RunasUser"] = v
	}
	if v, ok := d.GetOk("stderr_redirect_path"); ok {
		request["StderrRedirectPath"] = v
	}
	if v, ok := d.GetOk("stdout_redirect_path"); ok {
		request["StdoutRedirectPath"] = v
	}
	if v, ok := d.GetOk("task"); ok {
		request["Task"] = v
	}
	if v, ok := d.GetOk("thread"); ok {
		request["Thread"] = v
	}
	if v, ok := d.GetOk("variables"); ok {
		request["Variables"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-04-12"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ehpc_job_template", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TemplateId"]))

	return resourceAlicloudEhpcJobTemplateRead(d, meta)
}
func resourceAlicloudEhpcJobTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ehpcService := EhpcService{client}
	object, err := ehpcService.DescribeEhpcJobTemplate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ehpc_job_template ehpcService.DescribeEhpcJobTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("array_request", object["ArrayRequest"])
	d.Set("clock_time", object["ClockTime"])
	d.Set("command_line", object["CommandLine"])
	if v, ok := object["Gpu"]; ok && fmt.Sprint(v) != "0" {
		d.Set("gpu", formatInt(v))
	}
	d.Set("mem", object["Mem"])
	if v, ok := object["Node"]; ok && fmt.Sprint(v) != "0" {
		d.Set("node", formatInt(v))
	}
	d.Set("package_path", object["PackagePath"])
	if v, ok := object["Priority"]; ok {
		d.Set("priority", formatInt(v))
	}
	d.Set("queue", object["Queue"])
	d.Set("job_template_name", object["Name"])

	if v, ok := object["ReRunable"]; ok {
		v, _ := strconv.ParseBool(v.(string))
		err = d.Set("re_runable", v)
	}

	d.Set("runas_user", object["RunasUser"])
	d.Set("stderr_redirect_path", object["StderrRedirectPath"])
	d.Set("stdout_redirect_path", object["StdoutRedirectPath"])
	if v, ok := object["Task"]; ok && fmt.Sprint(v) != "0" {
		d.Set("task", formatInt(v))
	}
	if v, ok := object["Thread"]; ok && fmt.Sprint(v) != "0" {
		d.Set("thread", formatInt(v))
	}
	d.Set("variables", object["Variables"])
	return nil
}
func resourceAlicloudEhpcJobTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	request := map[string]interface{}{
		"TemplateId": d.Id(),
	}

	request["CommandLine"] = d.Get("command_line")

	request["Name"] = d.Get("job_template_name")

	if v, ok := d.GetOk("array_request"); ok {
		request["ArrayRequest"] = v
	}

	if v, ok := d.GetOk("clock_time"); ok {
		request["ClockTime"] = v
	}

	if v, ok := d.GetOk("gpu"); ok {
		request["Gpu"] = v
	}

	if v, ok := d.GetOk("mem"); ok {
		request["Mem"] = v
	}

	if v, ok := d.GetOk("node"); ok {
		request["Node"] = v
	}

	if v, ok := d.GetOk("package_path"); ok {
		request["PackagePath"] = v
	}

	if v, ok := d.GetOk("priority"); ok {
		request["Priority"] = v
	}

	if v, ok := d.GetOk("queue"); ok {
		request["Queue"] = v
	}

	if v, ok := d.GetOk("re_runable"); ok {
		request["ReRunable"] = v
	}

	if v, ok := d.GetOk("runas_user"); ok {
		request["RunasUser"] = v
	}

	if v, ok := d.GetOk("stderr_redirect_path"); ok {
		request["StderrRedirectPath"] = v
	}

	if v, ok := d.GetOk("stdout_redirect_path"); ok {
		request["StdoutRedirectPath"] = v
	}

	if v, ok := d.GetOk("task"); ok {
		request["Task"] = v
	}

	if v, ok := d.GetOk("thread"); ok {
		request["Thread"] = v
	}

	if v, ok := d.GetOk("variables"); ok {
		request["Variables"] = v
	}

	action := "EditJobTemplate"
	conn, err := client.NewEhpcClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-04-12"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
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

	return resourceAlicloudEhpcJobTemplateRead(d, meta)
}
func resourceAlicloudEhpcJobTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteJobTemplates"
	var response map[string]interface{}
	conn, err := client.NewEhpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Templates": fmt.Sprintf("[{\"Id\":\"%s\"}]", d.Id()),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-04-12"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
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

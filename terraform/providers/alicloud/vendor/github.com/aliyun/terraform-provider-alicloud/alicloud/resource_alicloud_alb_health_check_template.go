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

func resourceAlicloudAlbHealthCheckTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlbHealthCheckTemplateCreate,
		Read:   resourceAlicloudAlbHealthCheckTemplateRead,
		Update: resourceAlicloudAlbHealthCheckTemplateUpdate,
		Delete: resourceAlicloudAlbHealthCheckTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"health_check_codes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("health_check_protocol"); ok && v.(string) == "HTTP" {
						return false
					}
					return true
				},
			},
			"health_check_connect_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"health_check_host": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 80),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("health_check_protocol"); ok && v.(string) == "HTTP" {
						return false
					}
					return true
				},
			},
			"health_check_http_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP1.0", "HTTP1.1"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("health_check_protocol"); ok && v.(string) == "HTTP" {
						return false
					}
					return true
				},
			},
			"health_check_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 50),
			},
			"health_check_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"GET", "HEAD"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("health_check_protocol"); ok && v.(string) == "HTTP" {
						return false
					}
					return true
				},
			},
			"health_check_path": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\/[A-Za-z0-9\$\-_\.\+\/\&\~\@\:\$\^\'\,]{1,79}$`), "The URL must be 1 to 80 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), percent signs (%), question marks (?), number signs (#), and ampersands (&). "),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("health_check_protocol"); ok && v.(string) == "HTTP" {
						return false
					}
					return true
				},
			},
			"health_check_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "TCP"}, false),
			},
			"health_check_template_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z][A-Za-z0-9\\_\\.\\-]{1,127}$`), "The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter."),
			},
			"health_check_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 300),
			},
			"healthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(2, 10),
			},
			"unhealthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(2, 10),
			},
		},
	}
}

func resourceAlicloudAlbHealthCheckTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHealthCheckTemplate"
	request := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if m, ok := d.GetOk("health_check_codes"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("HealthCheckCodes.%d", k+1)] = v.(string)
		}
	}
	if v, ok := d.GetOk("health_check_connect_port"); ok {
		request["HealthCheckConnectPort"] = v
	}
	if v, ok := d.GetOk("health_check_host"); ok {
		request["HealthCheckHost"] = v
	}
	if v, ok := d.GetOk("health_check_http_version"); ok {
		request["HealthCheckHttpVersion"] = v
	}
	if v, ok := d.GetOk("health_check_interval"); ok {
		request["HealthCheckInterval"] = v
	}
	if v, ok := d.GetOk("health_check_method"); ok {
		request["HealthCheckMethod"] = v
	}
	if v, ok := d.GetOk("health_check_path"); ok {
		request["HealthCheckPath"] = v
	}
	if v, ok := d.GetOk("health_check_protocol"); ok {
		request["HealthCheckProtocol"] = v
	}
	request["HealthCheckTemplateName"] = d.Get("health_check_template_name")
	if v, ok := d.GetOk("health_check_timeout"); ok {
		request["HealthCheckTimeout"] = v
	}
	if v, ok := d.GetOk("healthy_threshold"); ok {
		request["HealthyThreshold"] = v
	}
	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		request["UnhealthyThreshold"] = v
	}
	request["ClientToken"] = buildClientToken("CreateHealthCheckTemplate")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IdempotenceProcessing", "QuotaExceeded.HealthCheckTemplatesNum", "SystemBusy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_health_check_template", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["HealthCheckTemplateId"]))

	return resourceAlicloudAlbHealthCheckTemplateRead(d, meta)
}
func resourceAlicloudAlbHealthCheckTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	object, err := albService.DescribeAlbHealthCheckTemplate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_health_check_template albService.DescribeAlbHealthCheckTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if healthCheckCodes, ok := object["HealthCheckCodes"]; ok {
		d.Set("health_check_codes", healthCheckCodes)
	}
	if v, ok := object["HealthCheckConnectPort"]; ok {
		d.Set("health_check_connect_port", formatInt(v))
	}
	if v, ok := object["HealthCheckInterval"]; ok && fmt.Sprint(v) != "0" {
		d.Set("health_check_interval", formatInt(v))
	}
	d.Set("health_check_protocol", object["HealthCheckProtocol"])
	if v, ok := object["HealthCheckProtocol"]; ok && fmt.Sprint(v) == "HTTP" {
		d.Set("health_check_method", object["HealthCheckMethod"])
		d.Set("health_check_path", object["HealthCheckPath"])
		d.Set("health_check_host", object["HealthCheckHost"])
		d.Set("health_check_http_version", object["HealthCheckHttpVersion"])
	}
	d.Set("health_check_template_name", object["HealthCheckTemplateName"])
	if v, ok := object["HealthCheckTimeout"]; ok && fmt.Sprint(v) != "0" {
		d.Set("health_check_timeout", formatInt(v))
	}
	if v, ok := object["HealthyThreshold"]; ok && fmt.Sprint(v) != "0" {
		d.Set("healthy_threshold", formatInt(v))
	}
	if v, ok := object["UnhealthyThreshold"]; ok && fmt.Sprint(v) != "0" {
		d.Set("unhealthy_threshold", formatInt(v))
	}
	return nil
}
func resourceAlicloudAlbHealthCheckTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"HealthCheckTemplateId": d.Id(),
	}
	if d.HasChange("health_check_codes") {
		update = true
		if m, ok := d.GetOk("health_check_codes"); ok {
			for k, v := range m.([]interface{}) {
				request[fmt.Sprintf("HealthCheckCodes.%d", k+1)] = v.(string)
			}
		}
	}
	if d.HasChange("health_check_connect_port") {
		update = true
		if v, ok := d.GetOkExists("health_check_connect_port"); ok {
			request["HealthCheckConnectPort"] = v
		}
	}
	if d.HasChange("health_check_host") {
		update = true
		if v, ok := d.GetOk("health_check_host"); ok {
			request["HealthCheckHost"] = v
		}
	}
	if d.HasChange("health_check_http_version") {
		update = true
		if v, ok := d.GetOk("health_check_http_version"); ok {
			request["HealthCheckHttpVersion"] = v
		}
	}
	if d.HasChange("health_check_interval") {
		update = true
		if v, ok := d.GetOk("health_check_interval"); ok {
			request["HealthCheckInterval"] = v
		}
	}
	if d.HasChange("health_check_method") {
		update = true
		if v, ok := d.GetOk("health_check_method"); ok {
			request["HealthCheckMethod"] = v
		}
	}
	if d.HasChange("health_check_path") {
		update = true
		if v, ok := d.GetOk("health_check_path"); ok {
			request["HealthCheckPath"] = v
		}
	}
	if d.HasChange("health_check_protocol") {
		update = true
		if v, ok := d.GetOk("health_check_protocol"); ok {
			request["HealthCheckProtocol"] = v
		}
	}
	if d.HasChange("health_check_template_name") {
		update = true
		request["HealthCheckTemplateName"] = d.Get("health_check_template_name")
	}
	if d.HasChange("health_check_timeout") {
		update = true
		if v, ok := d.GetOk("health_check_timeout"); ok {
			request["HealthCheckTimeout"] = v
		}
	}
	if d.HasChange("healthy_threshold") {
		update = true
		if v, ok := d.GetOk("healthy_threshold"); ok {
			request["HealthyThreshold"] = v
		}
	}
	if d.HasChange("unhealthy_threshold") {
		update = true
		if v, ok := d.GetOk("unhealthy_threshold"); ok {
			request["UnhealthyThreshold"] = v
		}
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
		action := "UpdateHealthCheckTemplateAttribute"
		request["ClientToken"] = buildClientToken("UpdateHealthCheckTemplateAttribute")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"IdempotenceProcessing", "IncorrectStatus.HealthCheckTemplate", "SystemBusy"}) || NeedRetry(err) {
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
	return resourceAlicloudAlbHealthCheckTemplateRead(d, meta)
}
func resourceAlicloudAlbHealthCheckTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteHealthCheckTemplates"
	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"HealthCheckTemplateIds.1": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = buildClientToken("DeleteHealthCheckTemplates")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IdempotenceProcessing", "IncorrectStatus.HealthCheckTemplate", "SystemBusy"}) || NeedRetry(err) {
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

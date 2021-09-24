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

func resourceAlicloudOosTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOosTemplateCreate,
		Read:   resourceAlicloudOosTemplateRead,
		Update: resourceAlicloudOosTemplateUpdate,
		Delete: resourceAlicloudOosTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_delete_executions": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"content": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"has_trigger": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"template_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudOosTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateTemplate"
	request := make(map[string]interface{})
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	request["Content"] = d.Get("content")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("tags"); ok {
		respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Tags"] = respJson
	}
	request["TemplateName"] = d.Get("template_name")
	if v, ok := d.GetOk("version_name"); ok {
		request["VersionName"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_template", action, AlibabaCloudSdkGoERROR)
	}
	responseTemplate := response["Template"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseTemplate["TemplateName"]))

	return resourceAlicloudOosTemplateRead(d, meta)
}
func resourceAlicloudOosTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}
	object, err := oosService.DescribeOosTemplate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oos_template oosService.DescribeOosTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("template_name", d.Id())
	d.Set("created_by", object["CreatedBy"])
	d.Set("created_date", object["CreatedDate"])
	d.Set("description", object["Description"])
	d.Set("has_trigger", object["HasTrigger"])
	d.Set("share_type", object["ShareType"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v))
	}
	d.Set("template_format", object["TemplateFormat"])
	d.Set("template_id", object["TemplateId"])
	d.Set("template_type", object["TemplateType"])
	d.Set("template_version", object["TemplateVersion"])
	d.Set("updated_by", object["UpdatedBy"])
	d.Set("updated_date", object["UpdatedDate"])
	return nil
}
func resourceAlicloudOosTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"TemplateName": d.Id(),
	}
	if d.HasChange("content") {
		update = true
	}
	request["Content"] = d.Get("content")
	request["RegionId"] = client.RegionId
	if d.HasChange("tags") {
		update = true
		respJson, err := convertMaptoJsonString(d.Get("tags").(map[string]interface{}))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_template", "UpdateTemplate", AlibabaCloudSdkGoERROR)
		}
		request["Tags"] = respJson
	}
	if d.HasChange("version_name") {
		update = true
		request["VersionName"] = d.Get("version_name")
	}
	if update {
		action := "UpdateTemplate"
		conn, err := client.NewOosClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	}
	return resourceAlicloudOosTemplateRead(d, meta)
}
func resourceAlicloudOosTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTemplate"
	var response map[string]interface{}
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TemplateName": d.Id(),
	}

	if v, ok := d.GetOkExists("auto_delete_executions"); ok {
		request["AutoDeleteExecutions"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Template"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

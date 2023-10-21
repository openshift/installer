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

func resourceAlicloudImpAppTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudImpAppTemplateCreate,
		Read:   resourceAlicloudImpAppTemplateRead,
		Update: resourceAlicloudImpAppTemplateUpdate,
		Delete: resourceAlicloudImpAppTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"component_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"component.live", "component.liveRecord", "component.liveBeauty", "component.rtc", "component.rtcRecord", "component.im", "component.whiteboard", "component.liveSecurity", "component.chatSecurity"}, false),
				},
				ForceNew: true,
			},
			"config_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"integration_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"paasSDK", "standardRoom"}, false),
			},
			"scene": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"business", "classroom"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudImpAppTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAppTemplate"
	request := make(map[string]interface{})
	conn, err := client.NewImpClient()
	if err != nil {
		return WrapError(err)
	}
	request["AppTemplateName"] = d.Get("app_template_name")
	request["ComponentList"] = convertListToJsonString(d.Get("component_list").([]interface{}))
	if v, ok := d.GetOk("integration_mode"); ok {
		request["IntegrationMode"] = v
	}
	if v, ok := d.GetOk("scene"); ok {
		request["Scene"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-06-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_imp_app_template", action, AlibabaCloudSdkGoERROR)
	}
	responseResult := response["Result"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseResult["AppTemplateId"]))

	return resourceAlicloudImpAppTemplateRead(d, meta)
}
func resourceAlicloudImpAppTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	impService := ImpService{client}
	object, err := impService.DescribeImpAppTemplate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_imp_app_template impService.DescribeImpAppTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("app_template_name", object["AppTemplateName"])
	if componentList, ok := object["ComponentList"]; ok && componentList != nil {
		d.Set("component_list", componentList)
	}
	if configList, ok := object["ConfigList"]; ok && configList != nil {
		configListMaps := make([]map[string]interface{}, 0)
		for _, configListItem := range configList.([]interface{}) {
			if configListItemMap, ok := configListItem.(map[string]interface{}); ok {
				configListMap := make(map[string]interface{}, 0)
				configListMap["key"] = configListItemMap["Key"]
				configListMap["value"] = configListItemMap["Value"]
				configListMaps = append(configListMaps, configListMap)
			}
		}
		d.Set("config_list", configListMaps)
	}

	d.Set("integration_mode", object["IntegrationMode"])
	d.Set("scene", object["Scene"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudImpAppTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"AppTemplateId": d.Id(),
	}
	if d.HasChange("app_template_name") {
		update = true
		request["AppTemplateName"] = d.Get("app_template_name")
	}
	if update {
		action := "UpdateAppTemplate"
		conn, err := client.NewImpClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-06-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("app_template_name")
	}
	update = false
	updateAppTemplateConfigReq := map[string]interface{}{
		"AppTemplateId": d.Id(),
	}
	if d.HasChanges("config_list") {
		update = true
		if v, ok := d.GetOk("config_list"); ok {
			configListMaps := make([]map[string]interface{}, 0)
			for _, configList := range v.(*schema.Set).List() {
				configListArg := configList.(map[string]interface{})
				configListMap := map[string]interface{}{
					"Key":   configListArg["key"],
					"Value": configListArg["value"],
				}
				configListMaps = append(configListMaps, configListMap)
			}
			if v, err := convertListMapToJsonString(configListMaps); err != nil {
				return WrapError(err)
			} else {
				updateAppTemplateConfigReq["ConfigList"] = v
			}
		}
	}
	if update {
		action := "UpdateAppTemplateConfig"
		conn, err := client.NewImpClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-06-30"), StringPointer("AK"), nil, updateAppTemplateConfigReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateAppTemplateConfigReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("key")
		d.SetPartial("value")
	}
	d.Partial(false)
	return resourceAlicloudImpAppTemplateRead(d, meta)
}
func resourceAlicloudImpAppTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAppTemplate"
	var response map[string]interface{}
	conn, err := client.NewImpClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AppTemplateId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-06-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidAppTemplateId.App.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

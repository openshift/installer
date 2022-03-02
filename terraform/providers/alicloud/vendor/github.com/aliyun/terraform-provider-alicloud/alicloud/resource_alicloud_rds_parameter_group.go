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

func resourceAlicloudRdsParameterGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRdsParameterGroupCreate,
		Read:   resourceAlicloudRdsParameterGroupRead,
		Update: resourceAlicloudRdsParameterGroupUpdate,
		Delete: resourceAlicloudRdsParameterGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"engine": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"mariadb", "mysql"}, false),
			},
			"engine_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"10.3", "5.1", "5.5", "5.6", "5.7", "8.0"}, false),
			},
			"param_detail": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"param_value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"parameter_group_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameter_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudRdsParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateParameterGroup"
	request := make(map[string]interface{})
	request["SourceIp"] = client.SourceIp
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	request["Engine"] = d.Get("engine")
	request["EngineVersion"] = d.Get("engine_version")
	list := d.Get("param_detail").(*schema.Set).List()
	paramMap := map[string]interface{}{}
	for _, v := range list {
		v := v.(map[string]interface{})
		paramMap[v["param_name"].(string)] = v["param_value"]
	}
	paramStr, _ := convertMaptoJsonString(paramMap)
	request["Parameters"] = paramStr
	if v, ok := d.GetOk("parameter_group_desc"); ok {
		request["ParameterGroupDesc"] = v
	}

	request["ParameterGroupName"] = d.Get("parameter_group_name")
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_parameter_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ParameterGroupId"]))

	return resourceAlicloudRdsParameterGroupRead(d, meta)
}
func resourceAlicloudRdsParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeRdsParameterGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rds_parameter_group rdsService.DescribeRdsParameterGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("engine", object["Engine"])
	d.Set("engine_version", object["EngineVersion"])
	if v, ok := object["ParamDetail"].(map[string]interface{})["ParameterDetail"].([]interface{}); ok {
		parameterDetail := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			parameterDetail = append(parameterDetail, map[string]interface{}{
				"param_name":  item["ParamName"],
				"param_value": item["ParamValue"],
			})
		}
		if err := d.Set("param_detail", parameterDetail); err != nil {
			return WrapError(err)
		}
	}
	d.Set("parameter_group_desc", object["ParameterGroupDesc"])
	d.Set("parameter_group_name", object["ParameterGroupName"])
	return nil
}
func resourceAlicloudRdsParameterGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ParameterGroupId": d.Id(),
		"SourceIp":         client.SourceIp,
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("param_detail") {
		update = true
		list := d.Get("param_detail").(*schema.Set).List()
		paramMap := map[string]interface{}{}
		for _, v := range list {
			v := v.(map[string]interface{})
			paramMap[v["param_name"].(string)] = v["param_value"]
		}
		paramStr, _ := convertMaptoJsonString(paramMap)
		request["Parameters"] = paramStr
	}
	if d.HasChange("parameter_group_desc") {
		update = true
		request["ParameterGroupDesc"] = d.Get("parameter_group_desc")
	}
	if d.HasChange("parameter_group_name") {
		update = true
		request["ParameterGroupName"] = d.Get("parameter_group_name")
	}
	if update {
		action := "ModifyParameterGroup"
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudRdsParameterGroupRead(d, meta)
}
func resourceAlicloudRdsParameterGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteParameterGroup"
	var response map[string]interface{}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ParameterGroupId": d.Id(),
		"SourceIp":         client.SourceIp,
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}

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

func resourceAlicloudRosStackGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRosStackGroupCreate,
		Read:   resourceAlicloudRosStackGroupRead,
		Update: resourceAlicloudRosStackGroupUpdate,
		Delete: resourceAlicloudRosStackGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_ids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"administration_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execution_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"operation_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operation_preferences": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameter_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"region_ids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stack_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stack_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"template_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudRosStackGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	var response map[string]interface{}
	action := "CreateStackGroup"
	request := make(map[string]interface{})
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("administration_role_name"); ok {
		request["AdministrationRoleName"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("execution_role_name"); ok {
		request["ExecutionRoleName"] = v
	}

	if v, ok := d.GetOk("parameters"); ok {
		parameters := make([]map[string]interface{}, len(v.(*schema.Set).List()))
		for i, j := range v.(*schema.Set).List() {
			parameters[i] = make(map[string]interface{})
			parameters[i]["ParameterKey"] = j.(map[string]interface{})["parameter_key"]
			parameters[i]["ParameterValue"] = j.(map[string]interface{})["parameter_value"]
		}
		request["Parameters"] = parameters

	}

	request["RegionId"] = client.RegionId
	request["StackGroupName"] = d.Get("stack_group_name")
	if v, ok := d.GetOk("template_body"); ok {
		request["TemplateBody"] = v
	}

	if v, ok := d.GetOk("template_url"); ok {
		request["TemplateURL"] = v
	}

	if v, ok := d.GetOk("template_version"); ok {
		request["TemplateVersion"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("CreateStackGroup")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ros_stack_group", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(request["StackGroupName"]))
	stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, rosService.RosStackGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudRosStackGroupRead(d, meta)
}
func resourceAlicloudRosStackGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	object, err := rosService.DescribeRosStackGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ros_stack_group rosService.DescribeRosStackGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("stack_group_name", d.Id())
	d.Set("administration_role_name", object["AdministrationRoleName"])
	d.Set("description", object["Description"])
	d.Set("execution_role_name", object["ExecutionRoleName"])

	parameters := make([]map[string]interface{}, 0)
	if parametersList, ok := object["Parameters"].([]interface{}); ok {
		for _, v := range parametersList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"parameter_key":   m1["ParameterKey"],
					"parameter_value": m1["ParameterValue"],
				}
				parameters = append(parameters, temp1)

			}
		}
	}
	if err := d.Set("parameters", parameters); err != nil {
		return WrapError(err)
	}
	d.Set("stack_group_id", object["StackGroupId"])
	d.Set("status", object["Status"])
	d.Set("template_body", object["TemplateBody"])
	return nil
}
func resourceAlicloudRosStackGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"StackGroupName": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("administration_role_name") {
		update = true
		request["AdministrationRoleName"] = d.Get("administration_role_name")
	}
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("execution_role_name") {
		update = true
		request["ExecutionRoleName"] = d.Get("execution_role_name")
	}
	if d.HasChange("parameters") {
		update = true
		parameters := make([]map[string]interface{}, len(d.Get("parameters").(*schema.Set).List()))
		for i, j := range d.Get("parameters").(*schema.Set).List() {
			parameters[i] = make(map[string]interface{})
			parameters[i]["ParameterKey"] = j.(map[string]interface{})["parameter_key"]
			parameters[i]["ParameterValue"] = j.(map[string]interface{})["parameter_value"]
		}
		request["Parameters"] = parameters

	}
	if d.HasChange("template_body") {
		update = true
		request["TemplateBody"] = d.Get("template_body")
	}
	if update {
		if _, ok := d.GetOk("account_ids"); ok {
			request["AccountIds"] = d.Get("account_ids")
		}
		if _, ok := d.GetOk("operation_description"); ok {
			request["OperationDescription"] = d.Get("operation_description")
		}
		if _, ok := d.GetOk("operation_preferences"); ok {
			request["OperationPreferences"] = d.Get("operation_preferences")
		}
		if _, ok := d.GetOk("region_ids"); ok {
			request["RegionIds"] = d.Get("region_ids")
		}
		if _, ok := d.GetOk("template_url"); ok {
			request["TemplateURL"] = d.Get("template_url")
		}
		if _, ok := d.GetOk("template_version"); ok {
			request["TemplateVersion"] = d.Get("template_version")
		}
		action := "UpdateStackGroup"
		conn, err := client.NewRosClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		request["ClientToken"] = buildClientToken("UpdateStackGroup")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, rosService.RosStackGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudRosStackGroupRead(d, meta)
}
func resourceAlicloudRosStackGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteStackGroup"
	var response map[string]interface{}
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"StackGroupName": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

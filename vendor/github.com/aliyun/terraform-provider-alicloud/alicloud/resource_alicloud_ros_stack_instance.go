package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRosStackInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRosStackInstanceCreate,
		Read:   resourceAlicloudRosStackInstanceRead,
		Update: resourceAlicloudRosStackInstanceUpdate,
		Delete: resourceAlicloudRosStackInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"operation_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operation_preferences": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameter_overrides": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"parameter_value": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
				Sensitive: true,
			},
			"stack_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][-_a-zA-Z0-9]{1,254}$`), "The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter."),
			},
			"stack_instance_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\d+$`), "It can only be made of numbers."),
			},
			"stack_instance_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timeout_in_minutes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"retain_stacks": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudRosStackInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateStackInstances"
	request := make(map[string]interface{})
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("parameter_overrides"); ok {
		for parameterOverridesPtr, parameterOverrides := range v.(*schema.Set).List() {
			parameterOverridesArg := parameterOverrides.(map[string]interface{})
			request["ParameterOverrides."+fmt.Sprint(parameterOverridesPtr+1)+".ParameterKey"] = parameterOverridesArg["parameter_key"]
			request["ParameterOverrides."+fmt.Sprint(parameterOverridesPtr+1)+".ParameterValue"] = parameterOverridesArg["parameter_value"]
		}
	}
	request["RegionId"] = client.RegionId
	request["StackGroupName"] = d.Get("stack_group_name")
	request["AccountIds"] = convertListToJsonString([]interface{}{d.Get("stack_instance_account_id")})
	request["RegionIds"] = convertListToJsonString([]interface{}{d.Get("stack_instance_region_id")})
	if v, ok := d.GetOk("timeout_in_minutes"); ok {
		request["TimeoutInMinutes"] = v
	}
	if v, ok := d.GetOk("operation_preferences"); ok {
		request["OperationPreferences"] = v
	}
	request["ClientToken"] = buildClientToken("CreateStackInstances")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ros_stack_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["StackGroupName"], ":", d.Get("stack_instance_account_id"), ":", d.Get("stack_instance_region_id")))
	rosService := RosService{client}
	stateConf := BuildStateConf([]string{}, []string{"CURRENT"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, rosService.RosStackInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudRosStackInstanceRead(d, meta)
}
func resourceAlicloudRosStackInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	object, err := rosService.DescribeRosStackInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ros_stack_instance rosService.DescribeRosStackInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("stack_group_name", parts[0])
	d.Set("stack_instance_account_id", parts[1])
	d.Set("stack_instance_region_id", parts[2])
	if parameterOverridesList, ok := object["ParameterOverrides"]; ok && parameterOverridesList != nil {
		parameterOverridesMaps := make([]map[string]interface{}, 0)
		for _, parameterOverridesListItem := range parameterOverridesList.([]interface{}) {
			if v, ok := parameterOverridesListItem.(map[string]interface{}); ok {
				parameterOverridesListItemMap := make(map[string]interface{})
				parameterOverridesListItemMap["parameter_key"] = v["ParameterKey"]
				parameterOverridesListItemMap["parameter_value"] = v["ParameterValue"]
				parameterOverridesMaps = append(parameterOverridesMaps, parameterOverridesListItemMap)
			}
		}
		d.Set("parameter_overrides", parameterOverridesMaps)
	}

	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudRosStackInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"StackGroupName": parts[0],
		"AccountIds":     convertListToJsonString([]interface{}{parts[1]}),
		"RegionIds":      convertListToJsonString([]interface{}{parts[2]}),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("parameter_overrides") {
		update = true
		if v, ok := d.GetOk("parameter_overrides"); ok {
			for parameterOverridesPtr, parameterOverrides := range v.(*schema.Set).List() {
				parameterOverridesArg := parameterOverrides.(map[string]interface{})
				request["ParameterOverrides."+fmt.Sprint(parameterOverridesPtr+1)+".ParameterKey"] = parameterOverridesArg["parameter_key"]
				request["ParameterOverrides."+fmt.Sprint(parameterOverridesPtr+1)+".ParameterValue"] = parameterOverridesArg["parameter_value"]
			}
		}
	}
	if d.HasChange("operation_description") {
		update = true
		if v, ok := d.GetOk("operation_description"); ok {
			request["OperationDescription"] = v
		}
	}
	if d.HasChange("operation_preferences") {
		update = true
		if v, ok := d.GetOk("operation_preferences"); ok {
			request["OperationPreferences"] = v
		}
	}
	if d.HasChange("timeout_in_minutes") {
		update = true
		if v, ok := d.GetOk("timeout_in_minutes"); ok {
			request["TimeoutInMinutes"] = v
		}
	}

	if update {
		action := "UpdateStackInstances"
		conn, err := client.NewRosClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateStackInstances")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"StackGroupOperationInProgress"}) {
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
	return resourceAlicloudRosStackInstanceRead(d, meta)
}
func resourceAlicloudRosStackInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteStackInstances"
	var response map[string]interface{}
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"StackGroupName": parts[0],
		"AccountIds":     convertListToJsonString([]interface{}{parts[1]}),
		"RegionIds":      convertListToJsonString([]interface{}{parts[2]}),
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOkExists("retain_stacks"); ok {
		request["RetainStacks"] = v
	} else {
		request["RetainStacks"] = false
	}
	request["ClientToken"] = buildClientToken("DeleteStackInstances")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"StackGroupOperationInProgress"}) {
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rosService.RosStackInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

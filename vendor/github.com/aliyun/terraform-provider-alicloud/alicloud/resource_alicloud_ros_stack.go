package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudRosStack() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRosStackCreate,
		Read:   resourceAlicloudRosStackRead,
		Update: resourceAlicloudRosStackUpdate,
		Delete: resourceAlicloudRosStackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_option": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deletion_protection": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled"}, false),
				Default:      "Disabled",
			},
			"disable_rollback": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"notification_urls": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
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
							Required: true,
						},
					},
				},
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"replacement_option": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"retain_all_resources": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"retain_resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"stack_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"stack_policy_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"stack_policy_during_update_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stack_policy_during_update_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stack_policy_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
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
			"timeout_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"use_previous_parameters": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudRosStackCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	var response map[string]interface{}
	action := "CreateStack"
	request := make(map[string]interface{})
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("create_option"); ok {
		request["CreateOption"] = v
	}

	if v, ok := d.GetOk("deletion_protection"); ok {
		request["DeletionProtection"] = v
	}

	if v, ok := d.GetOkExists("disable_rollback"); ok {
		request["DisableRollback"] = v
	}

	if v, ok := d.GetOk("notification_urls"); ok {
		request["NotificationURLs"] = v.(*schema.Set).List()
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

	if v, ok := d.GetOk("ram_role_name"); ok {
		request["RamRoleName"] = v
	}

	request["RegionId"] = client.RegionId
	request["StackName"] = d.Get("stack_name")
	if v, ok := d.GetOk("stack_policy_body"); ok {
		request["StackPolicyBody"] = v
	}

	if v, ok := d.GetOk("stack_policy_url"); ok {
		request["StackPolicyURL"] = v
	}

	if v, ok := d.GetOk("template_body"); ok {
		request["TemplateBody"] = v
	}

	if v, ok := d.GetOk("template_url"); ok {
		request["TemplateURL"] = v
	}

	if v, ok := d.GetOk("template_version"); ok {
		request["TemplateVersion"] = v
	}

	if v, ok := d.GetOk("timeout_in_minutes"); ok {
		request["TimeoutInMinutes"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("CreateStack")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ros_stack", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(response["StackId"]))
	stateConf := BuildStateConf([]string{}, []string{"CREATE_COMPLETE"}, d.Timeout(schema.TimeoutCreate), 100*time.Second, rosService.RosStackStateRefreshFunc(d.Id(), []string{"CREATE_FAILED", "CREATE_ROLLBACK_COMPLETE", "CREATE_ROLLBACK_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudRosStackUpdate(d, meta)
}
func resourceAlicloudRosStackRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	object, err := rosService.DescribeRosStack(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ros_stack rosService.DescribeRosStack Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("deletion_protection", object["DeletionProtection"])
	d.Set("disable_rollback", object["DisableRollback"])

	parameters := make([]map[string]interface{}, 0)
	if parametersList, ok := object["Parameters"].([]interface{}); ok {
		for _, v := range parametersList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"parameter_key":   m1["ParameterKey"],
					"parameter_value": m1["ParameterValue"],
				}
				if !strings.HasPrefix(v.(map[string]interface{})["ParameterKey"].(string), "ALIYUN::") {
					parameters = append(parameters, temp1)
				}
			}
		}
	}
	if err := d.Set("parameters", parameters); err != nil {
		return WrapError(err)
	}
	d.Set("ram_role_name", object["RamRoleName"])
	d.Set("stack_name", object["StackName"])
	d.Set("status", object["Status"])
	d.Set("timeout_in_minutes", formatInt(object["TimeoutInMinutes"]))

	listTagResourcesObject, err := rosService.ListTagResources(d.Id(), "stack")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))

	getStackPolicyObject, err := rosService.GetStackPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("stack_policy_body", getStackPolicyObject["StackPolicyBody"])
	return nil
}

func resourceAlicloudRosStackUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := rosService.SetResourceTags(d, "stack"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := map[string]interface{}{
		"StackId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("disable_rollback") {
		update = true
		request["DisableRollback"] = d.Get("disable_rollback")
	}
	if !d.IsNewResource() && d.HasChange("parameters") {
		update = true
		parameters := make([]map[string]interface{}, len(d.Get("parameters").(*schema.Set).List()))
		for i, v := range d.Get("parameters").(*schema.Set).List() {
			parameters[i] = make(map[string]interface{})
			parameters[i]["ParameterKey"] = v.(map[string]interface{})["parameter_key"]
			parameters[i]["ParameterValue"] = v.(map[string]interface{})["parameter_value"]
		}
		request["Parameters"] = parameters
	}
	if !d.IsNewResource() && d.HasChange("ram_role_name") {
		update = true
		request["RamRoleName"] = d.Get("ram_role_name")
	}
	if !d.IsNewResource() && d.HasChange("stack_policy_body") {
		update = true
		request["StackPolicyBody"] = d.Get("stack_policy_body")
	}
	if !d.IsNewResource() && d.HasChange("timeout_in_minutes") {
		update = true
		request["TimeoutInMinutes"] = d.Get("timeout_in_minutes")
	}
	if update {
		if _, ok := d.GetOk("replacement_option"); ok {
			request["ReplacementOption"] = d.Get("replacement_option")
		}
		if _, ok := d.GetOk("stack_policy_during_update_body"); ok {
			request["StackPolicyDuringUpdateBody"] = d.Get("stack_policy_during_update_body")
		}
		if _, ok := d.GetOk("stack_policy_during_update_url"); ok {
			request["StackPolicyDuringUpdateURL"] = d.Get("stack_policy_during_update_url")
		}
		if _, ok := d.GetOk("stack_policy_url"); ok {
			request["StackPolicyURL"] = d.Get("stack_policy_url")
		}
		if _, ok := d.GetOk("template_body"); ok {
			request["TemplateBody"] = d.Get("template_body")
		}
		if _, ok := d.GetOk("template_url"); ok {
			request["TemplateURL"] = d.Get("template_url")
		}
		if _, ok := d.GetOk("template_version"); ok {
			request["TemplateVersion"] = d.Get("template_version")
		}
		if _, ok := d.GetOkExists("use_previous_parameters"); ok {
			request["UsePreviousParameters"] = d.Get("use_previous_parameters")
		}
		action := "UpdateStack"
		conn, err := client.NewRosClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		request["ClientToken"] = buildClientToken("UpdateStack")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"UPDATE_COMPLETE"}, d.Timeout(schema.TimeoutUpdate), 100*time.Second, rosService.RosStackStateRefreshFunc(d.Id(), []string{"UPDATE_FAILED", "ROLLBACK_FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("disable_rollback")
		d.SetPartial("parameters")
		d.SetPartial("ram_role_name")
		d.SetPartial("stack_policy_body")
		d.SetPartial("timeout_in_minutes")
	}
	d.Partial(false)
	return resourceAlicloudRosStackRead(d, meta)
}
func resourceAlicloudRosStackDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	action := "DeleteStack"
	var response map[string]interface{}
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"StackId": d.Id(),
	}

	if v, ok := d.GetOk("ram_role_name"); ok {
		request["RamRoleName"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOkExists("retain_all_resources"); ok {
		request["RetainAllResources"] = v
	}
	if v, ok := d.GetOk("retain_resources"); ok {
		request["RetainResources"] = v.(*schema.Set).List()
	}
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
		if IsExpectedErrors(err, []string{"StackNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"DELETE_COMPLETE"}, d.Timeout(schema.TimeoutDelete), 100*time.Second, rosService.RosStackStateRefreshFunc(d.Id(), []string{"DELETE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

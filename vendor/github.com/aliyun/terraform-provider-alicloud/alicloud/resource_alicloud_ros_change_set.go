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

func resourceAlicloudRosChangeSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRosChangeSetCreate,
		Read:   resourceAlicloudRosChangeSetRead,
		Delete: resourceAlicloudRosChangeSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"change_set_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"change_set_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "UPDATE",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"disable_rollback": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"parameter_value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				ForceNew: true,
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"replacement_option": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"stack_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"stack_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"stack_policy_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"stack_policy_during_update_body": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"stack_policy_during_update_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"stack_policy_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"template_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"timeout_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"use_previous_parameters": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudRosChangeSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	var response map[string]interface{}
	action := "CreateChangeSet"
	request := make(map[string]interface{})
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	request["ChangeSetName"] = d.Get("change_set_name")
	if v, ok := d.GetOk("change_set_type"); ok {
		request["ChangeSetType"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOkExists("disable_rollback"); ok {
		request["DisableRollback"] = v
	}

	if v, ok := d.GetOk("notification_urls"); ok {
		request["NotificationURLs"] = v.(*schema.Set).List()
	}

	if v, ok := d.GetOk("parameters"); ok {
		request["Parameters"] = v.(*schema.Set).List()
	}

	if v, ok := d.GetOk("ram_role_name"); ok {
		request["RamRoleName"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("replacement_option"); ok {
		request["ReplacementOption"] = v
	}

	if v, ok := d.GetOk("resources_to_import"); ok {
		request["ResourcesToImport"] = v.(*schema.Set).List()
	}

	if v, ok := d.GetOk("stack_id"); ok {
		request["StackId"] = v
	}

	if v, ok := d.GetOk("stack_name"); ok {
		request["StackName"] = v
	}

	if v, ok := d.GetOk("stack_policy_body"); ok {
		request["StackPolicyBody"] = v
	}

	if v, ok := d.GetOk("stack_policy_during_update_body"); ok {
		request["StackPolicyDuringUpdateBody"] = v
	}

	if v, ok := d.GetOk("stack_policy_during_update_url"); ok {
		request["StackPolicyDuringUpdateURL"] = v
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

	if v, ok := d.GetOk("timeout_in_minutes"); ok {
		request["TimeoutInMinutes"] = v
	}

	if v, ok := d.GetOkExists("use_previous_parameters"); ok {
		request["UsePreviousParameters"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("CreateChangeSet")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ros_change_set", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(response["ChangeSetId"]))
	stateConf := BuildStateConf([]string{}, []string{"CREATE_COMPLETE"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, rosService.RosChangeSetStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudRosChangeSetRead(d, meta)
}
func resourceAlicloudRosChangeSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	object, err := rosService.DescribeRosChangeSet(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ros_change_set rosService.DescribeRosChangeSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("change_set_name", object["ChangeSetName"])
	d.Set("change_set_type", object["ChangeSetType"])
	d.Set("description", object["Description"])
	d.Set("disable_rollback", object["DisableRollback"])

	Parameters := object["Parameters"].([]interface{})
	parameters := make([]map[string]interface{}, len(Parameters))
	for i, v := range Parameters {
		parameters[i] = make(map[string]interface{})
		parameters[i]["parameter_key"] = v.(map[string]interface{})["ParameterKey"]
		parameters[i]["parameter_value"] = v.(map[string]interface{})["ParameterValue"]
	}
	if err := d.Set("parameters", parameters); err != nil {
		return WrapError(err)
	}
	d.Set("stack_id", object["StackId"])
	d.Set("stack_name", object["StackName"])
	d.Set("status", object["Status"])
	d.Set("template_body", object["TemplateBody"])
	d.Set("timeout_in_minutes", formatInt(object["TimeoutInMinutes"]))
	return nil
}
func resourceAlicloudRosChangeSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteChangeSet"
	var response map[string]interface{}
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ChangeSetId": d.Id(),
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
		if IsExpectedErrors(err, []string{"ChangeSetNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

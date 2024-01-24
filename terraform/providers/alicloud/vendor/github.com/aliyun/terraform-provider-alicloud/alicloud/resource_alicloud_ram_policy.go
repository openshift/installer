package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudRamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamPolicyCreate,
		Read:   resourceAlicloudRamPolicyRead,
		Update: resourceAlicloudRamPolicyUpdate,
		Delete: resourceAlicloudRamPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(26 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"default_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_document": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
				ConflictsWith: []string{"document", "version", "statement"},
			},
			"document": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
				Deprecated:    "Field 'document' has been deprecated from provider version 1.114.0. New field 'policy_document' instead.",
				ConflictsWith: []string{"policy_document", "version", "statement"},
			},
			"policy_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.114.0. New field 'policy_name' instead.",
				ConflictsWith: []string{"policy_name"},
			},
			"rotate_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "None",
				ValidateFunc: validation.StringInSlice([]string{"DeleteOldestNonDefaultVersionWhenLimitExceeded", "None"}, false),
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"statement": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'statement' has been deprecated from version 1.49.0, and use field 'document' to replace. ",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Allow", "Deny"}, false),
						},
						"action": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resource": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				ConflictsWith: []string{"document"},
			},
			"version": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "1",
				ConflictsWith: []string{"document"},
				// can only be '1' so far.
				ValidateFunc: validation.StringInSlice([]string{"1"}, false),
				Deprecated:   "Field 'version' has been deprecated from version 1.49.0, and use field 'document' to replace. ",
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePolicy"
	request := make(map[string]interface{})
	conn, err := client.NewRamClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("policy_document"); ok {
		request["PolicyDocument"] = v
	} else if v, ok := d.GetOk("document"); ok {
		request["PolicyDocument"] = v
	} else if v, ok := d.GetOk("statement"); ok {
		ramService := RamService{client}
		doc, err := ramService.AssemblePolicyDocument(v.(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return WrapError(err)
		}
		request["PolicyDocument"] = doc
	} else {
		return WrapError(Error("One of 'policy_document', 'document', 'statement'  must be specified."))

	}

	if v, ok := d.GetOk("policy_name"); ok {
		request["PolicyName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["PolicyName"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_policy", action, AlibabaCloudSdkGoERROR)
	}
	responsePolicy := response["Policy"].(map[string]interface{})
	d.SetId(fmt.Sprint(responsePolicy["PolicyName"]))

	return resourceAlicloudRamPolicyRead(d, meta)
}
func resourceAlicloudRamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	object, err := ramService.DescribeRamPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_policy ramService.DescribeRamPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("policy_name", d.Id())
	d.Set("name", d.Id())
	d.Set("default_version", object["Policy"].(map[string]interface{})["DefaultVersion"])
	d.Set("description", object["Policy"].(map[string]interface{})["Description"])
	d.Set("policy_document", object["DefaultPolicyVersion"].(map[string]interface{})["PolicyDocument"])
	d.Set("document", object["DefaultPolicyVersion"].(map[string]interface{})["PolicyDocument"])
	d.Set("version_id", object["DefaultPolicyVersion"].(map[string]interface{})["VersionId"])
	statement, version, err := ramService.ParsePolicyDocument(object["DefaultPolicyVersion"].(map[string]interface{})["PolicyDocument"].(string))
	if err != nil {
		return WrapError(err)
	}
	d.Set("version", version)
	d.Set("statement", statement)
	d.Set("attachment_count", object["Policy"].(map[string]interface{})["AttachmentCount"])
	d.Set("type", object["Policy"].(map[string]interface{})["PolicyType"])

	return nil
}
func resourceAlicloudRamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"PolicyName": d.Id(),
	}
	if d.HasChange("policy_document") {
		update = true
		request["PolicyDocument"] = d.Get("policy_document")
	}
	if d.HasChange("document") {
		update = true
		request["PolicyDocument"] = d.Get("document")
	}
	request["SetAsDefault"] = true
	if d.HasChange("statement") || d.HasChange("version") {
		ramService := RamService{client}
		document, err := ramService.AssemblePolicyDocument(d.Get("statement").(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return WrapError(err)
		}
		request["PolicyDocument"] = document
	}
	if update {
		if _, ok := d.GetOk("rotate_strategy"); ok {
			request["RotateStrategy"] = d.Get("rotate_strategy")
		}
		action := "CreatePolicyVersion"
		conn, err := client.NewRamClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudRamPolicyRead(d, meta)
}
func resourceAlicloudRamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePolicy"
	var response map[string]interface{}
	conn, err := client.NewRamClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PolicyName": d.Id(),
	}

	if d.Get("force").(bool) {
		listRequest := map[string]interface{}{
			"PolicyName": d.Id(),
			"PolicyType": "Custom",
		}
		listAction := "ListEntitiesForPolicy"
		response, err = conn.DoRequest(StringPointer(listAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, listRequest, &util.RuntimeOptions{})
		userResp, err := jsonpath.Get("$.Users.User", response)
		if len(userResp.([]interface{})) > 0 {
			for _, v := range userResp.([]interface{}) {
				userAction := "DetachPolicyFromUser"
				userRequest := map[string]interface{}{
					"PolicyName": d.Id(),
					"UserName":   v.(map[string]interface{})["UserName"],
					"PolicyType": "Custom",
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(userAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, userRequest, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, userRequest)
					return nil
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"EntityNotExist"}) {
						return nil
					}
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}
		groupResp, err := jsonpath.Get("$.Groups.Group", response)
		if len(groupResp.([]interface{})) > 0 {
			for _, v := range groupResp.([]interface{}) {
				groupAction := "DetachPolicyFromGroup"
				groupRequest := map[string]interface{}{
					"PolicyName": d.Id(),
					"GroupName":  v.(map[string]interface{})["GroupName"],
					"PolicyType": "Custom",
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(groupAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, groupRequest, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, groupRequest)
					return nil
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"EntityNotExist"}) {
						return nil
					}
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}
		roleResp, err := jsonpath.Get("$.Roles.Role", response)
		if len(roleResp.([]interface{})) > 0 {
			for _, v := range roleResp.([]interface{}) {
				roleAction := "DetachPolicyFromRole"
				roleRequest := map[string]interface{}{
					"PolicyName": d.Id(),
					"RoleName":   v.(map[string]interface{})["RoleName"],
					"PolicyType": "Custom",
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(roleAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, roleRequest, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, roleRequest)
					return nil
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"EntityNotExist"}) {
						return nil
					}
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		listVersionsRequest := map[string]interface{}{
			"PolicyName": d.Id(),
			"PolicyType": "Custom",
		}
		listVersionsAction := "ListPolicyVersions"
		response, err = conn.DoRequest(StringPointer(listVersionsAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, listVersionsRequest, &util.RuntimeOptions{})
		versionsResp, err := jsonpath.Get("$.PolicyVersions.PolicyVersion", response)

		// More than one means there are other versions besides the default version
		if len(versionsResp.([]interface{})) > 1 {
			for _, v := range versionsResp.([]interface{}) {
				if !v.(map[string]interface{})["IsDefaultVersion"].(bool) {
					versionAction := "DeletePolicyVersion"
					versionRequest := map[string]interface{}{
						"PolicyName": d.Id(),
						"VersionId":  v.(map[string]interface{})["VersionId"],
					}
					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
						response, err = conn.DoRequest(StringPointer(versionAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, versionRequest, &util.RuntimeOptions{})
						if err != nil {
							if NeedRetry(err) {
								wait()
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}
						addDebug(versionAction, response, versionRequest)
						return nil
					})
				}
			}
		}
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"DeleteConflict.Policy.Group", "DeleteConflict.Policy.User", "DeleteConflict.Policy.Version", "DeleteConflict.Role.Policy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

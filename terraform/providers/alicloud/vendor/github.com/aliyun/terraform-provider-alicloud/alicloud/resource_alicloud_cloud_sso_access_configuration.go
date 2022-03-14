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

func resourceAlicloudCloudSsoAccessConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudSsoAccessConfigurationCreate,
		Read:   resourceAlicloudCloudSsoAccessConfigurationRead,
		Update: resourceAlicloudCloudSsoAccessConfigurationUpdate,
		Delete: resourceAlicloudCloudSsoAccessConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_configuration_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-z0-9-]{1,32}$`), "The name of the resource. The name can be up to `32` characters long and can contain letters, digits, and hyphens (-)"),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"force_remove_permission_policies": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"permission_policies": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission_policy_document": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"permission_policy_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"permission_policy_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"System", "Inline"}, false),
						},
					},
				},
			},
			"relay_state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"session_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(900, 43200),
			},
		},
	}
}

func resourceAlicloudCloudSsoAccessConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAccessConfiguration"
	request := make(map[string]interface{})
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}

	request["AccessConfigurationName"] = d.Get("access_configuration_name")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["DirectoryId"] = d.Get("directory_id")

	if v, ok := d.GetOk("relay_state"); ok {
		request["RelayState"] = v
	}
	if v, ok := d.GetOk("session_duration"); ok {
		request["SessionDuration"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_sso_access_configuration", action, AlibabaCloudSdkGoERROR)
	}
	responseAccessConfiguration := response["AccessConfiguration"].(map[string]interface{})
	d.SetId(fmt.Sprint(request["DirectoryId"], ":", responseAccessConfiguration["AccessConfigurationId"]))

	return resourceAlicloudCloudSsoAccessConfigurationUpdate(d, meta)
}
func resourceAlicloudCloudSsoAccessConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}
	object, err := cloudssoService.DescribeCloudSsoAccessConfiguration(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_access_configuration cloudssoService.DescribeCloudSsoAccessConfiguration Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("access_configuration_id", parts[1])
	d.Set("directory_id", parts[0])
	d.Set("access_configuration_name", object["AccessConfigurationName"])
	d.Set("description", object["Description"])
	d.Set("relay_state", object["RelayState"])
	if v, ok := object["SessionDuration"]; ok && fmt.Sprint(v) != "0" {
		d.Set("session_duration", formatInt(v))
	}
	listPermissionPoliciesInAccessConfigurationObject, err := cloudssoService.ListPermissionPoliciesInAccessConfiguration(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if permissionPoliciesList, ok := listPermissionPoliciesInAccessConfigurationObject["PermissionPolicies"]; ok && permissionPoliciesList != nil {
		permissionPoliciesMaps := make([]map[string]interface{}, 0)
		for _, permissionPoliciesListItem := range permissionPoliciesList.([]interface{}) {
			if permissionPoliciesListItemMap, ok := permissionPoliciesListItem.(map[string]interface{}); ok {
				permissionPoliciesListItemMaps := map[string]interface{}{}
				permissionPoliciesListItemMaps["permission_policy_name"] = permissionPoliciesListItemMap["PermissionPolicyName"]
				permissionPoliciesListItemMaps["permission_policy_type"] = permissionPoliciesListItemMap["PermissionPolicyType"]
				permissionPoliciesListItemMaps["permission_policy_document"] = permissionPoliciesListItemMap["PermissionPolicyDocument"]
				permissionPoliciesMaps = append(permissionPoliciesMaps, permissionPoliciesListItemMaps)
			}
		}
		d.Set("permission_policies", permissionPoliciesMaps)
	}

	return nil
}
func resourceAlicloudCloudSsoAccessConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("permission_policies") {
		oraw, nraw := d.GetChange("permission_policies")
		removed := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		added := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()
		if len(removed) > 0 {
			for _, permissions := range removed {
				RemovePermissionPolicyFromAccessConfigurationReq := map[string]interface{}{
					"AccessConfigurationId": parts[1],
					"DirectoryId":           parts[0],
				}

				permissionsArg := permissions.(map[string]interface{})
				RemovePermissionPolicyFromAccessConfigurationReq["PermissionPolicyName"] = permissionsArg["permission_policy_name"]
				RemovePermissionPolicyFromAccessConfigurationReq["PermissionPolicyType"] = permissionsArg["permission_policy_type"]

				action := "RemovePermissionPolicyFromAccessConfiguration"
				conn, err := client.NewCloudssoClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, RemovePermissionPolicyFromAccessConfigurationReq, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationConflict.Task"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, RemovePermissionPolicyFromAccessConfigurationReq)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		if len(added) > 0 {
			for _, permissions := range added {
				AddPermissionPolicyFromAccessConfigurationReq := map[string]interface{}{
					"AccessConfigurationId": parts[1],
					"DirectoryId":           parts[0],
				}

				permissionsArg := permissions.(map[string]interface{})
				AddPermissionPolicyFromAccessConfigurationReq["PermissionPolicyName"] = permissionsArg["permission_policy_name"]
				AddPermissionPolicyFromAccessConfigurationReq["PermissionPolicyType"] = permissionsArg["permission_policy_type"]
				if AddPermissionPolicyFromAccessConfigurationReq["PermissionPolicyType"] == "Inline" {
					AddPermissionPolicyFromAccessConfigurationReq["InlinePolicyDocument"] = permissionsArg["permission_policy_document"]
				}

				action := "AddPermissionPolicyToAccessConfiguration"
				conn, err := client.NewCloudssoClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, AddPermissionPolicyFromAccessConfigurationReq, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationConflict.Task"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, AddPermissionPolicyFromAccessConfigurationReq)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		// Provisioning access configuration when permission policies has changed.
		cloudssoService := CloudssoService{client}
		objects, err := cloudssoService.DescribeCloudSsoAccessConfigurationProvisionings(fmt.Sprint(parts[0]), fmt.Sprint(parts[1]))
		if err != nil {
			return WrapError(err)
		}

		for _, object := range objects {
			err = cloudssoService.CloudssoServicAccessConfigurationProvisioning(fmt.Sprint(parts[0]), fmt.Sprint(parts[1]), fmt.Sprint(object["TargetType"]), fmt.Sprint(object["TargetId"]))
			if err != nil {
				return WrapError(err)
			}
		}
	}

	update := false
	updateAccessConfigurationReq := map[string]interface{}{
		"AccessConfigurationId": parts[1],
		"DirectoryId":           parts[0],
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			updateAccessConfigurationReq["NewDescription"] = v
		}
	}
	if d.HasChange("relay_state") {
		update = true
		if v, ok := d.GetOk("relay_state"); ok {
			updateAccessConfigurationReq["NewRelayState"] = v
		}
	}

	if d.HasChange("session_duration") {
		update = true
		if v, ok := d.GetOk("session_duration"); ok {
			updateAccessConfigurationReq["NewSessionDuration"] = v
		}
	}

	if update {
		action := "UpdateAccessConfiguration"
		conn, err := client.NewCloudssoClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, updateAccessConfigurationReq, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationConflict.Task"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateAccessConfigurationReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("description")
		d.SetPartial("relay_state")
		d.SetPartial("session_duration")
	}
	d.Partial(false)
	return resourceAlicloudCloudSsoAccessConfigurationRead(d, meta)
}
func resourceAlicloudCloudSsoAccessConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("permission_policies"); ok {
		removed := v.(*schema.Set).List()
		if len(removed) > 0 {
			for _, permissions := range removed {
				RemovePermissionPolicyFromAccessConfigurationReq := map[string]interface{}{
					"AccessConfigurationId": parts[1],
					"DirectoryId":           parts[0],
				}

				permissionsArg := permissions.(map[string]interface{})
				RemovePermissionPolicyFromAccessConfigurationReq["PermissionPolicyName"] = permissionsArg["permission_policy_name"]
				RemovePermissionPolicyFromAccessConfigurationReq["PermissionPolicyType"] = permissionsArg["permission_policy_type"]

				action := "RemovePermissionPolicyFromAccessConfiguration"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, RemovePermissionPolicyFromAccessConfigurationReq, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationConflict.Task"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, RemovePermissionPolicyFromAccessConfigurationReq)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}
	}

	action := "DeleteAccessConfiguration"
	request := map[string]interface{}{
		"AccessConfigurationId": parts[1],
		"DirectoryId":           parts[0],
	}

	if v, ok := d.GetOk("force_remove_permission_policies"); ok {
		request["ForceRemovePermissionPolicies"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"DeletionConflict.AccessConfiguration.Provisioning", "DeletionConflict.AccessConfiguration.AccessAssignment", "OperationConflict.Task", "DeletionConflict.AccessConfiguration.Task"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.AccessConfiguration"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

package alicloud

import (
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamRoleCreate,
		Read:   resourceAlicloudRamRoleRead,
		Update: resourceAlicloudRamRoleUpdate,
		Delete: resourceAlicloudRamRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ram_users": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:           schema.HashString,
				ConflictsWith: []string{"document"},
				Deprecated:    "Field 'ram_users' has been deprecated from version 1.49.0, and use field 'document' to replace. ",
			},
			"max_session_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
				ValidateFunc: validation.IntBetween(3600, 43200),
			},
			"services": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:           schema.HashString,
				ConflictsWith: []string{"document"},
				Deprecated:    "Field 'services' has been deprecated from version 1.49.0, and use field 'document' to replace. ",
			},
			"document": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"ram_users", "services", "version"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
				ValidateFunc: validation.ValidateJsonString,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request, err := buildAlicloudRamRoleCreateArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateRole(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_role", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ram.CreateRoleResponse)
	d.SetId(response.Role.RoleName)
	return resourceAlicloudRamRoleRead(d, meta)
}

func resourceAlicloudRamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"RoleName": d.Id(),
	}
	if d.HasChange("document") {
		update = true
		request["NewAssumeRolePolicyDocument"] = d.Get("document").(string)
	} else if d.HasChange("ram_users") || d.HasChange("services") || d.HasChange("version") {
		update = true
		document, err := ramService.AssembleRolePolicyDocument(d.Get("ram_users").(*schema.Set).List(), d.Get("services").(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return WrapError(err)
		}
		request["NewAssumeRolePolicyDocument"] = document
	}
	if d.HasChange("max_session_duration") {
		update = true
		request["NewMaxSessionDuration"] = d.Get("max_session_duration")
	}
	if d.HasChange("description") {
		update = true
		request["NewDescription"] = d.Get("description")
	}
	if update {
		action := "UpdateRole"
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

	return resourceAlicloudRamRoleRead(d, meta)
}

func resourceAlicloudRamRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}

	object, err := ramService.DescribeRamRole(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	role := object.Role
	rolePolicy, err := ramService.ParseRolePolicyDocument(role.AssumeRolePolicyDocument)
	if err != nil {
		return WrapError(err)
	}
	if len(rolePolicy.Statement) > 0 {
		principal := rolePolicy.Statement[0].Principal
		d.Set("services", principal.Service)
		d.Set("ram_users", principal.RAM)
	}
	d.Set("role_id", role.RoleId)
	d.Set("name", role.RoleName)
	d.Set("arn", role.Arn)
	d.Set("description", role.Description)
	d.Set("version", rolePolicy.Version)
	d.Set("document", role.AssumeRolePolicyDocument)
	d.Set("max_session_duration", role.MaxSessionDuration)
	return nil
}

func resourceAlicloudRamRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	ListPoliciesForRoleRequest := ram.CreateListPoliciesForRoleRequest()
	ListPoliciesForRoleRequest.RegionId = client.RegionId
	ListPoliciesForRoleRequest.RoleName = d.Id()

	if d.Get("force").(bool) {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForRole(ListPoliciesForRoleRequest)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), ListPoliciesForRoleRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(ListPoliciesForRoleRequest.GetActionName(), raw, ListPoliciesForRoleRequest.RpcRequest, ListPoliciesForRoleRequest)
		response, _ := raw.(*ram.ListPoliciesForRoleResponse)
		// Loop and remove the Policies from the Role
		if len(response.Policies.Policy) > 0 {
			for _, v := range response.Policies.Policy {
				request := ram.CreateDetachPolicyFromRoleRequest()
				request.RegionId = client.RegionId
				request.RoleName = v.PolicyName
				request.PolicyType = v.PolicyType
				request.RoleName = d.Id()
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DetachPolicyFromRole(request)
				})
				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}
	}

	deleteRoleRequest := ram.CreateDeleteRoleRequest()
	deleteRoleRequest.RegionId = client.RegionId
	deleteRoleRequest.RoleName = d.Id()
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteRole(deleteRoleRequest)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"DeleteConflict.Role.Policy"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(ListPoliciesForRoleRequest.GetActionName(), raw, ListPoliciesForRoleRequest.RpcRequest, ListPoliciesForRoleRequest)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), deleteRoleRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(ramService.WaitForRamRole(d.Id(), Deleted, DefaultTimeout))
}

func buildAlicloudRamRoleCreateArgs(d *schema.ResourceData, meta interface{}) (*ram.CreateRoleRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	request := ram.CreateCreateRoleRequest()
	request.RegionId = client.RegionId
	request.RoleName = d.Get("name").(string)

	ramUsers, usersOk := d.GetOk("ram_users")
	services, servicesOk := d.GetOk("services")
	document, documentOk := d.GetOk("document")

	if !usersOk && !servicesOk && !documentOk {
		return &ram.CreateRoleRequest{}, WrapError(Error("At least one of 'ram_users', 'services' or 'document' must be set."))
	}

	if documentOk {
		request.AssumeRolePolicyDocument = document.(string)
	} else {
		rolePolicyDocument, err := ramService.AssembleRolePolicyDocument(ramUsers.(*schema.Set).List(), services.(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return &ram.CreateRoleRequest{}, WrapError(err)
		}
		request.AssumeRolePolicyDocument = rolePolicyDocument
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}

	if v, ok := d.GetOk("max_session_duration"); ok {
		request.MaxSessionDuration = requests.NewInteger(v.(int))
	}

	return request, nil
}

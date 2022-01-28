package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudRamUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamUserCreate,
		Read:   resourceAlicloudRamUserRead,
		Update: resourceAlicloudRamUserUpdate,
		Delete: resourceAlicloudRamUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mobile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"comments": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
		},
	}
}

func resourceAlicloudRamUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}

	request := ram.CreateCreateUserRequest()
	request.RegionId = client.RegionId
	request.UserName = d.Get("name").(string)
	if v, ok := d.GetOk("display_name"); ok {
		request.DisplayName = v.(string)
	}
	if v, ok := d.GetOk("mobile"); ok {
		request.MobilePhone = v.(string)
	}
	if v, ok := d.GetOk("email"); ok {
		request.Email = v.(string)
	}
	if v, ok := d.GetOk("comments"); ok {
		request.Comments = v.(string)
	}

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateUser(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_user", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ram.CreateUserResponse)

	d.SetId(response.User.UserId)

	err = ramService.WaitForRamUser(d.Id(), Normal, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}
	return resourceAlicloudRamUserRead(d, meta)
}

func resourceAlicloudRamUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateUpdateUserRequest()
	request.RegionId = client.RegionId
	request.UserName = d.Get("name").(string)
	request.NewUserName = d.Get("name").(string)

	update := false

	if d.HasChange("name") && !d.IsNewResource() {
		ov, nv := d.GetChange("name")
		request.UserName = ov.(string)
		request.NewUserName = nv.(string)
		update = true
	}

	if d.HasChange("display_name") {
		request.NewDisplayName = d.Get("display_name").(string)
		update = true
	}

	if d.HasChange("mobile") {
		request.NewMobilePhone = d.Get("mobile").(string)
		update = true
	}

	if d.HasChange("email") {
		request.NewEmail = d.Get("email").(string)
		update = true
	}

	if d.HasChange("comments") {
		request.NewComments = d.Get("comments").(string)
		update = true
	}

	if update {

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.UpdateUser(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlicloudRamUserRead(d, meta)
}

func resourceAlicloudRamUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ramService := &RamService{client: client}
	object, err := ramService.DescribeRamUser(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.SetId(object.UserId)
	d.Set("name", object.UserName)
	d.Set("display_name", object.DisplayName)
	d.Set("mobile", object.MobilePhone)
	d.Set("email", object.Email)
	d.Set("comments", object.Comments)
	return nil
}

func resourceAlicloudRamUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ramService := &RamService{client: client}
	object, err := ramService.DescribeRamUser(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	userName := object.UserName
	request := ram.CreateListAccessKeysRequest()
	request.RegionId = client.RegionId
	request.UserName = userName

	if d.Get("force").(bool) {
		// list and delete access keys for this user
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListAccessKeys(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		listAccessKeysResponse, _ := raw.(*ram.ListAccessKeysResponse)
		if len(listAccessKeysResponse.AccessKeys.AccessKey) > 0 {
			for _, v := range listAccessKeysResponse.AccessKeys.AccessKey {
				request := ram.CreateDeleteAccessKeyRequest()
				request.RegionId = client.RegionId
				request.UserAccessKeyId = v.AccessKeyId
				request.UserName = userName
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DeleteAccessKey(request)
				})
				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}

		// list and delete policies for this user
		request := ram.CreateListPoliciesForUserRequest()
		request.RegionId = client.RegionId
		request.UserName = userName
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForUser(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		listPoliciesForUserResponse, _ := raw.(*ram.ListPoliciesForUserResponse)
		if len(listPoliciesForUserResponse.Policies.Policy) > 0 {
			for _, v := range listPoliciesForUserResponse.Policies.Policy {
				request := ram.CreateDetachPolicyFromUserRequest()
				request.RegionId = client.RegionId
				request.PolicyName = v.PolicyName
				request.PolicyType = v.PolicyType
				request.UserName = userName
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DetachPolicyFromUser(request)
				})
				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}

		// list and delete groups for this user
		listGroupsForUserRequest := ram.CreateListGroupsForUserRequest()
		listGroupsForUserRequest.RegionId = client.RegionId
		listGroupsForUserRequest.UserName = userName
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListGroupsForUser(listGroupsForUserRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), listGroupsForUserRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(listGroupsForUserRequest.GetActionName(), raw, listGroupsForUserRequest.RpcRequest, listGroupsForUserRequest)
		listGroupsForUserResponse, _ := raw.(*ram.ListGroupsForUserResponse)
		if len(listGroupsForUserResponse.Groups.Group) > 0 {
			for _, v := range listGroupsForUserResponse.Groups.Group {
				request := ram.CreateRemoveUserFromGroupRequest()
				request.RegionId = client.RegionId
				request.UserName = userName
				request.GroupName = v.GroupName
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.RemoveUserFromGroup(request)
				})
				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}

		// delete login profile for this user
		deleteLoginProfileRequest := ram.CreateDeleteLoginProfileRequest()
		deleteLoginProfileRequest.RegionId = client.RegionId
		deleteLoginProfileRequest.UserName = userName
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteLoginProfile(deleteLoginProfileRequest)
		})
		if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist.User.LoginProfile"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), deleteLoginProfileRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(deleteLoginProfileRequest.GetActionName(), raw)
		// unbind MFA device for this user
		unbindMFADeviceRequest := ram.CreateUnbindMFADeviceRequest()
		unbindMFADeviceRequest.UserName = userName
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.UnbindMFADevice(unbindMFADeviceRequest)
		})
		if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist", "EntityNotExist.User.MFADevice"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), unbindMFADeviceRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(unbindMFADeviceRequest.GetActionName(), raw, deleteLoginProfileRequest.RpcRequest, deleteLoginProfileRequest)
	}
	deleteUserRequest := ram.CreateDeleteUserRequest()
	deleteUserRequest.RegionId = client.RegionId
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		deleteUserRequest.UserName = userName
		return ramClient.DeleteUser(deleteUserRequest)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DeleteConflict.User.AccessKey", "DeleteConflict.User.Group", "DeleteConflict.User.Policy", "DeleteConflict.User.LoginProfile", "DeleteConflict.User.MFADevice"}) {
			return WrapError(Error("The user can not be deleted if he has any access keys, login profile, groups, policies, or MFA device attached. You can force the deletion of the user by setting force equals true."))
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(deleteUserRequest.GetActionName(), raw, deleteUserRequest.RpcRequest, deleteUserRequest)
	return WrapError(ramService.WaitForRamUser(d.Id(), Deleted, DefaultTimeout))
}

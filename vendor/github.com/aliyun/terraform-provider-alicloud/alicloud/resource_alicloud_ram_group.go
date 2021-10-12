package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRamGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamGroupCreate,
		Read:   resourceAlicloudRamGroupRead,
		Update: resourceAlicloudRamGroupUpdate,
		Delete: resourceAlicloudRamGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAlicloudRamGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramSercvice := RamService{client}

	request := ram.CreateCreateGroupRequest()
	request.RegionId = client.RegionId
	request.GroupName = d.Get("name").(string)
	if v, ok := d.GetOk("comments"); ok {
		request.Comments = v.(string)
	}
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ram.CreateGroupResponse)
	d.SetId(response.Group.GroupName)
	err = ramSercvice.WaitForRamGroup(d.Id(), Normal, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}
	return resourceAlicloudRamGroupRead(d, meta)
}

func resourceAlicloudRamGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateUpdateGroupRequest()
	request.RegionId = client.RegionId
	request.GroupName = d.Id()

	if d.HasChange("comments") {
		request.NewComments = d.Get("comments").(string)
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.UpdateGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlicloudRamGroupRead(d, meta)
}

func resourceAlicloudRamGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ramService := RamService{client}

	object, err := ramService.DescribeRamGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	group := object.Group
	d.Set("name", group.GroupName)
	d.Set("comments", group.Comments)
	return nil

}

func resourceAlicloudRamGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := &RamService{client}
	request := ram.CreateListUsersForGroupRequest()
	request.RegionId = client.RegionId
	request.GroupName = d.Id()

	if d.Get("force").(bool) {
		// list and delete users which in this group
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListUsersForGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		listUserResponse, _ := raw.(*ram.ListUsersForGroupResponse)
		users := listUserResponse.Users.User
		if len(users) > 0 {
			for _, v := range users {
				request := ram.CreateRemoveUserFromGroupRequest()
				request.UserName = v.UserName
				request.GroupName = d.Id()
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.RemoveUserFromGroup(request)
				})
				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}

		// list and detach policies which attach this group
		request := ram.CreateListPoliciesForGroupRequest()
		request.RegionId = client.RegionId
		request.GroupName = d.Id()
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		listPolicyResponse, _ := raw.(*ram.ListPoliciesForGroupResponse)
		policies := listPolicyResponse.Policies.Policy
		if len(policies) > 0 {
			for _, v := range policies {
				request := ram.CreateDetachPolicyFromGroupRequest()
				request.RegionId = client.RegionId
				request.PolicyType = v.PolicyType
				request.PolicyName = v.PolicyName
				request.GroupName = d.Id()
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DetachPolicyFromGroup(request)
				})
				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}
	}

	deleteGroupRequest := ram.CreateDeleteGroupRequest()
	deleteGroupRequest.RegionId = client.RegionId
	deleteGroupRequest.GroupName = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteGroup(deleteGroupRequest)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"DeleteConflict.Group.User", "DeleteConflict.Group.Policy"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(deleteGroupRequest.GetActionName(), raw, deleteGroupRequest.RpcRequest, deleteGroupRequest)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Group"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), deleteGroupRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(ramService.WaitForRamGroup(d.Id(), Deleted, DefaultTimeout))
}

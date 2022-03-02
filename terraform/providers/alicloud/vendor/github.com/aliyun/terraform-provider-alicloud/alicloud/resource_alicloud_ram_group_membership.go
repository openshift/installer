package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRamGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamGroupMembershipCreate,
		Read:   resourceAlicloudRamGroupMembershipRead,
		Update: resourceAlicloudRamGroupMembershipUpdate,
		Delete: resourceAlicloudRamGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceAlicloudRamGroupMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	group := d.Get("group_name").(string)
	users := expandStringList(d.Get("user_names").(*schema.Set).List())

	err := addUsersToGroup(client, users, group)
	if err != nil {
		return WrapError(err)
	}

	d.SetId(group)

	return resourceAlicloudRamGroupMembershipRead(d, meta)
}

func resourceAlicloudRamGroupMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	if d.HasChange("user_names") {
		d.SetPartial("user_names")
		o, n := d.GetChange("user_names")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		remove := expandStringList(oldSet.Difference(newSet).List())
		add := expandStringList(newSet.Difference(oldSet).List())
		group := d.Id()

		if err := removeUsersFromGroup(client, remove, group); err != nil {
			return WrapError(err)
		}

		if err := addUsersToGroup(client, add, group); err != nil {
			return WrapError(err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamGroupMembershipRead(d, meta)
}

func resourceAlicloudRamGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}

	object, err := ramService.DescribeRamGroupMembership(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	var users []string
	for _, v := range object.Users.User {
		users = append(users, v.UserName)
	}
	d.Set("group_name", d.Id())
	if err := d.Set("user_names", users); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAlicloudRamGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	users := expandStringList(d.Get("user_names").(*schema.Set).List())
	group := d.Id()

	if err := removeUsersFromGroup(client, users, group); err != nil {
		return WrapError(err)
	}
	return WrapError(ramService.WaitForRamGroupMembership(d.Id(), Deleted, DefaultTimeout))
}

func addUsersToGroup(client *connectivity.AliyunClient, users []string, group string) error {
	for _, u := range users {
		request := ram.CreateAddUserToGroupRequest()
		request.RegionId = client.RegionId
		request.UserName = u
		request.GroupName = group
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.AddUserToGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, u, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	}
	return nil
}

func removeUsersFromGroup(client *connectivity.AliyunClient, users []string, group string) error {
	for _, u := range users {
		request := ram.CreateRemoveUserFromGroupRequest()
		request.RegionId = client.RegionId
		request.UserName = u
		request.GroupName = group
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.RemoveUserFromGroup(request)
		})
		if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return WrapErrorf(err, DefaultErrorMsg, u, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	}
	return nil
}

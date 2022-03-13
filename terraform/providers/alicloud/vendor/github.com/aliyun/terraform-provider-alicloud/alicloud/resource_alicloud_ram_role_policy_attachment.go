package alicloud

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRamRolePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamRolePolicyAttachmentCreate,
		Read:   resourceAlicloudRamRolePolicyAttachmentRead,
		//Update: resourceAlicloudRamRolePolicyAttachmentUpdate,
		Delete: resourceAlicloudRamRolePolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"System", "Custom"}, false),
			},
		},
	}
}

func resourceAlicloudRamRolePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ram.CreateAttachPolicyToRoleRequest()
	request.RegionId = client.RegionId
	request.RoleName = d.Get("role_name").(string)
	request.PolicyType = d.Get("policy_type").(string)
	request.PolicyName = d.Get("policy_name").(string)

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.AttachPolicyToRole(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_role_policy_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	// In order to be compatible with previous Id (before 1.9.6) which format to role:<policy_name>:<policy_type>:<role_name>
	d.SetId(strings.Join([]string{"role", request.PolicyName, request.PolicyType, request.RoleName}, COLON_SEPARATED))

	return resourceAlicloudRamRolePolicyAttachmentRead(d, meta)
}

func resourceAlicloudRamRolePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}

	if split := strings.Split(d.Id(), ":"); len(split) != 4 {
		id := strings.Join([]string{"role", d.Get("policy_name").(string), d.Get("policy_type").(string), d.Get("role_name").(string)}, COLON_SEPARATED)
		d.SetId(id)
	}
	object, err := ramService.DescribeRamRolePolicyAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	d.Set("role_name", parts[3])
	d.Set("policy_name", object.PolicyName)
	d.Set("policy_type", object.PolicyType)
	return nil
}

func resourceAlicloudRamRolePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	id := strings.Join([]string{"role", d.Get("policy_name").(string), d.Get("policy_type").(string), d.Get("role_name").(string)}, COLON_SEPARATED)

	if d.Id() != id {
		d.SetId(id)
	}

	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return WrapError(err)
	}
	request := ram.CreateDetachPolicyFromRoleRequest()
	request.RegionId = client.RegionId
	request.PolicyName = parts[1]
	request.PolicyType = parts[2]
	request.RoleName = parts[3]

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.DetachPolicyFromRole(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(ramService.WaitForRamRolePolicyAttachment(d.Id(), Deleted, DefaultTimeout))
}

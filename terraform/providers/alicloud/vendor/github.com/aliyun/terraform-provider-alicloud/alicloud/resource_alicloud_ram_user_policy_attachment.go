package alicloud

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRamUserPolicyAtatchment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamUserPolicyAttachmentCreate,
		Read:   resourceAlicloudRamUserPolicyAttachmentRead,
		Delete: resourceAlicloudRamUserPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_name": {
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

func resourceAlicloudRamUserPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateAttachPolicyToUserRequest()
	request.RegionId = client.RegionId
	request.UserName = d.Get("user_name").(string)
	request.PolicyName = d.Get("policy_name").(string)
	request.PolicyType = d.Get("policy_type").(string)

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.AttachPolicyToUser(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_user_policy_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(strings.Join([]string{"user", request.PolicyName, request.PolicyType, request.UserName}, COLON_SEPARATED))
	return resourceAlicloudRamUserPolicyAttachmentRead(d, meta)
}

func resourceAlicloudRamUserPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	// In order to be compatible with previous Id (before 1.9.6) which format to user:<policy_name>:<policy_type>:<user_name>

	if split := strings.Split(d.Id(), ":"); len(split) != 4 {
		id := strings.Join([]string{"user", d.Get("policy_name").(string), d.Get("policy_type").(string), d.Get("user_name").(string)}, COLON_SEPARATED)
		d.SetId(id)
	}

	object, err := ramService.DescribeRamUserPolicyAttachment(d.Id())
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
	d.Set("user_name", parts[3])
	d.Set("policy_name", object.PolicyName)
	d.Set("policy_type", object.PolicyType)
	return nil
}

func resourceAlicloudRamUserPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}

	// In order to be compatible with previous Id (before 1.9.6) which format to user:<policy_name>:<policy_type>:<user_name>
	id := strings.Join([]string{"user", d.Get("policy_name").(string), d.Get("policy_type").(string), d.Get("user_name").(string)}, COLON_SEPARATED)

	if d.Id() != id {
		d.SetId(id)
	}
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return WrapError(err)
	}
	request := ram.CreateDetachPolicyFromUserRequest()
	request.RegionId = client.RegionId
	request.PolicyName = parts[1]
	request.PolicyType = parts[2]
	request.UserName = parts[3]

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.DetachPolicyFromUser(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(ramService.WaitForRamUserPolicyAttachment(d.Id(), Deleted, DefaultTimeout))

}

package alicloud

import (
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSecurityGroupCreate,
		Read:   resourceAliyunSecurityGroupRead,
		Update: resourceAliyunSecurityGroupUpdate,
		Delete: resourceAliyunSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "normal",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"normal", "enterprise"}, false),
			},
			"inner_access": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'inner_access' has been deprecated from provider version 1.55.3. Use 'inner_access_policy' replaces it.",
			},
			"inner_access_policy": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"inner_access"},
				ValidateFunc:  validation.StringInSlice([]string{"Accept", "Drop"}, false),
				// The InnerAccessPolicy attribute of enterprise level security group can't be modified.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("security_group_type").(string) == "enterprise"
				},
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ecs.CreateCreateSecurityGroupRequest()
	request.RegionId = client.RegionId

	if v := d.Get("name").(string); v != "" {
		request.SecurityGroupName = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	request.SecurityGroupType = d.Get("security_group_type").(string)

	if v := d.Get("vpc_id").(string); v != "" {
		request.VpcId = v
	}

	request.ResourceGroupId = d.Get("resource_group_id").(string)

	request.ClientToken = buildClientToken(request.GetActionName())

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateSecurityGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_security_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.CreateSecurityGroupResponse)
	d.SetId(response.SecurityGroupId)
	return resourceAliyunSecurityGroupUpdate(d, meta)
}

func resourceAliyunSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeSecurityGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_security_group ecsService.DescribeSecurityGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
	}

	d.Set("name", object.SecurityGroupName)
	d.Set("description", object.Description)
	d.Set("vpc_id", object.VpcId)
	d.Set("inner_access", object.InnerAccessPolicy == string(GroupInnerAccept))
	d.Set("inner_access_policy", object.InnerAccessPolicy)

	request := ecs.CreateDescribeSecurityGroupsRequest()
	request.RegionId = client.RegionId
	request.SecurityGroupId = d.Id()

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroups(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeSecurityGroupsResponse)
	if len(response.SecurityGroups.SecurityGroup) < 1 {
		return WrapErrorf(Error(GetNotFoundMessage("SecurityGroup", d.Id())), NotFoundMsg, ProviderERROR)
	}
	d.Set("security_group_type", response.SecurityGroups.SecurityGroup[0].SecurityGroupType)
	d.Set("resource_group_id", response.SecurityGroups.SecurityGroup[0].ResourceGroupId)

	tags, err := ecsService.ListTagResources(d.Id(), "securitygroup")
	if err != nil {
		return WrapError(err)
	} else {
		d.Set("tags", tagsToMap(tags))
	}

	return nil
}

func resourceAliyunSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		action := "JoinResourceGroup"
		request := map[string]interface{}{
			"ResourceType":    "securitygroup",
			"ResourceId":      d.Id(),
			"RegionId":        client.RegionId,
			"ResourceGroupId": d.Get("resource_group_id"),
		}
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		d.SetPartial("resource_group_id")
	}
	if err := setTags(client, TagResourceSecurityGroup, d); err != nil {
		return WrapError(err)
	} else {
		d.SetPartial("tags")
	}

	if d.HasChange("inner_access_policy") || d.HasChange("inner_access") || d.IsNewResource() && d.Get("security_group_type").(string) != "enterprise" {
		policy := GroupInnerAccept
		if v, ok := d.GetOk("inner_access_policy"); ok && v.(string) != "" {
			policy = GroupInnerAccessPolicy(v.(string))
		} else if v, ok := d.GetOkExists("inner_access"); ok && !v.(bool) {
			policy = GroupInnerDrop
		}

		request := ecs.CreateModifySecurityGroupPolicyRequest()
		request.RegionId = client.RegionId
		request.SecurityGroupId = d.Id()
		request.InnerAccessPolicy = string(policy)
		request.ClientToken = buildClientToken(request.GetActionName())

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifySecurityGroupPolicy(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("inner_access")
		d.SetPartial("inner_access_policy")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliyunSecurityGroupRead(d, meta)
	}

	update := false
	request := ecs.CreateModifySecurityGroupAttributeRequest()
	request.RegionId = client.RegionId
	request.SecurityGroupId = d.Id()
	if d.HasChange("name") {
		request.SecurityGroupName = d.Get("name").(string)
		update = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}
	if update {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifySecurityGroupAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("name")
		d.SetPartial("description")
	}

	d.Partial(false)

	return resourceAliyunSecurityGroupRead(d, meta)
}

func resourceAliyunSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	request := ecs.CreateDeleteSecurityGroupRequest()
	request.RegionId = client.RegionId
	request.SecurityGroupId = d.Id()

	err := resource.Retry(6*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteSecurityGroup(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"DependencyViolation"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR)
	}
	return WrapError(ecsService.WaitForSecurityGroup(d.Id(), Deleted, DefaultTimeoutMedium))

}

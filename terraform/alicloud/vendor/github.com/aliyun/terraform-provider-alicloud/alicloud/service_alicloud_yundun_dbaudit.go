package alicloud

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_dbaudit"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

type DbauditService struct {
	client *connectivity.AliyunClient
}

type PolicyRequired struct {
	PolicyName string
	PolicyType string
}

const (
	RoleDefaultDescription = "DbAudit perform the default role to access your other cloud resources"
	RoleName               = "AliyunDbAuditDefaultRole"
	AssumeRolePolicy       = `{
		"Statement": [
			{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": {
					"Service": [
						"dbaudit.aliyuncs.com"
					]
				}
			}
		],
		"Version": "1"
	}`
	DBauditResourceType = "INSTANCE"
)

var policyRequired = []PolicyRequired{
	{
		PolicyName: "AliyunDbAuditRolePolicy",
		PolicyType: "System",
	},
	{
		PolicyName: "AliyunLogFullAccess",
		PolicyType: "System",
	},
}

func (s *DbauditService) DescribeYundunDbauditInstance(id string) (v yundun_dbaudit.Instance, err error) {
	request := yundun_dbaudit.CreateDescribeInstancesRequest()
	var instanceIds []string
	instanceIds = append(instanceIds, id)
	request.InstanceId = &instanceIds
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.CurrentPage = requests.NewInteger(1)
	raw, err := s.client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
		return dbauditClient.DescribeInstances(request)
	})
	if err != nil {
		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*yundun_dbaudit.DescribeInstancesResponse)

	if len(response.Instances) == 0 || response.Instances[0].InstanceId != id {
		return v, WrapErrorf(Error(GetNotFoundMessage("Yundun_dbaudit Instance", id)), NotFoundMsg, ProviderERROR)
	}
	v = response.Instances[0]
	return
}

func (s *DbauditService) DescribeDbauditInstanceAttribute(id string) (v yundun_dbaudit.InstanceAttribute, err error) {
	request := yundun_dbaudit.CreateDescribeInstanceAttributeRequest()
	request.InstanceId = id

	raw, err := s.client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
		return dbauditClient.DescribeInstanceAttribute(request)
	})

	if err != nil {
		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*yundun_dbaudit.DescribeInstanceAttributeResponse)
	if response.InstanceAttribute.InstanceId != id {
		return v, WrapErrorf(Error(GetNotFoundMessage("Yundun_dbaudit Instance", id)), NotFoundMsg, ProviderERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	v = response.InstanceAttribute
	return v, WrapError(err)
}

func (s *DbauditService) StartDbauditInstance(instanceId string, vSwitchId string) error {
	request := yundun_dbaudit.CreateStartInstanceRequest()
	request.InstanceId = instanceId
	request.VswitchId = vSwitchId
	raw, err := s.client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
		return dbauditClient.StartInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *DbauditService) UpdateDbauditInstanceDescription(instanceId string, description string) error {
	request := yundun_dbaudit.CreateModifyInstanceAttributeRequest()
	request.InstanceId = instanceId
	request.Description = description
	raw, err := s.client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
		return dbauditClient.ModifyInstanceAttribute(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *DbauditService) UpdateInstanceSpec(schemaName string, specName string, d *schema.ResourceData, meta interface{}) error {
	request := bssopenapi.CreateModifyInstanceRequest()
	request.InstanceId = d.Id()

	request.ProductCode = "dbaudit"
	request.SubscriptionType = "Subscription"
	// only support upgrade
	request.ModifyType = "Upgrade"

	request.Parameter = &[]bssopenapi.ModifyInstanceParameter{
		{
			Code:  specName,
			Value: d.Get(schemaName).(string),
		},
	}

	request.RegionId = string(connectivity.Hangzhou)
	var response *bssopenapi.ModifyInstanceResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithBssopenapiClient(func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
			return bssopenapiClient.ModifyInstance(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				request.RegionId = string(connectivity.APSouthEast1)
				request.Domain = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*bssopenapi.ModifyInstanceResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if !response.Success {
		return WrapError(Error(response.Message))
	}
	return nil
}

func (s *DbauditService) DbauditInstanceRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeYundunDbauditInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil if nothing matched
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.InstanceStatus == failState {
				return object, object.InstanceStatus, WrapError(Error(FailedToReachTargetStatus, object.InstanceStatus))
			}
		}

		return object, object.InstanceStatus, nil
	}
}

func (s *DbauditService) createRole() error {
	createRoleRequest := ram.CreateCreateRoleRequest()
	createRoleRequest.RoleName = RoleName
	createRoleRequest.Description = RoleDefaultDescription
	createRoleRequest.AssumeRolePolicyDocument = AssumeRolePolicy
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateRole(createRoleRequest)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, RoleName, createRoleRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(createRoleRequest.GetActionName(), raw, createRoleRequest.RpcRequest, createRoleRequest)
	return nil
}

func (s *DbauditService) attachPolicy(policyToBeAttached []PolicyRequired) error {
	log.Printf("DEBUG attachPolicy policyRequred %v", policyToBeAttached)
	attachPolicyRequest := ram.CreateAttachPolicyToRoleRequest()
	for _, policy := range policyToBeAttached {
		log.Printf("DEBUG attach Policy in policyRequred %v", policy)
		attachPolicyRequest.RoleName = RoleName
		attachPolicyRequest.PolicyName = policy.PolicyName
		attachPolicyRequest.PolicyType = policy.PolicyType
		raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.AttachPolicyToRole(attachPolicyRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, RoleName, attachPolicyRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		if response, err := raw.(*ram.AttachPolicyToRoleResponse); !err || !response.IsSuccess() {
			log.Printf("AttachPolicyToRoleError [%v]", response)
			return errors.New("attach policy to role failed")
		}
		addDebug(attachPolicyRequest.GetActionName(), raw, attachPolicyRequest.RpcRequest, attachPolicyRequest)

	}
	return nil
}

func (s *DbauditService) ProcessRolePolicy() error {
	getRoleRequest := ram.CreateGetRoleRequest()
	getRoleRequest.RoleName = RoleName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetRole(getRoleRequest)
	})
	log.Printf("DEBUG ProcessRolePolicy create role %v", raw.(*ram.GetRoleResponse))
	addDebug(getRoleRequest.GetActionName(), raw, getRoleRequest.RpcRequest, getRoleRequest)
	response, _ := raw.(*ram.GetRoleResponse)
	if err != nil || response == nil || response.Role.RoleName != RoleName {
		if err := s.createRole(); err != nil {
			return WrapError(err)
		}
	}
	listPolicyForRoleRequest := ram.CreateListPoliciesForRoleRequest()
	listPolicyForRoleRequest.RoleName = RoleName
	raw, err = s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForRole(listPolicyForRoleRequest)
	})
	addDebug(listPolicyForRoleRequest.GetActionName(), raw, listPolicyForRoleRequest.RpcRequest, listPolicyForRoleRequest)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, RoleName, listPolicyForRoleRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	var policyToAttach []PolicyRequired
	if response, _ := raw.(*ram.ListPoliciesForRoleResponse); response != nil && response.IsSuccess() {
		for _, required := range policyRequired {
			contains := false
			for _, policy := range response.Policies.Policy {
				if required.PolicyName == policy.PolicyName {
					contains = true
				}
			}
			if !contains {
				policyToAttach = append(policyToAttach, required)
			}
		}
	}

	if policyToAttach != nil && len(policyToAttach) > 0 {
		return s.attachPolicy(policyToAttach)
	}

	return nil
}

func (s *DbauditService) WaitForYundunDbauditInstance(instanceId string, status Status, timeoutSenconds time.Duration) error {
	deadline := time.Now().Add(timeoutSenconds * time.Second)
	for {
		_, err := s.DescribeYundunDbauditInstance(instanceId)

		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, instanceId, GetFunc(1), timeoutSenconds, "", "", ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *DbauditService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType TagResourceType) (tags []yundun_dbaudit.TagResource, err error) {
	request := yundun_dbaudit.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = strings.ToUpper(string(resourceType))
	request.ResourceId = &[]string{resourceId}
	if resourceTags != nil && len(resourceTags) > 0 {
		var reqTags []yundun_dbaudit.ListTagResourcesTag
		for key, value := range resourceTags {
			reqTags = append(reqTags, yundun_dbaudit.ListTagResourcesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}

	var raw interface{}
	raw, err = s.client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
		return dbauditClient.ListTagResources(request)
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	response, _ := raw.(*yundun_dbaudit.ListTagResourcesResponse)

	return response.TagResources, nil
}

func (s *DbauditService) tagsToMap(tags []yundun_dbaudit.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *DbauditService) ignoreTag(t yundun_dbaudit.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *DbauditService) setInstanceTags(d *schema.ResourceData, resourceType TagResourceType) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

		if len(remove) > 0 {
			var tagKey []string
			for _, v := range remove {
				tagKey = append(tagKey, v.Key)
			}
			request := yundun_dbaudit.CreateUntagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = strings.ToUpper(string(resourceType))
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithDbauditClient(func(client *yundun_dbaudit.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := yundun_dbaudit.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = strings.ToUpper(string(resourceType))
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithDbauditClient(func(client *yundun_dbaudit.Client) (interface{}, error) {
				return client.TagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

	}

	return nil
}

func (s *DbauditService) diffTags(oldTags, newTags []yundun_dbaudit.TagResourcesTag) ([]yundun_dbaudit.TagResourcesTag, []yundun_dbaudit.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []yundun_dbaudit.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *DbauditService) tagsFromMap(m map[string]interface{}) []yundun_dbaudit.TagResourcesTag {
	result := make([]yundun_dbaudit.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, yundun_dbaudit.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *DbauditService) UpdateResourceGroup(resourceId, resourceGroupId string) error {
	request := yundun_dbaudit.CreateMoveResourceGroupRequest()
	request.RegionId = s.client.RegionId
	request.ResourceId = resourceId
	request.ResourceType = DBauditResourceType
	request.ResourceGroupId = resourceGroupId
	raw, err := s.client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
		return dbauditClient.MoveResourceGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

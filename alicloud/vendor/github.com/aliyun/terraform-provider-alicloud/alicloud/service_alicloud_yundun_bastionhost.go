package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_bastionhost"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

type YundunBastionhostService struct {
	client *connectivity.AliyunClient
}

type BastionhostPolicyRequired struct {
	PolicyName string
	PolicyType string
}

const (
	BastionhostRoleDefaultDescription = "Bastionhost will access other cloud resources by playing this role by default"
	BastionhostRoleName               = "AliyunBastionHostDefaultRole"
	BastionhostAssumeRolePolicy       = `{
		"Statement": [
			{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": {
					"Service": [
						"bastionhost.aliyuncs.com"
					]
				}
			}
		],
		"Version": "1"
	}`
	BastionhostResourceType = "INSTANCE"
)

var bastionhostpolicyRequired = []BastionhostPolicyRequired{
	{
		PolicyName: "AliyunBastionHostRolePolicy",
		PolicyType: "System",
	},
}

func (s *YundunBastionhostService) DescribeBastionhostInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeInstanceAttribute"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.InstanceAttribute", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.InstanceAttribute", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *YundunBastionhostService) StartBastionhostInstance(instanceId string, vSwitchId string, securityGroupIds []string) error {
	request := yundun_bastionhost.CreateStartInstanceRequest()
	request.InstanceId = instanceId
	request.VswitchId = vSwitchId
	request.SecurityGroupIds = &securityGroupIds
	raw, err := s.client.WithBastionhostClient(func(BastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
		return BastionhostClient.StartInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *YundunBastionhostService) UpdateBastionhostInstanceDescription(instanceId string, description string) error {
	request := yundun_bastionhost.CreateModifyInstanceAttributeRequest()
	request.InstanceId = instanceId
	request.Description = description
	raw, err := s.client.WithBastionhostClient(func(BastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
		return BastionhostClient.ModifyInstanceAttribute(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *YundunBastionhostService) UpdateBastionhostSecurityGroups(instanceId string, securityGroups []string) error {
	request := yundun_bastionhost.CreateConfigInstanceSecurityGroupsRequest()
	request.InstanceId = instanceId
	request.SecurityGroupIds = &securityGroups
	raw, err := s.client.WithBastionhostClient(func(BastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
		return BastionhostClient.ConfigInstanceSecurityGroups(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *YundunBastionhostService) UpdateInstanceSpec(schemaSpecMap map[string]string, d *schema.ResourceData, meta interface{}) error {
	request := bssopenapi.CreateModifyInstanceRequest()
	request.InstanceId = d.Id()

	request.ProductCode = "bastionhost"
	request.SubscriptionType = "Subscription"
	// only support upgrade
	request.ModifyType = "Upgrade"

	params := make([]bssopenapi.ModifyInstanceParameter, 0, len(schemaSpecMap))
	for schemaName, spec := range schemaSpecMap {
		params = append(params, bssopenapi.ModifyInstanceParameter{schemaName, d.Get(spec).(string)})
	}

	request.Parameter = &params
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

func (s *YundunBastionhostService) BastionhostInstanceRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeBastionhostInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil if nothing matched
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["InstanceStatus"]) == failState {
				return object, fmt.Sprint(object["InstanceStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["InstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["InstanceStatus"]), nil
	}
}

func (s *YundunBastionhostService) createRole() error {
	createRoleRequest := ram.CreateCreateRoleRequest()
	createRoleRequest.RoleName = BastionhostRoleName
	createRoleRequest.Description = BastionhostRoleDefaultDescription
	createRoleRequest.AssumeRolePolicyDocument = BastionhostAssumeRolePolicy
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateRole(createRoleRequest)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, BastionhostRoleName, createRoleRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(createRoleRequest.GetActionName(), raw, createRoleRequest.RpcRequest, createRoleRequest)
	return nil
}

func (s *YundunBastionhostService) attachPolicy(policyToBeAttached []BastionhostPolicyRequired) error {
	attachPolicyRequest := ram.CreateAttachPolicyToRoleRequest()
	for _, policy := range policyToBeAttached {
		attachPolicyRequest.RoleName = BastionhostRoleName
		attachPolicyRequest.PolicyName = policy.PolicyName
		attachPolicyRequest.PolicyType = policy.PolicyType
		raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.AttachPolicyToRole(attachPolicyRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, BastionhostRoleName, attachPolicyRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		if response, err := raw.(*ram.AttachPolicyToRoleResponse); !err || !response.IsSuccess() {
			return WrapError(errors.New("attach policy to role failed"))
		}
		addDebug(attachPolicyRequest.GetActionName(), raw, attachPolicyRequest.RpcRequest, attachPolicyRequest)

	}
	return nil
}

func (s *YundunBastionhostService) ProcessRolePolicy() error {
	getRoleRequest := ram.CreateGetRoleRequest()
	getRoleRequest.RoleName = BastionhostRoleName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetRole(getRoleRequest)
	})
	response, _ := raw.(*ram.GetRoleResponse)
	if err != nil || response == nil || response.Role.RoleName != BastionhostRoleName {
		if err := s.createRole(); err != nil {
			return err
		}
	}
	addDebug(getRoleRequest.GetActionName(), raw, getRoleRequest.RpcRequest, getRoleRequest)
	listPolicyForRoleRequest := ram.CreateListPoliciesForRoleRequest()
	listPolicyForRoleRequest.RoleName = BastionhostRoleName
	raw, err = s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForRole(listPolicyForRoleRequest)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, BastionhostRoleName, listPolicyForRoleRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(listPolicyForRoleRequest.GetActionName(), raw, listPolicyForRoleRequest.RpcRequest, listPolicyForRoleRequest)
	var policyToAttach []BastionhostPolicyRequired
	if response, _ := raw.(*ram.ListPoliciesForRoleResponse); response != nil && response.IsSuccess() {
		for _, required := range bastionhostpolicyRequired {
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

func (s *YundunBastionhostService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType TagResourceType) (tags []yundun_bastionhost.TagResource, err error) {
	request := yundun_bastionhost.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = strings.ToUpper(string(resourceType))
	request.ResourceId = &[]string{resourceId}
	if resourceTags != nil && len(resourceTags) > 0 {
		var reqTags []yundun_bastionhost.ListTagResourcesTag
		for key, value := range resourceTags {
			reqTags = append(reqTags, yundun_bastionhost.ListTagResourcesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}

	var raw interface{}

	raw, err = s.client.WithBastionhostClient(func(client *yundun_bastionhost.Client) (interface{}, error) {
		return client.ListTagResources(request)
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	response, _ := raw.(*yundun_bastionhost.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

func (s *YundunBastionhostService) tagsToMap(tags []yundun_bastionhost.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *YundunBastionhostService) ignoreTag(t yundun_bastionhost.TagResource) bool {
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

func (s *YundunBastionhostService) setInstanceTags(d *schema.ResourceData, resourceType TagResourceType) error {
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
			request := yundun_bastionhost.CreateUntagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = strings.ToUpper(string(resourceType))
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithBastionhostClient(func(client *yundun_bastionhost.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := yundun_bastionhost.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = strings.ToUpper(string(resourceType))
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithBastionhostClient(func(client *yundun_bastionhost.Client) (interface{}, error) {
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

func (s *YundunBastionhostService) diffTags(oldTags, newTags []yundun_bastionhost.TagResourcesTag) ([]yundun_bastionhost.TagResourcesTag, []yundun_bastionhost.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []yundun_bastionhost.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *YundunBastionhostService) tagsFromMap(m map[string]interface{}) []yundun_bastionhost.TagResourcesTag {
	result := make([]yundun_bastionhost.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, yundun_bastionhost.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *YundunBastionhostService) UpdateResourceGroup(resourceId, resourceGroupId string) error {
	request := yundun_bastionhost.CreateMoveResourceGroupRequest()
	request.RegionId = s.client.RegionId
	request.ResourceId = resourceId
	request.ResourceType = BastionhostResourceType
	request.ResourceGroupId = resourceGroupId
	raw, err := s.client.WithBastionhostClient(func(BastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
		return BastionhostClient.MoveResourceGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *YundunBastionhostService) DescribeBastionhostUserGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetUserGroup"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"InstanceId":  parts[0],
		"UserGroupId": parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost:UserGroup", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.UserGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.UserGroup", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *YundunBastionhostService) DescribeBastionhostUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetUser"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": parts[0],
		"UserId":     parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost:User", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.User", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.User", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *YundunBastionhostService) DescribeBastionhostHostGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetHostGroup"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"HostGroupId": parts[1],
		"InstanceId":  parts[0],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost:HostGroup", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HostGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostGroup", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *YundunBastionhostService) DescribeBastionhostUserAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListUsers"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"InstanceId":  parts[0],
		"UserGroupId": parts[1],
		"PageNumber":  1,
		"PageSize":    50,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
			if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus"}) {
				return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost:UserAttachment", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Users", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Users", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["UserId"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *YundunBastionhostService) DescribeBastionhostHost(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetHost"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"HostId":     parts[1],
		"InstanceId": parts[0],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(false)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost:Host", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Host", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Host", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *YundunBastionhostService) DescribeBastionhostHostAccount(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetHostAccount"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"HostAccountId": parts[1],
		"InstanceId":    parts[0],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost:HostAccount", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HostAccount", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostAccount", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *YundunBastionhostService) DescribeBastionhostHostAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListHosts"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"HostGroupId": parts[1],
		"InstanceId":  parts[0],
		"PageNumber":  1,
		"PageSize":    0,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
			if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus"}) {
				return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost:HostAttachment", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Hosts", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Hosts", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["HostId"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost", id)), NotFoundWithResponse, response)
	}
	return
}
func (s *YundunBastionhostService) DescribeBastionhostHostAccountUserAttachment(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListHostAccountsForUser"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"HostId":     parts[2],
		"InstanceId": parts[0],
		"UserId":     parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(10*time.Second, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost:HostAccountUserAttachment", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HostAccounts", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostAccounts", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost", id)), NotFoundWithResponse, response)
	}
	return v.([]interface{}), nil
}

func (s *YundunBastionhostService) DescribeBastionhostHostAccountUserGroupAttachment(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListHostAccountsForUserGroup"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"HostId":      parts[2],
		"InstanceId":  parts[0],
		"UserGroupId": parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(10*time.Second, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost:HostAccountUserGroupAttachment", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HostAccounts", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostAccounts", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost", id)), NotFoundWithResponse, response)
	}
	return v.([]interface{}), nil
}

func (s *YundunBastionhostService) DescribeBastionhostHostGroupAccountUserAttachment(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListHostGroupAccountNamesForUser"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"HostGroupId": parts[2],
		"InstanceId":  parts[0],
		"UserId":      parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(10*time.Second, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost:HostGroupAccountUserAttachment", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HostAccountNames", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostAccountNames", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost", id)), NotFoundWithResponse, response)
	}
	return v.([]interface{}), nil
}

func (s *YundunBastionhostService) DescribeBastionhostHostGroupAccountUserGroupAttachment(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListHostGroupAccountNamesForUserGroup"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"HostGroupId": parts[2],
		"InstanceId":  parts[0],
		"UserGroupId": parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(10*time.Second, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost:HostGroupAccountUserGroupAttachment", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HostAccountNames", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostAccountNames", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Bastionhost", id)), NotFoundWithResponse, response)
	}
	return v.([]interface{}), nil
}

func (s *YundunBastionhostService) EnableInstancePublicAccess(id string) (err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	action := "EnableInstancePublicAccess"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func (s *YundunBastionhostService) DisableInstancePublicAccess(id string) (err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	action := "DisableInstancePublicAccess"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

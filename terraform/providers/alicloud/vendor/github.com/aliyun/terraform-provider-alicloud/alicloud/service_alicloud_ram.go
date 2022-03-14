package alicloud

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type Effect string

const (
	Allow Effect = "Allow"
	Deny  Effect = "Deny"
)

type Principal struct {
	Service []string
	RAM     []string
}

type RolePolicyStatement struct {
	Effect    Effect
	Action    string
	Principal Principal
}

type RolePolicy struct {
	Statement []RolePolicyStatement
	Version   string
}

type PolicyStatement struct {
	Effect   Effect
	Action   interface{}
	Resource interface{}
}

type Policy struct {
	Statement []PolicyStatement
	Version   string
}

type RamService struct {
	client *connectivity.AliyunClient
}

func (s *RamService) ParseRolePolicyDocument(policyDocument string) (RolePolicy, error) {
	var policy RolePolicy
	err := json.Unmarshal([]byte(policyDocument), &policy)
	if err != nil {
		return RolePolicy{}, WrapError(err)
	}
	return policy, nil
}

func (s *RamService) ParsePolicyDocument(policyDocument string) (statement []map[string]interface{}, version string, err error) {
	policy := Policy{}
	err = json.Unmarshal([]byte(policyDocument), &policy)
	if err != nil {
		err = WrapError(err)
		return
	}

	version = policy.Version
	statement = make([]map[string]interface{}, 0, len(policy.Statement))
	for _, v := range policy.Statement {
		item := make(map[string]interface{})

		item["effect"] = v.Effect
		if val, ok := v.Action.([]interface{}); ok {
			item["action"] = val
		} else {
			item["action"] = []interface{}{v.Action}
		}

		if val, ok := v.Resource.([]interface{}); ok {
			item["resource"] = val
		} else {
			item["resource"] = []interface{}{v.Resource}
		}
		statement = append(statement, item)
	}
	return
}

func (s *RamService) AssembleRolePolicyDocument(ramUser, service []interface{}, version string) (string, error) {
	services := expandStringList(service)
	users := expandStringList(ramUser)

	statement := RolePolicyStatement{
		Effect: Allow,
		Action: "sts:AssumeRole",
		Principal: Principal{
			RAM:     users,
			Service: services,
		},
	}

	policy := RolePolicy{
		Version:   version,
		Statement: []RolePolicyStatement{statement},
	}

	data, err := json.Marshal(policy)
	if err != nil {
		return "", WrapError(err)
	}
	return string(data), nil
}

func (s *RamService) AssemblePolicyDocument(document []interface{}, version string) (string, error) {
	var statements []PolicyStatement

	for _, v := range document {
		doc := v.(map[string]interface{})

		actions := expandStringList(doc["action"].([]interface{}))
		resources := expandStringList(doc["resource"].([]interface{}))

		statement := PolicyStatement{
			Effect:   Effect(doc["effect"].(string)),
			Action:   actions,
			Resource: resources,
		}
		statements = append(statements, statement)
	}

	policy := Policy{
		Version:   version,
		Statement: statements,
	}

	data, err := json.Marshal(policy)
	if err != nil {
		return "", WrapError(err)
	}
	return string(data), nil
}

// Judge whether the role policy contains service "ecs.aliyuncs.com"
func (s *RamService) JudgeRolePolicyPrincipal(roleName string) error {
	request := ram.CreateGetRoleRequest()
	request.RegionId = s.client.RegionId
	request.RoleName = roleName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetRole(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, roleName, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ram.GetRoleResponse)
	policy, err := s.ParseRolePolicyDocument(resp.Role.AssumeRolePolicyDocument)
	if err != nil {
		return WrapError(err)
	}
	for _, v := range policy.Statement {
		for _, val := range v.Principal.Service {
			if strings.Trim(val, " ") == "ecs.aliyuncs.com" {
				return nil
			}
		}
	}
	return WrapError(fmt.Errorf("Role policy services must contains 'ecs.aliyuncs.com', Now is \n%v.", resp.Role.AssumeRolePolicyDocument))
}

func (s *RamService) GetIntersection(dataMap []map[string]interface{}, allDataMap map[string]interface{}) (allData []interface{}) {
	for _, v := range dataMap {
		if len(v) > 0 {
			for key := range allDataMap {
				if _, ok := v[key]; !ok {
					allDataMap[key] = nil
				}
			}
		}
	}

	for _, v := range allDataMap {
		if v != nil {
			allData = append(allData, v)
		}
	}
	return
}

func (s *RamService) DescribeRamUser(id string) (*ram.UserInGetUser, error) {
	user := &ram.UserInGetUser{}
	listUsersRequest := ram.CreateListUsersRequest()
	listUsersRequest.RegionId = s.client.RegionId
	listUsersRequest.MaxItems = requests.NewInteger(100)
	var userName string

	for {
		raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListUsers(listUsersRequest)
		})
		if err != nil {
			return user, WrapErrorf(err, DefaultErrorMsg, id, listUsersRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(listUsersRequest.GetActionName(), raw, listUsersRequest.RegionId, listUsersRequest)
		response, _ := raw.(*ram.ListUsersResponse)
		for _, user := range response.Users.User {
			// the d.Id() has changed from userName to userId since v1.44.0, add the logic for backward compatibility.
			if user.UserId == id || user.UserName == id {
				userName = user.UserName
				break
			}
		}
		if userName != "" || !response.IsTruncated {
			break
		}
		listUsersRequest.Marker = response.Marker
	}

	if userName == "" {
		return user, WrapErrorf(fmt.Errorf("there is no ram user with id or name is %s", id), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	getUserRequest := ram.CreateGetUserRequest()
	getUserRequest.RegionId = s.client.RegionId
	getUserRequest.UserName = userName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetUser(getUserRequest)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.User"}) {
			return user, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return user, WrapErrorf(err, DefaultErrorMsg, id, getUserRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(getUserRequest.GetActionName(), raw, getUserRequest.RpcRequest, getUserRequest)
	response, _ := raw.(*ram.GetUserResponse)

	return &response.User, nil
}

func (s *RamService) WaitForRamUser(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamUser(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.UserId == id {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, DefaultTimeoutMsg, id, GetFunc(1), ProviderERROR)
		}
	}
	return nil
}
func (s *RamService) DescribeRamGroupMembership(id string) (*ram.ListUsersForGroupResponse, error) {
	response := &ram.ListUsersForGroupResponse{}
	request := ram.CreateListUsersForGroupRequest()
	request.RegionId = s.client.RegionId
	request.GroupName = id
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListUsersForGroup(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response = raw.(*ram.ListUsersForGroupResponse)
	if len(response.Users.User) > 0 {
		return response, nil
	}
	return response, WrapErrorf(err, NotFoundMsg, ProviderERROR)
}

func (s *RamService) WaitForRamGroupMembership(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamGroupMembership(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.Itoa(len(object.Users.User)), status, ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamLoginProfile(id string) (*ram.GetLoginProfileResponse, error) {
	response := &ram.GetLoginProfileResponse{}
	request := ram.CreateGetLoginProfileRequest()
	request.RegionId = s.client.RegionId
	request.UserName = id

	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetLoginProfile(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.User.LoginProfile", "EntityNotExist.User"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response = raw.(*ram.GetLoginProfileResponse)
	return response, nil
}

func (s *RamService) WaitForRamLoginProfile(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamLoginProfile(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.LoginProfile.UserName == id {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.LoginProfile.UserName, id, ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamGroupPolicyAttachment(id string) (*ram.PolicyInListPoliciesForGroup, error) {
	response := &ram.PolicyInListPoliciesForGroup{}
	request := ram.CreateListPoliciesForGroupRequest()
	request.RegionId = s.client.RegionId
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return response, WrapError(err)
	}
	request.GroupName = parts[3]
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForGroup(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Group"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	listPoliciesForGroupResponse, _ := raw.(*ram.ListPoliciesForGroupResponse)
	if len(listPoliciesForGroupResponse.Policies.Policy) > 0 {
		for _, v := range listPoliciesForGroupResponse.Policies.Policy {
			if v.PolicyType == parts[2] && (v.PolicyName == parts[1] || strings.ToLower(v.PolicyName) == strings.ToLower(parts[1])) {
				return &v, nil
			}
		}
	}
	return response, WrapErrorf(err, NotFoundMsg, ProviderERROR)
}

func (s *RamService) WaitForRamGroupPolicyAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return WrapError(err)
	}
	for {
		object, err := s.DescribeRamGroupPolicyAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.PolicyName, parts[1], ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamAccountAlias(id string) (*ram.GetAccountAliasResponse, error) {
	response := &ram.GetAccountAliasResponse{}
	request := ram.CreateGetAccountAliasRequest()
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetAccountAlias(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response = raw.(*ram.GetAccountAliasResponse)

	return response, nil
}

func (s *RamService) DescribeRamAccessKey(id, userName string) (*ram.AccessKeyInListAccessKeys, error) {
	key := &ram.AccessKeyInListAccessKeys{}
	request := ram.CreateListAccessKeysRequest()
	request.RegionId = s.client.RegionId
	request.UserName = userName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListAccessKeys(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return key, WrapErrorf(Error(GetNotFoundMessage("RamAccessKey", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return key, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ram.ListAccessKeysResponse)
	for _, accessKey := range response.AccessKeys.AccessKey {
		if accessKey.AccessKeyId == id {
			return &accessKey, nil
		}
	}
	return key, WrapErrorf(Error(GetNotFoundMessage("RamAccessKey", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
}

func (s *RamService) WaitForRamAccessKey(id, useName string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamAccessKey(id, useName)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if string(status) == object.Status {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, status, ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewRamClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetPolicy"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PolicyName": id,
		"PolicyType": "Custom",
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("RamPolicy", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *RamService) DescribeRamRoleAttachment(id string) (*ecs.DescribeInstanceRamRoleResponse, error) {
	response := &ecs.DescribeInstanceRamRoleResponse{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return response, WrapError(err)
	}
	request := ecs.CreateDescribeInstanceRamRoleRequest()
	request.RegionId = s.client.RegionId
	request.InstanceIds = parts[1]
	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstanceRamRole(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"unexpected end of JSON input"}) {
				return resource.RetryableError(WrapError(err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidRamRole.NotFound"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response = raw.(*ecs.DescribeInstanceRamRoleResponse)
	instRoleSets := response.InstanceRamRoleSets.InstanceRamRoleSet
	if len(instRoleSets) > 0 {
		var instIds []string
		for _, item := range instRoleSets {
			if item.RamRoleName == parts[0] {
				instIds = append(instIds, item.InstanceId)
			}
		}
		ids := strings.Split(strings.TrimRight(strings.TrimLeft(strings.Replace(strings.Split(id, ":")[1], "\"", "", -1), "["), "]"), ",")
		sort.Strings(instIds)
		sort.Strings(ids)
		if reflect.DeepEqual(instIds, ids) {
			return response, nil
		}
	}
	return response, WrapErrorf(err, NotFoundMsg, ProviderERROR)
}

func (s *RamService) WaitForRamRoleAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamRoleAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.Itoa(object.TotalCount), status, ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamRole(id string) (*ram.GetRoleResponse, error) {
	response := &ram.GetRoleResponse{}
	request := ram.CreateGetRoleRequest()
	request.RegionId = s.client.RegionId
	request.RoleName = id
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetRole(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ram.GetRoleResponse)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *RamService) WaitForRamRole(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamRole(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Role.RoleName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Role.RoleName, id, ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamUserPolicyAttachment(id string) (*ram.PolicyInListPoliciesForUser, error) {
	response := &ram.PolicyInListPoliciesForUser{}
	request := ram.CreateListPoliciesForUserRequest()
	request.RegionId = s.client.RegionId
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return response, WrapError(err)
	}
	request.UserName = parts[3]
	var listPoliciesForUserResponse *ram.ListPoliciesForUserResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForUser(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		listPoliciesForUserResponse, _ = raw.(*ram.ListPoliciesForUserResponse)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if len(listPoliciesForUserResponse.Policies.Policy) > 0 {
		for _, v := range listPoliciesForUserResponse.Policies.Policy {
			if v.PolicyType == parts[2] && (v.PolicyName == parts[1] || strings.ToLower(v.PolicyName) == strings.ToLower(parts[1])) {
				return &v, nil
			}
		}
	}
	return response, WrapErrorf(err, NotFoundMsg, ProviderERROR)
}

func (s *RamService) WaitForRamUserPolicyAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return WrapError(err)
	}
	for {
		object, err := s.DescribeRamUserPolicyAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.PolicyName, parts[1], ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamRolePolicyAttachment(id string) (*ram.PolicyInListPoliciesForRole, error) {
	response := &ram.PolicyInListPoliciesForRole{}
	request := ram.CreateListPoliciesForRoleRequest()
	request.RegionId = s.client.RegionId
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return response, WrapError(err)
	}
	request.RoleName = parts[3]
	var listPoliciesForRoleResponse *ram.ListPoliciesForRoleResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForRole(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)

		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		listPoliciesForRoleResponse, _ = raw.(*ram.ListPoliciesForRoleResponse)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if len(listPoliciesForRoleResponse.Policies.Policy) > 0 {
		for _, v := range listPoliciesForRoleResponse.Policies.Policy {
			if v.PolicyType == parts[2] && (v.PolicyName == parts[1] || strings.ToLower(v.PolicyName) == strings.ToLower(parts[1])) {
				return &v, nil
			}
		}
	}
	return response, WrapErrorf(err, NotFoundMsg, ProviderERROR)
}

func (s *RamService) WaitForRamRolePolicyAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return WrapError(err)
	}
	for {
		object, err := s.DescribeRamRolePolicyAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {

				}
				return nil
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.PolicyName, parts[1], ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamGroup(id string) (*ram.GetGroupResponse, error) {
	response := &ram.GetGroupResponse{}
	request := ram.CreateGetGroupRequest()
	request.RegionId = s.client.RegionId
	request.GroupName = id
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetGroup(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Group"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*ram.GetGroupResponse)

	if response.Group.GroupName != id {
		return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *RamService) WaitForRamGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamGroup(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Group.GroupName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Group.GroupName, id, ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamAccountPasswordPolicy(id string) (*ram.GetPasswordPolicyResponse, error) {
	response := &ram.GetPasswordPolicyResponse{}
	request := ram.CreateGetPasswordPolicyRequest()
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetPasswordPolicy(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*ram.GetPasswordPolicyResponse)

	return response, nil
}

func (s *RamService) DescribeRamSecurityPreference(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewImsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetSecurityPreference"
	request := map[string]interface{}{}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-15"), StringPointer("AK"), nil, request, &runtime)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.SecurityPreference", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SecurityPreference", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *RamService) DescribeRamServiceLinkedRole(id string) (*ram.GetRoleResponse, error) {
	parts, _ := ParseResourceId(id, 2)
	id = parts[1]

	response := &ram.GetRoleResponse{}
	request := ram.CreateGetRoleRequest()
	request.RegionId = s.client.RegionId
	request.RoleName = id
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetRole(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ram.GetRoleResponse)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

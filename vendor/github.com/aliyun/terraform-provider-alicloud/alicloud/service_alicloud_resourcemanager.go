package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ResourcemanagerService struct {
	client *connectivity.AliyunClient
}

func (s *ResourcemanagerService) DescribeResourceManagerRole(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetRole"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"RoleName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerRole", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Role", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Role", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerResourceGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetResourceGroup"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"ResourceGroupId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.ResourceGroup"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerResourceGroup", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.ResourceGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ResourceGroup", response)
	}
	if vv, ok := v.(map[string]interface{})["Status"].(string); ok && vv == "PendingDelete" {
		log.Printf("[WARN] Removing ResourceManagerResourceGroup  %s because it's already gone", id)
		return v.(map[string]interface{}), WrapErrorf(Error(GetNotFoundMessage("ResourceManagerResourceGroup", id)), NotFoundMsg, ProviderERROR)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) ResourceManagerResourceGroupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeResourceManagerResourceGroup(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *ResourcemanagerService) DescribeResourceManagerFolder(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetFolder"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"FolderId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Folder", "EntityNotExists.ResourceDirectory"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerFolder", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Folder", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Folder", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerHandshake(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetHandshake"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"HandshakeId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Handshake"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerHandshake", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Handshake", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Handshake", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) GetPolicyVersion(id string, d *schema.ResourceData) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetPolicyVersion"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PolicyName": id,
		"PolicyType": "Custom",
	}
	if v, ok := d.GetOk("default_version"); ok {
		request["VersionId"] = v.(string)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExist.Policy.Version"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerPolicy", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.PolicyVersion", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PolicyVersion", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
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
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerPolicy", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Policy", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Policy", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerAccount(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetAccount"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"AccountId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Account", "EntityNotExists.ResourceDirectory"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("ResourceManager:Account", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Account", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Account", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerResourceDirectory(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetResourceDirectory"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceDirectoryNotInUse"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerResourceDirectory", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.ResourceDirectory", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ResourceDirectory", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerPolicyVersion(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetPolicyVersion"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PolicyName": parts[0],
		"VersionId":  parts[1],
		"PolicyType": "Custom",
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExist.Policy.Version"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerPolicyVersion", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.PolicyVersion", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PolicyVersion", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerPolicyAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListPolicyAttachments"
	parts, err := ParseResourceId(id, 5)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"PolicyName":      parts[0],
		"PolicyType":      parts[1],
		"PrincipalName":   parts[2],
		"PrincipalType":   parts[3],
		"ResourceGroupId": parts[4],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExists.ResourceGroup"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerPolicyAttachment", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.PolicyAttachments.PolicyAttachment", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PolicyAttachments.PolicyAttachment", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ResourceManager", id)), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerControlPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetControlPolicy"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"PolicyId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.ControlPolicy"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerControlPolicy", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.ControlPolicy", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ControlPolicy", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerControlPolicyAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListControlPolicyAttachmentsForTarget"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"TargetId": parts[1],
	}
	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Target"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerControlPolicyAttachment", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.ControlPolicyAttachments.ControlPolicyAttachment", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ControlPolicyAttachments.ControlPolicyAttachment", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ResourceManager", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if v.(map[string]interface{})["PolicyId"].(string) == parts[0] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("ResourceManager", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *ResourcemanagerService) ResourceManagerResourceDirectoryStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeResourceManagerResourceDirectory(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["ScpStatus"].(string) == failState {
				return object, object["ScpStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["ScpStatus"].(string)))
			}
		}
		return object, object["ScpStatus"].(string), nil
	}
}
func (s *ResourcemanagerService) GetPayerForAccount(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetPayerForAccount"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"AccountId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

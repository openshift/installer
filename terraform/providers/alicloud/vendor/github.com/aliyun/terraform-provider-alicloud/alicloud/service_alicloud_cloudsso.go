package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type CloudssoService struct {
	client *connectivity.AliyunClient
}

func (s *CloudssoService) GetDirectory(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetDirectory"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:Directory", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Directory", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Directory", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) GetMFAAuthenticationStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetMFAAuthenticationStatus"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:Directory", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) GetSCIMSynchronizationStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetSCIMSynchronizationStatus"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:Directory", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) GetDirectoryTasks(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTasks"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:Directory", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Tasks", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) GetExternalSAMLIdentityProvider(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetExternalSAMLIdentityProvider"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:Directory", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.SAMLIdentityProviderConfiguration", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SAMLIdentityProviderConfiguration", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) DescribeCloudSsoDirectory(id string) (object map[string]interface{}, err error) {

	getDirectoryObject, err := s.GetDirectory(id)
	getMFAAuthenticationStatusObject, err := s.GetMFAAuthenticationStatus(id)
	if err != nil {
		return nil, err
	}
	getSCIMSynchronizationStatusObject, err := s.GetSCIMSynchronizationStatus(id)
	if err != nil {
		return nil, err
	}
	getExternalSAMLIdentityProviderObject, err := s.GetExternalSAMLIdentityProvider(id)
	if err != nil {
		return nil, err
	}
	getDirectoryObject["MFAAuthenticationStatus"] = getMFAAuthenticationStatusObject["MFAAuthenticationStatus"]
	getDirectoryObject["SCIMSynchronizationStatus"] = getSCIMSynchronizationStatusObject["SCIMSynchronizationStatus"]
	getDirectoryObject["SAMLIdentityProviderConfiguration"] = getExternalSAMLIdentityProviderObject

	return getDirectoryObject, nil
}

func (s *CloudssoService) DescribeCloudSsoScimServerCredential(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListSCIMServerCredentials"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DirectoryId": parts[0],
	}
	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:SCIMServerCredential", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.SCIMServerCredentials", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SCIMServerCredentials", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["CredentialId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *CloudssoService) DescribeCloudSsoGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetGroup"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"GroupId":     parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Group"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:Group", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Group", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Group", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *CloudssoService) DescribeCloudSsoUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
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
		"DirectoryId": parts[0],
		"UserId":      parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.User"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:User", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *CloudssoService) ListMFADevicesForUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListMFADevicesForUser"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"UserId":      parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory", "EntityNotExists.User"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:User", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.MFADevices", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *CloudssoService) DescribeCloudSsoAccessConfiguration(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetAccessConfiguration"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AccessConfigurationId": parts[1],
		"DirectoryId":           parts[0],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.AccessConfiguration"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:AccessConfiguration", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.AccessConfiguration", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AccessConfiguration", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *CloudssoService) ListPermissionPoliciesInAccessConfiguration(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListPermissionPoliciesInAccessConfiguration"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AccessConfigurationId": parts[1],
		"DirectoryId":           parts[0],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.AccessConfiguration"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:AccessConfiguration", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *CloudssoService) DescribeCloudSsoUserAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListGroupMembers"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"GroupId":     parts[1],
		"MaxResults":  PageSizeLarge,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
			if IsExpectedErrors(err, []string{"EntityNotExists.Directory", "EntityNotExists.Group"}) {
				return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:UserAttachment", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.GroupMembers", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.GroupMembers", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["UserId"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *CloudssoService) DescribeCloudSsoAccessAssignment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListAccessAssignments"
	parts, err := ParseResourceId(id, 6)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DirectoryId":           parts[0],
		"AccessConfigurationId": parts[1],
		"TargetType":            parts[2],
		"TargetId":              parts[3],
		"PrincipalType":         parts[4],
		"PrincipalId":           parts[5],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:AccessAttachment", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.AccessAssignments", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AccessAssignments", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["PrincipalId"]) != parts[5] {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CloudssoService) GetTaskStatus(directoryId, taskId string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetTaskStatus"
	request := map[string]interface{}{
		"DirectoryId": directoryId,
		"TaskId":      taskId,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Task"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO", fmt.Sprint(directoryId, ":", taskId))), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, fmt.Sprint(directoryId, ":", taskId), action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.TaskStatus", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, fmt.Sprint(directoryId, ":", taskId), "$.TaskStatus", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *CloudssoService) CloudssoServiceAccessAssignmentStateRefreshFunc(directoryId, taskId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.GetTaskStatus(directoryId, taskId)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprintf(directoryId, ":", taskId)))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CloudssoService) DescribeCloudSsoAccessConfigurationProvisioning(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListAccessConfigurationProvisionings"
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AccessConfigurationId": parts[1],
		"DirectoryId":           parts[0],
		"TargetId":              parts[3],
		"TargetType":            parts[2],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.AccessConfigurationProvisioning"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:AccessConfigurationProvisioning", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.AccessConfigurationProvisionings", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AccessConfigurationProvisionings", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["TargetId"]) != parts[3] {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CloudssoService) CloudssoServiceAccessConfigurationProvisioningStateRefreshFunc(directoryId, taskId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.GetTaskStatus(directoryId, taskId)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprintf(directoryId, ":", taskId)))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CloudssoService) CloudssoServicAccessConfigurationProvisioning(directoryId, accessConfigurationId, targetType, targetId string) (err error) {
	var response map[string]interface{}
	action := "ProvisionAccessConfiguration"
	request := make(map[string]interface{})
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}

	request["DirectoryId"] = directoryId
	request["AccessConfigurationId"] = accessConfigurationId
	request["TargetType"] = targetType
	request["TargetId"] = targetId

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict.Task"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_sso_access_configuration_provisioning", action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Tasks", response)
	if err != nil || len(v.([]interface{})) < 1 {
		return WrapErrorf(err, IdMsg, fmt.Sprint(directoryId, ":", accessConfigurationId, ":", targetType, ":", targetId))
	}
	response = v.([]interface{})[0].(map[string]interface{})
	_, err = s.GetTaskStatus(fmt.Sprint(request["DirectoryId"]), fmt.Sprint(response["TaskId"]))
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	stateConf := BuildStateConf([]string{}, []string{"Success"}, 10*time.Minute, 5*time.Second, s.CloudssoServiceAccessConfigurationProvisioningStateRefreshFunc(fmt.Sprint(request["DirectoryId"]), fmt.Sprint(response["TaskId"]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, fmt.Sprint(directoryId, ":", accessConfigurationId, ":", targetType, ":", targetId))
	}
	return nil
}

func (s *CloudssoService) DescribeCloudSsoAccessConfigurationProvisionings(directoryId, accessConfigurationId string) (objects []map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}

	action := "ListAccessConfigurationProvisionings"
	request := map[string]interface{}{
		"AccessConfigurationId": accessConfigurationId,
		"DirectoryId":           directoryId,
		"MaxResults":            PageSizeMedium,
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
			return objects, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_sso_access_configuration", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AccessConfigurationProvisionings", response)
		if err != nil {
			return objects, WrapErrorf(err, FailedGetAttributeMsg, action, "$.AccessConfigurationProvisionings", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	return objects, nil
}

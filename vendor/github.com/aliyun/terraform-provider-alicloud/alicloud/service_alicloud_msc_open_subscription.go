package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type MscOpenSubscriptionService struct {
	client *connectivity.AliyunClient
}

func (s *MscOpenSubscriptionService) GetContact(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewMscopensubscriptionClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetContact"
	request := map[string]interface{}{
		"ContactId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2021-07-13"), StringPointer("AK"), request, nil, &runtime)
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
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"ResourceNotFound"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("MscSub:Contact", id)), NotFoundMsg, ProviderERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Contact", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Contact", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *MscOpenSubscriptionService) DescribeMscSubContact(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewMscopensubscriptionClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetContact"
	request := map[string]interface{}{
		"ContactId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2021-07-13"), StringPointer("AK"), request, nil, &runtime)
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
	if IsExpectedErrors(err, []string{"ResourceNotFound"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("MscSub:Contact", id)), NotFoundMsg, ProviderERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Contact", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Contact", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type SaeService struct {
	client *connectivity.AliyunClient
}

func (s *SaeService) DescribeSaeNamespace(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServerlessClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "/pop/v1/paas/namespace"
	request := map[string]*string{
		"NamespaceId": StringPointer(id),
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidNamespaceId.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("SAE:Namespace", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, "alicloud_sae_namespace", "GET "+action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	addDebug(action, response, request)
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"InvalidNamespaceId.NotFound"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("SAE:Namespace", id)), NotFoundMsg, ProviderERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", "GET "+action, response))
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *SaeService) DescribeSaeConfigMap(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServerlessClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "/pop/v1/sam/configmap/configMap"
	request := map[string]*string{
		"ConfigMapId": StringPointer(id),
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"NotFound.ConfigMap"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("SAE:ConfigMap", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	addDebug(action, response, request)

	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

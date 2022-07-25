package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ScdnService struct {
	client *connectivity.AliyunClient
}

func (s *ScdnService) DescribeScdnDomainCertificateInfo(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewScdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeScdnDomainCertificateInfo"
	request := map[string]interface{}{
		"DomainName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$.CertInfos.CertInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.CertInfos.CertInfo", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("SCDN", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DomainName"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("SCDN", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ScdnService) convertFunctionsToString(v []interface{}) (string, error) {
	arrayMaps := make([]interface{}, len(v))
	for i, vv := range v {
		item := vv.(map[string]interface{})
		temp := map[string]interface{}{
			"ConfigId":     item["config_id"],
			"FunctionName": item["function_name"],
			"Status":       item["status"],
		}

		functionArgs := item["function_args"].(*schema.Set).List()
		functionArgsMaps := make([]interface{}, len(functionArgs))
		for j, functionArgsVal := range functionArgs {
			functionArgsItem := functionArgsVal.(map[string]interface{})
			functionArgsTemp := map[string]interface{}{
				"ArgName":  functionArgsItem["arg_name"],
				"ArgValue": functionArgsItem["arg_value"],
			}

			functionArgsMaps[j] = functionArgsTemp
		}
		temp["function_args"] = functionArgsMaps
		arrayMaps[i] = temp
	}
	maps, err := json.Marshal(arrayMaps)
	if err != nil {
		return "", WrapError(err)
	}
	return string(maps), nil
}

func (s *ScdnService) DescribeScdnDomain(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewScdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeScdnDomainDetail"
	request := map[string]interface{}{
		"DomainName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("SCDN:Domain", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DomainDetail", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DomainDetail", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ScdnService) ScdnDomainStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeScdnDomain(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DomainStatus"]) == failState {
				return object, fmt.Sprint(object["DomainStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["DomainStatus"])))
			}
		}
		return object, fmt.Sprint(object["DomainStatus"]), nil
	}
}

func (s *ScdnService) DescribeScdnCertificateList(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewScdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeScdnCertificateList"
	request := map[string]interface{}{
		"DomainName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$.CertificateListModel.CertList.Cert", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.CertificateListModel.CertList.Cert", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("SCDN", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DomainName"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("SCDN", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ScdnService) DescribeScdnDomainCname(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewScdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeScdnDomainCname"
	request := map[string]interface{}{
		"DomainName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2017-11-15"), StringPointer("AK"), request, nil, &runtime)
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
	v, err := jsonpath.Get("$.CnameDatas.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.CnameDatas.Data", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("SCDN", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DomainName"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("SCDN", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ScdnService) DescribeScdnDomainConfigs(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewScdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeScdnDomainConfigs"
	request := map[string]interface{}{
		"DomainName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$.DomainConfigs.DomainConfig", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DomainConfigs.DomainConfig", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("SCDN", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DomainName"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("SCDN", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ScdnService) DescribeScdnDomainDetail(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewScdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeScdnDomainDetail"
	request := map[string]interface{}{
		"DomainName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("SCDN:Domain", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DomainDetail", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DomainDetail", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ScdnService) DescribeScdnDomainConfig(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewScdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	action := "DescribeScdnDomainConfigs"
	request := map[string]interface{}{
		"DomainName":    parts[0],
		"FunctionNames": parts[1],
	}
	if parts[2] != "" {
		request["ConfigId"] = parts[2]
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ScdnDomainConfig", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.DomainConfigs.DomainConfig", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DomainConfigs.DomainConfig", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("DCDN:DomainConfig", id)), NotFoundWithResponse, response)
	} else if len(v.([]interface{})) > 0 {
		object = v.([]interface{})[0].(map[string]interface{})
	}
	return object, nil
}

func (s *ScdnService) ScdnDomainConfigStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeScdnDomainConfig(id)
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

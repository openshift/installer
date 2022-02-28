package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type GaService struct {
	client *connectivity.AliyunClient
}

func (s *GaService) DescribeGaAccelerator(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewGaplusClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeAccelerator"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"AcceleratorId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"UnknownError"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Ga:Accelerator", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *GaService) GaAcceleratorStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaAccelerator(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeGaListener(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewGaplusClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeListener"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"ListenerId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotExist.Listener", "UnknownError"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("GaListener", id)), NotFoundMsg, ProviderERROR)
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

func (s *GaService) GaListenerStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaListener(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["State"].(string) == failState {
				return object, object["State"].(string), WrapError(Error(FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, object["State"].(string), nil
	}
}

func (s *GaService) DescribeGaBandwidthPackage(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewGaplusClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeBandwidthPackage"
	request := map[string]interface{}{
		"RegionId":           s.client.RegionId,
		"BandwidthPackageId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotExist.BandwidthPackage"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("GaBandwidthPackage", id)), NotFoundMsg, ProviderERROR)
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

func (s *GaService) GaBandwidthPackageStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBandwidthPackage(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["State"].(string) == failState {
				return object, object["State"].(string), WrapError(Error(FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, object["State"].(string), nil
	}
}

func (s *GaService) DescribeGaEndpointGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewGaplusClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeEndpointGroup"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"EndpointGroupId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotExist.EndPointGroup"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("GaEndpointGroup", id)), NotFoundMsg, ProviderERROR)
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
	if object["State"] == nil {
		err = WrapErrorf(Error(GetNotFoundMessage("GaEndpointGroup", id)), NotFoundMsg, ProviderERROR)
		return object, err
	}
	return object, nil
}

func (s *GaService) DescribeGaForwardingRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewGaplusClient()
	if err != nil {
		return nil, WrapError(err)
	}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	action := "ListForwardingRules"
	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"ListenerId":       parts[1],
		"AcceleratorId":    parts[0],
		"ForwardingRuleId": parts[2],
		"MaxResults":       PageSizeLarge,
	}
	request["ClientToken"] = buildClientToken(action)
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	err = resource.Retry(5*time.Second, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator"}) {
				return resource.RetryableError(WrapErrorf(Error(GetNotFoundMessage("ForwardingRule", id)), NotFoundMsg, ProviderERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR))
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return
	}
	v, err := jsonpath.Get("$.ForwardingRules", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ForwardingRules", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ForwardingRule", id)), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *GaService) DescribeGaIpSet(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewGaplusClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeIpSet"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"IpSetId":  id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"UnknownError"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("GaIpSet", id)), NotFoundMsg, ProviderERROR)
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
	if object["State"] == nil {
		err = WrapErrorf(Error(GetNotFoundMessage("GaIpSet", id)), NotFoundMsg, ProviderERROR)
		return object, err
	}
	return object, nil
}

func (s *GaService) DescribeGaBandwidthPackageAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	conn, err := s.client.NewGaplusClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeAccelerator"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"AcceleratorId": parts[0],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"UnknownError"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("GaBandwidthPackageAttachment", id)), NotFoundMsg, ProviderERROR)
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
	if object["BasicBandwidthPackage"] == nil || object["BasicBandwidthPackage"].(map[string]interface{})["InstanceId"] != parts[1] {
		return object, WrapErrorf(Error(GetNotFoundMessage("GaBandwidthPackageAttachment", id)), NotFoundMsg, ProviderERROR)
	}
	return object, nil
}

func (s *GaService) GaEndpointGroupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaEndpointGroup(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["State"].(string) == failState {
				return object, object["State"].(string), WrapError(Error(FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, object["State"].(string), nil
	}
}

func (s *GaService) GaForwardingRuleStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaForwardingRule(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["ForwardingRuleStatus"].(string) == failState {
				return object, object["ForwardingRuleStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["ForwardingRuleStatus"].(string)))
			}
		}
		return object, object["ForwardingRuleStatus"].(string), nil
	}
}

func (s *GaService) GaBandwidthPackageAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBandwidthPackageAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["State"].(string) == failState {
				return object, object["State"].(string), WrapError(Error(FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, object["State"].(string), nil
	}
}

func (s *GaService) GaIpSetStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaIpSet(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["State"].(string) == failState {
				return object, object["State"].(string), WrapError(Error(FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, object["State"].(string), nil
	}
}
func (s *GaService) DescribeAcceleratorAutoRenewAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewGaplusClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeAcceleratorAutoRenewAttribute"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"AcceleratorId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
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

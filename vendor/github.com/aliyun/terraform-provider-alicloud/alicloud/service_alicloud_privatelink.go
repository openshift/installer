package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type PrivatelinkService struct {
	client *connectivity.AliyunClient
}

func (s *PrivatelinkService) ListVpcEndpointServiceResources(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListVpcEndpointServiceResources"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ServiceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpointService", id)), NotFoundMsg, ProviderERROR)
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

func (s *PrivatelinkService) DescribePrivatelinkVpcEndpointService(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetVpcEndpointServiceAttribute"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ServiceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpointService", id)), NotFoundMsg, ProviderERROR)
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

func (s *PrivatelinkService) PrivatelinkVpcEndpointServiceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePrivatelinkVpcEndpointService(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["ServiceStatus"].(string) == failState {
				return object, object["ServiceStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["ServiceStatus"].(string)))
			}
		}
		return object, object["ServiceStatus"].(string), nil
	}
}

func (s *PrivatelinkService) DescribePrivatelinkVpcEndpointConnection(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListVpcEndpointConnections"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": parts[1],
		"ServiceId":  parts[0],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpointConnection", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Connections", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Connections", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("PrivateLink", id)), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *PrivatelinkService) PrivatelinkVpcEndpointConnectionStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePrivatelinkVpcEndpointConnection(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["ConnectionStatus"].(string) == failState {
				return object, object["ConnectionStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["ConnectionStatus"].(string)))
			}
		}
		return object, object["ConnectionStatus"].(string), nil
	}
}

func (s *PrivatelinkService) ListVpcEndpointSecurityGroups(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListVpcEndpointSecurityGroups"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpoint", id)), NotFoundMsg, ProviderERROR)
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

func (s *PrivatelinkService) ListVpcEndpointZones(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListVpcEndpointZones"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpoint", id)), NotFoundMsg, ProviderERROR)
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

func (s *PrivatelinkService) DescribePrivatelinkVpcEndpoint(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetVpcEndpointAttribute"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("GetVpcEndpointAttribute")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpoint", id)), NotFoundMsg, ProviderERROR)
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

func (s *PrivatelinkService) PrivatelinkVpcEndpointStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePrivatelinkVpcEndpoint(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["EndpointStatus"].(string) == failState {
				return object, object["EndpointStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["EndpointStatus"].(string)))
			}
		}
		return object, object["EndpointStatus"].(string), nil
	}
}

func (s *PrivatelinkService) DescribePrivatelinkVpcEndpointServiceResource(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListVpcEndpointServiceResources"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ServiceId": parts[0],
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpointServiceResource", id)), NotFoundMsg, ProviderERROR)
				return object, err
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return object, err
		}
		addDebug(action, response, request)
		v, err := jsonpath.Get("$.Resources", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Resources", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("PrivateLink", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if v.(map[string]interface{})["ResourceId"].(string) == parts[1] {
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
		return object, nil
	}
	return
}

func (s *PrivatelinkService) DescribePrivatelinkVpcEndpointServiceUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListVpcEndpointServiceUsers"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ServiceId": parts[0],
		"UserId":    parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpointServiceUser", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Users", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Users", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("PrivateLink", id)), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *PrivatelinkService) DescribePrivatelinkVpcEndpointZone(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPrivatelinkClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListVpcEndpointZones"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": parts[0],
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"EndpointNotFound"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("PrivatelinkVpcEndpointZone", id)), NotFoundMsg, ProviderERROR)
				return object, err
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return object, err
		}
		addDebug(action, response, request)
		v, err := jsonpath.Get("$.Zones", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Zones", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("PrivateLink", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if v.(map[string]interface{})["ZoneId"].(string) == parts[1] {
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
		return object, nil
	}
	return
}

func (s *PrivatelinkService) PrivatelinkVpcEndpointZoneStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePrivatelinkVpcEndpointZone(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["ZoneStatus"].(string) == failState {
				return object, object["ZoneStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["ZoneStatus"].(string)))
			}
		}
		return object, object["ZoneStatus"].(string), nil
	}
}

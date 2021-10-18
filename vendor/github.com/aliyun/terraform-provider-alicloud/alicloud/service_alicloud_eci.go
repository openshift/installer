package alicloud

import (
	"encoding/json"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type EciService struct {
	client *connectivity.AliyunClient
}

func (s *EciService) DescribeEciImageCache(id string) (object eci.DescribeImageCachesImageCache0, err error) {
	request := eci.CreateDescribeImageCachesRequest()
	request.RegionId = s.client.RegionId

	request.ImageCacheId = id

	raw, err := s.client.WithEciClient(func(eciClient *eci.Client) (interface{}, error) {
		return eciClient.DescribeImageCaches(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*eci.DescribeImageCachesResponse)

	if len(response.ImageCaches) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("EciImageCache", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return response.ImageCaches[0], nil
}

func (s *EciService) EciImageCacheStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEciImageCache(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *EciService) DescribeEciContainerGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewEciClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeContainerGroups"
	jsonId, err := json.Marshal([]string{id})
	if err != nil {
		return nil, WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":          s.client.RegionId,
		"ContainerGroupIds": string(jsonId),
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-08"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidParameter.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("EciContainerGroup", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.ContainerGroups", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ContainerGroups", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ECIOpenAPI", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["ContainerGroupId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("ECIOpenAPI", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *EciService) EciContainerGroupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEciContainerGroup(id)
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

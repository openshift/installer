package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ActiontrailService struct {
	client *connectivity.AliyunClient
}

func (s *ActiontrailService) DescribeActiontrailTrail(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewActiontrailClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeTrails"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"NameList": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.TrailList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TrailList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ActionTrail", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["Name"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("ActionTrail", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ActiontrailService) ActiontrailTrailStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeActiontrailTrail(id)
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

func (s *ActiontrailService) DescribeActiontrailHistoryDeliveryJob(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewActiontrailClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetDeliveryHistoryJob"
	request := map[string]interface{}{
		"JobId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"DeliveryHistoryJobNotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("ActionTrail:HistoryDeliveryJob", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *ActiontrailService) ActiontrailHistoryDeliveryJobStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeActiontrailHistoryDeliveryJob(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["JobStatus"]) == failState {
				return object, fmt.Sprint(object["JobStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["JobStatus"])))
			}
		}
		return object, fmt.Sprint(object["JobStatus"]), nil
	}
}

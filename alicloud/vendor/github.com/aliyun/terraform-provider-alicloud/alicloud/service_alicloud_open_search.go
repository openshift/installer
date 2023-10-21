package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type OpenSearchService struct {
	client *connectivity.AliyunClient
}

func (s *OpenSearchService) DescribeOpenSearchAppGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewOpensearchClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "/v4/openapi/app-groups/" + id
	body := map[string]interface{}{
		"appGroupIdentity": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2017-12-25"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), nil, nil, body, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug("GET "+action, response, body)
	if err != nil {
		if IsExpectedErrors(err, []string{"App.NotFound"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}

	v, err := jsonpath.Get("$.result", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.result", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *OpenSearchService) OpenSearchAppStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOpenSearchAppGroup(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["status"].(string) == failState {
				return object, object["status"].(string), WrapError(Error(FailedToReachTargetStatus, object["status"].(string)))
			}
		}
		return object, object["status"].(string), nil
	}
}

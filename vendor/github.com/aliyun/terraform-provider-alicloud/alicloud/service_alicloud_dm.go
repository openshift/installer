package alicloud

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DmService struct {
	client *connectivity.AliyunClient
}

func (s *DmService) DescribeDirectMailReceivers(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDmClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "QueryReceiverByParam"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"PageSize": 20,
		"PageNo":   1,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-11-23"), StringPointer("AK"), nil, request, &runtime)
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
		v, err := jsonpath.Get("$.data.receiver", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data.receiver", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("DirectMail", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if v.(map[string]interface{})["ReceiverId"].(string) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNo"] = request["PageNo"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("DirectMail", id)), NotFoundWithResponse, response)
	}
	return
}

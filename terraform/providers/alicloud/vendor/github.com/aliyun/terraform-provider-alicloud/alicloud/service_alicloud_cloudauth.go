package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type CloudauthService struct {
	client *connectivity.AliyunClient
}

func (s *CloudauthService) DescribeCloudauthFaceConfig(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudauthClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeFaceConfig"
	request := map[string]interface{}{
		"Lang": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-07"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$.Items", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Cloudauth", id)), NotFoundWithResponse, response)
	} else {
		for _, obj := range v.([]interface{}) {
			if fmt.Sprint(obj.(map[string]interface{})["BizType"]) == id {
				object = v.([]interface{})[0].(map[string]interface{})
				return object, nil
			}
		}
		return object, WrapErrorf(Error(GetNotFoundMessage("Cloudauth", id)), NotFoundWithResponse, response)
	}
}

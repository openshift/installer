package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type ImsService struct {
	client *connectivity.AliyunClient
}

func (s *ImsService) DescribeRamSamlProvider(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewImsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetSAMLProvider"
	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"SAMLProviderName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.SAMLProviderError"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("RamSamlProvider", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.SAMLProvider", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SAMLProvider", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

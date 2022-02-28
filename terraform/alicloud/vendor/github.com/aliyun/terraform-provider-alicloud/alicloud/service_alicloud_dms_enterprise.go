package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type Dms_enterpriseService struct {
	client *connectivity.AliyunClient
}

func (s *Dms_enterpriseService) DescribeDmsEnterpriseInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDmsenterpriseClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetInstance"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Host":     parts[0],
		"Port":     parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNoEnoughNumber"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DmsEnterpriseInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Instance", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Instance", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Dms_enterpriseService) DescribeDmsEnterpriseUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDmsenterpriseClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetUser"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Uid":      id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.User", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.User", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

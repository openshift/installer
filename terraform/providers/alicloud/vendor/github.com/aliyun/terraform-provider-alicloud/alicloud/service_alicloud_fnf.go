package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type FnfService struct {
	client *connectivity.AliyunClient
}

func (s *FnfService) DescribeFnfFlow(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewFnfClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeFlow"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Name":     id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-03-15"), StringPointer("AK"), request, nil, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"FlowNotExists"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("FnfFlow", id)), NotFoundMsg, ProviderERROR)
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

func (s *FnfService) DescribeFnfSchedule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewFnfClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeSchedule"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"FlowName":     parts[0],
		"ScheduleName": parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-03-15"), StringPointer("AK"), request, nil, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"FlowNotExists", "ScheduleNotExists"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("FnfSchedule", id)), NotFoundMsg, ProviderERROR)
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

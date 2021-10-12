package alicloud

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type Brain_industrialService struct {
	client *connectivity.AliyunClient
}

func (s *Brain_industrialService) DescribeBrainIndustrialPidOrganization(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewAistudioClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListPidOrganizations"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}
	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		err = Error("ListPidOrganizations failed for " + response["Message"].(string))
		return object, err
	}
	v, err := jsonpath.Get("$.OrganizationList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.OrganizationList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("BrainIndustrial", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if v.(map[string]interface{})["OrganizationId"].(string) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("BrainIndustrial", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *Brain_industrialService) DescribeBrainIndustrialPidProject(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewAistudioClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListPidProjects"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"CurrentPage": 1,
		"PageSize":    20,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return object, err
		}
		addDebug(action, response, request)
		if fmt.Sprintf(`%v`, response["Code"]) != "200" {
			err = Error("ListPidProjects failed for " + response["Message"].(string))
			return object, err
		}
		v, err := jsonpath.Get("$.PidProjectList", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PidProjectList", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("BrainIndustrial", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if v.(map[string]interface{})["PidProjectId"].(string) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("BrainIndustrial", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *Brain_industrialService) DescribeBrainIndustrialPidLoop(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewAistudioClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetLoop"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"PidLoopId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["Code"]), []string{"-106"}) {
		err = WrapErrorf(Error(GetNotFoundMessage("BrainIndustrialPidLoop", id)), NotFoundMsg, ProviderERROR)
		return object, err
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		err = Error("GetLoop failed for " + response["Message"].(string))
		return object, err
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

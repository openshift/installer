package alicloud

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type MaxcomputeService struct {
	client *connectivity.AliyunClient
}

func (s *MaxcomputeService) DescribeMaxcomputeProject(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewOdpsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetProject"
	request := map[string]interface{}{
		"RegionName":  s.client.RegionId,
		"ProjectName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-12"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["Code"]), []string{"102", "403"}) {
		err = WrapErrorf(Error(GetNotFoundMessage("MaxcomputeProject", id)), NotFoundMsg, ProviderERROR)
		return object, err
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		err = Error("GetProject failed for " + response["Message"].(string))
		return object, err
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

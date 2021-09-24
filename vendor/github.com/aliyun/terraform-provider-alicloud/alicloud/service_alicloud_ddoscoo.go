package alicloud

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type DdoscooService struct {
	client *connectivity.AliyunClient
}

func (s *DdoscooService) DescribeDdoscooInstance(id string) (v ddoscoo.Instance, err error) {
	request := ddoscoo.CreateDescribeInstancesRequest()
	request.RegionId = "cn-hangzhou"
	request.InstanceIds = &[]string{id}
	request.PageNumber = "1"
	request.PageSize = "10"

	raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.DescribeInstances(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return v, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}

		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ddoscoo.DescribeInstancesResponse)
	if len(response.Instances) == 0 || response.Instances[0].InstanceId != id {
		return v, WrapErrorf(Error(GetNotFoundMessage("Ddoscoo Instance", id)), NotFoundMsg, ProviderERROR)
	}

	v = response.Instances[0]
	return
}

// 创建实例后轮询查询直到创建成功
func (s *DdoscooService) DdosStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDdoscooInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		status := ""
		if object.Status == 1 {
			status = string(Available)
		} else {
			status = string(Unavailable)
		}
		return object, status, nil
	}
}

func (s *DdoscooService) DescribeDdoscooInstanceSpec(d *schema.ResourceData) (v ddoscoo.InstanceSpec, err error) {
	request := ddoscoo.CreateDescribeInstanceSpecsRequest()
	request.RegionId = "cn-hangzhou"
	id := d.Id()
	request.InstanceIds = &[]string{id}

	raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.DescribeInstanceSpecs(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound", "ddos_coop3301"}) {
			return v, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}

		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ddoscoo.DescribeInstanceSpecsResponse)
	if len(resp.InstanceSpecs) == 0 || resp.InstanceSpecs[0].InstanceId != id {
		return v, WrapErrorf(Error(GetNotFoundMessage("DdoscooInstanceSpec", id)), NotFoundMsg, ProviderERROR)
	}

	v = resp.InstanceSpecs[0]
	return v, WrapError(err)
}

func (s *DdoscooService) UpdateDdoscooInstanceName(instanceId string, name string) error {
	request := ddoscoo.CreateModifyInstanceRemarkRequest()
	request.RegionId = "cn-hangzhou"
	request.InstanceId = instanceId
	request.Remark = name

	raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.ModifyInstanceRemark(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *DdoscooService) UpdateInstanceSpec(schemaName string, specName string, d *schema.ResourceData, meta interface{}) error {
	request := bssopenapi.CreateModifyInstanceRequest()
	request.RegionId = s.client.RegionId
	if d.Get("product_type").(string) == "ddoscoo_intl" {
		request.RegionId = "ap-southeast-1"
	} else {
		request.RegionId = "cn-hangzhou"
	}
	request.InstanceId = d.Id()

	request.ProductCode = "ddos"
	request.ProductType = d.Get("product_type").(string)
	request.SubscriptionType = "Subscription"

	o, n := d.GetChange(schemaName)
	oi, _ := strconv.Atoi(o.(string))
	ni, _ := strconv.Atoi(n.(string))
	if ni < oi {
		request.ModifyType = "Downgrade"
	} else {
		request.ModifyType = "Upgrade"
	}

	request.Parameter = &[]bssopenapi.ModifyInstanceParameter{
		{
			Code:  specName,
			Value: d.Get(schemaName).(string),
		},
	}
	raw, err := s.client.WithBssopenapiClient(func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.ModifyInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*bssopenapi.ModifyInstanceResponse)
	if !response.Success {
		return WrapError(Error(response.Message))
	}
	return nil
}

func (s *DdoscooService) convertRulesToString(v []interface{}) (string, error) {
	arrayMaps := make([]ddoscoo.Rule, len(v))
	for i, vv := range v {
		item := vv.(map[string]interface{})
		arrayMaps[i] = ddoscoo.Rule{
			Priority:  item["priority"].(int),
			RegionId:  item["region_id"].(string),
			Status:    item["status"].(int),
			Type:      item["type"].(string),
			Value:     item["value"].(string),
			ValueType: item["value_type"].(int),
		}
	}
	maps, err := json.Marshal(arrayMaps)
	if err != nil {
		return "", WrapError(err)
	}
	return string(maps), nil
}

func (s *DdoscooService) DescribeDdoscooSchedulerRule(id string) (object ddoscoo.SchedulerRule, err error) {
	request := ddoscoo.CreateDescribeSchedulerRulesRequest()
	request.RegionId = s.client.RegionId

	request.RuleName = id
	request.PageSize = requests.NewInteger(10)

	raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.DescribeSchedulerRules(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ddoscoo.DescribeSchedulerRulesResponse)

	if len(response.SchedulerRules) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("DdoscooSchedulerRule", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return response.SchedulerRules[0], nil
}

func (s *DdoscooService) DescribeDdoscooDomainResource(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDdoscooClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDomainResource"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"Domain":     id,
		"PageNumber": 1,
		"PageSize":   10,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.WebRules", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.WebRules", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("DdosCoo", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["Domain"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("DdosCoo", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *DdoscooService) DescribeDdoscooPort(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDdoscooClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribePort"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"FrontendPort":     parts[1],
		"FrontendProtocol": parts[2],
		"InstanceId":       parts[0],
		"PageNumber":       1,
		"PageSize":         10,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.NetworkRules", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NetworkRules", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("DdosCoo", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["FrontendProtocol"].(string) != parts[2] {
			return object, WrapErrorf(Error(GetNotFoundMessage("DdosCoo", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

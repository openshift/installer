package alicloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type PvtzService struct {
	client *connectivity.AliyunClient
}

func (s *PvtzService) DescribePvtzZoneBasic(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPvtzClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeZoneInfo"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"ZoneId":   id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"Zone.Invalid.Id", "Zone.Invalid.UserId", "Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PvtzZone", id)), NotFoundMsg, ProviderERROR)
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

func (s *PvtzService) DescribePvtzZoneAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPvtzClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeZoneInfo"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"ZoneId":   id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"Zone.Invalid.Id", "Zone.Invalid.UserId", "Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PvtzZoneAttachment", id)), NotFoundMsg, ProviderERROR)
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

func (s *PvtzService) WaitForZoneAttachment(id string, vpcIdMap map[string]string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	var vpcId string
	for {
		object, err := s.DescribePvtzZoneAttachment(id)
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		equal := true
		vpcs := object["BindVpcs"].(map[string]interface{})["Vpc"].([]interface{})
		if len(vpcs) == len(vpcIdMap) {
			for _, vpc := range vpcs {
				vpc := vpc.(map[string]interface{})
				if _, ok := vpcIdMap[vpc["VpcId"].(string)]; !ok {
					equal = false
					vpcId = vpc["VpcId"].(string)
					break
				}
			}
		} else {
			equal = false
		}
		if equal {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", vpcId, ProviderERROR)
		}
	}
	return nil
}

func (s *PvtzService) DescribePvtzZoneRecord(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPvtzClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeZoneRecords"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"ZoneId":     parts[1],
		"PageNumber": 1,
		"PageSize":   20,
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"System.Busy", "ServiceUnavailable", "Throttling.User"}) {
					wait()
					return resource.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"Zone.Invalid.Id", "Zone.Invalid.UserId", "Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
					err = WrapErrorf(Error(GetNotFoundMessage("PvtzZoneRecord", id)), NotFoundMsg, ProviderERROR)
					return resource.NonRetryableError(err)
				}
				err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return
		}
		addDebug(action, response, request)
		v, err := jsonpath.Get("$.Records.Record", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Records.Record", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("PrivateZone", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(formatInt(v.(map[string]interface{})["RecordId"])) == parts[0] || fmt.Sprint(v.(map[string]interface{})["RecordId"]) == parts[0] {
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return
}

func (s *PvtzService) WaitForPvtzZone(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribePvtzZone(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object["ZoneId"] == id {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object["ZoneId"], id, ProviderERROR)
		}

	}
}

func (s *PvtzService) WaitForPvtzZoneAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribePvtzZoneAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object["ZoneId"] == id {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object["ZoneId"], id, ProviderERROR)
		}

	}
}

func (s *PvtzService) WaitForPvtzZoneRecord(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribePvtzZoneRecord(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strconv.FormatInt(object["RecordId"].(int64), 10) == parts[0] {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.FormatInt(object["RecordId"].(int64), 10), id, ProviderERROR)
		}

	}
}

func (s *PvtzService) DescribePvtzUserVpcAuthorization(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPvtzClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeUserVpcAuthorizations"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AuthType":         parts[1],
		"AuthorizedUserId": parts[0],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"System.Busy"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("PrivateZone:UserVpcAuthorization", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Users", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Users", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("PrivateZone", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["AuthType"]) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("PrivateZone", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *PvtzService) DescribePvtzEndpoint(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPvtzClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeResolverEndpoint"
	request := map[string]interface{}{
		"EndpointId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"ResolverEndpoint.NotExists"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("PrivateZone:Endpoint", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *PvtzService) PvtzEndpointStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePvtzEndpoint(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *PvtzService) DescribePvtzRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPvtzClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeResolverRule"
	request := map[string]interface{}{
		"RuleId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"ResolverRule.NotExists"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("PrivateZone:Rule", id)), NotFoundMsg, ProviderERROR)
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

func (s *PvtzService) DescribePvtzRuleAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPvtzClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeResolverRule"
	request := map[string]interface{}{
		"RuleId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"ResolverRule.NotExists"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("PrivateZone:RuleAttachment", id)), NotFoundMsg, ProviderERROR)
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

func (s *PvtzService) DescribeSyncEcsHostTask(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPvtzClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeSyncEcsHostTask"
	request := map[string]interface{}{
		"ZoneId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"System.Busy", "Ecs.SyncTask.NotExists", "ServiceUnavailable", "Throttling.User"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PvtzZone", id)), NotFoundMsg, ProviderERROR)
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

func (s *PvtzService) DescribePvtzZone(id string) (object map[string]interface{}, err error) {
	object, err = s.DescribePvtzZoneBasic(id)
	if err != nil {
		return nil, WrapError(err)
	}
	syncObj, err := s.DescribeSyncEcsHostTask(id)
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pvtz_zone pvtzService.DescribeSyncEcsHostTask Failed!!! %s", err)
			return object, nil
		}
		return nil, WrapError(err)
	}
	object["SyncHostTask"] = syncObj
	return object, nil
}

package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ClickhouseService struct {
	client *connectivity.AliyunClient
}

func (s *ClickhouseService) DescribeClickHouseDbCluster(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewClickhouseClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDBClusterAttribute"
	request := map[string]interface{}{
		"DBClusterId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DBCluster", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DBCluster", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ClickhouseService) ClickHouseDbClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeClickHouseDbCluster(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["DBClusterStatus"].(string) == failState {
				return object, object["DBClusterStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["DBClusterStatus"].(string)))
			}
		}
		return object, object["DBClusterStatus"].(string), nil
	}
}

func (s *ClickhouseService) DescribeClickHouseAccount(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewClickhouseClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeAccounts"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AccountName": parts[1],
		"DBClusterId": parts[0],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectAccountStatus", "IncorrectDBInstanceState"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("ClickHouse", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Accounts.Account", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Accounts.Account", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ClickHouse", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["AccountName"]) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("ClickHouse", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ClickhouseService) ClickhouseStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeClickHouseAccount(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		if object != nil {
			for _, failState := range failStates {
				if object["AccountStatus"].(string) == failState {
					return object, object["AccountStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["AccountStatus"].(string)))
				}
			}
			return object, object["AccountStatus"].(string), nil
		}

		return nil, "", WrapErrorf(Error(GetNotFoundMessage("ClickHouse:Account", id)), NotFoundMsg, ProviderERROR)
	}
}

func (s *ClickhouseService) DescribeDBClusterAccessWhiteList(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewClickhouseClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDBClusterAccessWhiteList"
	request := map[string]interface{}{
		"DBClusterId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DBClusterAccessWhiteList.IPArray", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ClickhouseService) DescribeClickHouseBackupPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewClickhouseClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeBackupPolicy"
	request := map[string]interface{}{
		"DBClusterId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("ClickHouse", id)), NotFoundWithResponse, response)
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

func (s *ClickhouseService) ClickHouseBackupPolicyStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeClickHouseBackupPolicy(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Switch"]) == failState {
				return object, fmt.Sprint(object["Switch"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Switch"])))
			}
		}
		return object, fmt.Sprint(object["Switch"]), nil
	}
}

package alicloud

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type HbrService struct {
	client *connectivity.AliyunClient
}

func (s *HbrService) convertDetailToString(details []interface{}) (string, error) {
	configRulesMap := make(map[string]interface{})
	for _, configRules := range details {
		configRulesArg := configRules.(map[string]interface{})
		configRulesMap["enableWriters"] = true
		configRulesMap["appConsistent"] = configRulesArg["app_consistent"]
		configRulesMap["snapshotGroup"] = configRulesArg["snapshot_group"]
		configRulesMap["EnableFsFreeze"] = configRulesArg["enable_fs_freeze"]
		configRulesMap["preScriptPath"] = configRulesArg["pre_script_path"]
		configRulesMap["postScriptPath"] = configRulesArg["post_script_path"]
		configRulesMap["timeoutInSeconds"] = configRulesArg["timeout_in_seconds"]
		configRulesMap["doCopy"] = configRulesArg["do_copy"]

		if configRulesArg["destination_region_id"] != "" {
			configRulesMap["destinationRegionId"] = configRulesArg["destination_region_id"]
		}
		configRulesMap["destinationRetention"] = configRulesArg["destination_retention"]
		diskListArg := make([]interface{}, 0)
		if configRulesArg["disk_id_list"] != nil {
			diskListArg = append(diskListArg, configRulesArg["disk_id_list"].([]interface{})...)
			configRulesMap["diskIdList"] = diskListArg
		}
	}
	if v, err := convertArrayObjectToJsonString(configRulesMap); err != nil {
		return "", WrapError(err)
	} else {
		return v, nil
	}
}

func (s *HbrService) DescribeHbrVault(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHbrClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVaults"
	request := map[string]interface{}{
		"VaultRegionId": s.client.RegionId,
		"VaultId":       id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Vaults.Vault", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Vaults.Vault", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["VaultId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *HbrService) HbrVaultStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHbrVault(id)
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

func (s *HbrService) DescribeHbrEcsBackupPlan(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHbrClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeBackupPlans"
	request := map[string]interface{}{
		"SourceType": "ECS_FILE",
	}
	filtersMapList := make([]map[string]interface{}, 0)
	filtersMapList = append(filtersMapList, map[string]interface{}{
		"Key":    "planId",
		"Values": []string{id},
	})
	request["Filters"] = filtersMapList
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.BackupPlans.BackupPlan", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.BackupPlans.BackupPlan", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["PlanId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *HbrService) DescribeHbrNasBackupPlan(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHbrClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeBackupPlans"
	request := map[string]interface{}{
		"SourceType": "NAS",
	}
	filtersMapList := make([]map[string]interface{}, 0)
	filtersMapList = append(filtersMapList, map[string]interface{}{
		"Key":    "planId",
		"Values": []string{id},
	})
	request["Filters"] = filtersMapList
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.BackupPlans.BackupPlan", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.BackupPlans.BackupPlan", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["PlanId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *HbrService) DescribeHbrOssBackupPlan(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHbrClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeBackupPlans"
	request := map[string]interface{}{
		"SourceType": "OSS",
		"PageNumber": 1,
		"PageSize":   100,
	}
	filtersMapList := make([]map[string]interface{}, 0)
	request["Filters"] = filtersMapList
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
		if fmt.Sprint(response["Success"]) == "false" {
			return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		v, err := jsonpath.Get("$.BackupPlans.BackupPlan", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.BackupPlans.BackupPlan", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["PlanId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
	}
	return
}

type InstallClientTaskStatus struct {
	InstanceId   string `json:"instanceId"`
	ClientId     string `json:"clientId"`
	ClientStatus string `json:"clientStatus"`
	ErrorCode    string `json:"errorCode"`
}

func (s *HbrService) DescribeHbrTask(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHbrClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeTask"
	request := map[string]interface{}{
		"TaskId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	if fmt.Sprint(response["Description"]) == "cancelled" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Result", response)
	}

	var taskStatus []InstallClientTaskStatus
	resultJson := v.(string)
	if resultJson == "" {
		return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
	}

	err = json.Unmarshal([]byte(resultJson), &taskStatus)

	if len(taskStatus) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
	}

	object = make(map[string]interface{})
	ty := reflect.TypeOf(taskStatus[0])
	tv := reflect.ValueOf(taskStatus[0])
	for i := 0; i < tv.NumField(); i++ {
		object[ty.Field(i).Name] = tv.Field(i).Interface()
	}
	return object, nil
}

func (s *HbrService) HbrTaskRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHbrTask(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["ClientStatus"]) == failState {
				return object, fmt.Sprint(object["ClientStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["ClientStatus"])))
			}
		}
		return object, fmt.Sprint(object["ClientStatus"]), nil
	}
}

func (s *HbrService) DescribeHbrEcsBackupClient(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHbrClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeBackupClients"
	request := map[string]interface{}{
		"ClientIds":  "[\"" + id + "\"]",
		"ClientType": "ECS_CLIENT",
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Clients", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Clients", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["ClientId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *HbrService) HbrEcsBackupClientStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHbrEcsBackupClient(id)
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
func (s *HbrService) DescribeHbrRestoreJob(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHbrClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeRestoreJobs2"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RestoreType": parts[1],
		"PageSize":    50,
		"PageNumber":  1,
	}
	filtersMapList := make([]map[string]interface{}, 0)
	request["Filters"] = filtersMapList
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
		if fmt.Sprint(response["Success"]) == "false" {
			return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		v, err := jsonpath.Get("$.RestoreJobs.RestoreJob", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.RestoreJobs.RestoreJob", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["RestoreId"]) == parts[0] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *HbrService) HbrRestoreJobStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHbrRestoreJob(id)
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
func (s *HbrService) DescribeHbrServerBackupPlan(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHbrClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeBackupPlans"
	request := map[string]interface{}{
		"SourceType": "UDM_ECS",
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
		if fmt.Sprint(response["Success"]) == "false" {
			return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		v, err := jsonpath.Get("$.BackupPlans.BackupPlan", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.BackupPlans.BackupPlan", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["PlanId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("HBR", id)), NotFoundWithResponse, response)
	}
	return
}

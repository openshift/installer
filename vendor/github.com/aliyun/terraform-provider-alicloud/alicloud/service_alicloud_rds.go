package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type RdsService struct {
	client *connectivity.AliyunClient
}

//
//       _______________                      _______________                       _______________
//       |              | ______param______\  |              |  _____request_____\  |              |
//       |   Business   |                     |    Service   |                      |    SDK/API   |
//       |              | __________________  |              |  __________________  |              |
//       |______________| \    (obj, err)     |______________|  \ (status, cont)    |______________|
//                           |                                    |
//                           |A. {instance, nil}                  |a. {200, content}
//                           |B. {nil, error}                     |b. {200, nil}
//                      					  |c. {4xx, nil}
//
// The API return 200 for resource not found.
// When getInstance is empty, then throw InstanceNotfound error.
// That the business layer only need to check error.
var DBInstanceStatusCatcher = Catcher{"OperationDenied.DBInstanceStatus", 60, 5}

func (s *RdsService) DescribeDBInstance(id string) (map[string]interface{}, error) {
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDBInstanceAttribute"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	var response map[string]interface{}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Items.DBInstanceAttribute", response)
	if err != nil {
		return nil, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.DBInstanceAttribute", response)
	}
	if len(v.([]interface{})) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBAccount", id)), NotFoundMsg, ProviderERROR)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *RdsService) DescribeTasks(id string) (map[string]interface{}, error) {
	action := "DescribeTasks"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, ProviderERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return response, nil
}

func (s *RdsService) DescribeDBReadonlyInstance(id string) (map[string]interface{}, error) {
	action := "DescribeDBInstanceAttribute"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	dBInstanceAttributes := response["Items"].(map[string]interface{})["DBInstanceAttribute"].([]interface{})
	if len(dBInstanceAttributes) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBInstance", id)), NotFoundMsg, ProviderERROR)
	}

	return dBInstanceAttributes[0].(map[string]interface{}), nil
}

func (s *RdsService) DescribeDBAccountPrivilege(id string) (map[string]interface{}, error) {
	var ds map[string]interface{}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return ds, WrapError(err)
	}
	action := "DescribeAccounts"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": parts[0],
		"AccountName":  parts[1],
		"SourceIp":     s.client.SourceIp,
	}
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	invoker := NewInvoker()
	invoker.AddCatcher(DBInstanceStatusCatcher)
	var response map[string]interface{}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	if err := invoker.Run(func() error {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		return nil
	}); err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return ds, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return ds, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	dBInstanceAccounts := response["Accounts"].(map[string]interface{})["DBInstanceAccount"].([]interface{})
	if len(dBInstanceAccounts) < 1 {
		return ds, WrapErrorf(Error(GetNotFoundMessage("DBAccountPrivilege", id)), NotFoundMsg, ProviderERROR)
	}
	return dBInstanceAccounts[0].(map[string]interface{}), nil
}

func (s *RdsService) DescribeDBDatabase(id string) (map[string]interface{}, error) {
	var ds map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return ds, WrapError(err)
	}
	dbName := parts[1]
	var response map[string]interface{}
	action := "DescribeDatabases"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": parts[0],
		"DBName":       dbName,
		"SourceIp":     s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError", "OperationDenied.DBInstanceStatus"}) {
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR))
			}
			if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidDBName.NotFound", "InvalidDBInstanceId.NotFoundError"}) {
				return resource.NonRetryableError(WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR))
		}
		addDebug(action, response, request)
		v, err := jsonpath.Get("$.Databases.Database", response)
		if err != nil {
			return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.Databases.Database", response))
		}
		if len(v.([]interface{})) < 1 {
			return resource.NonRetryableError(WrapErrorf(Error(GetNotFoundMessage("DBDatabase", dbName)), NotFoundMsg, ProviderERROR))
		}
		ds = v.([]interface{})[0].(map[string]interface{})
		return nil
	})
	return ds, err
}

func (s *RdsService) DescribeParameters(id string) (map[string]interface{}, error) {
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeParameters"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return response, err
}

func (s *RdsService) RefreshParameters(d *schema.ResourceData, attribute string) error {
	var param []map[string]interface{}
	documented, ok := d.GetOk(attribute)
	if !ok {
		d.Set(attribute, param)
		return nil
	}
	object, err := s.DescribeParameters(d.Id())
	if err != nil {
		return WrapError(err)
	}

	var parameters = make(map[string]interface{})
	dBInstanceParameters := object["RunningParameters"].(map[string]interface{})["DBInstanceParameter"].([]interface{})
	for _, i := range dBInstanceParameters {
		i := i.(map[string]interface{})
		if i["ParameterName"] != "" {
			parameter := map[string]interface{}{
				"name":  i["ParameterName"],
				"value": i["ParameterValue"],
			}
			parameters[i["ParameterName"].(string)] = parameter
		}
	}
	dBInstanceParameters = object["ConfigParameters"].(map[string]interface{})["DBInstanceParameter"].([]interface{})
	for _, i := range dBInstanceParameters {
		i := i.(map[string]interface{})
		if i["ParameterName"] != "" {
			parameter := map[string]interface{}{
				"name":  i["ParameterName"],
				"value": i["ParameterValue"],
			}
			parameters[i["ParameterName"].(string)] = parameter
		}
	}

	for _, parameter := range documented.(*schema.Set).List() {
		name := parameter.(map[string]interface{})["name"]
		for _, value := range parameters {
			if value.(map[string]interface{})["name"] == name {
				param = append(param, value.(map[string]interface{}))
				break
			}
		}
	}
	if err := d.Set(attribute, param); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *RdsService) ModifyParameters(d *schema.ResourceData, attribute string) error {
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	action := "ModifyParameter"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": d.Id(),
		"Forcerestart": d.Get("force_restart"),
		"SourceIp":     s.client.SourceIp,
	}
	config := make(map[string]string)
	allConfig := make(map[string]string)
	o, n := d.GetChange(attribute)
	os, ns := o.(*schema.Set), n.(*schema.Set)
	add := ns.Difference(os).List()
	if len(add) > 0 {
		for _, i := range add {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			config[key] = value
		}
		cfg, _ := json.Marshal(config)
		request["Parameters"] = string(cfg)
		// wait instance status is Normal before modifying
		if err := s.WaitForDBInstance(d.Id(), Running, DefaultLongTimeout); err != nil {
			return WrapError(err)
		}
		// Need to check whether some parameter needs restart
		if !d.Get("force_restart").(bool) {
			action := "DescribeParameterTemplates"
			request := map[string]interface{}{
				"RegionId":      s.client.RegionId,
				"DBInstanceId":  d.Id(),
				"Engine":        d.Get("engine"),
				"EngineVersion": d.Get("engine_version"),
				"ClientToken":   buildClientToken(action),
				"SourceIp":      s.client.SourceIp,
			}
			forceRestartMap := make(map[string]string)
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			addDebug(action, response, request)
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, s.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			templateRecords := response["Parameters"].(map[string]interface{})["TemplateRecord"].([]interface{})
			for _, para := range templateRecords {
				para := para.(map[string]interface{})
				if para["ForceRestart"] == "true" {
					forceRestartMap[para["ParameterName"].(string)] = para["ForceRestart"].(string)
				}
			}
			if len(forceRestartMap) > 0 {
				for key, _ := range config {
					if _, ok := forceRestartMap[key]; ok {
						return WrapError(fmt.Errorf("Modifying RDS instance's parameter '%s' requires setting 'force_restart = true'.", key))
					}
				}
			}
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		// wait instance parameter expect after modifying
		for _, i := range ns.List() {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			allConfig[key] = value
		}
		if err := s.WaitForDBParameter(d.Id(), DefaultTimeoutMedium, allConfig); err != nil {
			return WrapError(err)
		}
		// wait instance status is Normal after modifying
		if err := s.WaitForDBInstance(d.Id(), Running, DefaultLongTimeout); err != nil {
			return WrapError(err)
		}
	}
	d.SetPartial(attribute)
	return nil
}

func (s *RdsService) DescribeDBInstanceNetInfo(id string) ([]interface{}, error) {
	action := "DescribeDBInstanceNetInfo"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	var response map[string]interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
				return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
	}
	dBInstanceNetInfos := response["DBInstanceNetInfos"].(map[string]interface{})["DBInstanceNetInfo"].([]interface{})
	if len(dBInstanceNetInfos) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBInstanceNetInfo", id)), NotFoundMsg, ProviderERROR)
	}

	return dBInstanceNetInfos, nil
}

func (s *RdsService) DescribeDBConnection(id string) (map[string]interface{}, error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	object, err := s.DescribeDBInstanceNetInfo(parts[0])

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidCurrentConnectionString.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapError(err)
	}

	if object != nil {
		for _, o := range object {
			o := o.(map[string]interface{})
			if strings.HasPrefix(o["ConnectionString"].(string), parts[1]) {
				return o, nil
			}
		}
	}

	return nil, WrapErrorf(Error(GetNotFoundMessage("DBConnection", id)), NotFoundMsg, ProviderERROR)
}
func (s *RdsService) DescribeDBReadWriteSplittingConnection(id string) (map[string]interface{}, error) {
	object, err := s.DescribeDBInstanceNetInfo(id)
	if err != nil && !NotFoundError(err) {
		return nil, err
	}

	if object != nil {
		for _, conn := range object {
			conn := conn.(map[string]interface{})
			if conn["ConnectionStringType"] != "ReadWriteSplitting" {
				continue
			}
			if conn["MaxDelayTime"] == nil {
				continue
			}
			if _, err := strconv.Atoi(conn["MaxDelayTime"].(string)); err != nil {
				return nil, err
			}
			return conn, nil
		}
	}

	return nil, WrapErrorf(Error(GetNotFoundMessage("ReadWriteSplittingConnection", id)), NotFoundMsg, ProviderERROR)
}

func (s *RdsService) GrantAccountPrivilege(id, dbName string) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	action := "GrantAccountPrivilege"
	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"DBInstanceId":     parts[0],
		"AccountName":      parts[1],
		"DBName":           dbName,
		"AccountPrivilege": parts[2],
		"SourceIp":         s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	var response map[string]interface{}
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) || IsExpectedErrors(err, []string{"InvalidDB.NotFound"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if err := s.WaitForAccountPrivilege(id, dbName, Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return nil
}

func (s *RdsService) RevokeAccountPrivilege(id, dbName string) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	action := "RevokeAccountPrivilege"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": parts[0],
		"AccountName":  parts[1],
		"DBName":       dbName,
		"SourceIp":     s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) || NeedRetry(err) {
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"InvalidDB.NotFound"}) {
				log.Printf("[WARN] Resource alicloud_db_account_privilege RevokeAccountPrivilege Failed!!! %s", err)
				return nil
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	if err := s.WaitForAccountPrivilegeRevoked(id, dbName, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return nil
}

func (s *RdsService) ReleaseDBPublicConnection(instanceId, connection string) error {
	action := "ReleaseInstancePublicConnection"
	request := map[string]interface{}{
		"RegionId":                s.client.RegionId,
		"DBInstanceId":            instanceId,
		"CurrentConnectionString": connection,
		"SourceIp":                s.client.SourceIp,
	}
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return nil
}

func (s *RdsService) ModifyDBBackupPolicy(d *schema.ResourceData, updateForData, updateForLog bool) error {
	enableBackupLog := "1"

	backupPeriod := ""
	if v, ok := d.GetOk("preferred_backup_period"); ok && v.(*schema.Set).Len() > 0 {
		periodList := expandStringList(v.(*schema.Set).List())
		backupPeriod = fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
	} else {
		periodList := expandStringList(d.Get("backup_period").(*schema.Set).List())
		backupPeriod = fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
	}

	backupTime := "02:00Z-03:00Z"
	if v, ok := d.GetOk("preferred_backup_time"); ok && v.(string) != "02:00Z-03:00Z" {
		backupTime = v.(string)
	} else if v, ok := d.GetOk("backup_time"); ok && v.(string) != "" {
		backupTime = v.(string)
	}

	retentionPeriod := "7"
	if v, ok := d.GetOk("backup_retention_period"); ok && v.(int) != 7 {
		retentionPeriod = strconv.Itoa(v.(int))
	} else if v, ok := d.GetOk("retention_period"); ok && v.(int) != 0 {
		retentionPeriod = strconv.Itoa(v.(int))
	}

	logBackupRetentionPeriod := ""
	if v, ok := d.GetOk("log_backup_retention_period"); ok && v.(int) != 0 {
		logBackupRetentionPeriod = strconv.Itoa(v.(int))
	} else if v, ok := d.GetOk("log_retention_period"); ok && v.(int) != 0 {
		logBackupRetentionPeriod = strconv.Itoa(v.(int))
	}

	localLogRetentionHours := ""
	if v, ok := d.GetOk("local_log_retention_hours"); ok {
		localLogRetentionHours = strconv.Itoa(v.(int))
	}

	localLogRetentionSpace := ""
	if v, ok := d.GetOk("local_log_retention_space"); ok {
		localLogRetentionSpace = strconv.Itoa(v.(int))
	}

	highSpaceUsageProtection := d.Get("high_space_usage_protection").(string)

	if !d.Get("enable_backup_log").(bool) {
		enableBackupLog = "0"
	}

	if d.HasChange("log_backup_retention_period") {
		if d.Get("log_backup_retention_period").(int) > d.Get("backup_retention_period").(int) {
			logBackupRetentionPeriod = retentionPeriod
		}
	}

	logBackupFrequency := ""
	if v, ok := d.GetOk("log_backup_frequency"); ok {
		logBackupFrequency = v.(string)
	}
	compressType := ""
	if v, ok := d.GetOk("compress_type"); ok {
		compressType = v.(string)
	}

	archiveBackupRetentionPeriod := "0"
	if v, ok := d.GetOk("archive_backup_retention_period"); ok {
		archiveBackupRetentionPeriod = strconv.Itoa(v.(int))
	}

	archiveBackupKeepCount := "1"
	if v, ok := d.GetOk("archive_backup_keep_count"); ok {
		archiveBackupKeepCount = strconv.Itoa(v.(int))
	}

	archiveBackupKeepPolicy := "0"
	if v, ok := d.GetOk("archive_backup_keep_policy"); ok {
		archiveBackupKeepPolicy = v.(string)
	}

	instance, err := s.DescribeDBInstance(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if updateForData {
		conn, err := s.client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		action := "ModifyBackupPolicy"
		request := map[string]interface{}{
			"RegionId":              s.client.RegionId,
			"DBInstanceId":          d.Id(),
			"PreferredBackupPeriod": backupPeriod,
			"PreferredBackupTime":   backupTime,
			"BackupRetentionPeriod": retentionPeriod,
			"CompressType":          compressType,
			"BackupPolicyMode":      "DataBackupPolicy",
			"SourceIp":              s.client.SourceIp,
		}
		if instance["Engine"] == "SQLServer" && logBackupFrequency == "LogInterval" {
			request["LogBackupFrequency"] = logBackupFrequency
		}
		if instance["Engine"] == "MySQL" && instance["DBInstanceStorageType"] == "local_ssd" {
			request["ArchiveBackupRetentionPeriod"] = archiveBackupRetentionPeriod
			request["ArchiveBackupKeepCount"] = archiveBackupKeepCount
			request["ArchiveBackupKeepPolicy"] = archiveBackupKeepPolicy
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		if err := s.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
	}

	// At present, the sql server database does not support setting logBackupRetentionPeriod
	if updateForLog && instance["Engine"] != "SQLServer" {
		conn, err := s.client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		action := "ModifyBackupPolicy"
		request := map[string]interface{}{
			"RegionId":                 s.client.RegionId,
			"DBInstanceId":             d.Id(),
			"EnableBackupLog":          enableBackupLog,
			"LocalLogRetentionHours":   localLogRetentionHours,
			"LocalLogRetentionSpace":   localLogRetentionSpace,
			"HighSpaceUsageProtection": highSpaceUsageProtection,
			"BackupPolicyMode":         "LogBackupPolicy",
			"LogBackupRetentionPeriod": logBackupRetentionPeriod,
			"SourceIp":                 s.client.SourceIp,
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		if err := s.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func (s *RdsService) ModifyDBSecurityIps(instanceId, ips string) error {
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	action := "ModifySecurityIps"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": instanceId,
		"SecurityIps":  ips,
		"SourceIp":     s.client.SourceIp,
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	if err := s.WaitForDBInstance(instanceId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *RdsService) DescribeDBSecurityIps(instanceId string) ([]interface{}, error) {
	action := "DescribeDBInstanceIPArrayList"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": instanceId,
		"SourceIp":     s.client.SourceIp,
	}
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, instanceId, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return response["Items"].(map[string]interface{})["DBInstanceIPArray"].([]interface{}), nil
}

func (s *RdsService) GetSecurityIps(instanceId string) ([]string, error) {
	object, err := s.DescribeDBSecurityIps(instanceId)
	if err != nil {
		return nil, WrapError(err)
	}

	var ips, separator string
	ipsMap := make(map[string]string)
	for _, ip := range object {
		ip := ip.(map[string]interface{})
		if ip["DBInstanceIPArrayAttribute"] == "hidden" {
			continue
		}
		ips += separator + ip["SecurityIPList"].(string)
		separator = COMMA_SEPARATED
	}

	for _, ip := range strings.Split(ips, COMMA_SEPARATED) {
		ipsMap[ip] = ip
	}

	var finalIps []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			finalIps = append(finalIps, key)
		}
	}

	return finalIps, nil
}

func (s *RdsService) DescribeSecurityGroupConfiguration(id string) ([]string, error) {
	action := "DescribeSecurityGroupConfiguration"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	groupIds := make([]string, 0)
	ecsSecurityGroupRelations := response["Items"].(map[string]interface{})["EcsSecurityGroupRelation"].([]interface{})
	for _, v := range ecsSecurityGroupRelations {
		v := v.(map[string]interface{})
		groupIds = append(groupIds, v["SecurityGroupId"].(string))
	}
	return groupIds, nil
}

func (s *RdsService) DescribeDBInstanceSSL(id string) (map[string]interface{}, error) {
	action := "DescribeDBInstanceSSL"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return response, nil
}

func (s *RdsService) DescribeDBInstanceEncryptionKey(id string) (map[string]interface{}, error) {
	action := "DescribeDBInstanceEncryptionKey"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return response, nil
}

func (s *RdsService) DescribeHASwitchConfig(id string) (map[string]interface{}, error) {
	action := "DescribeHASwitchConfig"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return response, nil
}

func (s *RdsService) DescribeRdsTDEInfo(id string) (map[string]interface{}, error) {
	action := "DescribeDBInstanceTDE"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	statErr := s.WaitForDBInstance(id, Running, DefaultLongTimeout)
	if statErr != nil {
		return nil, WrapError(statErr)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return response, nil
}

func (s *RdsService) ModifySecurityGroupConfiguration(id string, groupid string) error {
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	action := "ModifySecurityGroupConfiguration"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	//openapi required that input "Empty" if groupid is ""
	if len(groupid) == 0 {
		groupid = "Empty"
	}
	request["SecurityGroupId"] = groupid
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return nil
}

// return multiIZ list of current region
func (s *RdsService) DescribeMultiIZByRegion() (izs []string, err error) {
	action := "DescribeRegions"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"SourceIp": s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "MultiIZByRegion", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	regions := response["Regions"].(map[string]interface{})["RDSRegion"].([]interface{})
	zoneIds := []string{}
	for _, r := range regions {
		r := r.(map[string]interface{})
		if r["RegionId"] == string(s.client.Region) && strings.Contains(r["ZoneId"].(string), MULTI_IZ_SYMBOL) {
			zoneIds = append(zoneIds, r["ZoneId"].(string))
		}
	}

	return zoneIds, nil
}

func (s *RdsService) DescribeBackupPolicy(id string) (map[string]interface{}, error) {
	action := "DescribeBackupPolicy"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return response, nil
}

func (s *RdsService) DescribeDbInstanceMonitor(id string) (monitoringPeriod int, err error) {
	action := "DescribeDBInstanceMonitor"
	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
		"SourceIp":     s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return 0, WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return 0, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	monPeriod, err := strconv.Atoi(response["Period"].(string))
	if err != nil {
		return 0, WrapError(err)
	}
	return monPeriod, nil
}

func (s *RdsService) DescribeSQLCollectorPolicy(id string) (map[string]interface{}, error) {
	action := "DescribeSQLCollectorPolicy"
	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
		"SourceIp":     s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return response, nil
}

func (s *RdsService) DescribeSQLCollectorRetention(id string) (map[string]interface{}, error) {
	action := "DescribeSQLCollectorRetention"
	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
		"SourceIp":     s.client.SourceIp,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return response, nil
}

// WaitForInstance waits for instance to given status
func (s *RdsService) WaitForDBInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBInstance(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object != nil && strings.ToLower(object["DBInstanceStatus"].(string)) == strings.ToLower(string(status)) {
			break
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object["DBInstanceStatus"], status, ProviderERROR)
		}
	}
	return nil
}

func (s *RdsService) RdsDBInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDBInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["DBInstanceStatus"] == failState {
				return object, object["DBInstanceStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["DBInstanceStatus"]))
			}
		}
		return object, object["DBInstanceStatus"].(string), nil
	}
}

func (s *RdsService) RdsTaskStateRefreshFunc(id string, taskAction string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeTasks(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		taskProgressInfos := object["Items"].(map[string]interface{})["TaskProgressInfo"].([]interface{})
		for _, t := range taskProgressInfos {
			t := t.(map[string]interface{})
			if t["TaskAction"] == taskAction {
				return object, t["Status"].(string), nil
			}
		}

		return object, "Pending", nil
	}
}

// WaitForDBParameter waits for instance parameter to given value.
// Status of DB instance is Running after ModifyParameters API was
// call, so we can not just wait for instance status become
// Running, we should wait until parameters have expected values.
func (s *RdsService) WaitForDBParameter(instanceId string, timeout int, expects map[string]string) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeParameters(instanceId)
		if err != nil {
			return WrapError(err)
		}
		var actuals = make(map[string]string)
		dBInstanceParameters := object["RunningParameters"].(map[string]interface{})["DBInstanceParameter"].([]interface{})
		for _, i := range dBInstanceParameters {
			i := i.(map[string]interface{})
			if i["ParameterName"] == nil || i["ParameterValue"] == nil {
				continue
			}
			actuals[i["ParameterName"].(string)] = i["ParameterValue"].(string)
		}
		dBInstanceParameters = object["ConfigParameters"].(map[string]interface{})["DBInstanceParameter"].([]interface{})
		for _, i := range dBInstanceParameters {
			i := i.(map[string]interface{})
			if i["ParameterName"] == nil || i["ParameterValue"] == nil {
				continue
			}
			actuals[i["ParameterName"].(string)] = i["ParameterValue"].(string)
		}

		match := true

		got_value := ""
		expected_value := ""

		for name, expect := range expects {
			if actual, ok := actuals[name]; ok {
				if expect != actual {
					match = false
					got_value = actual
					expected_value = expect
					break
				}
			} else {
				match = false
			}
		}

		if match {
			break
		}

		time.Sleep(DefaultIntervalShort * time.Second)

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, instanceId, GetFunc(1), timeout, got_value, expected_value, ProviderERROR)
		}
	}
	return nil
}

func (s *RdsService) WaitForDBConnection(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBConnection(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object != nil && object["ConnectionString"] != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object["ConnectionString"], id, ProviderERROR)
		}
	}
}

func (s *RdsService) WaitForDBReadWriteSplitting(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBReadWriteSplittingConnection(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if err == nil {
			break
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object["ConnectionString"], id, ProviderERROR)
		}
	}
	return nil
}

func (s *RdsService) WaitForAccountPrivilege(id, dbName string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBDatabase(parts[0] + ":" + dbName)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		ready := false
		if object != nil {
			accountPrivilegeInfos := object["Accounts"].(map[string]interface{})["AccountPrivilegeInfo"].([]interface{})
			for _, account := range accountPrivilegeInfos {
				// At present, postgresql response has a bug, DBOwner will be changed to ALL
				account := account.(map[string]interface{})
				if account["Account"] == parts[1] && (account["AccountPrivilege"] == parts[2] || (parts[2] == "DBOwner" && account["AccountPrivilege"] == "ALL")) {
					ready = true
					break
				}
			}
		}
		if status == Deleted && !ready {
			break
		}
		if ready {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", id, ProviderERROR)
		}
	}
	return nil
}

func (s *RdsService) WaitForAccountPrivilegeRevoked(id, dbName string, timeout int) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBDatabase(parts[0] + ":" + dbName)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
		}

		exist := false
		if object != nil {
			accountPrivilegeInfo := object["Accounts"].(map[string]interface{})["AccountPrivilegeInfo"].([]interface{})
			for _, account := range accountPrivilegeInfo {
				account := account.(map[string]interface{})
				if account["Account"] == parts[1] && (account["AccountPrivilege"] == parts[2] || (parts[2] == "DBOwner" && account["AccountPrivilege"] == "ALL")) {
					exist = true
					break
				}
			}
		}

		if !exist {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", dbName, ProviderERROR)
		}

	}
	return nil
}

func (s *RdsService) WaitForDBDatabase(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBDatabase(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object != nil && object["DBName"] == parts[1] {
			break
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object["DBName"], parts[1], ProviderERROR)
		}
	}
	return nil
}

// turn period to TimeType
func (s *RdsService) TransformPeriod2Time(period int, chargeType string) (ut int, tt common.TimeType) {
	if chargeType == string(Postpaid) {
		return 1, common.Day
	}

	if period >= 1 && period <= 9 {
		return period, common.Month
	}

	if period == 12 {
		return 1, common.Year
	}

	if period == 24 {
		return 2, common.Year
	}
	return 0, common.Day

}

// turn TimeType to Period
func (s *RdsService) TransformTime2Period(ut int, tt common.TimeType) (period int) {
	if tt == common.Year {
		return 12 * ut
	}

	return ut

}

func (s *RdsService) flattenDBSecurityIPs(list []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		i := i.(map[string]interface{})
		l := map[string]interface{}{
			"security_ips": i["SecurityIPList"],
		}
		result = append(result, l)
	}
	return result
}

func (s *RdsService) setInstanceTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		remove, add := diffRdsTags(o, n)

		if len(remove) > 0 {
			conn, err := s.client.NewRdsClient()
			if err != nil {
				return WrapError(err)
			}
			action := "UntagResources"
			request := map[string]interface{}{
				"ResourceId":   &[]string{d.Id()},
				"ResourceType": "INSTANCE",
				"TagKey":       &remove,
				"RegionId":     s.client.RegionId,
				"SourceIp":     s.client.SourceIp,
			}

			wait := incrementalWait(1*time.Second, 2*time.Second)
			runtime := util.RuntimeOptions{}
			err = resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
				if err != nil {
					if IsThrottling(err) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		if len(add) > 0 {
			conn, err := s.client.NewRdsClient()
			if err != nil {
				return WrapError(err)
			}
			action := "TagResources"
			request := map[string]interface{}{
				"ResourceId":   &[]string{d.Id()},
				"Tag":          &add,
				"ResourceType": "INSTANCE",
				"RegionId":     s.client.RegionId,
				"SourceIp":     s.client.SourceIp,
			}
			wait := incrementalWait(1*time.Second, 2*time.Second)
			runtime := util.RuntimeOptions{}
			err = resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
				if err != nil {
					if IsThrottling(err) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		d.SetPartial("tags")
	}

	return nil
}

func (s *RdsService) describeTags(d *schema.ResourceData) (tags []Tag, err error) {
	action := "DescribeTags"
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
		"RegionId":     s.client.RegionId,
		"SourceIp":     s.client.SourceIp,
	}
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	return s.respToTags(response["Items"].(map[string]interface{})["TagInfos"].([]interface{})), nil
}

func (s *RdsService) respToTags(tagSet []interface{}) (tags []Tag) {
	result := make([]Tag, 0, len(tagSet))
	for _, t := range tagSet {
		t := t.(map[string]interface{})
		tag := Tag{
			Key:   t["TagKey"].(string),
			Value: t["TagValue"].(string),
		}
		result = append(result, tag)
	}

	return result
}

func (s *RdsService) tagsToMap(tags []Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}

	return result
}

func (s *RdsService) ignoreTag(t Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *RdsService) tagsToString(tags []Tag) string {
	v, _ := json.Marshal(s.tagsToMap(tags))

	return string(v)
}

func (s *RdsService) DescribeDBProxy(id string) (map[string]interface{}, error) {
	action := "DescribeDBProxy"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return response, nil
}

func (s *RdsService) DescribeDBProxyEndpoint(id string, endpointName string) (map[string]interface{}, error) {
	action := "DescribeDBProxyEndpoint"
	request := map[string]interface{}{
		"RegionId":          s.client.RegionId,
		"DBInstanceId":      id,
		"DBProxyEndpointId": endpointName,
		"SourceIp":          s.client.SourceIp,
	}
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound", "Endpoint.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return response, nil
}

func (s *RdsService) DescribeRdsParameterGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewRdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeParameterGroup"
	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"ParameterGroupId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"ParamGroupsNotExistError"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("RdsParameterGroup", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.ParamGroup.ParameterGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ParamGroup.ParameterGroup", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("RDS", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["ParameterGroupId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("RDS", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *RdsService) DescribeRdsAccount(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewRdsClient()
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
		"RegionId":     s.client.RegionId,
		"SourceIp":     s.client.SourceIp,
		"AccountName":  parts[1],
		"DBInstanceId": parts[0],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("RdsAccount", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Accounts.DBInstanceAccount", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Accounts.DBInstanceAccount", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("RDS", id)), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *RdsService) RdsAccountStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeRdsAccount(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["AccountStatus"].(string) == failState {
				return object, object["AccountStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["AccountStatus"].(string)))
			}
		}
		return object, object["AccountStatus"].(string), nil
	}
}

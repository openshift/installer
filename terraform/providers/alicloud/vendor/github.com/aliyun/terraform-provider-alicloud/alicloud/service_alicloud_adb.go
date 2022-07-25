package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type AdbService struct {
	client *connectivity.AliyunClient
}

func (s *AdbService) DescribeAdbCluster(id string) (instance *adb.DBClusterInDescribeDBClusters, err error) {
	request := adb.CreateDescribeDBClustersRequest()
	request.RegionId = s.client.RegionId
	dbClusterIds := []string{}
	dbClusterIds = append(dbClusterIds, id)
	request.DBClusterIds = id
	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeDBClusters(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*adb.DescribeDBClustersResponse)
	if len(response.Items.DBCluster) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Cluster", id)), NotFoundMsg, ProviderERROR)
	}

	return &response.Items.DBCluster[0], nil
}

func (s *AdbService) DescribeAdbClusterAttribute(id string) (instance *adb.DBCluster, err error) {
	request := adb.CreateDescribeDBClusterAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeDBClusterAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*adb.DescribeDBClusterAttributeResponse)
	if len(response.Items.DBCluster) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Cluster", id)), NotFoundMsg, ProviderERROR)
	}

	return &response.Items.DBCluster[0], nil
}

func (s *AdbService) DescribeAdbAutoRenewAttribute(id string) (instance *adb.AutoRenewAttribute, err error) {
	request := adb.CreateDescribeAutoRenewAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterIds = id

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeAutoRenewAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*adb.DescribeAutoRenewAttributeResponse)
	if len(response.Items.AutoRenewAttribute) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Cluster", id)), NotFoundMsg, ProviderERROR)
	}

	return &response.Items.AutoRenewAttribute[0], nil
}

func (s *AdbService) WaitForAdbConnection(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAdbConnection(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted && object != nil && object.ConnectionString != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ConnectionString, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AdbService) DescribeAdbConnection(id string) (*adb.Address, error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(DefaultIntervalLong) * time.Second)
	for {
		object, err := s.DescribeAdbClusterNetInfo(parts[0])

		if err != nil {
			if NotFoundError(err) {
				return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return nil, WrapError(err)
		}

		if object != nil {
			for _, p := range object {
				if p.NetType == "Public" {
					return &p, nil
				}
			}
		}
		time.Sleep(DefaultIntervalMini * time.Second)
		if time.Now().After(deadline) {
			break
		}
	}

	return nil, WrapErrorf(Error(GetNotFoundMessage("DBConnection", id)), NotFoundMsg, ProviderERROR)
}

func (s *AdbService) DescribeAdbClusterNetInfo(id string) ([]adb.Address, error) {

	request := adb.CreateDescribeDBClusterNetInfoRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id
	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeDBClusterNetInfo(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*adb.DescribeDBClusterNetInfoResponse)
	if len(response.Items.Address) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBInstanceNetInfo", id)), NotFoundMsg, ProviderERROR)
	}

	return response.Items.Address, nil
}

func (s *AdbService) WaitForAdbAccount(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAdbAccount(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.AccountStatus == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AccountStatus, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *AdbService) DescribeAdbAccount(id string) (ds *adb.DBAccount, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := adb.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	invoker := NewInvoker()
	invoker.AddCatcher(DBInstanceStatusCatcher)
	var response *adb.DescribeAccountsResponse
	if err := invoker.Run(func() error {
		raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.DescribeAccounts(request)
		})
		if err != nil {
			return err
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		response, _ = raw.(*adb.DescribeAccountsResponse)
		return nil
	}); err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if len(response.AccountList.DBAccount) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBAccount", id)), NotFoundMsg, ProviderERROR)
	}
	return &response.AccountList.DBAccount[0], nil
}

// WaitForInstance waits for instance to given status
func (s *AdbService) WaitForAdbInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAdbCluster(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strings.ToLower(object.DBClusterStatus) == strings.ToLower(string(status)) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.DBClusterStatus, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *AdbService) setClusterTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

		if len(remove) > 0 {
			var tagKey []string
			for _, v := range remove {
				tagKey = append(tagKey, v.Key)
			}
			request := adb.CreateUntagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = "cluster"
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithAdbClient(func(client *adb.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := adb.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = "cluster"
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithAdbClient(func(client *adb.Client) (interface{}, error) {
				return client.TagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		d.SetPartial("tags")
	}

	return nil
}

func (s *AdbService) diffTags(oldTags, newTags []adb.TagResourcesTag) ([]adb.TagResourcesTag, []adb.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []adb.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *AdbService) tagsToMap(tags []adb.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *AdbService) tagsFromMap(m map[string]interface{}) []adb.TagResourcesTag {
	result := make([]adb.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, adb.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *AdbService) ignoreTag(t adb.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagValue)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *AdbService) DescribeTags(resourceId string, resourceType TagResourceType) (tags []adb.TagResource, err error) {
	request := adb.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	raw, err := s.client.WithAdbClient(func(client *adb.Client) (interface{}, error) {
		return client.ListTagResources(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*adb.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

// WaitForCluster waits for cluster to given status
func (s *AdbService) WaitForCluster(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAdbClusterAttribute(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strings.ToLower(object.DBClusterStatus) == strings.ToLower(string(status)) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.DBClusterStatus, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *AdbService) DescribeDBSecurityIps(clusterId string) (ips []string, err error) {

	request := adb.CreateDescribeDBClusterAccessWhiteListRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeDBClusterAccessWhiteList(request)
	})
	if err != nil {
		return ips, WrapErrorf(err, DefaultErrorMsg, clusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	resp, _ := raw.(*adb.DescribeDBClusterAccessWhiteListResponse)

	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, ip := range resp.Items.IPArray {
		if ip.DBClusterIPArrayAttribute != "hidden" {
			ipstr += separator + ip.SecurityIPList
			separator = COMMA_SEPARATED
		}
	}

	for _, ip := range strings.Split(ipstr, COMMA_SEPARATED) {
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

func (s *AdbService) ModifyDBSecurityIps(clusterId, ips string) error {

	request := adb.CreateModifyDBClusterAccessWhiteListRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId
	request.SecurityIps = ips

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.ModifyDBClusterAccessWhiteList(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, clusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForCluster(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *AdbService) DescribeAdbBackupPolicy(id string) (policy *adb.DescribeBackupPolicyResponse, err error) {

	request := adb.CreateDescribeBackupPolicyRequest()
	request.DBClusterId = id
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeBackupPolicy(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return policy, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return raw.(*adb.DescribeBackupPolicyResponse), nil
}

func (s *AdbService) ModifyAdbBackupPolicy(clusterId, backupTime, backupPeriod string) error {

	request := adb.CreateModifyBackupPolicyRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId
	request.PreferredBackupPeriod = backupPeriod
	request.PreferredBackupTime = backupTime

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.ModifyBackupPolicy(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, clusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForCluster(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *AdbService) AdbClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAdbClusterAttribute(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBClusterStatus == failState {
				return object, object.DBClusterStatus, WrapError(Error(FailedToReachTargetStatus, object.DBClusterStatus))
			}
		}
		return object, object.DBClusterStatus, nil
	}
}

func (s *AdbService) DescribeTask(id, taskId string) (*adb.DescribeTaskInfoResponse, error) {
	request := adb.CreateDescribeTaskInfoRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id
	request.TaskId = requests.Integer(taskId)

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeTaskInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*adb.DescribeTaskInfoResponse)

	return response, nil
}

func (s *AdbService) AdbTaskStateRefreshFunc(id, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeTask(id, taskId)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		return object, object.TaskInfo.Status, nil
	}
}

func (s *AdbService) DescribeAutoRenewAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewAdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeAutoRenewAttribute"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBClusterIds": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Items.AutoRenewAttribute", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.AutoRenewAttribute", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DBClusterId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *AdbService) DescribeDBClusterAccessWhiteList(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewAdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDBClusterAccessWhiteList"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("AnalyticdbForMysql3.0DbCluster", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Items.IPArray", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.IPArray", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), NotFoundWithResponse, response)
	} else {
		ipList := ""
		for _, item := range v.([]interface{}) {
			if item.(map[string]interface{})["DBClusterIPArrayAttribute"] == "hidden" {
				continue
			}
			ipList += item.(map[string]interface{})["SecurityIPList"].(string) + ","
		}
		v.([]interface{})[0].(map[string]interface{})["SecurityIPList"] = strings.TrimSuffix(ipList, ",")
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *AdbService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		conn, err := s.client.NewAdsClient()
		if err != nil {
			return WrapError(err)
		}

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if IsThrottling(err) {
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
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if IsThrottling(err) {
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

func (s *AdbService) DescribeAdbDbCluster(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewAdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDBClusterAttribute"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFoundError"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("AdbDbCluster", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Items.DBCluster", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.DBCluster", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DBClusterId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *AdbService) DescribeDBClusters(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewAdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDBClusters"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBClusterIds": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Items.DBCluster", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.DBCluster", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DBClusterId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *AdbService) AdbDbClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAdbDbCluster(id)
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

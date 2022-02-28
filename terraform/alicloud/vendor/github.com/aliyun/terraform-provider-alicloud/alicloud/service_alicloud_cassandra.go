package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cassandra"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type CassandraService struct {
	client *connectivity.AliyunClient
}

func (s *CassandraService) DescribeCassandraCluster(id string) (object cassandra.Cluster, err error) {
	cluster := cassandra.Cluster{}
	request := cassandra.CreateDescribeClusterRequest()
	request.RegionId = s.client.RegionId

	request.ClusterId = id

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeCluster(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Cluster.NotFound"}) {
			return cluster, WrapErrorf(Error(GetNotFoundMessage("Cassandra", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return cluster, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeClusterResponse)
	return response.Cluster, nil
}

func (s *CassandraService) CassandraClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCassandraCluster(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *CassandraService) DescribeCassandraDataCenter(id string) (object cassandra.DescribeDataCenterResponse, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := cassandra.CreateDescribeDataCenterRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = parts[1]
	request.DataCenterId = parts[0]

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeDataCenter(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Cluster.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("CassandraDataCenter", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeDataCenterResponse)
	return *response, nil
}

func (s *CassandraService) CassandraDataCenterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCassandraDataCenter(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *CassandraService) DescribeIpWhitelist(id string) (object cassandra.DescribeIpWhitelistResponse, err error) {
	request := cassandra.CreateDescribeIpWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeIpWhitelist(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeIpWhitelistResponse)
	return *response, nil
}

func (s *CassandraService) DescribeClusterDataCenter(id string) (object cassandra.DataCenter, err error) {
	request := cassandra.CreateDescribeDataCentersRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeDataCenters(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeDataCentersResponse)
	for _, dc := range response.DataCenters.DataCenter {
		if dc.CommodityInstance == id {
			return dc, nil
		}
	}
	return
}

func (s *CassandraService) DescribeSecurityGroups(id string) (object cassandra.DescribeSecurityGroupsResponse, err error) {
	request := cassandra.CreateDescribeSecurityGroupsRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeSecurityGroups(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeSecurityGroupsResponse)
	return *response, nil
}

func (s *CassandraService) DescribeCassandraEndpoints(id string) (object cassandra.DescribeContactPointsResponse, err error) {
	parts := strings.Split(id, ":")
	request := cassandra.CreateDescribeContactPointsRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = parts[len(parts)-1]

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeContactPoints(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeContactPointsResponse)
	return *response, nil
}

func (s *CassandraService) setInstanceTags(d *schema.ResourceData) error {
	if !d.HasChange("tags") {
		return nil
	}
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})

	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := cassandra.CreateUnTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.TagKey = &tagKey
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.UnTagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := cassandra.CreateTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.TagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.SetPartial("tags")
	return nil
}

func (s *CassandraService) diffTags(oldTags, newTags []cassandra.TagResourcesTag) ([]cassandra.TagResourcesTag, []cassandra.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []cassandra.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *CassandraService) tagsFromMap(m map[string]interface{}) []cassandra.TagResourcesTag {
	result := make([]cassandra.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, cassandra.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *CassandraService) tagsToMap(tags []cassandra.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *CassandraService) ignoreTag(t cassandra.Tag) bool {
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

func (s *CassandraService) DescribeAccounts(id string) (object cassandra.DescribeAccountsResponse, err error) {
	request := cassandra.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeAccounts(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeAccountsResponse)
	return *response, nil
}

func (s *CassandraService) DescribeCassandraBackupPlan(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCdsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeBackupPlans"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ClusterId": parts[0],
	}
	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Cluster.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Cassandra:BackupPlan", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.BackupPlans.BackupPlan", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.BackupPlans.BackupPlan", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Cassandra", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if v.(map[string]interface{})["DataCenterId"].(string) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("Cassandra", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

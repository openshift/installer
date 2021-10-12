package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type MongoDBService struct {
	client *connectivity.AliyunClient
}

func (s *MongoDBService) DescribeMongoDBInstance(id string) (instance dds.DBInstance, err error) {
	request := dds.CreateDescribeDBInstanceAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = id
	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.DescribeDBInstanceAttribute(request)
	})
	response, _ := raw.(*dds.DescribeDBInstanceAttributeResponse)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response == nil || len(response.DBInstances.DBInstance) == 0 {
		return instance, WrapErrorf(Error(GetNotFoundMessage("MongoDB Instance", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response.DBInstances.DBInstance[0], nil
}

func (s *MongoDBService) RdsMongodbDBInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMongoDBInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBInstanceStatus == failState {
				return object, object.DBInstanceStatus, WrapError(Error(FailedToReachTargetStatus, object.DBInstanceStatus))
			}
		}
		return object, object.DBInstanceStatus, nil
	}
}

func (s *MongoDBService) DescribeMongoDBSecurityIps(instanceId string) (ips []string, err error) {
	request := dds.CreateDescribeSecurityIpsRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = instanceId

	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.DescribeSecurityIps(request)
	})
	if err != nil {
		return ips, WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*dds.DescribeSecurityIpsResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, ip := range response.SecurityIpGroups.SecurityIpGroup {
		if ip.SecurityIpGroupAttribute == "hidden" {
			continue
		}
		ipstr += separator + ip.SecurityIpList
		separator = COMMA_SEPARATED
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

func (s *MongoDBService) ModifyMongoDBSecurityIps(d *schema.ResourceData, ips string) error {
	request := dds.CreateModifySecurityIpsRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = d.Id()
	request.SecurityIps = ips

	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.ModifySecurityIps(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, s.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *MongoDBService) DescribeMongoDBSecurityGroupId(id string) (*dds.DescribeSecurityGroupConfigurationResponse, error) {
	response := &dds.DescribeSecurityGroupConfigurationResponse{}
	request := dds.CreateDescribeSecurityGroupConfigurationRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = id

	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.DescribeSecurityGroupConfiguration(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*dds.DescribeSecurityGroupConfigurationResponse)

	return response, nil
}

func (server *MongoDBService) ModifyMongodbShardingInstanceNode(
	d *schema.ResourceData, nodeType MongoDBShardingNodeType, stateList, diffList []interface{}) error {
	client := server.client

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, server.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	//create node
	if len(stateList) < len(diffList) {
		createList := diffList[len(stateList):]
		diffList = diffList[:len(stateList)]

		for _, item := range createList {
			node := item.(map[string]interface{})

			request := dds.CreateCreateNodeRequest()
			request.RegionId = server.client.RegionId
			request.DBInstanceId = d.Id()
			request.NodeClass = node["node_class"].(string)
			request.NodeType = string(nodeType)
			if node["readonly_replicas"] != nil {
				request.ReadonlyReplicas = requests.NewInteger(node["readonly_replicas"].(int))
			}
			request.ClientToken = buildClientToken(request.GetActionName())

			if nodeType == MongoDBShardingNodeShard {
				request.NodeStorage = requests.NewInteger(node["node_storage"].(int))
			}

			raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
				return ddsClient.CreateNode(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)

			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, server.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapError(err)
			}
		}
	} else if len(stateList) > len(diffList) {
		deleteList := stateList[len(diffList):]
		stateList = stateList[:len(diffList)]

		for _, item := range deleteList {
			node := item.(map[string]interface{})

			request := dds.CreateDeleteNodeRequest()
			request.RegionId = server.client.RegionId
			request.DBInstanceId = d.Id()
			request.NodeId = node["node_id"].(string)
			request.ClientToken = buildClientToken(request.GetActionName())

			raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
				return ddsClient.DeleteNode(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}

			addDebug(request.GetActionName(), raw, request.RpcRequest, request)

			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, server.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapError(err)
			}
		}
	}

	//motify node
	for key := 0; key < len(stateList); key++ {
		state := stateList[key].(map[string]interface{})
		diff := diffList[key].(map[string]interface{})

		if state["node_class"] != diff["node_class"] ||
			state["node_storage"] != diff["node_storage"] {
			request := dds.CreateModifyNodeSpecRequest()
			request.RegionId = server.client.RegionId
			request.DBInstanceId = d.Id()
			request.NodeClass = diff["node_class"].(string)
			request.ClientToken = buildClientToken(request.GetActionName())

			if nodeType == MongoDBShardingNodeShard {
				request.NodeStorage = requests.NewInteger(diff["node_storage"].(int))
			}
			request.NodeId = state["node_id"].(string)

			raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
				return ddsClient.ModifyNodeSpec(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)

			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, server.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapError(err)
			}
		}
	}
	return nil
}

func (s *MongoDBService) DescribeMongoDBBackupPolicy(id string) (*dds.DescribeBackupPolicyResponse, error) {
	response := &dds.DescribeBackupPolicyResponse{}
	request := dds.CreateDescribeBackupPolicyRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = id
	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.DescribeBackupPolicy(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*dds.DescribeBackupPolicyResponse)
	return response, nil
}

func (s *MongoDBService) DescribeMongoDBTDEInfo(id string) (*dds.DescribeDBInstanceTDEInfoResponse, error) {

	response := &dds.DescribeDBInstanceTDEInfoResponse{}
	request := dds.CreateDescribeDBInstanceTDEInfoRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = id

	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.DescribeDBInstanceTDEInfo(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*dds.DescribeDBInstanceTDEInfoResponse)
	return response, nil
}

func (s *MongoDBService) DescribeDBInstanceSSL(id string) (*dds.DescribeDBInstanceSSLResponse, error) {
	response := &dds.DescribeDBInstanceSSLResponse{}
	request := dds.CreateDescribeDBInstanceSSLRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = id
	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.DescribeDBInstanceSSL(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*dds.DescribeDBInstanceSSLResponse)
	return response, nil
}

func (s *MongoDBService) MotifyMongoDBBackupPolicy(d *schema.ResourceData) error {
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0*time.Second, s.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}
	periodList := expandStringList(d.Get("backup_period").(*schema.Set).List())
	backupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
	backupTime := d.Get("backup_time").(string)

	request := dds.CreateModifyBackupPolicyRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = d.Id()
	request.PreferredBackupPeriod = backupPeriod
	request.PreferredBackupTime = backupTime
	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.ModifyBackupPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *MongoDBService) ResetAccountPassword(d *schema.ResourceData, password string) error {
	request := dds.CreateResetAccountPasswordRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = d.Id()
	request.AccountName = "root"
	request.AccountPassword = password
	raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.ResetAccountPassword(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return err
}

func (s *MongoDBService) setInstanceTags(d *schema.ResourceData) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})

	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := dds.CreateUntagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = "INSTANCE"
		request.TagKey = &tagKey
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.UntagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := dds.CreateTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.ResourceType = "INSTANCE"
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.TagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.SetPartial("tags")
	return nil
}

func (s *MongoDBService) tagsToMap(tags []dds.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *MongoDBService) ignoreTag(t dds.Tag) bool {
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

func (s *MongoDBService) tagsInAttributeToMap(tags []dds.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTagInAttribute(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *MongoDBService) ignoreTagInAttribute(t dds.Tag) bool {
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

func (s *MongoDBService) diffTags(oldTags, newTags []dds.TagResourcesTag) ([]dds.TagResourcesTag, []dds.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []dds.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *MongoDBService) tagsFromMap(m map[string]interface{}) []dds.TagResourcesTag {
	result := make([]dds.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, dds.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

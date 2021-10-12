package alicloud

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type KvstoreService struct {
	client *connectivity.AliyunClient
}

var KVstoreInstanceStatusCatcher = Catcher{"OperationDenied.KVstoreInstanceStatus", 60, 5}

func (s *KvstoreService) DescribeKVstoreInstance(id string) (*r_kvstore.DBInstanceAttribute, error) {
	instance := &r_kvstore.DBInstanceAttribute{}
	request := r_kvstore.CreateDescribeInstanceAttributeRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeInstanceAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return instance, WrapErrorf(Error(GetNotFoundMessage("KVstoreInstance", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeInstanceAttributeResponse)
	if len(response.Instances.DBInstanceAttribute) <= 0 {
		return instance, WrapErrorf(Error(GetNotFoundMessage("KVstoreInstance", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}

	return &response.Instances.DBInstanceAttribute[0], nil
}

func (s *KvstoreService) DescribeKVstoreBackupPolicy(id string) (*r_kvstore.DescribeBackupPolicyResponse, error) {
	response := &r_kvstore.DescribeBackupPolicyResponse{}
	request := r_kvstore.CreateDescribeBackupPolicyRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeBackupPolicy(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return response, WrapErrorf(Error(GetNotFoundMessage("KVstoreBackupPolicy", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*r_kvstore.DescribeBackupPolicyResponse)
	return response, nil
}

func (s *KvstoreService) WaitForKVstoreInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeKVstoreInstance(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.InstanceStatus == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceStatus, status, ProviderERROR)
		}
	}
	return nil
}

func (s *KvstoreService) RdsKvstoreInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeKVstoreInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.InstanceStatus == failState {
				return object, object.InstanceStatus, WrapError(Error(FailedToReachTargetStatus, object.InstanceStatus))
			}
		}
		return object, object.InstanceStatus, nil
	}
}

func (s *KvstoreService) WaitForKVstoreInstanceVpcAuthMode(id string, status string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeKVstoreInstance(id)
		if err != nil && !NotFoundError(err) {
			return err
		}
		if object.VpcAuthMode == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.VpcAuthMode, status, ProviderERROR)
		}
	}
	return nil
}

func (s *KvstoreService) DescribeParameters(id string) (*r_kvstore.DescribeParametersResponse, error) {
	response := &r_kvstore.DescribeParametersResponse{}
	request := r_kvstore.CreateDescribeParametersRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = id

	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeParameters(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return response, WrapErrorf(Error(GetNotFoundMessage("Parameters", id)), NotFoundMsg, ProviderERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*r_kvstore.DescribeParametersResponse)
	return response, nil
}

func (s *KvstoreService) ModifyInstanceConfig(id string, config string) error {
	request := r_kvstore.CreateModifyInstanceConfigRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id
	request.Config = config

	if err := s.WaitForKVstoreInstance(id, Normal, DefaultLongTimeout); err != nil {
		return WrapError(err)
	}
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.ModifyInstanceConfig(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *KvstoreService) setInstanceTags(d *schema.ResourceData) error {
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
			request := r_kvstore.CreateUntagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = strings.ToUpper(string(TagResourceInstance))
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := r_kvstore.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = strings.ToUpper(string(TagResourceInstance))
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
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

func (s *KvstoreService) tagsToMap(tags []r_kvstore.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *KvstoreService) tagsFromMap(m map[string]interface{}) []r_kvstore.TagResourcesTag {
	result := make([]r_kvstore.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, r_kvstore.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *KvstoreService) ignoreTag(t r_kvstore.TagResource) bool {
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

func (s *KvstoreService) diffTags(oldTags, newTags []r_kvstore.TagResourcesTag) ([]r_kvstore.TagResourcesTag, []r_kvstore.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []r_kvstore.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *KvstoreService) DescribeTags(resourceId string, resourceType TagResourceType) (tags []r_kvstore.TagResource, err error) {
	request := r_kvstore.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = strings.ToUpper(string(resourceType))
	request.ResourceId = &[]string{resourceId}
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.ListTagResources(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

func (s *KvstoreService) WaitForKVstoreAccount(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeKVstoreAccount(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object != nil && object.AccountStatus == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AccountStatus, status, ProviderERROR)
		}
	}
	return nil
}

func (s *KvstoreService) DescribeKVstoreAccount(id string) (*r_kvstore.Account, error) {
	ds := &r_kvstore.Account{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return ds, WrapError(err)
	}
	request := r_kvstore.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = parts[0]
	request.AccountName = parts[1]
	invoker := NewInvoker()
	invoker.AddCatcher(KVstoreInstanceStatusCatcher)
	var response *r_kvstore.DescribeAccountsResponse
	if err := invoker.Run(func() error {
		raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DescribeAccounts(request)
		})
		if err != nil {
			return err
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		response, _ = raw.(*r_kvstore.DescribeAccountsResponse)
		return nil
	}); err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return ds, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return ds, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if len(response.Accounts.Account) < 1 {
		return ds, WrapErrorf(Error(GetNotFoundMessage("KVstoreAccount", id)), NotFoundMsg, ProviderERROR)
	}
	return &response.Accounts.Account[0], nil
}

func (s *KvstoreService) DescribeKVstoreSecurityGroupId(id string) (*r_kvstore.DescribeSecurityGroupConfigurationResponse, error) {
	response := &r_kvstore.DescribeSecurityGroupConfigurationResponse{}
	request := r_kvstore.CreateDescribeSecurityGroupConfigurationRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id
	if err := s.WaitForKVstoreInstance(id, Normal, DefaultLongTimeout); err != nil {
		return response, WrapError(err)
	}
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeSecurityGroupConfiguration(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*r_kvstore.DescribeSecurityGroupConfigurationResponse)

	return response, nil
}

func (s *KvstoreService) DescribeDBInstanceNetInfo(id string) (*r_kvstore.NetInfoItemsInDescribeDBInstanceNetInfo, error) {
	response := &r_kvstore.DescribeDBInstanceNetInfoResponse{}
	request := r_kvstore.CreateDescribeDBInstanceNetInfoRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id
	if err := s.WaitForKVstoreInstance(id, Normal, DefaultLongTimeout); err != nil {
		return &response.NetInfoItems, WrapError(err)
	}
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeDBInstanceNetInfo(request)
	})
	if err != nil {
		return &response.NetInfoItems, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*r_kvstore.DescribeDBInstanceNetInfoResponse)

	return &response.NetInfoItems, nil
}

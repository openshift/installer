package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type SagService struct {
	client *connectivity.AliyunClient
}

func (s *SagService) DescribeCloudConnectNetwork(id string) (c smartag.CloudConnectNetwork, err error) {
	request := smartag.CreateDescribeCloudConnectNetworksRequest()
	request.RegionId = s.client.RegionId
	request.CcnId = id

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(ccnClient *smartag.Client) (interface{}, error) {
			return ccnClient.DescribeCloudConnectNetworks(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"CcnNotExist"}) {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeCloudConnectNetworksResponse)
	if len(response.CloudConnectNetworks.CloudConnectNetwork) <= 0 || response.CloudConnectNetworks.CloudConnectNetwork[0].CcnId != id {
		return c, WrapErrorf(Error(GetNotFoundMessage("CloudConnectNetwork ", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.CloudConnectNetworks.CloudConnectNetwork[0]
	return c, nil
}

func (s *SagService) DescribeCloudConnectNetworkGrant(id string) (c smartag.GrantRule, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeGrantRulesRequest()
	request.RegionId = s.client.RegionId
	request.AssociatedCcnId = parts[0]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(ccnClient *smartag.Client) (interface{}, error) {
			return ccnClient.DescribeGrantRules(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"CcnNotExist"}) {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeGrantRulesResponse)
	for _, value := range response.GrantRules.GrantRule {
		if value.CcnInstanceId == parts[0] && value.CenInstanceId == parts[1] {
			return value, nil
		}
	}
	return c, WrapErrorf(Error(GetNotFoundMessage("CloudConnectNetworkGrant", id)), NotFoundMsg, ProviderERROR)
}

func (s *SagService) DescribeCloudConnectNetworkAttachment(id string) (c smartag.SmartAccessGateway, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}
	ccnId := parts[0]
	sagId := parts[1]
	request := smartag.CreateDescribeSmartAccessGatewaysRequest()
	request.RegionId = s.client.RegionId
	request.SmartAGId = sagId

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(ccnClient *smartag.Client) (interface{}, error) {
			return ccnClient.DescribeSmartAccessGateways(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"SmartAccessGatewayNotExist"}) {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeSmartAccessGatewaysResponse)
	if len(response.SmartAccessGateways.SmartAccessGateway) <= 0 || response.SmartAccessGateways.SmartAccessGateway[0].AssociatedCcnId != ccnId || response.SmartAccessGateways.SmartAccessGateway[0].SmartAGId != sagId {
		return c, WrapErrorf(Error(GetNotFoundMessage("SmartAccessGatewayAttachment", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.SmartAccessGateways.SmartAccessGateway[0]
	return c, nil
}

func (s *SagService) DescribeSagAcl(id string) (c smartag.Acl, err error) {
	request := smartag.CreateDescribeACLsRequest()
	request.RegionId = s.client.RegionId
	request.AclIds = id

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeACLs(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"SagAclNotExist"}) {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeACLsResponse)
	if len(response.Acls.Acl) <= 0 || response.Acls.Acl[0].AclId != id {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag Acl", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.Acls.Acl[0]
	return c, nil
}

func (s *SagService) DescribeSagAclRule(id string) (c smartag.Acr, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeACLAttributeRequest()
	request.RegionId = s.client.RegionId
	request.AclId = parts[0]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeACLAttribute(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"SagAclRuleNotExist"}) {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeACLAttributeResponse)
	for _, value := range response.Acrs.Acr {
		if value.AcrId == parts[1] {
			return value, nil
		}
	}
	return c, WrapErrorf(Error(GetNotFoundMessage("Sag Acl Rule", id)), NotFoundMsg, ProviderERROR)
}

func (s *SagService) DescribeSagClientUser(id string) (c smartag.User, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeSmartAccessGatewayClientUsersRequest()
	request.RegionId = s.client.RegionId
	request.SmartAGId = parts[0]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeSmartAccessGatewayClientUsers(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"SagClientUserNotExist"}) {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeSmartAccessGatewayClientUsersResponse)
	for _, value := range response.Users.User {
		if value.UserName == parts[1] {
			return value, nil
		}
	}
	return c, WrapErrorf(Error(GetNotFoundMessage("Sag Client User", id)), NotFoundMsg, ProviderERROR)
}

func (s *SagService) DescribeSagSnatEntry(id string) (c smartag.SnatEntry, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeSnatEntriesRequest()
	request.RegionId = s.client.RegionId
	request.SmartAGId = parts[0]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeSnatEntries(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"SnatEntryiesNotExist"}) {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeSnatEntriesResponse)
	for _, value := range response.SnatEntries.SnatEntry {
		if value.InstanceId == parts[1] {
			return value, nil
		}
	}
	return c, WrapErrorf(Error(GetNotFoundMessage("sag_snat_entry", id)), NotFoundMsg, ProviderERROR)
}

func (s *SagService) DescribeSagDnatEntry(id string) (c smartag.DnatEntry, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeDnatEntriesRequest()
	request.RegionId = s.client.RegionId
	request.SagId = parts[0]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeDnatEntries(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DnatEntryiesNotExist"}) {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeDnatEntriesResponse)
	for _, value := range response.DnatEntries.DnatEntry {
		if value.DnatEntryId == parts[1] {
			return value, nil
		}
	}
	return c, WrapErrorf(Error(GetNotFoundMessage("sag_dnat_entry", id)), NotFoundMsg, ProviderERROR)
}

func (s *SagService) DescribeSagQos(id string) (c smartag.Qos, err error) {
	request := smartag.CreateDescribeQosesRequest()
	request.RegionId = s.client.RegionId
	request.QosIds = id

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeQoses(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"SagQosNotExist"}) {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeQosesResponse)
	if len(response.Qoses.Qos) <= 0 || response.Qoses.Qos[0].QosId != id {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag Qos", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.Qoses.Qos[0]
	return c, nil
}

func (s *SagService) DescribeSagQosPolicy(id string) (c smartag.QosPolicy, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeQosPoliciesRequest()
	request.RegionId = s.client.RegionId
	request.QosId = parts[0]
	request.QosPolicyId = parts[1]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeQosPolicies(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"SagQosPolicyNotExist"}) {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeQosPoliciesResponse)
	if len(response.QosPolicies.QosPolicy) <= 0 || response.QosPolicies.QosPolicy[0].QosPolicyId != parts[1] {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag Qos Policy", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.QosPolicies.QosPolicy[0]
	return c, nil
}

func (s *SagService) DescribeSagQosCar(id string) (c smartag.QosCar, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeQosCarsRequest()
	request.RegionId = s.client.RegionId
	request.QosId = parts[0]
	request.QosCarId = parts[1]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeQosCars(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"SagQosCarNotExist"}) {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeQosCarsResponse)
	if len(response.QosCars.QosCar) <= 0 || response.QosCars.QosCar[0].QosCarId != parts[1] {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag Qos Car", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.QosCars.QosCar[0]
	return c, nil
}

func (s *SagService) WaitForCloudConnectNetwork(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeCloudConnectNetwork(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.CcnId == id && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.CcnId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SagService) WaitForCloudConnectNetworkGrant(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	for {
		object, err := s.DescribeCloudConnectNetworkGrant(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.CenInstanceId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.CenInstanceId, parts[1], ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SagService) WaitForCloudConnectNetworkAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	ccnId := parts[0]
	for {
		object, err := s.DescribeCloudConnectNetworkAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.AssociatedCcnId == ccnId && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AssociatedCcnId, ccnId, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SagService) WaitForSagAcl(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSagAcl(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.AclId == id && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AclId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SagService) WaitForSagAclRule(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}

	for {
		object, err := s.DescribeSagAclRule(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.AcrId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AcrId, parts[1], ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SagService) WaitForSagQos(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeSagQos(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.QosId == id && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.QosId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SagService) WaitForSagQosPolicy(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeSagQosPolicy(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.QosPolicyId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.QosPolicyId, parts[1], ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SagService) WaitForSagQosCar(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSagQosCar(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.QosCarId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.QosCarId, parts[1], ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SagService) WaitForSagSnatEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	for {
		object, err := s.DescribeSagSnatEntry(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.InstanceId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceId, parts[1], ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SagService) WaitForSagDnatEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	for {
		object, err := s.DescribeSagDnatEntry(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.DnatEntryId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.DnatEntryId, parts[1], ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SagService) WaitForSagClientUser(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	for {
		object, err := s.DescribeSagClientUser(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.UserName == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.UserName, parts[1], ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

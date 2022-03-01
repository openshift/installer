package alicloud

import (
	"reflect"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type DnsService struct {
	client *connectivity.AliyunClient
}

func (s *DnsService) DescribeDns(id string) (*alidns.DescribeDomainInfoResponse, error) {
	response := &alidns.DescribeDomainInfoResponse{}
	request := alidns.CreateDescribeDomainInfoRequest()
	request.RegionId = s.client.RegionId
	request.DomainName = id

	raw, err := s.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.DescribeDomainInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ = raw.(*alidns.DescribeDomainInfoResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response.DomainName != id {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Dns", id)), NotFoundMsg, ProviderERROR)
	}
	return response, nil
}

func (dns *DnsService) DescribeDnsGroup(id string) (alidns.DomainGroup, error) {
	var group alidns.DomainGroup
	request := alidns.CreateDescribeDomainGroupsRequest()
	request.RegionId = dns.client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := dns.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		if err != nil {
			return group, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*alidns.DescribeDomainGroupsResponse)
		groups := response.DomainGroups.DomainGroup
		for _, domainGroup := range groups {
			if domainGroup.GroupId == id {
				return domainGroup, nil
			}
		}
		if len(groups) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return group, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return group, WrapErrorf(Error(GetNotFoundMessage("DnsGroup", id)), NotFoundMsg, ProviderERROR)
}

func (dns *DnsService) DescribeDnsRecord(id string) (*alidns.DescribeDomainRecordInfoResponse, error) {
	response := &alidns.DescribeDomainRecordInfoResponse{}
	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.RecordId = id
	request.RegionId = dns.client.RegionId
	raw, err := dns.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.DescribeDomainRecordInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser", "InvalidRR.NoExist"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*alidns.DescribeDomainRecordInfoResponse)
	if response.RecordId != id {
		return response, WrapErrorf(Error(GetNotFoundMessage("DnsRecord", id)), NotFoundMsg, ProviderERROR)
	}
	return response, nil
}

func (s *DnsService) DescribeDnsInstance(id string) (object alidns.DescribeDnsProductInstanceResponse, err error) {
	request := alidns.CreateDescribeDnsProductInstanceRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeDnsProductInstance(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDnsProduct"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DnsInstance", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeDnsProductInstanceResponse)
	return *response, nil
}

func (s *DnsService) DescribeAlidnsDomainAttachment(id string) (object alidns.DescribeInstanceDomainsResponse, err error) {
	request := alidns.CreateDescribeInstanceDomainsRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeInstanceDomains(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDnsProduct"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DnsDomainAttachment", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeInstanceDomainsResponse)

	if len(response.InstanceDomains) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("DnsDomainAttachment", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return *response, nil
}

func (s *DnsService) WaitForAlidnsDomainAttachment(id string, expected map[string]interface{}, isDelete bool, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlidnsDomainAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if isDelete {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		domainNames := make(map[string]interface{}, 0)
		for _, v := range object.InstanceDomains {
			domainNames[v.DomainName] = v.DomainName
		}

		exceptDomainNames := make(map[string]interface{}, 0)
		for _, v := range expected {
			for _, vv := range v.([]interface{}) {
				exceptDomainNames[vv.(string)] = vv.(string)
			}
		}

		if reflect.DeepEqual(domainNames, exceptDomainNames) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", expected, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *DnsService) ListTagResources(id string) (object alidns.ListTagResourcesResponse, err error) {
	request := alidns.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId

	request.ResourceType = "DOMAIN"
	request.ResourceId = &[]string{id}

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.ListTagResources(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.ListTagResourcesResponse)
	return *response, nil
}

func (s *DnsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]alidns.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, alidns.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := alidns.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := alidns.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *DnsService) DescribeDnsDomain(id string) (object alidns.DescribeDomainInfoResponse, err error) {
	request := alidns.CreateDescribeDomainInfoRequest()
	request.RegionId = s.client.RegionId

	request.DomainName = id

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeDomainInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DnsDomain", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeDomainInfoResponse)
	return *response, nil
}

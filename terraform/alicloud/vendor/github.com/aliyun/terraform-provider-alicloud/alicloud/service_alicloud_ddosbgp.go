package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddosbgp"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type DdosbgpService struct {
	client *connectivity.AliyunClient
}

func (s *DdosbgpService) DescribeDdosbgpInstance(id string) (v ddosbgp.Instance, err error) {
	request := ddosbgp.CreateDescribeInstanceListRequest()
	request.RegionId = s.client.RegionId
	request.DdosRegionId = s.client.RegionId
	request.InstanceIdList = "[\"" + id + "\"]"
	request.PageNo = "1"
	request.PageSize = "10"

	raw, err := s.client.WithDdosbgpClient(func(ddosbgpClient *ddosbgp.Client) (interface{}, error) {
		return ddosbgpClient.DescribeInstanceList(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return v, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}

		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ddosbgp.DescribeInstanceListResponse)
	if len(response.InstanceList) == 0 || response.InstanceList[0].InstanceId != id {
		return v, WrapErrorf(Error(GetNotFoundMessage("Ddosbgp Instance", id)), NotFoundMsg, ProviderERROR)
	}

	v = response.InstanceList[0]
	return
}

func (s *DdosbgpService) DescribeDdosbgpInstanceSpec(id string, region string) (v ddosbgp.InstanceSpec, err error) {
	request := ddosbgp.CreateDescribeInstanceSpecsRequest()
	request.InstanceIdList = "[\"" + id + "\"]"
	request.DdosRegionId = region
	request.RegionId = region

	raw, err := s.client.WithDdosbgpClient(func(ddosbgpClient *ddosbgp.Client) (interface{}, error) {
		return ddosbgpClient.DescribeInstanceSpecs(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound", "InvalidInstance"}) {
			return v, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}

		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ddosbgp.DescribeInstanceSpecsResponse)
	if len(resp.InstanceSpecs) == 0 || resp.InstanceSpecs[0].InstanceId != id {
		return v, WrapErrorf(Error(GetNotFoundMessage("DdosbgpInstanceSpec", id)), NotFoundMsg, ProviderERROR)
	}

	v = resp.InstanceSpecs[0]
	return v, WrapError(err)
}

func (s *DdosbgpService) WaitForDdosbgpInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDdosbgpInstance(id)

		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		} else if strings.ToLower(object.Status) == strings.ToLower(string(status)) {
			//TODO
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

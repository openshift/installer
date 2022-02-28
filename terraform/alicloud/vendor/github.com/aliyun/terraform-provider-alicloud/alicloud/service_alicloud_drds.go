package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DrdsService struct {
	client *connectivity.AliyunClient
}

func (s *DrdsService) DescribeDrdsInstance(id string) (*drds.DescribeDrdsInstanceResponse, error) {
	response := &drds.DescribeDrdsInstanceResponse{}
	request := drds.CreateDescribeDrdsInstanceRequest()
	request.RegionId = s.client.RegionId
	request.DrdsInstanceId = id
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstance(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDrdsInstanceId.NotFound"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*drds.DescribeDrdsInstanceResponse)
	if response.Data.Status == "5" {
		return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *DrdsService) DrdsInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDrdsInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Data.Status == failState {
				return object, object.Data.Status, WrapError(Error(FailedToReachTargetStatus, object.Data.Status))
			}
		}

		return object, object.Data.Status, nil
	}
}

func (s *DrdsService) WaitDrdsInstanceConfigEffect(id string, item map[string]string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		effected := false
		object, err := s.DescribeDrdsInstance(id)

		if err != nil {
			if NotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapError(err)
		}

		if value, ok := item["description"]; ok {
			if object.Data.Description == value {
				effected = true
			}
		}

		if effected {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Data, item, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

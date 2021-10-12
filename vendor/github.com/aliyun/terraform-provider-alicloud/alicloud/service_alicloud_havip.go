package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type HaVipService struct {
	client *connectivity.AliyunClient
}

func (s *HaVipService) DescribeHaVip(haVipId string) (v vpc.HaVip, err error) {
	request := vpc.CreateDescribeHaVipsRequest()
	request.RegionId = s.client.RegionId
	values := []string{haVipId}
	filter := []vpc.DescribeHaVipsFilter{{
		Key:   "HaVipId",
		Value: &values,
	},
	}
	request.Filter = &filter

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeHaVips(request)
		})
		if err != nil {
			return err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		resp, _ := raw.(*vpc.DescribeHaVipsResponse)
		if resp == nil || len(resp.HaVips.HaVip) <= 0 ||
			resp.HaVips.HaVip[0].HaVipId != haVipId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("HaVip", haVipId))
		}
		v = resp.HaVips.HaVip[0]
		return nil
	})
	return
}

func (s *HaVipService) WaitForHaVip(haVipId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		//wait the order effective
		havip, err := s.DescribeHaVip(haVipId)
		if err != nil {
			return err
		}
		if strings.ToLower(havip.Status) == strings.ToLower(string(status)) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("HaVip", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *HaVipService) DescribeHaVipAttachment(haVipId string, instanceId string) (err error) {
	invoker := NewInvoker()
	return invoker.Run(func() error {
		haVip, err := s.DescribeHaVip(haVipId)
		if err != nil {
			return err
		}
		for _, id := range haVip.AssociatedInstances.AssociatedInstance {
			if id == instanceId {
				return nil
			}
		}
		return GetNotFoundErrorFromString(GetNotFoundMessage("HaVipAttachment", haVipId+COLON_SEPARATED+instanceId))
	})
}

func (s *HaVipService) WaitForHaVipAttachment(haVipId string, instanceId string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		err := s.DescribeHaVipAttachment(haVipId, instanceId)

		if err != nil {
			if !NotFoundError(err) {
				return err
			}
		} else {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("HaVip Attachment", string("Unavailable")))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func getHaVipIdAndInstanceId(d *schema.ResourceData, meta interface{}) (string, string, error) {
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}

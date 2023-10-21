package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/market"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type MarketService struct {
	client *connectivity.AliyunClient
}

func (s *MarketService) DescribeMarketOrder(id string) (order *market.DescribeOrderResponse, err error) {
	request := market.CreateDescribeOrderRequest()
	request.OrderId = id
	raw, err := s.client.WithMarketClient(func(client *market.Client) (interface{}, error) {
		return client.DescribeOrder(request)
	})
	response, _ := raw.(*market.DescribeOrderResponse)
	if err != nil {
		if IsExpectedErrors(err, []string{"null"}) {
			return order, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return order, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response.OrderId == 0 {
		return order, WrapErrorf(Error(GetNotFoundMessage("Market Order", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func getRegionByImageIdPrefix(imageId string) string {
	switch imageId[:5] {

	case "m-m5e":
		return "cn-qingdao"
	case "m-2ze":
		return "cn-beijing"
	case "m-8vb":
		return "cn-zhangjiakou"
	case "m-hp3":
		return "cn-huhehaote"
	case "m-bp1":
		return "cn-hangzhou"
	case "m-uf6":
		return "cn-shanghai"
	case "m-wz9":
		return "cn-shenzhen"
	case "m-f8z":
		return "cn-heyuan"
	case "m-2vc":
		return "cn-chengdu"
	case "m-j6c":
		return "cn-hongkong"
	case "m-6we":
		return "ap-northeast-1"
	case "m-t4n":
		return "ap-southeast-1"
	case "m-p0w":
		return "ap-southeast-2"
	case "m-8ps":
		return "ap-southeast-3"
	case "m-k1a":
		return "ap-southeast-5"
	case "m-a2d":
		return "ap-south-1"
	case "m-0xi":
		return "us-east-1"
	case "m-rj9":
		return "us-west-1"
	case "m-d7o":
		return "eu-west-1"
	case "m-eb3":
		return "me-east-1"
	case "m-gw8":
		return "eu-central-1"
	default:
		return ""
	}
}

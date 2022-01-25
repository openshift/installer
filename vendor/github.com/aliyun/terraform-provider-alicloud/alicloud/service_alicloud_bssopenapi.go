package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const ModulesSizeLimit = 50

type BssopenapiService struct {
	client *connectivity.AliyunClient
}

func (b *BssopenapiService) GetInstanceTypePrice(productCode, productType string, modules interface{}) ([]float64, error) {
	var detailList []bssopenapi.ModuleDetail
	switch modules := modules.(type) {
	case []bssopenapi.GetPayAsYouGoPriceModuleList:
		request := bssopenapi.CreateGetPayAsYouGoPriceRequest()
		request.ProductCode = productCode
		request.ProductType = productType
		for {
			if len(modules) < ModulesSizeLimit {
				tmp := modules
				request.ModuleList = &tmp
			} else {
				tmp := modules[:ModulesSizeLimit]
				modules = modules[ModulesSizeLimit:]
				request.ModuleList = &tmp
			}
			data, err := b.getPayAsYouGoData(request)

			if err != nil {
				return nil, WrapError(err)
			}

			detailList = append(detailList, data.ModuleDetails.ModuleDetail...)

			if len(*request.ModuleList) < ModulesSizeLimit {
				break
			}
		}

	case []bssopenapi.GetSubscriptionPriceModuleList:
		request := bssopenapi.CreateGetSubscriptionPriceRequest()
		request.ProductCode = productCode
		request.ProductType = productType
		request.ServicePeriodQuantity = requests.NewInteger(1)
		request.ServicePeriodUnit = "Month"
		request.Quantity = requests.NewInteger(1)

		for {
			if len(modules) < ModulesSizeLimit {
				tmp := modules
				request.ModuleList = &tmp
			} else {
				tmp := modules[:ModulesSizeLimit]
				modules = modules[ModulesSizeLimit:]
				request.ModuleList = &tmp
			}
			data, err := b.getSubscriptionData(request)
			if err != nil {
				return nil, WrapError(err)
			}

			detailList = append(detailList, data.ModuleDetails.ModuleDetail...)

			if len(*request.ModuleList) < ModulesSizeLimit {
				break
			}
		}
	}

	var priceList []float64
	for _, module := range detailList {
		priceList = append(priceList, module.OriginalCost)
	}
	return priceList, nil
}

func (b *BssopenapiService) getSubscriptionData(request *bssopenapi.GetSubscriptionPriceRequest) (*bssopenapi.DataInGetSubscriptionPrice, error) {
	request.OrderType = "NewOrder"
	request.SubscriptionType = "Subscription"
	request.RegionId = string(connectivity.Hangzhou)
	var response *bssopenapi.GetSubscriptionPriceResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := b.client.WithBssopenapiClient(func(client *bssopenapi.Client) (interface{}, error) {
			return client.GetSubscriptionPrice(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				request.RegionId = string(connectivity.APSouthEast1)
				request.Domain = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*bssopenapi.GetSubscriptionPriceResponse)
		return nil
	})
	if err != nil {
		return nil, WrapError(err)
	}

	if !response.Success {
		return nil, WrapError(Error("Api:GetSubscriptionPrice  Modules:%v  RequestId:%s  Code:%s  Message:%s",
			request.ModuleList, response.RequestId, response.Code, response.Message))
	}

	if len(response.Data.ModuleDetails.ModuleDetail) == 0 {
		return nil, WrapError(Error("Api:GetSubscriptionPrice  Modules:%v  RequestId:%s  the moduleDetails length is 0!",
			request.ModuleList, response.RequestId))
	}
	return &response.Data, nil
}

func (b *BssopenapiService) getPayAsYouGoData(request *bssopenapi.GetPayAsYouGoPriceRequest) (*bssopenapi.DataInGetPayAsYouGoPrice, error) {
	request.SubscriptionType = "PayAsYouGo"
	request.RegionId = string(connectivity.Hangzhou)
	var response *bssopenapi.GetPayAsYouGoPriceResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := b.client.WithBssopenapiClient(func(client *bssopenapi.Client) (interface{}, error) {
			return client.GetPayAsYouGoPrice(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				request.RegionId = string(connectivity.APSouthEast1)
				request.Domain = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*bssopenapi.GetPayAsYouGoPriceResponse)
		return nil
	})

	if err != nil {
		return nil, WrapError(err)
	}

	if !response.Success {
		return nil, WrapError(Error("Api:GetPayAsYouGoPrice  Modules:%v  RequestId:%s  Code:%s  Message:%s",
			request.ModuleList, response.RequestId, response.Code, response.Message))
	}

	if len(response.Data.ModuleDetails.ModuleDetail) == 0 {
		return nil, WrapError(Error("Api:GetPayAsYouGoPrice  Modules:%v  RequestId:%s  the moduleDetails length is 0!",
			request.ModuleList, response.RequestId))
	}
	return &response.Data, nil
}

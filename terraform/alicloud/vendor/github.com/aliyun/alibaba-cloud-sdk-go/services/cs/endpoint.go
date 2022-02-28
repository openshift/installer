package cs

// EndpointMap Endpoint Data
var EndpointMap map[string]string

// EndpointType regional or central
var EndpointType = "regional"

// GetEndpointMap Get Endpoint Data Map
func GetEndpointMap() map[string]string {
	if EndpointMap == nil {
		EndpointMap = map[string]string{
			"cn-shanghai-internal-test-1": "cs.aliyuncs.com",
			"cn-beijing-gov-1":            "cs.aliyuncs.com",
			"cn-shenzhen-su18-b01":        "cs.aliyuncs.com",
			"cn-shanghai-inner":           "cs.aliyuncs.com",
			"cn-shenzhen-st4-d01":         "cs.aliyuncs.com",
			"cn-haidian-cm12-c01":         "cs.aliyuncs.com",
			"cn-hangzhou-internal-prod-1": "cs.aliyuncs.com",
			"cn-yushanfang":               "cs.aliyuncs.com",
			"cn-hongkong-finance-pop":     "cs.aliyuncs.com",
			"cn-qingdao-nebula":           "cs.aliyuncs.com",
			"cn-shanghai-finance-1":       "cs-vpc.cn-shanghai-finance-1.aliyuncs.com",
			"cn-beijing-finance-pop":      "cs.aliyuncs.com",
			"cn-wuhan":                    "cs.aliyuncs.com",
			"cn-zhangbei":                 "cs.aliyuncs.com",
			"cn-zhengzhou-nebula-1":       "cs.aliyuncs.com",
			"rus-west-1-pop":              "cs.aliyuncs.com",
			"cn-shanghai-et15-b01":        "cs.aliyuncs.com",
			"cn-hangzhou-bj-b01":          "cs.aliyuncs.com",
			"cn-hangzhou-internal-test-1": "cs.aliyuncs.com",
			"eu-west-1-oxs":               "cs.aliyuncs.com",
			"cn-zhangbei-na61-b01":        "cs.aliyuncs.com",
			"cn-beijing-finance-1":        "cs.aliyuncs.com",
			"cn-hangzhou-internal-test-3": "cs.aliyuncs.com",
			"cn-hangzhou-internal-test-2": "cs.aliyuncs.com",
			"cn-shenzhen-finance-1":       "cs-vpc.cn-shenzhen-finance-1.aliyuncs.com",
			"cn-hangzhou-test-306":        "cs.aliyuncs.com",
			"cn-huhehaote-nebula-1":       "cs.aliyuncs.com",
			"cn-shanghai-et2-b01":         "cs.aliyuncs.com",
			"cn-hangzhou-finance":         "cs-vpc.cn-hangzhou-finance.aliyuncs.com",
			"cn-beijing-nu16-b01":         "cs.aliyuncs.com",
			"cn-edge-1":                   "cs.aliyuncs.com",
			"cn-fujian":                   "cs.aliyuncs.com",
			"ap-northeast-2-pop":          "cs.aliyuncs.com",
			"cn-shenzhen-inner":           "cs.aliyuncs.com",
			"cn-zhangjiakou-na62-a01":     "cs.aliyuncs.com",
		}
	}
	return EndpointMap
}

// GetEndpointType Get Endpoint Type Value
func GetEndpointType() string {
	return EndpointType
}

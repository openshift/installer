package edas

// EndpointMap Endpoint Data
var EndpointMap map[string]string

// EndpointType regional or central
var EndpointType = "regional"

// GetEndpointMap Get Endpoint Data Map
func GetEndpointMap() map[string]string {
	if EndpointMap == nil {
		EndpointMap = map[string]string{
			"cn-shanghai-internal-test-1": "edas.aliyuncs.com",
			"cn-beijing-gov-1":            "edas.aliyuncs.com",
			"cn-shenzhen-su18-b01":        "edas.aliyuncs.com",
			"ap-south-1":                  "edas.ap-northeast-1.aliyuncs.com",
			"cn-shanghai-inner":           "edas.aliyuncs.com",
			"cn-shenzhen-st4-d01":         "edas.aliyuncs.com",
			"cn-haidian-cm12-c01":         "edas.aliyuncs.com",
			"cn-hangzhou-internal-prod-1": "edas.aliyuncs.com",
			"cn-yushanfang":               "edas.aliyuncs.com",
			"cn-hongkong-finance-pop":     "edas.aliyuncs.com",
			"cn-qingdao-nebula":           "edas.aliyuncs.com",
			"cn-shanghai-finance-1":       "edas.aliyuncs.com",
			"cn-beijing-finance-pop":      "edas.aliyuncs.com",
			"cn-wuhan":                    "edas.aliyuncs.com",
			"us-west-1":                   "edas.ap-northeast-1.aliyuncs.com",
			"cn-zhengzhou-nebula-1":       "edas.aliyuncs.com",
			"rus-west-1-pop":              "edas.ap-northeast-1.aliyuncs.com",
			"cn-shanghai-et15-b01":        "edas.aliyuncs.com",
			"cn-hangzhou-bj-b01":          "edas.aliyuncs.com",
			"eu-west-1":                   "edas.ap-northeast-1.aliyuncs.com",
			"cn-hangzhou-internal-test-1": "edas.aliyuncs.com",
			"eu-west-1-oxs":               "edas.ap-northeast-1.aliyuncs.com",
			"cn-zhangbei-na61-b01":        "edas.aliyuncs.com",
			"cn-beijing-finance-1":        "edas.aliyuncs.com",
			"cn-hangzhou-internal-test-3": "edas.aliyuncs.com",
			"cn-hangzhou-internal-test-2": "edas.aliyuncs.com",
			"cn-shenzhen-finance-1":       "edas.aliyuncs.com",
			"me-east-1":                   "edas.ap-northeast-1.aliyuncs.com",
			"cn-chengdu":                  "edas.aliyuncs.com",
			"cn-hangzhou-test-306":        "edas.aliyuncs.com",
			"cn-shanghai-et2-b01":         "edas.aliyuncs.com",
			"cn-hangzhou-finance":         "edas.aliyuncs.com",
			"cn-beijing-nu16-b01":         "edas.aliyuncs.com",
			"cn-edge-1":                   "edas.aliyuncs.com",
			"ap-southeast-3":              "edas.ap-northeast-1.aliyuncs.com",
			"cn-huhehaote":                "edas.aliyuncs.com",
			"ap-southeast-5":              "edas.ap-northeast-1.aliyuncs.com",
			"cn-fujian":                   "edas.aliyuncs.com",
			"ap-northeast-2-pop":          "edas.ap-northeast-1.aliyuncs.com",
			"cn-shenzhen-inner":           "edas.aliyuncs.com",
			"cn-zhangjiakou-na62-a01":     "edas.aliyuncs.com",
		}
	}
	return EndpointMap
}

// GetEndpointType Get Endpoint Type Value
func GetEndpointType() string {
	return EndpointType
}

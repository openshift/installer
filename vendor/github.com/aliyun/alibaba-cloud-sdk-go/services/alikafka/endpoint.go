package alikafka

// EndpointMap Endpoint Data
var EndpointMap map[string]string

// EndpointType regional or central
var EndpointType = "regional"

// GetEndpointMap Get Endpoint Data Map
func GetEndpointMap() map[string]string {
	if EndpointMap == nil {
		EndpointMap = map[string]string{
			"cn-shanghai-internal-test-1": "alikafka.aliyuncs.com",
			"cn-beijing-gov-1":            "alikafka.aliyuncs.com",
			"cn-shenzhen-su18-b01":        "alikafka.aliyuncs.com",
			"cn-wulanchabu":               "alikafka.aliyuncs.com",
			"cn-shanghai-inner":           "alikafka.aliyuncs.com",
			"cn-shenzhen-st4-d01":         "alikafka.aliyuncs.com",
			"cn-haidian-cm12-c01":         "alikafka.aliyuncs.com",
			"cn-hangzhou-internal-prod-1": "alikafka.aliyuncs.com",
			"cn-yushanfang":               "alikafka.aliyuncs.com",
			"cn-hongkong-finance-pop":     "alikafka.aliyuncs.com",
			"cn-qingdao-nebula":           "alikafka.aliyuncs.com",
			"cn-beijing-finance-pop":      "alikafka.aliyuncs.com",
			"cn-wuhan":                    "alikafka.aliyuncs.com",
			"cn-zhangbei":                 "alikafka.aliyuncs.com",
			"cn-zhengzhou-nebula-1":       "alikafka.aliyuncs.com",
			"rus-west-1-pop":              "alikafka.aliyuncs.com",
			"cn-shanghai-et15-b01":        "alikafka.aliyuncs.com",
			"cn-hangzhou-bj-b01":          "alikafka.aliyuncs.com",
			"cn-hangzhou-internal-test-1": "alikafka.aliyuncs.com",
			"eu-west-1-oxs":               "alikafka.aliyuncs.com",
			"cn-zhangbei-na61-b01":        "alikafka.aliyuncs.com",
			"cn-beijing-finance-1":        "alikafka.aliyuncs.com",
			"cn-hangzhou-internal-test-3": "alikafka.aliyuncs.com",
			"cn-hangzhou-internal-test-2": "alikafka.aliyuncs.com",
			"me-east-1":                   "alikafka.aliyuncs.com",
			"cn-hangzhou-test-306":        "alikafka.aliyuncs.com",
			"cn-huhehaote-nebula-1":       "alikafka.aliyuncs.com",
			"cn-shanghai-et2-b01":         "alikafka.aliyuncs.com",
			"cn-beijing-nu16-b01":         "alikafka.aliyuncs.com",
			"cn-edge-1":                   "alikafka.aliyuncs.com",
			"ap-southeast-2":              "alikafka.aliyuncs.com",
			"cn-fujian":                   "alikafka.aliyuncs.com",
			"ap-northeast-2-pop":          "alikafka.aliyuncs.com",
			"cn-shenzhen-inner":           "alikafka.aliyuncs.com",
			"cn-zhangjiakou-na62-a01":     "alikafka.aliyuncs.com",
		}
	}
	return EndpointMap
}

// GetEndpointType Get Endpoint Type Value
func GetEndpointType() string {
	return EndpointType
}

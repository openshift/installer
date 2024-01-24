package polardb

// EndpointMap Endpoint Data
var EndpointMap map[string]string

// EndpointType regional or central
var EndpointType = "regional"

// GetEndpointMap Get Endpoint Data Map
func GetEndpointMap() map[string]string {
	if EndpointMap == nil {
		EndpointMap = map[string]string{
			"cn-shanghai-internal-test-1": "polardb.aliyuncs.com",
			"cn-beijing-gov-1":            "polardb.aliyuncs.com",
			"cn-shenzhen-su18-b01":        "polardb.aliyuncs.com",
			"cn-beijing":                  "polardb.aliyuncs.com",
			"cn-wulanchabu":               "polardb.aliyuncs.com",
			"cn-shanghai-inner":           "polardb.aliyuncs.com",
			"cn-shenzhen-st4-d01":         "polardb.aliyuncs.com",
			"cn-haidian-cm12-c01":         "polardb.aliyuncs.com",
			"cn-hangzhou-internal-prod-1": "polardb.aliyuncs.com",
			"cn-north-2-gov-1":            "polardb.aliyuncs.com",
			"cn-yushanfang":               "polardb.aliyuncs.com",
			"cn-qingdao":                  "polardb.aliyuncs.com",
			"cn-hongkong-finance-pop":     "polardb.aliyuncs.com",
			"cn-qingdao-nebula":           "polardb.aliyuncs.com",
			"cn-shanghai":                 "polardb.aliyuncs.com",
			"cn-shanghai-finance-1":       "polardb.aliyuncs.com",
			"cn-hongkong":                 "polardb.aliyuncs.com",
			"cn-beijing-finance-pop":      "polardb.aliyuncs.com",
			"cn-wuhan":                    "polardb.aliyuncs.com",
			"us-west-1":                   "polardb.aliyuncs.com",
			"cn-zhangbei":                 "polardb.aliyuncs.com",
			"cn-shenzhen":                 "polardb.aliyuncs.com",
			"cn-zhengzhou-nebula-1":       "polardb.aliyuncs.com",
			"rus-west-1-pop":              "polardb.aliyuncs.com",
			"cn-shanghai-et15-b01":        "polardb.aliyuncs.com",
			"cn-hangzhou-bj-b01":          "polardb.aliyuncs.com",
			"cn-hangzhou-internal-test-1": "polardb.aliyuncs.com",
			"eu-west-1-oxs":               "polardb.aliyuncs.com",
			"cn-zhangbei-na61-b01":        "polardb.aliyuncs.com",
			"cn-beijing-finance-1":        "polardb.aliyuncs.com",
			"cn-hangzhou-internal-test-3": "polardb.aliyuncs.com",
			"cn-shenzhen-finance-1":       "polardb.aliyuncs.com",
			"cn-hangzhou-internal-test-2": "polardb.aliyuncs.com",
			"cn-hangzhou-test-306":        "polardb.aliyuncs.com",
			"cn-huhehaote-nebula-1":       "polardb.aliyuncs.com",
			"cn-shanghai-et2-b01":         "polardb.aliyuncs.com",
			"cn-hangzhou-finance":         "polardb.aliyuncs.com",
			"ap-southeast-1":              "polardb.aliyuncs.com",
			"cn-beijing-nu16-b01":         "polardb.aliyuncs.com",
			"cn-edge-1":                   "polardb.aliyuncs.com",
			"us-east-1":                   "polardb.aliyuncs.com",
			"cn-fujian":                   "polardb.aliyuncs.com",
			"ap-northeast-2-pop":          "polardb.aliyuncs.com",
			"cn-shenzhen-inner":           "polardb.aliyuncs.com",
			"cn-zhangjiakou-na62-a01":     "polardb.aliyuncs.com",
			"cn-hangzhou":                 "polardb.aliyuncs.com",
		}
	}
	return EndpointMap
}

// GetEndpointType Get Endpoint Type Value
func GetEndpointType() string {
	return EndpointType
}

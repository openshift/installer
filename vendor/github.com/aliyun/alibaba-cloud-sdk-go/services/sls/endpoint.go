package sls

// EndpointMap Endpoint Data
var EndpointMap map[string]string

// EndpointType regional or central
var EndpointType = "central"

// GetEndpointMap Get Endpoint Data Map
func GetEndpointMap() map[string]string {
	if EndpointMap == nil {
		EndpointMap = map[string]string{
			"cn-shenzhen":    "sls.cn-shenzhen.aliyuncs.com",
			"cn-shanghai":    "sls.cn-shanghai.aliyuncs.com",
			"cn-hongkong":    "sls.cn-hongkong.aliyuncs.com",
			"ap-southeast-1": "sls.ap-southeast-1.aliyuncs.com",
			"eu-central-1":   "sls.eu-central-1.aliyuncs.com",
			"cn-huhehaote":   "sls.cn-huhehaote.aliyuncs.com",
			"cn-zhangjiakou": "sls.cn-zhangjiakou.aliyuncs.com",
			"cn-hangzhou":    "sls.cn-hangzhou.aliyuncs.com",
		}
	}
	return EndpointMap
}

// GetEndpointType Get Endpoint Type Value
func GetEndpointType() string {
	return EndpointType
}

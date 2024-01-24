package alicloud

type ApiGatewayRequestConfig struct {
	Protocol   string `json:"requestProtocol"`
	Method     string `json:"requestHttpMethod"`
	Path       string `json:"requestPath"`
	Mode       string `json:"requestMode"`
	BodyFormat string `json:"bodyFormat"`
}

type ApiGatewayFunctionComputeConfig struct {
	Region       string `json:"fcRegionId"`
	ServiceName  string `json:"serviceName"`
	FunctionName string `json:"functionName"`
	Arn          string `json:"roleArn"`
}

type ApiGatewayVpcConfig struct {
	Name string `json:"name"`
}

type ApiGatewayServiceConfig struct {
	Protocol            string                          `json:"serviceProtocol"`
	Address             string                          `json:"serviceAddress"`
	Method              string                          `json:"serviceHttpMethod"`
	Path                string                          `json:"servicePath"`
	Timeout             int                             `json:"serviceTimeout"`
	ContentTypeCategory string                          `json:"contentTypeCatagory"`
	ContentTypeValue    string                          `json:"contentTypeValue"`
	MockEnable          string                          `json:"mock"`
	MockResult          string                          `json:"mockResult"`
	VpcEnable           string                          `json:"serviceVpcEnable"`
	FcConfig            ApiGatewayFunctionComputeConfig `json:"functionComputeConfig"`
	VpcConfig           ApiGatewayVpcConfig             `json:"vpcConfig"`
	AoneName            string                          `json:"aoneAppName"`
}

type ApiGatewayRequestParam struct {
	Type             string `json:"parameterType"`
	Name             string `json:"name"`
	ApiParameterName string `json:"apiParameterName"`
	Description      string `json:"description"`
	In               string `json:"location"`
	Required         string `json:"required"`
	DefualtValue     string `json:"defaultValue"`
}

type ApiGatewayServiceParam struct {
	Name    string `json:"serviceParameterName"`
	In      string `json:"location"`
	Type    string `json:"parameterType"`
	Catalog string `json:"parameterCatalog"`
}

type ApiGatewayParameterMap struct {
	RequestParamName string `json:"requestParameterName"`
	ServiceParamName string `json:"serviceParameterName"`
}

const (
	CatalogRequest           = "REQUEST"
	CatalogConstant          = "CONSTANT"
	CatalogSystem            = "SYSTEM"
	ResultType               = "JSON"
	ResultSample             = "Result Sample"
	Visibility               = "PRIVATE"
	AllowSignatureMethod     = "HmacSHA256"
	WebSocketApiType         = "COMMON"
	DeployCommonDescription  = "Terraform Deploy"
	StageNamePre             = "PRE"
	StageNameRelease         = "RELEASE"
	StageNameTest            = "TEST"
	AuthorizationDone        = "DONE"
	ApigatewayDefaultAddress = "http://www.aliyun.com"
	ApigatewayDefaultTimeout = 30
)

var ApiGatewayStageNames = []string{StageNamePre, StageNameRelease, StageNameTest}

package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliyunApigatewayApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewayApiCreate,
		Read:   resourceAliyunApigatewayApiRead,
		Update: resourceAliyunApigatewayApiUpdate,
		Delete: resourceAliyunApigatewayApiDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Required: true,
			},

			"auth_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"request_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"method": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Required: true,
						},
						"body_format": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"service_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"MOCK", "HTTP-VPC", "FunctionCompute", "HTTP"}, false),
			},

			"http_service_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"method": {
							Type:     schema.TypeString,
							Required: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"aone_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"http_vpc_service_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"method": {
							Type:     schema.TypeString,
							Required: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"aone_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"fc_service_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"function_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"arn_role": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"mock_service_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"result": {
							Type:     schema.TypeString,
							Required: true,
						},
						"aone_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"request_parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"required": {
							Type:     schema.TypeString,
							Required: true,
						},
						"in": {
							Type:     schema.TypeString,
							Required: true,
						},
						"in_service": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name_service": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"default_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"constant_parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"in": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"system_parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name_service": {
							Type:     schema.TypeString,
							Required: true,
						},
						"in": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"stage_names": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"PRE", "RELEASE", "TEST"}, false),
				},
				Optional: true,
			},

			"api_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"force_nonce_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliyunApigatewayApiCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request, err := buildAliyunApiArgs(d, meta)
	request.RegionId = client.RegionId
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.CreateApi(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apigateway_api", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cloudapi.CreateApiResponse)

	d.SetId(fmt.Sprintf("%s%s%s", request.GroupId, COLON_SEPARATED, response.ApiId))

	if l, ok := d.GetOk("stage_names"); ok {
		err = updateApiStages(d, l.(*schema.Set), meta)
		if err != nil {
			return WrapError(err)
		}
	}

	return resourceAliyunApigatewayApiRead(d, meta)
}

func resourceAliyunApigatewayApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	request := cloudapi.CreateDescribeApiRequest()
	request.RegionId = client.RegionId
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	object, err := cloudApiService.DescribeApiGatewayApi(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	stageNames, err := getStageNameList(d, cloudApiService)
	if err != nil {
		if !NotFoundError(err) {
			return WrapError(err)
		}
	}
	if err := d.Set("stage_names", stageNames); err != nil {
		return WrapError(err)
	}

	d.Set("api_id", object.ApiId)
	d.Set("group_id", object.GroupId)
	d.Set("name", object.ApiName)
	d.Set("description", object.Description)
	d.Set("auth_type", object.AuthType)
	d.Set("force_nonce_check", object.ForceNonceCheck)

	requestConfig := map[string]interface{}{}
	requestConfig["protocol"] = object.RequestConfig.RequestProtocol
	requestConfig["method"] = object.RequestConfig.RequestHttpMethod
	requestConfig["path"] = object.RequestConfig.RequestPath
	requestConfig["mode"] = object.RequestConfig.RequestMode
	if object.RequestConfig.BodyFormat != "" {
		requestConfig["body_format"] = object.RequestConfig.BodyFormat
	}
	if err := d.Set("request_config", []map[string]interface{}{requestConfig}); err != nil {
		return WrapError(err)
	}

	if object.ServiceConfig.Mock == "TRUE" {
		d.Set("service_type", "MOCK")
		MockServiceConfig := map[string]interface{}{}
		MockServiceConfig["result"] = object.ServiceConfig.MockResult
		MockServiceConfig["aone_name"] = object.ServiceConfig.AoneAppName
		if err := d.Set("mock_service_config", []map[string]interface{}{MockServiceConfig}); err != nil {
			return WrapError(err)
		}
	} else if object.ServiceConfig.ServiceVpcEnable == "TRUE" {
		d.Set("service_type", "HTTP-VPC")
		vpcServiceConfig := map[string]interface{}{}
		vpcServiceConfig["name"] = object.ServiceConfig.VpcConfig.Name
		vpcServiceConfig["path"] = object.ServiceConfig.ServicePath
		vpcServiceConfig["method"] = object.ServiceConfig.ServiceHttpMethod
		vpcServiceConfig["timeout"] = object.ServiceConfig.ServiceTimeout
		vpcServiceConfig["aone_name"] = object.ServiceConfig.AoneAppName
		if err := d.Set("http_vpc_service_config", []map[string]interface{}{vpcServiceConfig}); err != nil {
			return WrapError(err)
		}
	} else if object.ServiceConfig.ServiceProtocol == "FunctionCompute" {
		d.Set("service_type", "FunctionCompute")
		fcServiceConfig := map[string]interface{}{}
		fcServiceConfig["region"] = object.ServiceConfig.FunctionComputeConfig.RegionId
		fcServiceConfig["function_name"] = object.ServiceConfig.FunctionComputeConfig.FunctionName
		fcServiceConfig["service_name"] = object.ServiceConfig.FunctionComputeConfig.ServiceName
		fcServiceConfig["arn_role"] = object.ServiceConfig.FunctionComputeConfig.RoleArn
		fcServiceConfig["timeout"] = object.ServiceConfig.ServiceTimeout
		if err := d.Set("fc_service_config", []map[string]interface{}{fcServiceConfig}); err != nil {
			return WrapError(err)
		}
	} else {
		d.Set("service_type", "HTTP")
		httpServiceConfig := map[string]interface{}{}
		httpServiceConfig["address"] = object.ServiceConfig.ServiceAddress
		httpServiceConfig["path"] = object.ServiceConfig.ServicePath
		httpServiceConfig["method"] = object.ServiceConfig.ServiceHttpMethod
		httpServiceConfig["timeout"] = object.ServiceConfig.ServiceTimeout
		httpServiceConfig["aone_name"] = object.ServiceConfig.AoneAppName
		if err := d.Set("http_service_config", []map[string]interface{}{httpServiceConfig}); err != nil {
			return WrapError(err)
		}
	}

	requestParams := []map[string]interface{}{}
	for _, mapParam := range object.ServiceParametersMap.ServiceParameterMap {
		param := map[string]interface{}{}
		requestName := mapParam.RequestParameterName
		serviceName := mapParam.ServiceParameterName
		for _, serviceParam := range object.ServiceParameters.ServiceParameter {
			if serviceParam.ServiceParameterName == serviceName {
				param["name_service"] = serviceName
				param["in_service"] = strings.ToUpper(serviceParam.Location)
				break
			}
		}
		for _, requestParam := range object.RequestParameters.RequestParameter {
			if requestParam.ApiParameterName == requestName {
				param["name"] = requestName
				param["type"] = requestParam.ParameterType
				param["required"] = requestParam.Required
				param["in"] = requestParam.Location

				if requestParam.Description != "" {
					param["description"] = requestParam.Description
				}
				if requestParam.DefaultValue != "" {
					param["default_value"] = requestParam.DefaultValue
				}
				break
			}
		}
		requestParams = append(requestParams, param)
	}
	d.Set("request_parameters", requestParams)

	constantParams := []map[string]interface{}{}
	for _, constantParam := range object.ConstantParameters.ConstantParameter {
		param := map[string]interface{}{}
		param["name"] = constantParam.ServiceParameterName
		param["in"] = constantParam.Location
		param["value"] = constantParam.ConstantValue
		if constantParam.Description != "" {
			param["description"] = constantParam.Description
		}
		constantParams = append(constantParams, param)

	}
	d.Set("constant_parameters", constantParams)

	SystemParams := []map[string]interface{}{}
	for _, systemParam := range object.SystemParameters.SystemParameter {
		param := map[string]interface{}{}
		param["name"] = systemParam.ParameterName
		param["in"] = systemParam.Location
		param["name_service"] = systemParam.ServiceParameterName
	}
	d.Set("system_parameters", SystemParams)

	return nil
}

func resourceAliyunApigatewayApiUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cloudapi.CreateModifyApiRequest()
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	update := false

	d.Partial(true)

	if d.HasChange("name") || d.HasChange("description") || d.HasChange("auth_type") {
		update = true
	}
	request.ApiName = d.Get("name").(string)
	request.Description = d.Get("description").(string)
	request.AuthType = d.Get("auth_type").(string)

	if d.HasChange("force_nonce_check") {
		update = true
	}
	if v, exist := d.GetOk("force_nonce_check"); exist {
		request.ForceNonceCheck = requests.Boolean(strconv.FormatBool(v.(bool)))
	}

	var paramErr error
	var paramConfig string
	if d.HasChange("request_config") {
		update = true
	}
	paramConfig, paramErr = requestConfigToJsonStr(d.Get("request_config").([]interface{}))
	if paramErr != nil {
		return paramErr
	}
	request.RequestConfig = paramConfig

	if d.HasChange("service_type") || d.HasChange("http_service_config") || d.HasChange("http_vpc_service_config") || d.HasChange("mock_service_config") {
		update = true
	}
	serviceConfig, err := serviceConfigToJsonStr(d)
	if err != nil {
		return WrapError(err)
	}
	request.ServiceConfig = serviceConfig

	if d.HasChange("request_parameters") || d.HasChange("constant_parameters") || d.HasChange("system_parameters") {
		update = true
	}
	rps, sps, spm, err := setParameters(d)
	if err != nil {
		return WrapError(err)
	}
	request.RequestParameters = string(rps)
	request.ServiceParameters = string(sps)
	request.ServiceParametersMap = string(spm)

	if update {
		request.ResultType = ResultType
		request.ResultSample = ResultSample
		request.Visibility = Visibility
		request.AllowSignatureMethod = AllowSignatureMethod
		request.WebSocketApiType = WebSocketApiType

		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyApi(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		d.SetPartial("name")
		d.SetPartial("description")
		d.SetPartial("auth_type")
		d.SetPartial("service_type")
		d.SetPartial("http_service_config")
		d.SetPartial("http_vpc_service_config")
		d.SetPartial("fc_service_config")
		d.SetPartial("mock_service_config")
		d.SetPartial("request_parameters")
		d.SetPartial("constant_parameters")
		d.SetPartial("system_parameters")

	}

	if update || d.HasChange("stage_names") {
		if l, ok := d.GetOk("stage_names"); ok {
			err = updateApiStages(d, l.(*schema.Set), meta)
			if err != nil {
				return WrapError(err)
			}
		}
		d.SetPartial("stage_names")
	}

	d.Partial(false)

	return resourceAliyunApigatewayApiRead(d, meta)
}

func resourceAliyunApigatewayApiDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	request := cloudapi.CreateDeleteApiRequest()
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]

	for _, stageName := range ApiGatewayStageNames {
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			err := cloudApiService.AbolishApi(d.Id(), stageName)
			if err != nil {
				if IsExpectedErrors(err, []string{"ConcurrencyLockTimeout"}) {
					time.Sleep(3 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return WrapError(err)
		}
		_, err = cloudApiService.DescribeDeployedApi(d.Id(), stageName)
		if err != nil {
			if !NotFoundError(err) {
				return WrapError(err)
			}
		}
	}

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeleteApi(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"NotFoundApi"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(cloudApiService.WaitForApiGatewayApi(d.Id(), Deleted, DefaultTimeout))
}

func buildAliyunApiArgs(d *schema.ResourceData, meta interface{}) (*cloudapi.CreateApiRequest, error) {

	request := cloudapi.CreateCreateApiRequest()
	request.GroupId = d.Get("group_id").(string)
	request.Description = d.Get("description").(string)
	request.ApiName = d.Get("name").(string)
	request.AuthType = d.Get("auth_type").(string)
	if v, exist := d.GetOk("force_nonce_check"); exist {
		request.ForceNonceCheck = requests.Boolean(strconv.FormatBool(v.(bool)))
	}
	requestConfig, err := requestConfigToJsonStr(d.Get("request_config").([]interface{}))
	if err != nil {
		return request, WrapError(err)
	}
	request.RequestConfig = requestConfig

	serviceConfig, err := serviceConfigToJsonStr(d)
	if err != nil {
		return request, WrapError(err)
	}
	request.ServiceConfig = serviceConfig

	rps, sps, spm, err := setParameters(d)
	if err != nil {
		return request, WrapError(err)
	}

	request.RequestParameters = string(rps)
	request.ServiceParameters = string(sps)
	request.ServiceParametersMap = string(spm)

	request.ResultType = ResultType
	request.ResultSample = ResultSample
	request.Visibility = Visibility
	request.AllowSignatureMethod = AllowSignatureMethod
	request.WebSocketApiType = WebSocketApiType

	return request, WrapError(err)
}

func requestConfigToJsonStr(l []interface{}) (string, error) {
	config := l[0].(map[string]interface{})

	var requestConfig ApiGatewayRequestConfig
	requestConfig.Protocol = config["protocol"].(string)
	requestConfig.Path = config["path"].(string)
	requestConfig.Method = config["method"].(string)
	requestConfig.Mode = config["mode"].(string)

	if v, ok := config["body_format"]; ok {
		requestConfig.BodyFormat = v.(string)
	}
	configStr, err := json.Marshal(requestConfig)

	return string(configStr), WrapError(err)
}

func getHttpServiceConfig(d *schema.ResourceData) ([]byte, error) {
	var serviceConfig ApiGatewayServiceConfig
	var l []interface{}

	v, ok := d.GetOk("http_service_config")
	if !ok {
		return []byte{}, fmt.Errorf("Creating apigatway api error: http_service_config is null")
	}
	l = v.([]interface{})

	config := l[0].(map[string]interface{})
	serviceConfig.Protocol = "HTTP"
	serviceConfig.Address = config["address"].(string)
	serviceConfig.Path = config["path"].(string)
	serviceConfig.Method = config["method"].(string)
	serviceConfig.Timeout = config["timeout"].(int)
	serviceConfig.MockEnable = "FALSE"
	serviceConfig.VpcEnable = "FALSE"
	serviceConfig.ContentTypeCategory = "CLIENT"
	if v, ok := config["aone_name"]; ok {
		serviceConfig.AoneName = v.(string)
	}
	configStr, err := json.Marshal(serviceConfig)

	return configStr, WrapError(err)
}

func getHttpVpcServiceConfig(d *schema.ResourceData) ([]byte, error) {
	var serviceConfig ApiGatewayServiceConfig
	var l []interface{}

	v, ok := d.GetOk("http_vpc_service_config")
	if !ok {
		return []byte{}, WrapError(Error("Creating apigatway api error: http_vpc_service_config is null"))
	}
	l = v.([]interface{})

	config := l[0].(map[string]interface{})
	serviceConfig.Protocol = "HTTP"
	serviceConfig.VpcConfig.Name = config["name"].(string)
	serviceConfig.Path = config["path"].(string)
	serviceConfig.Method = config["method"].(string)
	serviceConfig.Timeout = config["timeout"].(int)
	serviceConfig.VpcEnable = "TRUE"
	serviceConfig.MockEnable = "FALSE"
	serviceConfig.ContentTypeCategory = "CLIENT"
	if v, ok := config["aone_name"]; ok {
		serviceConfig.AoneName = v.(string)
	}
	configStr, err := json.Marshal(serviceConfig)

	return configStr, WrapError(err)
}

func getFcServiceConfig(d *schema.ResourceData) ([]byte, error) {
	var serviceConfig ApiGatewayServiceConfig
	var l []interface{}

	v, ok := d.GetOk("fc_service_config")
	if !ok {
		return []byte{}, WrapError(Error("Creating apigatway api error: fc_service_config is null"))
	}
	l = v.([]interface{})

	config := l[0].(map[string]interface{})
	serviceConfig.Protocol = "FunctionCompute"
	serviceConfig.FcConfig.Region = config["region"].(string)
	serviceConfig.FcConfig.FunctionName = config["function_name"].(string)
	serviceConfig.FcConfig.ServiceName = config["service_name"].(string)
	serviceConfig.FcConfig.Arn = config["arn_role"].(string)
	serviceConfig.Timeout = config["timeout"].(int)
	serviceConfig.VpcEnable = "FALSE"
	serviceConfig.MockEnable = "FALSE"
	serviceConfig.ContentTypeCategory = "CLIENT"
	if v, ok := config["aone_name"]; ok {
		serviceConfig.AoneName = v.(string)
	}
	configStr, err := json.Marshal(serviceConfig)

	return configStr, WrapError(err)
}

func getMockServiceConfig(d *schema.ResourceData) ([]byte, error) {
	var serviceConfig ApiGatewayServiceConfig
	var l []interface{}

	v, ok := d.GetOk("mock_service_config")
	if !ok {
		return []byte{}, WrapError(Error("Creating apigatway api error: mock_service_config is null"))
	}
	l = v.([]interface{})

	config := l[0].(map[string]interface{})
	serviceConfig.Protocol = "HTTP"
	serviceConfig.Method = "GET"
	serviceConfig.MockResult = config["result"].(string)
	serviceConfig.MockEnable = "TRUE"
	serviceConfig.VpcEnable = "FALSE"
	serviceConfig.Timeout = ApigatewayDefaultTimeout
	serviceConfig.Address = ApigatewayDefaultAddress
	if v, ok := config["aone_name"]; ok {
		serviceConfig.AoneName = v.(string)
	}
	configStr, err := json.Marshal(serviceConfig)

	return configStr, WrapError(err)
}

func serviceConfigToJsonStr(d *schema.ResourceData) (string, error) {
	var err error
	var configStr []byte

	serviceType := d.Get("service_type").(string)

	switch serviceType {
	case "HTTP":
		configStr, err = getHttpServiceConfig(d)
		break
	case "HTTP-VPC":
		configStr, err = getHttpVpcServiceConfig(d)
		break
	case "FunctionCompute":
		configStr, err = getFcServiceConfig(d)
		break
	case "MOCK":
		configStr, err = getMockServiceConfig(d)
		break
	}
	if err != nil {
		return "", WrapError(err)
	}
	return string(configStr), nil
}

func setParameters(d *schema.ResourceData) (rps []byte, sps []byte, spm []byte, err error) {
	requestParameters := make([]ApiGatewayRequestParam, 0)
	serviceParameters := make([]ApiGatewayServiceParam, 0)
	serviceParamMaps := make([]ApiGatewayParameterMap, 0)

	requestParameters, serviceParameters, serviceParamMaps = setRequestParameters(d, requestParameters, serviceParameters, serviceParamMaps)
	requestParameters, serviceParameters, serviceParamMaps = setConstantParameters(d, requestParameters, serviceParameters, serviceParamMaps)
	requestParameters, serviceParameters, serviceParamMaps = setSystemParameters(d, requestParameters, serviceParameters, serviceParamMaps)

	rps, err = json.Marshal(requestParameters)
	if err != nil {
		err = WrapError(err)
		return
	}
	sps, err = json.Marshal(serviceParameters)
	if err != nil {
		err = WrapError(err)
		return
	}
	spm, err = json.Marshal(serviceParamMaps)
	if err != nil {
		err = WrapError(err)
		return
	}

	return rps, sps, spm, WrapError(err)
}

func setSystemParameters(d *schema.ResourceData, requestParameters []ApiGatewayRequestParam, serviceParameters []ApiGatewayServiceParam, serviceParamMaps []ApiGatewayParameterMap) ([]ApiGatewayRequestParam, []ApiGatewayServiceParam, []ApiGatewayParameterMap) {
	if l, ok := d.GetOk("system_parameters"); ok {

		for _, element := range l.(*schema.Set).List() {
			var requestParam ApiGatewayRequestParam
			var serviceParam ApiGatewayServiceParam
			var serviceParamMap ApiGatewayParameterMap

			request := element.(map[string]interface{})
			nameRequest := request["name"].(string)
			nameService := request["name_service"].(string)
			in := request["in"].(string)

			requestParam.Name = nameRequest
			requestParam.ApiParameterName = nameRequest
			requestParam.In = in
			requestParam.Required = "REQUIRED"
			requestParam.Type = "String"
			requestParameters = append(requestParameters, requestParam)

			serviceParam.Type = "String"
			serviceParam.In = in
			serviceParam.Name = nameService
			serviceParam.Catalog = CatalogSystem
			serviceParameters = append(serviceParameters, serviceParam)

			serviceParamMap.RequestParamName = nameRequest
			serviceParamMap.ServiceParamName = nameService
			serviceParamMaps = append(serviceParamMaps, serviceParamMap)
		}
	}

	return requestParameters, serviceParameters, serviceParamMaps
}

func setConstantParameters(d *schema.ResourceData, requestParameters []ApiGatewayRequestParam, serviceParameters []ApiGatewayServiceParam, serviceParamMaps []ApiGatewayParameterMap) ([]ApiGatewayRequestParam, []ApiGatewayServiceParam, []ApiGatewayParameterMap) {
	if l, ok := d.GetOk("constant_parameters"); ok {

		for _, element := range l.(*schema.Set).List() {
			var requestParam ApiGatewayRequestParam
			var serviceParam ApiGatewayServiceParam
			var serviceParamMap ApiGatewayParameterMap

			request := element.(map[string]interface{})
			name := request["name"].(string)
			in := request["in"].(string)
			value := request["value"].(string)

			requestParam.Name = name
			requestParam.Required = "REQUIRED"
			requestParam.ApiParameterName = name
			requestParam.In = in
			requestParam.Type = "String"
			if description, ok := request["description"]; !ok {
				requestParam.Description = description.(string)
			}
			requestParam.DefualtValue = value
			requestParameters = append(requestParameters, requestParam)

			serviceParam.In = in
			serviceParam.Name = name
			serviceParam.Catalog = CatalogConstant
			serviceParameters = append(serviceParameters, serviceParam)

			serviceParamMap.RequestParamName = name
			serviceParamMap.ServiceParamName = name
			serviceParamMaps = append(serviceParamMaps, serviceParamMap)
		}
	}

	return requestParameters, serviceParameters, serviceParamMaps
}

func setRequestParameters(d *schema.ResourceData, requestParameters []ApiGatewayRequestParam, serviceParameters []ApiGatewayServiceParam, serviceParamMaps []ApiGatewayParameterMap) ([]ApiGatewayRequestParam, []ApiGatewayServiceParam, []ApiGatewayParameterMap) {
	if l, ok := d.GetOk("request_parameters"); ok {

		for _, element := range l.(*schema.Set).List() {
			var requestParam ApiGatewayRequestParam
			var serviceParam ApiGatewayServiceParam
			var serviceParamMap ApiGatewayParameterMap

			request := element.(map[string]interface{})
			nameRequest := request["name"].(string)
			paramType := request["type"].(string)
			required := request["required"].(string)
			in := request["in"].(string)

			inService := request["in_service"].(string)
			nameService := request["name_service"].(string)

			if description, ok := request["description"]; ok {
				requestParam.Description = description.(string)
			}
			if defaultValue, ok := request["default_value"]; ok {
				requestParam.DefualtValue = defaultValue.(string)
			}

			requestParam.Name = nameRequest
			requestParam.Required = required
			requestParam.ApiParameterName = nameRequest
			requestParam.In = in
			requestParam.Type = paramType

			requestParameters = append(requestParameters, requestParam)

			serviceParam.Type = paramType
			serviceParam.In = inService
			serviceParam.Name = nameService
			serviceParam.Catalog = CatalogRequest
			serviceParameters = append(serviceParameters, serviceParam)

			serviceParamMap.RequestParamName = nameRequest
			serviceParamMap.ServiceParamName = nameService
			serviceParamMaps = append(serviceParamMaps, serviceParamMap)
		}
	}

	return requestParameters, serviceParameters, serviceParamMaps
}

func getStageNameList(d *schema.ResourceData, cloudApiService CloudApiService) ([]string, error) {
	var stageNames []string

	for _, stageName := range ApiGatewayStageNames {
		_, err := cloudApiService.DescribeDeployedApi(d.Id(), stageName)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return nil, WrapError(err)
		}
		stageNames = append(stageNames, stageName)
	}
	return stageNames, nil
}

func updateApiStages(d *schema.ResourceData, stageNames *schema.Set, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	for _, stageName := range ApiGatewayStageNames {
		if stageNames.Contains(stageName) {
			err := cloudApiService.DeployedApi(d.Id(), stageName)

			if err != nil {
				return WrapError(err)
			}

			_, err = cloudApiService.DescribeDeployedApi(d.Id(), stageName)
			if err != nil {
				return WrapError(err)
			}

		} else {
			err := cloudApiService.AbolishApi(d.Id(), stageName)
			if err != nil {
				return WrapError(err)
			}
			_, err = cloudApiService.DescribeDeployedApi(d.Id(), stageName)
			if err != nil {
				if !NotFoundError(err) {
					return WrapError(err)
				}
			}
		}
	}
	return nil
}

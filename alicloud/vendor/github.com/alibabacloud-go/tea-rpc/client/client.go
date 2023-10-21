// This file is auto-generated, don't edit it. Thanks.
/**
 * This is for RPC SDK
 */
package client

import (
	rpcutil "github.com/alibabacloud-go/tea-rpc-utils/service"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
)

/**
 * Model for initing client
 */
type Config struct {
	// accesskey id
	AccessKeyId *string `json:"accessKeyId,omitempty" xml:"accessKeyId,omitempty"`
	// accesskey secret
	AccessKeySecret *string `json:"accessKeySecret,omitempty" xml:"accessKeySecret,omitempty"`
	// security token
	SecurityToken *string `json:"securityToken,omitempty" xml:"securityToken,omitempty"`
	// http protocol
	Protocol *string `json:"protocol,omitempty" xml:"protocol,omitempty"`
	// region id
	RegionId *string `json:"regionId,omitempty" xml:"regionId,omitempty" pattern:"^[a-zA-Z0-9_-]+$"`
	// read timeout
	ReadTimeout *int `json:"readTimeout,omitempty" xml:"readTimeout,omitempty"`
	// connect timeout
	ConnectTimeout *int `json:"connectTimeout,omitempty" xml:"connectTimeout,omitempty"`
	// http proxy
	HttpProxy *string `json:"httpProxy,omitempty" xml:"httpProxy,omitempty"`
	// https proxy
	HttpsProxy *string `json:"httpsProxy,omitempty" xml:"httpsProxy,omitempty"`
	// credential
	Credential credential.Credential `json:"credential,omitempty" xml:"credential,omitempty"`
	// endpoint
	Endpoint *string `json:"endpoint,omitempty" xml:"endpoint,omitempty"`
	// proxy white list
	NoProxy *string `json:"noProxy,omitempty" xml:"noProxy,omitempty"`
	// max idle conns
	MaxIdleConns *int `json:"maxIdleConns,omitempty" xml:"maxIdleConns,omitempty"`
	// network for endpoint
	Network *string `json:"network,omitempty" xml:"network,omitempty" pattern:"^[a-zA-Z0-9_-]+$"`
	// user agent
	UserAgent *string `json:"userAgent,omitempty" xml:"userAgent,omitempty"`
	// suffix for endpoint
	Suffix *string `json:"suffix,omitempty" xml:"suffix,omitempty" pattern:"^[a-zA-Z0-9_-]+$"`
	// socks5 proxy
	Socks5Proxy *string `json:"socks5Proxy,omitempty" xml:"socks5Proxy,omitempty"`
	// socks5 network
	Socks5NetWork *string `json:"socks5NetWork,omitempty" xml:"socks5NetWork,omitempty"`
	// endpoint type
	EndpointType *string `json:"endpointType,omitempty" xml:"endpointType,omitempty"`
	// OpenPlatform endpoint
	OpenPlatformEndpoint *string `json:"openPlatformEndpoint,omitempty" xml:"openPlatformEndpoint,omitempty"`
	// Deprecated
	// credential type
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// source ip
	SourceIp *string `json:"sourceIp,omitempty" xml:"sourceIp,omitempty"`
	// secure transport
	SecureTransport *string `json:"secureTransport,omitempty" xml:"secureTransport,omitempty"`
}

func (s Config) String() string {
	return tea.Prettify(s)
}

func (s Config) GoString() string {
	return s.String()
}

func (s *Config) SetAccessKeyId(v string) *Config {
	s.AccessKeyId = &v
	return s
}

func (s *Config) SetAccessKeySecret(v string) *Config {
	s.AccessKeySecret = &v
	return s
}

func (s *Config) SetSecurityToken(v string) *Config {
	s.SecurityToken = &v
	return s
}

func (s *Config) SetProtocol(v string) *Config {
	s.Protocol = &v
	return s
}

func (s *Config) SetRegionId(v string) *Config {
	s.RegionId = &v
	return s
}

func (s *Config) SetReadTimeout(v int) *Config {
	s.ReadTimeout = &v
	return s
}

func (s *Config) SetConnectTimeout(v int) *Config {
	s.ConnectTimeout = &v
	return s
}

func (s *Config) SetHttpProxy(v string) *Config {
	s.HttpProxy = &v
	return s
}

func (s *Config) SetHttpsProxy(v string) *Config {
	s.HttpsProxy = &v
	return s
}

func (s *Config) SetCredential(v credential.Credential) *Config {
	s.Credential = v
	return s
}

func (s *Config) SetEndpoint(v string) *Config {
	s.Endpoint = &v
	return s
}

func (s *Config) SetNoProxy(v string) *Config {
	s.NoProxy = &v
	return s
}

func (s *Config) SetMaxIdleConns(v int) *Config {
	s.MaxIdleConns = &v
	return s
}

func (s *Config) SetNetwork(v string) *Config {
	s.Network = &v
	return s
}

func (s *Config) SetUserAgent(v string) *Config {
	s.UserAgent = &v
	return s
}

func (s *Config) SetSuffix(v string) *Config {
	s.Suffix = &v
	return s
}

func (s *Config) SetSocks5Proxy(v string) *Config {
	s.Socks5Proxy = &v
	return s
}

func (s *Config) SetSocks5NetWork(v string) *Config {
	s.Socks5NetWork = &v
	return s
}

func (s *Config) SetEndpointType(v string) *Config {
	s.EndpointType = &v
	return s
}

func (s *Config) SetOpenPlatformEndpoint(v string) *Config {
	s.OpenPlatformEndpoint = &v
	return s
}

func (s *Config) SetType(v string) *Config {
	s.Type = &v
	return s
}

func (s *Config) SetSourceIp(v string) *Config {
	s.SourceIp = &v
	return s
}

func (s *Config) SetSecureTransport(v string) *Config {
	s.SecureTransport = &v
	return s
}

type Client struct {
	Endpoint             *string
	RegionId             *string
	Protocol             *string
	UserAgent            *string
	EndpointRule         *string
	EndpointMap          map[string]*string
	Suffix               *string
	ReadTimeout          *int
	ConnectTimeout       *int
	HttpProxy            *string
	HttpsProxy           *string
	Socks5Proxy          *string
	Socks5NetWork        *string
	NoProxy              *string
	Network              *string
	ProductId            *string
	MaxIdleConns         *int
	EndpointType         *string
	OpenPlatformEndpoint *string
	SourceIp             *string
	SecureTransport      *string
	Credential           credential.Credential
}

/**
 * Init client with Config
 * @param config config contains the necessary information to create a client
 */
func NewClient(config *Config) (*Client, error) {
	client := new(Client)
	err := client.Init(config)
	return client, err
}

func (client *Client) Init(config *Config) (_err error) {
	if tea.BoolValue(util.IsUnset(tea.ToMap(config))) {
		_err = tea.NewSDKError(map[string]interface{}{
			"code":    "ParameterMissing",
			"message": "'config' can not be unset",
		})
		return _err
	}

	_err = util.ValidateModel(config)
	if _err != nil {
		return _err
	}
	if !tea.BoolValue(util.Empty(config.AccessKeyId)) && !tea.BoolValue(util.Empty(config.AccessKeySecret)) {
		if !tea.BoolValue(util.Empty(config.SecurityToken)) {
			config.Type = tea.String("sts")
		} else {
			config.Type = tea.String("access_key")
		}

		credentialConfig := &credential.Config{
			AccessKeyId:     config.AccessKeyId,
			Type:            config.Type,
			AccessKeySecret: config.AccessKeySecret,
			SecurityToken:   config.SecurityToken,
		}
		client.Credential, _err = credential.NewCredential(credentialConfig)
		if _err != nil {
			return _err
		}

	} else if !tea.BoolValue(util.IsUnset(config.Credential)) {
		client.Credential = config.Credential
	} else {
		_err = tea.NewSDKError(map[string]interface{}{
			"code":    "ParameterMissing",
			"message": "'accessKeyId' and 'accessKeySecret' or 'credential' can not be unset",
		})
		return _err
	}

	client.SourceIp = config.SourceIp
	client.SecureTransport = config.SecureTransport
	client.Network = config.Network
	client.Suffix = config.Suffix
	client.Endpoint = config.Endpoint
	client.Protocol = config.Protocol
	client.RegionId = config.RegionId
	client.UserAgent = config.UserAgent
	client.ReadTimeout = config.ReadTimeout
	client.ConnectTimeout = config.ConnectTimeout
	client.HttpProxy = config.HttpProxy
	client.HttpsProxy = config.HttpsProxy
	client.NoProxy = config.NoProxy
	client.Socks5Proxy = config.Socks5Proxy
	client.Socks5NetWork = config.Socks5NetWork
	client.MaxIdleConns = config.MaxIdleConns
	client.EndpointType = config.EndpointType
	client.OpenPlatformEndpoint = config.OpenPlatformEndpoint
	return nil
}

/**
 * Encapsulate the request and invoke the network
 * @param action api name
 * @param protocol http or https
 * @param method e.g. GET
 * @param version product version
 * @param authType when authType is Anonymous, the signature will not be calculate
 * @param pathname pathname of every api
 * @param query which contains request params
 * @param body content of request
 * @param runtime which controls some details of call api, such as retry times
 * @return the response
 */
func (client *Client) DoRequest(action *string, protocol *string, method *string, version *string, authType *string, query map[string]interface{}, body map[string]interface{}, runtime *util.RuntimeOptions) (_result map[string]interface{}, _err error) {
	_err = tea.Validate(runtime)
	if _err != nil {
		return _result, _err
	}
	_runtime := map[string]interface{}{
		"timeouted":      "retry",
		"readTimeout":    tea.IntValue(util.DefaultNumber(runtime.ReadTimeout, client.ReadTimeout)),
		"connectTimeout": tea.IntValue(util.DefaultNumber(runtime.ConnectTimeout, client.ConnectTimeout)),
		"httpProxy":      tea.StringValue(util.DefaultString(runtime.HttpProxy, client.HttpProxy)),
		"httpsProxy":     tea.StringValue(util.DefaultString(runtime.HttpsProxy, client.HttpsProxy)),
		"noProxy":        tea.StringValue(util.DefaultString(runtime.NoProxy, client.NoProxy)),
		"maxIdleConns":   tea.IntValue(util.DefaultNumber(runtime.MaxIdleConns, client.MaxIdleConns)),
		"retry": map[string]interface{}{
			"retryable":   tea.BoolValue(runtime.Autoretry),
			"maxAttempts": tea.IntValue(util.DefaultNumber(runtime.MaxAttempts, tea.Int(3))),
		},
		"backoff": map[string]interface{}{
			"policy": tea.StringValue(util.DefaultString(runtime.BackoffPolicy, tea.String("no"))),
			"period": tea.IntValue(util.DefaultNumber(runtime.BackoffPeriod, tea.Int(1))),
		},
		"ignoreSSL": tea.BoolValue(runtime.IgnoreSSL),
	}

	_resp := make(map[string]interface{})
	for _retryTimes := 0; tea.BoolValue(tea.AllowRetry(_runtime["retry"], tea.Int(_retryTimes))); _retryTimes++ {
		if _retryTimes > 0 {
			_backoffTime := tea.GetBackoffTime(_runtime["backoff"], tea.Int(_retryTimes))
			if tea.IntValue(_backoffTime) > 0 {
				tea.Sleep(_backoffTime)
			}
		}

		_resp, _err = func() (map[string]interface{}, error) {
			request_ := tea.NewRequest()
			request_.Protocol = util.DefaultString(client.Protocol, protocol)
			request_.Method = method
			request_.Pathname = tea.String("/")
			request_.Query = rpcutil.Query(tea.ToMap(map[string]interface{}{
				"Action":         tea.StringValue(action),
				"Format":         "json",
				"Timestamp":      tea.StringValue(rpcutil.GetTimestamp()),
				"Version":        tea.StringValue(version),
				"SignatureNonce": tea.StringValue(util.GetNonce()),
			}, query))
			if !tea.BoolValue(util.IsUnset(client.SourceIp)) {
				request_.Query["SourceIp"] = client.SourceIp
			}

			if !tea.BoolValue(util.IsUnset(client.SecureTransport)) {
				request_.Query["SecureTransport"] = client.SecureTransport
			}

			// endpoint is setted in product client
			request_.Headers = map[string]*string{
				"x-acs-version": version,
				"x-acs-action":  action,
				"host":          client.Endpoint,
				"user-agent":    client.GetUserAgent(),
			}
			if !tea.BoolValue(util.IsUnset(body)) {
				tmp := util.AnyifyMapValue(rpcutil.Query(body))
				request_.Body = tea.ToReader(util.ToFormString(tmp))
				request_.Headers["content-type"] = tea.String("application/x-www-form-urlencoded")
			}

			if !tea.BoolValue(util.EqualString(authType, tea.String("Anonymous"))) {
				accessKeyId, _err := client.GetAccessKeyId()
				if _err != nil {
					return _result, _err
				}

				accessKeySecret, _err := client.GetAccessKeySecret()
				if _err != nil {
					return _result, _err
				}

				securityToken, _err := client.GetSecurityToken()
				if _err != nil {
					return _result, _err
				}

				if !tea.BoolValue(util.Empty(securityToken)) {
					request_.Query["SecurityToken"] = securityToken
				}

				request_.Query["SignatureMethod"] = tea.String("HMAC-SHA1")
				request_.Query["SignatureVersion"] = tea.String("1.0")
				request_.Query["AccessKeyId"] = accessKeyId
				signedParam := tea.Merge(request_.Query,
					rpcutil.Query(body))
				request_.Query["Signature"] = rpcutil.GetSignatureV1(signedParam, request_.Method, accessKeySecret)
			}

			response_, _err := tea.DoRequest(request_, _runtime)
			if _err != nil {
				return _result, _err
			}
			obj, _err := util.ReadAsJSON(response_.Body)
			if _err != nil {
				return _result, _err
			}

			res := util.AssertAsMap(obj)
			if tea.BoolValue(util.Is4xx(response_.StatusCode)) || tea.BoolValue(util.Is5xx(response_.StatusCode)) {
				_err = tea.NewSDKError(map[string]interface{}{
					"code":       tea.ToString(DefaultAny(res["Code"], res["code"])),
					"statusCode": tea.IntValue(response_.StatusCode),
					"message":    "code: " + tea.ToString(tea.IntValue(response_.StatusCode)) + ", " + tea.ToString(DefaultAny(res["Message"], res["message"])) + " request id: " + tea.ToString(DefaultAny(res["RequestId"], res["requestId"])),
					"data":       res,
				})
				return _result, _err
			}

			_result = res
			return _result, _err
		}()
		if !tea.BoolValue(tea.Retryable(_err)) {
			break
		}
	}

	return _resp, _err
}

/**
 * Get user agent
 * @return user agent
 */
func (client *Client) GetUserAgent() (_result *string) {
	userAgent := util.GetUserAgent(client.UserAgent)
	_result = userAgent
	return _result
}

/**
 * Get accesskey id by using credential
 * @return accesskey id
 */
func (client *Client) GetAccessKeyId() (_result *string, _err error) {
	if tea.BoolValue(util.IsUnset(client.Credential)) {
		return _result, _err
	}

	accessKeyId, _err := client.Credential.GetAccessKeyId()
	if _err != nil {
		return _result, _err
	}

	_result = accessKeyId
	return _result, _err
}

/**
 * Get accesskey secret by using credential
 * @return accesskey secret
 */
func (client *Client) GetAccessKeySecret() (_result *string, _err error) {
	if tea.BoolValue(util.IsUnset(client.Credential)) {
		return _result, _err
	}

	secret, _err := client.Credential.GetAccessKeySecret()
	if _err != nil {
		return _result, _err
	}

	_result = secret
	return _result, _err
}

/**
 * Get security token by using credential
 * @return security token
 */
func (client *Client) GetSecurityToken() (_result *string, _err error) {
	if tea.BoolValue(util.IsUnset(client.Credential)) {
		return _result, _err
	}

	token, _err := client.Credential.GetSecurityToken()
	if _err != nil {
		return _result, _err
	}

	_result = token
	return _result, _err
}

/**
 * If the endpointRule and config.endpoint are empty, throw error
 * @param config config contains the necessary information to create a client
 */
func (client *Client) CheckConfig(config *Config) (_err error) {
	if tea.BoolValue(util.Empty(client.EndpointRule)) && tea.BoolValue(util.Empty(config.Endpoint)) {
		_err = tea.NewSDKError(map[string]interface{}{
			"code":    "ParameterMissing",
			"message": "'config.endpoint' can not be empty",
		})
		return _err
	}

	return _err
}

/**
 * If inputValue is not null, return it or return defaultValue
 * @param inputValue  users input value
 * @param defaultValue default value
 * @return the final result
 */
func DefaultAny(inputValue interface{}, defaultValue interface{}) (_result interface{}) {
	if tea.BoolValue(util.IsUnset(inputValue)) {
		_result = defaultValue
		return _result
	}

	_result = inputValue
	return _result
}

//The api client for prism's golang SDK
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	basic   = "basic"
	eTag    = "ETag"
	ifMatch = "If-Match"
)

var (
	jsonCheck               = regexp.MustCompile(`(?i:(?:application|text)/(?:vnd\.[^;]+\+)?json)`)
	xmlCheck                = regexp.MustCompile(`(?i:(?:application|text)/xml)`)
	uriCheck                = regexp.MustCompile(`/(?P<namespace>[-\w]+)/v\d+\.\d+(\.[a|b]\d+)?/(?P<suffix>.*)`)
	contentDispositionCheck = regexp.MustCompile("attachment;\\s*filename=\"(.*)\"")
	retryStatusList         = []int{408, 503, 504}
	userAgent               = "Nutanix-prism/v4.0.1-beta.1"
)

/*
  API client to handle the client-server communication, and is invariant across implementations.

    Scheme (optional) : URI scheme for connecting to the cluster (HTTP or HTTPS using SSL/TLS) (default : https)
    Host (required) : Host IPV4, IPV6 or FQDN for all http request made by this client (default : localhost)
    Port (optional) : Port for the host to connect to make all http request (default : 9440)
    Username (required) : Username to connect to a cluster
    Password (required) : Password to connect to a cluster
    Debug (optional) : flag to enable debug logging (default : empty)
    VerifySSL (optional) : Verify SSL certificate of cluster (default: true)
    MaxRetryAttempts (optional) : Maximum number of retry attempts to be made at a time (default: 5)
    ReadTimeout (optional) : Read timeout for all operations in milliseconds
    ConnectTimeout (optional) : Connection timeout for all operations in milliseconds
    RetryInterval (optional) : Interval between successive retry attempts (default: 3 sec)
    DownloadDirectory (optional) : Directory location on local for files to download (default: Current Directory)
    DownloadChunkSize (optional) : Chunk size in bytes for files to download (default: 8*1024 bytes)
    LoggerFile (optional) : Log file to write activity logs
*/
type ApiClient struct {
	Scheme            string `json:"scheme,omitempty"`
	Host              string `json:"host,omitempty"`
	Port              int    `json:"port,omitempty"`
	Username          string `json:"username,omitempty"`
	Password          string `json:"password,omitempty"`
	Debug             bool   `json:"debug,omitempty"`
	VerifySSL         bool
	Proxy             *Proxy
	MaxRetryAttempts  int           `json:"maxRetryAttempts,omitempty"`
	ReadTimeout       time.Duration `json:"readTimeout,omitempty"`
	ConnectTimeout    time.Duration `json:"connectTimeout,omitempty"`
	RetryInterval     time.Duration `json:"retryInterval,omitempty"`
	DownloadDirectory string        `json:"downloadDirectory,omitempty"`
	DownloadChunkSize int           `json:"downloadChunkSize,omitempty"`
	LoggerFile        string        `json:"loggerFile,omitempty"`
	defaultHeaders    map[string]string
	retryClient       *retryablehttp.Client
	httpClient        *http.Client
	dialer            *net.Dialer
	authentication    map[string]interface{}
	cookie            string
	refreshCookie     bool
	previousAuth      string
	basicAuth         *BasicAuth
	logger            *logrus.Logger

	// maxIdleConns controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	maxIdleConns int

	// maxIdleConnsPerHost, if non-zero, controls the maximum idle
	// (keep-alive) connections to keep per-host. If zero,
	// DefaultMaxIdleConnsPerHost is used.
	maxIdleConnsPerHost int

	// maxConnsPerHost optionally limits the total number of
	// connections per host, including connections in the dialing,
	// active, and idle states. On limit violation, dials will block.
	// Zero means no limit.
	maxConnsPerHost int

	// idleConnTimeout is the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing itself.
	// Zero means no limit.
	idleConnTimeout time.Duration

	// Timeout for the time spent during TLS handshake
	tlsHandshakeTimeout time.Duration
}

// Returns a newly generated ApiClient instance populated with default values
func NewApiClient() *ApiClient {

	basicAuth := new(BasicAuth)
	authentication := make(map[string]interface{})
	authentication["basicAuthScheme"] = basicAuth
	currentDirectory, _ := os.Getwd()

	a := &ApiClient{
		Scheme:              "https",
		Host:                "localhost",
		Port:                9440,
		Debug:               false,
		VerifySSL:           true,
		MaxRetryAttempts:    5,
		ReadTimeout:         30 * time.Second,
		ConnectTimeout:      30 * time.Second,
		RetryInterval:       3 * time.Second,
		DownloadDirectory:   currentDirectory,
		DownloadChunkSize:   8 * 1024,
		maxIdleConns:        10,
		maxIdleConnsPerHost: 10,
		maxConnsPerHost:     100,
		idleConnTimeout:     90 * time.Second,
		tlsHandshakeTimeout: 10 * time.Second,
		defaultHeaders:      make(map[string]string),
		refreshCookie:       true,
		basicAuth:           basicAuth,
		authentication:      authentication,
	}

	a.setupClient()
	return a
}

// Adds a default header to current api client instance for all the HTTP calls.
func (a *ApiClient) AddDefaultHeader(headerName string, headerValue string) {
	if headerName == "Authorization" {
		a.cookie = ""
	}

	a.defaultHeaders[headerName] = headerValue
}

// Makes the HTTP request with given options and returns response body as byte array.
func (a *ApiClient) CallApi(uri *string, httpMethod string, body interface{},
	queryParams url.Values, headerParams map[string]string, formParams url.Values,
	accepts []string, contentType []string, authNames []string) (interface{}, error) {
	path := a.Scheme + "://" + a.Host + ":" + strconv.Itoa(a.Port) + *uri

	if headerParams["Authorization"] != "" {
		a.previousAuth = headerParams["Authorization"]
	}

	if a.defaultHeaders["Authorization"] != "" {
		a.previousAuth = a.defaultHeaders["Authorization"]
	}

	// set Content-Type header
	if headerParams["Content-Type"] == "" {
		httpContentType := a.selectHeaderContentType(contentType)
		if httpContentType != "" {
			headerParams["Content-Type"] = httpContentType
		}
	}

	// set Accept header
	if headerParams["Accept"] == "" {
		httpHeaderAccept := a.selectHeaderAccept(accepts)
		if httpHeaderAccept == "" {
			httpHeaderAccept = "application/json"
		}

		headerParams["Accept"] = httpHeaderAccept
	}

	// set NTNX-Request-Id header
	_, requestIdHeaderExists := headerParams["NTNX-Request-Id"]
	_, requestIdDefaultHeaderExists := a.defaultHeaders["NTNX-Request-Id"]
	if !requestIdHeaderExists && !requestIdDefaultHeaderExists {
		headerParams["NTNX-Request-Id"] = uuid.New().String()
	}

	bodyValue := reflect.ValueOf(body)
	if bodyValue.IsValid() && !bodyValue.IsNil() {
		addEtagReferenceToHeader(body, headerParams)
	}

	request, err := a.prepareRequest(path, httpMethod, body, headerParams, queryParams, formParams, authNames)
	if err != nil {
		a.logger.Error(err.Error())
		return nil, err
	}

	a.setupClient()

	if a.Debug {
		printBody := true
		if headerParams["Content-Type"] == "application/octet-stream" {
			printBody = false
		}
		dump, err := httputil.DumpRequestOut(request, printBody)
		if err != nil {
			a.logger.Debugf("Error while logging request details: %s", err.Error())
			return nil, err
		}
		a.logger.Debug(string(dump))
	} else {
		a.logger.Infof("%s %s", request.Method, request.URL.String())
	}

	response, err := a.httpClient.Do(request)

	// Retry one more time without the cookie but with basic auth header
	if response != nil && response.StatusCode == 401 {
		a.logger.Debug("Retrying the request to refresh cookie...")
		request, _ = a.prepareRequest(path, httpMethod, body, headerParams, queryParams, formParams, authNames)
		a.refreshCookie = true
		if len(a.previousAuth) > 0 {
			request.Header["Authorization"] = []string{a.previousAuth}
		}
		delete(request.Header, "Cookie")

		dump, _ := httputil.DumpRequestOut(request, true)
		a.logger.Debug(string(dump))
		response, err = a.httpClient.Do(request)
	}

	if err != nil {
		a.logger.Error(err.Error())
		return nil, err
	}

	if a.Debug {
		printBody := true
		if response.Header.Get("Content-Type") == "application/octet-stream" {
			printBody = false
		}

		dump, err := httputil.DumpResponse(response, printBody)
		if err != nil {
			a.logger.Debugf("Error while logging response details: %s", err.Error())
			return nil, err
		}
		a.logger.Debug(string(dump))
	} else {
		a.logger.Infof("%s %s", response.Proto, response.Status)
	}

	if nil == response {
		msg := "Response is nil!"
		a.logger.Error(msg)
		return nil, ReportError(msg)
	}

	a.updateCookies(response)

	if response.StatusCode == 204 {
		return nil, nil
	}

	if response.Header.Get("Content-Type") == "application/octet-stream" {
		return a.downloadFile(response)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		a.logger.Error(err.Error())
		return nil, err
	}
	response.Body.Close()
	response.Body = ioutil.NopCloser(bytes.NewBuffer(responseBody))

	if !(200 <= response.StatusCode && response.StatusCode <= 209) {
		return nil, GenericOpenAPIError{
			Body:   responseBody,
			Status: response.Status,
		}
	} else {
		responseBody := addEtagReferenceToResponse(response.Header, responseBody)
		return responseBody, nil
	}
}

func (a *ApiClient) downloadFile(response *http.Response) (*string, error) {
	var filePath string
	if len(response.Header.Get("Content-Disposition")) != 0 {
		filename := contentDispositionCheck.FindStringSubmatch(response.Header.Get("Content-Disposition"))
		if len(filename) == 2 {
			filePath = filepath.Join(a.DownloadDirectory, filename[1])
		}
	} else {
		file, err := ioutil.TempFile(a.DownloadDirectory, "")
		if err != nil {
			a.logger.Errorf("Could not create a file on local for downloading: %s", err)
			return nil, err
		}
		filePath = file.Name()
		file.Close()
	}

	ext := filepath.Ext(filePath)
	ts := time.Now().UTC().Format("2006-01-02T15:04:05.000")
	filePath = fmt.Sprintf("%s_%s%s", filePath[:len(filePath)-len(ext)], ts, ext)

	a.logger.Infof("Writing response content to file at %s", filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	defer file.Close()
	if err != nil {
		a.logger.Errorf("Could not create a file on local for downloading: %s", filePath)
		return nil, err
	}

	buf := make([]byte, a.DownloadChunkSize)
	written, err := io.CopyBuffer(file, response.Body, buf)

	if err != nil {
		a.logger.Errorf("Something went wrong while downloading the file at path %s: %s", filePath, err)
		return nil, err
	}

	a.logger.Infof("%d bytes written to file at %s", written, filePath)

	return &filePath, nil
}

// Get all authentications (key: authentication name, value: authentication)
func (a *ApiClient) GetAuthentications() map[string]interface{} {
	return a.authentication
}

// Get authentication for the given auth name (eg : basic, oauth, bearer, apiKey)
func (a *ApiClient) GetAuthentication(authName string) interface{} {
	return a.authentication[authName]
}

// Helper method to set username for the first HTTP basic authentication.
func (a *ApiClient) SetUserName(username string) {
	a.Username = username
	a.basicAuth.UserName = username
}

// Helper method to set password for the first HTTP basic authentication
func (a *ApiClient) SetPassword(password string) {
	a.Password = password
	a.basicAuth.Password = password
}

// Helper method to set API key value for the first API key authentication
func (a *ApiClient) SetApiKey(key string) error {
	for _, value := range a.authentication {
		if auth, ok := value.(map[string]interface{}); ok {
			if apiKey, ok := auth["apiKey"].(*APIKey); ok {
				apiKey.Key = key
				return nil
			}
		}
	}

	return ReportError("no API key authentication configured!")
}

// Helper method to set API key prefix for the first API key authentication
func (a *ApiClient) SetApiKeyPrefix(apiKeyPrefix string) error {
	for _, value := range a.authentication {
		if auth, ok := value.(map[string]interface{}); ok {
			if apiKey, ok := auth["apiKey"].(*APIKey); ok {
				apiKey.Prefix = apiKeyPrefix
				return nil
			}
		}
	}

	return ReportError("no API key authentication configured!")
}

// Helper method to set access token for the first OAuth2 authentication.
func (a *ApiClient) SetAccessToken(accessToken string) error {
	for _, value := range a.authentication {
		if auth, ok := value.(*OAuth); ok {
			auth.AccessToken = accessToken
			return nil
		}
	}
	return ReportError("no OAuth2 authentication configured!")
}

// Helper method to set maximum retry attempts.
// After the initial instantiation of ApiClient, maximum retry attempts must be modified only via this method
func (a *ApiClient) SetMaxRetryAttempts(maxRetryAttempts int) {
	a.MaxRetryAttempts = maxRetryAttempts
}

func getValidTimeout(dur time.Duration, apiClient *ApiClient) time.Duration {
	if dur <= 0 {
		dur = 30 * time.Second
	} else if dur > (30 * time.Minute) {
		dur = 30 * time.Minute
	}

	return dur
}

// Helper method to enable/disable SSL verification. By default, SSL verification is enabled.
//
// Please note that disabling SSL verification is not recommended and should only be done for test purposes.
func (a *ApiClient) SetVerifySSL(verifySSL bool) {
	a.VerifySSL = verifySSL
}

// Helper method to set retry back off period.
// After the initial instantiation of ApiClient, back off period must be modified only via this method
func (a *ApiClient) SetRetryIntervalInMilliSeconds(ms int) {
	a.RetryInterval = time.Duration(ms) * time.Millisecond
}

func (a *ApiClient) setupClient() {
	var isRetryClientModified = false
	retryClientValue := reflect.ValueOf(a.retryClient)
	if !retryClientValue.IsValid() || retryClientValue.IsNil() {
		a.retryClient = retryablehttp.NewClient()
		isRetryClientModified = true
	}

	var transport = a.retryClient.HTTPClient.Transport.(*http.Transport)
	if isRetryClientModified || transport.TLSClientConfig == nil || transport.TLSClientConfig.InsecureSkipVerify != !a.VerifySSL ||
		a.dialer == nil || a.dialer.Timeout != a.ConnectTimeout {
		a.dialer = &net.Dialer{
			Timeout: getValidTimeout(a.ConnectTimeout, a),
		}
		transport := &http.Transport{
			DialContext:           a.dialer.DialContext,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: !a.VerifySSL},
			MaxIdleConns:          a.maxIdleConns,
			MaxIdleConnsPerHost:   a.maxIdleConnsPerHost,
			MaxConnsPerHost:       a.maxConnsPerHost,
			IdleConnTimeout:       a.idleConnTimeout,
			TLSHandshakeTimeout:   a.tlsHandshakeTimeout,
			ExpectContinueTimeout: 1 * time.Second,
		}
		if (a.Proxy != nil) && (*a.Proxy != Proxy{}) {
			path := a.Proxy.Host
			if a.Proxy.Port != 0 {
				path = path + ":" + strconv.Itoa(a.Proxy.Port)
			}
			transport.Proxy = http.ProxyURL(&url.URL{
				Scheme: a.Proxy.Scheme,
				User:   url.UserPassword(a.Proxy.Username, a.Proxy.Password),
				Host:   path,
			})
		}
		a.retryClient.HTTPClient.Transport = transport
		isRetryClientModified = true
	}

	if a.retryClient.RetryMax != a.MaxRetryAttempts ||
		a.retryClient.RetryWaitMax != a.RetryInterval ||
		a.retryClient.CheckRetry == nil {
		isRetryClientModified = true
	}

	a.retryClient.RetryMax = a.MaxRetryAttempts
	a.retryClient.RetryWaitMax = a.RetryInterval
	a.retryClient.CheckRetry = retryPolicy

	configureLogger(a)

	if isRetryClientModified {
		a.httpClient = a.retryClient.StandardClient()
	}

	a.httpClient.Timeout = getValidTimeout(a.ConnectTimeout, a) + a.tlsHandshakeTimeout + getValidTimeout(a.ReadTimeout, a)
}

func configureLogger(a *ApiClient) {
	a.retryClient.Logger = nil

	logLevel := logrus.InfoLevel
	if a.Debug {
		logLevel = logrus.DebugLevel
	}

	var output io.Writer
	if a.LoggerFile == "" {
		output = os.Stderr
	} else {
		f, _ := os.OpenFile(a.LoggerFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		output = io.MultiWriter(os.Stderr, f)
	}

	if a.logger == nil {
		a.logger = &logrus.Logger{
			Out:   output,
			Level: logLevel,
			Formatter: &myFormatter{
				logrus.TextFormatter{
					FullTimestamp:          true,
					TimestampFormat:        "2006-01-02 15:04:05.000",
					ForceColors:            true,
					DisableLevelTruncation: true,
				},
			},
		}
	} else {
		a.logger.Out = output
		a.logger.Level = logLevel
	}
}

// Custom formatter for logrus
type myFormatter struct {
	logrus.TextFormatter
}

// Format function implementation for the Formatter interface of logrus
func (f *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteString(entry.Time.Format(f.TimestampFormat))
	b.WriteByte(' ')
	b.WriteString(strings.ToUpper(entry.Level.String()))

	if entry.Message != "" {
		b.WriteString(" - ")
		b.WriteString(entry.Message)
	}

	if len(entry.Data) > 0 {
		b.WriteString(" || ")
	}

	for key, value := range entry.Data {
		b.WriteString(key)
		b.WriteByte('=')
		fmt.Fprint(b, value)
		b.WriteString(", ")
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

// Select a content type from the available list.
func (a *ApiClient) selectHeaderContentType(contentTypes []string) string {
	if len(contentTypes) == 0 {
		return ""
	}
	if contains(contentTypes, "application/json") {
		return "application/json"
	}
	return contentTypes[0] // use the first content type specified in 'consumes'
}

// Join all accept types and return
func (a *ApiClient) selectHeaderAccept(accepts []string) string {
	if len(accepts) == 0 {
		return ""
	}

	return strings.Join(accepts, ",")
}

// Build the request
func (a *ApiClient) prepareRequest(
	path string, method string,
	postBody interface{},
	headerParams map[string]string,
	queryParams url.Values,
	formParams url.Values,
	authNames []string) (localVarRequest *http.Request, err error) {

	var body *bytes.Buffer

	// Detect postBody type and post.
	postBodyValue := reflect.ValueOf(postBody)
	if postBodyValue.IsValid() && !postBodyValue.IsNil() {
		contentType := headerParams["Content-Type"]
		if contentType == "" {
			contentType = detectContentType(postBody)
			headerParams["Content-Type"] = contentType
		}

		if headerParams["Content-Type"] != "application/octet-stream" {
			body, err = setBody(postBody, contentType)
			if err != nil {
				return nil, err
			}
		}
	}

	if strings.HasPrefix(headerParams["Content-Type"], "application/x-www-form-urlencoded") && len(formParams) > 0 {
		if body != nil {
			return nil, errors.New("Cannot specify postBody and x-www-form-urlencoded form at the same time.")
		}
		body = &bytes.Buffer{}
		body.WriteString(formParams.Encode())
		// Set Content-Length
		headerParams["Content-Length"] = fmt.Sprintf("%d", body.Len())
	}

	// Setup path and query parameters
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Adding Query Param
	query := url.Query()
	for k, v := range queryParams {
		for _, iv := range v {
			query.Add(k, iv)
		}
	}

	// Encode the parameters.
	url.RawQuery = query.Encode()

	// Generate a new request
	if postBodyValue.IsValid() && !postBodyValue.IsNil() && headerParams["Content-Type"] == "application/octet-stream" {
		fileTypeBody, _ := postBody.(*os.File)
		localVarRequest, err = http.NewRequest(method, url.String(), fileTypeBody)
	} else if body != nil {
		localVarRequest, err = http.NewRequest(method, url.String(), body)
	} else {
		localVarRequest, err = http.NewRequest(method, url.String(), nil)
	}

	if err != nil {
		return nil, err
	}

	// add header parameters, if any
	if len(headerParams) > 0 {
		headers := http.Header{}
		for h, v := range headerParams {
			headers[h] = []string{v}
		}
		localVarRequest.Header = headers
	}

	// Add the user agent to the request.
	localVarRequest.Header["User-Agent"] = []string{userAgent}

	// Authentication
	a.SetUserName(a.Username)
	a.SetPassword(a.Password)
	for _, authName := range authNames {
		// Basic HTTP authentication
		if auth, ok := a.authentication[authName].(*BasicAuth); ok {
			if auth.UserName != "" && auth.Password != "" {
				localVarRequest.SetBasicAuth(auth.UserName, auth.Password)
			}
			// API Key authentication
		} else if auth, ok := a.authentication[authName].(map[string]interface{}); ok {
			var key string
			if apiKey, ok := auth["apiKey"].(*APIKey); ok && apiKey.Prefix != "" {
				key = apiKey.Prefix + " " + apiKey.Key
			} else {
				key = apiKey.Key
			}

			if auth["in"] == "header" {
				localVarRequest.Header[auth["name"].(string)] = []string{key}
			}
			if auth["in"] == "query" {
				queries := localVarRequest.URL.Query()
				queries.Add(auth["name"].(string), key)
				localVarRequest.URL.RawQuery = queries.Encode()
			}
			// OAuth or Bearer authentication
		} else if auth, ok := a.authentication[authName].(*OAuth); ok {
			localVarRequest.Header["Authorization"] = []string{"Bearer " + auth.AccessToken}
		} else {
			return nil, ReportError("unknown authentication type: %s", authName)
		}
	}

	for header, value := range a.defaultHeaders {
		localVarRequest.Header[header] = []string{value}
	}

	a.previousAuth = localVarRequest.Header.Get("Authorization")

	// Add the cookie to the request.
	if len(a.cookie) > 0 {
		localVarRequest.Header["Cookie"] = []string{a.cookie}
		delete(localVarRequest.Header, "Authorization")
	}

	// Add content length to request
	if localVarRequest.Header.Get("Content-Length") != "" {
		contentLengthInt64, _ := strconv.ParseInt(localVarRequest.Header.Get("Content-Length"), 10, 64)
		localVarRequest.ContentLength = contentLengthInt64
	}

	return localVarRequest, nil
}

// RetryPolicy provides a callback for Client.CheckRetry, specifies retry on
// error codes mentioned in RetryStatusList
func retryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	for _, status := range retryStatusList {
		if resp.StatusCode == status {
			return true, nil
		}
	}
	return false, nil
}

// Case-insensitive match for finding an item in the list.
func contains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if strings.ToLower(a) == strings.ToLower(needle) {
			return true
		}
	}
	return false
}

// Figure out `Request.Body` content type for request header
func detectContentType(body interface{}) string {
	contentType := "text/plain; charset=utf-8"
	kind := reflect.TypeOf(body).Kind()

	switch kind {
	case reflect.Struct, reflect.Map, reflect.Ptr:
		contentType = "application/json; charset=utf-8"
	case reflect.String:
		contentType = "text/plain; charset=utf-8"
	default:
		if b, ok := body.([]byte); ok {
			contentType = http.DetectContentType(b)
		} else if kind == reflect.Slice {
			contentType = "application/json; charset=utf-8"
		}
	}

	return contentType
}

// Read ETag and add it to If-Match header
func addEtagReferenceToHeader(body interface{}, headerParams map[string]string) {
	if reflect.ValueOf(body).Elem().Kind() == reflect.Struct {
		if reserved := reflect.ValueOf(body).Elem().FieldByName("Reserved_"); reserved.IsValid() {
			reservedMap := reserved.Interface().(map[string]interface{})
			if etag, etagOk := reservedMap[eTag].(string); etagOk {
				headerParams[ifMatch] = etag
			}
		}
	}
}

// Get ETag from an object if exists, otherwise returns empty string.
// The ETag is usually provided in the response of the GET API calls, which can further be used in other HTTP operations.
func (a *ApiClient) GetEtag(object interface{}) string {
	var reserved reflect.Value
	if reflect.TypeOf(object).Kind() == reflect.Struct {
		reserved = reflect.ValueOf(object).FieldByName("Reserved_")
	} else if reflect.TypeOf(object).Kind() == reflect.Interface || reflect.TypeOf(object).Kind() == reflect.Ptr {
		reserved = reflect.ValueOf(object).Elem().FieldByName("Reserved_")
	} else {
		a.logger.Warnf("Unrecognized input type %s for %s to retrieve etag!", reflect.TypeOf(object).Kind(), object)
		return ""
	}

	if reserved.IsValid() {
		etagKey := strings.ToLower(eTag)
		reservedMap := reserved.Interface().(map[string]interface{})
		for k, v := range reservedMap {
			if strings.ToLower(k) == etagKey {
				return v.(string)
			}
		}
	}

	return ""
}

// Read ETag and add it to response
func addEtagReferenceToResponse(headers http.Header, body []byte) []byte {
	if etag := headers.Get(eTag); etag != "" {
		responseMap := map[string]interface{}{}
		json.Unmarshal(body, &responseMap)
		if r, ok := responseMap["$reserved"].(map[string]interface{}); ok {
			r[eTag] = etag
			if d, ok := responseMap["data"].(map[string]interface{}); ok {
				if r2, ok := d["$reserved"].(map[string]interface{}); ok {
					r2[eTag] = etag
					m, _ := json.Marshal(responseMap)
					return m
				}
			} else if dList, ok := responseMap["data"].([]interface{}); ok {
				for _, d := range dList {
					if d, ok := d.(map[string]interface{}); ok {
						if r3, ok := d["$reserved"].(map[string]interface{}); ok {
							r3[eTag] = etag
						}
					}
				}
				m, _ := json.Marshal(responseMap)
				return m
			}
		}
	}
	return body
}

// Set request body from an interface{}
func setBody(body interface{}, contentType string) (bodyBuf *bytes.Buffer, err error) {
	if nil == bodyBuf {
		bodyBuf = &bytes.Buffer{}
	}

	if reader, ok := body.(io.Reader); ok {
		_, err = bodyBuf.ReadFrom(reader)
	} else if fp, ok := body.(**os.File); ok {
		_, err = bodyBuf.ReadFrom(*fp)
	} else if b, ok := body.([]byte); ok {
		_, err = bodyBuf.Write(b)
	} else if s, ok := body.(string); ok {
		_, err = bodyBuf.WriteString(s)
	} else if s, ok := body.(*string); ok {
		_, err = bodyBuf.WriteString(*s)
	} else if jsonCheck.MatchString(contentType) {
		err = json.NewEncoder(bodyBuf).Encode(body)
	} else if xmlCheck.MatchString(contentType) {
		err = xml.NewEncoder(bodyBuf).Encode(body)
	}

	if err != nil {
		return nil, err
	}

	if bodyBuf.Len() == 0 {
		err = fmt.Errorf("Invalid body type %s\n", contentType)
		return nil, err
	}
	return bodyBuf, nil
}

// Set Cookie information to reuse in subsequent requests for a valid response
func (a *ApiClient) updateCookies(response *http.Response) {
	if a.refreshCookie {
		cookiesList := response.Header["Set-Cookie"]
		if len(cookiesList) > 0 {
			cookieFromResponse := ""
			for _, value := range cookiesList {
				finalCookie := strings.SplitN(value, ";", 2)[0]
				if strings.Contains(finalCookie, "=") {
					finalCookie = strings.TrimSpace(finalCookie)
				} else {
					finalCookie = ""
				}

				if finalCookie != "" {
					cookieFromResponse += finalCookie + ";"
				}
			}

			// Remove trailing ";"
			if cookieFromResponse != "" {
				cookieFromResponse = strings.TrimSuffix(cookieFromResponse, ";")
			}

			a.cookie = cookieFromResponse
			a.refreshCookie = false
		}
	}
}

// Provides basic http authentication to a request passed via context using ContextBasicAuth
type BasicAuth struct {
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
}

/*
  Configuration for the Proxy Server that requests are to be routed through.

    Scheme: URI Scheme for connecting to the proxy ("http", "https" or "socks5")
    Host: Host of the proxy to which the client will connect to
    Port: Port of the proxy to which the client will connect to
    Username: Username to connect to the proxy
    Password: Password to connect to the proxy
*/
type Proxy struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Scheme   string `json:"scheme,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
}

// Provides API key based authentication to a request passed via context using ContextAPIKey
type APIKey struct {
	Key    string
	Prefix string
}

// Provides OAuth authentication
type OAuth struct {
	AccessToken string
}

// Provides access to the body (error), status and model on returned errors.
type GenericOpenAPIError struct {
	Body   []byte
	Model  interface{}
	Status string
}

// Returns non-empty string if there was an error.
func (e GenericOpenAPIError) Error() string {
	return string(e.Body)
}

// Returns deserialized response body if compatible with GenericOpenAPIError.Model
func (e GenericOpenAPIError) DeserializedModel() interface{} {
	err := json.Unmarshal(e.Body, e.Model)
	if err != nil {
		return nil
	}
	return e.Model
}

// Convert interface{} parameters to string, using a delimiter if format is provided.
func ParameterToString(obj interface{}, collectionFormat string) string {
	var delimiter string

	switch collectionFormat {
	case "pipes":
		delimiter = "|"
	case "ssv":
		delimiter = " "
	case "tsv":
		delimiter = "\t"
	case "csv":
		delimiter = ","
	}

	if reflect.TypeOf(obj).Kind() == reflect.Slice {
		return strings.Trim(strings.Replace(fmt.Sprint(obj), " ", delimiter, -1), "[]")
	} else if t, ok := obj.(time.Time); ok {
		return t.Format(time.RFC3339)
	}

	return fmt.Sprintf("%v", obj)
}

// Helper for converting interface{} parameters to json strings
func ParameterToJson(obj interface{}) (string, error) {
	jsonBuf, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonBuf), err
}

// Prevent trying to import "fmt"
func ReportError(format string, b ...interface{}) error {
	return fmt.Errorf(format, b...)
}

type LeveledLogrus struct {
	*logrus.Logger
}

func (l *LeveledLogrus) Error(msg string, keysAndValues ...interface{}) {
	l.WithFields(fields(keysAndValues)).Error(msg)
}

func (l *LeveledLogrus) Info(msg string, keysAndValues ...interface{}) {
	l.WithFields(fields(keysAndValues)).Info(msg)
}
func (l *LeveledLogrus) Debug(msg string, keysAndValues ...interface{}) {
	l.WithFields(fields(keysAndValues)).Debug(msg)
}

func (l *LeveledLogrus) Warn(msg string, keysAndValues ...interface{}) {
	l.WithFields(fields(keysAndValues)).Warn(msg)
}

func fields(keysAndValues []interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	for i := 0; i < len(keysAndValues)-1; i += 2 {
		fields[keysAndValues[i].(string)] = keysAndValues[i+1]
	}

	return fields
}

func (l *LeveledLogrus) setLoggerFilePath(filename string) error {
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		l.Error("Error opening log file", "error", err)
		return err
	}

	l.SetOutput(logFile)
	l.SetReportCaller(true)
	return nil
}

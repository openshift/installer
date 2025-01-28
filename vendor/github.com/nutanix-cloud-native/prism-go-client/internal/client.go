package internal

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/hashicorp/go-cleanhttp"
	"go.uber.org/zap"

	"github.com/nutanix-cloud-native/prism-go-client"
)

type Scheme string

const (
	defaultBaseURL  = "%s://%s/"
	mediaType       = "application/json"
	formEncodedType = "application/x-www-form-urlencoded"
	octetStreamType = "application/octet-stream"

	SchemeHTTP  Scheme = "http"
	SchemeHTTPS Scheme = "https"
)

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response, interface{})

// Client Config Configuration of the httpClient
type Client struct {
	credentials *prismgoclient.Credentials

	// HTTP httpClient used to communicate with the Nutanix API.
	httpClient *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for httpClient
	UserAgent string

	cookies []*http.Cookie

	// Optional function called after every successful request made.
	onRequestCompleted RequestCompletionCallback

	// absolutePath: for example api/nutanix/v3
	absolutePath string

	// error message, incase httpClient is in error state
	ErrorMsg string

	logger   *logr.Logger
	certpool *x509.CertPool
}

type ClientOption func(*Client) error

// WithLogger sets the logger for the httpClient
func WithLogger(logger *logr.Logger) ClientOption {
	return func(c *Client) error {
		c.logger = logger
		return nil
	}
}

// WithCredentials sets the credentials for the httpClient
func WithCredentials(credentials *prismgoclient.Credentials) ClientOption {
	return func(c *Client) error {
		c.credentials = credentials
		if c.credentials.Insecure {
			transport, ok := c.httpClient.Transport.(*http.Transport)
			if !ok {
				return fmt.Errorf("transport is not of type http.Transport: %T", c.httpClient.Transport)
			}
			transport.TLSClientConfig.InsecureSkipVerify = true
		}
		if c.credentials.ProxyURL != "" {
			c.logger.V(1).Info("Using proxy:", "proxy", c.credentials.ProxyURL)
			proxy, err := url.Parse(c.credentials.ProxyURL)
			if err != nil {
				return fmt.Errorf("error parsing proxy url: %s", err)
			}
			transport, ok := c.httpClient.Transport.(*http.Transport)
			if !ok {
				return fmt.Errorf("transport is not of type http.Transport: %T", c.httpClient.Transport)
			}
			transport.Proxy = http.ProxyURL(proxy)
		}
		return nil
	}
}

// WithUserAgent sets the user agent for the httpClient
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithBaseURL sets the base URL for the httpClient to communicate with
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		// if the BaseURL does not have a scheme, use https as a default scheme
		// the `url.Parse` function parses the base URL (i.e. Host) as a Path
		// if the URL does not have a scheme. Prefixing a scheme ensures the base URL
		// is parsed as a Host and not a Path.
		if !strings.HasPrefix(baseURL, string(SchemeHTTP)) {
			baseURL = fmt.Sprintf("%s://%s", SchemeHTTPS, baseURL)
		}
		u, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		if u.Path == "" {
			u.Path = "/"
		}
		c.BaseURL = u
		return nil
	}
}

// WithCookies sets the cookies for the httpClient
func WithCookies(cookies []*http.Cookie) ClientOption {
	return func(c *Client) error {
		c.cookies = cookies
		return nil
	}
}

// WithAbsolutePath sets the absolute path for the httpClient to communicate with
func WithAbsolutePath(absolutePath string) ClientOption {
	return func(c *Client) error {
		c.absolutePath = absolutePath
		return nil
	}
}

// WithCertificate adds the certificate to the certificate pool in tls config
func WithCertificate(cert *x509.Certificate) ClientOption {
	return func(c *Client) error {
		if cert == nil {
			return fmt.Errorf("certificate is nil")
		}
		c.certpool.AddCert(cert)
		return nil
	}
}

// WithRoundTripper overrides the transport for the httpClient
// Overriding transport is useful for testing against API Mocks
// This is not recommended for production use
func WithRoundTripper(transport http.RoundTripper) ClientOption {
	return func(c *Client) error {
		c.httpClient.Transport = transport
		return nil
	}
}

// NewClient returns a wrapper around http/https (as per isHTTP flag) httpClient with additions of proxy & session_auth if given
func NewClient(opts ...ClientOption) (*Client, error) {
	c := &Client{
		httpClient: cleanhttp.DefaultClient(),
	}

	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to get system cert pool: %s", err)
	}
	c.certpool = certPool

	c.httpClient.Transport = http.DefaultTransport
	c.httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{}
	c.httpClient.Transport.(*http.Transport).TLSClientConfig.RootCAs = c.certpool

	// If the user does not specify a logger, then we'll use zap for a default one
	// If the user specified a logger, then we'll use that logger
	zapLog, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	logger := zapr.NewLogger(zapLog)
	c.logger = &logger

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.UserAgent == "" {
		return nil, fmt.Errorf("userAgent argument must be passed")
	}
	if c.absolutePath == "" {
		return nil, fmt.Errorf("absolutePath argument must be passed")
	}
	if c.credentials == nil {
		return nil, fmt.Errorf("credentials argument must be passed")
	}
	if c.BaseURL == nil {
		c.logger.V(1).Info("BaseURL is not set. Using URL from credentials", "url", c.credentials.URL)
		if err := WithBaseURL(c.credentials.URL)(c); err != nil {
			return nil, err
		}
	}

	if c.credentials.SessionAuth {
		c.logger.V(1).Info("Using session_auth")
		if err := c.refreshCookies(context.Background()); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// NewRequest creates a request
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	req, err := c.NewUnAuthRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	if c.cookies != nil {
		decorateRequestWithCookies(req, c.cookies)
	} else {
		decorateRequestWithBasicAuthHeaders(req, c.credentials.Username, c.credentials.Password)
	}

	return req, nil
}

func (c *Client) refreshCookies(ctx context.Context) error {
	req, err := c.NewUnAuthRequest(http.MethodGet, "/users/me", nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	decorateRequestWithBasicAuthHeaders(req, c.credentials.Username, c.credentials.Password)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := CheckResponse(resp); err != nil {
		return err
	}

	c.cookies = resp.Cookies()

	return nil
}

func decorateRequestWithBasicAuthHeaders(req *http.Request, username, password string) {
	req.Header.Add("Authorization",
		"Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
}

func decorateRequestWithCookies(req *http.Request, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
}

func clearCookiesInRequest(req *http.Request) {
	req.Header.Del("Cookie")
}

// NewUnAuthRequest creates a request without authorisation headers
func (c *Client) NewUnAuthRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	// check if httpClient exists or not
	if c.httpClient == nil {
		return nil, fmt.Errorf(c.ErrorMsg)
	}

	// create main api url
	rel, err := url.Parse(c.absolutePath + urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// add api headers
	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

// NewUnAuthFormEncodedRequest returns content-type: application/x-www-form-urlencoded based unauth request
func (c *Client) NewUnAuthFormEncodedRequest(method, urlStr string, body map[string]string) (*http.Request, error) {
	// check if httpClient exists or not
	if c.httpClient == nil {
		return nil, fmt.Errorf(c.ErrorMsg)
	}
	// create main api url
	rel, err := url.Parse(c.absolutePath + urlStr)
	if err != nil {
		return nil, err
	}
	u := c.BaseURL.ResolveReference(rel)

	// create form data from body
	data := url.Values{}
	for k, v := range body {
		data.Set(k, v)
	}

	// create a new request based on encoded from data
	req, err := http.NewRequest(method, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// add api headers
	req.Header.Add("Content-Type", formEncodedType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

// NewUploadRequest Handles image uploads for image service
func (c *Client) NewUploadRequest(method, urlStr string, fileReader *os.File) (*http.Request, error) {
	// check if httpClient exists or not
	if c.httpClient == nil {
		return nil, fmt.Errorf(c.ErrorMsg)
	}
	rel, errp := url.Parse(c.absolutePath + urlStr)
	if errp != nil {
		return nil, errp
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), fileReader)
	if err != nil {
		return nil, err
	}

	// Set req.ContentLength and req.GetBody as internally there is no implementation of such for os.File type reader
	// http.NewRequest() only sets this values for types - bytes.Buffer, bytes.Reader and strings.Reader
	// Refer https://github.com/golang/go/blob/a0f77e56b7a7ecb92dca3e2afdd56ee773c2cb07/src/net/http/request.go#L896
	fileInfo, err := fileReader.Stat()
	if err != nil {
		return nil, err
	}
	req.ContentLength = fileInfo.Size()
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(fileReader), nil
	}

	req.Header.Add("Content-Type", octetStreamType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Authorization", "Basic "+
		base64.StdEncoding.EncodeToString([]byte(c.credentials.Username+":"+c.credentials.Password)))

	return req, nil
}

// NewUploadRequest handles image uploads for image service
func (c *Client) NewUnAuthUploadRequest(method, urlStr string, fileReader *os.File) (*http.Request, error) {
	// check if httpClient exists or not
	if c.httpClient == nil {
		return nil, fmt.Errorf(c.ErrorMsg)
	}
	rel, errp := url.Parse(c.absolutePath + urlStr)
	if errp != nil {
		return nil, errp
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), fileReader)
	if err != nil {
		return nil, err
	}

	// Set req.ContentLength and req.GetBody as internally there is no implementation of such for os.File type reader
	// http.NewRequest() only sets this values for types - bytes.Buffer, bytes.Reader and strings.Reader
	// Refer https://github.com/golang/go/blob/a0f77e56b7a7ecb92dca3e2afdd56ee773c2cb07/src/net/http/request.go#L896
	fileInfo, err := fileReader.Stat()
	if err != nil {
		return nil, err
	}
	req.ContentLength = fileInfo.Size()
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(fileReader), nil
	}

	req.Header.Add("Content-Type", octetStreamType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// OnRequestCompleted sets the DO API request completion callback
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

func (c *Client) clearCookies() {
	c.cookies = nil
}

// do will perform the request and handle retries. This is a recursive function
// that will retry the request if the response is a 401 and the client is configured
// to use session authentication. It will retry up to maxRetries times.
// do performs the request and handles retries. It is not exported.
func (c *Client) do(ctx context.Context, req *http.Request, v interface{}, retryCount int, maxRetries int) error {
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := CheckResponse(resp); err != nil {
		if resp.StatusCode == http.StatusUnauthorized && retryCount < maxRetries {
			c.clearCookies()
			if err := c.refreshCookies(ctx); err != nil {
				return err
			}
			clearCookiesInRequest(req)
			decorateRequestWithCookies(req, c.cookies)
			return c.do(ctx, req, v, retryCount+1, maxRetries)
		}
		return err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
		if err != nil {
			return fmt.Errorf("error unmarshalling json: %s", err)
		}
	}

	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp, v)
	}
	return nil
}

// Do performs request passed
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {
	noRetries := 0
	retryOnce := 1

	// It's recommend that the clients use session-auth cookies instead of basic auth for requests to Prism Central
	// where possible. Session tokens are alive for 15m and then replaced. The session-auth handling mechanism in
	// the client should be resilient to session token refresh. As long the credentials are accurate and session
	// auth mechanism is in play, the client upon receiving the first 401 should refresh the session cookie and
	// retry the request. Only if the requests fails the second time, should the client return a 401 back to the caller.
	if c.credentials.SessionAuth {
		return c.do(ctx, req, v, noRetries, retryOnce)
	}

	return c.do(ctx, req, v, noRetries, noRetries)
}

func searchSlice(slice []string, key string) bool {
	for _, v := range slice {
		if v == key {
			return true
		}
	}
	return false
}

// DoWithFilters performs request passed and filters entities in json response
func (c *Client) DoWithFilters(ctx context.Context, req *http.Request, v interface{}, filters []*prismgoclient.AdditionalFilter, baseSearchPaths []string) error {
	// check if httpClient exists or not
	if c.httpClient == nil {
		return fmt.Errorf(c.ErrorMsg)
	}
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if err := CheckResponse(resp); err != nil {
		return err
	}

	resp.Body, err = filter(resp.Body, filters, baseSearchPaths)
	if err != nil {
		return err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				fmt.Printf("Error io.Copy %s", err)
				return err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return fmt.Errorf("error unmarshalling json: %s", err)
			}
		}
	}

	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp, v)
	}

	return err
}

func filter(body io.ReadCloser, filters []*prismgoclient.AdditionalFilter, baseSearchPaths []string) (io.ReadCloser, error) {
	if filters == nil {
		return body, nil
	}
	if len(baseSearchPaths) == 0 {
		baseSearchPaths = []string{"$."}
	}

	var res map[string]interface{}
	b, err := io.ReadAll(body)
	if err != nil {
		return body, err
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	// Full search paths
	searchPaths := map[string][]string{}
	filterMap := map[string]*prismgoclient.AdditionalFilter{}
	for _, filter := range filters {
		filterMap[filter.Name] = filter
		// Build search paths by appending target search paths to base paths
		filterSearchPaths := []string{}
		for _, baseSearchPath := range baseSearchPaths {
			searchPath := fmt.Sprintf("%s.%s", baseSearchPath, filter.Name)
			filterSearchPaths = append(filterSearchPaths, searchPath)
		}
		searchPaths[filter.Name] = filterSearchPaths
	}

	// Entities that pass filters
	var filteredEntities []interface{}

	entities := res["entities"].([]interface{})
	for _, entity := range entities {
		filtersPassed := 0
	filter_loop:
		for filter, filterSearchPaths := range searchPaths {
			for _, searchPath := range filterSearchPaths {
				searchTarget := entity.(map[string]interface{})
				val, err := jsonpath.Get(searchPath, searchTarget)
				if err != nil {
					continue
				}
				// Stringify leaf value since we support only string values in filter
				value := fmt.Sprint(val)
				if searchSlice(filterMap[filter].Values, value) {
					filtersPassed++
					continue filter_loop
				}
			}
		}

		// Value must pass all filters since we perform logical AND b/w filters
		if filtersPassed == len(filters) {
			filteredEntities = append(filteredEntities, entity)
		}
	}
	res["entities"] = filteredEntities

	// Convert filtered result back to io.ReadCloser
	filteredBody, jsonErr := json.Marshal(res)
	if jsonErr != nil {
		return body, jsonErr
	}
	return io.NopCloser(bytes.NewReader(filteredBody)), nil
}

// CheckResponse checks errors if exist errors in request
func CheckResponse(r *http.Response) error {
	c := r.StatusCode

	if c >= 200 && c <= 299 {
		return nil
	}

	// Nutanix returns non-json response with code 401 when
	// invalid credentials are used
	if c == http.StatusUnauthorized {
		return fmt.Errorf("invalid Nutanix credentials")
	}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	rdr2 := io.NopCloser(bytes.NewBuffer(buf))

	r.Body = rdr2
	// if has entities -> return nil
	// if has message_list -> check_error["state"]
	// if has status -> check_error["status.state"]
	if len(buf) == 0 {
		return nil
	}

	var res map[string]interface{}
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return fmt.Errorf("unmarshalling error response %s for response body %s", err, string(buf))
	}

	errRes := &ErrorResponse{}
	if status, ok := res["status"]; ok {
		_, sok := status.(string)
		if sok {
			return nil
		}

		err = fillStruct(status.(map[string]interface{}), errRes)
	} else if _, ok := res["state"]; ok {
		err = fillStruct(res, errRes)
	} else if _, ok := res["entities"]; ok {
		return nil
	}

	if err != nil {
		return err
	}

	// karbon error check
	if messageInfo, ok := res["message_info"]; ok {
		return fmt.Errorf("error: %s", messageInfo)
	}

	// This check is also used for some foundation api errors
	if message, ok := res["message"]; ok {
		log.Print(message)
		return fmt.Errorf("error: %s", message)
	}
	if errRes.State != "ERROR" {
		return nil
	}

	pretty, err := json.MarshalIndent(errRes, "", "  ")
	if err != nil {
		return fmt.Errorf("status: %s, error-response: %+v, marshal error: %v", r.Status, errRes, err)
	}
	return fmt.Errorf("status: %s, error-response: %s", r.Status, string(pretty))
}

// ErrorResponse ...
type ErrorResponse struct {
	APIVersion  string            `json:"api_version,omitempty"`
	Code        int64             `json:"code,omitempty"`
	Kind        string            `json:"kind,omitempty"`
	MessageList []MessageResource `json:"message_list"`
	State       string            `json:"state"`
}

// MessageResource ...
type MessageResource struct {
	// Custom key-value details relevant to the status.
	Details interface{} `json:"details,omitempty"`

	// If state is ERROR, a message describing the error.
	Message string `json:"message"`

	// If state is ERROR, a machine-readable snake-cased *string.
	Reason string `json:"reason"`
}

func (r *ErrorResponse) Error() string {
	err := ""
	for key, value := range r.MessageList {
		err = fmt.Sprintf("%d: {message:%s, reason:%s }", key, value.Message, value.Reason)
	}

	return err
}

func fillStruct(data map[string]interface{}, result interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(j, result)
}

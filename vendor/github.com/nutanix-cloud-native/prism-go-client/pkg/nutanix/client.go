package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

const (
	// libraryVersion = "v3"
	defaultBaseURL = "%s://%s/"
	httpPrefix     = "http"
	httpsPrefix    = "https"
	// absolutePath   = "api/nutanix/" + libraryVersion
	// userAgent      = "nutanix/" + libraryVersion
	mediaType       = "application/json"
	formEncodedType = "application/x-www-form-urlencoded"
	octetStreamType = "application/octet-stream"
)

// Client Config Configuration of the client
type Client struct {
	Credentials *Credentials

	// HTTP client used to communicate with the Nutanix API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	Cookies []*http.Cookie

	// Optional function called after every successful request made.
	onRequestCompleted RequestCompletionCallback

	// absolutePath: for example api/nutanix/v3
	AbsolutePath string

	// error message, incase client is in error state
	ErrorMsg string
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response, interface{})

// Credentials needed username and password
type Credentials struct {
	URL                string
	Username           string
	Password           string
	Endpoint           string
	Port               string
	Insecure           bool
	SessionAuth        bool
	ProxyURL           string
	FoundationEndpoint string              // Required field for connecting to foundation VM APIs
	FoundationPort     string              // Port for connecting to foundation VM APIs
	RequiredFields     map[string][]string // RequiredFields is client to its required fields mapping for validations and usage in every client
}

// AdditionalFilter specification for client side filters
type AdditionalFilter struct {
	Name   string
	Values []string
}

// NewClient returns a wrapper around http/https (as per isHTTP flag) client with additions of proxy & session_auth if given
func NewClient(credentials *Credentials, userAgent string, absolutePath string, isHTTP bool) (*Client, error) {
	if userAgent == "" {
		return nil, fmt.Errorf("userAgent argument must be passed")
	}
	if absolutePath == "" {
		return nil, fmt.Errorf("absolutePath argument must be passed")
	}

	// create base client with basic configs
	baseClient, err := NewBaseClient(credentials, absolutePath, isHTTP)
	if err != nil {
		return nil, err
	}
	// add useragent
	baseClient.UserAgent = userAgent

	if credentials.ProxyURL != "" {
		log.Printf("[DEBUG] Using proxy: %s\n", credentials.ProxyURL)
		proxy, err := url.Parse(credentials.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing proxy url: %s", err)
		}

		// override transport config incase of using proxy
		transCfg := &http.Transport{
			// nolint:gas
			TLSClientConfig: &tls.Config{InsecureSkipVerify: credentials.Insecure}, // ignore expired SSL certificates
		}
		transCfg.Proxy = http.ProxyURL(proxy)
		baseClient.client.Transport = logging.NewTransport("Nutanix", transCfg)
	}

	if credentials.SessionAuth {
		log.Printf("[DEBUG] Using session_auth\n")

		ctx := context.TODO()
		req, err := baseClient.NewRequest(ctx, http.MethodGet, "/users/me", nil)
		if err != nil {
			return baseClient, err
		}

		resp, err := baseClient.client.Do(req)

		if err != nil {
			return baseClient, err
		}
		defer func() {
			if rerr := resp.Body.Close(); err == nil {
				err = rerr
			}
		}()

		err = CheckResponse(resp)

		baseClient.Cookies = resp.Cookies()
	}

	return baseClient, nil
}

// NewBaseClient returns a basic http/https client based on isHttp flag
func NewBaseClient(credentials *Credentials, absolutePath string, isHTTP bool) (*Client, error) {
	if absolutePath == "" {
		return nil, fmt.Errorf("absolutePath argument must be passed")
	}

	httpClient := http.DefaultClient

	transCfg := &http.Transport{
		//to skip/unskip SSL certificate validation
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: credentials.Insecure,
		},
	}
	httpClient.Transport = logging.NewTransport("Nutanix", transCfg)

	protocol := httpsPrefix
	if isHTTP {
		protocol = httpPrefix
	}

	baseURL, err := url.Parse(fmt.Sprintf(defaultBaseURL, protocol, credentials.URL))
	if err != nil {
		return nil, err
	}

	c := &Client{credentials, httpClient, baseURL, "", nil, nil, absolutePath, ""}

	return c, nil
}

// NewRequest creates a request
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	// check if client exists or not
	if c.client == nil {
		return nil, fmt.Errorf(c.ErrorMsg)
	}

	rel, errp := url.Parse(c.AbsolutePath + urlStr)
	if errp != nil {
		return nil, errp
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

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	if c.Cookies != nil {
		for _, i := range c.Cookies {
			req.AddCookie(i)
		}
	} else {
		req.Header.Add("Authorization", "Basic "+
			base64.StdEncoding.EncodeToString([]byte(c.Credentials.Username+":"+c.Credentials.Password)))
	}
	return req, nil
}

// NewRequest creates a request without authorisation headers
func (c *Client) NewUnAuthRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	// check if client exists or not
	if c.client == nil {
		return nil, fmt.Errorf(c.ErrorMsg)
	}

	//create main api url
	rel, err := url.Parse(c.AbsolutePath + urlStr)
	if err != nil {
		return nil, err
	}
	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		er := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, er
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	//add api headers
	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

// NewUnAuthFormEncodedRequest returns content-type: application/x-www-form-urlencoded based unauth request
func (c *Client) NewUnAuthFormEncodedRequest(ctx context.Context, method, urlStr string, body map[string]string) (*http.Request, error) {
	// check if client exists or not
	if c.client == nil {
		return nil, fmt.Errorf(c.ErrorMsg)
	}
	//create main api url
	rel, err := url.Parse(c.AbsolutePath + urlStr)
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

	//add api headers
	req.Header.Add("Content-Type", formEncodedType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

// NewUploadRequest Handles image uploads for image service
func (c *Client) NewUploadRequest(ctx context.Context, method, urlStr string, fileReader *os.File) (*http.Request, error) {
	// check if client exists or not
	if c.client == nil {
		return nil, fmt.Errorf(c.ErrorMsg)
	}
	rel, errp := url.Parse(c.AbsolutePath + urlStr)
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
		base64.StdEncoding.EncodeToString([]byte(c.Credentials.Username+":"+c.Credentials.Password)))

	return req, nil
}

// NewUploadRequest handles image uploads for image service
func (c *Client) NewUnAuthUploadRequest(ctx context.Context, method, urlStr string, fileReader *os.File) (*http.Request, error) {
	// check if client exists or not
	if c.client == nil {
		return nil, fmt.Errorf(c.ErrorMsg)
	}
	rel, errp := url.Parse(c.AbsolutePath + urlStr)
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

// Do performs request passed
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {
	// check if client exists or not
	if c.client == nil {
		return fmt.Errorf(c.ErrorMsg)
	}

	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	err = CheckResponse(resp)

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

func searchSlice(slice []string, key string) bool {
	for _, v := range slice {
		if v == key {
			return true
		}
	}
	return false
}

// DoWithFilters performs request passed and filters entities in json response
func (c *Client) DoWithFilters(ctx context.Context, req *http.Request, v interface{}, filters []*AdditionalFilter, baseSearchPaths []string) error {
	// check if client exists or not
	if c.client == nil {
		return fmt.Errorf(c.ErrorMsg)
	}
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	err = CheckResponse(resp)
	if err != nil {
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

func filter(body io.ReadCloser, filters []*AdditionalFilter, baseSearchPaths []string) (io.ReadCloser, error) {
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
	json.Unmarshal(b, &res)

	// Full search paths
	searchPaths := map[string][]string{}
	filterMap := map[string]*AdditionalFilter{}
	for _, filter := range filters {
		filterMap[filter.Name] = filter
		//Build search paths by appending target search paths to base paths
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
		return fmt.Errorf("invalid Nutanix Credentials")
	}

	buf, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}

	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

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

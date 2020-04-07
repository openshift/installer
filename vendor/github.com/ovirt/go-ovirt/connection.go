//
// Copyright (c) 2017 Joey <majunjiev@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Some codes of this file is from https://github.com/CpuID/ovirt-engine-sdk-go/blob/master/sdk/http/http.go.
// And I made some bug fixes, Thanks to @CpuID

package ovirtsdk

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"
)

// LogFunc represents a flexiable and injectable logger function which fits to most of logger libraries
type LogFunc func(format string, v ...interface{})

// Connection represents an HTTP connection to the engine server.
// It is intended as the entry point for the SDK, and it provides access to the `system` service and, from there,
// to the rest of the services provided by the API.
type Connection struct {
	url      *url.URL
	username string
	password string
	token    string
	insecure bool
	caFile   string
	caCert   []byte
	headers  map[string]string
	// Debug options
	logFunc LogFunc

	kerberos bool
	timeout  time.Duration
	compress bool
	// http client
	client *http.Client
	// SSO attributes
	ssoToken     string
	ssoTokenName string
}

// URL returns the base URL of this connection.
func (c *Connection) URL() string {
	return c.url.String()
}

// Test tests the connectivity with the server using the credentials provided in connection.
// If connectivity works correctly and the credentials are valid, it returns a nil error,
// or it will return an error containing the reason as the message.
// If the authentication fails because the oauth token is no longer valid it will
// try to re-authenticate, to renew the token.
func (c *Connection) Test() error {
	statusCode, err := c.testToken()
	if err != nil || statusCode == http.StatusUnauthorized {
		// failed, then clear state.
		c.ssoToken = ""
	}
	_, err = c.authenticate()
	return err
}

// testToken tries a minimal request, using the existing token. Returns the status
// code and an error.
func (c *Connection) testToken() (int, error) {
	// a simple http OPTIONS request is the lightest method to test auth.
	options, err := http.NewRequest(http.MethodOptions, c.url.String(), nil)
	// add auth token to the request
	options.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.ssoToken))
	if err != nil {
		// shouldn't fail to construct a request, but report anyway.
		return 0, err
	}
	res, err := c.client.Do(options)
	defer res.Body.Close()

	return res.StatusCode, err
}

func (c *Connection) getHref(object Href) (string, bool) {
	return object.Href()
}

// IsLink indicates if the given object is a link.
// An object is a link if it has an `href` attribute.
func (c *Connection) IsLink(object Href) bool {
	_, ok := c.getHref(object)
	return ok
}

// FollowLink follows the `href` attribute of the given object, retrieves the target object and returns it.
func (c *Connection) FollowLink(object Href) (interface{}, error) {
	if !c.IsLink(object) {
		return nil, errors.New("Can't follow link because object don't have any")
	}
	href, ok := c.getHref(object)
	if !ok {
		return nil, errors.New("Can't follow link because the 'href' attribute does't have a value")
	}
	useURL, err := url.Parse(c.URL())
	if err != nil {
		return nil, errors.New("Failed to parse connection url")
	}
	prefix := useURL.Path
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}
	if !strings.HasPrefix(href, prefix) {
		return nil, fmt.Errorf("The URL '%v' isn't compatible with the base URL of the connection", href)
	}
	path := href[len(prefix):]
	service, err := NewSystemService(c, "").Service(path)
	if err != nil {
		return nil, err
	}

	serviceValue := reflect.ValueOf(service)
	// `object` is ptr, so use Elem() to get struct value
	hrefObjectValue := reflect.ValueOf(object).Elem()
	var requestCaller reflect.Value
	// If it's TypeStructSlice (list)
	if hrefObjectValue.FieldByName("slice").IsValid() {
		// Call List() method
		requestCaller = serviceValue.MethodByName("List").Call([]reflect.Value{})[0]
	} else {
		requestCaller = serviceValue.MethodByName("Get").Call([]reflect.Value{})[0]
	}
	callerResponse := requestCaller.MethodByName("Send").Call([]reflect.Value{})[0]
	// Get the method index, which is not the Must version
	methodIndex := 0
	callerResponseType := callerResponse.Type()
	for i := 0; i < callerResponseType.NumMethod(); i++ {
		if strings.HasPrefix(callerResponseType.Method(i).Name, "Must") {
			methodIndex = i
			break
		}
	}
	methodIndex = 1 - methodIndex
	// Retrieve the data
	returnedValues := callerResponse.Method(methodIndex).Call([]reflect.Value{})

	result, ok := returnedValues[0].Interface(), returnedValues[1].Bool()
	if !ok {
		return nil, errors.New("The data retrieved not exists")
	}
	return result, nil
}

// authenticate uses OAuth to do authentication
func (c *Connection) authenticate() (string, error) {
	if c.ssoToken == "" {
		token, err := c.getAccessToken()
		if err != nil {
			return "", err
		}
		c.ssoToken = token
	}
	return c.ssoToken, nil
}

// Close releases the resources used by this connection.
func (c *Connection) Close() error {
	return c.CloseIfRevokeSSOToken(true)
}

// CloseIfRevokeSSOToken releases the resources used by this connection.
// logout parameter specifies if token should be revoked, and so user should be logged out.
func (c *Connection) CloseIfRevokeSSOToken(logout bool) error {
	if logout {
		return c.revokeAccessToken()
	}
	return nil
}

// getAccessToken obtains the access token from SSO to be used for bearer authentication.
func (c *Connection) getAccessToken() (string, error) {
	if c.ssoToken == "" {
		// Build the URL and parameters required for the request:
		url, parameters := c.buildSsoAuthRequest()
		// Send the response and wait for the request:
		ssoResp, err := c.getSsoResponse(url, parameters)
		if err != nil {
			return "", err
		}
		// Top level array already handled in getSsoResponse() generically.
		if ssoResp.SsoError != "" {
			return "", &AuthError{
				baseError{
					Msg: fmt.Sprintf("Error during SSO authentication %s : %s", ssoResp.SsoErrorCode, ssoResp.SsoError),
				},
			}
		}
		c.ssoToken = ssoResp.AccessToken
	}
	return c.ssoToken, nil
}

// Revoke the SSO access token.
func (c *Connection) revokeAccessToken() error {
	// Build the URL and parameters required for the request:
	url, parameters := c.buildSsoRevokeRequest()

	// Send the response and wait for the request:
	ssoResp, err := c.getSsoResponse(url, parameters)
	if err != nil {
		return err
	}

	if ssoResp.SsoError != "" {
		return &AuthError{
			baseError{
				Msg: fmt.Sprintf("Error during SSO token revoke %s : %s", ssoResp.SsoErrorCode, ssoResp.SsoError),
			},
		}
	}

	return nil
}

type ssoResponseJSONParent struct {
	children []ssoResponseJSON
}

type ssoResponseJSON struct {
	AccessToken  string `json:"access_token"`
	SsoError     string `json:"error"`
	SsoErrorCode string `json:"error_code"`
}

// Execute a get request to the SSO server and return the response.
func (c *Connection) getSsoResponse(inputURL *url.URL, parameters map[string]string) (*ssoResponseJSON, error) {
	// Configure TLS parameters:
	var tlsConfig *tls.Config
	if inputURL.Scheme == "https" {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: c.insecure,
		}
		if len(c.caFile) > 0 {
			if _, err := os.Stat(c.caFile); os.IsNotExist(err) {
				return nil, fmt.Errorf("The CA File '%s' doesn't exist", c.caFile)
			}
			caCerts, err := ioutil.ReadFile(c.caFile)
			if err != nil {
				return nil, err
			}
			pool, err := createCertPool(caCerts)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse CA Certificate in file '%s'", c.caFile)
			}
			tlsConfig.RootCAs = pool
		} else if len(c.caCert) > 0 {
			pool, err := createCertPool(c.caCert)
			if err != nil {
				return nil, err
			}
			tlsConfig.RootCAs = pool
		}
	}

	c.client = &http.Client{
		Timeout: c.timeout,
		Transport: &http.Transport{
			// Close the http connection after calling resp.Body.Close()
			DisableKeepAlives:  true,
			DisableCompression: !c.compress,
			TLSClientConfig:    tlsConfig,
		},
	}

	// POST request body:
	formValues := make(url.Values)
	for k1, v1 := range parameters {
		formValues[k1] = []string{v1}
	}
	// Build the net/http request:
	req, err := http.NewRequest("POST", inputURL.String(), strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}

	// Add request headers:
	req.Header.Add("User-Agent", fmt.Sprintf("GoSDK/%s", SDK_VERSION))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	// Send the request and wait for the response:
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse and return the JSON response:
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonObj ssoResponseJSON
	err = json.Unmarshal(body, &jsonObj)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse non-array sso with response %v", string(body))
	}
	// Unmarshal successfully
	if jsonObj.AccessToken != "" || jsonObj.SsoError != "" || jsonObj.SsoErrorCode != "" {
		return &jsonObj, nil
	}
	// Maybe it's array encapsulated, try the other approach.
	var jsonObjList ssoResponseJSONParent
	err = json.Unmarshal(body, &jsonObjList)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse array sso with response %v", string(body))
	}
	if len(jsonObjList.children) > 0 {
		jsonObj.AccessToken = jsonObjList.children[0].AccessToken
		jsonObj.SsoError = jsonObjList.children[0].SsoError
	}

	// Maybe it's revoke access token response, which is empty
	return &jsonObj, nil
}

// buildSsoAuthRequest builds a the URL and parameters to acquire the access token from SSO.
func (c *Connection) buildSsoAuthRequest() (*url.URL, map[string]string) {
	// Compute the entry point and the parameters:
	parameters := map[string]string{
		"scope": "ovirt-app-api",
	}

	var entryPoint string
	if c.kerberos {
		entryPoint = "token-http-auth"
		parameters["grant_type"] = "urn:ovirt:params:oauth:grant-type:http"
	} else {
		entryPoint = "token"
		parameters["grant_type"] = "password"
		parameters["username"] = c.username
		parameters["password"] = c.password
	}

	// Compute the URL:
	var ssoURL url.URL = *c.url
	ssoURL.Path = fmt.Sprintf("/ovirt-engine/sso/oauth/%s", entryPoint)

	// Return the URL and the parameters:
	return &ssoURL, parameters
}

// buildSsoRevokeRequest builds a the URL and parameters to revoke the SSO access token.
// string = the URL of the SSO service
// map = hash containing the parameters required to perform the revoke
func (c *Connection) buildSsoRevokeRequest() (*url.URL, map[string]string) {
	// Compute the parameters:
	parameters := map[string]string{
		"scope": "",
		"token": c.token,
	}

	// Compute the URL:
	var ssoRevokeURL url.URL = *c.url
	ssoRevokeURL.Path = "/ovirt-engine/services/sso-logout"

	// Return the URL and the parameters:
	return &ssoRevokeURL, parameters
}

// SystemService returns a reference to the root of the services tree.
func (c *Connection) SystemService() *SystemService {
	return NewSystemService(c, "")
}

// NewConnectionBuilder creates the `ConnectionBuilder struct instance
func NewConnectionBuilder() *ConnectionBuilder {
	return &ConnectionBuilder{
		conn: &Connection{
			ssoTokenName: "access_token"},
		err: nil}
}

// ConnectionBuilder represents a builder for the `Connection` struct
type ConnectionBuilder struct {
	conn *Connection
	err  error
}

// URL sets the url field for `Connection` instance
func (connBuilder *ConnectionBuilder) URL(urlStr string) *ConnectionBuilder {
	// If already has errors, just return
	if connBuilder.err != nil {
		return connBuilder
	}

	// Save the URL:
	useURL, err := url.Parse(urlStr)
	if err != nil {
		connBuilder.err = err
		return connBuilder
	}
	connBuilder.conn.url = useURL
	return connBuilder
}

// Username sets the username field for `Connection` instance
func (connBuilder *ConnectionBuilder) Username(username string) *ConnectionBuilder {
	// If already has errors, just return
	if connBuilder.err != nil {
		return connBuilder
	}

	connBuilder.conn.username = username
	return connBuilder
}

// Password sets the password field for `Connection` instance
func (connBuilder *ConnectionBuilder) Password(password string) *ConnectionBuilder {
	// If already has errors, just return
	if connBuilder.err != nil {
		return connBuilder
	}

	connBuilder.conn.password = password
	return connBuilder
}

// Insecure sets the insecure field for `Connection` instance
func (connBuilder *ConnectionBuilder) Insecure(insecure bool) *ConnectionBuilder {
	// If already has errors, just return
	if connBuilder.err != nil {
		return connBuilder
	}
	connBuilder.conn.insecure = insecure
	return connBuilder
}

// LogFunc sets the logging function field for `Connection` instance
func (connBuilder *ConnectionBuilder) LogFunc(logFunc LogFunc) *ConnectionBuilder {
	// If already has errors, just return
	if connBuilder.err != nil {
		return connBuilder
	}
	connBuilder.conn.logFunc = logFunc
	return connBuilder
}

// Timeout sets the timeout field for `Connection` instance
func (connBuilder *ConnectionBuilder) Timeout(timeout time.Duration) *ConnectionBuilder {
	// If already has errors, just return
	if connBuilder.err != nil {
		return connBuilder
	}
	connBuilder.conn.timeout = timeout
	return connBuilder
}

// CAFile sets the caFile field for `Connection` instance
func (connBuilder *ConnectionBuilder) CAFile(caFilePath string) *ConnectionBuilder {
	// If already has errors, just return
	if connBuilder.err != nil {
		return connBuilder
	}
	connBuilder.conn.caFile = caFilePath
	return connBuilder
}

// CACert sets the caCert field for `Connection` instance
func (connBuilder *ConnectionBuilder) CACert(caCert []byte) *ConnectionBuilder {
	// If already has errors, just return
	if connBuilder.err != nil {
		return connBuilder
	}
	connBuilder.conn.caCert = caCert
	return connBuilder
}

// Headers sets a map of custom HTTP headers to be added to each HTTP request
func (connBuilder *ConnectionBuilder) Headers(headers map[string]string) *ConnectionBuilder {
	// If already has errors, just return
	if connBuilder.err != nil {
		return connBuilder
	}

	if connBuilder.conn.headers == nil {
		connBuilder.conn.headers = map[string]string{}
	}

	for hk, hv := range headers {
		connBuilder.conn.headers[hk] = hv
	}
	return connBuilder
}

// Kerberos sets the kerberos field for `Connection` instance
func (connBuilder *ConnectionBuilder) Kerberos(kerbros bool) *ConnectionBuilder {
	// If already has errors, just return
	if connBuilder.err != nil {
		return connBuilder
	}
	// TODO: kerbros==true is not implemented
	if kerbros == true {
		connBuilder.err = errors.New("Kerberos is not currently implemented")
		return connBuilder
	}
	connBuilder.conn.kerberos = kerbros
	return connBuilder
}

// Compress sets the compress field for `Connection` instance
func (connBuilder *ConnectionBuilder) Compress(compress bool) *ConnectionBuilder {
	// If already has errors, just return
	if connBuilder.err != nil {
		return connBuilder
	}
	connBuilder.conn.compress = compress
	return connBuilder
}

// Build constructs the `Connection` instance
func (connBuilder *ConnectionBuilder) Build() (*Connection, error) {
	// If already has errors, just return
	if connBuilder.err != nil {
		return nil, connBuilder.err
	}

	// Check parameters
	if connBuilder.conn.url == nil {
		return nil, errors.New("The URL must not be empty")
	}
	if len(connBuilder.conn.username) == 0 {
		return nil, errors.New("The Username must not be empty")
	}
	if len(connBuilder.conn.password) == 0 {
		return nil, errors.New("The Password must not be empty")
	}

	// Construct http.Client
	var tlsConfig *tls.Config
	if connBuilder.conn.url.Scheme == "https" {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: connBuilder.conn.insecure,
		}
		if len(connBuilder.conn.caFile) > 0 {
			// Check if the CA File specified exists.
			if _, err := os.Stat(connBuilder.conn.caFile); os.IsNotExist(err) {
				return nil, fmt.Errorf("The ca file '%s' doesn't exist", connBuilder.conn.caFile)
			}
			caCerts, err := ioutil.ReadFile(connBuilder.conn.caFile)
			if err != nil {
				return nil, err
			}
			pool, err := createCertPool(caCerts)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse CA Certificate in file '%s'", connBuilder.conn.caFile)
			}
			tlsConfig.RootCAs = pool
		} else if len(connBuilder.conn.caCert) > 0 {
			pool, err := createCertPool(connBuilder.conn.caCert)
			if err != nil {
				return nil, err
			}
			tlsConfig.RootCAs = pool
		}
	}
	connBuilder.conn.client = &http.Client{
		Timeout: connBuilder.conn.timeout,
		Transport: &http.Transport{
			// Close the http connection after calling resp.Body.Close()
			DisableKeepAlives:  true,
			DisableCompression: !connBuilder.conn.compress,
			TLSClientConfig:    tlsConfig,
		},
	}
	return connBuilder.conn, nil
}

func createCertPool(caCerts []byte) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caCerts) {
		return nil, fmt.Errorf("Failed to parse CA Certificate")
	}
	return pool, nil
}

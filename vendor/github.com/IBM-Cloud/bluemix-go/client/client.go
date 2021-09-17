//Package client provides a generic client to be used by all services
package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	gohttp "net/http"
	"path"
	"strings"
	"sync"
	"time"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

//TokenProvider ...
type TokenProvider interface {
	RefreshToken() (string, error)
	GetPasscode() (string, error)
	AuthenticatePassword(string, string) error
	AuthenticateAPIKey(string) error
}

/*type PaginatedResourcesHandler interface {
    Resources(rawResponse []byte, curPath string) (resources []interface{}, nextPath string, err error)
}

//HandlePagination ...
type HandlePagination func(c *Client, path string, paginated PaginatedResourcesHandler, cb func(interface{}) bool) (resp *gohttp.Response, err error)
*/

//Client is the base client for all service api client
type Client struct {
	Config         *bluemix.Config
	DefaultHeader  gohttp.Header
	ServiceName    bluemix.ServiceName
	TokenRefresher TokenProvider
	//HandlePagination HandlePagination

	headerLock sync.Mutex
}

//Config stores any generic service client configurations
type Config struct {
	Config   *bluemix.Config
	Endpoint string
}

//New ...
func New(c *bluemix.Config, serviceName bluemix.ServiceName, refresher TokenProvider) *Client {
	return &Client{
		Config:         c,
		ServiceName:    serviceName,
		TokenRefresher: refresher,
		//HandlePagination: pagination,
		DefaultHeader: getDefaultAuthHeaders(serviceName, c),
	}
}

//SendRequest ...
func (c *Client) SendRequest(r *rest.Request, respV interface{}) (*gohttp.Response, error) {

	retries := *c.Config.MaxRetries
	if retries < 1 {
		return c.MakeRequest(r, respV)
	}
	wait := *c.Config.RetryDelay

	return c.tryHTTPRequest(retries, wait, r, respV)
}

// MakeRequest ...
func (c *Client) MakeRequest(r *rest.Request, respV interface{}) (*gohttp.Response, error) {
	httpClient := c.Config.HTTPClient
	if httpClient == nil {
		httpClient = gohttp.DefaultClient
	}
	restClient := &rest.Client{
		DefaultHeader: c.DefaultHeader,
		HTTPClient:    httpClient,
	}
	resp, err := restClient.Do(r, respV, nil)
	// The response returned by go HTTP client.Do() could be nil if request timeout.
	// For convenience, we ensure that response returned by this method is always not nil.
	if resp == nil {
		return new(gohttp.Response), err
	}
	if err != nil {
		if ( resp.StatusCode == 401 || resp.StatusCode == 403) && c.TokenRefresher != nil {
			log.Println("Authentication failed. Trying token refresh")
			c.headerLock.Lock()
			defer c.headerLock.Unlock()
			_, err := c.TokenRefresher.RefreshToken()
			switch err.(type) {
			case nil:
				restClient.DefaultHeader = getDefaultAuthHeaders(c.ServiceName, c.Config)
				for k := range c.DefaultHeader {
					r.Del(k)
				}
				c.DefaultHeader = restClient.DefaultHeader
				resp, err := restClient.Do(r, respV, nil)
				if resp == nil {
					return new(gohttp.Response), err
				}
				if err != nil {
					err = bmxerror.WrapNetworkErrors(resp.Request.URL.Host, err)
				}
				return resp, err
			case *bmxerror.InvalidTokenError:
				return resp, bmxerror.NewRequestFailure("InvalidToken", fmt.Sprintf("%v", err), 401)
			default:
				return resp, fmt.Errorf("Authentication failed, Unable to refresh auth token: %v. Try again later", err)
			}
		}

	}
	return resp, err
}

func (c *Client) tryHTTPRequest(retries int, wait time.Duration, r *rest.Request, respV interface{}) (*gohttp.Response, error) {

	resp, err := c.MakeRequest(r, respV)
	if err != nil {
		if !isRetryable(err) {
			if resp == nil {
				return new(gohttp.Response), err
			}
			return resp, err
		}
		if retries--; retries >= 0 {
			time.Sleep(wait)
			return c.tryHTTPRequest(
				retries, wait, r, respV)
		}
	}
	if resp == nil {
		return new(gohttp.Response), err
	}
	return resp, err
}

//Get ...
func (c *Client) Get(path string, respV interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.GetRequest(c.URL(path))
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.SendRequest(r, respV)
}

//Put ...
func (c *Client) Put(path string, data interface{}, respV interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.PutRequest(c.URL(path)).Body(data)
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.SendRequest(r, respV)
}

//Patch ...
func (c *Client) Patch(path string, data interface{}, respV interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.PatchRequest(c.URL(path)).Body(data)
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.SendRequest(r, respV)
}

//Post ...
func (c *Client) Post(path string, data interface{}, respV interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.PostRequest(c.URL(path)).Body(data)
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}

	return c.SendRequest(r, respV)
}

//PostWithForm ...
func (c *Client) PostWithForm(path string, form interface{}, respV interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.PostRequest(c.URL(path))
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	addToRequestForm(form, r)

	return c.SendRequest(r, respV)
}

//Delete ...
func (c *Client) Delete(path string, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.DeleteRequest(c.URL(path))
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.SendRequest(r, nil)
}

//DeleteWithResp ...
func (c *Client) DeleteWithResp(path string, respV interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.DeleteRequest(c.URL(path))
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.SendRequest(r, respV)
}

//DeleteWithBody ...
func (c *Client) DeleteWithBody(path string, data interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.DeleteRequest(c.URL(path)).Body(data)
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.SendRequest(r, nil)
}

func addToRequestHeader(h interface{}, r *rest.Request) {
	switch v := h.(type) {
	case map[string]string:
		for key, value := range v {
			r.Set(key, value)
		}
	}
}

func addToRequestForm(h interface{}, r *rest.Request) {
	switch v := h.(type) {
	case map[string]string:
		for key, value := range v {
			r.Field(key, value)
		}
	}
}

/*//GetPaginated ...
func (c *Client) GetPaginated(path string, paginated PaginatedResourcesHandler, cb func(interface{}) bool) (resp *gohttp.Response, err error) {
    return c.HandlePagination(c, path, paginated, cb)
}*/

type PaginatedResourcesHandler interface {
	Resources(rawResponse []byte, curPath string) (resources []interface{}, nextPath string, err error)
}

func (c *Client) GetPaginated(path string, paginated PaginatedResourcesHandler, cb func(interface{}) bool) (resp *gohttp.Response, err error) {
	for path != "" {
		var raw json.RawMessage
		resp, err = c.Get(path, &raw)
		if err != nil {
			return
		}

		var resources []interface{}
		var nextPath string
		resources, nextPath, err = paginated.Resources([]byte(raw), path)
		if err != nil {
			err = fmt.Errorf("%s: Error parsing JSON", err.Error())
			return
		}

		for _, resource := range resources {
			if !cb(resource) {
				return
			}
		}

		path = nextPath
	}
	return
}

//URL ...
func (c *Client) URL(path string) string {
	return *c.Config.Endpoint + cleanPath(path)
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return path.Clean(p)
}

const (
	userAgentHeader       = "User-Agent"
	authorizationHeader   = "Authorization"
	uaaAccessTokenHeader  = "X-Auth-Uaa-Token"
	userAccessTokenHeader = "X-Auth-User-Token"
	iamRefreshTokenHeader = "X-Auth-Refresh-Token"
	crRefreshTokenHeader  = "RefreshToken"
)

func getDefaultAuthHeaders(serviceName bluemix.ServiceName, c *bluemix.Config) gohttp.Header {
	h := gohttp.Header{}
	switch serviceName {
	case bluemix.MccpService, bluemix.AccountService:
		h.Set(userAgentHeader, http.UserAgent())
		h.Set(authorizationHeader, c.UAAAccessToken)
	case bluemix.ContainerService:
		h.Set(userAgentHeader, http.UserAgent())
		h.Set(authorizationHeader, c.IAMAccessToken)
		h.Set(iamRefreshTokenHeader, c.IAMRefreshToken)
		h.Set(uaaAccessTokenHeader, c.UAAAccessToken)
	case bluemix.VpcContainerService:
		h.Set(userAgentHeader, http.UserAgent())
		h.Set(authorizationHeader, c.IAMAccessToken)
		h.Set(iamRefreshTokenHeader, c.IAMRefreshToken)
	case bluemix.SchematicsService:
		h.Set(userAgentHeader, http.UserAgent())
		h.Set(authorizationHeader, c.IAMAccessToken)
		h.Set(iamRefreshTokenHeader, c.IAMRefreshToken)
	case bluemix.ContainerRegistryService:
		h.Set(authorizationHeader, c.IAMAccessToken)
		h.Set(crRefreshTokenHeader, c.IAMRefreshToken)
	case bluemix.IAMPAPService, bluemix.AccountServicev1, bluemix.ResourceCatalogrService, bluemix.ResourceControllerService, bluemix.ResourceControllerServicev2, bluemix.ResourceManagementService, bluemix.ResourceManagementServicev2, bluemix.IAMService, bluemix.IAMUUMService, bluemix.IAMUUMServicev2, bluemix.IAMPAPServicev2, bluemix.CseService:
		h.Set(authorizationHeader, c.IAMAccessToken)
	case bluemix.UserManagement:
		h.Set(userAgentHeader, http.UserAgent())
		h.Set(authorizationHeader, c.IAMAccessToken)
	case bluemix.CisService:
		h.Set(userAgentHeader, http.UserAgent())
		h.Set(userAccessTokenHeader, c.IAMAccessToken)
	case bluemix.GlobalSearchService, bluemix.GlobalTaggingService:
		h.Set(userAgentHeader, http.UserAgent())
		h.Set(authorizationHeader, c.IAMAccessToken)
		h.Set(iamRefreshTokenHeader, c.IAMRefreshToken)
	case bluemix.ICDService:
		h.Set(userAgentHeader, http.UserAgent())
		h.Set(authorizationHeader, c.IAMAccessToken)
	case bluemix.CertificateManager:
		h.Set(userAgentHeader, http.UserAgent())
		h.Set(authorizationHeader, c.IAMAccessToken)
	case bluemix.HPCService:
		h.Set(authorizationHeader, c.IAMAccessToken)
	case bluemix.FunctionsService:
		h.Set(userAgentHeader, http.UserAgent())
		h.Set(authorizationHeader, c.IAMAccessToken)

	default:
		log.Println("Unknown service - No auth headers set")
	}
	return h
}

func isTimeout(err error) bool {
	if bmErr, ok := err.(bmxerror.RequestFailure); ok {
		switch bmErr.StatusCode() {
		case 408, 504, 599, 429, 500, 502, 520, 503, 403:
			return true
		}
	}

	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}

	if netErr, ok := err.(*net.OpError); ok && netErr.Timeout() {
		return true
	}

	if netErr, ok := err.(net.UnknownNetworkError); ok && netErr.Timeout() {
		return true
	}

	return false
}

func isRetryable(err error) bool {
	return isTimeout(err)
}

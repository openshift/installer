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

	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
)

const (
	// libraryVersion = "v3"
	defaultBaseURL = "https://%s/"
	// absolutePath   = "api/nutanix/" + libraryVersion
	// userAgent      = "nutanix/" + libraryVersion
	mediaType = "application/json"
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
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response, interface{})

// Credentials needed username and password
type Credentials struct {
	URL         string
	Username    string
	Password    string
	Endpoint    string
	Port        string
	Insecure    bool
	SessionAuth bool
	ProxyURL    string
}

// NewClient returns a new Nutanix API client.
func NewClient(credentials *Credentials, userAgent string, absolutePath string) (*Client, error) {
	if userAgent == "" {
		return nil, fmt.Errorf("userAgent argument must be passed")
	}
	if absolutePath == "" {
		return nil, fmt.Errorf("absolutePath argument must be passed")
	}

	transCfg := &http.Transport{
		// nolint:gas
		TLSClientConfig: &tls.Config{InsecureSkipVerify: credentials.Insecure}, // ignore expired SSL certificates
	}

	if credentials.ProxyURL != "" {
		log.Printf("[DEBUG] Using proxy: %s\n", credentials.ProxyURL)
		proxy, err := url.Parse(credentials.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing proxy url: %s", err)
		}

		transCfg.Proxy = http.ProxyURL(proxy)
	}

	httpClient := http.DefaultClient

	httpClient.Transport = logging.NewTransport("Nutanix", transCfg)

	baseURL, err := url.Parse(fmt.Sprintf(defaultBaseURL, credentials.URL))

	if err != nil {
		return nil, err
	}

	c := &Client{credentials, httpClient, baseURL, userAgent, nil, nil, absolutePath}

	if credentials.SessionAuth {
		log.Printf("[DEBUG] Using session_auth\n")

		ctx := context.TODO()
		req, err := c.NewRequest(ctx, http.MethodGet, "/users/me", nil)
		if err != nil {
			return c, err
		}

		resp, err := c.client.Do(req)

		if err != nil {
			return c, err
		}
		defer func() {
			if rerr := resp.Body.Close(); err == nil {
				err = rerr
			}
		}()

		err = CheckResponse(resp)

		c.Cookies = resp.Cookies()
	}

	return c, nil
}

// NewRequest creates a request
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
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

// NewUploadRequest Handles image uploads for image service
func (c *Client) NewUploadRequest(ctx context.Context, method, urlStr string, body []byte) (*http.Request, error) {
	rel, errp := url.Parse(c.AbsolutePath + urlStr)
	if errp != nil {
		return nil, errp
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := bytes.NewBuffer(body)

	req, err := http.NewRequest(method, u.String(), buf)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Accept", "application/octet-stream")
	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Authorization", "Basic "+
		base64.StdEncoding.EncodeToString([]byte(c.Credentials.Username+":"+c.Credentials.Password)))

	return req, nil
}

// OnRequestCompleted sets the DO API request completion callback
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

// Do performs request passed
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {
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
	log.Print("[DEBUG] after json.Unmarshal")

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

	log.Print("[DEBUG] after bunch of switch cases")
	if err != nil {
		return err
	}
	log.Print("[DEBUG] first nil check")

	// karbon error check
	if messageInfo, ok := res["message_info"]; ok {
		return fmt.Errorf("error: %s", messageInfo)
	}
	if message, ok := res["message"]; ok {
		log.Print(message)
		return fmt.Errorf("error: %s", message)
	}
	if errRes.State != "ERROR" {
		return nil
	}

	log.Print("[DEBUG] after errRes.State")
	pretty, _ := json.MarshalIndent(errRes, "", "  ")
	return fmt.Errorf("error: %s", string(pretty))
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
	Details map[string]interface{} `json:"details,omitempty"`

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

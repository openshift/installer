// Copyright 2023 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dcl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	apihttp "google.golang.org/api/transport/http"
)

// SendRequest applies the credentials in the provided Config to a request with the specified
// verb, url, and body.  It returns the Response from the server if the request returns a
// 2XX success code, or the *googleapi.Error if it returns any other code. The retry is
// optional; if supplied HTTP errors that are deemed temporary will be retried according
// to the policy implemented by the retry.
func SendRequest(ctx context.Context, c *Config, verb, url string, body *bytes.Buffer, retryProvider RetryProvider) (*RetryDetails, error) {
	hdrs := http.Header{}
	for h, v := range c.header {
		for _, s := range v {
			hdrs.Add(h, s)
		}
	}
	hdrs.Set("User-Agent", c.UserAgent())
	hdrs.Set("Content-Type", c.contentType)

	u, err := AddQueryParams(url, c.queryParams)
	if err != nil {
		return nil, err
	}

	hasUserProjectOverride, billingProject := UserProjectOverride(c, u)
	if hasUserProjectOverride {
		hdrs.Set("X-Goog-User-Project", billingProject)
	}

	mtls, err := GetMTLSEndpoint(u)
	if err != nil {
		return nil, err
	}

	options := []option.ClientOption{
		option.WithScopes(Scopes...),
		internaloption.WithDefaultEndpoint(u),
		internaloption.WithDefaultMTLSEndpoint(mtls),
	}
	for _, o := range c.clientOptions {
		options = append(options, o)
	}

	httpClient, endpoint, err := apihttp.NewClient(ctx, options...)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		u = endpoint
	}

	if _, ok := httpClient.Transport.(loggingTransport); !ok {
		// In cases where the config has been created using WithHTTPClient() we want to
		// replace the default transport with our logging transport only once.
		httpClient = &http.Client{
			Transport: loggingTransport{
				underlyingTransport: httpClient.Transport,
				logger:              c.Logger,
			},
			CheckRedirect: httpClient.CheckRedirect,
			Jar:           httpClient.Jar,
			Timeout:       httpClient.Timeout,
		}
	}

	if body == nil {
		// A nil value indicates an empty request body.
		body = &bytes.Buffer{}
	}
	bodyBytes := body.Bytes()
	req, err := http.NewRequestWithContext(ctx, verb, u, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header = hdrs

	var res *http.Response
	if retryProvider == nil {
		res, err = httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		err = googleapi.CheckResponse(res)
		if err != nil {
			// If this is an error, we will not be returning the
			// body, so we should close it.
			googleapi.CloseBody(res)
			return nil, err
		}
		return &RetryDetails{Request: req, Response: res}, nil
	}

	// The start time of request retries is used to determine if an HTTP error is still retryable.
	start := time.Now()
	err = Do(ctx, func(ctx context.Context) (*RetryDetails, error) {
		// Reset req body before http call.
		req.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		res, err = httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		if err := googleapi.CheckResponse(res); err != nil {
			// If this is an error, we will not be returning the
			// body, so we should close it.
			googleapi.CloseBody(res)
			if IsRetryableRequestError(c, err, false, start) {
				return nil, OperationNotDone{Err: err}
			}
			return nil, err
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		return &RetryDetails{Request: req.Clone(ctx), Response: res}, err
	}, retryProvider)
	if err != nil {
		return nil, err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	return &RetryDetails{Request: req, Response: res}, nil
}

// AddQueryParams adds the specified query parameters to the specified url.
func AddQueryParams(rawurl string, params map[string]string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// ParseResponse reads a JSON response into a Go struct
func ParseResponse(resp *http.Response, ptr interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(ptr)
}

// IsRetryableRequestError returns true if an error is determined to be
// a common retryable error based on heuristics about GCP API behaviours.
// The start time is used to determine if errors with custom timeouts should be retried.
func IsRetryableRequestError(c *Config, err error, retryNotFound bool, start time.Time) bool {
	// Return transient errors that should be retried.
	if IsRetryableHTTPError(err, c.codeRetryability, start) || (retryNotFound && IsNotFound(err)) {
		c.Logger.Infof("Error appears retryable: %s", err)
		return true
	}

	if IsNonRetryableHTTPError(err, c.codeRetryability, start) {
		c.Logger.Infof("Error appears not to be retryable: %s", err)
		return false
	}

	// Assume other errors are retryable.
	c.Logger.Warningf("Unexpected HTTP error, assuming retryable: %s", err)
	return true
}

// Nprintf takes in a format string (with format {{key}} instead of %s) and a params map.
// Returns filled string.
func Nprintf(format string, params map[string]interface{}) string {
	pq := strings.Split(format, "?")
	path := pq[0]
	query := ""
	if len(pq) == 2 {
		query = pq[1]
	} else if len(pq) > 2 {
		return "error: too many path separators."
	}
	for key, val := range params {
		r := regexp.MustCompile(`{{\s?` + regexp.QuoteMeta(key) + `\s?}}`)
		path = r.ReplaceAllString(path, fmt.Sprintf("%v", val))
	}
	for key, val := range params {
		r := regexp.MustCompile(`{{\s?` + regexp.QuoteMeta(key) + `\s?}}`)
		query = r.ReplaceAllString(query, url.QueryEscape(fmt.Sprintf("%v", val)))
	}
	if query != "" {
		return path + "?" + query
	}
	return path
}

// URL takes in a partial URL, default base path, optional user-specified base-path and a params map.
func URL(urlpath, basePath, userPath string, params map[string]interface{}) string {
	if userPath != "" {
		if strings.HasSuffix(userPath, "/") {
			userPath = userPath[:len(userPath)-1]
		}
		return Nprintf(strings.Join([]string{userPath, urlpath}, "/"), params)
	}
	if strings.HasSuffix(basePath, "/") {
		basePath = strings.TrimSuffix(basePath, "/")
	}
	return Nprintf(strings.Join([]string{basePath, urlpath}, "/"), params)
}

// ResponseBodyAsJSON reads the response body from a *RetryDetails and returns
// it as unstructured JSON in a map[string]interface{}.
func ResponseBodyAsJSON(retry *RetryDetails) (map[string]interface{}, error) {
	defer retry.Response.Body.Close()
	b, err := ioutil.ReadAll(retry.Response.Body)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetMTLSEndpoint returns the API endpoint used for mTLS authentication.
func GetMTLSEndpoint(baseEndpoint string) (string, error) {
	u, err := url.Parse(baseEndpoint)
	if err != nil {
		return "", err
	}
	portParts := strings.Split(u.Host, ":")
	if len(portParts) == 0 || portParts[0] == "" {
		return "", fmt.Errorf("api endpoint %q is missing host", u.String())
	}
	domainParts := strings.Split(portParts[0], ".")
	if len(domainParts) > 1 {
		u.Host = fmt.Sprintf("%s.mtls.%s", domainParts[0], strings.Join(domainParts[1:], "."))
	} else {
		u.Host = fmt.Sprintf("%s.mtls", domainParts[0])
	}
	if len(portParts) > 1 {
		u.Host = fmt.Sprintf("%s:%s", u.Host, portParts[1])
	}
	return u.String(), nil
}

// UserProjectOverride returns true if user project override should be used and the project that should be set.
func UserProjectOverride(c *Config, url string) (bool, string) {
	if !c.userOverrideProject {
		return false, ""
	}

	if c.billingProject != "" {
		return true, c.billingProject
	}

	r := regexp.MustCompile(`projects/([a-z0-9A-Z-:_]*)/`)
	g := r.FindStringSubmatch(url)
	if g != nil && len(g) > 1 {
		return true, g[1]
	}

	// This URL does not contain a project and no project was found in the URL.
	// This most likely means a non-project resource was used accidentally.
	return false, ""
}

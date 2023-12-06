package core

// (C) Copyright IBM Corp. 2023.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"sync"
	"time"
)

// MCSPAuthenticator uses an apikey to obtain an access token,
// and adds the access token to requests via an Authorization header
// of the form:  "Authorization: Bearer <access-token>"
type MCSPAuthenticator struct {

	// [Required] The apikey used to fetch the bearer token from the token server.
	ApiKey string

	// [Required] The endpoint base URL for the token server.
	URL string

	// [Optional] A flag that indicates whether verification of the token server's SSL certificate
	// should be disabled; defaults to false.
	DisableSSLVerification bool

	// [Optional] A set of key/value pairs that will be sent as HTTP headers in requests
	// made to the token server.
	Headers map[string]string

	// [Optional] The http.Client object used to invoke token server requests.
	// If not specified by the user, a suitable default Client will be constructed.
	Client     *http.Client
	clientInit sync.Once

	// The cached token and expiration time.
	tokenData *mcspTokenData

	// Mutex to make the tokenData field thread safe.
	tokenDataMutex sync.Mutex
}

var mcspRequestTokenMutex sync.Mutex
var mcspNeedsRefreshMutex sync.Mutex

const (
	mcspAuthOperationPath = "/siusermgr/api/1.0/apikeys/token"
)

// MCSPAuthenticatorBuilder is used to construct an MCSPAuthenticator instance.
type MCSPAuthenticatorBuilder struct {
	MCSPAuthenticator
}

// NewMCSPAuthenticatorBuilder returns a new builder struct that
// can be used to construct an MCSPAuthenticator instance.
func NewMCSPAuthenticatorBuilder() *MCSPAuthenticatorBuilder {
	return &MCSPAuthenticatorBuilder{}
}

// SetApiKey sets the ApiKey field in the builder.
func (builder *MCSPAuthenticatorBuilder) SetApiKey(s string) *MCSPAuthenticatorBuilder {
	builder.MCSPAuthenticator.ApiKey = s
	return builder
}

// SetURL sets the URL field in the builder.
func (builder *MCSPAuthenticatorBuilder) SetURL(s string) *MCSPAuthenticatorBuilder {
	builder.MCSPAuthenticator.URL = s
	return builder
}

// SetDisableSSLVerification sets the DisableSSLVerification field in the builder.
func (builder *MCSPAuthenticatorBuilder) SetDisableSSLVerification(b bool) *MCSPAuthenticatorBuilder {
	builder.MCSPAuthenticator.DisableSSLVerification = b
	return builder
}

// SetHeaders sets the Headers field in the builder.
func (builder *MCSPAuthenticatorBuilder) SetHeaders(headers map[string]string) *MCSPAuthenticatorBuilder {
	builder.MCSPAuthenticator.Headers = headers
	return builder
}

// SetClient sets the Client field in the builder.
func (builder *MCSPAuthenticatorBuilder) SetClient(client *http.Client) *MCSPAuthenticatorBuilder {
	builder.MCSPAuthenticator.Client = client
	return builder
}

// Build() returns a validated instance of the MCSPAuthenticator with the config that was set in the builder.
func (builder *MCSPAuthenticatorBuilder) Build() (*MCSPAuthenticator, error) {

	// Make sure the config is valid.
	err := builder.MCSPAuthenticator.Validate()
	if err != nil {
		return nil, err
	}

	return &builder.MCSPAuthenticator, nil
}

// client returns the authenticator's http client after potentially initializing it.
func (authenticator *MCSPAuthenticator) client() *http.Client {
	authenticator.clientInit.Do(func() {
		if authenticator.Client == nil {
			authenticator.Client = DefaultHTTPClient()
			authenticator.Client.Timeout = time.Second * 30

			// If the user told us to disable SSL verification, then do it now.
			if authenticator.DisableSSLVerification {
				transport := &http.Transport{
					// #nosec G402
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}
				authenticator.Client.Transport = transport
			}
		}
	})
	return authenticator.Client
}

// newMCSPAuthenticatorFromMap constructs a new MCSPAuthenticator instance from a map.
func newMCSPAuthenticatorFromMap(properties map[string]string) (authenticator *MCSPAuthenticator, err error) {
	if properties == nil {
		return nil, fmt.Errorf(ERRORMSG_PROPS_MAP_NIL)
	}

	disableSSL, err := strconv.ParseBool(properties[PROPNAME_AUTH_DISABLE_SSL])
	if err != nil {
		disableSSL = false
	}

	authenticator, err = NewMCSPAuthenticatorBuilder().
		SetApiKey(properties[PROPNAME_APIKEY]).
		SetURL(properties[PROPNAME_AUTH_URL]).
		SetDisableSSLVerification(disableSSL).
		Build()

	return
}

// AuthenticationType returns the authentication type for this authenticator.
func (*MCSPAuthenticator) AuthenticationType() string {
	return AUTHTYPE_MCSP
}

// Authenticate adds the Authorization header to the request.
// The value will be of the form: "Authorization: Bearer <bearer-token>""
func (authenticator *MCSPAuthenticator) Authenticate(request *http.Request) error {
	token, err := authenticator.GetToken()
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer "+token)
	return nil
}

// getTokenData returns the tokenData field from the authenticator.
func (authenticator *MCSPAuthenticator) getTokenData() *mcspTokenData {
	authenticator.tokenDataMutex.Lock()
	defer authenticator.tokenDataMutex.Unlock()

	return authenticator.tokenData
}

// setTokenData sets the given mcspTokenData to the tokenData field of the authenticator.
func (authenticator *MCSPAuthenticator) setTokenData(tokenData *mcspTokenData) {
	authenticator.tokenDataMutex.Lock()
	defer authenticator.tokenDataMutex.Unlock()

	authenticator.tokenData = tokenData
	GetLogger().Info("setTokenData: expiration=%d, refreshTime=%d",
		authenticator.tokenData.Expiration, authenticator.tokenData.RefreshTime)
}

// Validate the authenticator's configuration.
//
// Ensures that the ApiKey and URL properties are both specified.
func (authenticator *MCSPAuthenticator) Validate() error {

	if authenticator.ApiKey == "" {
		return fmt.Errorf(ERRORMSG_PROP_MISSING, "ApiKey")
	}

	if authenticator.URL == "" {
		return fmt.Errorf(ERRORMSG_PROP_MISSING, "URL")
	}

	return nil
}

// GetToken: returns an access token to be used in an Authorization header.
// Whenever a new token is needed (when a token doesn't yet exist, needs to be refreshed,
// or the existing token has expired), a new access token is fetched from the token server.
func (authenticator *MCSPAuthenticator) GetToken() (string, error) {
	if authenticator.getTokenData() == nil || !authenticator.getTokenData().isTokenValid() {
		// synchronously request the token
		err := authenticator.synchronizedRequestToken()
		if err != nil {
			return "", err
		}
	} else if authenticator.getTokenData().needsRefresh() {
		// If refresh needed, kick off a go routine in the background to get a new token.
		//nolint: errcheck
		go authenticator.invokeRequestTokenData()
	}

	// return an error if the access token is not valid or was not fetched
	if authenticator.getTokenData() == nil || authenticator.getTokenData().AccessToken == "" {
		return "", fmt.Errorf("Error while trying to get access token")
	}

	return authenticator.getTokenData().AccessToken, nil
}

// synchronizedRequestToken: synchronously checks if the current token in cache
// is valid. If token is not valid or does not exist, it will fetch a new token.
func (authenticator *MCSPAuthenticator) synchronizedRequestToken() error {
	mcspRequestTokenMutex.Lock()
	defer mcspRequestTokenMutex.Unlock()
	// if cached token is still valid, then just continue to use it
	if authenticator.getTokenData() != nil && authenticator.getTokenData().isTokenValid() {
		return nil
	}

	return authenticator.invokeRequestTokenData()
}

// invokeRequestTokenData: requests a new token from the access server and
// unmarshals the token information to the tokenData cache. Returns
// an error if the token was unable to be fetched, otherwise returns nil
func (authenticator *MCSPAuthenticator) invokeRequestTokenData() error {
	tokenResponse, err := authenticator.RequestToken()
	if err != nil {
		return err
	}

	GetLogger().Info("invokeRequestTokenData(): RequestToken returned tokenResponse:\n%+v", *tokenResponse)
	tokenData, err := newMCSPTokenData(tokenResponse)
	if err != nil {
		tokenData = &mcspTokenData{}
	}

	authenticator.setTokenData(tokenData)

	return nil
}

// RequestToken fetches a new access token from the token server.
func (authenticator *MCSPAuthenticator) RequestToken() (*MCSPTokenServerResponse, error) {

	builder := NewRequestBuilder(POST)
	_, err := builder.ResolveRequestURL(authenticator.URL, mcspAuthOperationPath, nil)
	if err != nil {
		return nil, err
	}

	builder.AddHeader(CONTENT_TYPE, APPLICATION_JSON)
	builder.AddHeader(Accept, APPLICATION_JSON)
	requestBody := fmt.Sprintf(`{"apikey":"%s"}`, authenticator.ApiKey)
	_, _ = builder.SetBodyContentString(requestBody)

	// Add user-defined headers to request.
	for headerName, headerValue := range authenticator.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	req, err := builder.Build()
	if err != nil {
		return nil, err
	}

	// If debug is enabled, then dump the request.
	if GetLogger().IsLogLevelEnabled(LevelDebug) {
		buf, dumpErr := httputil.DumpRequestOut(req, req.Body != nil)
		if dumpErr == nil {
			GetLogger().Debug("Request:\n%s\n", RedactSecrets(string(buf)))
		} else {
			GetLogger().Debug(fmt.Sprintf("error while attempting to log outbound request: %s", dumpErr.Error()))
		}
	}

	GetLogger().Debug("Invoking MCSP 'get token' operation: %s", builder.URL)
	resp, err := authenticator.client().Do(req)
	if err != nil {
		return nil, err
	}
	GetLogger().Debug("Returned from MCSP 'get token' operation, received status code %d", resp.StatusCode)

	// If debug is enabled, then dump the response.
	if GetLogger().IsLogLevelEnabled(LevelDebug) {
		buf, dumpErr := httputil.DumpResponse(resp, req.Body != nil)
		if dumpErr == nil {
			GetLogger().Debug("Response:\n%s\n", RedactSecrets(string(buf)))
		} else {
			GetLogger().Debug(fmt.Sprintf("error while attempting to log inbound response: %s", dumpErr.Error()))
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		buff := new(bytes.Buffer)
		_, _ = buff.ReadFrom(resp.Body)

		// Create a DetailedResponse to be included in the error below.
		detailedResponse := &DetailedResponse{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			RawResult:  buff.Bytes(),
		}

		errorMsg := string(detailedResponse.RawResult)
		if errorMsg == "" {
			errorMsg =
				fmt.Sprintf("unexpected status code %d received from MCSP token server %s", detailedResponse.StatusCode, builder.URL)
		}
		return nil, NewAuthenticationError(detailedResponse, fmt.Errorf(errorMsg))
	}

	tokenResponse := &MCSPTokenServerResponse{}
	_ = json.NewDecoder(resp.Body).Decode(tokenResponse)
	defer resp.Body.Close() // #nosec G307

	return tokenResponse, nil
}

// MCSPTokenServerResponse : This struct models a response received from the token server.
type MCSPTokenServerResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
}

// mcspTokenData : This struct represents the cached information related to a fetched access token.
type mcspTokenData struct {
	AccessToken string
	RefreshTime int64
	Expiration  int64
}

// newMCSPTokenData: constructs a new mcspTokenData instance from the specified
// MCSPTokenServerResponse instance.
func newMCSPTokenData(tokenResponse *MCSPTokenServerResponse) (*mcspTokenData, error) {
	if tokenResponse == nil || tokenResponse.Token == "" {
		return nil, fmt.Errorf("Error while trying to parse access token!")
	}

	// Need to crack open the access token (a JWT) to get the expiration and issued-at times
	// so that we can compute the refresh time.
	claims, err := parseJWT(tokenResponse.Token)
	if err != nil {
		return nil, err
	}

	// Compute the adjusted refresh time (expiration time - 20% of timeToLive)
	timeToLive := claims.ExpiresAt - claims.IssuedAt
	expireTime := claims.ExpiresAt
	refreshTime := expireTime - int64(float64(timeToLive)*0.2)

	tokenData := &mcspTokenData{
		AccessToken: tokenResponse.Token,
		Expiration:  expireTime,
		RefreshTime: refreshTime,
	}

	GetLogger().Info("newMCSPTokenData: expiration=%d, refreshTime=%d", tokenData.Expiration, tokenData.RefreshTime)

	return tokenData, nil
}

// isTokenValid: returns true iff the mcspTokenData instance represents a valid (non-expired) access token.
func (tokenData *mcspTokenData) isTokenValid() bool {
	if tokenData.AccessToken != "" && GetCurrentTime() < tokenData.Expiration {
		GetLogger().Info("isTokenValid: Token is valid!")
		return true
	}
	GetLogger().Info("isTokenValid: Token is NOT valid!")
	GetLogger().Info("isTokenValid: expiration=%d, refreshTime=%d", tokenData.Expiration, tokenData.RefreshTime)
	GetLogger().Info("GetCurrentTime(): %d\n", GetCurrentTime())
	return false
}

// needsRefresh: synchronously returns true iff the currently stored access token should be refreshed. This method also
// updates the refresh time if it determines the token needs refreshed to prevent other threads from
// making multiple refresh calls.
func (tokenData *mcspTokenData) needsRefresh() bool {
	mcspNeedsRefreshMutex.Lock()
	defer mcspNeedsRefreshMutex.Unlock()

	// Advance refresh by one minute
	if tokenData.RefreshTime >= 0 && GetCurrentTime() > tokenData.RefreshTime {
		tokenData.RefreshTime = GetCurrentTime() + 60
		return true
	}

	return false
}

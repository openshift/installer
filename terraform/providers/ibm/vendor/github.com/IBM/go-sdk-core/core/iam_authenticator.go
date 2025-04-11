package core

// (C) Copyright IBM Corp. 2019.
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
	"strconv"
	"time"
)

// IamAuthenticator-related constants.
const (
	DEFAULT_IAM_URL             = "https://iam.cloud.ibm.com/identity/token"
	DEFAULT_CONTENT_TYPE        = "application/x-www-form-urlencoded"
	REQUEST_TOKEN_GRANT_TYPE    = "urn:ibm:params:oauth:grant-type:apikey"
	REQUEST_TOKEN_RESPONSE_TYPE = "cloud_iam"
)

// IamAuthenticator : This authenticator will automatically fetch an access token for the
// configured apikey.  Outbound REST requests invoked by the BaseService are then authenticated
// by adding a Bearer-type Authorization header containing the access token.
type IamAuthenticator struct {

	// [Required] The apikey used to fetch the access token from the IAM token server.
	ApiKey string

	// [Optional] The URL representing the IAM token server's endpoint.
	// If not specified, a suitable default value is used.
	URL string

	// [Optional] The ClientId and ClientSecret fields are used to form a "basic auth" Authorization header
	// for interactions with the IAM token server.
	// If neither field is specified, then no Authorization header will be sent with token server requests.
	// These fields are optional, but must be specified together.
	// Default: "", ""
	ClientId     string
	ClientSecret string

	// [Optional] A flag that indicates whether verificaton of the server's SSL certificate should be disabled or not.
	// Default: false
	DisableSSLVerification bool

	// [Optional] A set of key/value pairs that will be sent as HTTP headers in requests
	// made to the token server.
	Headers map[string]string

	// [Optional] The http.Client object used to invoke token server requests.
	// If not specified by the user, a suitable default Client will be constructed.
	Client *http.Client

	// The cached token and expiration time.
	tokenData *iamTokenData
}

// NewIamAuthenticator : Constructs a new IamAuthenticator instance.
func NewIamAuthenticator(apikey string, url string, clientId string, clientSecret string,
	disableSSLVerification bool, headers map[string]string) (*IamAuthenticator, error) {

	authenticator := &IamAuthenticator{
		ApiKey:                 apikey,
		URL:                    url,
		ClientId:               clientId,
		ClientSecret:           clientSecret,
		DisableSSLVerification: disableSSLVerification,
		Headers:                headers,
	}

	// Make sure the config is valid.
	err := authenticator.Validate()
	if err != nil {
		return nil, err
	}

	return authenticator, nil
}

// NewIamAuthenticatorFromMap : Constructs a new IamAuthenticator instance from a map.
func newIamAuthenticatorFromMap(properties map[string]string) (*IamAuthenticator, error) {
	if properties == nil {
		return nil, fmt.Errorf(ERRORMSG_PROPS_MAP_NIL)
	}

	disableSSL, err := strconv.ParseBool(properties[PROPNAME_AUTH_DISABLE_SSL])
	if err != nil {
		disableSSL = false
	}
	return NewIamAuthenticator(properties[PROPNAME_APIKEY], properties[PROPNAME_AUTH_URL],
		properties[PROPNAME_CLIENT_ID], properties[PROPNAME_CLIENT_SECRET],
		disableSSL, nil)
}

func (IamAuthenticator) AuthenticationType() string {
	return AUTHTYPE_IAM
}

// Authenticate: Performs the authentication on the specified Request by adding a Bearer-type Authorization header
// containing the access token fetched from the token server.
func (authenticator IamAuthenticator) Authenticate(request *http.Request) error {
	token, err := authenticator.getToken()
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", fmt.Sprintf(`Bearer %s`, token))
	return nil
}

// Validate: validates the configuration of the IamAuthenticator instance.
func (this IamAuthenticator) Validate() error {
	if this.ApiKey == "" {
		return fmt.Errorf(ERRORMSG_PROP_MISSING, "ApiKey")
	}

	if HasBadFirstOrLastChar(this.ApiKey) {
		return fmt.Errorf(ERRORMSG_PROP_INVALID, "ApiKey")
	}

	// Validate ClientId and ClientSecret.  They must both be specified togther or neither should be specified.
	if this.ClientId == "" && this.ClientSecret == "" {
		// Do nothing as this is the valid scenario
	} else {
		// Since it is NOT the case that both properties are empty, make sure BOTH are specified.
		if this.ClientId == "" {
			return fmt.Errorf(ERRORMSG_PROP_MISSING, "ClientId")
		}

		if this.ClientSecret == "" {
			return fmt.Errorf(ERRORMSG_PROP_MISSING, "ClientSecret")
		}
	}

	return nil
}

// getToken: returns an access token to be used in an Authorization header.
// Whenever a new token is needed (when a token doesn't yet exist, or the existing token has expired),
// a new access token is fetched from the token server.
func (authenticator *IamAuthenticator) getToken() (string, error) {
	if authenticator.tokenData == nil || !authenticator.tokenData.isTokenValid() {
		tokenResponse, err := authenticator.requestToken()
		if err != nil {
			return "", err
		}

		authenticator.tokenData, err = newIamTokenData(tokenResponse)
		if err != nil {
			return "", err
		}
	}

	return authenticator.tokenData.AccessToken, nil
}

// requestToken: fetches a new access token from the token server.
func (authenticator *IamAuthenticator) requestToken() (*iamTokenServerResponse, error) {
	// Use the default IAM URL if one was not specified by the user.
	url := authenticator.URL
	if url == "" {
		url = DEFAULT_IAM_URL
	}

	builder := NewRequestBuilder(POST)
	_, err := builder.ConstructHTTPURL(url, nil, nil)
	if err != nil {
		return nil, err
	}

	builder.AddHeader(CONTENT_TYPE, DEFAULT_CONTENT_TYPE).
		AddHeader(Accept, APPLICATION_JSON).
		AddFormData("grant_type", "", "", REQUEST_TOKEN_GRANT_TYPE).
		AddFormData("apikey", "", "", authenticator.ApiKey).
		AddFormData("response_type", "", "", REQUEST_TOKEN_RESPONSE_TYPE)

	// Add user-defined headers to request.
	for headerName, headerValue := range authenticator.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	req, err := builder.Build()
	if err != nil {
		return nil, err
	}

	// If client id and secret were configured by the user, then set them on the request
	// as a basic auth header.
	if authenticator.ClientId != "" && authenticator.ClientSecret != "" {
		req.SetBasicAuth(authenticator.ClientId, authenticator.ClientSecret)
	}

	// If the authenticator does not have a Client, create one now.
	if authenticator.Client == nil {
		authenticator.Client = &http.Client{
			Timeout: time.Second * 30,
		}

		// If the user told us to disable SSL verification, then do it now.
		if authenticator.DisableSSLVerification {
			transport := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			authenticator.Client.Transport = transport
		}
	}

	resp, err := authenticator.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp != nil {
			buff := new(bytes.Buffer)
			_, _ = buff.ReadFrom(resp.Body)
			return nil, fmt.Errorf(buff.String())
		}
	}

	tokenResponse := &iamTokenServerResponse{}
	_ = json.NewDecoder(resp.Body).Decode(tokenResponse)
	defer resp.Body.Close()
	return tokenResponse, nil
}

// iamTokenServerResponse : This struct models a response received from the token server.
type iamTokenServerResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Expiration   int64  `json:"expiration"`
}

// iamTokenData : This struct represents the cached information related to a fetched access token.
type iamTokenData struct {
	AccessToken string
	RefreshTime int64
}

// newIamTokenData: constructs a new IamTokenData instance from the specified IamTokenServerResponse instance.
func newIamTokenData(tokenResponse *iamTokenServerResponse) (*iamTokenData, error) {
	// Compute the adjusted refresh time (expiration time - 20% of timeToLive)
	timeToLive := tokenResponse.ExpiresIn
	expireTime := tokenResponse.Expiration
	refreshTime := expireTime - int64(float64(timeToLive)*0.2)

	tokenData := &iamTokenData{
		AccessToken: tokenResponse.AccessToken,
		RefreshTime: refreshTime,
	}

	return tokenData, nil
}

// isTokenValid: returns true iff the IamTokenData instance represents a valid (non-expired) access token.
func (this *iamTokenData) isTokenValid() bool {
	if this.AccessToken != "" && GetCurrentTime() < this.RefreshTime {
		return true
	}
	return false
}

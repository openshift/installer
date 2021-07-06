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
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Constants for CP4D
const (
	PRE_AUTH_PATH = "/v1/preauth/validateAuth"
)

// CloudPakForDataAuthenticator : This authenticator will automatically fetch an access token for the
// user-specified username and password.  Outbound REST requests invoked by the BaseService are then authenticated
// by adding a Bearer-type Authorization header containing the access token.
type CloudPakForDataAuthenticator struct {
	// [Required] The URL representing the token server's endpoing.
	URL string

	// [Required] The username and password used to compute the basic auth Authorization header
	// to be sent with requests to the token server.
	Username string
	Password string

	// [Optional] A flag that indicates whether SSL hostname verification should be disabled or not.
	// Default: false
	DisableSSLVerification bool

	// [Optional] A set of key/value pairs that will be sent as HTTP headers in requests
	// made to the token server.
	Headers map[string]string

	// [Optional] The http.Client object used to invoke token server requests.
	// If not specified by the user, a suitable default Client will be constructed.
	Client *http.Client

	// The cached token and expiration time.
	tokenData *cp4dTokenData
}

// NewCloudPakForDataAuthenticator : Constructs a new CloudPakForDataAuthenticator instance.
func NewCloudPakForDataAuthenticator(url string, username string, password string,
	disableSSLVerification bool, headers map[string]string) (*CloudPakForDataAuthenticator, error) {

	authenticator := &CloudPakForDataAuthenticator{
		Username:               username,
		Password:               password,
		URL:                    url,
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

// newCloudPakForDataAuthenticatorFromMap : Constructs a new CloudPakForDataAuthenticator instance from a map.
func newCloudPakForDataAuthenticatorFromMap(properties map[string]string) (*CloudPakForDataAuthenticator, error) {
	if properties == nil {
		return nil, fmt.Errorf(ERRORMSG_PROPS_MAP_NIL)
	}

	disableSSL, err := strconv.ParseBool(properties[PROPNAME_AUTH_DISABLE_SSL])
	if err != nil {
		disableSSL = false
	}
	return NewCloudPakForDataAuthenticator(properties[PROPNAME_AUTH_URL],
		properties[PROPNAME_USERNAME], properties[PROPNAME_PASSWORD],
		disableSSL, nil)
}

func (CloudPakForDataAuthenticator) AuthenticationType() string {
	return AUTHTYPE_CP4D
}

// Validate: validates the configuration.
func (authenticator CloudPakForDataAuthenticator) Validate() error {

	if authenticator.Username == "" {
		return fmt.Errorf(ERRORMSG_PROP_MISSING, "Username")
	}

	if authenticator.Password == "" {
		return fmt.Errorf(ERRORMSG_PROP_MISSING, "Password")
	}

	if authenticator.URL == "" {
		return fmt.Errorf(ERRORMSG_PROP_MISSING, "URL")
	}

	return nil
}

// Authenticate: performs the authentication on the specified Request by adding a Bearer-type Authorization header
// containing the access token fetched from the token server.
func (authenticator CloudPakForDataAuthenticator) Authenticate(request *http.Request) error {
	token, err := authenticator.getToken()
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", fmt.Sprintf(`Bearer %s`, token))
	return nil
}

// getToken: returns an access token to be used in an Authorization header.
// Whenever a new token is needed (when a token doesn't yet exist, or the existing token has expired),
// a new access token is fetched from the token server.
func (authenticator *CloudPakForDataAuthenticator) getToken() (string, error) {
	if authenticator.tokenData == nil || !authenticator.tokenData.isTokenValid() {
		tokenResponse, err := authenticator.requestToken()
		if err != nil {
			return "", err
		}

		authenticator.tokenData, err = newCp4dTokenData(tokenResponse)
		if err != nil {
			return "", err
		}
	}

	return authenticator.tokenData.AccessToken, nil
}

// requestToken: fetches a new access token from the token server.
func (authenticator *CloudPakForDataAuthenticator) requestToken() (*cp4dTokenServerResponse, error) {
	// If the user-specified URL does not end with the required path,
	// then add it now.
	url := authenticator.URL
	if !strings.HasSuffix(url, PRE_AUTH_PATH) {
		url = fmt.Sprintf("%s%s", url, PRE_AUTH_PATH)
	}

	builder, err := NewRequestBuilder(GET).ConstructHTTPURL(url, nil, nil)
	if err != nil {
		return nil, err
	}

	// Add user-defined headers to request.
	for headerName, headerValue := range authenticator.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	req, err := builder.Build()
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(authenticator.Username, authenticator.Password)

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

	tokenResponse := &cp4dTokenServerResponse{}
	_ = json.NewDecoder(resp.Body).Decode(tokenResponse)
	defer resp.Body.Close()
	return tokenResponse, nil
}

// cp4dTokenServerResponse : This struct models a response received from the token server.
type cp4dTokenServerResponse struct {
	Username    string   `json:"username,omitempty"`
	Role        string   `json:"role,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Subject     string   `json:"sub,omitempty"`
	Issuer      string   `json:"iss,omitempty"`
	Audience    string   `json:"aud,omitempty"`
	UID         string   `json:"uid,omitempty"`
	AccessToken string   `json:"accessToken,omitempty"`
	MessageCode string   `json:"_messageCode_,omitempty"`
	Message     string   `json:"message,omitempty"`
}

// cp4dTokenData : This struct represents the cached information related to a fetched access token.
type cp4dTokenData struct {
	AccessToken string
	RefreshTime int64
}

// newCp4dTokenData: constructs a new Cp4dTokenData instance from the specified Cp4dTokenServerResponse instance.
func newCp4dTokenData(tokenResponse *cp4dTokenServerResponse) (*cp4dTokenData, error) {
	// Need to crack open the access token (a JWToken) to get the expiration and issued-at times.
	claims := &jwt.StandardClaims{}
	if token, _ := jwt.ParseWithClaims(tokenResponse.AccessToken, claims, nil); token == nil {
		return nil, fmt.Errorf("Error while trying to parse access token!")
	}

	// Compute the adjusted refresh time (expiration time - 20% of timeToLive)
	timeToLive := claims.ExpiresAt - claims.IssuedAt
	expireTime := claims.ExpiresAt
	refreshTime := expireTime - int64(float64(timeToLive)*0.2)

	tokenData := &cp4dTokenData{
		AccessToken: tokenResponse.AccessToken,
		RefreshTime: refreshTime,
	}
	return tokenData, nil
}

// isTokenValid: returns true iff the Cp4dTokenData instance represents a valid (non-expired) access token.
func (this *cp4dTokenData) isTokenValid() bool {
	if this.AccessToken != "" && GetCurrentTime() < this.RefreshTime {
		return true
	}
	return false
}

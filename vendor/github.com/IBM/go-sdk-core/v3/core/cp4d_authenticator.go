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
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Constants for CP4D
const (
	PRE_AUTH_PATH = "/v1/preauth/validateAuth"
)

// CloudPakForDataAuthenticator uses a username and password pair to obtain a
// suitable bearer token, and adds the bearer token to requests via an
// Authorization header of the form:
//
// 		Authorization: Bearer <bearer-token>
//
type CloudPakForDataAuthenticator struct {
	// The URL representing the Cloud Pak for Data token service endpoint [required].
	URL string

	// The username used to obtain a bearer token [required].
	Username string

	// The password used to obtain a bearer token [required].
	Password string

	// A flag that indicates whether verification of the server's SSL certificate
	// should be disabled; defaults to false [optional].
	DisableSSLVerification bool

	// Default headers to be sent with every CP4D token request [optional].
	Headers map[string]string

	// The http.Client object used to invoke token server requests [optional]. If
	// not specified, a suitable default Client will be constructed.
	Client *http.Client

	// The cached token and expiration time.
	tokenData *cp4dTokenData
}

var cp4dRequestTokenMutex sync.Mutex
var cp4dNeedsRefreshMutex sync.Mutex

// NewCloudPakForDataAuthenticator constructs a new CloudPakForDataAuthenticator
// instance.
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

// AuthenticationType returns the authentication type for this authenticator.
func (CloudPakForDataAuthenticator) AuthenticationType() string {
	return AUTHTYPE_CP4D
}

// Validate the authenticator's configuration.
//
// Ensures the username, password, and url are not Nil. Additionally, ensures
// they do not contain invalid characters.
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

// Authenticate adds the bearer token (obtained from the token server) to the
// specified request.
//
// The CP4D bearer token will be added to the request's headers in the form:
//
// 		Authorization: Bearer <bearer-token>
//
func (authenticator *CloudPakForDataAuthenticator) Authenticate(request *http.Request) error {
	token, err := authenticator.getToken()
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", fmt.Sprintf(`Bearer %s`, token))
	return nil
}

// getToken: returns an access token to be used in an Authorization header.
// Whenever a new token is needed (when a token doesn't yet exist, needs to be refreshed,
// or the existing token has expired), a new access token is fetched from the token server.
func (authenticator *CloudPakForDataAuthenticator) getToken() (string, error) {
	if authenticator.tokenData == nil || !authenticator.tokenData.isTokenValid() {
		// synchronously request the token
		err := authenticator.synchronizedRequestToken()
		if err != nil {
			return "", err
		}
	} else if authenticator.tokenData.needsRefresh() {
		// If refresh needed, kick off a go routine in the background to get a new token
		ch := make(chan error)
		go func() {
			ch <- authenticator.getTokenData()
		}()
		select {
		case err := <-ch:
			if err != nil {
				return "", err
			}
		default:
		}
	}

	// return an error if the access token is not valid or was not fetched
	if authenticator.tokenData == nil || authenticator.tokenData.AccessToken == "" {
		return "", fmt.Errorf("Error while trying to get access token")
	}

	return authenticator.tokenData.AccessToken, nil
}

// synchronizedRequestToken: synchronously checks if the current token in cache
// is valid. If token is not valid or does not exist, it will fetch a new token
// and set the tokenRefreshTime
func (authenticator *CloudPakForDataAuthenticator) synchronizedRequestToken() error {
	cp4dRequestTokenMutex.Lock()
	defer cp4dRequestTokenMutex.Unlock()
	// if cached token is still valid, then just continue to use it
	if authenticator.tokenData != nil && authenticator.tokenData.isTokenValid() {
		return nil
	}

	return authenticator.getTokenData()
}

// getTokenData: requests a new token from the access server and
// unmarshals the token information to the tokenData cache. Returns
// an error if the token was unable to be fetched, otherwise returns nil
func (authenticator *CloudPakForDataAuthenticator) getTokenData() error {
	tokenResponse, err := authenticator.requestToken()
	if err != nil {
		return err
	}

	authenticator.tokenData, err = newCp4dTokenData(tokenResponse)
	if err != nil {
		return err
	}

	return nil
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
	Expiration  int64
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
		Expiration:  expireTime,
		RefreshTime: refreshTime,
	}
	return tokenData, nil
}

// isTokenValid: returns true iff the Cp4dTokenData instance represents a valid (non-expired) access token.
func (this *cp4dTokenData) isTokenValid() bool {
	if this.AccessToken != "" && GetCurrentTime() < this.Expiration {
		return true
	}
	return false
}

// needsRefresh: synchronously returns true iff the currently stored access token should be refreshed. This method also
// updates the refresh time if it determines the token needs refreshed to prevent other threads from
// making multiple refresh calls.
func (this *cp4dTokenData) needsRefresh() bool {
	cp4dNeedsRefreshMutex.Lock()
	defer cp4dNeedsRefreshMutex.Unlock()

	// Advance refresh by one minute
	if this.RefreshTime >= 0 && GetCurrentTime() > this.RefreshTime {
		this.RefreshTime += 60
		return true
	}

	return false

}
